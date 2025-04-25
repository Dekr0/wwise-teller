package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const bkhdFieldSize = 4 * 5

type BKHD struct {
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

func NewBKHD() *BKHD {
	return &BKHD{}
}

func (b *BKHD) Encode() []byte {
	size := uint32(bkhdFieldSize + len(b.Undefined))
	w := wio.NewWriter(uint64(chunkHeaderSize + size))
	w.AppendBytes([]byte{'B', 'K', 'H', 'D'})
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
	return w.Bytes()
}
