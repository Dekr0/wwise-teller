package wwise

import (
	"context"
	"fmt"
	"slices"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

// 1 mean no alignment
var PossibleDataAlignments []uint8 = []uint8{ 1, 2, 4, 8, 16, 24, 32, 64 }

type DIDX struct {
	I                    uint8
	T                  []byte
	Alignment            uint8 // guess
	AvailableAlignment []uint8
	MediaIndexs        []MediaIndex
	MediaIndexsMap     map[uint32]*MediaIndex
}

func NewDIDX(I uint8, T []byte, num uint32) *DIDX {
	return &DIDX{
		I,
		T,
		1,
		[]uint8{},
		make([]MediaIndex, 0, num),
		make(map[uint32]*MediaIndex, num),
	}
}

func (d *DIDX) Encode(ctx context.Context) ([]byte, error) {
	size := uint32(len(d.MediaIndexs) * SizeOfMediaIndex)
	w := wio.NewWriter(uint64(SizeOfChunkHeader + size))
	w.AppendBytes(d.T)
	w.Append(size)
	for _, m := range d.MediaIndexs { w.Append(m) }
	assert.Equal(
		int(size),
		w.Len() - 4 - 4,
		"(DIDX) The size of encoded data does not match with calculated size",
	)
	return w.Bytes(), nil
}

func (d *DIDX) Tag() []byte {
	return d.T
}

func (d *DIDX) Idx() uint8 {
	return d.I
}

func (d *DIDX) GuessAlignment() {
	d.Alignment = 1
	for _, a := range PossibleDataAlignments {
		algined := true
		for _, m := range d.MediaIndexs {
			if m.Offset % uint32(a) != 0 {
				algined = false
				break
			}
		}
		if algined {
			d.AvailableAlignment = append(d.AvailableAlignment, a)
			d.Alignment = max(a, d.Alignment)
		}
	}
}

func (d *DIDX) Append(sid uint32, size uint32) error {
	if len(d.MediaIndexs) == 0 {
		d.MediaIndexs = append(d.MediaIndexs, MediaIndex{sid, 0, size})
		return nil
	}
	last := d.MediaIndexs[len(d.MediaIndexs) - 1]
	if _, in := d.MediaIndexsMap[sid]; in {
		return fmt.Errorf("Source ID %d alreay exists.", sid)
	}
	d.MediaIndexs = append(d.MediaIndexs, MediaIndex{
		sid,
		last.Offset + last.Size,
		size,
	})
	d.MediaIndexsMap[sid] = &d.MediaIndexs[len(d.MediaIndexs) - 1]
	return nil
}

func (d *DIDX) Remove(sid uint32) {
	if len(d.MediaIndexs) == 0 {
		return
	}
	if _, in := d.MediaIndexsMap[sid]; in {
		delete(d.MediaIndexsMap, sid)
	}
	d.MediaIndexs = slices.DeleteFunc(d.MediaIndexs, func(m MediaIndex) bool {
		return m.Sid == sid
	})
}

const SizeOfMediaIndex = 12

type MediaIndex struct {
	Sid uint32 // sid
	Offset uint32 // U32
	Size uint32 // U32
}
