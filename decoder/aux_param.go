package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeAuxParam(r io.Reader, o order) (a *wwise.AuxParam) {
	a = wwise.AllocAuxParam()
	a.SettingVector = uio.U8P(r, o)
	if wwise.HasAux(*a) {
		for i := range a.AuxIds {
			a.AuxIds[i] = uio.U32P(r, o)
		}
	}
	a.ReflectionAux = uio.U32P(r, o)
	return a
}
