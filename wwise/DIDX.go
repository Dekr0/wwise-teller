package wwise

import (
	"github.com/Dekr0/wwise-teller/reader"
	"github.com/Dekr0/wwise-teller/assert"
)

type DIDX struct {
	UniqueNumMedias uint32 // U32
	MediaIndexs []*MediaIndex
}

func NewDIDX(uniqueNumMedias uint32) *DIDX {
	return &DIDX{
		uniqueNumMedias,
		make([]*MediaIndex, 0, uniqueNumMedias),
	}
}

func (d *DIDX) Encode() []byte {
	assert.AssertEqual(
		uint32(len(d.MediaIndexs)),
		d.UniqueNumMedias,
		"# of media index specified in the header doesn't equal to # of media index in the slice.",
	)
	chunkSize := uint32(d.UniqueNumMedias * MEDIA_INDEX_FIELD_SIZE)
	bw := reader.NewFixedSizeBlobWriter(uint64(CHUNK_HEADER_SIZE + chunkSize))
	bw.AppendBytes([]byte{'D', 'I', 'D', 'X'})
	bw.Append(chunkSize)
	for _, m := range d.MediaIndexs { bw.Append(m) }
	assert.AssertEqual(
		int(chunkSize),
		bw.Len() - 4 - 4,
		"(DIDX) The size of encoded data does not match with calculated size",
	)
	return bw.GetBlob()
}

const MEDIA_INDEX_FIELD_SIZE = 12

type MediaIndex struct {
	AudioSrcId uint32 // sid
	DATAOffset uint32 // U32
	DATABlobSize uint32 // U32
}
