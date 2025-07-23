package wwise

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)


type FxChunk struct {
	// UniqueNumFx uint8
	BitsFxByPass uint8 // > 145 is called bBypassAll
	FxChunkItems []FxChunkItem
}

func (f *FxChunk) BypassFx(set bool) {
	if set {
		f.BitsFxByPass = 1
		for i := range f.FxChunkItems {
			f.FxChunkItems[i].BitVector = 1
		}
	} else {
		f.BitsFxByPass = 0
		for i := range f.FxChunkItems {
			f.FxChunkItems[i].BitVector = 0
		}
	}
}

func (f *FxChunk) Clone() FxChunk {
	return FxChunk{f.BitsFxByPass, slices.Clone(f.FxChunkItems)}
}

func (f *FxChunk) Encode(v int) []byte {
	if len(f.FxChunkItems) <= 0 {
		return []byte{0}
	}
	size := f.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(f.FxChunkItems)))
	w.AppendByte(f.BitsFxByPass)
	for _, i := range f.FxChunkItems {
		w.Append(i.UniqueFxIndex)
		w.Append(i.FxId)
		if v <= 145 {
			w.Append(i.BitIsShareSet)
			w.Append(i.BitIsRendered)
		} else {
			w.Append(i.BitVector)
		}
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunk) Size(v int) uint32 {
	if len(f.FxChunkItems) <= 0 {
		return 1
	}
	if v <= 145 {
		return uint32(1 + 1 + len(f.FxChunkItems) * SizeOfFxChunk_LE145)
	}
	return uint32(1 + 1 + len(f.FxChunkItems) * SizeOfFxChunk_G145)
}

const SizeOfFxChunk_LE145 = 7
const SizeOfFxChunk_G145 = 6
type FxChunkItem struct {
	UniqueFxIndex uint8 // u8i
	FxId          uint32 // tid
	// The following is exclusive
	// <= 145
	BitIsShareSet uint8 // U8x
	BitIsRendered uint8 // U8x unused (effects can't render)
	// > 145
	BitVector     uint8
}

func (f *FxChunkItem) Bypass(set bool) {
	if set {
		f.BitVector = 1
	} else {
		f.BitVector = 0
	}
}

type FxChunkMetadata struct {
	BitIsOverrideParentMetadata uint8
	// UniqueNumFxMetadata uint8
	FxMetaDataChunkItems []FxChunkMetadataItem
}

func (f *FxChunkMetadata) Clone() FxChunkMetadata {
	return FxChunkMetadata{f.BitIsOverrideParentMetadata, slices.Clone(f.FxMetaDataChunkItems)}
}

func (f *FxChunkMetadata) Encode(v int) []byte {
	size := f.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendByte(f.BitIsOverrideParentMetadata)
	w.AppendByte(uint8(len(f.FxMetaDataChunkItems)))
	for _, i := range f.FxMetaDataChunkItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunkMetadata) Size(int) uint32 {
	return uint32(1 + 1 + len(f.FxMetaDataChunkItems) * SizeOfFxChunkMetadata)
}

const SizeOfFxChunkMetadata = 6
type FxChunkMetadataItem struct {
	UniqueFxIndex uint8  // u8i
	FxId          uint32 // tid
	BitIsShareSet uint8  // U8x
}
