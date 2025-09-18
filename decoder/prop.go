package decoder

import (
	"io"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecodeProp(r io.Reader, o order) (p *wwise.Prop) {
	p = wwise.AllocProp(uio.U8P(r, o))

	for i := range p.Ids {
		p.Ids[i] = uio.U8P(r, o)
	}
	for i := range p.Vals {
		p.Vals[i] = uio.F32P(r, o)
	}

	return p
}

func AllocDecodeRProp(r io.Reader, o order) (p *wwise.RProp) {
	p = wwise.AllocRProp(uio.U8P(r, o))

	for i := range p.Ids {
		p.Ids[i] = uio.U8P(r, o)
	}

	for i := range p.Ids {
		p.Mins[i] = uio.F32P(r, o)
		p.Maxs[i] = uio.F32P(r, o)
	}

	return p 
}
