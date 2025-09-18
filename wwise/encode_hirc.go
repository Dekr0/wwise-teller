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
		tName := GetHircTypeName(t)
		switch t {
		case HircTypeState:
			s := h.State(id)
			if err := AssertState(s); err != nil {
				panic(fmt.Errorf("(%s %d): %w", tName, hid, err))
			}
			size += SizeOfState(s)
		case HircTypeSound:
			s := h.Sound(id, version)
			if err := AssertSound(s, version); err != nil {
				panic(fmt.Errorf("(%s %d): %w", tName, hid, err))
			}
			size += SizeOfSound(s, version)
		case HircTypeEvent:
			e := h.Event(id)
			if err := AssertEvent(e); err != nil {
				panic(fmt.Errorf("(%s %d): %w", tName, id, err))
			}
			size += SizeOfEvent(e)
		case HircTypeActorMixer:
			a := h.ActorMixer(id, version)
			if err := AssertActorMixer(a, version); err != nil {
				panic(fmt.Errorf("(%s %d): %w", tName, hid, err))
			}
			size += SizeOfActorMixer(a, version)
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

		const errMsg = "Failed to encode %s %d: %w"
		hid, t := hierarchy.Id, hierarchy.Type
		name := GetHircTypeName(t)
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

			state := h.State(internalId)
			if err = EncodeState(&be, state); err != nil {
				panic(fmt.Errorf(errMsg, name, hid, err))
			}

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf(errMsg, name, hid, err)
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

			sound := h.Sound(internalId, e.Version)
			if err := EncodeSound(&be, sound); err != nil {
				panic(fmt.Errorf(errMsg, name, hid, err))
			}

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf(errMsg, name, hid, err)
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

			event := h.Event(internalId)
			if err := EncodeEvent(&be, event); err != nil {
				panic(fmt.Errorf(errMsg, name, hid, err))
			}

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf("Failed to encode %s %d: %w", name, hid, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		case HircTypeActorMixer:
			bufWriter := pool.Get().(*bytes.Buffer)

			be := HircEncoderCtx{
				Encoder: &uio.EncoderCtx{
					Writer: bufWriter,
					Order: e.Encoder.Order,
					Count: 0,
				},
				Version: e.Version,
			}

			a := h.ActorMixer(internalId, e.Version)
			if err := EncodeActorMixer(&be, a); err != nil {
				panic(fmt.Errorf(errMsg, name, hid, err))
			}

			encoded := bufWriter.Bytes()

			if err = e.Bytes(encoded); err != nil {
				return fmt.Errorf(errMsg, name, hid, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		default:
			if err = h.EncodeEncodedHierarchy(e, t, internalId); err != nil {
				return fmt.Errorf(errMsg, name, hid, err)
			}
		}
	}

	return nil
}
