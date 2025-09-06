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

// Use for pre-allocation
func (b *BKHD) Size() u32 {
	return BaseSizeBKHD + u32(len(b.Data))
}

func (b *BKHD) Encode(w io.Writer, o bin.ByteOrder) (err error) {
	if err = bin.Write(w, o, b.Version); err != nil { return err }
	if err = bin.Write(w, o, b.Id); err != nil { return err }
	if err = bin.Write(w, o, b.Language); err != nil { return err }
	joint := (u32(b.DeviceAllocated) << 16) | (u32(b.Alignment))
	if err = bin.Write(w, o, joint); err != nil { return err }
	if err = bin.Write(w, o, b.Project); err != nil { return err }
	_, err = w.Write(b.Data)
	return err
}
