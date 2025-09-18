package wwise

import (
	"fmt"
	"slices"

	uio "github.com/Dekr0/unwise/io"
)

const SizeOfMediaIndex = 12

type MediaIndexS struct {
	SourceId u32
	Offset   u32
	Size     u32
}

type AudioStore struct {
	SourceIds []u32
	Offsets   map[u32]u32
	Sizes     map[u32]u32
	AudioData map[u32][]byte
}

// Has side effect
func ComputeDIDXOffset(d *AudioStore) {
	sourceIds := d.SourceIds
	offsets := d.Offsets
	sizes := d.Sizes
	
	offsetNotAlign := u64(0)
	for _, sourceId := range sourceIds {
		_, in := offsets[sourceId]
		if !in {
			panic(fmt.Sprintf("Source %d does not have an offset value", offsetNotAlign))
		}

		size, in := sizes[sourceId]
		if !in {
			panic(fmt.Sprintf("Source %d does not have a size value", size))
		}

		offsets[sourceId] = u32(offsetNotAlign)
		offsetNotAlign += u64(size)
	}
}

// Has no side effect
func VerifyDIDXDATA(d *AudioStore) {
	sourceIds := d.SourceIds
	audioDataIndices := d.AudioData
	offsets := d.Offsets
	sizes := d.Sizes

	if len(sourceIds) != len(audioDataIndices) {
		panic(fmt.Sprintf(
			"# of source ids (%d) does not equal to # of audio data (%d)", 
			len(sourceIds), len(audioDataIndices),
		))
	}

	if len(audioDataIndices) != len(offsets) {
		panic(fmt.Sprintf(
			"# of offset values (%d) does not equal to # of source ids (%d)",
			len(offsets), len(sourceIds),
		))
	}

	if len(offsets) != len(sizes) {
		panic(fmt.Sprintf(
			"# of size values (%d) does not equal to # of source ids (%d)",
			len(sizes), len(sourceIds),
		))
	}

	offsetNotAlignChecker := u64(0)
	for i, sourceId := range d.SourceIds {
		audioData, in := audioDataIndices[sourceId]
		if !in {
			panic(fmt.Sprintf(
				"Source %d does not have an associated audio data",
				sourceId,
			))
		}

		offset, in := offsets[sourceId]
		if !in {
			panic(fmt.Sprintf(
				"Source %d does not have an associated offset value",
				sourceId,
			))
		}

		size, in := sizes[sourceId]
		if !in {
			panic(fmt.Sprintf(
				"Source %d does not have an associated size value",
				sourceId,
			))
		}

		if u64(offset) != offsetNotAlignChecker {
			panic(fmt.Sprintf(
				"Expecting source %d (linear index at %d) has an offset of %d but receive %d",
				sourceId, i, offsetNotAlignChecker, offset,
			))
		}

		audioDataSize := len(audioData)
		if int(size) != audioDataSize {
			panic(fmt.Sprintf(
				"Source %d (linear index at %d) has a size value of %d but its audio data has a size value of %d",
				sourceId, i, size, audioDataSize,

			))
		}
	}
}

// Has no side effect
func EncodeDIDX(d *AudioStore, e *uio.EncoderCtx) (err error) {
	sourceIds := d.SourceIds

	chunkHeader := ChunkHeader{ [4]byte{'D', 'I', 'D', 'X'}, 12 * u32(len(sourceIds)) }
	if err = uio.EncodeStruct(e, chunkHeader, SizeOfChunkHeader); err != nil {
		return fmt.Errorf("Failed to encode DIDX chunk header: %w", err)
	}

	offsets := d.Offsets
	sizes := d.Sizes

	var payload MediaIndexS
	for i, sourceId := range sourceIds {
		offset, in := offsets[sourceId]
		if !in {
			panic(fmt.Sprintf("Source %d does not have a offset value", sourceId))
		}

		size, in := sizes[sourceId]
		if !in {
			panic(fmt.Sprintf("Source %d does not have a size value", sourceId))
		}

		payload.SourceId = sourceId
		payload.Offset = offset
		payload.Size = size

		err = uio.EncodeStruct(e, payload, SizeOfMediaIndex)
		if err != nil {
			return fmt.Errorf("Failed to encode media index (linear index %d) of source %d: %w", i, sourceId, err)
		}
	}

	return nil
}

// Has no side effect
func EncodeDATANotAlign(d *AudioStore, e *uio.EncoderCtx) (err error) {
	sourceIds := d.SourceIds
	audioData := d.AudioData

	for i, sourceId := range sourceIds {
		audioData, in := audioData[sourceId]
		if !in {
			panic(fmt.Sprintf("Source %d does not have an associated audio data", sourceId))
		}
		if err = uio.EncodeBytes(e, audioData); err != nil {
			return fmt.Errorf(
				"Failed to write audio data of source %d (linear index %d): %w",
				sourceId, i, err,
			)
		}
	}

	return err
}

func AllocAudioStore(size u32) *AudioStore {
	return &AudioStore{
		SourceIds: make([]u32, 0, size),
		Offsets: make(map[u32]u32, size),
		Sizes: make(map[u32]u32, size),
	}
}

// Has side effect
// Use this if assuming there will be no duplicate in DIDX entry (e.g., at 
// decoding phase)
func (d *AudioStore) AddMediaIndex(m MediaIndexS) {
	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	if slices.Contains(d.SourceIds, sourceId) {
		panic(fmt.Sprintf("Media index with %d already exist.", sourceId))
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
}

// Has side effect
func (d *AudioStore) AddMediaIndexCheck(m MediaIndexS) error {
	sourceId := m.SourceId
	offset := m.Offset
	size := m.Size

	if slices.Contains(d.SourceIds, sourceId) {
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
func (d *AudioStore) NumMediaIndex() u32 {
	return u32(len(d.SourceIds))
}

// No side effect
func (d *AudioStore) HasMediaIndex(sourceId u32) bool {
	return slices.Contains(d.SourceIds, sourceId)
}

// No side effect
// Use HasMediaIndex before MediaIndex
func (d *AudioStore) MediaIndex(sourceId u32) (offset u32, size u32) {
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
func (d *AudioStore) MediaIndexCheck(sourceId u32) (
	offset u32, size u32, in bool,
) {
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
// Use this when omiting all alignment at the decoding phase, or use it with 
// HasMediaIndex
func (d *AudioStore) SetMediaIndex(m MediaIndexS) {
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
