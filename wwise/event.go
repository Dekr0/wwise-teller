package wwise

import (
	"fmt"
	"sync"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type EventData struct {
	NumActionIds   uio.V128
	ActionIds    []uint32
}

type EventComponent struct {
	dMu sync.Mutex

	EventData map[u32]*EventData
}

// --- allocation / freeing --- //

func AllocEventData(numActionIds *uio.V128) (data *EventData) {
	return &EventData{
		NumActionIds: *numActionIds,
		ActionIds: make([]uint32, numActionIds.V, numActionIds.V),
	}
}

// --- function --- //

func AddEventData(e *EventComponent, internalId u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}

	e.dMu.Lock()
	defer e.dMu.Unlock()
	if _, in := e.EventData[internalId]; in {
		panic(MonotonicIdCollision)
	}
	e.EventData[internalId] = data
}

func AssertEvent(e *EventComponent, version u32, internalId u32, id u32) {
	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}
	if eventData.NumActionIds.V != u64(len(eventData.ActionIds)) {
		panic(fmt.Sprintf("Event %d: # of action ids does not equal to actual # of stored action ids", id))
	}
}

func SizeOfEvent(e *EventComponent, version u32, internalId u32, id u32) (size u32) {
	size = SizeOfHierarchyId

	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}

	b, v := eventData.NumActionIds.B, eventData.NumActionIds.V

	size += u32(len(b))
	size += u32(v) * Size32

	return size
}

func EncodeEvent(
	eCtx       *HircEncoderCtx,
	e          *EventComponent,
	internalId  u32, 
	id          u32,
) {
	AssertEvent(e, eCtx.Version, internalId, id)

	var err error

	size := SizeOfEvent(e, eCtx.Version, internalId, id)
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
