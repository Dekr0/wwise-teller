package wwise

import (
	"context"
	bin "encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sort"
)

type Bank struct {
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
	ctx  context.Context, 
	w    io.Writer,
	o    order, 
	b   *Bank, 
	opt *EncodeBankOpt,
) (err error) {
	if opt == nil {
		return fmt.Errorf("Must provided a bank encoding option")
	}

	err = EncodeBKHD(b.BKHD, w, o)
	if err != nil {
		return fmt.Errorf("Failed to encode BKHD chunk: %w", err)
	}
	slog.Info("Encoded BKHD")

	if b.DIDXDATA.AudioData != nil {
		ComputeDIDXOffset(b.DIDXDATA)
		VerifyDIDXDATA(b.DIDXDATA)

		err = EncodeDIDX(b.DIDXDATA, w, o)
		if err != nil {
			return fmt.Errorf("Failed to encode DIDX chunk: %w", err)
		}
		slog.Info("Encoded DIDX")

		err = EncodeDATANotAlign(b.DIDXDATA, w)
		if err != nil {
			return fmt.Errorf("Failed to encode DATA chunk without alignment: %w", err)
		}
		slog.Info("Encoded DATA")
	} else {
		err = EncodeDIDX(b.DIDXDATA, w, o)
		if err != nil {
			return err
		}
		slog.Info("Encoded DIDX")

		chunk, in := b.EncodedChunk["DATA"]
		if !in {
			slog.Warn("Encoded DATA chunk is missing")
		} else {
			chunkHeader := ChunkHeader{
				[4]byte{'D', 'A', 'T', 'A'}, 
				u32(len(chunk)),
			}
			if err = bin.Write(w, o, chunkHeader); err != nil {
				return fmt.Errorf("Failed to write DATA chunk header: %w", err)
			}

			_, err = w.Write(chunk)
			if err != nil {
				return fmt.Errorf("Failed to write encoded DATA chunk: %w", err)
			}
			slog.Info("Encoded DATA")
		}
	}

	{ // Temporary
		chunk, in := b.EncodedChunk["HIRC"]
		if !in {
			slog.Warn("Encoded HIRC chunk is missing")
		} else {
			chunkHeader := ChunkHeader{
				[4]byte{'H', 'I', 'R', 'C'}, 
				u32(len(chunk)),
			}
			if err = bin.Write(w, o, chunkHeader); err != nil {
				return fmt.Errorf("Failed to write HIRC chunk header: %w", err)
			}

			_, err = w.Write(chunk)
			if err != nil {
				return fmt.Errorf("Failed to write encoded HIRC chunk: %w", err)
			}
			slog.Info("Encoded HIRC")
		}
	}
	
	// Write the rest of encoded chunks in the order appeared in the decoding 
	// phase.
	type ChunkPosition struct {
		ChunkName string
		Position  u8
	}

	chunkPositions := make([]ChunkPosition, 0, len(b.ChunkPosition))
	for chunkName, pos := range b.ChunkPosition {
		if chunkName == "META" && !IsIncludeEncodedMETA(opt) {
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
				panic(fmt.Sprintf(
					"Chunk %s and chunk %s occupy the same chunk position %d",
					chunkName, chunkPositions[i].ChunkName, chunkPositions[i].Position,
				))
			}

			chunkPositions = slices.Insert(
				chunkPositions, i, ChunkPosition{ chunkName, pos },
			)
		}
	}

	for _, chunkPos := range chunkPositions {
		chunkName := chunkPos.ChunkName
		chunk, in := b.EncodedChunk[chunkName]

		chunkHeader := ChunkHeader{
			[4]byte{chunkName[0], chunkName[1], chunkName[2], chunkName[3]}, 
			u32(len(chunk)),
		}
		if err = bin.Write(w, o, chunkHeader); err != nil {
			return fmt.Errorf("Faile to write %s chunk header: %w", chunkName, err)
		}

		if !in {
			slog.Warn(fmt.Sprintf("Encoded chunk of %s is missing", chunkName))
			continue
		}

		if _, err = w.Write(chunk); err != nil {
			return fmt.Errorf("Failed to write encoded chunk of %s: %w", chunkName, err)
		}

		slog.Info(fmt.Sprintf("Encoded %s", chunkName))
	}

	return nil
}

// No side effect
// Thread safe
func HasChunk(b *Bank, name string) (in bool) {
	_, in = b.ChunkPosition[name]
	return in 
}

// No side effect
// Thread safe
func PopEncodedChunk(bnk *Bank, name string) (in bool, chunk []byte) {
	chunk, in = bnk.EncodedChunk[name]
	if !in {
		return in, nil
	}

	delete(bnk.EncodedChunk, name)

	return in, chunk 
}

// Has side effect
// Thread safe
func RegBKHD(bnk *Bank, bkhd *BKHD) {
	if bkhd == nil {
		panic("bkhd is nil")
	}

	if _, in := bnk.ChunkPosition["BKHD"]; in {
		panic(fmt.Sprintf("Duplicated BKHD chunk"))
	}
	bnk.ChunkPosition["BKHD"] = 0

	bnk.BKHD = bkhd
}

// Has side effect
// Thread safe
func RegDIDXDATA(bnk *Bank, didxdata *DIDXDATA, pos u8) {
	if didxdata == nil {
		panic("didxdata is nil")
	}

	if _, in := bnk.ChunkPosition["DIDX"]; in {
		panic(fmt.Sprintf("Duplicated DIDX chunk"))
	}
	bnk.ChunkPosition["DIDX"] = pos

	bnk.DIDXDATA = didxdata
}

// Has side effect
// Thread safe
func NewEncodedChunk(bnk *Bank, chunkName string, pos u8, encoded []byte) {
	if encoded == nil {
		panic("Encoded slice is nil")
	}

	if _, in := bnk.ChunkPosition[chunkName]; in {
		panic(fmt.Sprintf("Duplicate %s chunk", chunkName))
	}

	if _, in := bnk.EncodedChunk[chunkName]; in {
		panic(fmt.Sprintf("Duplicate encoded chunk %s", chunkName))
	}

	bnk.ChunkPosition[chunkName] = pos
	bnk.EncodedChunk[chunkName] = encoded
}
