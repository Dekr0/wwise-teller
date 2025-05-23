package wwise

import (
	"log/slog"

	"github.com/Dekr0/wwise-teller/wio"
)

type SwitchCntr struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	GroupType uint8 // U8x
	GroupID uint32 // tid
	DefaultSwitch uint32 // tid
	IsContinuousValidation uint8 // U8x
	Container *Container

	// NumSwitchGroups uint32 // u32

	SwitchGroups []*SwitchGroupItem

	// NumSwitchParams uint32 // u32

	SwitchParams []*SwitchParam
}

func (s *SwitchCntr) Encode() []byte {
	baseParamData := s.BaseParam.Encode()
	cntrData := s.Container.Encode()
	switchGroupDataSize := uint32(4)
	for _, i := range s.SwitchGroups {
		switchGroupDataSize += i.Size()
	}
	dataSize := 4 + uint32(len(baseParamData)) + 1 + 4 + 4 + 1 + 
				uint32(len(cntrData)) + switchGroupDataSize + 4 + 
				uint32(len(s.SwitchParams)) * SizeOfSwitchParam
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeSwitchCntr))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(baseParamData)
	w.AppendByte(s.GroupType)
	w.Append(s.GroupID)
	w.Append(s.DefaultSwitch)
	w.AppendByte(s.IsContinuousValidation)
	w.AppendBytes(cntrData)
	w.Append(uint32(len(s.SwitchGroups)))
	for _, i := range s.SwitchGroups {
		w.AppendBytes(i.Encode())
	}
	w.Append(uint32(len(s.SwitchParams)))
	for _, i := range s.SwitchParams {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (s *SwitchCntr) BaseParameter() *BaseParameter { return s.BaseParam }

func (s *SwitchCntr) HircID() (uint32, error) { return s.Id, nil }

func (s *SwitchCntr) HircType() HircType { return HircTypeSwitchCntr }

func (s *SwitchCntr) IsCntr() bool { return true }

func (s *SwitchCntr) NumLeaf() int { return len(s.Container.Children) }

func (s *SwitchCntr) ParentID() uint32 { return s.BaseParam.DirectParentId }

func (s *SwitchCntr) AddLeaf(o HircObj) {
	slog.Warn("Adding new leaf is not implemented for switch container.")
}

func (s *SwitchCntr) RemoveLeaf(o HircObj) {
	slog.Warn("Removing old leaf is not implemented for switch container.")
}

