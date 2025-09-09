package wwise

import (
	bin "encoding/binary"
	"io"
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

func EncodeBKHD(b *BKHD, w io.Writer, o bin.ByteOrder) (err error) {
	chunkHeader := ChunkHeader{ [4]byte{'B', 'K', 'H', 'D'}, BKHDSize(b)}
	if err = bin.Write(w, o, chunkHeader); err != nil {
		return err
	}

	payload := BKHDEncodePayload{
		Version: b.Version,
		Id: b.Id,
		Language: b.Language,
		DeviceAllocatedWithAlignment: (u32(b.Alignment) << 16) | u32(b.DeviceAllocated),
		Project: b.Project,
	}
	if err = bin.Write(w, o, payload); err != nil {
		return err
	}

	_, err = w.Write(b.Data) 

	return err
}
