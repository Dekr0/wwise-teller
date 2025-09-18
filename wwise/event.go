package wwise

import (
	"fmt"
	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type EventH struct {
	Id         u32
	EventData *EventData
}

type EventData struct {
	NumActionIds   uio.V128
	ActionIds    []uint32
}

type EventComponent struct {
	EventData map[u32]*EventData
}

// --- allocation / freeing --- //

func AllocEventComponent(numEvent u32) *EventComponent {
	if numEvent <= 0 {
		return &EventComponent{
			EventData: make(map[u32]*EventData),
		}
	}
	return &EventComponent{
		EventData: make(map[u32]*EventData, numEvent),
	}
}

func AllocEventData(numActionIds *uio.V128) (data *EventData) {
	return &EventData{
		NumActionIds: *numActionIds,
		ActionIds: make([]uint32, numActionIds.V, numActionIds.V),
	}
}

// --- assertion --- //

// Has no side effect
func AssertEvent(eventData *EventData) error {
	if eventData.NumActionIds.V != u64(len(eventData.ActionIds)) {
		return fmt.Errorf("# of action ids does not equal to actual # of stored action ids")
	}
	return nil
}

// --- sizing --- //

// Has no side effect
func SizeOfEvent(eventData *EventData) (size u32) {
	size = SizeOfHierarchyId
	b, v := eventData.NumActionIds.B, eventData.NumActionIds.V
	size += u32(len(b))
	size += u32(v) * Size32
	return size
}

// --- encoding --- //

// Has no side effect
func EncodeEvent(
	eCtx  *HircEncoderCtx,
	event *EventH,
	size   u32,
) {
	var err error
	id := event.Id
	eventData := event.EventData
	header := HierarchyHeader{ HircTypeEvent, size }
	if err = eCtx.Struct(header, SizeOfHierarchyHeader); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode hierarchy header: %w", id, err))
	}
	curr := eCtx.Encoder.Count
	if err = eCtx.Primitive(id); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode id: %w", id, err))
	}
	err = eCtx.Bytes(eventData.NumActionIds.B)
	if err != nil {
		panic(fmt.Errorf("(Event %d) Failed to write number of action ids: %w", id, err))
	}
	actionIds := eventData.ActionIds
	for i, actionId := range actionIds {
		if err = eCtx.Primitive(actionId); err != nil {
			panic(fmt.Errorf("(Event %d) Failed to encode %d-th action id %d", id, i, actionId))
		}
	}
	if err := eCtx.Expect(curr, size); err != nil {
		panic(fmt.Errorf("(Event %d) %w", id, err))
	}
}

// --- component assertion wrapper --- //

func (e *EventComponent) AssertEventById(internalId u32) error {
	eventData, in := e.EventData[internalId]
	if !in {
		return fmt.Errorf("Failed to locate event data")
	}
	return AssertEvent(eventData)
}

// --- component sizing wrapper --- //

func (e *EventComponent) SizeOfEventById(internalId u32) u32 {
	eventData := e.GetEventData(internalId)
	return SizeOfEvent(eventData)
}

// --- component getter and setter --- //

// Has no side effect
func (e *EventComponent) HasEventData(internalId u32) (in bool) {
	_, in = e.EventData[internalId]
	return in 
}

// Has no side effect
func (e *EventComponent) GetEventData(internalId u32) (data *EventData) {
	data, in := e.EventData[internalId]
	if !in {
		panic("Failed to locate event data")
	}
	return data
}

// Has side effect
func (e *EventComponent) AddEventData(internalId u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	if _, in := e.EventData[internalId]; in {
		panic(MonotonicIdCollision)
	}
	e.EventData[internalId] = data
}

// --- HIRC component wrapper --- //

func (h *HIRC) GatherEventData(internalId u32) *EventH {
	node := h.Hierarchy.GetHierarchyNode(internalId)
	data := h.EventComponet.GetEventData(internalId)
	return &EventH{ node.Id, data }
}
