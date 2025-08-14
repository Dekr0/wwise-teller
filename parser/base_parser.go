package parser

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseBaseParam(r *wio.Reader, v int) *wwise.BaseParameter {
	b := wwise.BaseParameter{}
	b.BitIsOverrideParentFx = r.U8Unsafe()
	ParseFxChunk(r, &b.FxChunk, v)
	ParseFxChunkMetadata(r, &b.FxChunkMetadata, v)
	if v <= 145 {
		b.BitOverrideAttachmentParams = r.U8Unsafe()
	}
	b.OverrideBusId = r.U32Unsafe()
	b.DirectParentId = r.U32Unsafe()
	b.ByBitVectorA = r.U8Unsafe()
	ParsePropBundle(r, &b.PropBundle, v)
	ParseRangePropBundle(r, &b.RangePropBundle, v)
	ParsePositioningParam(r, &b.PositioningParam, v)
	ParseAuxParam(r, &b.AuxParam, v)
	ParseAdvanceSetting(r, &b.AdvanceSetting, v)
	ParseStateProp(r, &b.StateProp, v)
	ParseStateGroup(r, &b.StateGroup, v)
	ParseRTPC(r, &b.RTPC, v)
	return &b
}

func ParsePropBundle(r *wio.Reader, p *wwise.PropBundle, v int) {
	CProps := r.U8Unsafe()
	p.PropValues = make([]wwise.PropValue, CProps)
	for i := range CProps {
		p.PropValues[i].P = r.U8Unsafe()
	}
	for i := range CProps {
		p.PropValues[i].V = r.ReadNUnsafe(4, 0)
	}
}

func ParseStatePropBundle(r *wio.Reader, p *wwise.StatePropBundle, v int) {
	CProps := r.U16Unsafe()
	p.StatePropValues = make([]wwise.StatePropValue, CProps)
	for i := range CProps {
		p.StatePropValues[i].P = wwise.StatePropType(r.U16Unsafe())
	}
	for i := range CProps {
		p.StatePropValues[i].V = r.ReadNUnsafe(4, 0)
	}
}

func ParseRangePropBundle(r *wio.Reader, p *wwise.RangePropBundle, v int) {
	CProps := r.U8Unsafe()
	p.RangeValues = make([]wwise.RangeValue, CProps)
	for i := range p.RangeValues {
		p.RangeValues[i].P = r.U8Unsafe()
	}
	for i := range p.RangeValues {
		p.RangeValues[i].Min = r.ReadNUnsafe(4, 0)
		p.RangeValues[i].Max = r.ReadNUnsafe(4, 0)
	}
}

func ParsePositioningParam(r *wio.Reader, p *wwise.PositioningParam, v int) {
	p.BitsPositioning = r.U8Unsafe()
	p.FallbackBitsPositioning = p.BitsPositioning
	if !p.OverrideParentAndHasListenerRelativeRouting() {
		return
	}
	p.Bits3D = r.U8Unsafe()
	p.FallbackBits3D = p.Bits3D
	if !p.Has3DAutomation() {
		return
	}
	p.PathMode = r.U8Unsafe()
	p.FallbackPathMode = p.PathMode
	p.TransitionTime = r.I32Unsafe()
	p.FallbackTransitionTime = p.TransitionTime
	NumPositionVertices := r.U32Unsafe()
	p.PositionVertices = make([]wwise.PositionVertex, NumPositionVertices)
	for i := range p.PositionVertices {
		p.PositionVertices[i].X = r.F32Unsafe()
		p.PositionVertices[i].Y = r.F32Unsafe()
		p.PositionVertices[i].Z = r.F32Unsafe()
		p.PositionVertices[i].Duration = r.I32Unsafe()
	}
	p.FallbackPositionVertices = slices.Clone(p.PositionVertices)
	NumPositionPlayListItem := r.U32Unsafe()
	p.PositionPlayListItems = make([]wwise.PositionPlayListItem, NumPositionPlayListItem)
	for i := range p.PositionPlayListItems {
		p.PositionPlayListItems[i].UniqueVerticesOffset = r.U32Unsafe()
		p.PositionPlayListItems[i].INumVertices =  r.U32Unsafe()
	}
	p.FallbackPositionPlayListItems = slices.Clone(p.FallbackPositionPlayListItems)
	p.Ak3DAutomationParams = make([]wwise.Ak3DAutomationParam, NumPositionPlayListItem)
	for i := range p.Ak3DAutomationParams {
		p.Ak3DAutomationParams[i].XRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].YRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].ZRange = r.F32Unsafe()
	}
	p.FallbackAk3DAutomationParams = slices.Clone(p.Ak3DAutomationParams)
}

func ParseAuxParam(r *wio.Reader, a *wwise.AuxParam, v int) {
	a.AuxBitVector = r.U8Unsafe()
	if a.HasAux() {
		a.AuxIds[0] = r.U32Unsafe()
		a.AuxIds[1] = r.U32Unsafe()
		a.AuxIds[2] = r.U32Unsafe()
		a.AuxIds[3] = r.U32Unsafe()
		a.RestoreAuxIds[0] = a.AuxIds[0]
		a.RestoreAuxIds[1] = a.AuxIds[1]
		a.RestoreAuxIds[2] = a.AuxIds[2]
		a.RestoreAuxIds[3] = a.AuxIds[3]
	}
	a.ReflectionAuxBus = r.U32Unsafe()
	a.RestoreReflectionAuxBus = a.ReflectionAuxBus
}

func ParseAdvanceSetting(r *wio.Reader, a *wwise.AdvanceSetting, v int) {
	a.AdvanceSettingBitVector = r.U8Unsafe()
	a.VirtualQueueBehavior = r.U8Unsafe()
	a.MaxNumInstance = r.U16Unsafe()
	a.BelowThresholdBehavior =  r.U8Unsafe()
	a.HDRBitVector = r.U8Unsafe()
}

func ParseStateProp(r *wio.Reader, s *wwise.StateProp, v int) {
	s.NumStateProps = r.VarUnsafe()
	s.StatePropItems = make([]wwise.StatePropItem, s.NumStateProps.Value, s.NumStateProps.Value)
	for i := range s.StatePropItems {
		s.StatePropItems[i].PropertyId = r.VarUnsafe()
		s.StatePropItems[i].AccumType = wwise.RTPCAccumType(r.U8Unsafe())
		s.StatePropItems[i].InDb = r.U8Unsafe()
	}
}

func ParseStateGroup(r *wio.Reader, s *wwise.StateGroup, v int) {
	s.NumStateGroups = r.VarUnsafe()
	s.StateGroupItems = make([]wwise.StateGroupItem, s.NumStateGroups.Value, s.NumStateGroups.Value)
	for i := range s.StateGroupItems {
		item := &s.StateGroupItems[i]
		item.StateGroupID = r.U32Unsafe()
		item.StateSyncType = r.U8Unsafe()
		item.NumStates = r.VarUnsafe()
		item.States = make([]wwise.StateGroupItemState, item.NumStates.Value, item.NumStates.Value)
		for i := range item.States {
			item.States[i].StateID = r.U32Unsafe()
			if v <= 145 {
				item.States[i].StateInstanceID = r.U32Unsafe()
			} else {
				ParseStatePropBundle(r, &item.States[i].StatePropBundle, v)
			}
		}
	}
}

func ParseRTPC(r *wio.Reader, rtpc *wwise.RTPC, v int) {
	NumRTPC := r.U16Unsafe() // NumCurves in > 141
	rtpc.RTPCItems = make([]wwise.RTPCItem, NumRTPC, NumRTPC)
	for i := range rtpc.RTPCItems {
		item := &rtpc.RTPCItems[i]
		item.RTPCID = r.U32Unsafe()
		item.RTPCType = r.U8Unsafe()
		item.RTPCAccum = wwise.RTPCAccumType(r.U8Unsafe())
		item.ParamID = r.VarUnsafe()
		item.RTPCCurveID = r.U32Unsafe()
		item.Scaling = wwise.CurveScalingType(r.U8Unsafe())
		NumRTPCGraphPoints := r.U16Unsafe()
		item.RTPCGraphPointsX = make([]float32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		item.RTPCGraphPointsY = make([]float32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		item.RTPCGraphPointsInterp = make([]uint32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		ParseRTPCGraphPoints(
			r,
			item.RTPCGraphPointsX,
			item.RTPCGraphPointsY,
			item.RTPCGraphPointsInterp,
			v,
		)
	}
}

func ParseRTPCGraphPoints(
	r *wio.Reader,
	xs []float32,
	ys []float32,
	interps []uint32,
	v int,
) {
	for i := range xs {
		xs[i] = r.F32Unsafe()
		ys[i] = r.F32Unsafe()
		interps[i] = r.U32Unsafe()
	}
}
