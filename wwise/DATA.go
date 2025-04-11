package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
)

type DATA struct {
	Blob []byte
}

func NewDATA(blob []byte) *DATA {
	return &DATA{Blob: blob}
}

func (d *DATA) Encode() []byte {
	chunkSize := uint32(len(d.Blob))
	bw := reader.NewFixedSizeBlobWriter(uint64(CHUNK_HEADER_SIZE + chunkSize))
	bw.AppendBytes([]byte{'D', 'A', 'T', 'A'})
	bw.Append(chunkSize)
	bw.AppendBytes(d.Blob)
	assert.AssertEqual(
		int(chunkSize),
		bw.Len() - 4 - 4,
		"(DATA) The size of encoded data does not equal to calculated size.",
	)

	return bw.GetBlob()
}
