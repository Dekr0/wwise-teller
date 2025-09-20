package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
func AllocDecodeSound(r io.Reader, o order, version u32, size u32, out any) {
	s := out.(*wwise.Sound)

	s.Id = uio.U32P(r, o)

	s.SourceData = *AllocDecodeSourceData(r, o, version)

	if wwise.SourceHasParam(s.SourceData) && wwise.SourceHasPluginID(s.SourceData) {
		s.PluginParam = *AllocDecodePluginParam(r, o, version)
	}

	AllocDecodeBaseParameter(r, o, version, &s.BaseParameter)
}
