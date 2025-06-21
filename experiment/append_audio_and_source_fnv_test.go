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
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestAppendAudioUsingFNV(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	bnk, err := parser.ParseBank("../tests/st_bnks/wep_cr1_adjudicator.st_bnk", ctx)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	audioDataDesert, err := os.ReadFile("../tests/wems/reflection_close_desert_00.wem")
	if err != nil {
		t.Fatal(err)
	}
	desertSid, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioDataDesert, desertSid)
	for i := range hirc.HircObjs {
		switch s := hirc.HircObjs[i].(type) {
		case *wwise.Sound:
			if s.BaseParam.DirectParentId == 274049716 {
				s.BankSourceData.PluginID = wwise.VORBIS
				s.BankSourceData.SourceID = desertSid
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

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	data, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, "../tests/patch")
}

func TestAppendAudioAndSoundUsingFNV(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	bnk, err := parser.ParseBank("../tests/st_bnks/wep_cr1_adjudicator.st_bnk", ctx)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	audioDataDesert, err := os.ReadFile("../tests/wems/reflection_close_desert_00.wem")
	if err != nil {
		t.Fatal(err)
	}
	desertSid, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioDataDesert, desertSid)
	for i := range hirc.HircObjs {
		switch s := hirc.HircObjs[i].(type) {
		case *wwise.Sound:
			if s.BaseParam.DirectParentId == 274049716 {
				s.BankSourceData.PluginID = wwise.VORBIS
				s.BankSourceData.SourceID = desertSid
				s.BankSourceData.InMemoryMediaSize = uint32(len(audioDataDesert))
			}
		case *wwise.RanSeqCntr:
			if s.Id == 274049716 {
			}
			if s.Id == 435636362 {
				s.BaseParam.PropBundle.AddBaseProp()	
				s.BaseParam.PropBundle.ChangeBaseProp(0, wwise.PropTypeMakeUpGain)
			}
			if s.Id == 98920475 {
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
	if err != nil {
		t.Fatal(err)
	}
	urban0Sid, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioDataUrban0, urban0Sid)

	soundAId, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	soundA := wwise.Sound{
		Id: soundAId,
		BankSourceData: wwise.BankSourceData{
			PluginID: wwise.VORBIS,
			StreamType: wwise.STREAM_TYPE_BNK,
			SourceID: urban0Sid,
			InMemoryMediaSize: uint32(len(audioDataUrban0)),
			SourceBits: 0,
		},
		BaseParam: ref.BaseParam.Clone(false),
	}

	audioDataUrban1, err := os.ReadFile("../tests/wems/reflection_close_urban_01.wem")
	if err != nil {
		t.Fatal(err)
	}
	urban1Sid, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioDataUrban1, urban1Sid)

	soundBId, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	soundB := wwise.Sound{
		Id: soundBId,
		BankSourceData: wwise.BankSourceData{
			PluginID: wwise.VORBIS,
			StreamType: wwise.STREAM_TYPE_BNK,
			SourceID: urban1Sid,
			InMemoryMediaSize: uint32(len(audioDataUrban1)),
			SourceBits: 0,
		},
		BaseParam: ref.BaseParam.Clone(false),
	}

	hirc.AppendNewSoundToRanSeqContainer(&soundA, 274049716)
	hirc.AppendNewSoundToRanSeqContainer(&soundB, 274049716)

	fmt.Println("desertSid", desertSid)
	fmt.Println("urban0Sid", urban0Sid)
	fmt.Println("urban1Sid", urban1Sid)
	fmt.Println("soundAId", soundAId)
	fmt.Println("soundBId", soundBId)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	data, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, "../tests/patch")
}
