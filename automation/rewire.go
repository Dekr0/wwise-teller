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
	RewireTypeSound RewireType = 0
	RewireTypeCntr  RewireType = 1
)

type RewireType uint8;

type CSVHeader struct {
	Workspace  string
	Type       RewireType
	Conversion string
	Format     waapi.ConversionFormatType
	Project    string
}

func RewireWithNewSources(
	ctx         context.Context,
	bnk        *wwise.Bank,
	mappingFile string,
	dry         bool,
) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	iniDir := filepath.Base(mappingFile)
	header := CSVHeader{
		Workspace: iniDir,
		Type: 0,
	}
	if err = ParseRewireHeader(&header, reader); err != nil {
		return err
	}

	wavsMapSound := make(map[string][]*wwise.Sound, 32)
	if header.Type == RewireTypeSound {
		err = ParseSoundMapping(reader, &header, h, wavsMapSound)
		if err != nil {
			return err
		}
		if len(wavsMapSound) <= 0 {
			return nil
		}
	}

	wemsMapSound := make(map[string][]*wwise.Sound, len(wavsMapSound))
	wsource, err := waapi.CreateConversionList(ctx, wavsMapSound, wemsMapSound, header.Conversion, dry)
	if len(wavsMapSound) != len(wemsMapSound) {
		panic("Panic Trap")
	}
	if err != nil {
		return err
	}

	if dry {
		for wem, sounds := range wemsMapSound {
			fmt.Println(wem)
			for _, s := range sounds {
				fmt.Println(s.Id)
			}
			fmt.Println()
		}
		return nil
	}

	if err := waapi.WwiseConversion(ctx, wsource, header.Project); err != nil {
		return err
	}

	wemsMapAudioData := make(map[string][]byte, len(wemsMapSound))
	errorWems := make([]string, 0, len(wemsMapSound))
	for wem := range wemsMapSound {
		if _, in := wemsMapAudioData[wem]; in {
			panic("Panic Trap")
		}
		wemsMapAudioData[wem], err = os.ReadFile(wem)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to read audio data from %s", wem))
			slog.Error(fmt.Sprintf("Discard all rewiring related to file %s", wem))
			errorWems = append(errorWems, wem)
		}
	}
	for _, errorWem := range errorWems {
		delete(wemsMapSound, errorWem)
	}

	q, closeDb, commit, txRollback, err := db.CreateDefaultConnWithTxQuery(ctx)
	if err != nil {
		return err
	}
	defer closeDb()

	wemsMapMediaIndexs := make(map[string]wwise.MediaIndex, len(wemsMapSound))
	for wem := range wemsMapSound {
		sid, err := TrySid(ctx, q)
		if err != nil {
			txRollback()
			return fmt.Errorf("Failed to allocate a new source ID: %w", err)
		}
		if _, in := wemsMapMediaIndexs[wem]; in { panic("Panic Trap") }
		if audioData, in := wemsMapAudioData[wem]; !in { 
			panic("Panic Trap") 
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
			panic("Panic Trap")
		} else {
			if err := bnk.AppendAudio(audioData, m.Sid); err != nil {
				// Database is out of sync
				panic("Panic Trap")
			}
		}
	}

	for wem, sounds := range wemsMapSound {
		m, in := wemsMapMediaIndexs[wem]
		if !in { panic("Panic Trap") }
		for _, sound := range sounds {
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
func ParseRewireHeader(header *CSVHeader, reader *csv.Reader) error {
	conversionRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain conversion setting: %w", err)
	}
	if !strings.EqualFold(conversionRow[0], "conversion") {
		return fmt.Errorf("Expecting Row 1, Column 1 to be `conversion`")
	}
	if len(conversionRow) < 2 {
		return fmt.Errorf("Missing conversion setting value")
	}
	header.Conversion = conversionRow[1]

	formatRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain format setting: %w", err)
	}
	if !strings.EqualFold(formatRow[0], "format") {
		return fmt.Errorf("Expecting Row 2, Column 1 to be `format`")
	}
	if len(formatRow) < 2 {
		return fmt.Errorf("Missing format setting value")
	}
	format, err := strconv.Atoi(formatRow[1])
	if err != nil {
		return fmt.Errorf("Failed to parse format setting value: %w", err)
	}
	if format < int(waapi.ConversionFormatTypePCM) || 
	   format > int(waapi.ConversionFormatTypeWEMOpus) {
		return fmt.Errorf("Invalid format setting value %d", format)
	}
	header.Format = waapi.ConversionFormatType(format)

	projRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain project path: %w", err)
	}
	if !strings.EqualFold(projRow[0], "project") {
		return fmt.Errorf("Expecting Row 3, Column 1 to be `project`")
	}
	if len(projRow) < 2 {
		return fmt.Errorf("Missing project path")
	}
	if !filepath.IsAbs(projRow[1]) {
		return fmt.Errorf("Project path %s is not in absolute path", projRow[1])
	}
	if _, err := os.Lstat(projRow[1]); err != nil {
		return err
	}
	header.Project = projRow[1]

	rewireTypeRow, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain rewire type value: %w", err)
	}
	if !strings.EqualFold(rewireTypeRow[0], "type") {
		return fmt.Errorf("Expecting Row, 4, Column 1 to be `type`")
	}
	if len(rewireTypeRow) < 2 {
		return fmt.Errorf("Missing rewire type")
	}
	rewireType, err := strconv.Atoi(rewireTypeRow[1])
	if err != nil {
		return fmt.Errorf("Failed to parse rewire type: %w", err)
	}
	if rewireType < int(RewireTypeSound) || rewireType > int(RewireTypeCntr) {
		return fmt.Errorf("Invalid rewire type %d", rewireType)
	}
	header.Type = RewireType(rewireType)

	workspace_line, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain workspace directory: %w", err)
	}
	if !strings.EqualFold(workspace_line[0], "workspace") {
		return fmt.Errorf("Expecting Row 5, Column 1 to be `workspace`")
	}
	// If workspace value is provided, overwrite default
	if len(workspace_line) >= 2 && len(workspace_line[1]) > 0 {
		overwrite := workspace_line[1]
		if !filepath.IsAbs(overwrite) {
			return fmt.Errorf("Workspace path %s is not in full path format", overwrite)
		}
		_, err := os.Lstat(overwrite)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Workspace %s does not exist", overwrite)
			}
			return fmt.Errorf("Failed to obtain information of workspace path %s: %w", overwrite, err)
		}
		header.Workspace = overwrite
	}
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
	header *CSVHeader,
	h *wwise.HIRC,
	wavsMapping map[string][]*wwise.Sound,
) error {
	rowNum := 2
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
			rowNum += 1
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
			rowNum += 1
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
			rowNum += 1
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
		rowNum += 1
	}

	return nil
}
