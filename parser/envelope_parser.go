package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseModulator(size uint32, r *wio.Reader, v int) *wwise.Modulator {
	assert.Equal(0, r.Pos(), "Envelope modulator parser position doesn't start at position 0.")
	begin := r.Pos()
	e := wwise.Modulator{ Id: r.U32Unsafe() }
	ParsePropBundle(r, &e.PropBundle, v)
	ParseRangePropBundle(r, &e.RangePropBundle, v)
	ParseRTPC(r, &e.RTPC, v)
	e.PropBundle.Modulator = true
	e.RangePropBundle.Modulator = true
	e.RTPC.Modulator = true
	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return &e
}

func ParseLFOModulator(size uint32, r *wio.Reader, v int) *wwise.Modulator {
	e := ParseModulator(size, r, v)
	e.ModulatorType = wwise.HircTypeLFOModulator
	return e
}

func ParseEnvelopeModulator(size uint32, r *wio.Reader, v int) *wwise.Modulator {
	e := ParseModulator(size, r, v)
	e.ModulatorType = wwise.HircTypeEnvelopeModulator
	return e
}

func ParseTimeModulator(size uint32, r *wio.Reader, v int) *wwise.Modulator {
	e := ParseModulator(size, r, v)
	e.ModulatorType = wwise.HircTypeTimeModulator
	return e
}
