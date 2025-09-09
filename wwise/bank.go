package wwise

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sort"
	"sync"
)

type Bank struct {
	mu sync.Mutex

	ChunkPosition map[string]u8
	EncodedChunk  map[string][]byte

	BKHD     *BKHD
	DIDXDATA *DIDXDATA
	HIRC     *HIRC
}

type ChunkHeader struct {
	ChunkName [4]byte
	ChunkSize    u32
}

func NewBank() *Bank {
	return &Bank{
		ChunkPosition: make(map[string]u8, 11),
		EncodedChunk: make(map[string][]byte, 7),
	}
}

type EncodeBankOpt struct {
	option u8
}

const MaskMETA u8 = 0b1000_0000

func IncludeEncodedMETA(o *EncodeBankOpt) {
	o.option |= MaskMETA
}

func IsIncludeEncodedMETA(o *EncodeBankOpt) bool {
	return o.option & MaskMETA > 0
}

func ExcludeEncodedMETA(o *EncodeBankOpt) {
	o.option = o.option | (^MaskMETA)
}

// The encoded chunk will follow convention / order imposed by Wwise authoring 
// tool.
// BKHD -> DIDX -> DATA -> HIRC
func EncodeBank(
	ctx context.Context, 
	w    io.Writer,
	o    order, 
	b   *Bank, 
	opt *EncodeBankOpt,
) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if opt == nil {
		opt = &EncodeBankOpt{}
		IncludeEncodedMETA(opt)
	}

	err = EncodeBKHD(b.BKHD, w, o)
	if err != nil {
		return err
	}

	if b.DIDXDATA.AudioData != nil {
		ComputeDIDXOffset(b.DIDXDATA)

		err = VerifyDIDXDATA(b.DIDXDATA)
		if err != nil {
			return err
		}

		err = EncodeDIDX(b.DIDXDATA, w, o)
		if err != nil {
			return err
		}
		err = EncodeDATANotAlign(b.DIDXDATA, w)
		if err != nil {
			return err
		}
	} else {
		err = EncodeDIDX(b.DIDXDATA, w, o)
		if err != nil {
			return err
		}

		chunk, in := b.EncodedChunk["DATA"]
		if !in {
			slog.Warn("DATA chunk is missing")
		}

		_, err = w.Write(chunk)
		if err != nil {
			return err
		}
	}

	// Temporary
	chunk, in := b.EncodedChunk["HIRC"]
	if !in {
		slog.Warn("HIRC chunk is missing")
	}
	_, err = w.Write(chunk)
	if err != nil {
		return err
	}
	
	// Write the rest of encoded chunks in the order appeared in the decoding 
	// phase.
	type ChunkPosition struct {
		ChunkName string
		Position  u8
	}

	chunkPositions := make([]ChunkPosition, 0, len(b.ChunkPosition))
	for chunkName, pos := range b.ChunkPosition {
		if chunkName == "META" && IsIncludeEncodedMETA(opt) {
			continue
		}
		switch chunkName {
		case ChunkNameBKHD:
		case ChunkNameDIDX:
		case ChunkNameDATA:
		case ChunkNameHIRC:
		default:
			i, found := sort.Find(len(chunkPositions), func(i int) int {
				if pos < chunkPositions[i].Position {
					return -1
				}
				if pos == chunkPositions[i].Position {
					return 0
				}
				return 1
			})

			if found {
				return fmt.Errorf(
					"Chunk %s and chunk %s occupy the same chunk position %d",
					chunkName, chunkPositions[i].ChunkName, chunkPositions[i].Position,
					)
			}

			chunkPositions = slices.Insert(
				chunkPositions, i, ChunkPosition{ chunkName, pos },
			)
		}
	}

	for _, chunkPos := range chunkPositions {
		chunkName := chunkPos.ChunkName
		chunk, in := b.EncodedChunk[chunkName]
		if !in {
			slog.Warn(fmt.Sprintf("Chunk %s is missing", chunkPos.ChunkName))
		}
		if _, err = w.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}

// No side effect
// Thread safe
func HasChunk(b *Bank, name string) (in bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, in = b.ChunkPosition[name]
	return in 
}

// No side effect
// Thread safe
func PopEncodedChunk(bnk *Bank, name string) (in bool, chunk []byte) {
	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	chunk, in = bnk.EncodedChunk[name]
	if !in {
		return in, nil
	}

	delete(bnk.EncodedChunk, name)

	return in, chunk 
}

// Has side effect
// Thread safe
func RegBKHD(bnk *Bank, bkhd *BKHD) error {
	if bkhd == nil {
		panic("bkhd is nil")
	}

	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	if _, in := bnk.ChunkPosition["BKHD"]; in {
		return fmt.Errorf("Duplicated BKHD chunk")
	}
	bnk.ChunkPosition["BKHD"] = 0

	bnk.BKHD = bkhd

	return nil
}

// Has side effect
// Thread safe
func RegDIDXDATA(bnk *Bank, didxdata *DIDXDATA, pos u8) error {
	if didxdata == nil {
		panic("didxdata is nil")
	}

	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	if _, in := bnk.ChunkPosition["DIDX"]; in {
		return fmt.Errorf("Duplicated DIDX chunk")
	}
	bnk.ChunkPosition["DIDX"] = pos

	bnk.DIDXDATA = didxdata

	return nil
}

// Has side effect
// Thread safe
func NewEncodedChunk(bnk *Bank, chunkName string, pos u8, encoded []byte) error {
	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	if encoded == nil {
		panic("Encoded slice is nil")
	}

	if _, in := bnk.ChunkPosition[chunkName]; in {
		return fmt.Errorf("Duplicate %s chunk", chunkName)
	}

	if _, in := bnk.EncodedChunk[chunkName]; in {
		panic(fmt.Sprintf("Duplicate encoded chunk %s", chunkName))
	}

	bnk.ChunkPosition[chunkName] = pos
	bnk.EncodedChunk[chunkName] = encoded

	return nil
}

// Has side effect
// Thread safe
func NewAudioSources(
	didxdata *DIDXDATA, 
	newSourceIds []u32, 
	audioData [][]byte,
) (ok []u32, fail []u32, err error) {
	return ok, fail, err
}
