package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type BaseParameter struct {
	BitIsOverrideParentFx uint8
	FxChunk *FxChunk
	FxChunkMetadata *FxChunkMetadata
	BitOverrideAttachmentParams uint8
	OverrideBusId uint32
	DirectParentId uint32
	ByBitVectorA uint8
	PropBundle *PropBundle
	RangePropBundle *RangePropBundle
	PositioningParam *PositioningParam
	AuxParam *AuxParam
	AdvanceSetting *AdvanceSetting
	StateProp *StateProp
	StateGroup *StateGroup
	RTPC *RTPC
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
}

func (b *BaseParameter) PriorityApplyDistFactor() bool {
	return wio.GetBit(b.ByBitVectorA, 1)
}

func (b *BaseParameter) SetPriorityApplyDistFactor(set bool) {
	b.ByBitVectorA = wio.SetBit(b.ByBitVectorA, 1, set) 
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

type FxChunk struct {
	// UniqueNumFx uint8
	BitsFxByPass uint8
	FxChunkItems []*FxChunkItem
}

func NewFxChunk() *FxChunk {
	return &FxChunk{0, []*FxChunkItem{}}
}

func (f *FxChunk) Encode() []byte {
	if len(f.FxChunkItems) <= 0 {
		return []byte{0}
	}
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(f.FxChunkItems)))
	w.AppendByte(f.BitsFxByPass)
	for _, i := range f.FxChunkItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunk) Size() uint32 {
	if len(f.FxChunkItems) <= 0 {
		return 1
	}
	return uint32(1 + 1 + len(f.FxChunkItems) * SizeOfFxChunk)
}

const SizeOfFxChunk = 7
type FxChunkItem struct {
	UniqueFxIndex uint8 // u8i
	FxId uint32 // tid
	BitIsShareSet uint8 // U8x
	BitIsRendered uint8 // U8x
}

type FxChunkMetadata struct {
	BitIsOverrideParentMetadata uint8
	// UniqueNumFxMetadata uint8
	FxMetaDataChunkItems []FxChunkMetadataItem
}

func NewFxChunkMetadata() *FxChunkMetadata {
	return &FxChunkMetadata{0, []FxChunkMetadataItem{}}
}

func (f *FxChunkMetadata) Encode() []byte {
	size := f.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(f.BitIsOverrideParentMetadata)
	w.AppendByte(uint8(len(f.FxMetaDataChunkItems)))
	for _, i := range f.FxMetaDataChunkItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (f *FxChunkMetadata) Size() uint32 {
	return uint32(1 + 1 + len(f.FxMetaDataChunkItems) * SizeOfFxChunkMetadata)
}

const SizeOfFxChunkMetadata = 6
type FxChunkMetadataItem struct {
	UniqueFxIndex uint8 // u8i
	FxId uint32 // tid
	BitIsShareSet uint8 // U8x
}

type AuxParam struct {
	AuxBitVector     uint8    // U8x
	AuxIds           []uint32 // 4 * tid
	RestoreAuxIds    []uint32
	ReflectionAuxBus uint32   // tid
}

func NewAuxParam() *AuxParam {
	return &AuxParam{0, []uint32{}, []uint32{}, 0}
}

func (a *AuxParam) OverrideAuxSends() bool {
	return wio.GetBit(a.AuxBitVector, 2)
}

func (a *AuxParam) SetOverrideAuxSends(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 2, set)
}

func (a *AuxParam) HasAux() bool {
	return a.AuxBitVector & 0b00001000 != 0
}

func (a *AuxParam) SetHasAux(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 3, set)
	if !a.HasAux() {
		a.AuxIds = []uint32{}
	} else {
		a.AuxIds = []uint32{0, 0, 0, 0}
	}
}

func (a *AuxParam) OverrideReflectionAuxBus() bool {
	return wio.GetBit(a.AuxBitVector, 4)
}

func (a *AuxParam) SetOverrideReflectionAuxBus(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 4, set)
}

func (a *AuxParam) Encode() []byte {
	a.assert()

	if !a.HasAux() {
		size := a.Size()
		w := wio.NewWriter(uint64(size))
		w.AppendByte(a.AuxBitVector)
		w.Append(a.ReflectionAuxBus)
		return w.BytesAssert(int(size))
	}

	size := a.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(a.AuxBitVector)
	for _, id := range a.AuxIds { w.Append(id) }
	w.Append(a.ReflectionAuxBus)

	return w.BytesAssert(int(size))
}

func (a *AuxParam) Size() uint32 {
	if !a.HasAux() {
		return 5
	}
	return uint32(1 + 4 * 4 + 4)
}

func (a *AuxParam) assert() {
	if !a.HasAux() {
		assert.True(
			len(a.AuxIds) <= 0, 
			"Aux bit vector indicate there is no auxiliary send but # of Aux IDs" +
			" is not equal to 0.",
		)
		return
	}
	assert.True(
		len(a.AuxIds) == 4,
		"Each auxiliary parameter should only have 4 auxiliary IDs",
	)
}

const (
	VirtualQueueBehaviorFromBeginning = 0
	VirtualQueueBehaviorPlayFromElapsedTime = 1
	VirtualQueueBehaviorResume = 2
)
var VirtualQueueBehaviorString []string = []string{
	"From Beginning", "Play From Elapsed Time", "Resume",
}

const (
	BelowThresholdBehaviorContinueToPlay = 0
	BelowThresholdBehaviorKillVoice = 1
	BelowThresholdBehaviorSendToVirtualVoice = 2
)
var BelowThresholdBehaviorString []string = []string{
	"Continue To Play", "Kill Voice", "Send To Virtual Voice",
}

const SizeOfAdvanceSetting = 6
type AdvanceSetting struct {
	AdvanceSettingBitVector uint8 // U8x
	VirtualQueueBehavior uint8 // U8x
	MaxNumInstance uint16 // u16
	BelowThresholdBehavior uint8 // U8x
	HDRBitVector uint8 // U8x
}

func (a *AdvanceSetting) KillNewest() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 0)
}

func (a *AdvanceSetting) SetKillNewest(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 0, set)
}

func (a *AdvanceSetting) UseVirtualBehavior() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 1)
}

func (a *AdvanceSetting) SetUseVirtualBehavior(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 1, set)
}

func (a *AdvanceSetting) IgnoreParentMaxNumInst() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 3)
}

func (a *AdvanceSetting) SetIgnoreParentMaxNumInst(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 3, set)
}

func (a *AdvanceSetting) IsVVoicesOptOverrideParent() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 4)
}

func (a *AdvanceSetting) SetVVoicesOptOverrideParent(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 4, set)
}

func (a *AdvanceSetting) OverrideHDREnvelope() bool {
	return wio.GetBit(a.HDRBitVector, 0)
}

func (a *AdvanceSetting) SetOverrideHDREnvelope(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 0, set)
}

func (a *AdvanceSetting) OverrideAnalysis() bool {
	return wio.GetBit(a.HDRBitVector, 1)
}

func (a *AdvanceSetting) SetOverrideAnalysis(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 1, set)
}

func (a *AdvanceSetting) NormalizeLoudness() bool {
	return wio.GetBit(a.HDRBitVector, 2)
}

func (a *AdvanceSetting) SetNormalizeLoudness(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 2, set)
}

func (a *AdvanceSetting) EnableEnvelope() bool {
	return wio.GetBit(a.HDRBitVector, 3)
}

func (a *AdvanceSetting) SetEnableEnvelope(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 3, set)
}

type StateProp struct {
	// NumStateProps uint8
	StatePropItems []StatePropItem
}

func NewStateProp() *StateProp {
	return &StateProp{[]StatePropItem{}}
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
	PropertyId uint8 // var (at least 1 byte / 8 bits)
	AccumType uint8 // U8x
	InDb uint8 // U8x
}

type StateGroup struct {
	// NumStateGroups uint8
	StateGroupItems []*StateGroupItem
}

func NewStateGroup() *StateGroup {
	return &StateGroup{[]*StateGroupItem{}}
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
	States []*StateGroupItemState // NumStates * sizeof(StateGroupItemState)
}

func NewStateGroupItem() *StateGroupItem {
	return &StateGroupItem{0, 0, []*StateGroupItemState{}}
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

var SourceType []string = []string{
	"DATA",
	"Streaming",
	"Prefetch Streaming",
}
type BankSourceData struct {
	PluginID uint32 // U32
	StreamType uint8 // U8x
	SourceID uint32 // tid
	InMemoryMediaSize uint32 // U32
	SourceBits uint8 // U8x
	PluginParam *PluginParam
}

func (b *BankSourceData) PluginType() uint32 {
	return (b.PluginID >> 0) & 0x000F
}

func (b *BankSourceData) Company() uint32 {
	return (b.PluginID >> 4) & 0x03FF
}

func (b *BankSourceData) LanguageSpecific() bool {
	return b.SourceBits & 0b00000001 != 0
}

func (b *BankSourceData) Prefetch() bool {
	return b.SourceBits & 0b00000010 != 0
}

func (b *BankSourceData) NonCacheable() bool {
	return b.SourceBits & 0b00001000 != 0
}

func (b *BankSourceData) HasSource() bool {
	return b.SourceBits & 0b10000000 != 0
}

func (b *BankSourceData) HasParam() bool {
	return (b.PluginID & 0x0F) == 2
}

func (b *BankSourceData) ChangeSource(sid uint32, inMemoryMediaSize uint32) {
	b.SourceID = sid
	b.InMemoryMediaSize = inMemoryMediaSize
}

func (b *BankSourceData) Encode() []byte {
	b.assert()
	size := b.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(b.PluginID)
	w.AppendByte(b.StreamType)
	w.Append(b.SourceID)
	w.Append(b.InMemoryMediaSize)
	w.Append(b.SourceBits)
	if b.PluginParam != nil {
		w.AppendBytes(b.PluginParam.Encode())
	}
	return w.BytesAssert(int(size))
}

func (b *BankSourceData) Size() uint32 {
	size := uint32(4 + 1 + 4 + 4 + 1)
	if b.PluginParam != nil {
		size += b.PluginParam.Size()
	}
	return size
}

func (b *BankSourceData) assert() {
	if !b.HasParam() {
		assert.Nil(b.PluginParam,
			"Plugin type indicate that there's no plugin parameter data.",
		)
		return
	}
	// This make no sense???
	if b.PluginID == 0 {
		assert.Nil(b.PluginParam,
			"Plugin type indicate that there's no plugin parameter data.",
		)
	}
}

type PluginParam struct {
	PluginParamSize uint32 // U32
	PluginParamData []byte
}

func (p *PluginParam) Encode() []byte {
	assert.Equal(
		int(p.PluginParamSize),
		len(p.PluginParamData),
		"Plugin parameter size doesn't equal # of bytes in plugin parameter data",
	)
	size := 4 + len(p.PluginParamData)
	w := wio.NewWriter(uint64(size))
	w.Append(p.PluginParamSize)
	w.AppendBytes(p.PluginParamData)
	return w.BytesAssert(size)
}

func (p *PluginParam) Size() uint32 {
	return uint32(4 + len(p.PluginParamData))
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

const SizeOfPlayListSetting = 24

const (
	TransitionModeDisable = 0
	TransitionModeCrossFadeAmp = 1
	TransitionModeCrossFadePower = 2
	TransitionModeDelay = 3
	TransitionModeSampleAccurate = 4
	TransitionModeTriggerRate = 5
)
var TransitionModeString []string = []string{
	"Disabled",
	"Cross Fade Amp",
	"Cross Fade Power",
	"Delay",
	"Sample Accurate",
	"Trigger Rate",
}

const (
	RandomModeNormal = 0
	RandomModeShuffle = 1
)
var RandomModeString []string = []string{"Normal", "Shuffle"}

const (
	ModeRandom = 0
	ModeSequence = 1
)
var PlayListModeString []string = []string{"Random", "Sequence"}

type PlayListSetting struct {
	LoopCount uint16 // u16
	LoopModMin uint16 // u16
	LoopModMax uint16 // u16
	TransitionTime float32 // f32
	TransitionTimeModMin float32 // f32
	TransitionTimeModMax float32 // f32
	AvoidRepeatCount uint16 // u16
	TransitionMode uint8 // U8x
	RandomMode uint8 // U8x
	Mode uint8 // U8x

	// _bIsUsingWeight
	// bResetPlayListAtEachPlay
	// bIsRestartBackward
	// bIsContinuous
	// bIsGlobal
	PlayListBitVector uint8 // U8x
}

func (p *PlayListSetting) UsingWeight() bool {
	return wio.GetBit(p.PlayListBitVector, 0)
}

func (p *PlayListSetting) SetUsingWeight(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 0, set)
}

func (p *PlayListSetting) ResetPlayListAtEachPlay() bool {
	return wio.GetBit(p.PlayListBitVector, 1)
}

func (p *PlayListSetting) SetResetPlayListAtEachPlay(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 1, set)
}

func (p *PlayListSetting) RestartBackward() bool {
	return wio.GetBit(p.PlayListBitVector, 2)
}

func (p *PlayListSetting) SetRestartBackward(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 2, set)
}

func (p *PlayListSetting) Continuous() bool {
	return wio.GetBit(p.PlayListBitVector, 3)
}

func (p *PlayListSetting) SetContinuous(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 3, set)
}

func (p *PlayListSetting) Global() bool {
	return wio.GetBit(p.PlayListBitVector, 4)
}

func (p *PlayListSetting) SetGlobal(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 4, set)
}

const SizeOfPlayListItem = 8
type PlayListItem struct {
	UniquePlayID uint32 // tid
	Weight int32 // s32
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

type Layer struct {
	Id uint32 // tid
	InitialRTPC *RTPC
	RTPCId uint32 // tid
	RTPCType uint8 // U8x

	// NumAssoc uint32 // u32

	LayerRTPCs []*LayerRTPC
}

func (l *Layer) Encode() []byte {
	size := l.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(l.Id)
	w.AppendBytes(l.InitialRTPC.Encode())
	w.Append(l.RTPCId)
	w.AppendByte(l.RTPCType)
	w.Append(uint32(len(l.LayerRTPCs)))
	for _, i := range l.LayerRTPCs {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (l *Layer) Size() uint32 {
	size := uint32(4 + l.InitialRTPC.Size() + 4 + 1 + 4)
	for _, i := range l.LayerRTPCs {
		size += i.Size()
	}
	return size
}

type LayerRTPC struct {
	AssociatedChildID uint32 // tid

	// NumRTPCGraphPoints / CurveSize uint32 // u32

	RTPCGraphPoints []*RTPCGraphPoint
}

func (l *LayerRTPC) Encode() []byte {
	size := l.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(l.AssociatedChildID)
	w.Append(uint32(len(l.RTPCGraphPoints)))
	for _, i := range l.RTPCGraphPoints {
		w.Append(i.Enocde())
	}
	return w.BytesAssert(int(size))
}

func (l *LayerRTPC) Size() uint32 {
	return uint32(4 + 4 + len(l.RTPCGraphPoints) * SizeOfRTPCGraphPoint)
}
