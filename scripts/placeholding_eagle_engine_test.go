package scripts_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestPlaceholdingEagleEngine(t *testing.T) {
	paths := []string{
		"../tests/default_st_bnks/content_audio_stratagems_eagle_airstrike.st_bnk",
		"../tests/default_st_bnks/content_audio_stratagems_eagle_500kg_bomb.st_bnk",
		"../tests/default_st_bnks/content_audio_stratagems_eagle_110m_rockets.st_bnk",
		"../tests/default_st_bnks/content_audio_stratagems_eagle_strafing_run.st_bnk",
		"../tests/default_st_bnks/content_audio_stratagems_eagle_napalm_strike.st_bnk",
		"../tests/default_st_bnks/content_audio_stratagems_eagle_airstrike_smoke.st_bnk",
	}
	bnks := make([]*wwise.Bank, 0, len(paths))
	for _, path := range paths {
		bnk, err := parser.ParseBank(path, t.Context(), false)
		bnks = append(bnks, bnk)
		if err != nil {
			t.Fatal(err)
		}
	}
	v, in := bnks[0].HIRC().ActorMixerHirc.Load(uint32(66293742))
	if !in {
		t.Fatalf("No reference Sound object with ID %d", 66293742)
	}
	newEngineSoundIds := []uint32{}
	newEngineSourceIds := []uint32{}
	for range 4 {
		soundId, err := utils.ShortID()
		if err != nil {
			t.Fatal(err)
		}
		newEngineSoundIds = append(newEngineSoundIds, soundId)
		sourceId, err := utils.ShortID()
		if err != nil {
			t.Fatal(err)
		}
		newEngineSourceIds = append(newEngineSourceIds, sourceId)
	}
	refSound := v.(*wwise.Sound)
	newSounds := []*wwise.Sound{}
	for i, soundId := range newEngineSoundIds {
		newSounds = append(newSounds, &wwise.Sound{
			Id: soundId,
			BankSourceData: wwise.BankSourceData{
				PluginID: wwise.VORBIS,
				StreamType: wwise.SourceTypeDATA,
				SourceID: newEngineSourceIds[i],
				InMemoryMediaSize: 0,
				SourceBits: 0,
			},
			BaseParam: refSound.BaseParam.Clone(false),
		})
	}
	for i, bnk := range bnks {
		for i, sourceId := range newEngineSourceIds {
			if err := bnk.AppendAudio([]byte{}, sourceId); err != nil {
				t.Fatal(err)
			}
			if err := bnk.HIRC().AppendNewSoundToRanSeqContainer(newSounds[i], 164074760, false); err != nil {
				t.Fatal(paths[i], err)
			}
		}
		data, err := bnk.Encode(t.Context(), false, false)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Base(paths[i]), data, 0666); err != nil {
			t.Fatal(err)
		}
		for _, newSound := range newSounds {
			newSound.BaseParam.DirectParentId = 0
		}
	}
}
