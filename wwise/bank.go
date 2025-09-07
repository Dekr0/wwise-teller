package wwise

import (
	"fmt"
	"sync"
)

type Bank struct {
	mu sync.Mutex

	ChunkPosition map[string]u8
	EncodedChunk  map[string][]byte

	BKHD *BKHD
	HIRC *HIRC
}

func NewBank() *Bank {
	return &Bank{
		ChunkPosition: make(map[string]u8, 11),
	}
}

func BankHasChunk(b *Bank, name string) (in bool) {
	_, in = b.ChunkPosition[name]
	return in 
}

func BankAddBKHD(bnk *Bank, bkhd *BKHD) {
	if bkhd == nil {
		panic("bkhd is nil")
	}

	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	if _, in := bnk.ChunkPosition["BKHD"]; in {
		panic("Duplicated BKHD chunk")
	}
	bnk.ChunkPosition["BKHD"] = 0

	bnk.BKHD = bkhd
}

func BankAddEncodedChunk(bnk *Bank, chunkName string, pos u8, encoded []byte) error {
	bnk.mu.Lock()
	defer bnk.mu.Unlock()

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
