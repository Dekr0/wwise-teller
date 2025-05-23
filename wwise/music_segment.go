package wwise

import "github.com/Dekr0/wwise-teller/wio"

type MusicSegment struct {
	HircObj

	Id   uint32
	data []byte
}

func (h *MusicSegment) Encode() []byte {
	dataSize := h.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicSegment))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendBytes(h.data)
	return w.BytesAssert(int(size))
}

func (h *MusicSegment) DataSize() uint32 {
	return uint32(4 + len(h.data))
}

func (h *MusicSegment) BaseParameter() *BaseParameter { return nil }

func (h *MusicSegment) HircType() HircType { return HircTypeMusicSegment }

func (h *MusicSegment) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicSegment) IsCntr() bool { return false }

func (h *MusicSegment) NumLeaf() int { return 0 }

func (h *MusicSegment) ParentID() int { return 0 }

func (h *MusicSegment) AddLeaf(o HircObj) { panic("") }

func (h *MusicSegment) RemoveLeaf(o HircObj) { panic("") }
