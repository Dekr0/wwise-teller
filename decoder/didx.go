package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func DecodeDIDX(
	inReader  io.Reader,
	chunkSize u32,
	o         order,
) (
	d *wwise.DIDXDATA, err error,
) {
	num := chunkSize / wwise.SizeOfMediaIndex

	r := io.LimitReader(inReader, int64(chunkSize))

	d = wwise.NewDIDXDATA(chunkSize)

	var m wwise.MediaIndexEntry
	for range num {
		err = uio.Struct(r, o, &m)
		if err != nil {
			return nil, err
		}
		wwise.NewMediaIndex(d, m)
	}

	return d, nil
}
