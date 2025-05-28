package wwise

import "github.com/Dekr0/wwise-teller/wio"

type Action struct {
	HircObj

	Id      uint32
	Data    []byte
}

func (a *Action) Encode() []byte {
	dataSize := a.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeAction))
	w.Append(dataSize)
	w.Append(a.Id)
	w.AppendBytes(a.Data)
	return w.BytesAssert(int(size))
}

func (a *Action) DataSize() uint32 {
	return uint32(4 + len(a.Data))
}

func (a *Action) BaseParameter() *BaseParameter { return nil }

func (a *Action) HircType() HircType { return HircTypeAction }

func (a *Action) HircID() (uint32, error) { return a.Id, nil }

func (a *Action) IsCntr() bool { return false }

func (a *Action) NumLeaf() int { return 0 }

func (a *Action) ParentID() int { return 0 }

func (a *Action) AddLeaf(o HircObj) { panic("Panic Trap") }

func (a *Action) RemoveLeaf(o HircObj) { panic("Panic Trap") }
