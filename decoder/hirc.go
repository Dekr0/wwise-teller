package decoder

import (
	"bytes"
	"context"
	"io"
	"sync/atomic"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

type HierarchyDecoder func(io.Reader, order, u32, *wwise.HIRC, u32)

type DecoderOptHIRC struct {
	NumRoutine   u8
	Exclude    []u8
}

func DecodeHIRC(
	ctx      context.Context, 
	decoder  DecoderOptHIRC,
	r        io.Reader, 
	o        order,
	size     u32, 
	ver      u32, 
) (h *wwise.HIRC, err error) {
	p := 0

	numHirc := uio.U32PT(r, o, &p)

	sem := make(chan struct{}, decoder.NumRoutine)

	h = wwise.NewHIRC(numHirc)

	dispatch := uint32(0)
	parsed := atomic.Uint32{}

	for parsed.Load() < numHirc {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if dispatch >= numHirc {
			continue
		}

		t := uio.U8PT(r, o, &p)
		size := uio.U32PT(r, o, &p)

		buffer := make([]byte, 0, size)
		if _, err = r.Read(buffer); err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buffer)

		var decoder HierarchyDecoder
		switch t {
		case wwise.HircTypeState:
			decoder = DecodeState
		case wwise.HircTypeEvent:
			decoder = DecodeEvent
		}

		select {
		case sem <- struct{}{}:
			go decoder(reader, o, ver, h, size)
		default:
			decoder(reader, o, ver, h, size)
		}
		dispatch++
	}

	return h, nil
}
