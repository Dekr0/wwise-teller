package wwise

import (
	"encoding/binary"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

type ActionType              uint16
type ActionParamType         uint8
type ActionSpecificParamType uint8

const ActionTypeCount = 42
var ActionTypeName []string = []string{
	"",
	"Stop",
	"Pause",
	"Resume",
	"Play",
	"Play And Continue", // early (removed in later versions)
	"Mute",
	"Mute",
	"Set Pitch", // AkPropID_Pitch
	"Set Pitch", // AkPropID_Pitch
	"Set Volume", // (none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150) 
	"Set Volume", // (none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150)
	"Set Bus Volume", // AkPropID_BusVolume
	"Set Bus Volume", // AkPropID_BusVolume
	"Set LPF", // AkPropID_LPF
	"Set LPF", // AkPropID_LPF
	"Use State",
	"Use State",
	"Set State",
	"Set Game Parameter",
	"Set Game Parameter",
	"Event", // not in v150
	"Event", // not in v150
	"Event", // not in v150
	"Set Switch",
	"Bypass FX",
	"Bypass FX",
	"Break",
	"Trigger",
	"Seek",
	"Release",
	"Set Prop HPF", // AkPropID_HPF
	"Play Event",
	"Reset Play list",
	"Play Event Unknown", // normally not defined
	"Set Prop HPF", // AkPropID_HPF
	"Set FX",
	"Set FX",
	"Bypass FX",
	"Bypass FX",
	"Bypass FX",
	"Bypass FX",
	"Bypass FX",
}

const (
	TypeActionNoParam        ActionParamType = 0
	TypeActionActiveParam    ActionParamType = 1
	TypeActionPlayParam      ActionParamType = 2
	TypeActionSetValueParam  ActionParamType = 3
	TypeActionSetStateParam  ActionParamType = 4
	TypeActionSetSwitchParam ActionParamType = 5
	TypeActionSetRTPCParam   ActionParamType = 6
	TypeActionSetFXParam     ActionParamType = 7
	TypeActionBypassFXParam  ActionParamType = 8
	TypeActionSeekParam      ActionParamType = 9
	TypeActionReleaseParam   ActionParamType = 10
	TypeActionPlayEventParam ActionParamType = 11
)

const (
	TypeActionNoSpecificParam               ActionSpecificParamType = 0
	TypeActionStopSpecificParam             ActionSpecificParamType = 1
	TypeActionPauseSpecificParam            ActionSpecificParamType = 2
	TypeActionResumeSpecificParam           ActionSpecificParamType = 3
	TypeActionSetPropSpecificParam          ActionSpecificParamType = 4
	TypeActionSetGameParameterSpecificParam ActionSpecificParamType = 5
	TypeActionResetPlayListSpecificParam    ActionSpecificParamType = 6
)

type ActionParam interface {
	Type()             ActionParamType
	SpecificParam()    ActionSpecificParam
	Size()             uint32
	Encode()         []byte
	Clone()            ActionParam
}

type ActionSpecificParam interface {
	Type()     ActionSpecificParamType
	Size()     uint32
	Encode() []byte
	Clone()    ActionSpecificParam
}

type Action struct {
	HircObj

	Id              uint32
	ActionType      ActionType
	IdExt           uint32
	IdExt4          uint8
	PropBundle      PropBundle
	RangePropBundle RangePropBundle
	ActionParam     ActionParam
}

func (a *Action) Clone(id uint32, target uint32) Action {
	return Action{
		Id: id,
		ActionType: a.ActionType,
		IdExt: target,
		IdExt4: a.IdExt4,
		PropBundle: a.PropBundle.Clone(),
		RangePropBundle: a.RangePropBundle.Clone(),
		ActionParam: a.ActionParam.Clone(),
	}
}

func (a *Action) Type() uint16 {
	return (uint16(a.ActionType) & 0xFF00) >> 8
}

func (a *Action) SetType(t uint16) {
	a.ActionType = ActionType((uint16(a.ActionType) & 0x00FF) | (t << 8))
}

func (a *Action) Encode() []byte {
	dataSize := a.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeAction))
	w.Append(dataSize)
	w.Append(a.Id)
	w.Append(a.ActionType)
	w.Append(a.IdExt)
	w.Append(a.IdExt4)
	w.AppendBytes(a.PropBundle.Encode())
	w.AppendBytes(a.RangePropBundle.Encode())
	w.AppendBytes(a.ActionParam.Encode())
	return w.BytesAssert(int(size))
}

func (a *Action) DataSize() uint32 {
	return 4 + 2 + 4 + 1 + a.PropBundle.Size() + a.RangePropBundle.Size() + a.ActionParam.Size()
}

func (a *Action) BaseParameter() *BaseParameter { return nil }

func (a *Action) HircType() HircType { return HircTypeAction }

func (a *Action) HircID() (uint32, error) { return a.Id, nil }

func (a *Action) IsCntr() bool { return false }

func (a *Action) NumLeaf() int { return 0 }

func (a *Action) ParentID() uint32 { return 0 }

func (a *Action) AddLeaf(o HircObj) { panic("Panic Trap") }

func (a *Action) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (a *Action) Leafs() []uint32 { return []uint32{} }

func (a *Action) IsBus() bool {
	return wio.GetBit(a.IdExt4, 0)
}

// Action Specific Param

type ActionNoSpecificParam struct {}

func (p *ActionNoSpecificParam) Type() ActionSpecificParamType {
	return TypeActionNoSpecificParam
}

func (p *ActionNoSpecificParam) Encode() []byte { return []byte{} }

func (p *ActionNoSpecificParam) Size() uint32 { return 0 }

func (p *ActionNoSpecificParam) Clone() ActionSpecificParam { return &ActionNoSpecificParam{} }

type ActionStopSpecificParam struct {
	BitVector uint8
}

func (p *ActionStopSpecificParam) Size() uint32 { return 1 }

func (p *ActionStopSpecificParam) Type() ActionSpecificParamType {
	return TypeActionStopSpecificParam
}

func (p *ActionStopSpecificParam) Encode() []byte { return []byte{p.BitVector} }

func (p *ActionStopSpecificParam) Clone() ActionSpecificParam {
	return &ActionStopSpecificParam{p.BitVector}
}

func (p *ActionStopSpecificParam) ApplyToStateTransition() bool {
	return wio.GetBit(p.BitVector, 1)
}

func (p *ActionStopSpecificParam) ApplyToDynamicSequence() bool {
	return wio.GetBit(p.BitVector, 2)
}

type ActionPauseSpecificParam struct {
	BitVector uint8
}

func (p *ActionPauseSpecificParam) Type() ActionSpecificParamType { 
	return TypeActionPauseSpecificParam
}

func (p *ActionPauseSpecificParam) Size() uint32 { return 1 }

func (p *ActionPauseSpecificParam) Encode() []byte { return []byte{p.BitVector} }

func (p *ActionPauseSpecificParam) Clone() ActionSpecificParam {
	return &ActionPauseSpecificParam{p.BitVector}
}

func (p *ActionPauseSpecificParam) IncludePendingResume() bool {
	return wio.GetBit(p.BitVector, 0)
}

func (p *ActionPauseSpecificParam) ApplyToStateTransition() bool {
	return wio.GetBit(p.BitVector, 1)
}

func (p *ActionPauseSpecificParam) ApplyToDynamicSequence() bool {
	return wio.GetBit(p.BitVector, 2)
}

type ActionResumeSpecificParam struct {
	BitVector uint8
}

func (p *ActionResumeSpecificParam) Type() ActionSpecificParamType { 
	return TypeActionResumeSpecificParam
}

func (p *ActionResumeSpecificParam) Size() uint32 { return 1 }

func (p *ActionResumeSpecificParam) Encode() []byte { return []byte{p.BitVector} }

func (p *ActionResumeSpecificParam) Clone() ActionSpecificParam {
	return &ActionResumeSpecificParam{p.BitVector}
}

func (p *ActionResumeSpecificParam) MasterResume() bool {
	return wio.GetBit(p.BitVector, 0)
}

func (p *ActionResumeSpecificParam) ApplyToStateTransition() bool {
	return wio.GetBit(p.BitVector, 1)
}

func (p *ActionResumeSpecificParam) ApplyToDynamicSequence() bool {
	return wio.GetBit(p.BitVector, 2)
}

type ActionSetPropSpecificParam struct {
	Prop uint8
	Base float32
	Min  float32
	Max  float32
}

func (p *ActionSetPropSpecificParam) Type() ActionSpecificParamType { 
	return TypeActionSetPropSpecificParam
}

func (p *ActionSetPropSpecificParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

func (p *ActionSetPropSpecificParam) Clone() ActionSpecificParam {
	return &ActionSetPropSpecificParam{
		p.Prop, p.Base, p.Min, p.Max,
	}
}

func (p *ActionSetPropSpecificParam) Size() uint32 { return 13 }

type ActionSetGameParameterSpecificParam struct {
	ByPassTransition uint8
	EnumValueMeaning uint8
	Base             float32
	Min              float32
	Max              float32
}

func (p *ActionSetGameParameterSpecificParam) Type() ActionSpecificParamType { 
	return TypeActionSetGameParameterSpecificParam
}

func (p *ActionSetGameParameterSpecificParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

func (p *ActionSetGameParameterSpecificParam) Size() uint32 { return 14 }

func (p *ActionSetGameParameterSpecificParam) Clone() ActionSpecificParam {
	return &ActionSetGameParameterSpecificParam{
		p.ByPassTransition, p.EnumValueMeaning, p.Base, p.Min, p.Max,
	}
}

type ActionResetPlayListSpecificParam struct {}

func (p *ActionResetPlayListSpecificParam) Type() ActionSpecificParamType { 
	return TypeActionResetPlayListSpecificParam
}

func (p *ActionResetPlayListSpecificParam) Encode() []byte { return []byte{} }

func (p *ActionResetPlayListSpecificParam) Size() uint32 { return 0 }

func (p *ActionResetPlayListSpecificParam) Clone() ActionSpecificParam {
	return &ActionResetPlayListSpecificParam{}
}

// End of Action Specific Param

// Action Parameter

type ActionNoParam struct {}

func (p *ActionNoParam) Size() uint32 { return 0 }

func (p *ActionNoParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionNoParam) Type() ActionParamType { return TypeActionNoParam }

func (p *ActionNoParam) Encode() []byte { return  []byte{} }

func (p *ActionNoParam) Clone() ActionParam {
	return &ActionNoParam{}
}

type ActionActiveParam struct {
	EnumFadeCurve    uint8
	AkSpecificParam  ActionSpecificParam
	ExceptParams   []ExceptParam
}

const SizeOfExceptParam = 5
type ExceptParam struct {
	ID    uint32
	IsBus uint8
}

func (p *ActionActiveParam) Type() ActionParamType { return TypeActionActiveParam }

func (p *ActionActiveParam) Size() uint32 {
	return 1 + p.AkSpecificParam.Size() + 1 + 5 * uint32(len(p.ExceptParams))
}

func (p *ActionActiveParam) SpecificParam() ActionSpecificParam { 
	return p.AkSpecificParam
}

func (p *ActionActiveParam) Clone() ActionParam {
	return &ActionActiveParam{
		p.EnumFadeCurve,
		p.AkSpecificParam.Clone(),
		slices.Clone(p.ExceptParams),
	}
}

func (p *ActionActiveParam) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.EnumFadeCurve)
	w.AppendBytes(p.AkSpecificParam.Encode())
	w.Append(uint8(len(p.ExceptParams)))
	for _, e := range p.ExceptParams {
		w.Append(e)
	}
	return w.BytesAssert(int(size))
}

type ActionPlayParam struct {
	EnumFadeCurve uint8
	BankID        uint32
}

func (p *ActionPlayParam) SpecificParam() ActionSpecificParam {
	return &ActionNoSpecificParam{}
}

func (p *ActionPlayParam) Size() uint32 {
	return 5
}

func (p *ActionPlayParam) Type() ActionParamType {
	return TypeActionPlayParam
}

func (p *ActionPlayParam) Clone() ActionParam {
	return &ActionPlayParam{p.EnumFadeCurve, p.BankID}
}

func (p *ActionPlayParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

type ActionSetValueParam struct {
	EnumFadeCurve uint8
	AkSpecificParam ActionSpecificParam
	ExceptParams  []ExceptParam
}

func (p *ActionSetValueParam) Size() uint32 {
	return 1 + p.AkSpecificParam.Size() + 1 + uint32(len(p.ExceptParams)) * SizeOfExceptParam
}

func (p *ActionSetValueParam) SpecificParam() ActionSpecificParam { return p.AkSpecificParam }

func (p *ActionSetValueParam) Type() ActionParamType { return TypeActionSetValueParam }

func (p *ActionSetValueParam) Clone() ActionParam {
	return &ActionSetValueParam{
		p.EnumFadeCurve,
		p.AkSpecificParam.Clone(),
		slices.Clone(p.ExceptParams),
	}
}

func (p *ActionSetValueParam) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.EnumFadeCurve)
	w.AppendBytes(p.AkSpecificParam.Encode())
	w.Append(uint8(len(p.ExceptParams)))
	for _, e := range p.ExceptParams {
		w.Append(e)
	}
	return w.BytesAssert(int(size))
}

type ActionSetStateParam struct {
	StateGroupID  uint32
	TargetStateID uint32
}

func (p *ActionSetStateParam) Type() ActionParamType { return TypeActionSetStateParam }

func (p *ActionSetStateParam) Size() uint32 { return 8 }

func (p *ActionSetStateParam) Clone() ActionParam {
	return &ActionSetStateParam{p.StateGroupID, p.TargetStateID}
}

func (p *ActionSetStateParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

func (p *ActionSetStateParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

type ActionSetSwitchParam struct {
	SwitchGroupID uint32
	SwitchStateID uint32
}

func (p *ActionSetSwitchParam) Type() ActionParamType { return TypeActionSetSwitchParam }

func (p *ActionSetSwitchParam) Size() uint32 { return 8 }

func (p *ActionSetSwitchParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionSetSwitchParam) Clone() ActionParam {
	return &ActionSetSwitchParam{p.SwitchGroupID, p.SwitchStateID}
}

func (p *ActionSetSwitchParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

type ActionSetRTPCParam struct {
	RTPCID    uint32
	RTPCValue float32
}

func (p *ActionSetRTPCParam) Type() ActionParamType { return TypeActionSetRTPCParam }

func (p *ActionSetRTPCParam) Size() uint32 { return 8 }

func (p *ActionSetRTPCParam) Clone() ActionParam {
	return &ActionSetRTPCParam{p.RTPCID, p.RTPCValue}
}

func (p *ActionSetRTPCParam) Encode() []byte {
	b, _ := binary.Append(nil, wio.ByteOrder, p)
	return b
}

func (p *ActionSetRTPCParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

type ActionSetFXParam struct {
	IsAudioDeviceElement uint8
	SlotIndex            uint8
	FXID                 uint32
	IsShared             uint8
	ExceptParams         []ExceptParam
}

func (p *ActionSetFXParam) Type() ActionParamType { return TypeActionSetFXParam }

func (p *ActionSetFXParam) Size() uint32 {
	return 1 + 1 + 4 + 1 + 1 + uint32(len(p.ExceptParams)) * SizeOfExceptParam
}

func (p *ActionSetFXParam) Clone() ActionParam {
	return &ActionSetFXParam{
		p.IsAudioDeviceElement,
		p.SlotIndex,
		p.FXID,
		p.IsShared,
		slices.Clone(p.ExceptParams),
	}
}

func (p *ActionSetFXParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionSetFXParam) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.IsAudioDeviceElement)
	w.Append(p.SlotIndex)
	w.Append(p.FXID)
	w.Append(p.IsShared)
	w.Append(uint8(len(p.ExceptParams)))
	for _, e := range p.ExceptParams {
		w.Append(e)
	}
	return w.BytesAssert(int(size))
}

type ActionByPassFXParam struct {
	IsByPass       uint8
	ByFxSolt       uint8
	ExceptParams []ExceptParam
}

func (p *ActionByPassFXParam) Size() uint32 {
	return 1 + 1 + 1 + uint32(len(p.ExceptParams)) * SizeOfExceptParam
}

func (p *ActionByPassFXParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionByPassFXParam) Type() ActionParamType { return TypeActionBypassFXParam }

func (p *ActionByPassFXParam) Clone() ActionParam {
	return &ActionByPassFXParam{p.IsByPass, p.ByFxSolt, slices.Clone(p.ExceptParams)}
}

func (p *ActionByPassFXParam) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.IsByPass)
	w.Append(p.ByFxSolt)
	w.Append(uint8(len(p.ExceptParams)))
	for _, e := range p.ExceptParams {
		w.Append(e)
	}
	return w.BytesAssert(int(size))
}

type ActionSeekParam struct {
	IsSeekRelativeDuration uint8
	SeekValue              float32
	SeekValueMin           float32
	SeekValueMax           float32
	SnapToNearestMark      uint8
	ExceptParams           []ExceptParam
}

func (p *ActionSeekParam) Size() uint32 {
	return 1 + 4 * 3 + 1 + 1 + uint32(len(p.ExceptParams))
}

func (p *ActionSeekParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionSeekParam) Type() ActionParamType { return TypeActionSeekParam }

func (p *ActionSeekParam) Clone() ActionParam {
	return &ActionSeekParam{
		p.IsSeekRelativeDuration,
		p.SeekValue,
		p.SeekValueMin,
		p.SeekValueMax,
		p.SnapToNearestMark,
		slices.Clone(p.ExceptParams),
	}
}

func (p *ActionSeekParam) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(p.IsSeekRelativeDuration)
	w.Append(p.SeekValue)
	w.Append(p.SeekValueMin)
	w.Append(p.SeekValueMax)
	w.Append(p.SnapToNearestMark)
	w.Append(uint8(len(p.ExceptParams)))
	for _, e := range p.ExceptParams {
		w.Append(e)
	}
	return w.BytesAssert(int(size))
}

type ActionReleaseParam struct {}

func (p *ActionReleaseParam) Size() uint32 { return 0 }

func (p *ActionReleaseParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionReleaseParam) Type() ActionParamType { return TypeActionReleaseParam }

func (p *ActionReleaseParam) Clone() ActionParam {
	return &ActionReleaseParam{}
}

func (p *ActionReleaseParam) Encode() []byte { return  []byte{} }

type ActionPlayEventParam struct {}

func (p *ActionPlayEventParam) Size() uint32 { return 0 }

func (p *ActionPlayEventParam) SpecificParam() ActionSpecificParam { return &ActionNoSpecificParam{} }

func (p *ActionPlayEventParam) Type() ActionParamType { return TypeActionPlayEventParam }

func (p *ActionPlayEventParam) Clone() ActionParam {
	return &ActionPlayEventParam{}
}

func (p *ActionPlayEventParam) Encode() []byte { return  []byte{} }
