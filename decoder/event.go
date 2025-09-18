package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
func AllocDecodeEvent(r io.Reader, o order, ver u32, size u32) any {
	id := uio.U32P(r, o)

	data := wwise.AllocEventData(uio.VV128P(r, o))
	for i := range data.ActionIds {
		data.ActionIds[i] = uio.U32P(r, o)
	}

	e := &wwise.Event{Id: id, EventData: data}

	return e
}
