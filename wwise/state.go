package wwise

import "github.com/Dekr0/wwise-teller/wio"

type State struct {
	HircObj

	Id   uint32
	Data []byte
}

func (s *State) Encode() []byte {
	dataSize := s.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeState))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(s.Data)
	return w.BytesAssert(int(size))
}

func (s *State) DataSize() uint32 {
	return uint32(4 + len(s.Data))
}

func (s *State) BaseParameter() *BaseParameter { return nil }

func (s *State) HircType() HircType { return HircTypeState }

func (s *State) HircID() (uint32, error) { return s.Id, nil }

func (s *State) IsCntr() bool { return false }

func (s *State) NumLeaf() int { return 0 }

func (s *State) ParentID() int { return 0 }

func (s *State) AddLeaf(o HircObj) { panic("") }

func (s *State) RemoveLeaf(o HircObj) { panic("") }
