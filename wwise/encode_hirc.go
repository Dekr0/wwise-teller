package wwise

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	uio "github.com/Dekr0/unwise/io"
)

type EncodeHircOpt struct {}

type HircEncoderCtx struct {
	Encoder *uio.EncoderCtx
	Version  u32
}

func (e *HircEncoderCtx) Count() u32 {
	return e.Encoder.Count
}

func (e *HircEncoderCtx) Primitive(data any) (err error) {
	return uio.Encode(e.Encoder, data)
}

func (e *HircEncoderCtx) Bytes(data []byte) (err error) {
	return uio.EncodeBytes(e.Encoder, data)
}

func (e *HircEncoderCtx) Struct(data any, size u32) (err error) {
	return uio.EncodeStruct(e.Encoder, data, size)
}

func (e *HircEncoderCtx) Expect(prev u32, expect u32) (err error) {
	return uio.AssertEncodeLimit(e.Encoder, prev, expect)
}

// Has no side effect
func SizeOfHIRC(h *HIRC, version u32) (size u32) {
	size = SizeOfHierarchyNumCounter

	hierarchy := &h.Hierarchy

	internalIds := hierarchy.InternalIds
	nodes := hierarchy.Nodes
	encodedNodes := hierarchy.EncodedNodes

	size += SizeOfHierarchyHeader * u32(len(internalIds))
	for _, id := range internalIds {
		hierarchy, in := nodes[id]
		if !in {
			panic(fmt.Errorf("Internal id %d does not have hierarchy info", id))
		}
		hid, t := hierarchy.Id, hierarchy.Type
		switch t {
		case HircTypeState:
			if err := h.StateComponent.AssertStateById(id); err != nil {
				panic(fmt.Errorf("(State %d): %w", hid, err))
			}
			// TODO: refactor
			size += h.StateComponent.SizeOfStateById(id)
		case HircTypeSound:
			s := h.SoundH(id, version)
			if err := AssertSound(s, version); err != nil {
				panic(fmt.Errorf("(Sound %d): %w", hid, err))
			}
			size += SizeOfSound(s, version)
		case HircTypeEvent:
			if err := h.EventComponet.AssertEventById(id); err != nil {
				panic(fmt.Errorf("(Event %d): %w", id, err))
			}
			// TODO: refactor
			size += h.EventComponet.SizeOfEventById(id)
		default:
			encodedData, in := encodedNodes[id]
			if !in {
				panic(fmt.Errorf("Internal id %d does not have encoded data for %s %d", id, GetHircTypeName(t), hid))
			}
			size += u32(len(encodedData))
		}
	}

	return size
}

// Has no side effect
func EncodeHirc(
	ctx      context.Context,
	e       *HircEncoderCtx,
	h       *HIRC, 
	opt     *EncodeHircOpt,
) (err error) {
	size := SizeOfHIRC(h, e.Version)

	chunkHeader := ChunkHeader{ [4]byte([]byte(ChunkNameHIRC)), size }
	if err := e.Struct(chunkHeader, SizeOfChunkHeader); err != nil {
		return fmt.Errorf("Failed to encode HIRC chunk header: %w", err)
	}

	hierarchy := &h.Hierarchy

	internalIds := hierarchy.InternalIds
	nodes := hierarchy.Nodes

	if err := e.Primitive(u32(len(internalIds))); err != nil {
		return fmt.Errorf("Failed to encode number of Hierarchy: %w", err)
	}

	pool := sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}

	for _, internalId := range internalIds {
		hierarchy, in := nodes[internalId]
		if !in {
			panic(fmt.Errorf("Internal id %d does not have hierarchy", internalId))
		}

		hierarchyId, t := hierarchy.Id, hierarchy.Type
		switch t {
		case HircTypeState:
			bufWriter := pool.Get().(*bytes.Buffer)

			be := HircEncoderCtx{
				Encoder: &uio.EncoderCtx{
					Writer: bufWriter,
					Order: e.Encoder.Order,
					Count: 0,
				},
				Version: e.Version,
			}

			state := h.GatherStateData(internalId)
			size := SizeOfState(state.StateProps)
			EncodeState(&be, state, size)

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf("Failed to encode State %d: %w", hierarchyId, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		case HircTypeSound:
			bufWriter := pool.Get().(*bytes.Buffer)

			be := HircEncoderCtx{
				Encoder: &uio.EncoderCtx{
					Writer: bufWriter,
					Order: e.Encoder.Order,
					Count: 0,
				},
				Version: e.Version,
			}

			h.EncodeSound(&be, internalId)

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf("Failed to encode Sound %d: %w", hierarchyId, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		case HircTypeEvent:
			bufWriter := pool.Get().(*bytes.Buffer)

			be := HircEncoderCtx{
				Encoder: &uio.EncoderCtx{
					Writer: bufWriter,
					Order: e.Encoder.Order,
					Count: 0,
				},
				Version: e.Version,
			}

			event := h.GatherEventData(internalId)
			size := SizeOfEvent(event.EventData)
			EncodeEvent(&be, event, size)

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf("Failed to encode Event %d: %w", hierarchyId, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		default:
			ts := GetHircTypeName(t)
			if err = h.EncodeEncodedHierarchy(e, t, internalId); err != nil {
				return fmt.Errorf("Failed to encode %s %d: %w", ts, hierarchyId, err)
			}
		}
	}

	return nil
}
