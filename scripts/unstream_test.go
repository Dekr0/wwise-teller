package scripts

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestTuneEagleBank(t *testing.T) {
	bnks := []string{
		"../scripts/output/content_audio_stratagems_eagle_airstrike.st_bnk",
		"../scripts/output/content_audio_stratagems_eagle_500kg_bomb.st_bnk",
		"../scripts/output/content_audio_stratagems_eagle_110m_rockets.st_bnk",
		"../scripts/output/content_audio_stratagems_eagle_strafing_run.st_bnk",
		"../scripts/output/content_audio_stratagems_eagle_napalm_strike.st_bnk",
		"../scripts/output/content_audio_stratagems_eagle_airstrike_smoke.st_bnk",
	}
	soundIds := []uint32{
		66293742,
		439413948,
		1035272917,
	}
	engineID := uint32(164074760)
	coreLoopID := uint32(418261766)
	impactID := uint32(351516354)
	bg := context.Background()
	for _, bnk := range bnks {
		b, err := parser.ParseBank(bnk, bg, false)
		if err != nil {
			t.Fatal(err)
		}
		hirc := b.HIRC()
		if hirc == nil {
			t.Fatal()
		}
		for _, soundId := range soundIds {
			v, ok := hirc.ActorMixerHirc.Load(soundId)
			if !ok {
				t.Fatal()
			}
			sound := v.(*wwise.Sound)
			sound.BankSourceData.StreamType = wwise.SourceTypeDATA
		}

		v, ok := hirc.ActorMixerHirc.Load(engineID)
		if !ok {
			t.Fatal()
		}
		r := v.(*wwise.RanSeqCntr)
		r.BaseParam.AdvanceSetting.SetIgnoreParentMaxNumInst(true)
		r.BaseParam.AdvanceSetting.MaxNumInstance = 32

		if filepath.Base(bnk) == "content_audio_stratagems_eagle_strafing_run.st_bnk" {
			t.Log("Tunning strafing run")
			v, ok = hirc.ActorMixerHirc.Load(coreLoopID)
			if !ok {
				t.Fatal()
			}
			r = v.(*wwise.RanSeqCntr)
			r.BaseParam.AdvanceSetting.SetIgnoreParentMaxNumInst(true)
			r.BaseParam.AdvanceSetting.MaxNumInstance = 32
			v, ok = hirc.ActorMixerHirc.Load(impactID)
			if !ok {
				t.Fatal()
			}
			r = v.(*wwise.RanSeqCntr)
			r.BaseParam.AdvanceSetting.SetIgnoreParentMaxNumInst(true)
			r.BaseParam.AdvanceSetting.MaxNumInstance = 32
		}

		data, err := b.Encode(bg, false, false)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Base(bnk), data, 0666)
		if err != nil {
			t.Fatal(err)
		}
	}
}
