// TODO:
// - Replacing audio data
// 		- Target import and automation should be first class citizen
package wwise

import (
	"context"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type DATAAppendOnly struct {
	I   uint8
	T []byte
	B []byte
}

func (d *DATAAppendOnly) Encode(ctx context.Context) ([]byte, error) {
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

func (d *DATAAppendOnly) Tag() []byte {
	return d.T
}

func (d *DATAAppendOnly) Idx() uint8 {
	return d.I
}

type DATA struct {
	I             uint8
	T           []byte
	Audios    [][]byte
	AudiosMap     map[uint32][]byte
}

func (d *DATA) Encode(ctx context.Context) ([]byte, error) {
	size := d.Size()
	bw := wio.NewWriter(uint64(SizeOfChunkHeader + size))
	bw.AppendBytes(d.T)
	bw.Append(size)
	for _, audio := range d.Audios {
		bw.AppendBytes(audio)
	}
	assert.Equal(
		int(size),
		bw.Len() - 4 - 4,
		"(DATA) The size of encoded data does not equal to calculated size.",
	)

	return bw.Bytes(), nil
}

func (d *DATA) Size() uint32 {
	size := uint32(0)
	for _, audio := range d.Audios {
		size += uint32(len(audio))
	}
	return size
}

func (d *DATA) Tag() []byte {
	return d.T
}

func (d *DATA) Idx() uint8 {
	return d.I
}
