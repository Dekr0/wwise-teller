package wwise

import "github.com/Dekr0/wwise-teller/wio"

type Event struct {
	HircObj

	Id   uint32
	Data []byte
}

func (h *Event) Encode() []byte {
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeEvent))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendBytes(h.Data)
	return w.BytesAssert(int(size))
}

func (h *Event) DataSize() uint32 {
	return uint32(4 + len(h.Data))
}

func (h *Event) BaseParameter() *BaseParameter { return nil }

func (h *Event) HircType() HircType { return HircTypeEvent }

func (h *Event) HircID() (uint32, error) { return h.Id, nil }

func (h *Event) IsCntr() bool { return false }

func (h *Event) NumLeaf() int { return 0 }

func (h *Event) ParentID() int { return 0 }

func (h *Event) AddLeaf(o HircObj) { panic("") }

func (h *Event) RemoveLeaf(o HircObj) { panic("") }
