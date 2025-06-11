package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type FxShareSet struct {
	Id              uint32
	PluginTypeId    uint32
	// Present if PluginID >= 0
	PluginParam     *PluginParam
	// uNumBankData uint8
	MediaMap        []MediaMapItem
	RTPC            RTPC
	StateProp       StateProp
	StateGroup      StateGroup
	// NumValues    uint16
	PluginProps     []PluginProp
}

func (h *FxShareSet) HasParam() bool {
	return h.PluginTypeId >= 0
}

func (h *FxShareSet) assert() {
	if !h.HasParam() {
		assert.Nil(h.PluginParam,
			"Plugin Type ID indicate that there's no plugin parameter data.",
		)
	}
}

func (h *FxShareSet) Encode() []byte {
	h.assert()
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.Append(HircTypeFxShareSet)
	w.Append(dataSize)
	w.Append(h.Id)
	w.Append(h.PluginTypeId)
	if h.PluginParam != nil {
		w.AppendBytes(h.PluginParam.Encode())
	}
	w.Append(uint8(len(h.MediaMap)))
	for _, i := range h.MediaMap {
		w.Append(i)
	}
	w.AppendBytes(h.RTPC.Encode())
	w.AppendBytes(h.StateProp.Encode())
	w.AppendBytes(h.StateGroup.Encode())
	w.Append(uint16(len(h.PluginProps)))
	for _, p := range h.PluginProps {
		w.Append(p)
	}
	return w.BytesAssert(int(size))
}

func (h *FxShareSet) DataSize() uint32 {
	size := 8 + 1 + uint32(len(h.MediaMap)) * SizeOfMediaMapItem + h.RTPC.Size() + h.StateProp.Size() + h.StateGroup.Size() + 2 + uint32(len(h.PluginProps)) * SizeOfPluginProp
	if h.PluginParam != nil {
		size += h.PluginParam.Size()
	}
	return size
}

func (h *FxShareSet) PluginType() PluginType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginTypeInvalid
	}
	return PluginType((h.PluginTypeId >> 0) & 0x000F)
}

func (h *FxShareSet) PluginCompany () PluginCompanyType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginCompanyTypeInvalid
	}
	return PluginCompanyType((h.PluginTypeId >> 4) & 0x03FF)
}

func (h *FxShareSet) BaseParameter() *BaseParameter { return nil }

func (h *FxShareSet) HircType() HircType { return HircTypeFxShareSet }

func (h *FxShareSet) HircID() (uint32, error) { return h.Id, nil }

func (h *FxShareSet) IsCntr() bool { return false }

func (h *FxShareSet) NumLeaf() int { return 0 }

func (h *FxShareSet) ParentID() uint32 { return 0 }

func (h *FxShareSet) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *FxShareSet) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (h *FxShareSet) Leafs() []uint32 { return []uint32{} }
