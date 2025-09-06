package wwise

import (
	"fmt"
	"sync"
)

type DIDX struct {
	mu sync.Mutex
	
	SourceIds map[u32]struct{}
	Offset    map[u32]u32
	Size      map[u32]u32
}

func NewDIDX(size u32) *DIDX {
	return &DIDX{
		SourceIds: make(map[u32]struct{}, size),
		Offset: make(map[u32]u32, size),
		Size: make(map[u32]u32, size),
	}
}

// Write test to make sure offset and size syncs up with source id.
func AddNewMediaIndex(d *DIDX, sourceId u32, offset u32, size u32) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; in {
		return fmt.Errorf("Media index with %d already exist.", sourceId)
	}

	d.SourceIds[sourceId] = struct{}{}
	d.Offset[sourceId] = offset
	d.Size[sourceId] = size

	return nil
}

// Write test to make sure offset and size syncs up with source id.
func MediaIndex(d *DIDX, sourceId u32) (offset u32, size u32, in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in = d.SourceIds[sourceId]; !in {
		return 0, 0, in
	}

	offset, _ = d.Offset[sourceId]
	size, _ = d.Size[sourceId]

	return offset, size, in
}
