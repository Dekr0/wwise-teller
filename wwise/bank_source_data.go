package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const (
	PCM      =  0x00010001
	ADPCM    =  0x00020001
	VORBIS   =  0x00040001
	WEM_OPUS =  0x00140001
)

const (
	STREAM_TYPE_BNK = 0x00
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
	CacheID           uint32 // > 150
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

func (b *BankSourceData) Encode(v int) []byte {
	b.assert()
	size := b.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(b.PluginID)
	w.Append(b.StreamType)
	w.Append(b.SourceID)
	if v > 150 {
		w.Append(b.CacheID)
	}
	w.Append(b.InMemoryMediaSize)
	w.Append(b.SourceBits)
	if b.PluginParam != nil {
		w.AppendBytes(b.PluginParam.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (b *BankSourceData) Size(v int) uint32 {
	size := uint32(4 + 1 + 4 + 4 + 1)
	if v > 150 {
		size += 4
	}
	if b.PluginParam != nil {
		size += b.PluginParam.Size(v)
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

