package automation

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
	header := FilesToSoundsHeader{
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
