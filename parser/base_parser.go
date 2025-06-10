package parser

import (
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseBaseParam(r *wio.Reader) *wwise.BaseParameter {
	b := wwise.BaseParameter{}
	b.BitIsOverrideParentFx = r.U8Unsafe()
	ParseFxChunk(r, &b.FxChunk)
	ParseFxChunkMetadata(r, &b.FxChunkMetadata)
	b.BitOverrideAttachmentParams = r.U8Unsafe()
	b.OverrideBusId = r.U32Unsafe()
	b.DirectParentId = r.U32Unsafe()
	b.ByBitVectorA = r.U8Unsafe()
	ParsePropBundle(r, &b.PropBundle)
	ParseRangePropBundle(r, &b.RangePropBundle)
	ParsePositioningParam(r, &b.PositioningParam)
	ParseAuxParam(r, &b.AuxParam)
	ParseAdvanceSetting(r, &b.AdvanceSetting)
	ParseStateProp(r, &b.StateProp)
	ParseStateGroup(r, &b.StateGroup)
	ParseRTPC(r, &b.RTPC)
	return &b
}

func ParsePropBundle(r *wio.Reader, p *wwise.PropBundle) {
	CProps := r.U8Unsafe()
	p.PropValues = make([]wwise.PropValue, CProps)
	for i := range CProps {
		p.PropValues[i].P = wwise.PropType(r.U8Unsafe())
	}
	for i := range CProps {
		p.PropValues[i].V = r.ReadNUnsafe(4, 0)
	}
}

func ParseRangePropBundle(r *wio.Reader, p *wwise.RangePropBundle) {
	CProps := r.U8Unsafe()
	p.RangeValues = make([]wwise.RangeValue, CProps)
	for i := range p.RangeValues {
		p.RangeValues[i].P = wwise.PropType(r.U8Unsafe())
	}
	for i := range p.RangeValues {
		p.RangeValues[i].Min = r.ReadNUnsafe(4, 0)
		p.RangeValues[i].Max = r.ReadNUnsafe(4, 0)
	}
}

func ParsePositioningParam(r *wio.Reader, p *wwise.PositioningParam) {
	p.BitsPositioning = r.U8Unsafe()
	if !p.HasPositioningAnd3D() {
		return
	}
	p.Bits3D = r.U8Unsafe()
	if !p.HasAutomation() {
		return
	}
	p.PathMode = r.U8Unsafe()
	p.TransitionTime = r.I32Unsafe()
	NumPositionVertices := r.U32Unsafe()
	p.PositionVertices = make([]wwise.PositionVertex, NumPositionVertices)
	for i := range p.PositionVertices {
		p.PositionVertices[i].X = r.F32Unsafe()
		p.PositionVertices[i].Y = r.F32Unsafe()
		p.PositionVertices[i].Z = r.F32Unsafe()
		p.PositionVertices[i].Duration = r.I32Unsafe()
	}
	NumPositionPlayListItem := r.U32Unsafe()
	p.PositionPlayListItems = make([]wwise.PositionPlayListItem, NumPositionPlayListItem)
	for i := range p.PositionPlayListItems {
		p.PositionPlayListItems[i].UniqueVerticesOffset = r.U32Unsafe()
		p.PositionPlayListItems[i].INumVertices =  r.U32Unsafe()
	}
	p.Ak3DAutomationParams = make([]wwise.Ak3DAutomationParam, NumPositionPlayListItem)
	for i := range p.Ak3DAutomationParams {
		p.Ak3DAutomationParams[i].XRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].YRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].ZRange = r.F32Unsafe()
	}
}

func ParseAuxParam(r *wio.Reader, a *wwise.AuxParam) {
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

func ParseAdvanceSetting(r *wio.Reader, a *wwise.AdvanceSetting) {
	a.AdvanceSettingBitVector = r.U8Unsafe()
	a.VirtualQueueBehavior = r.U8Unsafe()
	a.MaxNumInstance = r.U16Unsafe()
	a.BelowThresholdBehavior =  r.U8Unsafe()
	a.HDRBitVector = r.U8Unsafe()
}

func ParseStateProp(r *wio.Reader, s *wwise.StateProp) {
	NumStateProps := r.U8Unsafe()
	s.StatePropItems = make([]wwise.StatePropItem, NumStateProps)
	for i := range s.StatePropItems {
		s.StatePropItems[i].PropertyId = wwise.RTPCParameterType(r.U8Unsafe())
		s.StatePropItems[i].AccumType = wwise.RTPCAccumType(r.U8Unsafe())
		s.StatePropItems[i].InDb = r.U8Unsafe()
	}
}

func ParseStateGroup(r *wio.Reader, s *wwise.StateGroup) {
	NumStateGroups := r.U8Unsafe()
	s.StateGroupItems = make([]wwise.StateGroupItem, NumStateGroups)
	for i := range s.StateGroupItems {
		item := &s.StateGroupItems[i]
		item.StateGroupID = r.U32Unsafe()
		item.StateSyncType = r.U8Unsafe()
		NumStates := r.U8Unsafe()
		item.States = make([]wwise.StateGroupItemState, NumStates)
		for i := range item.States {
			item.States[i].StateID = r.U32Unsafe()
			item.States[i].StateInstanceID = r.U32Unsafe()
		}
	}
}

func ParseRTPC(r *wio.Reader, rtpc *wwise.RTPC) {
	NumRTPC := r.U16Unsafe()
	rtpc.RTPCItems = make([]wwise.RTPCItem, NumRTPC, NumRTPC)
	for i := range rtpc.RTPCItems {
		item := &rtpc.RTPCItems[i]
		item.RTPCID = r.U32Unsafe()
		item.RTPCType = r.U8Unsafe()
		item.RTPCAccum = wwise.RTPCAccumType(r.U8Unsafe())
		item.ParamID = wwise.RTPCParameterType(r.U8Unsafe())
		item.RTPCCurveID = r.U32Unsafe()
		item.Scaling = wwise.CurveScalingType(r.U8Unsafe())
		NumRTPCGraphPoints := r.U16Unsafe()
		item.RTPCGraphPoints = make([]wwise.RTPCGraphPoint, NumRTPCGraphPoints, NumRTPCGraphPoints)
		ParseRTPCGraphPoints(r, item.RTPCGraphPoints)
	}
}

func ParseRTPCGraphPoints(r *wio.Reader, pts []wwise.RTPCGraphPoint) {
	for i := range pts {
		pts[i].From = r.F32Unsafe()
		pts[i].To = r.F32Unsafe()
		pts[i].Interp = r.U32Unsafe()
	}
}
