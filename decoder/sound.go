package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

// Expect r to be SectionReader, LimitReader, bytes.Reader, or bytes.Buffer 
// since these type of reader will do automatic bound checking.
func AllocDecodeSound(r io.Reader, o order, version u32, size u32) any {
	id := uio.U32P(r, o)

	sourceData := AllocDecodeSourceData(r, o, version)

	var pluginParam *wwise.PluginParam
	if wwise.SourceHasParam(sourceData) && wwise.SourceHasPluginID(sourceData) {
		pluginParam = AllocDecodePluginParam(r, o, version)
	}

	b := wwise.BaseParameter{}

	AllocDecodeBaseParameter(r, o, version, &b)

	sound := &wwise.Sound{
		Id: id,
		SourceData: sourceData,
		PluginParam: pluginParam,
		BaseParameter: &b,
	}

	return sound
}
