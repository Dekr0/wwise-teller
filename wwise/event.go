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

// --- function --- //

// Has side effect
func AddEventData(e *EventComponent, internalId u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	if _, in := e.EventData[internalId]; in {
		panic(MonotonicIdCollision)
	}
	e.EventData[internalId] = data
}

// Has no side effect
func AssertEvent(e *EventComponent, internalId u32, id u32) {
	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}
	if eventData.NumActionIds.V != u64(len(eventData.ActionIds)) {
		panic(fmt.Sprintf("Event %d: # of action ids does not equal to actual # of stored action ids", id))
	}
}

// Has no side effect
func SizeOfEvent(e *EventComponent, internalId u32, id u32) (size u32) {
	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}
	size = SizeOfHierarchyId
	b, v := eventData.NumActionIds.B, eventData.NumActionIds.V
	size += u32(len(b))
	size += u32(v) * Size32
	return size
}

// Has no side effect
func EncodeEvent(
	eCtx       *HircEncoderCtx,
	e          *EventComponent,
	internalId  u32, 
	id          u32,
) {
	var err error
	size := SizeOfEvent(e, internalId, id)
	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}
	header := HierarchyHeader{ HircTypeEvent, size }
	if err = uio.EncodeStruct(eCtx.Encoder, header, SizeOfHierarchyHeader); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode hierarchy header: %w", id, err))
	}
	curr := eCtx.Encoder.Count
	if err = HIRCEncode(eCtx, id); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode id: %w", id, err))
	}
	err = HIRCEncodeBytes(eCtx, eventData.NumActionIds.B)
	if err != nil {
		panic(fmt.Errorf("(Event %d) Failed to write number of action ids: %w", id, err))
	}
	actionIds := eventData.ActionIds
	for i, actionId := range actionIds {
		if err = HIRCEncode(eCtx, actionId); err != nil {
			panic(fmt.Errorf("(Event %d) Failed to encode %d-th action id %d", id, i, actionId))
		}
	}
	if err := HIRCAssertEncodeLimit(eCtx, curr, size); err != nil {
		panic(fmt.Errorf("(Event %d) %w", id, err))
	}
}
