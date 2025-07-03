package experiment

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestAppendAudio(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	bnk, err := parser.ParseBank("../tests/st_bnks/wep_cr1_adjudicator.st_bnk", ctx, false)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	audioData, err := os.ReadFile("../tests/wems/reflection_close_desert_00.wem.wem")
	if err != nil {
		t.Fatal(err)
	}

	bnk.AppendAudio(audioData, 26007159)
	for _, h := range hirc.HircObjs {
		switch s := h.(type) {
		case *wwise.Sound:
			if s.BaseParam.DirectParentId == 274049716 {
				s.BankSourceData.PluginID = 0x00040001
				s.BankSourceData.SourceID = 26007159
				s.BankSourceData.InMemoryMediaSize = uint32(len(audioData))
			}
		case *wwise.RanSeqCntr:
			if s.Id == 274049716 {
			}
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	data, err := bnk.Encode(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, "../tests/patch")
}

func TestAppendAudioAndSound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	bnk, err := parser.ParseBank("../tests/st_bnks/wep_cr1_adjudicator.st_bnk", ctx, false)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	audioDataDesert, err := os.ReadFile("../tests/wems/reflection_close_desert_00.wem")
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioDataDesert, 26007159)
	for i := range hirc.HircObjs {
		switch s := hirc.HircObjs[i].(type) {
		case *wwise.Sound:
			if s.BaseParam.DirectParentId == 274049716 {
				s.BankSourceData.PluginID = wwise.VORBIS
				s.BankSourceData.SourceID = 26007159
				s.BankSourceData.InMemoryMediaSize = uint32(len(audioDataDesert))
			}
		case *wwise.RanSeqCntr:
			if s.Id == 274049716 {
				fmt.Println("274049716")
			}
			if s.Id == 435636362 {
				fmt.Println("435636362")
				s.BaseParam.PropBundle.AddBaseProp()	
				s.BaseParam.PropBundle.ChangeBaseProp(0, wwise.PropTypeMakeUpGain)
			}
			if s.Id == 98920475 {
				fmt.Println("98920475")
			}
		}
	}

	idx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 80388110
	})
	ref := hirc.HircObjs[idx].(*wwise.Sound)

	audioDataUrban0, err := os.ReadFile("../tests/wems/reflection_close_urban_00.wem")
	bnk.AppendAudio(audioDataUrban0, 568558003)
	soundA := wwise.Sound{
		Id: 415678546,
		BankSourceData: wwise.BankSourceData{
			PluginID: wwise.VORBIS,
			StreamType: wwise.STREAM_TYPE_BNK,
			SourceID: 568558003,
			InMemoryMediaSize: uint32(len(audioDataUrban0)),
			SourceBits: 0,
		},
		BaseParam: ref.BaseParam.Clone(false),
	}

	audioDataUrban1, err := os.ReadFile("../tests/wems/reflection_close_urban_01.wem")
	bnk.AppendAudio(audioDataUrban1, 107266693)
	soundB := wwise.Sound{
		Id: 234198813,
		BankSourceData: wwise.BankSourceData{
			PluginID: wwise.VORBIS,
			StreamType: wwise.STREAM_TYPE_BNK,
			SourceID: 107266693,
			InMemoryMediaSize: uint32(len(audioDataUrban1)),
			SourceBits: 0,
		},
		BaseParam: ref.BaseParam.Clone(false),
	}

	hirc.AppendNewSoundToRanSeqContainer(&soundA, 274049716, false)
	hirc.AppendNewSoundToRanSeqContainer(&soundB, 274049716, false)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	data, err := bnk.Encode(ctx, false)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, "../tests/patch")
}
