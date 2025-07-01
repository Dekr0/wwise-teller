package wwise

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/Dekr0/wwise-teller/wio"
)

const MaxEncodeRoutine = 6

// # of hierarchy object (uint32)
const SizeOfHIRCHeader = 4

type HIRC struct {
	I uint8
	T []byte

	// Retain the tree structure that comes from the decoding with minimal
	// modification
	HircObjs []HircObj

	// 1
	AudioDevices sync.Map

	// Master Mixer Hierarchy 8
	Buses             sync.Map
	BusRoots        []*BusHircNode

	Attenuations      sync.Map
	FxShareSets       sync.Map // Share the same ID with Aux Bus that use this FX
	FxCustoms         sync.Map // Share the same ID with Aux Bus that use this FX
	AuxBuses          sync.Map
	LFOModulators     sync.Map
	EnvelopeModulator sync.Map
	TimeModulator     sync.Map

	// Sound
	// Random / Sequence Container
	// Switch Container
	// Actor Mixer
	// Layer Container
	// Dialogue Event?
	ActorMixerHirc          sync.Map // 6
	ActorMixerRoots         []*ActorMixerHircNode
	ActorMixerHircNodesMap  map[uint32]*ActorMixerHircNode

	// Game Sync 3
	Actions sync.Map
	Events  sync.Map
	States  sync.Map

	// Music Segment
	// Music Track
	// Music Switch Container
	// Music Random Sequence Container
	MusicHirc        sync.Map // 4
	MusicHircRoots   []MusicHircNode
}

func NewHIRC(I uint8, T []byte, numHircItem uint32) *HIRC {
	return &HIRC{
		I:        I,
		T:        T,
		HircObjs: make([]HircObj, numHircItem),
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
		case <- ctx.Done():
			return nil, ctx.Err()
		case r := <- c:
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

type HircCount struct {
	Attenuations    uint16
	ActorMixerHircs uint16
	ActorMixerRoots uint16
	Buses           uint16
	FxS             uint16
	Events          uint16
	Modulators      uint16
	MusicHircs      uint16
	MusicHircRoots  uint16
	States          uint16
}

// Use for deterministic allocation initially of different types of filter array
func (h *HIRC) HierarchyCount() HircCount {
	c := HircCount{}
	for _, o := range h.HircObjs {
		if ActorMixerHircType(o) {
			c.ActorMixerHircs += 1
			if ContainerActorMixerHircType(o) {
				c.ActorMixerRoots += 1
			}
		} else if MusicHircType(o) {
			c.MusicHircs += 1
			if ContainerMusicHircType(o) {
				c.MusicHircRoots += 1
			}
		} else if BusHircType(o) {
			c.Buses += 1
		} else if FxHircType(o) {
			c.FxS += 1
		} else if ModulatorType(o) {
			c.Modulators += 1
		} else {
			switch o.HircType() {
			case HircTypeAttenuation:
				c.Attenuations += 1
			case HircTypeEvent:
				c.Events += 1
			case HircTypeState:
				c.States += 1
			}
		}
	}
	return c
}

func (h *HIRC) ChangeRoot(id, newRootID, oldRootID uint32, syncUITree bool) {
	if newRootID == 0 {
		h.RemoveRoot(id, oldRootID, syncUITree)
		return
	}

	// Use this for now because using Map will have issue due to how Wwise 
	// bundle hierarchy objects
	idx := slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
		i, err := h.HircID()
		if err != nil {
			return false
		}
		return i == id
	})
	if idx == -1 {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", id))
	}
	leaf := h.HircObjs[idx]
	if b := leaf.BaseParameter(); b == nil {
		panic(fmt.Sprintf("%d is not containable", id))
	}
	_, err := leaf.HircID()
	if err != nil {
		panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", id))
	}

	idx = slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
		i, err := h.HircID()
		if err != nil {
			return false
		}
		return i == newRootID
	})
	if idx == -1 {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", newRootID))
	}
	newRoot := h.HircObjs[idx]
	_, err = newRoot.HircID()
	if err != nil {
		panic(fmt.Sprintf("ID %d has an associated hiearchy but its HircID interface returns error.", newRootID))
	}
	if !newRoot.IsCntr() {
		panic(fmt.Sprintf("ID %d is not a container", newRootID))
	}

	if oldRootID != 0 {
		idx = slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
			i, err := h.HircID()
			if err != nil {
				return false
			}
			return i == oldRootID
		})
		if idx == -1 {
			panic(fmt.Sprintf("ID %d has no associated hierarchy.", oldRootID))
		}
		oldRoot := h.HircObjs[idx]
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

	if syncUITree {
		h.BuildTree()
	}
}

func (h *HIRC) RemoveRoot(id, oldRootID uint32, syncUITree bool) {
	idx := slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
		i, err := h.HircID()
		if err != nil {
			return false
		}
		return i == id
	})
	if idx == -1 {
		panic(fmt.Sprintf("ID %d has no associated hierarchy.", id))
	}
	leaf := h.HircObjs[idx]
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
		idx := slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
			i, err := h.HircID()
			if err != nil {
				return false
			}
			return i == oldRootID
		})
		if idx == -1 {
			panic(fmt.Sprintf("ID %d has no associated hierarchy.", oldRootID))
		}
		oldRoot := h.HircObjs[idx]
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

		// Find the first action hierarchy, and put it next to that action
		idx = slices.IndexFunc(h.HircObjs, func(h HircObj) bool {
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

// Prototyping
func (h *HIRC) AppendNewSoundToRanSeqContainer(s *Sound, rsId uint32, syncUITree bool) error {
	if s.BaseParam.DirectParentId != 0 {
		return fmt.Errorf("Sound %d already has a parent", s.Id)
	}
	idx := slices.IndexFunc(h.HircObjs, func(o HircObj) bool {
		id, err := o.HircID()
		if err != nil {
			return false
		}
		return rsId == id
	})
	if idx == -1 {
		return fmt.Errorf("No random / sequence container has ID of %d", rsId)
	}
	r := h.HircObjs[idx].(*RanSeqCntr)
	r.AddLeaf(s)
	r.AddLeafToPlayList(len(r.Container.Children) - 1)
	h.HircObjs = slices.Insert(h.HircObjs, idx, HircObj(s))
	_, in := h.ActorMixerHirc.LoadOrStore(s.Id, s)
	if in {
		panic(fmt.Sprintf("Sound object %d already exist!", s.Id))
	}
	if syncUITree {
		h.BuildTree()
	}
	return nil
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

// Prototyping
func (h *HIRC) AppendNewRanSeqCntrToActorMixer(r *RanSeqCntr, actorId uint32, syncUITree bool) error {
	if r.BaseParam.DirectParentId != 0 {
		return fmt.Errorf("Random / Sequence Container %d already has a parent", r.Id)
	}
	idx := slices.IndexFunc(h.HircObjs, func(o HircObj) bool {
		id, err := o.HircID()
		if err != nil {
			return false
		}
		return actorId == id
	})
	if idx == -1 {
		return fmt.Errorf("No Actor Mixer has ID of %d", actorId)
	}
	mixer := h.HircObjs[idx].(*ActorMixer)
	mixer.AddLeaf(r)
	h.HircObjs = slices.Insert(h.HircObjs, idx, HircObj(r))
	_, in := h.ActorMixerHirc.LoadOrStore(r.Id, r)
	if in {
		panic(fmt.Sprintf("Randome / Sequence object %d already exist!", r.Id))
	}
	if syncUITree {
		h.BuildTree()
	}
	return nil
}

// Prototyping
func (h *HIRC) AppendNewActionToEvent(a *Action, eventID uint32) error {
	idx := slices.IndexFunc(h.HircObjs, func(o HircObj) bool {
		id, err := o.HircID()
		if err != nil {
			return false
		}
		return eventID == id
	})
	if idx == -1 {
		return fmt.Errorf("No Event has ID of %d", eventID)
	}
	event := h.HircObjs[idx].(*Event)
	event.NewAction(a.Id)
	h.HircObjs = slices.Insert(h.HircObjs, idx, HircObj(a))
	_, in := h.ActorMixerHirc.LoadOrStore(a.Id, a)
	if in {
		panic(fmt.Sprintf("Action object %d already exist!", a.Id))
	}
	return nil
}
