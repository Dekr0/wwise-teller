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

func AddNewMediaIndex(d *DIDX, sourceId u32, offset u32, size u32) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; in {
		return fmt.Errorf("Media index with %d already exist.", sourceId)
	}

	d.SourceIds[sourceId] = struct{}{}

	if _, in := d.Offset[sourceId]; in {
		panic(fmt.Sprintf("Media index with %d does not exist but it has offset value", sourceId))
	}
	d.Offset[sourceId] = offset

	if _, in := d.Size[sourceId]; in {
		panic(fmt.Sprintf("Media index with %d does not exist but it has size value", sourceId))
	}
	d.Size[sourceId] = size

	return nil
}

func HasSource(d *DIDX, sourceId u32) (in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, in = d.SourceIds[sourceId]
	return in
}

// Use HasSource before MediaIndex
func MediaIndex(d *DIDX, sourceId u32) (offset u32, size u32) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; !in {
		panic(fmt.Sprintf("No media index with %d.", sourceId))
	}

	offset, in := d.Offset[sourceId]
	if !in {
		panic(fmt.Sprintf("No offset value associated with source id %d", sourceId))
	}

	size, in = d.Offset[sourceId]
	if !in {
		panic(fmt.Sprintf("No size value associated with source id %d", sourceId))
	}

	return offset, size
}

func MediaIndexCheck(d *DIDX, sourceId u32) (offset u32, size u32, in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; !in {
		return offset, size, in
	}

	offset, in = d.Offset[sourceId]
	if !in {
		panic(fmt.Sprintf("No offset value associated with source id %d", sourceId))
	}

	size, in = d.Offset[sourceId]
	if !in {
		panic(fmt.Sprintf("No size value associated with source id %d", sourceId))
	}

	return offset, size, in
}
