package decoder

import (
	"bytes"
	"context"
	bin "encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

type HierarchyDecoder func(io.Reader, order, u32, *wwise.HIRC, u32)

type HircDecodeOption struct {
	NumRoutine   u8
	Exclude    []u8
}

type DecoderJob struct {
	DataSize  u32
	Decoder   HierarchyDecoder
	Data    []byte
}

func Decoder(
	ctx      context.Context,
	h       *wwise.HIRC,
	o        order,
	version  u32,
	jobs     <-chan DecoderJob,
	finished *atomic.Uint32,
) {
	for {
		select {
		case <- ctx.Done():
			slog.Info("Decoder exit")
			return
		case j := <- jobs:
			reader := bytes.NewReader(j.Data)
			j.Decoder(reader, o, version, h, j.DataSize)
			finished.Add(1)
		}
	}
}

func DecodeHIRC(
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

	h = wwise.NewHIRC(numHirc)

	decodingCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var jobs chan DecoderJob
	var finished atomic.Uint32
	if opt.NumRoutine > 0 {
		jobs = make(chan DecoderJob)
		for range opt.NumRoutine {
			go Decoder(decodingCtx, h, o, version, jobs, &finished)
		}
	}

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
			decoder = DecodeState
		case wwise.HircTypeEvent:
			decoder = DecodeEvent
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
			wwise.NewEncodedHierarchy(h, id, t, buffer)
			dispatch++
			finished.Add(1)
			continue
		}

		if jobs != nil {
			jobs <- DecoderJob{size, decoder, buffer}
		} else {
			reader := bytes.NewReader(buffer)
			decoder(reader, o, version, h, size)
		}

		dispatch++
	}

	if dispatch < numHirc && eof {
		return nil, fmt.Errorf(
			"Hierarchy decoding process encounter EOF after dispatching %d decoding routine. The total # of hierarchy is %d",
			dispatch, numHirc,
		)
	}

	if jobs != nil {
		for finished.Load() < numHirc {
			select {
			case <- ctx.Done():
				return nil, fmt.Errorf(
					"Failed to finish decoding HIRC due to context cancel: %w", 
					ctx.Err(),
					)
			default:
			}
		}
	}

	return h, nil
}
