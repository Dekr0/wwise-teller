package wwise

import (
	"github.com/Dekr0/wwise-teller/wio"
)

const CurveInterpolationCount = 10
var CurveInterpolationName []string = []string{
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
