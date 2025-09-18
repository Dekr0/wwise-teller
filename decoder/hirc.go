package decoder

import (
	"bytes"
	"context"
	bin "encoding/binary"
	"fmt"
	"io"
	"log/slog"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

type HierarchyDecoder func(io.Reader, order, u32, u32) any

type HircDecodeOption struct {
	Exclude    []u8
}

func AllocDecodeHIRC(
	ctx       context.Context, 
	opt      *HircDecodeOption,
	inReader  io.Reader, 
	o         order,
	size      u32, 
	version   u32, 
) (h *wwise.HIRC, err error) {
	if opt == nil {
		return nil, fmt.Errorf("Must provide HIRC decoder option")
	}

	r := io.LimitReader(inReader, int64(size))

	numHirc, err := uio.U32(r, o)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode # of hierarchies: %w", err)
	}

	h = wwise.AllocHIRC(numHirc)

	dispatch := uint32(0)

	eof := false
	for dispatch < numHirc && !eof {
		select {
		case <- ctx.Done():
			return nil, fmt.Errorf("Hierarchy decoding process cancel: %w", ctx.Err())
		default:
		}

		t, err := uio.U8(r, o)
		if err != nil {
			if err == io.EOF {
				eof = true
				break
			}
			return nil, fmt.Errorf("Failed to decode hierarchy type: %w", err)
		}
		
		size, err := uio.U32(r, o)
		if err != nil {
			if err == io.EOF {
				eof = true
				break
			}
			return nil, fmt.Errorf("Failed to decode hierarchy data size: %w", err)
		}

		buffer := make([]byte, size, size)
		_, err = io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				eof = true
				break
			}
			return nil, fmt.Errorf(
				"Failed to read %d bytes of hierarchy data for buffering decoding: %w", 
				size, err,
			)
		}

		slog.Debug(fmt.Sprintf("Decoding a %s", wwise.GetHircTypeName(t)), "dispath", dispatch, "size", size)

		var decoder HierarchyDecoder
		switch t {
		case wwise.HircTypeState:
			decoder = AllocDecodeState
		case wwise.HircTypeSound:
			decoder = AllocDecodeSound
		case wwise.HircTypeEvent:
			decoder = AllocDecodeEvent
		}

		if decoder == nil {
			reader := bytes.NewReader(buffer)
			var id u32
			if err = bin.Read(reader, o, &id); err != nil {
				return nil, fmt.Errorf(
					"Failed to decode %s hierarchy id at position %d: %w",
					wwise.GetHircTypeName(t), dispatch, err,
				)
			}
			h.Hierarchy.AddEncodedHierarchyNode(id, t, buffer)
			dispatch++
			continue
		}

		reader := bytes.NewReader(buffer)
		res := decoder(reader, o, version, size)
		switch t := res.(type) {
		case *wwise.StateH:
			h.AddState(t.Id, t.StateProps)
		case *wwise.SoundH:
			h.AddSound(t, version)
		case *wwise.EventH:
			h.AddEvent(t.Id, t.EventData)
		}
		dispatch++
	}

	if dispatch < numHirc && eof {
		return nil, fmt.Errorf(
			"Hierarchy decoding process encounter EOF after dispatching %d decoding routine. The total # of hierarchy is %d",
			dispatch, numHirc,
		)
	}

	return h, nil
}
