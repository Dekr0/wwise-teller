package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseBankSourceData(r *wio.Reader) wwise.BankSourceData {
	b := wwise.BankSourceData{
		PluginID:          r.U32Unsafe(),
		StreamType:        wwise.SourceType(r.U8Unsafe()),
		SourceID:          r.U32Unsafe(),
		InMemoryMediaSize: r.U32Unsafe(),
		SourceBits:        r.U8Unsafe(),
		PluginParam:       nil,
	}
	if !b.HasParam() {
		return b
	}
	if b.PluginID == 0 {
		return b
	}
	b.PluginParam = &wwise.PluginParam{
		PluginParamSize: r.U32Unsafe(), PluginParamData: []byte{},
	}
	if b.PluginParam.PluginParamSize <= 0 {
		return b
	}
	begin := r.Pos()
	b.PluginParam.PluginParamData = r.ReadNUnsafe(
		uint64(b.PluginParam.PluginParamSize), 0,
	)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(b.PluginParam.PluginParamSize, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal size in "+
			"source plugin parameter header",
	)
	return b
}
