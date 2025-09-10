package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
func DecodeEvent(r io.Reader, o order, ver u32, h *wwise.HIRC, size u32) {
	id := uio.U32P(r, o)

	data := &wwise.EventData{}
	data.NumActionIds = *uio.VV128P(r, o)
	v := data.NumActionIds.V
	data.ActionIds = make([]u32, v, v)
	for i := range v {
		data.ActionIds[i] = uio.U32P(r, o)
	}

	wwise.NewEvent(h, id, data)
}
