package experiment

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
)

var TestBankDir = "../tests/st_bnks"
var TestWemsDir = "../tests/wems/"
var TestPatchDir = "../tests/patch"

func TestAppendRanSeqCntr(t *testing.T) {
	const bank = "content_audio_stratagems_sentry_machine_gun.st_bnk"
	const audio  = "2a72_fire_echo_mid.wem"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	bnk, err := parser.ParseBank(filepath.Join(TestBankDir, bank), ctx, false)
	if err != nil {
		cancel()
		t.Fatal(err)
	}

	h := bnk.HIRC()

	// Create a new source
	newSourceID, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}
	audioData, err := os.ReadFile(filepath.Join(TestWemsDir, audio))
	if err != nil {
		t.Fatal(err)
	}
	bnk.AppendAudio(audioData, newSourceID)
	
	// Create a new sound, and wire up the new source to the new sound
	newSoundID, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}

	const refSoundId = 966968349
	refSoundIdx := slices.IndexFunc(h.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == refSoundId
	})
	if refSoundIdx == -1 {
		t.Fatalf("Sound %d does not exist.", refSoundId)
	}
	refSound := h.HircObjs[refSoundIdx].(*wwise.Sound)

	newSound := wwise.Sound{
		Id: newSoundID,
		BankSourceData: wwise.BankSourceData{
			PluginID: wwise.VORBIS,
			StreamType: wwise.SourceTypeDATA,
			SourceID: newSourceID,
			InMemoryMediaSize: uint32(len(audioData)),
			SourceBits: 0,
		},
		BaseParam: refSound.BaseParam.Clone(false),
	}

	// Create random sequence container
	const refRanSeqCntrId = 543920301
	newRanSeqCntrID, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}

	refRanSeqCntrIdx := slices.IndexFunc(h.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == refRanSeqCntrId
	})
	if refRanSeqCntrIdx == -1 {
		t.Fatalf("Random / Sequence Container %d does not exist.", refRanSeqCntrId)
	}
	refRanSeqCntr := h.HircObjs[refRanSeqCntrIdx].(*wwise.RanSeqCntr)
	newRanSeqCntr := refRanSeqCntr.Clone(newRanSeqCntrID, false)

	// Append random sequence container to actor mixer
	const actorMixerId = 455686496
	if err := h.AppendNewRanSeqCntrToActorMixer(&newRanSeqCntr, actorMixerId, false); err != nil {
		t.Fatal(err)
	}

	// Append the new sound to the new sequence container
	if err := h.AppendNewSoundToRanSeqContainer(&newSound, newRanSeqCntrID, false); err != nil {
		t.Fatal(err)
	}

	// Create new action
	newActionId, err := utils.ShortID()
	if err != nil {
		t.Fatal(err)
	}

	const refActionId = 644905578
	refActionIdx := slices.IndexFunc(h.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == refActionId
	})
	if refActionIdx == -1 {
		t.Fatalf("Action %d does not exist.", refActionId)
	}
	refAction := h.HircObjs[refActionIdx].(*wwise.Action)
	newAction := refAction.Clone(newActionId, newRanSeqCntrID)

	// Append new action on Event
	const eventID = 3683182397
	if err := h.AppendNewActionToEvent(&newAction, eventID); err != nil {
		t.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	data, err := bnk.Encode(ctx, true, false)
	if err != nil {
		t.Fatal(err)
	}
	helldivers.GenHelldiversPatchStable(data, bnk.META().B, "../tests/patch")
}
