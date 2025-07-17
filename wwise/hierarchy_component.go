// TODO
// - Channel Mask to string name
package wwise

import (
	"encoding/binary"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type BaseParameter struct {
	BitIsOverrideParentFx uint8
	FxChunk FxChunk
	FxChunkMetadata FxChunkMetadata
	BitOverrideAttachmentParams uint8 // <= 145
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
	cb := BaseParameter{
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
	return &cb
}

func (b *BaseParameter) Encode(v int) []byte {
	size := b.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendByte(b.BitIsOverrideParentFx)
	w.AppendBytes(b.FxChunk.Encode(v))
	w.AppendBytes(b.FxChunkMetadata.Encode(v))
	if v <= 145 {
		w.AppendByte(b.BitOverrideAttachmentParams)
	}
	w.Append(b.OverrideBusId)
	w.Append(b.DirectParentId)
	w.AppendByte(b.ByBitVectorA)
	w.AppendBytes(b.PropBundle.Encode(v))
	w.AppendBytes(b.RangePropBundle.Encode(v))
	w.AppendBytes(b.PositioningParam.Encode(v))
	w.AppendBytes(b.AuxParam.Encode(v))
	w.Append(b.AdvanceSetting)
	w.AppendBytes(b.StateProp.Encode(v))
	w.AppendBytes(b.StateGroup.Encode(v))
	w.AppendBytes(b.RTPC.Encode(v))
	return w.BytesAssert(int(size))
}

func (b *BaseParameter) Size(v int) uint32 {
	size := 1 + b.FxChunk.Size(v) + b.FxChunkMetadata.Size(v) + 1 + 4 + 4 + 1 + b.PropBundle.Size(v) + b.RangePropBundle.Size(v) + b.PositioningParam.Size(v) + b.AuxParam.Size(v) + SizeOfAdvanceSetting + b.StateProp.Size(v) + b.StateGroup.Size(v) + b.RTPC.Size(v)
	if v <= 145 {
		return size
	}
	return size - 1
}

func (b *BaseParameter) PriorityOverrideParent() bool {
	return wio.GetBit(b.ByBitVectorA, 0)
}

func (b *BaseParameter) SetPriorityOverrideParent(set bool, v int) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 0, set) 
	b.PropBundle.Add(TPriority, v)
}

func (b *BaseParameter) PriorityApplyDistFactor() bool {
	return wio.GetBit(b.ByBitVectorA, 1)
}

func (b *BaseParameter) SetPriorityApplyDistFactor(set bool, v int) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 1, set) 
	b.PropBundle.Add(TPriorityDistanceOffset, v)
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

func (b *BaseParameter) SetOverrideAuxSends(set bool, v int) {
	b.AuxParam.SetOverrideAuxSends(set)
	if !b.AuxParam.OverrideAuxSends() {
		b.PropBundle.RemoveAllUserAuxSendVolumeProp(v)
	}
}

func (b *BaseParameter) SetOverrideReflectionAuxBus(set bool, v int) {
	b.AuxParam.SetOverrideReflectionAuxBus(set)
	if !b.AuxParam.OverrideReflectionAuxBus() {
		b.PropBundle.Remove(TReflectionBusVolume, v)
	}
}

func (b *BaseParameter) SetEnableEnvelope(set bool, v int) {
	b.AdvanceSetting.SetEnableEnvelope(set)
	if b.AdvanceSetting.EnableEnvelope() {
		b.PropBundle.AddHDRActiveRange(v)
	}
}

type StateProp struct {
	NumStateProps    wio.Var
	StatePropItems []StatePropItem
}

func (s *StateProp) Clone() StateProp {
	return StateProp{
		NumStateProps: wio.Var{Bytes: slices.Clone(s.NumStateProps.Bytes), Value: s.NumStateProps.Value}, 
		StatePropItems: slices.Clone(s.StatePropItems),
	}
}

func (s *StateProp) Encode(v int) []byte {
	size := s.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendBytes(s.NumStateProps.Bytes)
	for _, si := range s.StatePropItems {
		w.AppendBytes(si.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (s *StateProp) Size(v int) uint32 {
	size := uint32(len(s.NumStateProps.Bytes))
	for _, i := range s.StatePropItems {
		size += i.Size(v)
	}
	return size
}

type StatePropItem struct {
	PropertyId wio.Var // var (at least 1 byte / 8 bits)
	AccumType  RTPCAccumType // U8x
	InDb       uint8 // U8x
}

func (s *StatePropItem) Encode(v int) []byte {
	b := slices.Clone(s.PropertyId.Bytes)
	b, _ = binary.Append(b, wio.ByteOrder, s.AccumType)
	b, _ = binary.Append(b, wio.ByteOrder, s.InDb)
	return b
}

func (s *StatePropItem) Size(v int) uint32 {
	return uint32(len(s.PropertyId.Bytes)) + 2
}

type StateGroup struct {
	NumStateGroups    wio.Var
	StateGroupItems []StateGroupItem
}

func (s *StateGroup) Clone() StateGroup {
	cs := StateGroup{
		NumStateGroups: wio.Var{
			Bytes: slices.Clone(s.NumStateGroups.Bytes),
			Value: s.NumStateGroups.Value,
		}, 
		StateGroupItems: make([]StateGroupItem, len(s.StateGroupItems)),
	}
	for i := range s.StateGroupItems {
		cs.StateGroupItems[i].StateGroupID = s.StateGroupItems[i].StateGroupID
		cs.StateGroupItems[i].StateSyncType = s.StateGroupItems[i].StateSyncType
		cs.StateGroupItems[i].States = slices.Clone(s.StateGroupItems[i].States)
	}
	return cs
}

func (s *StateGroup) Encode(v int) []byte {
	size := s.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendBytes(s.NumStateGroups.Bytes)
	for _, i := range s.StateGroupItems {
		w.AppendBytes(i.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (s *StateGroup) Size(v int) uint32 {
	size := uint32(len(s.NumStateGroups.Bytes))
	for _, i := range s.StateGroupItems {
		size += i.Size(v)
	}
	return size
}

type StateGroupItem struct {
	StateGroupID uint32 // tid
	StateSyncType uint8 // U8x
	NumStates wio.Var
	States []StateGroupItemState
}

func (s * StateGroupItem) Encode(v int) []byte {
	size := s.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(s.StateGroupID)
	w.AppendByte(s.StateSyncType)
	w.AppendBytes(s.NumStates.Bytes)
	for _, state := range s.States {
		w.Append(state.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (s *StateGroupItem) Size(v int) uint32 {
	size := 4 + 1 + uint32(len(s.NumStates.Bytes))
	for _, s := range s.States {
		size += s.Size(v)
	}
	return size
}

type StateGroupItemState struct {
	StateID uint32 // tid
	// The following section is exclusive
	// <= 145
	StateInstanceID uint32 // tid
	// > 145
	StatePropBundle StatePropBundle
}

func (s *StateGroupItemState) Size(v int) uint32 {
	if v <= 145 {
		return 8
	} else {
		return 4 + s.StatePropBundle.Size(v)
	}
}

func (s *StateGroupItemState) Encode(v int) []byte {
	size := s.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(s.StateID)
	if v <= 145 {
		w.Append(s.StateInstanceID)
	} else {
		w.AppendBytes(s.StatePropBundle.Encode(v))
	}
	return w.BytesAssert(int(size))
}

type Container struct {
	// NumChild u32
	Children []uint32 // NUmChild * sizeof(tid)
}

func NewCntrChildren() *Container {
	return &Container{[]uint32{}}
}

func (c *Container) Encode(v int) []byte {
	size := 4 + 4 * len(c.Children)
	w := wio.NewWriter(uint64(size))
	w.Append(uint32(len(c.Children)))
	w.Append(c.Children)
	return w.BytesAssert(int(size))
}

func (c *Container) Size(v int) uint32 {
	return uint32(4 + 4 * len(c.Children))
}
