package wwise

import (
	"context"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const bkhdFieldSize = 4 * 5

type BKHD struct {
	I uint8
	T []byte
	BankGenerationVersion uint32 // u32
	SoundbankID uint32 // sid
	LanguageID uint32 // sid

	/** Alignment and DeviceAllocated combine into one single U32 */
	Alignment uint16
	DeviceAllocated uint16

	ProjectID uint32 // u32
	Undefined []byte
	
	oldData []byte /* for testing */
}

func NewBKHD(I uint8, T []byte) *BKHD {
	return &BKHD{I: I, T: T}
}

func (b *BKHD) Encode(ctx context.Context) ([]byte, error) {
	size := uint32(bkhdFieldSize + len(b.Undefined))
	w := wio.NewWriter(uint64(chunkHeaderSize + size))
	w.AppendBytes(b.T)
	w.Append(size);
	w.Append(b.BankGenerationVersion); 
	w.Append(b.SoundbankID)
	w.Append(b.LanguageID)
	altValues := (uint32(b.DeviceAllocated) << 16) | (uint32(b.Alignment))
	w.Append(altValues)
	w.Append(b.ProjectID)
	w.AppendBytes(b.Undefined)
	assert.Equal(
		int(size),
		w.Len() - 4 - 4,
		"(BKHD) The size of encoded data does not equal to calculated size.",
	)
	return w.Bytes(), nil
}

func (b *BKHD) Tag() []byte {
	return b.T
}

func (b *BKHD) Idx() uint8 {
	return b.I
}
