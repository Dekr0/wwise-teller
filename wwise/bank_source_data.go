package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type SourceType uint8

const (
	SourceTypeDATA              SourceType = 0
	SourceTypeStreaming         SourceType = 1
	SourceTypePrefetchStreaming SourceType = 2
	SourceTypeCount             SourceType = 3
)
var SourceTypeNames []string = []string{
	"DATA",
	"Streaming",
	"Prefetch Streaming",
}

type BankSourceData struct {
	PluginID          uint32 // U32
	StreamType        SourceType // U8x
	SourceID          uint32 // tid
	InMemoryMediaSize uint32 // U32
	SourceBits        uint8 // U8x
	PluginParam      *PluginParam
}

func (b *BankSourceData) PluginType() uint32 {
	return (b.PluginID >> 0) & 0x000F
}

func (b *BankSourceData) Company() uint32 {
	return (b.PluginID >> 4) & 0x03FF
}

func (b *BankSourceData) LanguageSpecific() bool {
	return b.SourceBits & 0b00000001 != 0
}

func (b *BankSourceData) Prefetch() bool {
	return b.SourceBits & 0b00000010 != 0
}

func (b *BankSourceData) NonCacheable() bool {
	return b.SourceBits & 0b00001000 != 0
}

func (b *BankSourceData) HasSource() bool {
	return b.SourceBits & 0b10000000 != 0
}

func (b *BankSourceData) HasParam() bool {
	return (b.PluginID & 0x0F) == 2
}

func (b *BankSourceData) ChangeSource(sid uint32, inMemoryMediaSize uint32) {
	b.SourceID = sid
	b.InMemoryMediaSize = inMemoryMediaSize
}

func (b *BankSourceData) Encode() []byte {
	b.assert()
	size := b.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(b.PluginID)
	w.Append(b.StreamType)
	w.Append(b.SourceID)
	w.Append(b.InMemoryMediaSize)
	w.Append(b.SourceBits)
	if b.PluginParam != nil {
		w.AppendBytes(b.PluginParam.Encode())
	}
	return w.BytesAssert(int(size))
}

func (b *BankSourceData) Size() uint32 {
	size := uint32(4 + 1 + 4 + 4 + 1)
	if b.PluginParam != nil {
		size += b.PluginParam.Size()
	}
	return size
}

func (b *BankSourceData) assert() {
	if !b.HasParam() {
		assert.Nil(b.PluginParam,
			"Plugin type indicate that there's no plugin parameter data.",
		)
		return
	}
	// This make no sense???
	if b.PluginID == 0 {
		assert.Nil(b.PluginParam,
			"Plugin type indicate that there's no plugin parameter data.",
		)
	}
}

