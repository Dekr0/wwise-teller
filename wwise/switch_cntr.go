package wwise

import (
	"encoding/binary"
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
	Container Container

	// NumSwitchGroups uint32 // u32

	SwitchGroups []SwitchGroupItem

	// NumSwitchParams uint32 // u32

	SwitchParams []SwitchParam
}

type SwitchGroupItem struct {
	SwitchID uint32 // sid

	// ulNumItems uint32 // u32

	NodeList []uint32 // tid
}

func (s *SwitchGroupItem) Size(int) uint32 {
	return uint32(4 + 4 + len(s.NodeList) * 4)
}

func (s *SwitchGroupItem) Encode(v int) []byte { 
	size := uint64(4 + 4 + len(s.NodeList) * 4)
	w := wio.NewWriter(size)
	w.Append(s.SwitchID)
	w.Append(uint32(len(s.NodeList)))
	for _, i := range s.NodeList {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

const SizeOfSwitchParam_LE150 = 14
const SizeOfSwitchParam_G150 = 13

type SwitchParam struct {
	NodeId uint32 // tid
	PlayBackBitVector uint8 // U8x
	ModeBitVector uint8 // U8x <= 150
	FadeOutTime int32 // s32
	FadeInTime int32 // s32
}

func (s *SwitchParam) Size(v int) uint32 {
	if v <= 150 {
		return 14
	} else {
		return 13
	}
}

func (s *SwitchParam) Encode(v int) []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, s.NodeId)
	b, _ = binary.Append(b, wio.ByteOrder, s.PlayBackBitVector)
	if v <= 150 {
		b, _ = binary.Append(b, wio.ByteOrder, s.ModeBitVector)
	}
	b, _ = binary.Append(b, wio.ByteOrder, s.FadeOutTime)
	b, _ = binary.Append(b, wio.ByteOrder, s.FadeInTime)
	return b
}

func (s *SwitchCntr) Size(v int) uint32 {
	size := 4 + s.BaseParam.Size(v)
	size += 1 + 4 + 4 + 1
	size += s.Container.Size(v)
	size += 4
	for _, i := range s.SwitchGroups {
		size += i.Size(v)
	}
	size += 4
	if v <= 150 {
		size += uint32(len(s.SwitchParams)) * SizeOfSwitchParam_LE150
	} else {
		size += uint32(len(s.SwitchParams)) * SizeOfSwitchParam_G150
	}
	return size
}

func (s *SwitchCntr) Encode(v int) []byte {
	dataSize := s.Size(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeSwitchCntr))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(s.BaseParam.Encode(v))
	w.AppendByte(s.GroupType)
	w.Append(s.GroupID)
	w.Append(s.DefaultSwitch)
	w.AppendByte(s.IsContinuousValidation)
	w.AppendBytes(s.Container.Encode(v))
	w.Append(uint32(len(s.SwitchGroups)))
	for _, i := range s.SwitchGroups {
		w.AppendBytes(i.Encode(v))
	}
	w.Append(uint32(len(s.SwitchParams)))
	for _, i := range s.SwitchParams {
		w.AppendBytes(i.Encode(v))
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

func (s *SwitchCntr) Leafs() []uint32 { return s.Container.Children }
