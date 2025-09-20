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

type HierarchyDecoder func(io.Reader, order, u32, u32, any)

type HircDecodeOption struct {
	Exclude    []u8
}

func AllocDecodeHIRC(
	ctx       context.Context, 
	opt      *HircDecodeOption,
	spec     *wwise.HIRCSpaceSpec,
	inReader  io.Reader, 
	o         order,
	size      u32, 
	version   u32, 
) (h *wwise.HIRC, err error) {
	if opt == nil {
		return nil, fmt.Errorf("Must provide HIRC decoder option")
	}
	if spec == nil {
		return nil, fmt.Errorf("Must provide hierarchy space allocation specification")
	}

	r := io.LimitReader(inReader, int64(size))

	numHirc, err := uio.U32(r, o)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode # of hierarchies: %w", err)
	}

	h = wwise.AllocHIRC(spec)

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
		var res any
		switch t {
		case wwise.HircTypeState:
			res = &wwise.State{}
			decoder = AllocDecodeState
		case wwise.HircTypeSound:
			res = &wwise.Sound{}
			decoder = AllocDecodeSound
		case wwise.HircTypeEvent:
			res = &wwise.Event{}
			decoder = AllocDecodeEvent
		case wwise.HircTypeActorMixer:
			res = &wwise.ActorMixer{}
			decoder = AllocDecodeActorMixer
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
		decoder(reader, o, version, size, res)
		switch t := res.(type) {
		case *wwise.State:
			h.AddState(*t)
		case *wwise.Sound:
			h.AddSound(*t, version)
		case *wwise.Event:
			h.AddEvent(*t)
		case *wwise.ActorMixer:
		 	h.AddActorMixer(*t, version)
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

func PrefetchHIRCMetadata(
	ctx       context.Context, 
	r         io.ReadSeeker, 
	out      *wwise.HierarchyStat,
	o         order,
	size      u32, 
) error {
	if out == nil {
		return fmt.Errorf("Must provide hierarchy statistic")
	}

	numHirc, err := uio.U32(r, o)
	if err != nil {
		return fmt.Errorf("Failed to decode # of hierarchies: %w", err)
	}

	eof := false
	dispatch := uint32(0)

	for dispatch < numHirc && !eof {
		select {
		case <- ctx.Done():
			return fmt.Errorf("Hierarchy space estimation process cancel: %w", ctx.Err())
		default:
		}
		t, err := uio.U8(r, o)
		if err != nil {
			if err == io.EOF {
				eof = true
				break
			}
			return fmt.Errorf("Failed to decode hierarchy type: %w", err)
		}
		size, err := uio.U32(r, o)
		if err != nil {
			if err == io.EOF {
				eof = true
				break
			}
			return fmt.Errorf("Failed to decode hierarchy data size: %w", err)
		}
		switch t {
			case wwise.HircTypeState:
				out.State++
			case wwise.HircTypeSound:
				out.Sound++
			case wwise.HircTypeAction:
				out.Action++
			case wwise.HircTypeEvent:
				out.Event++
			case wwise.HircTypeRanSeqCntr:
				out.RanSeqCntr++
			case wwise.HircTypeSwitchCntr:
				out.SwitchCntr++
			case wwise.HircTypeActorMixer:
				out.ActorMixer++
			case wwise.HircTypeBus:
				out.Bus++
			case wwise.HircTypeLayerCntr:
				out.LayerCntr++
			case wwise.HircTypeMusicSegment:
				out.MusicSegment++
			case wwise.HircTypeMusicTrack:
				out.MusicTrack++
			case wwise.HircTypeMusicSwitchCntr:
				out.MusicSwitchCntr++
			case wwise.HircTypeMusicRanSeqCntr:
				out.MusicRanSeqCntr++
			case wwise.HircTypeAttenuation:
				out.Attenuation++
			case wwise.HircTypeDialogueEvent:
				out.DialogueEvent++
			case wwise.HircTypeFxShareSet:
				out.FxShareSet++
			case wwise.HircTypeFxCustom:
				out.FxCustom++
			case wwise.HircTypeAuxBus:
				out.AuxBus++
			case wwise.HircTypeLFOModulator:
				out.LFOModulator++
			case wwise.HircTypeEnvelopeModulator:
				out.EnvelopeModulator++
			case wwise.HircTypeAudioDevice:
				out.AudioDevice++
			case wwise.HircTypeTimeModulator:
				out.TimeModulator++
		}
		if _, err := r.Seek(int64(size), io.SeekCurrent); err != nil {
			return fmt.Errorf("Failed to skip %d bytes ahead to the next hierarchy: %w", size, err)
		}
	}
	return nil
}
