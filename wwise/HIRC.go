package wwise

import (
	"slices"
	"sync"
)


type EncodeHircOpt struct {
	NumRoutine u8
}

type HierarchyHeader struct {
	Type HircType
	Size u32
}

type Hierarchy struct {
	Id   u32
	Type HircType
}

type HIRC struct {
	monoId u32 // a monotonic id counter that only increase

	// This is intended to be used in decoding phase
	hierarchyDMu sync.Mutex
	InternalIds  []u32
	Hierarchies  map[u32]*Hierarchy

	EventComponet  EventComponent
	StateComponent StateComponent

	// This is intended to be used in decoding phase
	encodedHierarchyDMu sync.Mutex
	EncodedHierarchy    map[u32][]byte
}

func AllocHIRC(numHirc u32) *HIRC {
	return &HIRC{
		monoId: 0,
		Hierarchies: make(map[u32]*Hierarchy, numHirc),
		EventComponet: EventComponent{
			// Estimate 25% of hierarchies will be Event
			EventData: make(map[u32]*EventData, numHirc / 4),
		},
		StateComponent: StateComponent{
			// TODO: Estimation
			StateProps: make(map[u32]*StateProps),
		},
		EncodedHierarchy: make(map[u32][]byte),
	}
}

// Has side effect
func AddHierarchy(h *HIRC, id u32, t HircType) (internalId u32) {
	h.hierarchyDMu.Lock()
	defer h.hierarchyDMu.Unlock()

	internalId = h.monoId

	if _, in := h.Hierarchies[internalId]; in {
		panic(MonotonicIdCollision)
	}
	// Get rid off this once I figure out the tree traversal algorithm 
	if slices.Contains(h.InternalIds, internalId) {
		panic(MonotonicIdCollision)
	}

	h.InternalIds = append(h.InternalIds, internalId)
	h.Hierarchies[internalId] = &Hierarchy{ id, t }

	h.monoId++

	return internalId
}

// Has side effect
func AddState(h *HIRC, id u32, data *StateProps) {
	if data == nil {
		panic("State property is nil")
	}

	internalId := AddHierarchy(h, id, HircTypeState)
	AddStateData(&h.StateComponent, internalId, data)
}

// Has side effect
func AddEvent(h *HIRC, id u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	internalId := AddHierarchy(h, id, HircTypeEvent)

	AddEventData(&h.EventComponet, internalId, data)
}

// Has side effect
func AddEncodedHierarchy(h *HIRC, id u32, t HircType, encoded []byte) {
	if encoded == nil {
		panic("Encoded hierarchy data is nil")
	}

	internalId := AddHierarchy(h, id, t)

	h.encodedHierarchyDMu.Lock()
	defer h.encodedHierarchyDMu.Unlock()
	if _, in := h.EncodedHierarchy[internalId]; in {
		panic(MonotonicIdCollision)
	}
	h.EncodedHierarchy[internalId] = encoded
}
