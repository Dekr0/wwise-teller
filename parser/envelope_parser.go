package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseEnvelopeModulator(size uint32, r *wio.Reader) *wwise.EnvelopeModulator {
	assert.Equal(0, r.Pos(), "Envelope modulator parser position doesn't start at position 0.")
	begin := r.Pos()
	e := wwise.EnvelopeModulator{ Id: r.U32Unsafe() }
	ParsePropBundle(r, &e.PropBundle)
	ParseRangePropBundle(r, &e.RangePropBundle)
	ParseRTPC(r, &e.RTPC)
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
