package wwise

import (
	"sync"

	"github.com/Dekr0/unwise/io"
)

type ActionType = u16

type ActionBasicData struct {
	Type   ActionType 
	Target u32 // IdExt
	IdExt4 u8
}

type ExceptionParam struct {
	Size    io.V128
	Id    []u32
	IsBus []u8
}

type ActionParamActive struct {
	LerpType LerpType
}

type ActionParamPlay struct {
	LerpType LerpType
	BankId   u32
	BankType u32 // >= 144
}

type ActionParamSetValue struct {
	LerpType LerpType
}

type ActionParamSetState struct {
	StateGroupId  u32
	TargetStateId u32
}

type ActionParamSetSwitch struct {
	SwitchGroupId u32
	SwitchStateId u32
}

type ActionParamSetRTPC struct {
	RTPCId    u32
	RTPCValue f32
}

type ActionParamSetFx struct {
	IsAudioDeviceElement u8
	SlotIdx              u8
	FxId                 u32
	IsShared             u8
}

type ActionParamByPassFx struct {
	IsByPass     u8
	BypassFxSlot u8
}

type ActionParamSeek struct {
	IsRelativeSeek   u8
	SeekValue        f32
	SeekValueMin     f32
	SeekValueMax     f32
	SnapToNearstMark u8
}

type ActionSpecificParamStop struct {
	BitVector u8
}

type ActionSpecificParamPause struct {
	BitVector u8
}

type ActionSpecificParamResume struct {
	BitVector u8
}

type ActionSpecificParamSetProp struct {
	Prop u8
	Base u32
	Min  f32
	Max  f32
}

type ActionSpecificParamSetGameParameter struct {
	ByPassTransition u8
	GameParamMeaning u8
	Base             f32
	Min              f32
	Max              f32
}

type ActionComponent struct {
	mu sync.Mutex

	// Components that are shared across all actions
	ActionBasicData map[u32]*ActionBasicData
	Properties      map[u32]*Prop
	RProperties     map[u32]*Prop

	// Components that are shared across among some actions but not all of them
	ExceptionParams map[u32]*ExceptionParam

	// Action Parameter
	ActionActiveParam   map[u32]*ActionParamActive
	ActionPlayParam     map[u32]*ActionParamPlay
	ActionSetValueParam map[u32]*ActionParamSetValue

	// Action Specific Parameter
	ActionSpecificParamStop             map[u32]*ActionSpecificParamStop
	ActionSpecificParamPause            map[u32]*ActionSpecificParamPause
	ActionSpecificParamResume           map[u32]*ActionSpecificParamResume
	ActionSpecificParamSetProp          map[u32]*ActionSpecificParamSetProp
	ActionSpecificParamSetGameParameter map[u32]*ActionSpecificParamSetGameParameter
}
