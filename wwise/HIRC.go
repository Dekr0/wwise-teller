package wwise

import (
	"context"
	"io"
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

// Has side effect
func NewHierarchy(h *HIRC, id u32, t HircType) (internalId u32) {
	internalId = h.monoId
	if _, in := h.Hierarchies[internalId]; in {
		panic("Implementation error of monotonic internal id: duplication detected")
	}
	h.Hierarchies[h.monoId] = &Hierarchy{ id, t }
	h.monoId++
	return internalId
}

// Has side effect
func NewState(h *HIRC, id u32, data *StateProp) {
	if data == nil {
		panic("State property is nil")
	}
	internalId := NewHierarchy(h, id, HircTypeState)
	s := &h.StateComponent
	if _, in := s.StateProps[internalId]; in {
		panic("Implementation error of monotonic internal id: duplication detected")
	}
	s.StateProps[internalId] = data
}

// Has side effect
func NewEvent(h *HIRC, id u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	internalId := NewHierarchy(h, id, HircTypeEvent)
	e := &h.EventComponet
	if _, in := e.EventData[internalId]; in {
		panic("Implementation error of monotonic internal id: duplication detected")
	}
	e.EventData[internalId] = data
}

func EncodeHirc(
	ctx      context.Context,
	w        io.Writer,
	o        order,
	version  u32,
	h       *HIRC, 
	opt     *EncodeHircOpt,
) (err error) {
	return nil
}
