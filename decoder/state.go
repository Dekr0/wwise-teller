package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
// Has side effect on HIRC
func AllocDecodeState(r io.Reader, o order, ver u32, size u32) any {
	id := uio.U32P(r, o)

	data := wwise.AllocStateProps(uio.U16P(r, o))

	for i := range data.Ids {
		data.Ids[i] = uio.U16P(r, o)
		data.Vals[i] = uio.F32P(r, o)
	}

	res := &wwise.StateH{ Id: id, StateProps: data }

	return res
}
