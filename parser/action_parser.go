package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

type ActionDispatch struct {
	ParamParser         ActionParamParser
	SpecificParamParser ActionSpecificParser
}

type ActionSpecificParser func(*wio.Reader, int) wwise.ActionSpecificParam

type ActionParamParser func(*wio.Reader, ActionSpecificParser, int) wwise.ActionParam

var ActionDispatchMux_L150 []uint16 = []uint16{
	27,
	0,  // 0x0100: CAkActionStop,
	1,  // 0x0200: CAkActionPause,
	2,  // 0x0300: CAkActionResume,
	4,  // 0x0400: CAkActionPlay,
	4,  // 0x0500: CAkActionPlayAndContinue, #early (removed in later versions)
	5,  // 0x0600: CAkActionMute,
	5,  // 0x0700: CAkActionMute,
	10, // 0x0800: CAkActionSetAkProp, #AkPropID_Pitch
	10, // 0x0900: CAkActionSetAkProp, #AkPropID_Pitch
	10, // 0x0A00: CAkActionSetAkProp, #(none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150) 
	10, // 0x0B00: CAkActionSetAkProp, #(none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150)
	10, // 0x0C00: CAkActionSetAkProp, #AkPropID_BusVolume
	10, // 0x0D00: CAkActionSetAkProp, #AkPropID_BusVolume
	10, // 0x0E00: CAkActionSetAkProp, #AkPropID_LPF
	10, // 0x0F00: CAkActionSetAkProp, #AkPropID_LPF
	11, // 0x1000: CAkActionUseState,
	11, // 0x1100: CAkActionUseState,
	12, // 0x1200: CAkActionSetState,
	13, // 0x1300: CAkActionSetGameParameter,
	13, // 0x1400: CAkActionSetGameParameter,
	14, // 0x1500: CAkActionEvent, #not in v150
	14, // 0x1600: CAkActionEvent, #not in v150
	14, // 0x1700: CAkActionEvent, #not in v150
	27, // 0x1800
	16, // 0x1900: CAkActionSetSwitch,
	19, // 0x1A00: CAkActionBypassFX,
	19, // 0x1B00: CAkActionBypassFX,
	20, // 0x1C00: CAkActionBreak,
	21, // 0x1D00: CAkActionTrigger,
	22, // 0x1E00: CAkActionSeek,
	23, // 0x1F00: CAkActionRelease,
	10, // 0x2000: CAkActionSetAkProp, #AkPropID_HPF
	24, // 0x2100: CAkActionPlayEvent,
	25, // 0x2200: CAkActionResetPlaylist,
	26, // 0x2300: CAkActionPlayEventUnknown, #normally not defined
	27, // 0x2400
	27, // 0x2500
	27, // 0x2600
	27, // 0x2700
	27, // 0x2800
	27, // 0x2900
	27, // 0x2A00
	27, // 0x2B00
	27, // 0x2C00
	27, // 0x2D00
	27, // 0x2E00
	27, // 0x2F00
	10, // 0x3000: CAkActionSetAkProp, #AkPropID_HPF
	18, // 0x3100: CAkActionSetFX,
	18, // 0x3200: CAkActionSetFX,
	19, // 0x3300: CAkActionBypassFX,
	19, // 0x3400: CAkActionBypassFX,
	19, // 0x3500: CAkActionBypassFX,
	19, // 0x3600: CAkActionBypassFX,
	19, // 0x3700: CAkActionBypassFX,
}

var ActionDispatchMux_GE150 []uint16 = []uint16{
	27,
	0,  // 0x0100: CAkActionStop,
	1,  // 0x0200: CAkActionPause,
	2,  // 0x0300: CAkActionResume,
	4,  // 0x0400: CAkActionPlay,
	4,  // 0x0500: CAkActionPlayAndContinue, #early (removed in later versions)
	5,  // 0x0600: CAkActionMute,
	5,  // 0x0700: CAkActionMute,
	10, // 0x0800: CAkActionSetAkProp, #AkPropID_Pitch
	10, // 0x0900: CAkActionSetAkProp, #AkPropID_Pitch
	10, // 0x0A00: CAkActionSetAkProp, #(none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150) 
	10, // 0x0B00: CAkActionSetAkProp, #(none) / AkPropID_Volume (~v145) / AkPropID_FirstRtpc (v150)
	10, // 0x0C00: CAkActionSetAkProp, #AkPropID_BusVolume
	10, // 0x0D00: CAkActionSetAkProp, #AkPropID_BusVolume
	10, // 0x0E00: CAkActionSetAkProp, #AkPropID_LPF
	10, // 0x0F00: CAkActionSetAkProp, #AkPropID_LPF
	11, // 0x1000: CAkActionUseState,
	11, // 0x1100: CAkActionUseState,
	12, // 0x1200: CAkActionSetState,
	13, // 0x1300: CAkActionSetGameParameter,
	13, // 0x1400: CAkActionSetGameParameter,
	14, // 0x1500: CAkActionEvent, #not in v150
	14, // 0x1600: CAkActionEvent, #not in v150
	14, // 0x1700: CAkActionEvent, #not in v150
	27, // 0x1800
	16, // 0x1900: CAkActionSetSwitch,
	20, // 0x1A00: CAkActionBreak,
	21, // 0x1B00: CAkActionTrigger,
	20, // 0x1C00: CAkActionBreak,
	21, // 0x1D00: CAkActionTrigger,
	22, // 0x1E00: CAkActionSeek,
	23, // 0x1F00: CAkActionRelease,
	10, // 0x2000: CAkActionSetAkProp, #AkPropID_HPF
	24, // 0x2100: CAkActionPlayEvent,
	25, // 0x2200: CAkActionResetPlaylist,
	26, // 0x2300: CAkActionPlayEventUnknown, #normally not defined
	27, // 0x2400
	27, // 0x2500
	27, // 0x2600
	27, // 0x2700
	27, // 0x2800
	27, // 0x2900
	27, // 0x2A00
	27, // 0x2B00
	27, // 0x2C00
	27, // 0x2D00
	27, // 0x2E00
	27, // 0x2F00
	10, // 0x3000: CAkActionSetAkProp, #AkPropID_HPF
	18, // 0x3100: CAkActionSetFX,
	18, // 0x3200: CAkActionSetFX,
	19, // 0x3300: CAkActionBypassFX,
	19, // 0x3400: CAkActionBypassFX,
	19, // 0x3500: CAkActionBypassFX,
	19, // 0x3600: CAkActionBypassFX,
	19, // 0x3700: CAkActionBypassFX,
}

// No change for SetExceptParams for version bumping from 141 to latest

// Done migrating
// ParseActionActiveParam
// ParseActionPlayParam
// ParseActionSetValueParam
// ParseActionSetSwitchParam
// ParseActionSetRTPCParam  
// ParseActionSetFXParam    
// ParseActionByPassFXParam 
// ParseActionSeekParam
// ParseActionReleaseParam
// ParseActionPlayEventParam

// Done migrating
// ParseActionNoSpecificParam
// ParseActionStopSpecificParam
// ParseActionPauseSpecificParam
// ParseActionResumeSpecificParam
// ParseActionSetPropSpecificParam
// ParseActionSetGameParameterSpecificParam

var ActionDispatchLUT []ActionDispatch = []ActionDispatch{
	{ParseActionActiveParam,    ParseActionStopSpecificParam            }, // 0 Stop
	{ParseActionActiveParam,    ParseActionPauseSpecificParam           }, // 1 Pause
	{ParseActionActiveParam,    ParseActionResumeSpecificParam          }, // 2 Resume
	{ParseActionPlayParam,      ParseActionNoSpecificParam              }, // 3 Play
	{ParseActionPlayParam,      ParseActionNoSpecificParam              }, // 4 Play And Continue
	{ParseActionSetValueParam , ParseActionNoSpecificParam              }, // 5 Mute
	{ParseActionSetValueParam , ParseActionNoSpecificParam              }, // 6 Set Pitch
	{ParseActionSetValueParam , ParseActionNoSpecificParam              }, // 7 Set Volume
	{ParseActionSetValueParam , ParseActionNoSpecificParam              }, // 8 Set LFE
	{ParseActionSetValueParam , ParseActionNoSpecificParam              }, // 9 Set LPF
	{ParseActionSetValueParam , ParseActionSetPropSpecificParam         }, // 10 Set Prop
	{ParseActionNoParam       , ParseActionNoSpecificParam              }, // 11 Use State
	{ParseActionSetStateParam , ParseActionNoSpecificParam              }, // 12 Set State
	{ParseActionSetValueParam , ParseActionSetGameParameterSpecificParam}, // 13 Set Game Parameter
	{ParseActionNoParam       , ParseActionNoSpecificParam              }, // 14 Event
	{ParseActionNoParam       , ParseActionNoSpecificParam              }, // 15 Duck
	{ParseActionSetSwitchParam, ParseActionNoSpecificParam              }, // 16 Set Switch
	{ParseActionSetRTPCParam  , ParseActionNoSpecificParam              }, // 17 Set RTPC
	{ParseActionSetFXParam    , ParseActionNoSpecificParam              }, // 18 Set FX
	{ParseActionByPassFXParam , ParseActionNoSpecificParam              }, // 19 Bypass FX
	{ParseActionNoParam       , ParseActionNoSpecificParam              }, // 20 Break
	{ParseActionNoParam       , ParseActionNoSpecificParam              }, // 21 Trigger
	{ParseActionSeekParam     , ParseActionNoSpecificParam              }, // 22 Seek
	{ParseActionReleaseParam  , ParseActionNoSpecificParam              }, // 23 Release
	{ParseActionPlayEventParam, ParseActionNoSpecificParam              }, // 24 Play Event
	{ParseActionActiveParam   , ParseActionResetPlayListSpecificParam   }, // 25 Reset Playlist
	{ParseActionPlayParam     , ParseActionNoSpecificParam              }, // 26 Play Event Unknown
	{ParseActionParamPanic    , ParseActionSpecificParamPanic           }, // 27 Invalid
}

func ParseActionParamPanic(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	panic("Unknown action parameter. Please use Hex editor investigate this")
}

func ParseActionSpecificParamPanic(r *wio.Reader, v int) wwise.ActionSpecificParam {
	panic("Unknonw action specific parameter. Please use Hex editor investigate this")
}

func ParseAction(size uint32, r *wio.Reader, v int) *wwise.Action {
	assert.Equal(0, r.Pos(), "Switch container parser position doesn't start at 0.")
	begin := r.Pos()

	a := wwise.Action{
		Id: r.U32Unsafe(),
		ActionType: wwise.ActionType(r.U16Unsafe()),
		IdExt: r.U32Unsafe(),
		IdExt4: r.U8Unsafe(),
	}
	ParsePropBundle(r, &a.PropBundle, v)
	ParseRangePropBundle(r, &a.RangePropBundle, v)

	var lutIdx uint16 = 0
	if v >= 150 {
		lutIdx = ActionDispatchMux_GE150[(a.ActionType & 0xFF00) >> 8]
	} else {
		lutIdx = ActionDispatchMux_L150[(a.ActionType & 0xFF00) >> 8]
	}
	dispatch := ActionDispatchLUT[lutIdx]
	a.ActionParam = dispatch.ParamParser(r, dispatch.SpecificParamParser, v)

	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)

	return &a
}

func ParseActionNoSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionNoSpecificParam{} 
}

func ParseActionStopSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionStopSpecificParam{BitVector: r.U8Unsafe()}
}

func ParseActionPauseSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionPauseSpecificParam{BitVector: r.U8Unsafe()}
}

func ParseActionResumeSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionResumeSpecificParam{BitVector: r.U8Unsafe()}
}

func ParseActionSetPropSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionSetPropSpecificParam{
		Prop: r.U8Unsafe(),
		Base: r.F32Unsafe(),
		Min: r.F32Unsafe(),
		Max: r.F32Unsafe(),
	}
}

func ParseActionSetGameParameterSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionSetGameParameterSpecificParam{
		ByPassTransition: r.U8Unsafe(),
		EnumValueMeaning: r.U8Unsafe(),
		Base: r.F32Unsafe(),
		Min: r.F32Unsafe(),
		Max: r.F32Unsafe(),
	}
}

func ParseActionResetPlayListSpecificParam(r *wio.Reader, v int) wwise.ActionSpecificParam {
	return &wwise.ActionResetPlayListSpecificParam{}
}

func ParseActionExceptParams(r *wio.Reader, e []wwise.ExceptParam, v int) {
	for i := range e {
		e[i].ID = r.U32Unsafe()
		e[i].IsBus = r.U8Unsafe()
	}
}

func ParseActionNoParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionNoParam{}
}

func ParseActionActiveParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	param := wwise.ActionActiveParam{
		EnumFadeCurve: wwise.InterpCurveType(r.U8Unsafe()),
		AkSpecificParam: p(r, v),
		ExceptionListSize: r.VarUnsafe(),
	}
	param.ExceptParams = make([]wwise.ExceptParam, param.ExceptionListSize.Value, param.ExceptionListSize.Value)
	ParseActionExceptParams(r, param.ExceptParams, v)
	return &param
}

func ParseActionPlayParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	actionPlayParam := &wwise.ActionPlayParam{
		EnumFadeCurve: wwise.InterpCurveType(r.U8Unsafe()),
		BankID: r.U32Unsafe(),
	}
	if v >= 144 {
		actionPlayParam.BankType = r.U32Unsafe()
	}
	return actionPlayParam
}

func ParseActionSetValueParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	param := &wwise.ActionSetValueParam{
		EnumFadeCurve: wwise.InterpCurveType(r.U8Unsafe()),
		AkSpecificParam: p(r, v),
		ExceptionListSize: r.VarUnsafe(),
	}
	param.ExceptParams = make([]wwise.ExceptParam, param.ExceptionListSize.Value, param.ExceptionListSize.Value)
	ParseActionExceptParams(r, param.ExceptParams, v)
	return param
}

func ParseActionSetStateParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionSetStateParam{
		StateGroupID: r.U32Unsafe(),
		TargetStateID: r.U32Unsafe(),
	}
}

func ParseActionSetSwitchParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionSetSwitchParam{
		SwitchGroupID: r.U32Unsafe(),
		SwitchStateID: r.U32Unsafe(),
	}
}

func ParseActionSetRTPCParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionSetRTPCParam{
		RTPCID: r.U32Unsafe(),
		RTPCValue: r.F32Unsafe(),
	}
}

func ParseActionSetFXParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	param := wwise.ActionSetFXParam{
		IsAudioDeviceElement: r.U8Unsafe(),
		SlotIndex: r.U8Unsafe(),
		FXID: r.U32Unsafe(),
		IsShared: r.U8Unsafe(),
		ExceptionListSize: r.VarUnsafe(),
	}
	param.ExceptParams = make([]wwise.ExceptParam, param.ExceptionListSize.Value, param.ExceptionListSize.Value)
	ParseActionExceptParams(r, param.ExceptParams, v)
	return &param
}

func ParseActionByPassFXParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	param := wwise.ActionByPassFXParam{
		IsByPass: r.U8Unsafe(),
		ByFxSolt: r.U8Unsafe(),
		ExceptionListSize: r.VarUnsafe(),
	}
	param.ExceptParams = make([]wwise.ExceptParam, param.ExceptionListSize.Value, param.ExceptionListSize.Value)
	ParseActionExceptParams(r, param.ExceptParams, v)
	return &param
}

func ParseActionSeekParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	param := wwise.ActionSeekParam{
		IsSeekRelativeDuration: r.U8Unsafe(),
		SeekValue: r.F32Unsafe(),
		SeekValueMin: r.F32Unsafe(),
		SeekValueMax: r.F32Unsafe(),
		SnapToNearestMark: r.U8Unsafe(),
		ExceptionListSize: r.VarUnsafe(),
	}
	param.ExceptParams = make([]wwise.ExceptParam, param.ExceptionListSize.Value, param.ExceptionListSize.Value)
	ParseActionExceptParams(r, param.ExceptParams, v)
	return &param
}

func ParseActionReleaseParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionReleaseParam{}
}

func ParseActionPlayEventParam(r *wio.Reader, p ActionSpecificParser, v int) wwise.ActionParam {
	return &wwise.ActionPlayEventParam{}
}
