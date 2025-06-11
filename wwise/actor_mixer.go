package wwise

import (
	"fmt"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type ActorMixer struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container *Container
}

func (a *ActorMixer) Encode() []byte {
	dataSize := a.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeActorMixer))
	w.Append(dataSize)
	w.Append(a.Id)
	w.AppendBytes(a.BaseParam.Encode())
	w.AppendBytes(a.Container.Encode())
	return w.BytesAssert(int(size))
}

func (a *ActorMixer) DataSize() uint32 {
	return uint32(4 + a.BaseParam.Size() + a.Container.Size())
}

func (a *ActorMixer) BaseParameter() *BaseParameter { return a.BaseParam }

func (a *ActorMixer) HircID() (uint32, error) { return a.Id, nil }

func (a *ActorMixer) HircType() HircType { return HircTypeActorMixer }

func (a *ActorMixer) IsCntr() bool { return true }

func (a *ActorMixer) NumLeaf() int { return len(a.Container.Children) }

func (a *ActorMixer) ParentID() uint32 { return a.BaseParam.DirectParentId }

func (a *ActorMixer) AddLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a hierarchy without a hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("%d is not containable")
	}
	if b.DirectParentId != 0 {
		panic(fmt.Sprintf("%d is already attach to root %d. AddLeaf is an atomic operation.", id, b.DirectParentId))
	}
	if slices.Contains(a.Container.Children, id) {
		panic(fmt.Sprintf("%d is already in actor mixer %d", id, a.Id))
	}
	a.Container.Children = append(a.Container.Children, id)
	b.DirectParentId = a.Id
}

func (a *ActorMixer) RemoveLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a hierarchy without a hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("%d is not containable")
	}
	l := len(a.Container.Children)
	a.Container.Children = slices.DeleteFunc(
		a.Container.Children, 
		func(c uint32) bool {
			return c == id
		},
	)
	if l <= len(a.Container.Children) {
		panic(fmt.Sprintf("%d is not in actor mixer %d", id, a.Id))
	}
	b.DirectParentId = 0
}

func (a *ActorMixer) Leafs() []uint32 { return a.Container.Children }
