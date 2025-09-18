package decoder

import (
	"fmt"
	"io"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecodePluginParam(
	r            io.Reader,
	o            order,
	version      u32,
) (pluginParam *wwise.PluginParam) {
	pluginParam = wwise.AllocPluginParam(uio.U32P(r, o))
	if pluginParam.Size > 0 {
		pluginParam.Data = make([]byte, pluginParam.Size, pluginParam.Size)
		if _, err := io.ReadFull(r, pluginParam.Data); err != nil {
			panic(fmt.Errorf("Failed to read plugin parameter data: %w", err))
		}
	}
	return pluginParam
}

func AllocDecodeFXs(r io.Reader, o order, version u32) (f *wwise.FXs) {
	numFXs := uio.U8P(r, o)
	f = wwise.AllocFXs(numFXs)
	if numFXs <= 0 {
		return f
	}
	f.BypassAll = uio.U8P(r, o)
	for i := range f.FXs {
		f.FXs[i].Idx = uio.U8P(r, o)
		f.FXs[i].Id  = uio.U32P(r, o)
		if version <= 145 {
			f.FXs[i].IsShareSet = uio.U8P(r, o)
			f.FXs[i].IsRender   = uio.U8P(r, o)
		} else {
			f.FXs[i].BitVector = uio.U8P(r, o)
		}
	}
	return f
}

func AllocDecodeFXMetadatas(r io.Reader, o order) (
	f *wwise.FxMetadatas,
) {
	f = wwise.AllocFxMetadatas(uio.U8P(r, o))
	for i := range f.Idx {
		f.Idx[i] = uio.U8P(r, o)
		f.Id[i] = uio.U32P(r, o)
		f.IsShareSet[i] = uio.U8P(r, o)
	}
	return f
}

