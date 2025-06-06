// TODO:
// - Make proper enum type
package wwise

import (
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

const RTPCParameterIDCount = 60
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

const RTPCAccumTypeCount = 7
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

const CurveScalingTypeCount = 5
var CurveScalingTypeName []string = []string{
	"None",
	"Unsupported",
	"dB",
	"Log",
	"dB To Lin",
}

type RTPC struct {
	// NumRTPC uint16 // u16
	RTPCItems []RTPCItem
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
	RTPCID uint32 // tid
	RTPCType uint8 // U8x
	RTPCAccum uint8 // U8x
	ParamID uint8 // var (assume at least 1 byte / 8 bits, can be more)
	RTPCCurveID uint32 // sid
	Scaling uint8 // U8x
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
	w.AppendByte(r.RTPCAccum)
	w.AppendByte(r.ParamID)
	w.Append(r.RTPCCurveID)
	w.AppendByte(r.Scaling)
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
