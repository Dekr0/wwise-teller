package wwise

import (
	"fmt"
	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type Event struct {
	Id         u32
	EventData  EventData
}

type EventData struct {
	NumActionIds   uio.V128
	ActionIds    []uint32
}

type EventComponent struct {
	EventData map[u32]EventData
}

// --- allocation / freeing --- //

func AllocEventComponent(numEvent u32) *EventComponent {
	if numEvent <= 0 {
		return &EventComponent{
			EventData: make(map[u32]EventData),
		}
	}
	return &EventComponent{
		EventData: make(map[u32]EventData, numEvent),
	}
}

func AllocEventData(numActionIds uio.V128) (data *EventData) {
	return &EventData{
		NumActionIds: numActionIds,
		ActionIds: make([]uint32, numActionIds.V, numActionIds.V),
	}
}

// --- assertion --- //

// Has no side effect
func AssertEvent(e Event) error {
	eventData := e.EventData
	if eventData.NumActionIds.V != u64(len(eventData.ActionIds)) {
		return fmt.Errorf("# of action ids does not equal to actual # of stored action ids")
	}
	return nil
}

// --- sizing --- //

// Has no side effect
func SizeOfEvent(e Event) (size u32) {
	size = SizeOfHierarchyId
	eventData := e.EventData
	b, v := eventData.NumActionIds.B, eventData.NumActionIds.V
	size += u32(len(b))
	size += u32(v) * Size32
	return size
}

// --- encoding --- //

// Has no side effect
func EncodeEvent(
	eCtx  *HircEncoderCtx,
	event Event,
) error {
	var err error
	id := event.Id
	size := SizeOfEvent(event)
	eventData := event.EventData
	header := HierarchyHeader{ HircTypeEvent, size }
	if err = eCtx.Struct(header, SizeOfHierarchyHeader); err != nil {
		return fmt.Errorf("Failed to encode hierarchy header: %w", err)
	}
	curr := eCtx.Encoder.Count
	if err = eCtx.Primitive(id); err != nil {
		return fmt.Errorf("Failed to encode id: %w", err)
	}
	err = eCtx.Bytes(eventData.NumActionIds.B)
	if err != nil {
		return fmt.Errorf("Failed to write number of action ids: %w", err)
	}
	actionIds := eventData.ActionIds
	for i, actionId := range actionIds {
		if err = eCtx.Primitive(actionId); err != nil {
			return fmt.Errorf("Failed to encode %d-th action id %d", i, actionId)
		}
	}
	return eCtx.Expect(curr, size)
}

// --- component getter and setter --- //

// Has no side effect
func (e *EventComponent) HasEventData(internalId u32) (in bool) {
	_, in = e.EventData[internalId]
	return in 
}

// Has no side effect
func (e *EventComponent) GetEventData(internalId u32) (data EventData) {
	data, in := e.EventData[internalId]
	if !in {
		panic("Failed to locate event data")
	}
	return data
}

// Has side effect
func (e *EventComponent) AddEventData(internalId u32, data EventData) {
	if _, in := e.EventData[internalId]; in {
		panic(MonotonicIdCollision)
	}
	e.EventData[internalId] = data
}

// --- HIRC component wrapper --- //

func (h *HIRC) Event(internalId u32, inOut *Event) {
	inOut.Id = h.Hierarchy.GetHierarchyNode(internalId).Id
	inOut.EventData = h.EventComponet.GetEventData(internalId)
}
