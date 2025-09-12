package decoder

import (
	"fmt"
	"io"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecodeDIDX(
	inReader  io.Reader,
	chunkSize u32,
	o         order,
) (
	d *wwise.AudioStore, err error,
) {
	num := chunkSize / wwise.SizeOfMediaIndex

	r := io.LimitReader(inReader, int64(chunkSize))

	d = wwise.AllocAudioStore(chunkSize)

	var m wwise.MediaIndexS
	for i := range num {
		err = uio.Struct(r, o, &m)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode media index entry at %d: %w", i, err)
		}
		wwise.AddMediaIndex(d, m)
	}

	return d, nil
}
