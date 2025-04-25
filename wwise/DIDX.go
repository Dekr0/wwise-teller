package wwise

import (
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/assert"
)

type DIDX struct {
	MediaIndexs []*MediaIndex
}

func NewDIDX(num uint32) *DIDX {
	return &DIDX{make([]*MediaIndex, 0, num)}
}

func (d *DIDX) Encode() []byte {
	size := uint32(len(d.MediaIndexs) * mediaIndexFieldSize)
	w := wio.NewWriter(uint64(chunkHeaderSize + size))
	w.AppendBytes([]byte{'D', 'I', 'D', 'X'})
	w.Append(size)
	for _, m := range d.MediaIndexs { w.Append(m) }
	assert.Equal(
		int(size),
		w.Len() - 4 - 4,
		"(DIDX) The size of encoded data does not match with calculated size",
	)
	return w.Bytes()
}

const mediaIndexFieldSize = 12

type MediaIndex struct {
	Sid uint32 // sid
	Offset uint32 // U32
	Size uint32 // U32
}
