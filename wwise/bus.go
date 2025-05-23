package wwise

import "github.com/Dekr0/wwise-teller/wio"

type Bus struct {
	HircObj

	Id   uint32
	data []byte
}

func (h *Bus) Encode() []byte {
	dataSize := h.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeBus))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendBytes(h.data)
	return w.BytesAssert(int(size))
}

func (h *Bus) DataSize() uint32 {
	return uint32(4 + len(h.data))
}

func (h *Bus) BaseParameter() *BaseParameter { return nil }

func (h *Bus) HircType() HircType { return HircTypeEvent }

func (h *Bus) HircID() (uint32, error) { return h.Id, nil }

func (h *Bus) IsCntr() bool { return false }

func (h *Bus) NumLeaf() int { return 0 }

func (h *Bus) ParentID() int { return 0 }

func (h *Bus) AddLeaf(o HircObj) { panic("") }

func (h *Bus) RemoveLeaf(o HircObj) { panic("") }

