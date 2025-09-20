package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeActorMixer(r io.Reader, o order, version u32, size u32, inOut any) {
	a := inOut.(*wwise.ActorMixer)

	a.Id = uio.U32P(r, o)

	AllocDecodeBaseParameter(r, o, version, &a.BaseParameter)

	a.Container = *AllocDecodeContainer(r, o)
}
