package wwise

import (
	"fmt"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type RanSeqCntr struct {
	HircObj
	Id uint32
	BaseParam BaseParameter
	Container Container
	PlayListSetting PlayListSetting

	// NumPlayListItem u16

	PlayListItems []PlayListItem 
}

func (r *RanSeqCntr) Clone(id uint32, withParent bool) RanSeqCntr {
	return RanSeqCntr{
		Id: id,
		Container: Container{make([]uint32, 0, 16)},
		PlayListSetting: r.PlayListSetting.Clone(),
	}
}

func (r *RanSeqCntr) Encode() []byte {
	dataSize := r.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeRanSeqCntr))
	w.Append(dataSize)
	w.Append(r.Id)
	w.AppendBytes(r.BaseParam.Encode())
	w.Append(r.PlayListSetting)
	w.AppendBytes(r.Container.Encode())
	w.Append(uint16(len(r.PlayListItems)))
	for _, i := range r.PlayListItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (r *RanSeqCntr) DataSize() uint32 {
	return uint32(4 + r.BaseParam.Size() + r.Container.Size() + SizeOfPlayListSetting + 2 + uint32(len(r.PlayListItems)) * SizeOfPlayListItem)
}

func (r *RanSeqCntr) BaseParameter() *BaseParameter { return &r.BaseParam }

func (r *RanSeqCntr) HircID() (uint32, error) { return r.Id, nil }

func (r *RanSeqCntr) HircType() HircType { return HircTypeRanSeqCntr }

func (r *RanSeqCntr) IsCntr() bool { return true }

func (r *RanSeqCntr) NumLeaf() int { return len(r.Container.Children) }

func (r *RanSeqCntr) ParentID() uint32 { return r.BaseParam.DirectParentId }

func (r *RanSeqCntr) AddLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a hierarchy without a hierarchy ID.")
	}
	if slices.Contains(r.Container.Children, id) {
		panic(fmt.Sprintf("%d is already in random / sequence container %d", id, r.Id))
	} 
	b := o.BaseParameter()
	if b == nil {
		panic("%d is not containable.")
	}
	if b.DirectParentId != 0 {
		panic(fmt.Sprintf("%d is already attach to root %d. AddLeaf is an atomic operation.", id, b.DirectParentId))
	}
	r.Container.Children = append(r.Container.Children, id)
	if slices.ContainsFunc(
		r.PlayListItems,
		func(p PlayListItem) bool {
			return p.UniquePlayID == id
		},
	) {
		panic(fmt.Sprintf("%d is in playlist item of random / sequence container %d", id, r.Id))
	}
	b.DirectParentId = r.Id
}

func (r *RanSeqCntr) RemoveLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a hierarchy without a hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("%d is not containable.")
	}
	l := len(r.Container.Children)
	r.Container.Children = slices.DeleteFunc(
		r.Container.Children,
		func(c uint32) bool {
			return c == id
		},
	)
	if l <= len(r.Container.Children) {
		panic(fmt.Sprintf("%d is not in random / sequence container %d", id, r.Id))
	}
	r.PlayListItems = slices.DeleteFunc(
		r.PlayListItems,
		func(p PlayListItem) bool {
			return p.UniquePlayID == id
		},
	)
	b.DirectParentId = 0
}

func (r *RanSeqCntr) Leafs() []uint32 { return r.Container.Children }

func (r *RanSeqCntr) AddLeafToPlayList(i int) {
	if slices.ContainsFunc(r.PlayListItems, func(p PlayListItem) bool {
		return p.UniquePlayID == r.Container.Children[i]
	}) {
		return
	}
	r.PlayListItems = append(r.PlayListItems, PlayListItem{
		r.Container.Children[i], 50000,
	})
}

func (r *RanSeqCntr) MovePlayListItem(a int, b int) {
	r.PlayListItems[b], r.PlayListItems[a] = r.PlayListItems[a], r.PlayListItems[b]
}

func (r *RanSeqCntr) RemoveLeafFromPlayList(i int) {
	r.PlayListItems = slices.Delete(r.PlayListItems, i, i + 1)
}

func (r *RanSeqCntr) RemoveLeafsFromPlayList(ids []uint32) {
	for _, id := range ids {
		r.PlayListItems = slices.DeleteFunc(
			r.PlayListItems,
			func(p PlayListItem) bool { return id == p.UniquePlayID },
		)
	}
}
