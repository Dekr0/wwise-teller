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

type ChunkHeader struct {
	ChunkName [4]byte
	ChunkSize    u32
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
	ctx      context.Context, 
	w        io.Writer,
	o        order, 
	b       *Bank, 
	bankOpt *EncodeBankOpt,
	hircOpt *EncodeHircOpt,
) (err error) {
	if bankOpt == nil {
		return fmt.Errorf("Must provided a bank encoding option")
	}

	if b.BKHD == nil {
		panic(fmt.Errorf("A sound bank without BKHD passed the decoding phase"))
	}

	err = EncodeBKHD(b.BKHD, w, o)
	if err != nil {
		return fmt.Errorf("Failed to encode BKHD chunk: %w", err)
	}
	slog.Info("Encoded BKHD")

	if b.DIDXDATA != nil {
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
	} else {
		slog.Warn(fmt.Sprintf("Sound bank %d does not have DIDX chunk (or DATA chunk as well)", b.BKHD.Id))
	}

	if b.HIRC != nil {
		err := EncodeHirc(ctx, w, o, b.BKHD.Version, b.HIRC, hircOpt)
		if err != nil {
			return fmt.Errorf("Failed to encode HIRC chunk: %w", err)
		}
	} else {
		slog.Warn("Sound bank ")
	}
	
	// Write the rest of encoded chunks in the order appeared in the decoding 
	// phase.
	type ChunkPosition struct {
		ChunkName string
		Position  u8
	}

	chunkPositions := make([]ChunkPosition, 0, len(b.ChunkPosition))
	for chunkName, pos := range b.ChunkPosition {
		if chunkName == "META" && !IsIncludeEncodedMETA(bankOpt) {
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
