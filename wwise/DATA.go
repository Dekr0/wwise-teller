// TODO:
// - Replacing audio data
// 		- Target import and automation should be first class citizen
package wwise

import (
	"context"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type DATA struct {
	I uint8
	T []byte
	B []byte
}

func NewDATA(I uint8, T []byte, b []byte) *DATA {
	return &DATA{I: I, T: T, B: b}
}

func (d *DATA) Encode(ctx context.Context) ([]byte, error) {
	size := uint32(len(d.B))
	bw := wio.NewWriter(uint64(SizeOfChunkHeader + size))
	bw.AppendBytes(d.T)
	bw.Append(size)
	bw.AppendBytes(d.B)
	assert.Equal(
		int(size),
		bw.Len() - 4 - 4,
		"(DATA) The size of encoded data does not equal to calculated size.",
	)

	return bw.Bytes(), nil
}

func (d *DATA) Tag() []byte {
	return d.T
}

func (d *DATA) Idx() uint8 {
	return d.I
}
