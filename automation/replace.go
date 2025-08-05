package automation

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ReplaceAudioSources(
	ctx context.Context, bnk *wwise.Bank, mappingFile string, dry bool,
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

	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	iniDir := filepath.Base(mappingFile)
	header := WavSoundMapHeader{
		Workspace: iniDir,
		Type: 0,
	}
	row := uint16(0)
	if err = ParseFilesToSoundsHeader(&header, reader, &row); err != nil {
		return err
	}

	wavsMapSounds := make(map[string][]*wwise.Sound, 32)
	if header.Type == FilesToSoundsTypeGrain {
		err = ParseSoundMapping(reader, &header, h, wavsMapSounds, &row)
		if err != nil {
			return err
		}
		if len(wavsMapSounds) <= 0 {
			return nil
		}
	}

	wemsMapSounds := make(map[string][]*wwise.Sound, len(wavsMapSounds))
	wsource, err := waapi.CreateConversionList(ctx, wavsMapSounds, wemsMapSounds, header.Conversion, dry)
	if len(wavsMapSounds) != len(wemsMapSounds) {
		panic("# of wavs file in indexing map does not equal # of output wem in indexing map")
	}
	if err != nil {
		return err
	}
	stagingDir := filepath.Dir(wsource)
	defer os.RemoveAll(stagingDir)

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

	for wem, sounds := range wemsMapSounds {
		audioData, in := wemsMapAudioData[wem]
		if !in {
			panic(fmt.Sprintf("No audio data is mapped to %s", wem))
		}
		for _, sound := range sounds {
			if err := bnk.ReplaceAudio(audioData, sound.BankSourceData.SourceID); err != nil {
				return err
			}
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
			sound.BankSourceData.InMemoryMediaSize = uint32(len(audioData))
		}
	}
	bnk.ComputeDIDXOffset()
	if err := bnk.CheckDIDXDATA(); err != nil {
		return err
	}
	return nil
}

type WavSoundMap struct {
	Header     WavSoundMapHeader
	Wavs     []string
	SoundIDs [][]uint32
}

func (p *WavSoundMap) Encode(out string) error {
	var b strings.Builder
	b.WriteString(p.Header.Conversion + "\n")
	b.WriteString(strconv.FormatUint(uint64(p.Header.Format), 10) + "\n")
	b.WriteString(strconv.FormatUint(uint64(p.Header.Type), 10) + "\n")
	b.WriteString(p.Header.Workspace + "\n")

	for i := range p.Wavs {
		b.WriteString(fmt.Sprintf("%s,%d,", p.Wavs[i], len(p.SoundIDs[i])))
		for j := range p.SoundIDs[i] {
			b.WriteString(strconv.FormatUint(uint64(p.SoundIDs[i][j]), 10) + ",")
			if j < len(p.SoundIDs) - 1 {
				b.WriteByte(',')
			}
		}
		if i < len(p.Wavs) - 1 {
			b.WriteByte('\n')
		}
	}

	return os.WriteFile(out, []byte(b.String()), 0666) 
}

func DecodeReplaceAudioSourcesScript(p *WavSoundMap, fspec string) error {
	f, err := os.Open(fspec)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	iniDir := filepath.Base(fspec)
	p.Header.Workspace = iniDir
	p.Header.Type = 0
	
	rowNum := uint16(0)
	if err = ParseFilesToSoundsHeader(&p.Header, reader, &rowNum); err != nil {
		return err
	}

	var row []string
	var input, ext string
	for {
		row, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if len(row) < 2 {
			slog.Error(fmt.Sprintf("Expecting two columns, file name and # of targeting sound IDs, at row %d", rowNum))
			rowNum += 1
		}

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

		if !filepath.IsAbs(input) {
			input = filepath.Join(p.Header.Workspace, input)
		}

		_, err = os.Lstat(input)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to obtain info of wave file %s", input))
			rowNum += 1
		}

		count, err := strconv.ParseUint(row[1], 10, 8)
		if err != nil {
			slog.Error("Failed to parse # of targeting sound IDs", "error", err)
			continue
		}
		if count > uint64(len(row) - 2) {
			slog.Error(fmt.Sprintf("Expecting %d of targeting sound IDs but receiving %d", count, len(row)- 2))
			continue
		}

		idx := slices.Index(p.Wavs, input)
		columnNum := 3
		if idx == -1 {
			soundIDs := make([]uint32, 0, count)
			for _, i := range row[2:] {
				id, err := strconv.ParseUint(i, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse sound ID at column %d", columnNum), "error", err)
					columnNum += 1
					continue
				}
				soundIDs = append(soundIDs, uint32(id))
				columnNum += 1
			}
			p.SoundIDs = append(p.SoundIDs, soundIDs)
		} else {
			slog.Warn(fmt.Sprintf("Duplicated input file %s at row %d, column 1", input, rowNum))
			for _, i := range row[2:] {
				id, err := strconv.ParseUint(i, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse sound ID at column %d", columnNum), "error", err)
					columnNum += 1
					continue
				}
				p.SoundIDs[idx] = append(p.SoundIDs[idx], uint32(id))
				columnNum += 1
			}
		}

		rowNum += 1
	}

	return nil
}
