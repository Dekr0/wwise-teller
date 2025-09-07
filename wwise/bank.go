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
	DIDX *DIDX
	HIRC *HIRC
}

func NewBank() *Bank {
	return &Bank{
		ChunkPosition: make(map[string]u8, 11),
		EncodedChunk: make(map[string][]byte, 7),
	}
}

// No side effect
// Thread safe
func BankHasChunk(b *Bank, name string) (in bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, in = b.ChunkPosition[name]
	return in 
}

// Has side effect
// Thread safe
func BankAddBKHD(bnk *Bank, bkhd *BKHD) error {
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
func BankAddDIDX(bnk *Bank, didx *DIDX, pos u8) error {
	if didx == nil {
		panic("didx is nil")
	}

	bnk.mu.Lock()
	defer bnk.mu.Unlock()

	if _, in := bnk.ChunkPosition["DIDX"]; in {
		return fmt.Errorf("Duplicated DIDX chunk")
	}
	bnk.ChunkPosition["DIDX"] = pos

	bnk.DIDX = didx

	return nil
}

// Has side effect
// Thread safe
func BankAddEncodedChunk(bnk *Bank, chunkName string, pos u8, encoded []byte) error {
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
