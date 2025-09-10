package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
func DecodeState(r io.Reader, o order, ver u32, h *wwise.HIRC, size u32) {
	id := uio.U32P(r, o)

	numStateProps := uio.U16P(r, o)

	data := &wwise.StateProps{
		Ids: make([]u16, numStateProps, numStateProps),
		Vals: make([]f32, numStateProps, numStateProps),
	}

	for i := range numStateProps {
		data.Ids[i] = uio.U16P(r, o)
		data.Vals[i] = uio.F32P(r, o)
	}

	wwise.NewState(h, id, data)
}
