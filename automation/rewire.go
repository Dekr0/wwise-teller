package automation

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dekr0/wwise-teller/db"
	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wwise"
)

const (
	FilesToSoundsTypeGrain FilesToSoundsType = 0
	FilesToSoundsTypeGroup FilesToSoundsType = 1
)

type FilesToSoundsType uint8;

type FilesToSoundsHeader struct {
	Workspace  string
	Conversion string
	Type       FilesToSoundsType
	Format     waapi.ConversionFormatType
}

func RewireWithNewSources(
	ctx         context.Context,
	bnk        *wwise.Bank,
	mappingFile string,
	dry         bool,
) error {
	if bnk.DIDX() == nil {
		return wwise.NoDIDX
	}
	if bnk.DATA() == nil {
		return wwise.NoDATA
	}
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	proj, err := waapi.GetProject()
	if err != nil {
		return err
	}
	err = db.CheckDatabaseEnv()
	if err != nil {
		return err
	}

	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	iniDir := filepath.Base(mappingFile)
	header := FilesToSoundsHeader{
		Workspace: iniDir,
		Type: 0,
	}
	row := uint16(0)
	if err = ParseFilesToSoundsHeader(&header, reader, &row); err != nil {
		return err
	}

	wavsMapSound := make(map[string][]*wwise.Sound, 32)
	if header.Type == FilesToSoundsTypeGrain {
		err = ParseSoundMapping(reader, &header, h, wavsMapSound, &row)
		if err != nil {
			return err
		}
		if len(wavsMapSound) <= 0 {
			return nil
		}
	}

	wemsMapSounds := make(map[string][]*wwise.Sound, len(wavsMapSound))
	wsource, err := waapi.CreateConversionList(ctx, wavsMapSound, wemsMapSounds, header.Conversion, dry)
	if len(wavsMapSound) != len(wemsMapSounds) {
		panic("# of wavs file in indexing map does not equal # of output wem in indexing map")
	}
	if err != nil {
		return err
	}

	if dry {
		for wem, sounds := range wemsMapSounds {
			fmt.Println(wem)
			for _, s := range sounds {
				fmt.Println(s.Id)
			}
			fmt.Println()
		}
		return nil
	}

	if err := waapi.WwiseConversion(ctx, wsource, proj); err != nil {
		return err
	}

	wemsMapAudioData := make(map[string][]byte, len(wemsMapSounds))
	errorWems := make([]string, 0, len(wemsMapSounds))
	for wem := range wemsMapSounds {
		if _, in := wemsMapAudioData[wem]; in {
			panic(fmt.Sprintf("Detect duplicated wem file %s when storing audio data", wem))
		}
		wemsMapAudioData[wem], err = os.ReadFile(wem)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read audio data from %s", wem))
			slog.Error(fmt.Sprintf("Discard all rewiring related to file %s", wem))
			errorWems = append(errorWems, wem)
		}
	}
	for _, errorWem := range errorWems {
		delete(wemsMapSounds, errorWem)
	}

	db.WriteLock.Lock()
	defer db.WriteLock.Unlock()
	q, closeDb, commit, txRollback, err := db.CreateDefaultConnWithTxQuery(ctx)
	if err != nil {
		return err
	}
	defer closeDb()

	wemsMapMediaIndexs := make(map[string]wwise.MediaIndex, len(wemsMapSounds))
	for wem := range wemsMapSounds {
		sid, err := db.TrySid(ctx, q)
		if err != nil {
			txRollback()
			return fmt.Errorf("Failed to allocate a new source ID: %w", err)
		}
		if _, in := wemsMapMediaIndexs[wem]; in {
			panic(fmt.Sprintf("Detect duplicated wem file %s when storing media index", wem))
		}
		if audioData, in := wemsMapAudioData[wem]; !in { 
			panic(fmt.Sprintf("Cannot find audio data with wem file %s", wem)) 
		} else {
			wemsMapMediaIndexs[wem] = wwise.MediaIndex{Sid: sid, Size: uint32(len(audioData))}
		}
	}

	if err := commit(); err != nil {
		slog.Error(err.Error())
		txRollback()
		return err
	}

	for wem, m := range wemsMapMediaIndexs {
		if audioData, in := wemsMapAudioData[wem]; !in {
			panic(fmt.Sprintf("Cannot find audio data with wem file %s", wem)) 
		} else {
			if err := bnk.AppendAudio(audioData, m.Sid); err != nil {
				// Database is out of sync
				panic(fmt.Sprintf("Source ID %d collision after allocating source ID and commit it into the database", m.Sid))
			}
		}
	}

	for wem, sounds := range wemsMapSounds {
		m, in := wemsMapMediaIndexs[wem]
		if !in {
			panic(fmt.Sprintf("Cannot find audio data with wem file %s", wem)) 
		}
		for _, sound := range sounds {
			switch header.Format {
			case waapi.ConversionFormatTypePCM:
				sound.BankSourceData.PluginID = wwise.PCM
			case waapi.ConversionFormatTypeADPCM:
				sound.BankSourceData.PluginID = wwise.ADPCM
			case waapi.ConversionFormatTypeVORBIS:
				sound.BankSourceData.PluginID = wwise.VORBIS
			case waapi.ConversionFormatTypeWEMOpus:
				sound.BankSourceData.PluginID = wwise.WEM_OPUS
			default:
				panic(fmt.Sprintf("Unsupported conversion format %d", header.Format))
			}
			sound.BankSourceData.SourceID = m.Sid
			sound.BankSourceData.InMemoryMediaSize = m.Size
		}
	}

	return nil
}

// Example
// conversion,VORBIS Quality High
// format,2
// project,absolute_path
// type,0
// workspace,workspace_relative_path_or_absolute_path
func ParseFilesToSoundsHeader(
	header *FilesToSoundsHeader,
	reader *csv.Reader,
	row *uint16,
) error {
	var err error
	header.Conversion, err = CheckConversionRow(reader, *row)
	if err != nil {
		return err
	}
	*row += 1

	header.Format, err = CheckFormatRow(reader, *row)
	if err != nil {
		return err
	}
	*row += 1

	header.Type, err = CheckRewireTypeRow(reader, *row)
	if err != nil {
		return err
	}
	*row += 1

	header.Workspace, err = CheckWorkspaceRow(reader, header.Workspace, *row)
	if err != nil {
		return err
	}
	*row += 1

	return nil
}

// Assumption
// Skip if a row has the following error
// - less than 2 columns
// - an input is not in wave format
// - an input does not exist
// - not an unsigned integer for # of sound IDs specified
// - provided # of sound IDs is less than # of sound IDs specified
// - not an unsigned integer for a sound ID
// It will append .wav extension if no extension is provided
// It will use workspace to construct full path if a relative path is provided
// Skip a sound ID if it doesn't exist
// Merge duplicate input
func ParseSoundMapping(
	reader *csv.Reader,
	header *FilesToSoundsHeader,
	h *wwise.HIRC,
	wavsMapping map[string][]*wwise.Sound,
	rowNum *uint16,
) error {
	var err   error
	var row []string
	var input, ext string
	for {
		row, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Syntax
		if len(row) < 2 {
			slog.Error(fmt.Sprintf("Expecting two columns, filename and # of targeting sound IDs, at row %d", rowNum))
			*rowNum += 1
			continue
		}

		// Extension
		input = row[0]
		ext = filepath.Ext(input)
		if ext == "" {
			ext = ".wav"
			input += ext
		}
		if ext != ".wav" {
			slog.Error("Wave file is the only supported file format.")
			*rowNum += 1
			continue
		}

		// Full path
		if !filepath.IsAbs(input) {
			input = filepath.Join(header.Workspace, input)
		}

		// Existence
		_, err = os.Lstat(input)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to obtain info of wave file %s", input), "error", err)
			*rowNum += 1
			continue
		}

		// Counter
		count, err := strconv.ParseUint(row[1], 10, 8)
		if err != nil {
			slog.Error("Failed to parse # of targeting sound IDs", "error", err)
			continue
		}
		if count > uint64(len(row) - 2) {
			slog.Error(fmt.Sprintf("Expecting %d of targeting sound IDs but receiving %d", count, len(row)- 2))
			continue
		}

		if sounds, in := wavsMapping[input]; !in {
			// Parse IDs
			columnNum := 3
			targets := make([]*wwise.Sound, 0, count)
			for _, i := range row[2:] {
				id, err := strconv.ParseUint(i, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse sound ID at column %d", columnNum), "error", err)
					columnNum += 1
					continue
				}
				if v, ok := h.ActorMixerHirc.Load(uint32(id)); !ok {
					slog.Error(fmt.Sprintf("ID %d does not have an associated sound object", id))
				} else {
					// no duplication check because it will map to the same sid 
					// at the end
					targets = append(targets, v.(*wwise.Sound))
				}
				columnNum += 1
			}
			wavsMapping[input] = targets
		} else {
			slog.Warn(fmt.Sprintf("Duplicated input file %s at row %d, column 1", input, rowNum))
			// Parse IDs
			columnNum := 3
			for _, i := range row[2:] {
				id, err := strconv.ParseUint(i, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse sound ID at column %d", columnNum), "error", err)
					columnNum += 1
					continue
				}
				if v, ok := h.ActorMixerHirc.Load(uint32(id)); !ok {
					slog.Error(fmt.Sprintf("ID %d does not have an associated sound object", id))
				} else {
					// no duplication check because it will map to the same sid 
					// at the end
					sounds = append(sounds, v.(*wwise.Sound))
				}
				columnNum += 1
			}
			wavsMapping[input] = sounds
		}
		*rowNum += 1
	}

	return nil
}

func CheckConversionRow(reader *csv.Reader, row uint16) (string, error) {
	conversionRow, err := reader.Read()
	if err != nil {
		return "", fmt.Errorf("Failed to obtain conversion setting: %w", err)
	}
	if !strings.EqualFold(conversionRow[0], "conversion") {
		return "", fmt.Errorf("Expecting Row %d, Column 1 to be `conversion`", row)
	}
	if len(conversionRow) < 2 {
		return "", fmt.Errorf("Missing conversion setting value")
	}
	return conversionRow[1], nil
}

func CheckFormatRow(reader *csv.Reader, row uint16) (waapi.ConversionFormatType, error) {
	formatRow, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain format setting: %w", err)
	}
	if !strings.EqualFold(formatRow[0], "format") {
		return 0, fmt.Errorf("Expecting Row %d, Column 1 to be `format`", row)
	}
	if len(formatRow) < 2 {
		return 0, fmt.Errorf("Missing format setting value")
	}
	format, err := strconv.Atoi(formatRow[1])
	if err != nil {
		return 0, fmt.Errorf("Failed to parse format setting value: %w", err)
	}
	if format < int(waapi.ConversionFormatTypePCM) || 
	   format > int(waapi.ConversionFormatTypeWEMOpus) {
		return 0, fmt.Errorf("Invalid format setting value %d", format)
	}
	return waapi.ConversionFormatType(format), nil
}

func CheckRewireTypeRow(reader *csv.Reader, row uint16) (FilesToSoundsType, error) {
	rewireTypeRow, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("Failed to obtain rewire type value: %w", err)
	}
	if !strings.EqualFold(rewireTypeRow[0], "type") {
		return 0, fmt.Errorf("Expecting Row %d, Column 1 to be `type`", row)
	}
	if len(rewireTypeRow) < 2 {
		return 0, fmt.Errorf("Missing rewire type")
	}
	rewireType, err := strconv.Atoi(rewireTypeRow[1])
	if err != nil {
		return 0, fmt.Errorf("Failed to parse rewire type: %w", err)
	}
	if rewireType < int(FilesToSoundsTypeGrain) || rewireType > int(FilesToSoundsTypeGroup) {
		return 0, fmt.Errorf("Invalid rewire type %d", rewireType)
	}
	return FilesToSoundsType(rewireType), nil
}

func CheckWorkspaceRow(reader *csv.Reader, init string, row uint16) (string, error) {
	workspaceRow, err := reader.Read()
	if err != nil {
		return "", fmt.Errorf("Failed to obtain workspace directory: %w", err)
	}
	if !strings.EqualFold(workspaceRow[0], "workspace") {
		return "", fmt.Errorf("Expecting Row %d, Column 1 to be `workspace`", row)
	}
	// If workspace value is provided, overwrite default
	if len(workspaceRow) >= 2 && len(workspaceRow[1]) > 0 {
		overwrite := workspaceRow[1]
		if !filepath.IsAbs(overwrite) {
			slog.Error(fmt.Sprintf("Workspace path %s is not an absolute path.", overwrite))
			return init, nil
		}
		stat, err := os.Lstat(overwrite)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Error(fmt.Sprintf("Workspace %s does not exist.", overwrite))
				return init, nil 
			}
			slog.Error(fmt.Sprintf("Failed to obtain information of workspace path %s", overwrite), "error", err)
			return init, nil
		}
		if !stat.IsDir() {
			slog.Error(fmt.Sprintf("Workspace %s is not a directory", overwrite))
			return init, nil
		}
		return overwrite, nil
	}
	return init, nil
}
