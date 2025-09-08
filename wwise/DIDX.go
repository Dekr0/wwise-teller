package wwise

import (
	"fmt"
	"slices"
	"sync"
)

const SizeOfMediaIndex = 12

type MediaIndexEntry struct {
	SourceId u32
	Offset   u32
	Size     u32
}

type DIDXDATA struct {
	mu sync.Mutex
	
	SourceIds []u32
	Offsets   map[u32]u32
	Sizes     map[u32]u32
	AudioData map[u32][]byte
}

// Only use this for DecodeDATA!
func LockDIDXDATA(d *DIDXDATA) {
	d.mu.Lock()
}

// Only use this for DecodeDATA!
func UnlockDIDXDATA(d *DIDXDATA) {
	d.mu.Unlock()
}

func NewDIDXDATA(size u32) *DIDXDATA {
	return &DIDXDATA{
		SourceIds: make([]u32, 0, size),
		Offsets: make(map[u32]u32, size),
		Sizes: make(map[u32]u32, size),
	}
}

// Has side effect
// Thread safe
// Use this if assuming there will be no duplicate in DIDX entry (e.g., at 
// decoding phase)
func NewMediaIndex(d *DIDXDATA, m MediaIndexEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	d.SourceIds = append(d.SourceIds, m.SourceId) 
	d.Offsets[sourceId] = offset
	d.Sizes[sourceId] = size
}

// Has side effect
// Thread safe
func NewMediaIndexCheck(d *DIDXDATA, m MediaIndexEntry) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	if !slices.Contains(d.SourceIds, sourceId) {
		return fmt.Errorf("Media index with %d already exist.", sourceId)
	}

	d.SourceIds = append(d.SourceIds, m.SourceId)

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
func NumMediaIndex(d *DIDXDATA) u32 {
	d.mu.Lock()
	defer d.mu.Unlock()

	return u32(len(d.SourceIds))
}

// No side effect
// Thread safe
func HasMediaIndex(d *DIDXDATA, sourceId u32) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return slices.Contains(d.SourceIds, sourceId)
}

// No side effect
// Thread safe
// Use HasMediaIndex before MediaIndex
func MediaIndex(d *DIDXDATA, sourceId u32) (offset u32, size u32) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !slices.Contains(d.SourceIds, sourceId) {
		panic(fmt.Sprintf("No media index with %d.", sourceId))
	}

	offset, in := d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("%d has meida index but it has no offset value", sourceId))
	}

	size, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("%d has meida index but it has no size value", sourceId))
	}

	return offset, size
}

// No side effect
// Thread safe
func MediaIndexCheck(d *DIDXDATA, sourceId u32) (offset u32, size u32, in bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !slices.Contains(d.SourceIds, sourceId) {
		return offset, size, false
	}

	offset, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("%d has meida index but it has no offset value", sourceId))
	}

	size, in = d.Offsets[sourceId]
	if !in {
		panic(fmt.Sprintf("%d has meida index but it has no size value", sourceId))
	}

	return offset, size, in
}

// Has side effect
// Thread safe
// Use this when omiting all alignment at the decoding phase, or use it with 
// HasMediaIndex
func UpdateMediaIndex(d *DIDXDATA, m MediaIndexEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	if !slices.Contains(d.SourceIds, sourceId) {
		panic(fmt.Sprintf("Source id %d has no media index", sourceId))
	}

	if _, in := d.Offsets[sourceId]; !in {
		panic(fmt.Sprintf("Source id %d does not have offset value", sourceId))
	}

	if _, in := d.Sizes[sourceId]; !in {
		panic(fmt.Sprintf("Source id %d does not have size value", sourceId))
	}

	d.Offsets[sourceId] = offset
	d.Sizes[sourceId] = size
}
