package wwise

import "github.com/Dekr0/wwise-teller/wio"

type MusicSwitchCntr struct {
	Id                  uint32
	OverwriteFlags      uint8
	BaseParam           BaseParameter
	Children            Container
	MeterInfo           MeterInfo
	// NumStingers      uint32
	Stingers            []Stinger
	// NumRules         uint32
	TransitionRules     []MusicTransitionRule

	ContinuePlayBack    uint8

	DecisionTreeData    []byte
}

func (h *MusicSwitchCntr) Encode() []byte {
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicSwitchCntr))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendByte(h.OverwriteFlags)
	w.AppendBytes(h.BaseParam.Encode())
	w.AppendBytes(h.Children.Encode())
	w.Append(h.MeterInfo)
	w.Append(uint32(len(h.Stingers)))
	for _, s := range h.Stingers {
		w.Append(s)
	}
	w.Append(uint32(len(h.TransitionRules)))
	for _, t := range h.TransitionRules {
		w.AppendBytes(t.Encode())
	}
	w.AppendByte(h.ContinuePlayBack)
	w.AppendBytes(h.DecisionTreeData)
	return w.BytesAssert(int(size))
}

func (h *MusicSwitchCntr) DataSize() uint32 {
	size := 4 + 1 +
		h.BaseParam.Size() + 
		h.Children.Size() + 
		SizeOfMeterInfo + 
		4 + uint32(len(h.Stingers)) * SizeOfStinger
	size += 4
	for _, t := range h.TransitionRules {
		size += t.Size()
	}
	size += 1 + uint32(len(h.DecisionTreeData))
	return size
}

func (h *MusicSwitchCntr) BaseParameter() *BaseParameter { return &h.BaseParam }

func (h *MusicSwitchCntr) HircType() HircType { return HircTypeMusicSwitchCntr }

func (h *MusicSwitchCntr) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicSwitchCntr) IsCntr() bool { return true }

func (h *MusicSwitchCntr) NumLeaf() int { return len(h.Children.Children) }

func (h *MusicSwitchCntr) ParentID() uint32 { return h.BaseParam.DirectParentId }

func (h *MusicSwitchCntr) AddLeaf(o HircObj) { panic("") }

func (h *MusicSwitchCntr) RemoveLeaf(o HircObj) { panic("") }
