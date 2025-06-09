// TODO:
// - Make proper enum type
package wwise

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type InterpCurveType uint8 

const (
	InterpCurveTypeLog3      InterpCurveType = 0
	InterpCurveTypeSine      InterpCurveType = 1
	InterpCurveTypeLog1      InterpCurveType = 2
	InterpCurveTypeInvSCurve InterpCurveType = 3
	InterpCurveTypeLinear    InterpCurveType = 4
	InterpCurveTypeSCurve    InterpCurveType = 5
	InterpCurveTypeExp1      InterpCurveType = 6
	InterpCurveTypInvSine    InterpCurveType = 7
	InterpCurveTypeExp3      InterpCurveType = 8
	InterpCurveTypeConst     InterpCurveType = 9
	InterpCurveTypeCount     InterpCurveType = 10
)
var InterpCurveTypeName []string = []string{
  	"Logarithmic (Base 3)",
  	"Sine (Constant Power Fade Out)",
  	"Logarithmic (Base 1.41)",
  	"Inverted S-Curve",
  	"Linear",
  	"S-Curve",
  	"Exponential (Base 1.41)",
  	"Sine (Constant Power Fade In)",
  	"Exponential (Base 3)",
  	"Constant",
}

// Target toward v144 but HD2 is using 141???
const RTPCTypeCount = 5
var RTPCTypeName []string = []string{
	"Game Parameter",
	"MIDI Parameter",
	"Switch",
	"State",
	"Modulator",
	// "Count",
	// "Max Number",
}


type RTPCParameterType uint8 // var (at least 1 byte / 8 bits)

const (
	RTPCParameterTypeVolume RTPCParameterType = 0 
	RTPCParameterTypeLFE RTPCParameterType = 1 
	RTPCParameterTypePitch RTPCParameterType = 2 
	RTPCParameterTypeLPF RTPCParameterType = 3 
	RTPCParameterTypeHPF RTPCParameterType = 4 
	RTPCParameterTypeBusVolume RTPCParameterType = 5 
	RTPCParameterTypeInitialDelay RTPCParameterType = 6 
	RTPCParameterTypeMakeUpGain RTPCParameterType = 7 
	RTPCParameterTypeDeprecatedFeedbackVolume RTPCParameterType = 8 // #140~~ Deprecated_RTPC_FeedbackVolume
	RTPCParameterTypeDeprecatedFeedbackLowpass RTPCParameterType = 9 // #140~~ Deprecated_RTPC_FeedbackLowpass
	RTPCParameterTypeDeprecatedFeedbackPitch RTPCParameterType = 10 // ##140~~ Deprecated_RTPC_FeedbackPitch
	RTPCParameterTypeMidiTransposition RTPCParameterType = 11 
	RTPCParameterTypeMidiVelocityOffset RTPCParameterType = 12 
	RTPCParameterTypePlaybackSpeed RTPCParameterType = 13 
	RTPCParameterTypeMuteRatio RTPCParameterType = 14 
	RTPCParameterTypePlayMechanismSpecialTransitionsValue RTPCParameterType = 15 
	RTPCParameterTypeMaxNumberInstances RTPCParameterType = 16 
	RTPCParameterTypePriority RTPCParameterType = 17 // OVERRIDABLE_PARAMS_STAR RTPCParameterType
	RTPCParameterTypePositionPANX2D RTPCParameterType = 18 
	RTPCParameterTypePositionPANY2D RTPCParameterType = 19 
	RTPCParameterTypePositionPANX3D RTPCParameterType = 20 
	RTPCParameterTypePositionPANY3D RTPCParameterType = 21 
	RTPCParameterTypePositionPANZ3D RTPCParameterType = 22 
	RTPCParameterTypePositioningTypeBlend RTPCParameterType = 23 
	RTPCParameterTypePositioningDivergenceCenterPCT RTPCParameterType = 24 
	RTPCParameterTypePositioningConeAttenuationONOFF RTPCParameterType = 25 
	RTPCParameterTypePositioningConeAttenuation RTPCParameterType = 26 
	RTPCParameterTypePositioningConeLPF RTPCParameterType = 27 
	RTPCParameterTypePositioningConeHPF RTPCParameterType = 28 
	RTPCParameterTypeBypassFX0 RTPCParameterType = 29 
	RTPCParameterTypeBypassFX1 RTPCParameterType = 30 
	RTPCParameterTypeBypassFX2 RTPCParameterType = 31 
	RTPCParameterTypeBypassFX3 RTPCParameterType = 32 
	RTPCParameterTypeBypassAllFX RTPCParameterType = 33 
	RTPCParameterTypeHDRBusThreshold RTPCParameterType = 34 
	RTPCParameterTypeHDRBusReleaseTime RTPCParameterType = 35 
	RTPCParameterTypeHDRBusRatio RTPCParameterType = 36 
	RTPCParameterTypeHDRActiveRange RTPCParameterType = 37 
	RTPCParameterTypeGameAuxSendVolume RTPCParameterType = 38 
	RTPCParameterTypeUserAuxSendVolume0 RTPCParameterType = 39 
	RTPCParameterTypeUserAuxSendVolume1 RTPCParameterType = 40 
	RTPCParameterTypeUserAuxSendVolume2 RTPCParameterType = 41 
	RTPCParameterTypeUserAuxSendVolume3 RTPCParameterType = 42 
	RTPCParameterTypeOutputBusVolume RTPCParameterType = 43 
	RTPCParameterTypeOutputBusHPF RTPCParameterType = 44 
	RTPCParameterTypeOutputBusLPF RTPCParameterType = 45 
	RTPCParameterTypePositioningEnableAttenuation RTPCParameterType = 46 
	RTPCParameterTypeReflectionsVolume RTPCParameterType = 47 
	RTPCParameterTypeUserAuxSendLPF0 RTPCParameterType = 48 
	RTPCParameterTypeUserAuxSendLPF1 RTPCParameterType = 49 
	RTPCParameterTypeUserAuxSendLPF2 RTPCParameterType = 50 
	RTPCParameterTypeUserAuxSendLPF3 RTPCParameterType = 51 
	RTPCParameterTypeUserAuxSendHPF0 RTPCParameterType = 52 
	RTPCParameterTypeUserAuxSendHPF1 RTPCParameterType = 53 
	RTPCParameterTypeUserAuxSendHPF2 RTPCParameterType = 54 
	RTPCParameterTypeUserAuxSendHPF3 RTPCParameterType = 55 
	RTPCParameterTypeGameAuxSendLPF RTPCParameterType = 56 
	RTPCParameterTypeGameAuxSendHPF RTPCParameterType = 57 
	RTPCParameterTypePositionPANZ2D RTPCParameterType = 58 
	RTPCParameterTypeBypassAllMetadata RTPCParameterType = 59 
	RTPCParameterTypeCount = 60
)

var RTPCParameterIDName []string = []string{
	"Volume",
	"LFE",
	"Pitch",
	"LPF",
	"HPF",
	"Bus Volume",
	"Initial Delay",
	"Make Up Gain",
	"Deprecated Feedback Volume", // #140~~ Deprecated_RTPC_FeedbackVolume
	"Deprecated Feedback Lowpass", // #140~~ Deprecated_RTPC_FeedbackLowpass
	"Deprecated Feedback Pitch", // ##140~~ Deprecated_RTPC_FeedbackPitch
	"Midi Transposition",
	"Midi Velocity Offset",
	"Playback Speed",
	"Mute Ratio",
	"Play Mechanism Special Transitions Value",
	"Max Number Instances",
	"Priority", // OVERRIDABLE_PARAMS_START
	"Position PAN X 2D",
	"Position PAN Y 2D",
	"Position PAN X 3D",
	"Position PAN Y 3D",
	"Position PAN Z 3D",
	"Positioning Type Blend",
	"Positioning Divergence Center PCT",
	"Positioning Cone Attenuation ON OFF",
	"Positioning Cone Attenuation",
	"Positioning Cone LPF",
	"Positioning Cone HPF",
	"Bypass FX0",
	"Bypass FX1",
	"Bypass FX2",
	"Bypass FX3",
	"Bypass All FX",
	"HDR Bus Threshold",
	"HDR Bus Release Time",
	"HDR Bus Ratio",
	"HDR Active Range",
	"Game Aux Send Volume",
	"User Aux Send Volume 0",
	"User Aux Send Volume 1",
	"User Aux Send Volume 2",
	"User Aux Send Volume 3",
	"Output Bus Volume",
	"Output Bus HPF",
	"Output Bus LPF",
	"Positioning Enable Attenuation",
	"Reflections Volume",
	"User Aux Send LPF 0",
	"User Aux Send LPF 1",
	"User Aux Send LPF 2",
	"User Aux Send LPF 3",
	"User Aux Send HPF 0",
	"User Aux Send HPF 1",
	"User Aux Send HPF 2",
	"User Aux Send HPF 3",
	"Game Aux Send LPF",
	"Game Aux Send HPF",
	"Position PAN Z 2D",
	"Bypass All Metadata",
	// #0x3C: "MaxNumRTPC,

	// 0x3D: "Unknown/Custom?", #AC Valhalla
	// 0x3E: "Unknown/Custom?", #AC Valhalla (found near "DB" scaling, some volume?)
	// 0x3F: "Unknown/Custom?", #AC Valhalla
}

type RTPCAccumType uint8

const (
	RTPCAccumTypeNone RTPCAccumType = 0
	RTPCAccumTypeExclusive RTPCAccumType = 1
	RTPCAccumTypeAdditive RTPCAccumType = 2
	RTPCAccumTypeMultiply RTPCAccumType = 3
	RTPCAccumTypeBoolean RTPCAccumType = 4
	RTPCAccumTypeMaximum RTPCAccumType = 5
	RTPCAccumTypeFilter RTPCAccumType = 6
	RTPCAccumTypeCount RTPCAccumType = 7
)

var RTPCAccumTypeName []string = []string{
	"None",
	"Exclusive",
	"Additive",
	"Multiply",
	"Boolean",
	"Maximum",
	"Filter",
	// "Max Number / Count",
}

type CurveScalingType uint8

const (
	CurveScalingTypeNone        CurveScalingType = 0
	CurveScalingTypeUnsupported CurveScalingType = 1
	CurveScalingTypeDb          CurveScalingType = 2
	CurveScalingTypeLog         CurveScalingType = 3
	CurveScalingTypeDbToLin     CurveScalingType = 4
	CurveScalingTypeCount       CurveScalingType = 5
)

var CurveScalingTypeName []string = []string{
	"None",
	"Unsupported",
	"dB",
	"Log",
	"dB To Lin",
}

type RTPC struct {
	// NumRTPC uint16   // <= 141
	// NumCurves uint16 // < 141
	RTPCItems []RTPCItem
}

func (r *RTPC) RemoveRTPCItem(i int) {
	r.RTPCItems = slices.Delete(r.RTPCItems, i, i + 1)
}

func NewRTPC() *RTPC {
	return &RTPC{[]RTPCItem{}}
}

func (r *RTPC) Encode() []byte {
	size := r.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(uint16(len(r.RTPCItems)))
	for _, i := range r.RTPCItems {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (r *RTPC) Size() uint32 {
	size := uint32(2)
	for _, i := range r.RTPCItems {
		size += i.Size()
	}
	return size
}

type RTPCItem struct {
	RTPCID      uint32 // tid
	RTPCType    uint8 // U8x
	RTPCAccum   RTPCAccumType // U8x
	ParamID     RTPCParameterType
	RTPCCurveID uint32 // sid
	Scaling     CurveScalingType // U8x
	// NumRTPCGraphPoints / ulSize uint16 // u16
	RTPCGraphPoints []RTPCGraphPoint 
}

func NewRTPCItem() *RTPCItem {
	return &RTPCItem{0, 0, 0, 0, 0, 0, []RTPCGraphPoint{}}
}

func (r *RTPCItem) Encode() []byte {
	size := r.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(r.RTPCID)
	w.AppendByte(r.RTPCType)
	w.Append(r.RTPCAccum)
	w.Append(r.ParamID)
	w.Append(r.RTPCCurveID)
	w.Append(r.Scaling)
	w.Append(uint16(len(r.RTPCGraphPoints)))
	for _, i := range r.RTPCGraphPoints {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (r *RTPCItem) Size() uint32 {
	return uint32(4 + 1 + 1 + 1 + 4 + 1 + 2 + len(r.RTPCGraphPoints) * SizeOfRTPCGraphPoint)
}

const RTPCInterpSampleRate = 128

const SizeOfRTPCGraphPoint = 12
type RTPCGraphPoint struct {
	From           float32 // f32 
	To             float32 // f32
	Interp         uint32 // U32
	SamplePointsX  []float32
	SamplePointsY  []float32
}

func (r *RTPCGraphPoint) Sample() {}

func (r *RTPCGraphPoint) Encode() []byte {
	w := wio.NewWriter(SizeOfRTPCGraphPoint)
	w.Append(r.From)
	w.Append(r.To)
	w.Append(r.Interp)
	return w.BytesAssert(SizeOfRTPCGraphPoint)
}
