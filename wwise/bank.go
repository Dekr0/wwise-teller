package wwise

import (
	"fmt"
)

type Bank struct {
	BKHD       *BKHD
	AudioStore *AudioStore
	HIRC       *HIRC
	Chunk       ChunkComponent
}

type ChunkComponent struct {
	Position map[string]u8
	Encoded  map[string][]byte
}

func AllocChunkComponent() *ChunkComponent {
	return &ChunkComponent{
		Position: make(map[string]u8, 11),
		Encoded: make(map[string][]byte, 7),
	}
}

func AllocBank() *Bank {
	return &Bank{
		Chunk: *AllocChunkComponent(),
	}
}

// No side effect
func HasChunk(c *ChunkComponent, name string) (in bool) {
	_, in = c.Position[name]
	return in 
}

// Has side effect
func AddChunkPosition(c *ChunkComponent, name string, pos u8) {
	if _, in := c.Position[name]; in {
		panic(fmt.Sprintf("Postion value for chunk %s already exist: %d", name, pos))
	}
	c.Position[name] = pos
}

// Has side effect
func PopEncodedChunk(c *ChunkComponent, name string) (in bool, chunk []byte) {
	chunk, in = c.Encoded[name]
	if !in {
		return in, nil
	}

	delete(c.Encoded, name)

	return in, chunk 
}

// Has side effect
func RegBKHD(bnk *Bank, bkhd *BKHD) {
	if bkhd == nil {
		panic("bkhd is nil")
	}

	if HasChunk(&bnk.Chunk, ChunkNameBKHD) {
		panic(fmt.Sprintf("Duplicated BKHD chunk"))
	}

	AddChunkPosition(&bnk.Chunk, ChunkNameBKHD, 0)

	bnk.BKHD = bkhd
}

// Has side effect
func RegDIDXDATA(bnk *Bank, audioStore *AudioStore, pos u8) {
	if audioStore == nil {
		panic("didxdata is nil")
	}

	if HasChunk(&bnk.Chunk, ChunkNameDIDX) {
		panic(fmt.Sprintf("Duplicated DIDX chunk"))
	}

	AddChunkPosition(&bnk.Chunk, ChunkNameDIDX, pos)

	bnk.AudioStore = audioStore
}

// Has side effect
func RegHIRC(bnk *Bank, hirc *HIRC, pos u8) {
	if hirc == nil {
		panic("hirc is nil")
	}

	if HasChunk(&bnk.Chunk, ChunkNameHIRC) {
		panic(fmt.Sprintf("Duplicated HIRC chunk"))
	}

	AddChunkPosition(&bnk.Chunk, ChunkNameHIRC, pos)

	bnk.HIRC = hirc
}

// Has side effect
func AddEncodedChunk(c *ChunkComponent, chunkName string, pos u8, encoded []byte) {
	if encoded == nil {
		panic("Encoded slice is nil")
	}
	if _, in := c.Position[chunkName]; in {
		panic(fmt.Sprintf("Duplicate %s chunk", chunkName))
	}
	c.Position[chunkName] = pos
	if _, in := c.Encoded[chunkName]; in {
		panic(fmt.Sprintf("Duplicate encoded chunk %s", chunkName))
	}
	c.Encoded[chunkName] = encoded
}
