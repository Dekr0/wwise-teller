package main

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestRemoveActorMixerCntrChild2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	bnk, err := parser.ParseBank("./tests/st_bnk/content_audio_p2_peacemaker.st_bnk", ctx)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	hirc := bnk.HIRC()

	idx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		return h.HircType() == 0x03
	})

	l := len(hirc.HircObjs)

	v, in := hirc.HircObjsMap.Load(uint32(99586918))
	if !in {
		t.Fatalf("ID %d does not have an associated switch container", 99586918)
	}
	switchOne := v.(*wwise.SwitchCntr)
	v, in = hirc.HircObjsMap.Load(uint32(338060418))
	if !in {
		t.Fatalf("ID %d does not have an associated switch container", 338060418)
	}
	switchTwo := v.(*wwise.SwitchCntr)

	v, in = hirc.HircObjsMap.Load(uint32(662359126))
	if !in {
		t.Fatalf("ID %d does not have an associated actor mixer", 662359126)
	}
	actorMixer := v.(*wwise.ActorMixer)

	hirc.RemoveRoot(99586918, 662359126)
	hirc.RemoveRoot(338060418, 662359126)

	if switchOne.BaseParam.DirectParentId != 0 {
		t.Fatalf("Switch container 99586918 parent ID is not zero")
	}
	if switchTwo.BaseParam.DirectParentId != 0 {
		t.Fatalf("Switch container 99586918 parent ID is not zero")
	}
	if len(actorMixer.Container.Children) != 0 {
		t.Fatalf("Actor mixer 662359126 still have children")
	}

	if len(hirc.HircObjs) != l {
		t.Fatalf("Expected: %d, Got: %d", l, len(hirc.HircObjs))
	}

	newIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		return h.HircType() == 0x03
	})

	if idx != newIdx {
		t.Fatalf("Expected: %d, Got: %d", idx, newIdx)
	}

	idx = slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 99586918
	})
	if idx != 400 {
		t.Fatalf("Expected: %d, Got: %d", 400, idx)
	}
	idx = slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 338060418
	})
	if idx != 401 {
		t.Fatalf("Expected: %d, Got: %d", 401, idx)
	}

	got, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(ref)-8 != len(got) {
		t.Fatalf("Expected: %d, Got: %d", len(ref), len(got))
	}
}

func TestRemoveRanSeqCntrChild(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	bnk, err := parser.ParseBank("./tests/st_bnk/content_audio_p2_peacemaker.st_bnk", ctx)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	hirc := bnk.HIRC()

	soundIDs := []uint32{
		5565038,
		66205529,
		66509257,
		135118204,
		145912119,
		169359783,
	}

	idx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		return h.HircType() == 0x03
	})

	v, in := hirc.HircObjsMap.Load(uint32(351417743))
	if !in {
		t.Fatalf("%d does not have an associated random / sequence container object.", 351417743)
	}
	ranSeqCntr := v.(*wwise.RanSeqCntr)

	for _, soundID := range soundIDs {
		v, in := hirc.HircObjsMap.Load(soundID)
		if !in {
			t.Fatalf("%d does not have an associated sound object.", soundID)
		}
		sound := v.(*wwise.Sound)
		hirc.RemoveRoot(soundID, 351417743)
		if sound.BaseParam.DirectParentId != 0 {
			t.Fatalf("Sound %d parent ID is non zero", soundID)
		}
	}

	for _, soundID := range soundIDs {
		if slices.Contains(ranSeqCntr.Container.Children, soundID) {
			t.Fatalf("Sound %d is still in random sequence container %d", soundID, 351417743)
		}
		if slices.ContainsFunc(ranSeqCntr.PlayListItems, func(p *wwise.PlayListItem) bool {
			return p.UniquePlayID == soundID
		}) {
			t.Fatalf("Sound %d is still in the playlist item of random sequence container %d", soundID, 351417743)
		}
	}

	l := len(hirc.HircObjs)

	if len(hirc.HircObjs) != l {
		t.Fatalf("Expected: %d, Got: %d", l, len(hirc.HircObjs))
	}

	newIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		return h.HircType() == 0x03
	})

	if idx != newIdx {
		t.Fatalf("Expected: %d, Got: %d", idx, newIdx)
	}

	for i, soundID := range soundIDs {
		sIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
			id, err := h.HircID()
			if err != nil {
				return false
			}
			return id == soundID
		})
		expect := newIdx - len(soundIDs) + i
		if sIdx != expect {
			t.Fatalf("Expect %d, Got %d", expect, sIdx)
		}
	}

	got, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(ref)-(len(soundIDs)*4+(len(soundIDs)*(4+4))) != len(got) {
		t.Fatalf("Expected: %d, Got: %d", len(ref), len(got))
	}
}

func TestChangeSoundRoot(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	bnk, err := parser.ParseBank("./tests/st_bnk/content_audio_p2_peacemaker.st_bnk", ctx)
	if err != nil {
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	ref, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	v, in := hirc.HircObjsMap.Load(uint32(862008135))
	if !in {
		t.Fatalf("%d does not have an associated random / sequence container", 862008135)
	}
	oldSeq := v.(*wwise.RanSeqCntr)
	prevOldSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 862008135
	})
	if prevOldSeqIdx == -1 {
		t.Fatalf("%d is not in HIRC", 862008135)
	}

	v, in = hirc.HircObjsMap.Load(uint32(114819736))
	if !in {
		t.Fatalf("%d does not have an associated random / sequence container", 114819736)
	}
	newSeq := v.(*wwise.RanSeqCntr)
	prevNewSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 114819736
	})
	if prevNewSeqIdx == -1 {
		t.Fatalf("%d is not in HIRC", 114819736)
	}
	prevNewSeqNumLeaf := newSeq.NumLeaf()

	soundIDs := []uint32{
		69870573,
		268669903,
		486185470,
		794408455,
		809283960,
		946668925,
	}
	oldSoundIdxs := make([]int, len(soundIDs))
	for i, soundID := range soundIDs {
		v, in := hirc.HircObjsMap.Load(soundID)
		if !in {
			t.Fatalf("%d does not have an associated sound", soundID)
		}
		sound := v.(*wwise.Sound)

		oldSoundIdxs[i] = hirc.TreeArrIdx(soundID)
		if oldSoundIdxs[i] == -1 {
			t.Fatalf("%d is not in HIRC", soundID)
		}

		hirc.ChangeRoot(soundID, newSeq.Id, oldSeq.Id)
		if sound.ParentID() != newSeq.Id {
			t.Fatalf("Expect %d, Got %d", newSeq.Id, sound.ParentID())
		}
	}

	if len(oldSeq.Container.Children) != 0 {
		t.Fatalf("Expect %d, Got %d", 0, len(oldSeq.Container.Children))
	}
	if len(oldSeq.PlayListItems) != 0 {
		t.Fatalf("Expect %d, Got %d", 0, len(oldSeq.PlayListItems))
	}
	if len(newSeq.Container.Children) != prevNewSeqNumLeaf+len(soundIDs) {
		t.Fatalf("Expect %d, Got %d", prevNewSeqNumLeaf+len(soundIDs), len(newSeq.Container.Children))
	}

	newOldSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 862008135
	})
	if newOldSeqIdx != prevOldSeqIdx {
		t.Fatalf("Expect %d, Got %d", prevOldSeqIdx-len(soundIDs), newOldSeqIdx)
	}

	newNewSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 114819736
	})
	if newNewSeqIdx != prevNewSeqIdx+len(soundIDs) {
		t.Fatalf("Expect %d, Got %d", prevNewSeqIdx+len(soundIDs), newNewSeqIdx)
	}

	for i, soundID := range soundIDs {
		if !slices.Contains(newSeq.Container.Children, soundID) {
			t.Fatalf("%d is not in random / sequence container %d", soundID, newSeq.Id)
		}

		idx := hirc.TreeArrIdx(soundID)
		if oldSoundIdxs[i]-1 != idx {
			t.Fatalf("Expect %d, Got %d", oldSoundIdxs[i]-1, idx)
		}
	}

	got, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(ref)-len(soundIDs)*(4+4) != len(got) {
		t.Fatalf("Expected: %d, Got: %d", len(ref), len(got))
	}
}

func TestChangeSoundPartial(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	bnk, err := parser.ParseBank("./tests/st_bnk/content_audio_p2_peacemaker.st_bnk", ctx)
	if err != nil {
		t.Fatal(err)
	}
	hirc := bnk.HIRC()

	ref, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	v, in := hirc.HircObjsMap.Load(uint32(862008135))
	if !in {
		t.Fatalf("%d does not have an associated random / sequence container", 862008135)
	}
	oldSeq := v.(*wwise.RanSeqCntr)
	prevOldSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 862008135
	})
	if prevOldSeqIdx == -1 {
		t.Fatalf("%d is not in HIRC", 862008135)
	}

	v, in = hirc.HircObjsMap.Load(uint32(114819736))
	if !in {
		t.Fatalf("%d does not have an associated random / sequence container", 114819736)
	}
	newSeq := v.(*wwise.RanSeqCntr)
	prevNewSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 114819736
	})
	if prevNewSeqIdx == -1 {
		t.Fatalf("%d is not in HIRC", 114819736)
	}
	prevNewSeqNumLeaf := newSeq.NumLeaf()

	soundIDs := []uint32{
		69870573,
		268669903,
		486185470,
	}
	for _, soundID := range soundIDs {
		v, in := hirc.HircObjsMap.Load(soundID)
		if !in {
			t.Fatalf("%d does not have an associated sound", soundID)
		}
		sound := v.(*wwise.Sound)
		hirc.ChangeRoot(soundID, newSeq.Id, oldSeq.Id)
		if sound.ParentID() != newSeq.Id {
			t.Fatalf("Expect %d, Got %d", newSeq.Id, sound.ParentID())
		}
	}

	if len(oldSeq.Container.Children) != 3 {
		t.Fatalf("Expect %d, Got %d", 3, len(oldSeq.Container.Children))
	}
	if len(oldSeq.PlayListItems) != 3 {
		t.Fatalf("Expect %d, Got %d", 3, len(oldSeq.PlayListItems))
	}
	if len(newSeq.Container.Children) != prevNewSeqNumLeaf+len(soundIDs) {
		t.Fatalf("Expect %d, Got %d", prevNewSeqNumLeaf+len(soundIDs), len(newSeq.Container.Children))
	}
	for _, soundID := range soundIDs {
		if !slices.Contains(newSeq.Container.Children, soundID) {
			t.Fatalf("%d is not in random / sequence container %d", soundID, newSeq.Id)
		}
	}

	newOldSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 862008135
	})
	if newOldSeqIdx != prevOldSeqIdx {
		t.Fatalf("Expect %d, Got %d", prevOldSeqIdx-len(soundIDs), newOldSeqIdx)
	}

	newNewSeqIdx := slices.IndexFunc(hirc.HircObjs, func(h wwise.HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == 114819736
	})
	if newNewSeqIdx != prevNewSeqIdx+len(soundIDs) {
		t.Fatalf("Expect %d, Got %d", prevNewSeqIdx+len(soundIDs), newNewSeqIdx)
	}

	got, err := bnk.Encode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(ref)-len(soundIDs)*(4+4) != len(got) {
		t.Fatalf("Expected: %d, Got: %d", len(ref), len(got))
	}
}

func TestChangeSoundMixerSeq(t *testing.T) {

}

func TestMain(t *testing.T) {
	main()
}
