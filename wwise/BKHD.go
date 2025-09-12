package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

const BaseSizeBKHD = 20

type BKHD struct {
	Version          u32
	Id               u32
	Language         u32
	// Alignment and DeviceAllocated combine into one single u32
	DeviceAllocated  u16
	Alignment        u16
	Project          u32
	Data           []u8
}

type BKHDEncodePayload struct {
	Version                      u32
	Id                           u32
	Language                     u32
	DeviceAllocatedWithAlignment u32
	Project                      u32
}

// Use for pre-allocation
func BKHDSize(b *BKHD) u32 {
	return BaseSizeBKHD + u32(len(b.Data))
}

func EncodeBKHD(b *BKHD, e *uio.EncoderCtx) (err error) {
	chunkHeader := ChunkHeader{ [4]byte{'B', 'K', 'H', 'D'}, BKHDSize(b) }
	err = uio.EncodeStruct(e, chunkHeader, SizeOfChunkHeader)
	if err != nil {
		return fmt.Errorf("Failed to encode BKHD chunk header: %w", err)
	}

	payload := BKHDEncodePayload{
		Version: b.Version,
		Id: b.Id,
		Language: b.Language,
		DeviceAllocatedWithAlignment: (u32(b.Alignment) << 16) | u32(b.DeviceAllocated),
		Project: b.Project,
	}
	err = uio.EncodeStruct(e, payload, Size32 * 5)
	if err != nil {
		return fmt.Errorf("Failed to encode BKHD basic field: %w", err)
	}

	err = uio.EncodeBytes(e, b.Data) 
	if err != nil {
		err = fmt.Errorf("Failed to write BKHD encoded data: %w", err)
	}

	return err
}
