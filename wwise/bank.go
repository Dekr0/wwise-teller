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

func (b *Bank) DATAAppendOnly() *DATAAppendOnly {
	for _, chunk := range b.Chunks {
		if bytes.Compare(chunk.Tag(), []byte{'D', 'A', 'T', 'A'}) == 0 {
			switch c := chunk.(type) {
			case *DATAAppendOnly:
				return c
			}
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

func (bnk *Bank) Encode(ctx context.Context, diffTest bool) ([]byte, error) {
	if diffTest {
		bnk.ComputeDIDXOffset()
		if bnk.DIDX() != nil && bnk.DATA() != nil {
			if err := bnk.CheckDIDXDATA(); err != nil {
				return nil, err
			}
		}
	}

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
				fmt.Sprintf("Encoded %s section.", res.b[0:4]),
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
	if _, in := data.AudiosMap[sid]; in {
		return fmt.Errorf("ID %d already has an associate audio data.", sid)
	}
	err := didx.Append(sid, uint32(len(audioData)))
	if err != nil {
		return err
	}
	data.Audios = append(data.Audios, audioData)
	data.AudiosMap[sid] = audioData
	return nil
}

func (b *Bank) ReplaceAudio(audioData []byte, sid uint32) error {
	didx := b.DIDX()
	if didx == nil {
		return NoDIDX
	}
	data := b.DATA()
	if data == nil {
		return NoDATA
	}
	_, in := data.AudiosMap[sid]
	if !in {
		return fmt.Errorf("No audio data has ID %d in index", sid)
	}
	_, in = didx.MediaIndexsMap[sid]
	if !in {
		return fmt.Errorf("No media index has ID %d in index", sid)
	}
	didx.MediaIndexsMap[sid].Size = uint32(len(audioData))
	audioIdx := slices.IndexFunc(didx.MediaIndexs, func(m MediaIndex) bool {
		return m.Sid == sid
	})
	if audioIdx == -1 {
		return fmt.Errorf("Failed to index media index using source ID %d", sid)
	}
	data.Audios[audioIdx] = audioData
	data.AudiosMap[sid] = audioData

	return nil
}

func (b *Bank) ComputeDIDXOffset() {
	didx := b.DIDX()
	if didx == nil {
		return
	}
	data := b.DATA()
	if data == nil {
		return
	}
	offset := uint64(0)
	for i, entry := range didx.MediaIndexs {
		didx.MediaIndexs[i].Offset = uint32(offset)
		offset += uint64(entry.Size)
	}
}

func (b *Bank) CheckDIDXDATA() error {
	didx := b.DIDX()
	if didx == nil {
		return NoDIDX
	}
	data := b.DATA()
	if data == nil {
		return NoDATA
	}
	offset := uint64(0)
	if len(didx.MediaIndexs) != len(data.Audios) {
		return fmt.Errorf(
			"# of Media Index (%d) doesnt' equal # of audios data (%d)", len(didx.MediaIndexs), len(data.Audios),
		)
	}
	for i, entry := range didx.MediaIndexs {
		if uint32(len(data.Audios[i])) != entry.Size {
			return fmt.Errorf("Audio source size at index %d does not equal to size of Media Index (%d) at index %d.", i, entry.Size, i)
		}
		if _, in := didx.MediaIndexsMap[entry.Sid]; !in {
			return fmt.Errorf("Media index %d is not in index.", entry.Sid)
		}
		if _, in := data.AudiosMap[entry.Sid]; !in {
			return fmt.Errorf("Media index %d cannot find audio data in audio data index", entry.Sid)
		}
		if offset != uint64(entry.Offset) {
			return fmt.Errorf("Expecting media index (%d) at index %d has offset of %d but received %d", entry.Sid, i, offset, entry.Offset)
		}
		offset += uint64(entry.Size)
	}
	return nil
}
