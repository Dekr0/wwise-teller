package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParsePluginParam(r *wio.Reader, p *wwise.PluginParam, pluginId uint32) {
	p.PluginParamSize = r.U32Unsafe()
	p.PluginParamData = &wwise.FxPlaceholder{Data: []byte{}}
	if p.PluginParamSize > 0 {
		switch pluginId {
		case 0x00690003:
			p.PluginParamData = &wwise.ParametricEQ{}
			ParseParametricEQ(r, p.PluginParamSize, p.PluginParamData.(*wwise.ParametricEQ))
		case 0x006C0003:
			p.PluginParamData = &wwise.Compressor{}
			ParseCompressor(r, p.PluginParamSize, p.PluginParamData.(*wwise.Compressor))
		case 0x006D0003:
			p.PluginParamData = &wwise.Expander{}
			ParseExpander(r, p.PluginParamSize, p.PluginParamData.(*wwise.Expander))
		case 0x006E0003:
			p.PluginParamData = &wwise.PeakLimiter{}
			ParsePeakLimiter(r, p.PluginParamSize, p.PluginParamData.(*wwise.PeakLimiter))
		case 0x00810003:
			p.PluginParamData = &wwise.MeterFX{}
			ParseMeterFX(r, p.PluginParamSize, p.PluginParamData.(*wwise.MeterFX))
		case 0x008B0003:
			p.PluginParamData = &wwise.GainFX{}
		 	ParseGainFX(r, p.PluginParamSize, p.PluginParamData.(*wwise.GainFX))
		default:
			p.PluginParamData = &wwise.FxPlaceholder{
				Data: r.ReadNUnsafe(uint64(p.PluginParamSize), 0),
			}
		}
	}
}

func ParseParametricEQ(r *wio.Reader, size uint32, p *wwise.ParametricEQ) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.EQBand = [3]wwise.ParametricEQBand(make([]wwise.ParametricEQBand, 3, 3))
	for i := range p.EQBand {
		p.EQBand[i].FilterType = wwise.EQFilterType(r.U32Unsafe())
		p.EQBand[i].Gain = r.F32Unsafe()
		p.EQBand[i].Frequency = r.F32Unsafe()
		p.EQBand[i].QFactor = r.F32Unsafe()
		p.EQBand[i].OnOff = r.U8Unsafe()
	}
	p.OutputLevel = r.F32Unsafe()
	p.ProcessLFE = r.U8Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}

func ParseMeterFX(r *wio.Reader, size uint32, p *wwise.MeterFX) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.Attack = r.F32Unsafe()
	p.Release = r.F32Unsafe()
	p.Min = r.F32Unsafe()
	p.Max = r.F32Unsafe()
	p.Hold = r.F32Unsafe()
	if size >= 0x1c {
		i := r.U8Unsafe()
		p.InfiniteHold = &i
	}
	if size != 0x19 {
		mode := wwise.MeterMode(r.U8Unsafe())
		p.Mode = &mode
	}
	if size > 0x1A {
		scope := wwise.MeterScope(r.U8Unsafe())
		p.Scope = &scope
	}
	p.ApplyDownstreamVolume = r.U8Unsafe()
	p.GameParamID = r.U32Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}

func ParsePeakLimiter(r *wio.Reader, size uint32, p *wwise.PeakLimiter) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.Threshold = r.F32Unsafe()
	p.Ratio = r.F32Unsafe()
	p.LookAhead = r.F32Unsafe()
	p.Release = r.F32Unsafe()
	p.OutputLevel = r.F32Unsafe()
	p.ProcessLFE = r.U8Unsafe()
	p.ChannelLink = r.U8Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}

func ParseGainFX(r *wio.Reader, size uint32, p *wwise.GainFX) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.FullbainGain = r.F32Unsafe()
	p.LFEGain = r.F32Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}

func ParseCompressor(r *wio.Reader, size uint32, p *wwise.Compressor) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.Threshold = r.F32Unsafe()
	p.Ratio = r.F32Unsafe()
	p.Attack = r.F32Unsafe()
	p.Release = r.F32Unsafe()
	p.OutputGain = r.F32Unsafe()
	p.ProcessLFE = r.U8Unsafe()
	p.ChannelLink = r.U8Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}

func ParseExpander(r *wio.Reader, size uint32, p *wwise.Expander) {
	begin := r.Pos()
	expectedEnd := begin + uint64(size)
	p.Threshold = r.F32Unsafe()
	p.Ratio = r.F32Unsafe()
	p.Attack = r.F32Unsafe()
	p.Release = r.F32Unsafe()
	p.OutputGain = r.F32Unsafe()
	p.ProcessLFE = r.U8Unsafe()
	p.ChannelLink = r.U8Unsafe()
	p.Data = r.ReadNUnsafe(expectedEnd - r.Pos(), 0)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
}
