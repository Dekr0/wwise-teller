package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeActorMixer(r io.Reader, o order, version u32, size u32) any {
	id := uio.U32P(r, o)

	b := &wwise.BaseParameter{}

	AllocDecodeBaseParameter(r, o, version, b)

	c := AllocDecodeContainer(r, o)

	return &wwise.ActorMixer{
		Id: id, BaseParameter: b, Container: c,
	}
}
