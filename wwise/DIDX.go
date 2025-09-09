package wwise

import (
	bin "encoding/binary"
	"fmt"
	"io"
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

// Has side effect
// Thread safe
func ComputeDIDXOffset(d *DIDXDATA) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceIds := d.SourceIds
	offsets := d.Offsets
	sizes := d.Sizes
	
	offset := u64(0)
	for _, sourceId := range sourceIds {
		_, in := offsets[sourceId]
		if !in {
			panic(fmt.Sprintf("%d does not have an offset value", offset))
		}

		size, in := sizes[sourceId]
		if !in {
			panic(fmt.Sprintf("%d does not have a size value", size))
		}

		offsets[sourceId] = u32(offset)
		offset += u64(size)
	}
}

func VerifyDIDXDATA(d *DIDXDATA) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceIds := d.SourceIds
	audioDataIndices := d.AudioData
	offsets := d.Offsets
	sizes := d.Sizes

	if len(sourceIds) != len(audioDataIndices) {
		return fmt.Errorf(
			"# of source ids (%d) does not equal to # of audio data (%d)", 
			len(sourceIds), len(audioDataIndices),
		)
	}

	if len(audioDataIndices) != len(offsets) {
		return fmt.Errorf(
			"# of offset values (%d) does not equal to # of source ids (%d)",
			len(offsets), len(sourceIds),
		)
	}

	if len(offsets) != len(sizes) {
		return fmt.Errorf(
			"# of size values (%d) does not equal to # of source ids (%d)",
			len(sizes), len(sourceIds),
		)
	}

	offsetChecker := u64(0)
	for i, sourceId := range d.SourceIds {
		audioData, in := audioDataIndices[sourceId]
		if !in {
			return fmt.Errorf(
				"Source id %d does not have an associated audio data",
				sourceId,
			)
		}

		offset, in := offsets[sourceId]
		if !in {
			return fmt.Errorf(
				"Source id %d does not have an associated offset value",
				sourceId,
			)
		}

		size, in := sizes[sourceId]
		if !in {
			return fmt.Errorf(
				"Source id %d does not have an associated size value",
				sourceId,
			)
		}

		if u64(offset) != offsetChecker {
			return fmt.Errorf(
				"Expecting media index (index %d) with source id %d has an offset of %d but receive %d",
				i, sourceId, offsetChecker, offset,
			)
		}

		audioDataSize := len(audioData)
		if int(size) != audioDataSize {
			return fmt.Errorf(
				"Media index (index %d) with source id %d has a size value of %d but its audio data has a size value of %d",
				i, sourceId, size, audioDataSize,
			)
		}
	}

	return nil
}

func EncodeDIDX(d *DIDXDATA, w io.Writer, o order) (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	sourceIds := d.SourceIds
	offsets := d.Offsets
	sizes := d.Sizes

	var payload MediaIndexEntry
	for _, sourceId := range sourceIds {
		offset, in := offsets[sourceId]
		if !in {
			return fmt.Errorf("Source id %d does not have a offset value", sourceId)
		}

		size, in := sizes[sourceId]
		if !in {
			return fmt.Errorf("Source id %d does not have a size value", sourceId)
		}

		payload.SourceId = sourceId
		payload.Offset = offset
		payload.Size = size

		err = bin.Write(w, o, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func EncodeDATANotAlign(d *DIDXDATA, w io.Writer) (err error) {
	d.mu.Lock()
	d.mu.Unlock()

	sourceIds := d.SourceIds
	audioData := d.AudioData

	for _, sourceId := range sourceIds {
		audioData, in := audioData[sourceId]
		if !in {
			return fmt.Errorf("Source id %d does not have an associated audio data", sourceId)
		}
		if _, err := w.Write(audioData); err != nil {
			return err
		}
	}

	return err
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
