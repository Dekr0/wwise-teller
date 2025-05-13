package wwise

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const maxEncodeRoutine = 8

// # of hierarchy object (uint32)
const sizeOfHIRCHeader = 4

type HircType uint8

const (
	HircTypeSound HircType = 0x02
	HircTypeRanSeqCntr HircType = 0x05
	HircTypeSwitchCntr HircType = 0x06
	HircTypeActorMixer HircType = 0x07
	HircTypeLayerCntr HircType = 0x09
)

var KnownHircType []HircType = []HircType{
	HircTypeSound,
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
}

var HircTypeName []string = []string{
	"",
	"State",
	"Sound",
	"Action",
	"Event",
	"Random / Sequence Container",
	"Switch Container",
	"Actor Mixer",
	"Bus",
	"Layer Container",
	"Music Segment",
	"Music Track",
	"Music Switch Container",
	"Music Random / Sequence Container",
	"Attenuation",
	"Dialogue Event",
    "FX Share Set",
    "FX Custom",
    "Aux Bus",
    "LFO Modulator",
    "Envelope Modulator",
    "Audio Device",
    "Time Modulator",
}

type HIRC struct {
	I uint8
	T []byte

	// Retain the tree structure that comes from the decoding with minimal 
	// modification
	HircObjs    []HircObj
	HircObjsMap map[uint32]HircObj
	
	// Map for different types of hierarchy objects. Each object is a pointer 
	// to a specific hierarchy object, which is also in `HircObjs`.
	ActorMixers map[uint32]*ActorMixer
	LayerCntrs  map[uint32]*LayerCntr
	SwitchCntrs map[uint32]*SwitchCntr
	RanSeqCntrs map[uint32]*RanSeqCntr
	Sounds      map[uint32]*Sound
}

func NewHIRC(I uint8, T []byte, numHircItem uint32) *HIRC {
	return &HIRC{
		I: I,
		T: T,
		HircObjs: make([]HircObj, numHircItem),
		HircObjsMap: make(map[uint32]HircObj),
		ActorMixers: make(map[uint32]*ActorMixer),
		LayerCntrs: make(map[uint32]*LayerCntr),
		SwitchCntrs: make(map[uint32]*SwitchCntr),
		RanSeqCntrs: make(map[uint32]*RanSeqCntr),
		Sounds: make(map[uint32]*Sound),
	}
}

func (h *HIRC) encode(ctx context.Context) ([]byte, error) {
	type result struct {
		i int
		b []byte
	}
	slices.Reverse(h.HircObjs)

	// No initialization since I want it to crash and catch encoding bugs
	results := make([][]byte, len(h.HircObjs))

	// sync signal
	c := make(chan *result, maxEncodeRoutine)

	// limit # of go routines running at the same time
	sem := make(chan struct{} , maxEncodeRoutine)

	done := 0
	i := 0
	for done < len(h.HircObjs) {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case r := <- c:
			results[r.i] = r.b
			done += 1
		case sem <- struct{}{}:
			if i < len(h.HircObjs) {
				j := i
				go func() {
					c <- &result{j, h.HircObjs[j].Encode()}
					<- sem
				}()	
				i += 1
			}
		default:
			if i < len(h.HircObjs) {
				results[i] = h.HircObjs[i].Encode()
				done += 1
				i += 1
			}
		}
	}
	
	return bytes.Join(results, []byte{}), nil
}

func (h *HIRC) Encode(ctx context.Context) ([]byte, error) {
	b, err := h.encode(ctx)
	if err != nil {
		return nil, err
	}

	dataSize := uint32(sizeOfHIRCHeader + len(b))
	size := chunkHeaderSize + dataSize
	w := wio.NewWriter(uint64(chunkHeaderSize + dataSize))
	w.AppendBytes(h.T)
	w.Append(dataSize)
	w.Append(uint32(len(h.HircObjs)))
	w.AppendBytes(b)
	return w.BytesAssert(int(size)), nil 
}

func (h *HIRC) Tag() []byte {
	return h.T
}

func (h *HIRC) Idx() uint8 {
	return h.I
}

func (h *HIRC) ChangeRoot(l, r HircObj, orid uint32) {
	lid, err := l.HircID()
	if err != nil {
		panic("Passing a leaf object without a hierarchy ID.")
	}
	b := l.BaseParameter() 
	if b == nil {
		panic(fmt.Sprintf("Leaf object %d is not a containable object.", lid))
	}
	rid, err := r.HircID()
	if err != nil {
		panic("Passing a root object without a hierarchy ID.")
	}
	if !r.IsCntr() {
		panic(fmt.Sprintf("Root object %d is not a container type object.", rid))
	}

	b.DirectParentId = rid

	r.AddLeaf(l)

	if orid != 0 {
		op, in := h.HircObjsMap[orid]
		if !in {
			panic(fmt.Sprintf("Old parent ID %d does not exist.", orid))
		}
		op.RemoveLeaf(l)
	}
}

type HircObj interface {
	Encode() []byte
	BaseParameter() (*BaseParameter)
	HircID() (uint32, error)
	HircType() HircType 
	IsCntr() bool
	NumLeaf() int
	ParentID() uint32
	// Modify DirectParentId, 
	// pre condition: o.DirectParentId == 0
	// post condition: o.DirectParentId == HircObj.HircID
	AddLeaf(o HircObj) 
	// Modify DirectParentId, 
	// pre condition: o.DirectParentId == HircObj.HircID
	// post condition: DirectParentId = 0
	RemoveLeaf(o HircObj) 
}

const sizeOfHircObjHeader = 1 + 4

type HircObjHeader struct {
	Type HircType // U8x
	Size uint32 // U32
}

type ActorMixer struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container *Container
}

func (a *ActorMixer) Encode() []byte {
	dataSize := a.DataSize()
	size := sizeOfHircObjHeader + dataSize
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
		panic("Passing a leaf object without hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("Leaf object %d is not a containable object.")
	}
	if b.DirectParentId != 0 {
		panic(fmt.Sprintf("Leaf object %d is already attach to root %d. AddLeaf is an atomic operation.", id, b.DirectParentId))
	}
	if slices.Contains(a.Container.Children, id) {
		panic(fmt.Sprintf("Leaf object %d is already in actor mixer %d", id, a.Id))
	}
	a.Container.Children = append(a.Container.Children, id)
	b.DirectParentId = a.Id
}

func (a *ActorMixer) RemoveLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a leaf object without hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("Leaf object %d is not a containable object.")
	}
	l := len(a.Container.Children)
	a.Container.Children = slices.DeleteFunc(
		a.Container.Children, 
		func(c uint32) bool {
			return c == id
		},
	)
	if l >= len(a.Container.Children) {
		panic(fmt.Sprintf("Leaf object is not actor mixer %d", id))
	}
	b.DirectParentId = 0
}

type LayerCntr struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container *Container

	// NumLayers uint32 // u32
	
	Layers []*Layer
	IsContinuousValidation uint8 // U8x
}

func (l *LayerCntr) Encode() []byte {
	dataSize := l.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeLayerCntr))
	w.Append(dataSize)
	w.Append(l.Id)
	w.AppendBytes(l.BaseParam.Encode())
	w.AppendBytes(l.Container.Encode())
	w.Append(uint32(len(l.Layers)))
	for _, i := range l.Layers {
		w.AppendBytes(i.Encode())
	}
	w.AppendByte(l.IsContinuousValidation)
	return w.BytesAssert(int(size))
}

func (l *LayerCntr) DataSize() uint32 {
	size := 4 + l.BaseParam.Size() + l.Container.Size() + 4
	for _, i := range l.Layers {
		size += i.Size()
	}
	return size + 1
}

func (l *LayerCntr) BaseParameter() *BaseParameter {
	return l.BaseParam
}

func (l *LayerCntr) HircID() (uint32, error) {
	return l.Id, nil
}

func (l *LayerCntr) HircType() HircType {
	return HircTypeLayerCntr
}

func (l *LayerCntr) IsCntr() bool { return true }

func (l *LayerCntr) NumLeaf() int { return len(l.Container.Children) }

func (l *LayerCntr) ParentID() uint32 { return l.BaseParam.DirectParentId }

func (l *LayerCntr) AddLeaf(o HircObj) {
	slog.Warn("Adding new leaf is not implemented for layer container.")
}

func (l *LayerCntr) RemoveLeaf(o HircObj) {
	slog.Warn("Removing old leaf is not implemented for layer container.")
}

type RanSeqCntr struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container *Container
	PlayListSetting *PlayListSetting

	// NumPlayListItem u16

	PlayListItems []*PlayListItem 
}

func (r *RanSeqCntr) Encode() []byte {
	dataSize := r.DataSize()
	size := sizeOfHircObjHeader + dataSize
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
	return uint32(4 + r.BaseParam.Size() + r.Container.Size() + sizeOfPlayListSetting + 2 + uint32(len(r.PlayListItems)) * sizeOfPlayListItem)
}

func (r *RanSeqCntr) BaseParameter() *BaseParameter { return r.BaseParam }

func (r *RanSeqCntr) HircID() (uint32, error) { return r.Id, nil }

func (r *RanSeqCntr) HircType() HircType { return HircTypeRanSeqCntr }

func (r *RanSeqCntr) IsCntr() bool { return true }

func (r *RanSeqCntr) NumLeaf() int { return len(r.Container.Children) }

func (r *RanSeqCntr) ParentID() uint32 { return r.BaseParam.DirectParentId }

func (r *RanSeqCntr) AddLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a leaf object without hierarchy ID.")
	}
	if slices.Contains(r.Container.Children, id) {
		panic(fmt.Sprintf("Leaf object %d is already in random / sequence container %d", id, r.Id))
	} 
	b := o.BaseParameter()
	if b == nil {
		panic("Leaf object %d is not a containable object.")
	}
	if b.DirectParentId != 0 {
		panic(fmt.Sprintf("Leaf object %d is already attach to root %d. AddLeaf is an atomic operation.", id, b.DirectParentId))
	}
	r.Container.Children = append(r.Container.Children, id)
	if slices.ContainsFunc(
		r.PlayListItems,
		func(p *PlayListItem) bool {
			return p.UniquePlayID == id
		},
	) {
		panic(fmt.Sprintf("Newly added leaf object %d is in playlist item of random / sequence container %d", id, r.Id))
	}
	b.DirectParentId = r.Id
}

func (r *RanSeqCntr) RemoveLeaf(o HircObj) {
	id, err := o.HircID()
	if err != nil {
		panic("Passing a leaf object without hierarchy ID.")
	}
	b := o.BaseParameter()
	if b == nil {
		panic("Leaf object %d is not a containable object.")
	}
	l := len(r.Container.Children)
	r.Container.Children = slices.DeleteFunc(
		r.Container.Children,
		func(c uint32) bool {
			return c == id
		},
	)
	if l >= len(r.Container.Children) {
		panic(fmt.Sprintf("Leaf object %d is not in random / sequence container %d", id, r.Id))
	}
	r.PlayListItems = slices.DeleteFunc(
		r.PlayListItems,
		func(p *PlayListItem) bool {
			return p.UniquePlayID == id
		},
	)
	b.DirectParentId = 0
}

func (r *RanSeqCntr) AddLeafToPlayList(i int) {
	if slices.ContainsFunc(r.PlayListItems, func(p *PlayListItem) bool {
		return p.UniquePlayID == r.Container.Children[i]
	}) {
		return
	}
	r.PlayListItems = append(r.PlayListItems, &PlayListItem{
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
			func(p *PlayListItem) bool { return id == p.UniquePlayID },
		)
	}
}

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
				uint32(len(s.SwitchParams)) * sizeOfSwitchParam
	size := sizeOfHircObjHeader + dataSize
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

func (r *SwitchCntr) IsCntr() bool { return true }

func (s *SwitchCntr) NumLeaf() int { return len(s.Container.Children) }

func (s *SwitchCntr) ParentID() uint32 { return s.BaseParam.DirectParentId }

func (s *SwitchCntr) AddLeaf(o HircObj) {
	slog.Warn("Adding new leaf is not implemented for switch container.")
}

func (s *SwitchCntr) RemoveLeaf(o HircObj) {
	slog.Warn("Removing old leaf is not implemented for switch container.")
}

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
	size := sizeOfHircObjHeader + dataSize
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

func (s *Sound) IsCntr() bool { return true }

func (s *Sound) NumLeaf() int { return 0 }

func (s *Sound) ParentID() uint32 { return s.BaseParam.DirectParentId }

func (s *Sound) AddLeaf(o HircObj) {
	panic("Sound is not a container type hierarchy object.")
}

func (s *Sound) RemoveLeaf(o HircObj) {
	panic("Sound is not a container type hierarchy object.")
}

type Unknown struct {
	HircObj
	Header *HircObjHeader
	b []byte
}

func NewUnknown(t HircType, s uint32, b []byte) *Unknown {
	return &Unknown{
		Header: &HircObjHeader{Type: t, Size: s},
		b: b,
	}
}

func (u *Unknown) Encode() []byte {
	assert.Equal(
		u.Header.Size,
		uint32(len(u.b)),
		"Header size does not equal to actual data size",
	)

	bw := wio.NewWriter(uint64(sizeOfHircObjHeader + len(u.b)))
	
	/* Header */
	bw.Append(u.Header)
	bw.AppendBytes(u.b)

	return bw.Bytes() 
}

func (u *Unknown) BaseParameter() *BaseParameter { return nil }

func (u *Unknown) HircID() (uint32, error) {
	return 0, fmt.Errorf("Hierarchy object type %d has yet implement GetHircID.", u.Header.Type)
}

func (u *Unknown) HircType() HircType { return u.Header.Type }

func (u *Unknown) IsCntr() bool { return true }

func (u *Unknown) NumLeaf() int { return 0 }

func (u *Unknown) ParentID() uint32 { return 0 }

func (u *Unknown) AddLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling AddLeaf.")
}

func (u *Unknown) RemoveLeaf(o HircObj) {
	panic("Work in development hierarchy object type is calling RemoveLeaf.")
}
