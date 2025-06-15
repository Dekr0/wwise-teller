package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseFxCustom(size uint32, r *wio.Reader) *wwise.FxCustom {
	assert.Equal(0, r.Pos(), "Fx Custom parser position doesn't start at 0.")
	begin := r.Pos()
	f := wwise.FxCustom{
		Id: r.U32Unsafe(),
		PluginTypeId: r.U32Unsafe(),
	}
	if f.HasParam() {
		f.PluginParam = &wwise.PluginParam{}
		ParsePluginParam(r, f.PluginParam, f.PluginTypeId)
	}
	f.MediaMap = make([]wwise.MediaMapItem, r.U8Unsafe())
	for i := range f.MediaMap {
		f.MediaMap[i].Index = r.U8Unsafe()
		f.MediaMap[i].SourceId = r.U32Unsafe()
	}
	ParseRTPC(r, &f.RTPC)
	ParseStateProp(r, &f.StateProp)
	ParseStateGroup(r, &f.StateGroup)
	f.PluginProps = make([]wwise.PluginProp, r.U16Unsafe())
	for i := range f.PluginProps {
		f.PluginProps[i].PropertyID = wwise.RTPCParameterType(r.U8Unsafe())
		f.PluginProps[i].RTPCAccum = wwise.RTPCAccumType(r.U8Unsafe())
		f.PluginProps[i].Value = r.F32Unsafe()
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return &f
}

func ParseFxShareSet(size uint32, r *wio.Reader) *wwise.FxShareSet {
	assert.Equal(0, r.Pos(), "Fx Share Set parser position doesn't start at 0.")
	begin := r.Pos()
	f := wwise.FxShareSet{
		Id: r.U32Unsafe(),
		PluginTypeId: r.U32Unsafe(),
	}
	if f.HasParam() {
		f.PluginParam = &wwise.PluginParam{}
		ParsePluginParam(r, f.PluginParam, f.PluginTypeId)
	}
	f.MediaMap = make([]wwise.MediaMapItem, r.U8Unsafe())
	for i := range f.MediaMap {
		f.MediaMap[i].Index = r.U8Unsafe()
		f.MediaMap[i].SourceId = r.U32Unsafe()
	}
	ParseRTPC(r, &f.RTPC)
	ParseStateProp(r, &f.StateProp)
	ParseStateGroup(r, &f.StateGroup)
	f.PluginProps = make([]wwise.PluginProp, r.U16Unsafe())
	for i := range f.PluginProps {
		f.PluginProps[i].PropertyID = wwise.RTPCParameterType(r.U8Unsafe())
		f.PluginProps[i].RTPCAccum = wwise.RTPCAccumType(r.U8Unsafe())
		f.PluginProps[i].Value = r.F32Unsafe()
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return &f
}

func ParseFxChunk(r *wio.Reader, f *wwise.FxChunk) {
	UniqueNumFx := r.U8Unsafe()
	if UniqueNumFx <= 0 {
		f.BitsFxByPass = 0
		f.FxChunkItems = make([]wwise.FxChunkItem, 0)
		return
	}
	f.BitsFxByPass = r.U8Unsafe()
	f.FxChunkItems = make([]wwise.FxChunkItem, UniqueNumFx)
	for i := range f.FxChunkItems {
		f.FxChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxChunkItems[i].FxId = r.U32Unsafe()
		f.FxChunkItems[i].BitIsShareSet = r.U8Unsafe()
		f.FxChunkItems[i].BitIsRendered = r.U8Unsafe()
	}
}

func ParseFxChunkMetadata(r *wio.Reader, f *wwise.FxChunkMetadata)  {
	f.BitIsOverrideParentMetadata = r.U8Unsafe()
	UniqueNumFxMetadata := r.U8Unsafe()
	f.FxMetaDataChunkItems = make([]wwise.FxChunkMetadataItem, UniqueNumFxMetadata)
	for i := range f.FxMetaDataChunkItems {
		f.FxMetaDataChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxMetaDataChunkItems[i].FxId = r.U32Unsafe()
		f.FxMetaDataChunkItems[i].BitIsShareSet = r.U8Unsafe()
	}
}
