package wwise

import (
	"sync"

	"github.com/Dekr0/unwise/io"
)

type EventData struct {
	NumActionIds   io.V128
	ActionIds    []uint32
}

type EventComponet struct {
	mu sync.Mutex

	EventData map[u32]*EventData
}
