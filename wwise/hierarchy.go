package wwise

import (
	"fmt"
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
)

type HIRC struct {
 	/* Used for memory allocation upfront during encoding */
	/*
	Future notes: Should I track individual byte change whenever a hierarchy obj
 	 changes?
	*/
	oldChunkSize uint32

	/*
	Currently, I don't know the algorithm of how Wwise encode its hierarchy tree.
	It's probably some sort of modified DFS since there are lot of places where 
    child nodes come first, and then the parent node of those child nodes come 
	right after it. So far now, I will book keep the linear order of hierarchy 
	tree as I parse them linearly through.
	*/
	HircObjs    []HircObj
	
	ActorMixers map[uint32]*ActorMixer
	LayerCntrs  map[uint32]*LayerCntr
	SwitchCntrs map[uint32]*SwitchCntr
	RanSeqCntrs map[uint32]*RanSeqCntr
	Sounds      map[uint32]*Sound
}

func NewHIRC(
	chunkSize uint32,
	numHircItem uint32,
) *HIRC {
	return &HIRC{
		oldChunkSize: chunkSize,
		HircObjs: make([]HircObj, numHircItem),
		ActorMixers: make(map[uint32]*ActorMixer),
		LayerCntrs: make(map[uint32]*LayerCntr),
		SwitchCntrs: make(map[uint32]*SwitchCntr),
		RanSeqCntrs: make(map[uint32]*RanSeqCntr),
		Sounds: make(map[uint32]*Sound),
	}
}

func (h *HIRC) Encode() ([]byte, error) {
	bw := reader.NewFixedSizeBlobWriter(uint64(CHUNK_HEADER_SIZE + h.oldChunkSize))
	bw.AppendBytes([]byte{'H', 'I', 'R', 'C'})
	/* gap chunk size field, and come back later */
	bw.AppendBytes([]byte{0, 0, 0, 0})
	bw.Append(uint32(len(h.HircObjs)))
	for _, hircObj := range h.HircObjs {
		bw.AppendBytes(hircObj.Encode())
	}

	assert.AssertTrue(bw.Len() - CHUNK_HEADER_SIZE > 0, "HIRC chunk size is less than or equal 0!")
	chunkSize := uint32(bw.Len() - CHUNK_HEADER_SIZE)
	hircBlob := bw.GetBlob()

	bw = reader.NewFixedSizeBlobWriter(4)
	bw.Append(chunkSize)
	chunkSizeBlob := bw.GetBlob()

	for i, _byte := range chunkSizeBlob {
		hircBlob[4 + i] = _byte
	}

	return hircBlob, nil
}

type HircObj interface {
	Encode() []byte
	GetHircID() (uint32, error)
	GetHircType() uint8 
}

const HIRC_OBJ_HEADER_SIZE = 1 + 4

type HircObjHeader struct {
	HircType uint8  // U8x
	HircSize uint32 // U32
}

type ActorMixer struct {
	HircId uint32
	BaseParam *BaseParameter
	Children *CntrChildren
}

func (s *ActorMixer) Encode() []byte {
	blob := s.BaseParam.Encode()
	blob = append(blob, s.Children.Encode()...)
	dataSize := uint32(4 + len(blob))
	blobSize := 1 + 4 + dataSize
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(s.GetHircType())
	bw.Append(dataSize)
	bw.Append(s.HircId)
	bw.AppendBytes(blob)
	return bw.Flush(int(blobSize))
}

func (s *ActorMixer) GetHircID() (uint32, error) {
	return s.HircId, nil
}

func (s *ActorMixer) GetHircType() uint8 {
	return 0x07
}

type LayerCntr struct {

}

func (s *LayerCntr) Encode() []byte {
	return []byte{}
}

func (s *LayerCntr) GetHircID() (uint32, error) {
	return 1, nil
}

func (s *LayerCntr) GetHircType() uint8 {
	return 0x09
}

type RanSeqCntr struct {

}

func (s *RanSeqCntr) Encode() []byte {
	return []byte{}
}

func (s *RanSeqCntr) GetHircID() (uint32, error) {
	return 1, nil
}

func (s *RanSeqCntr) GetHircType() uint8 {
	return 0x05
}

type SwitchCntr struct {

}

func (s *SwitchCntr) Encode() []byte {
	return []byte{}
}

func (s *SwitchCntr) GetHircID() (uint32, error) {
	return 1, nil
}

func (s *SwitchCntr) GetHircType() uint8 {
	return 0x07
}

type Sound struct {

}

func (s *Sound) Encode() []byte {
	return []byte{}
}

func (s *Sound) GetHircID() (uint32, error) {
	return 1, nil
}

func (s *Sound) GetHircType() uint8 {
	return 0x02
}

type Unknown struct {
	Header *HircObjHeader
	Blob   []byte
}

func (u *Unknown) Encode() []byte {
	assert.AssertEqual(
		u.Header.HircSize,
		uint32(len(u.Blob)),
		"Header size does not equal to actual data size",
	)

	bw := reader.NewFixedSizeBlobWriter(uint64(HIRC_OBJ_HEADER_SIZE + len(u.Blob)))
	
	/* Header */
	bw.Append(u.Header)
	bw.AppendBytes(u.Blob)

	return bw.GetBlob() 
}

func (u *Unknown) GetHircID() (uint32, error) {
	return 0, fmt.Errorf("Hierarchy object type %d has yet implement GetHircID.", u.Header.HircType)
}

func (u *Unknown) GetHircType() uint8 {
	return u.Header.HircType
}
