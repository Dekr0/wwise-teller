package wwise

import (
	"bytes"
	bin "encoding/binary"
	"io"

	"github.com/Dekr0/wwise-teller/wio"
)

const BKHDSize = 4 * 5

type BKHD struct {
	Version           uint32
	SoundbankID       uint32
	LanguageID        uint32

	// Alignment and `Device Allocated` combine into one single 32-bits integer
	Alignment         uint16
	DeviceAllocated   uint16

	ProjectID         uint32
	Rest            []byte
}

func (b *BKHD) Size() uint32 {
	return BKHDSize + uint32(len(b.Rest))
}

func (b *BKHD) Encode(o bin.ByteOrder) []byte {
	size := b.Size()
	buff := bytes.NewBuffer(make([]byte, 0, size)) 
	b.Write(buff, o, size)
	return buff.Bytes()
}

func (b *BKHD) Write(w io.Writer, o bin.ByteOrder, size uint32) {
	p := uint32(0)
	wio.MustWriteTrack(w, o, b.Version, "Bank Version", &p, 4)
	wio.MustWriteTrack(w, o, b.SoundbankID, "Sound bank ID", &p, 4)
	wio.MustWriteTrack(w, o, b.LanguageID, "Language ID", &p, 4)
	combined := uint32(b.DeviceAllocated) << 16 | uint32(b.Alignment)
	wio.MustWriteTrack(w, o, combined, "Device Allocated and Alignment", &p, 4)
	wio.MustWriteTrack(w, o, b.ProjectID, "Project ID", &p, 4)
	wio.MustWriteTrack(w, o, b.Rest, "Rest of BKHD", &p, uint32(len(b.Rest)))
	wio.AssertFit(size, p)
}
