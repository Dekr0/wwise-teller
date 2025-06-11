package wwise

import (
	"fmt"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type Unknown struct {
	HircObj
	Header *HircObjHeader
	Data   []byte
}

func NewUnknown(t HircType, s uint32, b []byte) *Unknown {
	return &Unknown{
		Header: &HircObjHeader{Type: t, Size: s},
		Data:   b,
	}
}

func (u *Unknown) Encode() []byte {
	assert.Equal(
		u.Header.Size,
		uint32(len(u.Data)),
		"Header size does not equal to actual data size",
	)

	bw := wio.NewWriter(uint64(SizeOfHircObjHeader + len(u.Data)))

	/* Header */
	bw.Append(u.Header)
	bw.AppendBytes(u.Data)

	return bw.Bytes()
}

func (u *Unknown) BaseParameter() *BaseParameter { return nil }

func (u *Unknown) HircID() (uint32, error) {
	return 0, fmt.Errorf("Hierarchy object type %d has yet implement GetHircID.", u.Header.Type)
}

func (u *Unknown) HircType() HircType { return u.Header.Type }

func (u *Unknown) IsCntr() bool { return false }

func (u *Unknown) NumLeaf() int { return 0 }

func (u *Unknown) ParentID() uint32 { return 0 }

func (u *Unknown) AddLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling AddLeaf.")
}

func (u *Unknown) RemoveLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling RemoveLeaf.")
}

func (u *Unknown) Leafs() []uint32 { return []uint32{} }
