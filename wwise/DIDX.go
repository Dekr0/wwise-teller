package wwise

import (
	"context"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type DIDX struct {
	I uint8
	T []byte
	MediaIndexs []*MediaIndex
}

func NewDIDX(I uint8, T []byte, num uint32) *DIDX {
	return &DIDX{I, T, make([]*MediaIndex, 0, num)}
}

func (d *DIDX) Encode(ctx context.Context) ([]byte, error) {
	size := uint32(len(d.MediaIndexs) * mediaIndexFieldSize)
	w := wio.NewWriter(uint64(chunkHeaderSize + size))
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

const mediaIndexFieldSize = 12

type MediaIndex struct {
	Sid uint32 // sid
	Offset uint32 // U32
	Size uint32 // U32
}
