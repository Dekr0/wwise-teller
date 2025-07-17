package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseBankSourceData(r *wio.Reader, v int) wwise.BankSourceData {
	var b wwise.BankSourceData = wwise.BankSourceData{}
	b.PluginID = r.U32Unsafe()
	b.StreamType = wwise.SourceType(r.U8Unsafe())
	b.SourceID = r.U32Unsafe()
	if v <= 150 {
		b.InMemoryMediaSize = r.U32Unsafe()
	} else {
		b.CacheID = r.U32Unsafe()
		b.InMemoryMediaSize = r.U32Unsafe()
	}
	b.SourceBits = r.U8Unsafe()
	b.PluginParam = nil
	if !b.HasParam() {
		return b
	}
	if b.PluginID == 0 {
		return b
	}
	b.PluginParam = &wwise.PluginParam{
		PluginParamSize: r.U32Unsafe(), PluginParamData: &wwise.FxPlaceholder{Data: []byte{}},
	}
	if b.PluginParam.PluginParamSize <= 0 {
		return b
	}
	begin := r.Pos()
	ParsePluginParam(r, b.PluginParam, b.PluginID, v)
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
