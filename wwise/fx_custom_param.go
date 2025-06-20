package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)


type PluginParam struct {
	PluginParamSize uint32 // U32
	PluginParamData FxParam
}

func (p *PluginParam) Encode() []byte {
	assert.Equal(
		p.PluginParamSize,
		p.PluginParamData.Size(),
		"Plugin parameter size doesn't equal # of bytes in plugin parameter data",
	)
	size := 4 + p.PluginParamData.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.PluginParamSize)
	w.AppendBytes(p.PluginParamData.Encode())
	return w.BytesAssert(int(size))
}

func (p *PluginParam) Size() uint32 {
	return 4 + p.PluginParamData.Size()
}

type FxParam interface {
	Encode() []byte
	Size()     uint32
}

type FxPlaceholder struct {
	Data []byte
}

func (f *FxPlaceholder) Encode() []byte {
	return f.Data
}

func (f *FxPlaceholder) Size() uint32 {
	return uint32(len(f.Data))
}

type SourceSine struct {
	Frequency   float32
	Gain        float32
	Duration    float32
	ChannelMask uint32
	Data        []byte
}

func (f *SourceSine) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Frequency)
	w.Append(f.Gain)
	w.Append(f.Duration)
	w.Append(f.ChannelMask)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *SourceSine) Size() uint32 {
	return 4 * 4 + uint32(len(f.Data))
}

type SourceSlience struct {
	Duration              float32
	RandomizedLengthMinus float32
	RandomizedLengthPlus  float32
	Data                  []byte
}

func (f *SourceSlience) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Duration)
	w.Append(f.RandomizedLengthMinus)
	w.Append(f.RandomizedLengthPlus)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *SourceSlience) Size() uint32 {
	return 3 * 4 + uint32(len(f.Data))
}

type ToneGen struct {
	Gain         float32
	StartFreq    float32
	StopFreq     float32
	StartFreqMin float32
	Data         []byte
}

type ParametricEQ struct {
	EQBand      [3]ParametricEQBand
	OutputLevel    float32
	ProcessLFE     uint8
	Data         []byte
}

func (f *ParametricEQ) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	for _, b := range f.EQBand {
		w.Append(b)
	}
	w.Append(f.OutputLevel)
	w.Append(f.ProcessLFE)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *ParametricEQ) Size() uint32 {
	return 56 + uint32(len(f.Data))
}

type EQFilterType uint32
const (
	EQFilterTypeLowPass   EQFilterType = 0
	EQFilterTypeHiPass    EQFilterType = 1
	EQFilterTypeBandPass  EQFilterType = 2
	EQFilterTypeNotch     EQFilterType = 3
	EQFilterTypeLowShelf  EQFilterType = 4
	EQFilterTypeHiShelf   EQFilterType = 5
	EQFilterTypePeakingEQ EQFilterType = 6
	EQFilterTypeCount     EQFilterType = 7
)
var EQFilterNames []string = []string{
	"Low Pass",
	"High Pass",
	"Band Pass",
	"Notch",
	"Low Shelf",
	"High Shelf",
	"Peaking EQ",
}

const SizeOfFxParametricEQBand = 17
type ParametricEQBand struct {
	FilterType EQFilterType
	Gain       float32
	Frequency  float32
	QFactor    float32
	OnOff      uint8
}

type MeterMode uint8
const (
	MeterModePeak  MeterMode = 0
	MeterModeRMS   MeterMode = 1
	MeterModeCount MeterMode = 2
)

var MeterModeNames []string = []string{
	"Peak", "RMS",
}

type MeterScope uint8
const (
	MeterScopeGlobal     MeterScope = 0
	MeterScopeGameObject MeterScope = 1
	MeterScopeCount      MeterScope = 2
)

var MeterScopeNames []string = []string{
	"Global", "Game Object",
}

type MeterFX struct {
	Attack                float32 // RTPC
	Release               float32 // RTPC
	Min                   float32 // RTPC
	Max                   float32 // RTPC
	Hold                  float32 // RTPC

	InfiniteHold         *uint8 // RTPC
	Mode                 *MeterMode // Non RTPC
	Scope                *MeterScope // Non RTPC
	ApplyDownstreamVolume uint8 // Non RTPC
	GameParamID           uint32 // Non RTPC
	Data                []byte
}

func (f *MeterFX) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Attack)
	w.Append(f.Release)
	w.Append(f.Min)
	w.Append(f.Max)
	w.Append(f.Hold)
	if f.InfiniteHold != nil {
		w.Append(*f.InfiniteHold)
	}
	if f.Mode != nil {
		w.Append(*f.Mode)
	}
	if f.Scope != nil {
		w.Append(*f.Scope)
	}
	w.Append(f.ApplyDownstreamVolume)
	w.Append(f.GameParamID)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *MeterFX) Size() uint32 {
	size := 25 + uint32(len(f.Data))
	if f.InfiniteHold != nil {
		size += 1
	}
	if f.Mode != nil {
		size += 1
	}
	if f.Scope != nil {
		size += 1
	}
	return size
}

type PeakLimiter struct {
	Threshold   float32 // RTPC
	Ratio       float32 // RTPC
	LookAhead   float32 // Non RTPC
	Release     float32 // RTPC
	OutputLevel float32 // RTPC in db
	ProcessLFE  uint8   // Non RTPC
	ChannelLink uint8   // Non RTPC
	Data        []byte
}

func (f *PeakLimiter) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Threshold)
	w.Append(f.Ratio)
	w.Append(f.LookAhead)
	w.Append(f.Release)
	w.Append(f.OutputLevel)
	w.Append(f.ProcessLFE)
	w.Append(f.ChannelLink)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *PeakLimiter) Size() uint32 {
	return 22 + uint32(len(f.Data))
}

type GainFX struct {
	FullbandGain float32
	LFEGain      float32
	Data         []byte
}

func (f *GainFX) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.FullbandGain)
	w.Append(f.LFEGain)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *GainFX) Size() uint32 {
	return 8 + uint32(len(f.Data))
}

type Compressor struct {
	Threshold   float32;
	Ratio       float32;
	Attack      float32;
	Release     float32;
	OutputGain  float32;
	ProcessLFE  uint8;
	ChannelLink uint8;
	Data        []byte
}

func (f *Compressor) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Threshold)
	w.Append(f.Ratio)
	w.Append(f.Attack)
	w.Append(f.Release)
	w.Append(f.OutputGain)
	w.Append(f.ProcessLFE)
	w.Append(f.ChannelLink)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *Compressor) Size() uint32 {
	return 22 + uint32(len(f.Data))
}

type Expander struct {
	Threshold   float32;
	Ratio       float32;
	Attack      float32;
	Release     float32;
	OutputGain  float32;
	ProcessLFE  uint8;
	ChannelLink uint8;
	Data        []byte
}

func (f *Expander) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Threshold)
	w.Append(f.Ratio)
	w.Append(f.Attack)
	w.Append(f.Release)
	w.Append(f.OutputGain)
	w.Append(f.ProcessLFE)
	w.Append(f.ChannelLink)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *Expander) Size() uint32 {
	return 22 + uint32(len(f.Data))
}

type ConvolutionReverb struct {
	PreDelay         float32
	FrontRearDelay   float32
	StereoWidth      float32
	InputCenterLevel float32
	InputLFELevel    float32
	InputStereoWidth *float32
	FrontLevel       float32
	RearLevel        float32
	CenterLevel      float32
	LFELevel         float32
	DryLevel         float32
	WetLevel         float32
}
