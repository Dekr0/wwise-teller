package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeStateProp(r io.Reader, o order) (s *wwise.StateProp) {
	numStateProp := uio.VV128P(r, o)
	s = wwise.AllocStateProp(*numStateProp)
	for i := range s.PropId {
		s.PropId[i] = *uio.VV128P(r, o)
		s.AccumType[i] = uio.U8P(r, o)
		s.InDb[i] = uio.U8P(r, o)
	}
	return s
}
