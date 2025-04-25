package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type DATA struct {
	b []byte
}

func NewDATA(b []byte) *DATA {
	return &DATA{b: b}
}

func (d *DATA) Encode() []byte {
	size := uint32(len(d.b))
	bw := wio.NewWriter(uint64(chunkHeaderSize + size))
	bw.AppendBytes([]byte{'D', 'A', 'T', 'A'})
	bw.Append(size)
	bw.AppendBytes(d.b)
	assert.Equal(
		int(size),
		bw.Len() - 4 - 4,
		"(DATA) The size of encoded data does not equal to calculated size.",
	)

	return bw.Bytes()
}
