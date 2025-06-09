package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)


type PluginParam struct {
	PluginParamSize uint32 // U32
	PluginParamData []byte
}

func (p *PluginParam) Encode() []byte {
	assert.Equal(
		int(p.PluginParamSize),
		len(p.PluginParamData),
		"Plugin parameter size doesn't equal # of bytes in plugin parameter data",
	)
	size := 4 + len(p.PluginParamData)
	w := wio.NewWriter(uint64(size))
	w.Append(p.PluginParamSize)
	w.AppendBytes(p.PluginParamData)
	return w.BytesAssert(size)
}

func (p *PluginParam) Size() uint32 {
	return uint32(4 + len(p.PluginParamData))
}


type FxSourceSineParam struct {
	Frequency   float32
	Gain        float32
	Duration    float32
	ChannelMask uint32
	Data        []byte
}

func (f *FxSourceSineParam) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Frequency)
	w.Append(f.Gain)
	w.Append(f.Duration)
	w.Append(f.ChannelMask)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *FxSourceSineParam) Size() uint32 {
	return 4 * 4 + uint32(len(f.Data))
}

type FxSourceSlienceParam struct {
	Duration              float32
	RandomizedLengthMinus float32
	RandomizedLengthPlus  float32
	Data                  []byte
}

func (f *FxSourceSlienceParam) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(f.Duration)
	w.Append(f.RandomizedLengthMinus)
	w.Append(f.RandomizedLengthPlus)
	w.AppendBytes(f.Data)
	return w.BytesAssert(int(size))
}

func (f *FxSourceSlienceParam) Size() uint32 {
	return 3 * 4 + uint32(len(f.Data))
}

type FxToneGenParam struct {
	Gain         float32
	StartFreq    float32
	StopFreq     float32
	StartFreqMin float32
}
