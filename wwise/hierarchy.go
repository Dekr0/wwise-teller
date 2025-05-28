package wwise

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const MaxEncodeRoutine = 6

// # of hierarchy object (uint32)
const SizeOfHIRCHeader = 4

type HircType uint8

const (
	HircTypeState           HircType = 0x01 // ???
	HircTypeSound           HircType = 0x02
	HircTypeAction          HircType = 0x03 // ???
	HircTypeEvent           HircType = 0x04 // ???
	HircTypeRanSeqCntr      HircType = 0x05
	HircTypeSwitchCntr      HircType = 0x06
	HircTypeActorMixer      HircType = 0x07
	HircTypeBus             HircType = 0x08 // ???
	HircTypeLayerCntr       HircType = 0x09
	HircTypeMusicSegment    HircType = 0x0a
	HircTypeMusicTrack      HircType = 0x0b // ???
	HircTypeMusicSwitchCntr HircType = 0x0c // ???
	HircTypeMusicRanSeqCntr HircType = 0x0d // ???
	HircTypeAttenuation     HircType = 0x0e // ???
	HircTypeDialogueEvent   HircType = 0x0f // ???
	HircTypeFxShareSet      HircType = 0x10 // ???
	HircTypeFxCustom        HircType = 0x11 // ???
	HircTypeAuxBus          HircType = 0x12 // ???
	HircTypeLFOModulator    HircType = 0x13 // ???
	HircEnvelopeModulator   HircType = 0x14 // ???
	HircAudioDevice         HircType = 0x15 // ???
	HircTimeModulator       HircType = 0x16 // ???
)

var KnownHircType []HircType = []HircType{
	0x00,
	HircTypeSound,
	HircTypeEvent,
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
	HircTypeMusicSegment,
	HircTypeMusicTrack,
	HircTypeMusicSwitchCntr,
	HircTypeMusicRanSeqCntr,
}

var ContainerHircType []HircType = []HircType{
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
	HircTypeMusicSegment,
	HircTypeMusicRanSeqCntr,
	HircTypeMusicSwitchCntr,
}

var HircTypeName []string = []string{
	"",
	"State",
	"Sound",
	"Action",
	"Event",
	"Random / Sequence Container",
	"Switch Container",
	"Actor Mixer",
	"Bus",
	"Layer Container",
	"Music Segment",
	"Music Track",
	"Music Switch Container",
	"Music Random / Sequence Container",
	"Attenuation",
	"Dialogue Event",
	"FX Share Set",
	"FX Custom",
	"Aux Bus",
	"LFO Modulator",
	"Envelope Modulator",
	"Audio Device",
	"Time Modulator",
}

type HIRC struct {
	I uint8
	T []byte

	// Retain the tree structure that comes from the decoding with minimal
	// modification
	HircObjs    []HircObj
	HircObjsMap *sync.Map

	// Map for different types of hierarchy objects. Each object is a pointer
	// to a specific hierarchy object, which is also in `HircObjs`.
	ActorMixers     *sync.Map
	Events          *sync.Map
	LayerCntrs      *sync.Map
	MusicSegments   *sync.Map
	MusicTracks     *sync.Map
	MusicRanSeqCntr *sync.Map
	MusicSwitchCntr *sync.Map
	SwitchCntrs     *sync.Map
	RanSeqCntrs     *sync.Map
	Sounds          *sync.Map
}

func NewHIRC(I uint8, T []byte, numHircItem uint32) *HIRC {
	return &HIRC{
		I:               I,
		T:               T,
		HircObjs:        make([]HircObj, numHircItem),
		HircObjsMap:     &sync.Map{},
		ActorMixers:     &sync.Map{},
		Events:          &sync.Map{},
		LayerCntrs:      &sync.Map{},
		MusicSegments:   &sync.Map{},
		MusicTracks:     &sync.Map{},
		MusicRanSeqCntr: &sync.Map{},
		MusicSwitchCntr: &sync.Map{},
		SwitchCntrs:     &sync.Map{},
		RanSeqCntrs:     &sync.Map{},
		Sounds:          &sync.Map{},
	}
}

func (h *HIRC) encode(ctx context.Context) ([]byte, error) {
	type result struct {
		i int
		b []byte
	}

	// No initialization since I want it to crash and catch encoding bugs
	results := make([][]byte, len(h.HircObjs))

	// sync signal
	c := make(chan *result, MaxEncodeRoutine)

	// limit # of go routines running at the same time
	sem := make(chan struct{}, MaxEncodeRoutine)

	done := 0
	i := 0
	for done < len(h.HircObjs) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case r := <-c:
			results[r.i] = r.b
			done += 1
		case sem <- struct{}{}:
			if i < len(h.HircObjs) {
				j := i
				go func() {
					c <- &result{j, h.HircObjs[j].Encode()}
					<-sem
				}()
				i += 1
			}
		default:
			if i < len(h.HircObjs) {
				results[i] = h.HircObjs[i].Encode()
				done += 1
				i += 1
			}
		}
	}

	return bytes.Join(results, []byte{}), nil
}

func (h *HIRC) Encode(ctx context.Context) ([]byte, error) {
	b, err := h.encode(ctx)
	if err != nil {
		return nil, err
	}

	dataSize := uint32(SizeOfHIRCHeader + len(b))
	size := SizeOfChunkHeader + dataSize
	w := wio.NewWriter(uint64(SizeOfChunkHeader + dataSize))
	w.AppendBytes(h.T)
	w.Append(dataSize)
	w.Append(uint32(len(h.HircObjs)))
	w.AppendBytes(b)
	return w.BytesAssert(int(size)), nil
}

func (h *HIRC) Tag() []byte {
	return h.T
}

func (h *HIRC) Idx() uint8 {
	return h.I
}

func (h *HIRC) ChangeRoot(id, newRootID, oldRootID uint32) {
	if newRootID == 0 {
		h.RemoveRoot(id, oldRootID)
		return
	}

	v, in := h.HircObjsMap.Load(id)
	if !in {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", id))
	}
	leaf := v.(HircObj)
	if b := leaf.BaseParameter(); b == nil {
		panic(fmt.Sprintf("%d is not containable", id))
	}
	_, err := leaf.HircID()
	if err != nil {
		panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", id))
	}

	v, in = h.HircObjsMap.Load(newRootID)
	if !in {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", newRootID))
	}
	newRoot := v.(HircObj)
	_, err = newRoot.HircID()
	if err != nil {
		panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", newRootID))
	}
	if !newRoot.IsCntr() {
		panic(fmt.Sprintf("ID %d is not a container", newRootID))
	}

	if oldRootID != 0 {
		v, in = h.HircObjsMap.Load(oldRootID)
		if !in {
			panic(fmt.Sprintf("ID %d has no associated hierarchy.", oldRootID))
		}
		oldRoot := v.(HircObj)
		_, err = oldRoot.HircID()
		if err != nil {
			panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", oldRootID))
		}
		if !oldRoot.IsCntr() {
			panic(fmt.Sprintf("ID %d is not a container", oldRootID))
		}
		oldRoot.RemoveLeaf(leaf)
		if leaf.BaseParameter().DirectParentId != 0 {
			panic(fmt.Sprintf("RemoveLeaf contract break. %d parent ID is non zero after removing from %d", id, oldRootID))
		}
	}

	newRoot.AddLeaf(leaf)

	l := len(h.HircObjs)
	h.HircObjs = slices.DeleteFunc(h.HircObjs, func(h HircObj) bool {
		tid, err := h.HircID()
		if err != nil {
			return false
		}
		return tid == id
	})
	if l <= len(h.HircObjs) {
		panic(fmt.Sprintf("%d does not exists in the HIRC", id))
	}

	newRootIdx := slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
		id, err := h.HircID()
		if err != nil {
			return false
		}
		return id == newRootID
	})
	if newRootIdx == -1 {
		panic(fmt.Sprintf("%d does not exists in the HIRC", newRootID))
	}

	l = len(h.HircObjs)
	h.HircObjs = slices.Insert(h.HircObjs, newRootIdx, leaf)
	if l >= len(h.HircObjs) {
		panic(fmt.Sprintf("%d does not added back to the HIRC", id))
	}

	if leaf.BaseParameter().DirectParentId != newRootID {
		panic(fmt.Sprintf("AddLeaf contract break. %d parent ID is non zero after attaching to %d", id, newRootID))
	}
}

func (h *HIRC) RemoveRoot(id, oldRootID uint32) {
	v, in := h.HircObjsMap.Load(id)
	if !in {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", id))
	}
	leaf := v.(HircObj)
	if b := leaf.BaseParameter(); b == nil {
		panic(fmt.Sprintf("%d is not containable", id))
	}
	if leaf.HircType() != HircTypeSound {
		slog.Warn("Parent rewiring is only supported for Sound right now")
		return
	}
	_, err := leaf.HircID()
	if err != nil {
		panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", id))
	}

	if oldRootID != 0 {
		v, in = h.HircObjsMap.Load(oldRootID)
		if !in {
			panic(fmt.Sprintf("ID %d has no associated hierarchy.", oldRootID))
		}
		oldRoot := v.(HircObj)
		_, err = oldRoot.HircID()
		if err != nil {
			panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", oldRootID))
		}
		if !oldRoot.IsCntr() {
			panic(fmt.Sprintf("ID %d is not a container", oldRootID))
		}

		oldRoot.RemoveLeaf(leaf)
		l := len(h.HircObjs)
		h.HircObjs = slices.DeleteFunc(h.HircObjs, func(h HircObj) bool {
			tid, err := h.HircID()
			if err != nil {
				return false
			}
			return tid == id
		})
		if l <= len(h.HircObjs) {
			panic(fmt.Sprintf("%d does not exists in the HIRC", id))
		}

		idx := slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
			return h.HircType() == HircTypeAction
		})
		l = len(h.HircObjs)
		if idx != -1 {
			h.HircObjs = slices.Insert(h.HircObjs, idx, leaf)
		} else {
			h.HircObjs = append(h.HircObjs, leaf)
		}
		if l >= len(h.HircObjs) {
			panic(fmt.Sprintf("%d does not added back to the HIRC", id))
		}
	}

	if leaf.BaseParameter().DirectParentId != 0 {
		panic(fmt.Sprintf("RemoveLeaf contract break. %d parent ID is non zero after removing from %d", id, oldRootID))
	}
}

func (h *HIRC) TreeArrIdx(tid uint32) int {
	return slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
		if id, err := h.HircID(); err != nil {
			return false
		} else {
			return id == tid
		}
	})
}

type HircObj interface {
	Encode() []byte
	BaseParameter() *BaseParameter
	HircID() (uint32, error)
	HircType() HircType
	IsCntr() bool
	NumLeaf() int
	ParentID() uint32
	// Modify DirectParentId,
	// pre condition: o.DirectParentId == 0
	// post condition: o.DirectParentId == HircObj.HircID
	AddLeaf(o HircObj)
	// Modify DirectParentId,
	// pre condition: o.DirectParentId == HircObj.HircID
	// post condition: DirectParentId = 0
	RemoveLeaf(o HircObj)
}

const SizeOfHircObjHeader = 1 + 4

type HircObjHeader struct {
	Type HircType // U8x
	Size uint32   // U32
}

type Unknown struct {
	HircObj
	Header *HircObjHeader
	Data   []byte
}

func NewUnknown(t HircType, s uint32, b []byte) *Unknown {
	return &Unknown{
		Header: &HircObjHeader{Type: t, Size: s},
		Data:   b,
	}
}

func (u *Unknown) Encode() []byte {
	assert.Equal(
		u.Header.Size,
		uint32(len(u.Data)),
		"Header size does not equal to actual data size",
	)

	bw := wio.NewWriter(uint64(SizeOfHircObjHeader + len(u.Data)))

	/* Header */
	bw.Append(u.Header)
	bw.AppendBytes(u.Data)

	return bw.Bytes()
}

func (u *Unknown) BaseParameter() *BaseParameter { return nil }

func (u *Unknown) HircID() (uint32, error) {
	return 0, fmt.Errorf("Hierarchy object type %d has yet implement GetHircID.", u.Header.Type)
}

func (u *Unknown) HircType() HircType { return u.Header.Type }

func (u *Unknown) IsCntr() bool { return false }

func (u *Unknown) NumLeaf() int { return 0 }

func (u *Unknown) ParentID() uint32 { return 0 }

func (u *Unknown) AddLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling AddLeaf.")
}

func (u *Unknown) RemoveLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling RemoveLeaf.")
}
