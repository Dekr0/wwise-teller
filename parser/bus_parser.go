package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseBus(size uint32, r *wio.Reader, v int) *wwise.Bus {
	assert.Equal(0, r.Pos(), "Bus parser position doesn't start at 0.")
	begin := r.Pos()
	bus := wwise.Bus{
		Id: r.U32Unsafe(),
		CanSetHDR: -1,
		OverrideBusId: r.U32Unsafe(),
	}
	if bus.OverrideBusId == 0 {
		bus.DeviceShareSetID = r.U32Unsafe()
	}
	ParsePropBundle(r, &bus.PropBundle, v)
	ParsePositioningParam(r, &bus.PositioningParam, v)
	ParseAuxParam(r, &bus.AuxParam, v)

	bus.VirtualBehaviorBitVector = r.U8Unsafe()
	bus.MaxNumInstance = r.U16Unsafe()
	bus.ChannelConf = r.U32Unsafe()
	bus.HDRBitVector = r.U8Unsafe()
	bus.RecoveryTime = r.I32Unsafe()
	bus.MaxDuckVolume = r.F32Unsafe()

	bus.DuckInfoList = make([]wwise.DuckInfo, r.U32Unsafe())
	for i := range bus.DuckInfoList {
		bus.DuckInfoList[i].BusID = r.U32Unsafe()
		bus.DuckInfoList[i].DuckVolume = r.F32Unsafe()
		bus.DuckInfoList[i].FadeOutTime = r.I32Unsafe()
		bus.DuckInfoList[i].FadeInTime = r.I32Unsafe()
		bus.DuckInfoList[i].EnumFadeCurve = wwise.InterpCurveType(r.U8Unsafe())
		bus.DuckInfoList[i].TargetProp = wwise.PropType(r.U8Unsafe())
	}

	ParseFxChunk(r, &bus.BusFxParam.FxChunk, v)

	if v <= 145 {
		bus.BusFxParam.FxID_0 = r.U32Unsafe()
		bus.BusFxParam.IsShareSet_0 = r.U8Unsafe()
	}

	if v <= 145 {
		bus.OverrideAttachmentParams = r.U8Unsafe()
	}

	bus.BusFxMetadataParam.FxChunkMetadataItems = make([]wwise.FxChunkMetadataItem, r.U8Unsafe())
	for i := range bus.BusFxMetadataParam.FxChunkMetadataItems {
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].UniqueFxIndex = r.U8Unsafe()
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].FxId = r.U32Unsafe()
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].BitIsShareSet = r.U8Unsafe()
	}

	ParseRTPC(r, &bus.BusRTPC, v)
	ParseStateProp(r, &bus.StateProp, v)
	ParseStateGroup(r, &bus.StateGroup, v)

	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return &bus
}

func ParseAuxBus(size uint32, r *wio.Reader, v int) *wwise.AuxBus {
	assert.Equal(0, r.Pos(), "Aux Bus parser position doesn't start at 0.")
	begin := r.Pos()
	bus := wwise.AuxBus{
		Id: r.U32Unsafe(),
		OverrideBusId: r.U32Unsafe(),
	}
	if bus.OverrideBusId == 0 {
		bus.DeviceShareSetID = r.U32Unsafe()
	}
	ParsePropBundle(r, &bus.PropBundle, v)
	ParsePositioningParam(r, &bus.PositioningParam, v)
	ParseAuxParam(r, &bus.AuxParam, v)

	bus.VirtualBehaviorBitVector = r.U8Unsafe()
	bus.MaxNumInstance = r.U16Unsafe()
	bus.ChannelConf = r.U32Unsafe()
	bus.HDRBitVector = r.U8Unsafe()
	bus.RecoveryTime = r.I32Unsafe()
	bus.MaxDuckVolume = r.F32Unsafe()

	bus.DuckInfoList = make([]wwise.DuckInfo, r.U32Unsafe())
	for i := range bus.DuckInfoList {
		bus.DuckInfoList[i].BusID = r.U32Unsafe()
		bus.DuckInfoList[i].DuckVolume = r.F32Unsafe()
		bus.DuckInfoList[i].FadeOutTime = r.I32Unsafe()
		bus.DuckInfoList[i].FadeInTime = r.I32Unsafe()
		bus.DuckInfoList[i].EnumFadeCurve = wwise.InterpCurveType(r.U8Unsafe())
		bus.DuckInfoList[i].TargetProp = wwise.PropType(r.U8Unsafe())
	}

	ParseFxChunk(r, &bus.BusFxParam.FxChunk, v)

	if v <= 145 {
		bus.BusFxParam.FxID_0 = r.U32Unsafe()
		bus.BusFxParam.IsShareSet_0 = r.U8Unsafe()
	}

	if v <= 145 {
		bus.OverrideAttachmentParams = r.U8Unsafe()
	}

	bus.BusFxMetadataParam.FxChunkMetadataItems = make([]wwise.FxChunkMetadataItem, r.U8Unsafe())
	for i := range bus.BusFxMetadataParam.FxChunkMetadataItems {
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].UniqueFxIndex = r.U8Unsafe()
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].FxId = r.U32Unsafe()
		bus.BusFxMetadataParam.FxChunkMetadataItems[i].BitIsShareSet = r.U8Unsafe()
	}

	ParseRTPC(r, &bus.BusRTPC, v)
	ParseStateProp(r, &bus.StateProp, v)
	ParseStateGroup(r, &bus.StateGroup, v)

	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return &bus
}

func ParseAudioDevice(size uint32, r *wio.Reader, v int) {}
