package wwise

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)


type FxChunk struct {
	// UniqueNumFx uint8
	BitsFxByPass uint8
	FxChunkItems []FxChunkItem
}

func (f *FxChunk) Clone() FxChunk {
	return FxChunk{f.BitsFxByPass, slices.Clone(f.FxChunkItems)}
}

func (f *FxChunk) Encode() []byte {
	if len(f.FxChunkItems) <= 0 {
		return []byte{0}
	}
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(f.FxChunkItems)))
	w.AppendByte(f.BitsFxByPass)
	for _, i := range f.FxChunkItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunk) Size() uint32 {
	if len(f.FxChunkItems) <= 0 {
		return 1
	}
	return uint32(1 + 1 + len(f.FxChunkItems) * SizeOfFxChunk)
}

const SizeOfFxChunk = 7
type FxChunkItem struct {
	UniqueFxIndex uint8 // u8i
	FxId          uint32 // tid
	BitIsShareSet uint8 // U8x
	BitIsRendered uint8 // U8x unused (effects can't render)
}

type FxChunkMetadata struct {
	BitIsOverrideParentMetadata uint8
	// UniqueNumFxMetadata uint8
	FxMetaDataChunkItems []FxChunkMetadataItem
}

func (f *FxChunkMetadata) Clone() FxChunkMetadata {
	return FxChunkMetadata{f.BitIsOverrideParentMetadata, slices.Clone(f.FxMetaDataChunkItems)}
}

func (f *FxChunkMetadata) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(f.BitIsOverrideParentMetadata)
	w.AppendByte(uint8(len(f.FxMetaDataChunkItems)))
	for _, i := range f.FxMetaDataChunkItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunkMetadata) Size() uint32 {
	return uint32(1 + 1 + len(f.FxMetaDataChunkItems) * SizeOfFxChunkMetadata)
}

const SizeOfFxChunkMetadata = 6
type FxChunkMetadataItem struct {
	UniqueFxIndex uint8  // u8i
	FxId          uint32 // tid
	BitIsShareSet uint8  // U8x
}
