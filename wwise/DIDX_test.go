package wwise

import (
	"testing"
)

func TestMediaIndexEncode(t *testing.T) {
	d := DIDX{2, make([]*MediaIndex, 0, 2)}

	d.MediaIndexs = append(d.MediaIndexs, &MediaIndex{
		AudioSrcId: 1,
		DATAOffset: 0,
		DATABlobSize: 4,
	})

	d.MediaIndexs = append(d.MediaIndexs, &MediaIndex{
		AudioSrcId: 2,
		DATAOffset: 4,
		DATABlobSize: 4,
	})

	b := d.Encode()
	if uint32(len(b)) != 4 + 4 + d.UniqueNumMedias * MEDIA_INDEX_FIELD_SIZE {
		t.Fail()
	}
}
