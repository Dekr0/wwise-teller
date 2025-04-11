package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
)

const BKHD_FIELD_SIZE = 4 * 5

type BKHD struct {
	BankGenerationVersion uint32 // u32
	SoundbankID uint32 // sid
	LanguageID uint32 // sid

	/** Alignment and DeviceAllocated combine into one single U32 */
	Alignment uint16
	DeviceAllocated uint16

	ProjectID uint32 // u32
	Undefined []byte
	
	originalData []byte /* for testing */
}

func NewBKHD() *BKHD {
	return &BKHD{}
}

func (d *BKHD) Encode() []byte {
	chunkSize := uint32(BKHD_FIELD_SIZE + len(d.Undefined))
	bw := reader.NewFixedSizeBlobWriter(uint64(CHUNK_HEADER_SIZE + chunkSize))
	bw.AppendBytes([]byte{'B', 'K', 'H', 'D'})
	bw.Append(chunkSize);
	bw.Append(d.BankGenerationVersion); 
	bw.Append(d.SoundbankID)
	bw.Append(d.LanguageID)
	altValues := (uint32(d.DeviceAllocated) << 16) | (uint32(d.Alignment))
	bw.Append(altValues)
	bw.Append(d.ProjectID)
	bw.AppendBytes(d.Undefined)
	assert.AssertEqual(
		int(chunkSize),
		bw.Len() - 4 - 4,
		"(BKHD) The size of encoded data does not equal to calculated size.",
	)

	return bw.GetBlob()
}
