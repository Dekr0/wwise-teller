package wwise

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type Event struct {
	HircObj

	Id        uint32
	ActionIDs []uint32
}

func (h *Event) NewAction(actionID uint32) {
	if slices.Contains(h.ActionIDs, actionID) {
		return
	}
	h.ActionIDs = append(h.ActionIDs, actionID)
}

func (h *Event) Encode() []byte {
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeEvent))
	w.Append(dataSize)
	w.Append(h.Id)
	w.Append(uint8(len(h.ActionIDs)))
	for _, i := range h.ActionIDs {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (h *Event) DataSize() uint32 {
	return 4 + 1 + 4 * uint32(len(h.ActionIDs))
}

func (h *Event) BaseParameter() *BaseParameter { return nil }

func (h *Event) HircType() HircType { return HircTypeEvent }

func (h *Event) HircID() (uint32, error) { return h.Id, nil }

func (h *Event) IsCntr() bool { return false }

func (h *Event) NumLeaf() int { return 0 }

func (h *Event) ParentID() uint32 { return 0 }

func (h *Event) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *Event) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (h *Event) Leafs() []uint32 { return []uint32{} }
