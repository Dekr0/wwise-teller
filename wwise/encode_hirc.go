package wwise

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	uio "github.com/Dekr0/unwise/io"
)

type HircEncoderCtx struct {
	Encoder *uio.EncoderCtx
	Version  u32
}

func HIRCEncode(e *HircEncoderCtx, data any) (err error) {
	return uio.Encode(e.Encoder, data)
}

func HIRCEncodeBytes(e *HircEncoderCtx, data []byte) (err error) {
	return uio.EncodeBytes(e.Encoder, data)
}

func HIRCEncodeStruct(e *HircEncoderCtx, data any, size u32) (err error) {
	return uio.EncodeStruct(e.Encoder, data, size)
}

func HIRCAssertEncodeLimit(e *HircEncoderCtx, prev u32, expect u32) (err error) {
	return uio.AssertEncodeLimit(e.Encoder, prev, expect)
}

func SizeOfHIRC(h *HIRC, version u32) (size u32) {
	size = SizeOfHierarchyNumCounter

	internalIds := h.InternalIds
	hierarchies := h.Hierarchies
	encodedHierarchy := h.EncodedHierarchy

	size += SizeOfHierarchyHeader * u32(len(internalIds))
	for _, id := range internalIds {
		hierarchy, in := hierarchies[id]
		if !in {
			panic(fmt.Errorf("Internal id %d does not have hierarchy info", id))
		}
		hid, t := hierarchy.Id, hierarchy.Type
		switch t {
		case HircTypeState:
			size += SizeOfState(&h.StateComponent, version, id, hid)
		case HircTypeEvent:
			size += SizeOfEvent(&h.EventComponet, version, id, hid)
		default:
			encodedData, in := encodedHierarchy[id]
			if !in {
				panic(fmt.Errorf("Internal id %d does not have encoded data for %s %d", id, GetHircTypeName(t), hid))
			}
			size += u32(len(encodedData))
		}
	}

	return size
}

func EncodeHirc(
	ctx      context.Context,
	e       *HircEncoderCtx,
	h       *HIRC, 
	opt     *EncodeHircOpt,
) (err error) {
	size := SizeOfHIRC(h, e.Version)

	chunkHeader := ChunkHeader{ [4]byte{ 'H', 'I', 'R', 'C' }, size }
	if err := HIRCEncodeStruct(e, chunkHeader, SizeOfChunkHeader); err != nil {
		return fmt.Errorf("Failed to encode HIRC chunk header: %w", err)
	}

	if err := HIRCEncode(e, u32(len(h.InternalIds))); err != nil {
		return fmt.Errorf("Failed to encode number of Hierarchy: %w", err)
	}

	pool := sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}

	internalIds := h.InternalIds
	hierarchies := h.Hierarchies
	encodedHierarchy := h.EncodedHierarchy
	for _, internalId := range internalIds {
		hierarchy, in := hierarchies[internalId]
		if !in {
			panic(fmt.Errorf("Internal id %d does not have hierarchy", internalId))
		}

		hierarchyId, t := hierarchy.Id, hierarchy.Type
		switch t {
		case HircTypeState:
			bufWriter := pool.Get().(*bytes.Buffer)

			EncodeState(e, &h.StateComponent, internalId, hierarchyId)

			encoded := bufWriter.Bytes()

			if err = HIRCEncodeBytes(e, encoded); err != nil {
				return fmt.Errorf("Failed to encode State %d: %w", hierarchyId, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		case HircTypeEvent:
			bufWriter := pool.Get().(*bytes.Buffer)

			EncodeEvent(e, &h.EventComponet, internalId, hierarchyId)

			encoded := bufWriter.Bytes()

			if err = HIRCEncodeBytes(e, encoded); err != nil {
				return fmt.Errorf("Failed to encode Event %d: %w", hierarchyId, err)
			}

			bufWriter.Reset()
			pool.Put(bufWriter)
		default:
			hierarchyName := GetHircTypeName(t)
			encodedChunk, in := encodedHierarchy[internalId]
			if !in {
				panic(fmt.Sprintf("%s %d does not have encoded chunk", 
					hierarchyName, hierarchyId,
				))
			}

			header := HierarchyHeader{ t, u32(len(encodedChunk)) }
			if err = HIRCEncodeStruct(e, header, SizeOfHierarchyHeader); err != nil {
				return fmt.Errorf("Failed to encode %s %d header: %w", 
					hierarchyName, hierarchyId, err,
				)
			}

			if err = HIRCEncodeBytes(e, encodedChunk); err != nil {
				return fmt.Errorf("Failed to write encoded chunk of %s %d: %w",
					hierarchyName, hierarchyId, err,
				)
			}
		}
	}

	return nil
}
