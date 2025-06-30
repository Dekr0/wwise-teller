package wwise

import "github.com/Dekr0/wwise-teller/wio"

type GroupType uint8
const (
	GroupTypeSwitch GroupType = 0
	GroupTypeState GroupType = 1
	GroupTypeCount GroupType = 2
)
var GroupTypeName = []string{
	"Switch", "State",
}

const SizeOfStateProp = 6
type State struct {
	HircObj

	StateID    uint32
	// cProps u16
	StateProps []struct{
		PID uint16
		Val float32
	}
}

func (s *State) Encode() []byte {
	dataSize := s.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeState))
	w.Append(dataSize)
	w.Append(s.StateID)
	w.Append(uint16(len(s.StateProps)))
	for _, sp := range s.StateProps {
		w.Append(sp)
	}
	return w.BytesAssert(int(size))
}

func (s *State) DataSize() uint32 {
	return 4 + 2 + uint32(len(s.StateProps)) * SizeOfStateProp
}

func (s *State) BaseParameter() *BaseParameter { return nil }

func (s *State) HircType() HircType { return HircTypeState }

func (s *State) HircID() (uint32, error) { return s.StateID, nil }

func (s *State) IsCntr() bool { return false }

func (s *State) NumLeaf() int { return 0 }

func (s *State) ParentID() uint32 { return 0 }

func (s *State) AddLeaf(o HircObj) { panic("State object cannot add leaf") }

func (s *State) RemoveLeaf(o HircObj) { panic("State object cannot remove leaf") }

func (s *State) Leafs() []uint32 { return []uint32{} }
