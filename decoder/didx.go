package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func DecodeDIDX(
	inReader io.Reader,
	o        order,
) (
	d *wwise.DIDX, err error,
) {
	chunkSize, err := uio.U32(inReader, o)
	if err != nil {
		return nil, err
	}

	num := chunkSize / wwise.SizeOfMediaIndex

	r := io.LimitReader(inReader, int64(chunkSize))

	d = wwise.NewDIDX(chunkSize)

	var m wwise.MediaIndexEntry
	for range num {
		err = uio.Struct(r, o, &m)
		if err != nil {
			return nil, err
		}
		wwise.AddNewMediaIndex(d, m)
	}

	return d, nil
}
