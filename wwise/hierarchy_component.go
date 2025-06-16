// TODO
// - Channel Mask to string name
package wwise

import (
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type BaseParameter struct {
	BitIsOverrideParentFx uint8
	FxChunk FxChunk
	FxChunkMetadata FxChunkMetadata
	BitOverrideAttachmentParams uint8
	OverrideBusId uint32
	DirectParentId uint32
	ByBitVectorA uint8
	PropBundle PropBundle
	RangePropBundle RangePropBundle
	PositioningParam PositioningParam
	AuxParam AuxParam
	AdvanceSetting AdvanceSetting
	StateProp StateProp
	StateGroup StateGroup
	RTPC RTPC
}

func (b *BaseParameter) Clone(withParent bool) *BaseParameter {
	cb := &BaseParameter{
		BitIsOverrideParentFx: b.BitIsOverrideParentFx,
		FxChunk: b.FxChunk.Clone(),
		FxChunkMetadata: b.FxChunkMetadata.Clone(),
		BitOverrideAttachmentParams: b.BitOverrideAttachmentParams,
		OverrideBusId: b.OverrideBusId,
		DirectParentId: 0,
		ByBitVectorA: b.ByBitVectorA,
		PropBundle: b.PropBundle.Clone(),
		RangePropBundle: b.RangePropBundle.Clone(),
		PositioningParam: b.PositioningParam.Clone(),
		AuxParam: b.AuxParam.Clone(),
		AdvanceSetting: b.AdvanceSetting.Clone(),
		StateProp: b.StateProp.Clone(),
		StateGroup:b.StateGroup.Clone() ,
		RTPC: b.RTPC.Clone(),
	}
	if withParent {
		cb.DirectParentId = b.DirectParentId
	}
	return cb
}

func (b *BaseParameter) Encode() []byte {
	size := b.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(b.BitIsOverrideParentFx)
	w.AppendBytes(b.FxChunk.Encode())
	w.AppendBytes(b.FxChunkMetadata.Encode())
	w.AppendByte(b.BitOverrideAttachmentParams)
	w.Append(b.OverrideBusId)
	w.Append(b.DirectParentId)
	w.AppendByte(b.ByBitVectorA)
	w.AppendBytes(b.PropBundle.Encode())
	w.AppendBytes(b.RangePropBundle.Encode())
	w.AppendBytes(b.PositioningParam.Encode())
	w.AppendBytes(b.AuxParam.Encode())
	w.Append(b.AdvanceSetting)
	w.AppendBytes(b.StateProp.Encode())
	w.AppendBytes(b.StateGroup.Encode())
	w.AppendBytes(b.RTPC.Encode())
	return w.BytesAssert(int(size))
}

func (b *BaseParameter) Size() uint32 {
	return 1 + b.FxChunk.Size() + b.FxChunkMetadata.Size() + 1 + 4 + 4 + 1 + b.PropBundle.Size() + b.RangePropBundle.Size() + b.PositioningParam.Size() + b.AuxParam.Size() + SizeOfAdvanceSetting + b.StateProp.Size() + b.StateGroup.Size() + b.RTPC.Size()
}

func (b *BaseParameter) PriorityOverrideParent() bool {
	return wio.GetBit(b.ByBitVectorA, 0)
}

func (b *BaseParameter) SetPriorityOverrideParent(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 0, set) 
	b.PropBundle.AddPriority()
}

func (b *BaseParameter) PriorityApplyDistFactor() bool {
	return wio.GetBit(b.ByBitVectorA, 1)
}

func (b *BaseParameter) SetPriorityApplyDistFactor(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 1, set) 
	b.PropBundle.AddPriorityApplyDistFactor()
}

func (b *BaseParameter) OverrideMidiEventsBehavior() bool {
	return wio.GetBit(b.ByBitVectorA, 2)
}

func (b *BaseParameter) SetOverrideMidiEventsBehavior(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 2, set) 
}

func (b *BaseParameter) OverrideMidiNoteTracking() bool {
	return wio.GetBit(b.ByBitVectorA, 3)
}

func (b *BaseParameter) SetOverrideMidiNoteTracking(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 3, set) 
}

func (b *BaseParameter) EnableMidiNoteTracking() bool {
	return wio.GetBit(b.ByBitVectorA, 4)
}

func (b *BaseParameter) SetEnableMidiNoteTracking(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 4, set) 
}

func (b *BaseParameter) MidiBreakLoopOnNoteOff() bool {
	return wio.GetBit(b.ByBitVectorA, 5)
}

func (b *BaseParameter) SetMidiBreakLoopOnNoteOff(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 5, set) 
}

func (b *BaseParameter) SetOverrideAuxSends(set bool) {
	b.AuxParam.SetOverrideAuxSends(set)
	if !b.AuxParam.OverrideAuxSends() {
		b.PropBundle.RemoveAllUserAuxSendVolumeProp()
	}
}

func (b *BaseParameter) SetOverrideReflectionAuxBus(set bool) {
	b.AuxParam.SetOverrideReflectionAuxBus(set)
	if !b.AuxParam.OverrideReflectionAuxBus() {
		b.PropBundle.Remove(PropTypeReflectionBusVolume)
	}
}

func (b *BaseParameter) SetEnableEnvelope(set bool) {
	b.AdvanceSetting.SetEnableEnvelope(set)
	if b.AdvanceSetting.EnableEnvelope() {
		b.PropBundle.AddHDRActiveRange()
	}
}

type StateProp struct {
	// NumStateProps uint8
	StatePropItems []StatePropItem
}

func (s *StateProp) Clone() StateProp {
	return StateProp{slices.Clone(s.StatePropItems)}
}

func (s *StateProp) Encode() []byte {
	size := s.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(s.StatePropItems)))
	for _, si := range s.StatePropItems {
		w.Append(si)
	}
	return w.BytesAssert(int(size))
}

func (s *StateProp) Size() uint32 {
	return uint32(1 + len(s.StatePropItems) * SizeOfStatePropItem)
}

const SizeOfStatePropItem = 3
type StatePropItem struct {
	PropertyId RTPCParameterType // var (at least 1 byte / 8 bits)
	AccumType  RTPCAccumType // U8x
	InDb       uint8 // U8x
}

type StateGroup struct {
	// NumStateGroups uint8
	StateGroupItems []StateGroupItem
}

func (s *StateGroup) Clone() StateGroup {
	cs := StateGroup{make([]StateGroupItem, len(s.StateGroupItems))}
	for i := range s.StateGroupItems {
		cs.StateGroupItems[i].StateGroupID = s.StateGroupItems[i].StateGroupID
		cs.StateGroupItems[i].StateSyncType = s.StateGroupItems[i].StateSyncType
		cs.StateGroupItems[i].States = slices.Clone(s.StateGroupItems[i].States)
	}
	return cs
}

func NewStateGroup() *StateGroup {
	return &StateGroup{[]StateGroupItem{}}
}

func (s *StateGroup) Encode() []byte {
	size := s.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(s.StateGroupItems)))
	for _, i := range s.StateGroupItems {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (s *StateGroup) Size() uint32 {
	size := uint32(1)
	for _, i := range s.StateGroupItems {
		size += i.Size()
	}
	return size
}

type StateGroupItem struct {
	StateGroupID uint32 // tid
	StateSyncType uint8 // U8x
	// NumStates uint8 // var (assume at least 1 byte / 8 bits, can be more)
	States []StateGroupItemState // NumStates * sizeof(StateGroupItemState)
}

func NewStateGroupItem() *StateGroupItem {
	return &StateGroupItem{0, 0, []StateGroupItemState{}}
}

func (s * StateGroupItem) Encode() []byte {
	size := s.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(s.StateGroupID)
	w.AppendByte(s.StateSyncType)
	w.AppendByte(uint8(len(s.States)))
	for _, state := range s.States {
		w.Append(state)
	}
	return w.BytesAssert(int(size))
}

func (s *StateGroupItem) Size() uint32 {
	return uint32(4 + 1 + 1 + SizeOfStateGroupItem * len(s.States))
}

const SizeOfStateGroupItem = 8
type StateGroupItemState struct {
	StateID uint32 // tid
	StateInstanceID uint32 // tid
}

type Container struct {
	// NumChild u32
	Children []uint32 // NUmChild * sizeof(tid)
}

func NewCntrChildren() *Container {
	return &Container{[]uint32{}}
}

func (c *Container) Encode() []byte {
	size := 4 + 4 * len(c.Children)
	w := wio.NewWriter(uint64(size))
	w.Append(uint32(len(c.Children)))
	w.Append(c.Children)
	return w.BytesAssert(int(size))
}

func (c *Container) Size() uint32 {
	return uint32(4 + 4 * len(c.Children))
}

type SwitchGroupItem struct {
	SwitchID uint32 // sid

	// ulNumItems uint32 // u32

	NodeList []uint32 // tid
}

func (s *SwitchGroupItem) Size() uint32 {
	return uint32(4 + 4 + len(s.NodeList) * 4)
}

func (s *SwitchGroupItem) Encode() []byte { 
	size := uint64(4 + 4 + len(s.NodeList) * 4)
	w := wio.NewWriter(size)
	w.Append(s.SwitchID)
	w.Append(uint32(len(s.NodeList)))
	for _, i := range s.NodeList {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

const SizeOfSwitchParam = 14
type SwitchParam struct {
	NodeId uint32 // tid
	PlayBackBitVector uint8 // U8x
	ModeBitVector uint8 // U8x
	FadeOutTime int32 // s32
	FadeInTime int32 // s32
}
