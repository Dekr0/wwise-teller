// TODO
// - Rework of encoding (detail see TODO.md)
package wwise

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
)

var BankVersion = -1

var NoHIRC = errors.New("This sound bank does not have HIRC chunk.")
var NoDIDX = errors.New("This sound bank does not have DIDX chunk.")
var NoDATA = errors.New("This sound bank does not have DATA chunk.")

const SizeOfChunkHeader = 4 + 4

type Chunk interface {
	Encode(ctx context.Context) ([]byte, error)
	Tag() []byte
	Idx() uint8 // for maintaining the order of each chunk section
}

type Bank struct {
	Chunks  []Chunk

	// Experiment
	Sources        []Source
	SourcesMutex     sync.Mutex
	SourcesMap   map[uint32]uint32
}

func NewBank() Bank {
	return Bank{
		make([]Chunk, 0, 4),
		[]Source{},
		sync.Mutex{},
		make(map[uint32]uint32),
	}
}

func (b *Bank) AddChunk(c Chunk) error {
	if slices.ContainsFunc(b.Chunks, func(tc Chunk) bool {
		if bytes.Compare(tc.Tag(), c.Tag()) == 0 {
			return true
		}
		return false
	}) {
		return fmt.Errorf("Chunk %s already exists", c.Tag())
	}
	b.Chunks = append(b.Chunks, c)
	return nil
}

func (b *Bank) BKHD() *BKHD {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'B', 'K', 'H', 'D'}) == 0 {
			return chunk.(*BKHD)
		}
	}
	return nil
}

func (b *Bank) DIDX() *DIDX {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'D', 'I', 'D', 'X'}) == 0 {
			return chunk.(*DIDX)
		}
	}
	return nil
}

func (b *Bank) DATA() *DATA {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'D', 'A', 'T', 'A'}) == 0 {
			return chunk.(*DATA)
		}
	}
	return nil
}

func (b *Bank) HIRC() *HIRC {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'H', 'I', 'R', 'C'}) == 0 {
			return chunk.(*HIRC)
		}
	}
	return nil
}

func (b *Bank) META() *META {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'M', 'E', 'T', 'A'}) == 0 {
			return chunk.(*META)
		}
	}
	return nil
}

type EncodedChunk struct {
	i uint8
	b []byte
	e error
}

func CreateEncodeClosure(
	ctx context.Context, c chan *EncodedChunk, cu Chunk,
) func() {
	return func() {
		slog.Debug(fmt.Sprintf("Start encoding %s section", cu.Tag()))
		data, err := cu.Encode(ctx)
		c <- &EncodedChunk{cu.Idx(), data, err}
	}
}

func (bnk *Bank) Encode(ctx context.Context) ([]byte, error) {
	c := make(chan *EncodedChunk, len(bnk.Chunks))

	// No initialization since I want it to crash and catch encoding bugs
	chunks := make([][]byte, len(bnk.Chunks))

	i := 0
	for _, cu := range bnk.Chunks {
		if bytes.Compare(cu.Tag(), []byte{'M', 'E', 'T', 'A'}) == 0 {
			continue
		}
		go CreateEncodeClosure(ctx, c, cu)()
		i += 1
	}

	for i > 0 {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case res := <- c:
			if res.e != nil {
				return nil, res.e
			}
			chunks[res.i] = res.b
			slog.Info(
				fmt.Sprintf("Encoded %s section", res.b[0:4]),
				"size", len(res.b[8:]),
			)
			i -= 1
		}
	}

	return bytes.Join(chunks, []byte{}), nil
}

func (b *Bank) AppendAudio(audioData []byte, sid uint32) error {
	didx := b.DIDX()
	if didx == nil {
		return NoDIDX
	}
	data := b.DATA()
	if data == nil {
		return NoDATA
	}
	err := didx.Append(sid, uint32(len(audioData)))
	if err != nil {
		return err
	}
	data.B = append(data.B, audioData...)
	return nil
}
