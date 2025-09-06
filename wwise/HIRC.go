package wwise

import (
	"sync"
)

type Hierarchy struct {
	Id   u32
	Type HircType
}

type HIRC struct {
	mu sync.Mutex

	monoId u32 // a monotonic id counter that only increase

	Hierarchies map[u32]*Hierarchy

	EventComponet  EventComponet
	StateComponent StateComponent
}

func NewHIRC(numHirc u32) *HIRC {
	return &HIRC{
		monoId: 0,
		Hierarchies: make(map[u32]*Hierarchy, numHirc),
		EventComponet: EventComponet{
			// Estimate 25% of hierarchies will be Event
			EventData: make(map[u32]*EventData, numHirc / 4),
		},
		StateComponent: StateComponent{
			// TODO: Estimation
			StateProps: make(map[u32]*StateProp),
		},
	}
}

// Thread safe setter and getter
func HIRCNewHierarchy(h *HIRC, id u32, t HircType) (internalId u32) {
	h.mu.Lock()
	internalId = h.monoId
	h.Hierarchies[h.monoId] = &Hierarchy{ id, t }
	h.monoId++
	h.mu.Unlock()
	return internalId
}

func HIRCNewState(h *HIRC, id u32, data *StateProp) {
	internalId := HIRCNewHierarchy(h, id, HircTypeState)

	s := &h.StateComponent
	s.mu.Lock()
	s.StateProps[internalId] = data
	s.mu.Unlock()
}

func HIRCNewEvent(h *HIRC, id u32, data *EventData) {
	internalId := HIRCNewHierarchy(h, id, HircTypeEvent)

	e := &h.EventComponet
	e.mu.Lock()
	e.EventData[internalId] = data
	e.mu.Unlock()
}
