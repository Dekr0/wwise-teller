package wwise

import (
	"bytes"
	"context"
	bin "encoding/binary"
	"fmt"
	"io"
	"sync"
)

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
	w        io.Writer,
	o        order,
	version  u32,
	h       *HIRC, 
	opt     *EncodeHircOpt,
) (err error) {
	size := SizeOfHIRC(h, version)

	chunkHeader := ChunkHeader{ [4]byte{ 'H', 'I', 'R', 'C' }, size }
	if err := bin.Write(w, o, chunkHeader); err != nil {
		return fmt.Errorf("Failed to encode HIRC chunk header: %w", err)
	}

	if err := bin.Write(w, o, u32(len(h.InternalIds))); err != nil {
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

			EncodeState(w, o, &h.StateComponent, version, internalId, hierarchyId)

			encoded := bufWriter.Bytes()

			n, err := w.Write(encoded)
			if err != nil {
				return fmt.Errorf("Failed to encode State %d: %w", hierarchyId, err)
			}
			if n != len(encoded) {
				return fmt.Errorf(
					"Failed to encode State %d: # (%d) of bytes written does not equal to # (%d) of bytes from buffer writer",
					hierarchyId, n, len(encoded),
				)
			}
			bufWriter.Reset()
			pool.Put(bufWriter)
		case HircTypeEvent:
			bufWriter := pool.Get().(*bytes.Buffer)

			EncodeEvent(w, o, &h.EventComponet, version, internalId, hierarchyId)

			encoded := bufWriter.Bytes()

			n, err := w.Write(encoded)
			if err != nil {
				return fmt.Errorf("Failed to encode Event %d: %w", hierarchyId, err)
			}
			if n != len(encoded) {
				return fmt.Errorf(
					"Failed to encode Event %d: # (%d) of bytes written does not equal to # (%d) of bytes from buffer writer",
					hierarchyId, n, len(encoded),
				)
			}
		default:
			hierarchyName := GetHircTypeName(t)
			encodedChunk, in := encodedHierarchy[internalId]
			if !in {
				panic(fmt.Sprintf("%s %d does not have encoded chunk", 
					hierarchyName, hierarchyId,
				))
			}

			header := HierarchyHeader{ t, u32(len(encodedChunk)) }
			if err := bin.Write(w, o, header); err != nil {
				return fmt.Errorf("Failed to encode %s %d header: %w", 
					hierarchyName, hierarchyId, err,
				)
			}

			n, err := w.Write(encodedChunk)
			if err != nil {
				return fmt.Errorf("Failed to write encoded chunk of %s %d: %w",
					hierarchyName, hierarchyId, err,
				)
			}

			if n != len(encodedChunk) {
				return fmt.Errorf(
					"Failed to encode %s %d: # (%d) of bytes written does not equal to # (%d) of bytes from buffer writer",
					hierarchyName, hierarchyId, n, len(encodedChunk),
				)
			}
		}
	}

	return nil
}
