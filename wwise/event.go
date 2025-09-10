package wwise

import (
	bin "encoding/binary"
	"fmt"
	"io"
	"sync"

	uio "github.com/Dekr0/unwise/io"
)

type EventData struct {
	NumActionIds   uio.V128
	ActionIds    []uint32
}

type EventComponent struct {
	dMu sync.Mutex

	EventData map[u32]*EventData
}

func NewEventData(e *EventComponent, internalId u32, data *EventData) {
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
	w           io.Writer,
	o           order,
	e          *EventComponent,
	version     u32,
	internalId  u32, 
	id          u32,
) {
	AssertEvent(e, version, internalId, id)

	var err error

	pos := u32(0)

	size := SizeOfEvent(e, version, internalId, id)
	eventData, in := e.EventData[internalId]
	if !in {
		panic(fmt.Sprintf("Event %d does not have event data", id))
	}

	header := HierarchyHeader{ HircTypeEvent, size }
	if err = bin.Write(w, o, header); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode hierarchy header: %w", id, err))
	}

	if err = bin.Write(w, o, id); err != nil {
		panic(fmt.Errorf("(Event %d) Failed to encode id: %w", id, err))
	}
	pos += SizeOfHierarchyId

	b := eventData.NumActionIds.B
	n, err := w.Write(b)
	if err != nil {
		panic(fmt.Errorf("(Event %d) Failed to write number of action ids: %w", id, err))
	}
	if n != len(b) {
		panic(fmt.Errorf("(Event %d) # of bytes write does not equal to # of bytes stored", id))
	}
	pos += u32(len(b))

	actionIds := eventData.ActionIds
	for i, actionId := range actionIds {
		if err = bin.Write(w, o, actionId); err != nil {
			panic(fmt.Errorf("(Event %d) Failed to encode %d-th action id %d", id, i, actionId))
		}
		pos += Size32
	}

	if pos != size {
		panic(fmt.Errorf("(Event %d) Expect encoded data has a size of %d but receive %d", id, size, pos))
	}
}
