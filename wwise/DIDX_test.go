package wwise

import (
	"testing"
)

func TestMediaIndexEncode(t *testing.T) {
	d := DIDX{make([]*MediaIndex, 0, 2)}

	d.MediaIndexs = append(d.MediaIndexs, &MediaIndex{
		Sid: 1,
		Offset: 0,
		Size: 4,
	})

	d.MediaIndexs = append(d.MediaIndexs, &MediaIndex{
		Sid: 2,
		Offset: 4,
		Size: 4,
	})

	b := d.Encode()
	if len(b) != 4 + 4 + len(d.MediaIndexs) * mediaIndexFieldSize {
		t.Fail()
	}
}
