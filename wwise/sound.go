package wwise

import (
	"github.com/Dekr0/wwise-teller/wio"
)

type Sound struct {
	HircObj
	Id uint32
	BankSourceData *BankSourceData
	BaseParam *BaseParameter
}

func (s *Sound) Encode() []byte {
	b := s.BankSourceData.Encode()
	b = append(b, s.BaseParam.Encode()...)
	dataSize := uint32(4 + len(b))
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeSound))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(b)
	return w.BytesAssert(int(size))
}

func (s *Sound) BaseParameter() *BaseParameter { return s.BaseParam }

func (s *Sound) HircID() (uint32, error) { return s.Id, nil }

func (s *Sound) HircType() HircType { return HircTypeSound }

func (s *Sound) IsCntr() bool { return false }

func (s *Sound) NumLeaf() int { return 0 }

func (s *Sound) ParentID() uint32 { return s.BaseParam.DirectParentId }

func (s *Sound) AddLeaf(o HircObj) {
	panic("Sound is not a container type hierarchy object.")
}

func (s *Sound) RemoveLeaf(o HircObj) {
	panic("Sound is not a container type hierarchy object.")
}

