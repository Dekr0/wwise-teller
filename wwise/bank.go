package wwise

import (
	"fmt"
)

type Bank struct {
	ChunkPosition map[string]u8
	EncodedChunk  map[string][]byte

	BKHD     *BKHD
	DIDXDATA *AudioStore
	HIRC     *HIRC
}

func AllocBank() *Bank {
	return &Bank{
		ChunkPosition: make(map[string]u8, 11),
		EncodedChunk: make(map[string][]byte, 7),
	}
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
func RegDIDXDATA(bnk *Bank, didxdata *AudioStore, pos u8) {
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
func RegHIRC(bnk *Bank, hirc *HIRC, pos u8) {
	if hirc == nil {
		panic("hirc is nil")
	}

	if _, in := bnk.ChunkPosition["HIRC"]; in {
		panic(fmt.Sprintf("Duplicated HIRC chunk"))
	}
	bnk.ChunkPosition["HIRC"] = pos

	bnk.HIRC = hirc
}

// Has side effect
// Thread safe
func AddEncodedChunk(bnk *Bank, chunkName string, pos u8, encoded []byte) {
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
