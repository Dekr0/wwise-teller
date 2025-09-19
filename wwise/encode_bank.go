package wwise

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	uio "github.com/Dekr0/unwise/io"
)

type ChunkHeader struct {
	ChunkName [4]byte
	ChunkSize    u32
}

const SizeOfChunkHeader = 8

// The encoded chunk will follow convention / order imposed by Wwise authoring 
// tool.
// BKHD -> DIDX -> DATA -> INIT -> STMG -> HIRC
func EncodeBank(
	ctx      context.Context, 
	e       *uio.EncoderCtx,
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

	err = EncodeBKHD(b.BKHD, e)
	if err != nil {
		return fmt.Errorf("Failed to encode BKHD chunk: %w", err)
	}
	slog.Info("Encoded BKHD")

	if b.AudioStore != nil {
		if b.AudioStore.AudioData != nil {
			ComputeDIDXOffset(b.AudioStore)
			VerifyDIDXDATA(b.AudioStore)

			err = EncodeDIDX(b.AudioStore, e)
			if err != nil {
				return fmt.Errorf("Failed to encode DIDX chunk: %w", err)
			}
			slog.Info("Encoded DIDX")

			err = EncodeDATANotAlign(b.AudioStore, e)
			if err != nil {
				return fmt.Errorf("Failed to encode DATA chunk without alignment: %w", err)
			}
			slog.Info("Encoded DATA")
		} else {
			err = EncodeDIDX(b.AudioStore, e)
			if err != nil {
				return err
			}
			slog.Info("Encoded DIDX")

			if err = EncodeEncodedChunk(&b.Chunk, TagDATA, e); err != nil {
				return err
			}
		}
	} else {
		slog.Warn(fmt.Sprintf("Sound bank %d does not have DIDX chunk and / or DATA chunk", b.BKHD.Id))
	}

	if err := EncodeEncodedChunk(&b.Chunk, TagINIT, e); err != nil {
		return err
	}

	if err := EncodeEncodedChunk(&b.Chunk, TagSTMG, e); err != nil {
		return err
	}

	if b.HIRC != nil {
		err := EncodeHirc(ctx, &HircEncoderCtx{e, b.BKHD.Version}, b.HIRC, hircOpt)
		if err != nil {
			return fmt.Errorf("Failed to encode HIRC chunk: %w", err)
		}
	} else {
		if err = EncodeEncodedChunk(&b.Chunk, TagHIRC, e); err != nil {
			return fmt.Errorf("Faile to encode HIRC chunk: %w", err)
		}
	}

	return EncodeRemainEncodedChunk(&b.Chunk, e, bankOpt)
}

func EncodeEncodedChunk(c *ChunkComponent, tag Tag, e *uio.EncoderCtx) (err error) {
 	chunk, in := c.Encoded[tag]
 	if !in {
 		slog.Warn(fmt.Sprintf("Encoded %s chunk is missing", tag))
 	} else {
 		chunkHeader := ChunkHeader{
 			[4]byte([]byte(tag)), 
 			u32(len(chunk)),
 		}
 		if err = uio.EncodeStruct(e, chunkHeader, SizeOfChunkHeader); err != nil {
 			return fmt.Errorf("Failed to write %s chunk header: %w", tag, err)
 		}

 		if err = uio.EncodeBytes(e, chunk); err != nil {
 			return fmt.Errorf("Failed to write encoded %s chunk: %w", tag, err)
 		}
 		slog.Info(fmt.Sprintf("Encoded %s", tag))
 	}
 	return nil
}

func EncodeRemainEncodedChunk(
	c *ChunkComponent, 
	e *uio.EncoderCtx, 
	o *EncodeBankOpt,
) (err error) {
	// Write the rest of encoded chunks in the order appeared in the decoding 
	// phase.
	type ChunkPosition struct {
		ChunkName string
		Position  u8
	}

	chunkPositions := make([]ChunkPosition, 0, len(c.Position))
	for chunkName, v := range c.Position {
		if chunkName == "META" && !IsIncludeEncodedMETA(o) {
			continue
		}

		switch chunkName {
		case TagBKHD:
		case TagDIDX:
		case TagDATA:
		case TagINIT:
		case TagSTMG:
		case TagHIRC:
		default:
			i, found := sort.Find(len(chunkPositions), func(i int) int {
				if v < chunkPositions[i].Position {
					return -1
				}
				if v == chunkPositions[i].Position {
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
				chunkPositions, i, ChunkPosition{ chunkName, v },
				)
		}
	}

	for _, chunkPos := range chunkPositions {
		chunkName := chunkPos.ChunkName
		chunk, in := c.Encoded[chunkName]

		chunkHeader := ChunkHeader{
			[4]byte{chunkName[0], chunkName[1], chunkName[2], chunkName[3]}, 
			u32(len(chunk)),
		}
		if err = uio.EncodeStruct(e, chunkHeader, SizeOfChunkHeader); err != nil {
			return fmt.Errorf("Faile to write %s chunk header: %w", chunkName, err)
		}

		if !in {
			slog.Warn(fmt.Sprintf("Encoded chunk of %s is missing", chunkName))
			continue
		}

		if err = uio.EncodeBytes(e, chunk); err != nil {
			return fmt.Errorf("Failed to write encoded chunk of %s: %w", chunkName, err)
		}

		slog.Info(fmt.Sprintf("Encoded %s", chunkName))
	}

	return nil
}
