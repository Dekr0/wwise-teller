package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
// Has side effect on HIRC
func AllocDecodeState(r io.Reader, o order, ver u32, size u32, out any) {
	s := out.(*wwise.State)
	s.Id = uio.U32P(r, o)

	s.StateProps = *wwise.AllocStateProps(uio.U16P(r, o))

	for i := range s.StateProps.Ids {
		s.StateProps.Ids[i] = uio.U16P(r, o)
		s.StateProps.Vals[i] = uio.F32P(r, o)
	}
}
