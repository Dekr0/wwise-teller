package wwise

import (
	"fmt"
	"sync"
)

const SizeOfMediaIndex = 12

type MediaIndexEntry struct {
	SourceId u32
	Offset   u32
	Size     u32
}

type DIDX struct {
	mu sync.Mutex
	
	SourceIds map[u32]struct{}
	Offsets   map[u32]u32
	Sizes     map[u32]u32
}

func NewDIDX(size u32) *DIDX {
	return &DIDX{
		SourceIds: make(map[u32]struct{}, size),
		Offsets: make(map[u32]u32, size),
		Sizes: make(map[u32]u32, size),
	}
}

// Has side effect
// Thread safe
// Use this if assuming there will be no duplicate in DIDX entry
func AddNewMediaIndex(d *DIDX, m MediaIndexEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	d.SourceIds[sourceId] = struct{}{}
	d.Offsets[sourceId] = offset
	d.Sizes[size] = size
}

// Has side effect
// Thread safe
func AddNewMediaIndexCheck(d *DIDX, m MediaIndexEntry) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	if _, in := d.SourceIds[sourceId]; in {
		return fmt.Errorf("Media index with %d already exist.", sourceId)
	}

	d.SourceIds[sourceId] = struct{}{}

	if _, in := d.Offsets[sourceId]; in {
		panic(fmt.Sprintf("Media index with %d does not exist but it has offset value", sourceId))
	}
	d.Offsets[sourceId] = offset

	if _, in := d.Sizes[sourceId]; in {
		panic(fmt.Sprintf("Media index with %d does not exist but it has size value", sourceId))
	}
	d.Sizes[sourceId] = size

	return nil
}

// No side effect
// Thread safe
func HasMediaIndex(d *DIDX, sourceId u32) (in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, in = d.SourceIds[sourceId]
	return in
}

// No side effect
// Thread safe
// Use HasMediaIndex before MediaIndex
func MediaIndex(d *DIDX, sourceId u32) (offset u32, size u32) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; !in {
		panic(fmt.Sprintf("No media index with %d.", sourceId))
	}

	offset, in := d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("No offset value associated with source id %d", sourceId))
	}

	size, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("No size value associated with source id %d", sourceId))
	}

	return offset, size
}

// No side effect
// Thread safe
func MediaIndexCheck(d *DIDX, sourceId u32) (offset u32, size u32, in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.SourceIds[sourceId]; !in {
		return offset, size, in
	}

	offset, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("No offset value associated with source id %d", sourceId))
	}

	size, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("No size value associated with source id %d", sourceId))
	}

	return offset, size, in
}
