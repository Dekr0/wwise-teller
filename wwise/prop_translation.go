package wwise

import "fmt"

type PropType uint8

// Translation Layer
const (
	TVolume PropType = iota
	TLFE
	TPitch
	TLPF
	THPF
	TMakeUpGain
	TInitialDelay

	TDelayTime

	TUserAuxSendVolume0
	TUserAuxSendVolume1
	TUserAuxSendVolume2
	TUserAuxSendVolume3
	TUserAuxSendLPF0
	TUserAuxSendLPF1
	TUserAuxSendLPF2
	TUserAuxSendLPF3
	TUserAuxSendHPF0
	TUserAuxSendHPF1
	TUserAuxSendHPF2
	TUserAuxSendHPF3

	TGameAuxSendVolume
	TGameAuxSendLPF
	TGameAuxSendHPF

	TBusVolume
	TOutputBusVolume
	TOutputBusHPF
	TOutputBusLPF
	TReflectionBusVolume

	THDRBusThreshold
    THDRBusRatio
    THDRBusReleaseTime
	THDRActiveRange
    THDRBusGameParam
    THDRBusGameParamMin
    THDRBusGameParamMax

    TMidiTrackingRootNote
    TMidiPlayOnNoteType
	TMidiTransposition
	TMidiVelocityOffset
    TMidiKeyRangeMin
    TMidiKeyRangeMax
    TMidiVelocityRangeMin
    TMidiVelocityRangeMax
    TMidiChannelMask
    TMidiTempoSource
    TMidiTargetNode

	TPANLR
	TPANFR

    TPositioningPanX2D
    TPositioningPanY2D
    TPositioningPanZ2D
    TPositioningPanX3D
    TPositioningPanY3D
    TPositioningPanZ3D
    TPositioningCenterPercent
    TPositioningTypeBlend
    TPositioningEnableAttenuation
    TPositioningConeAttenuationOnOff
    TPositioningConeAttenuation
    TPositioningConeLPF
    TPositioningConeHPF

    TBypassFX
    TBypassAllFX
    TBypassAllMetadata
    TMuteRatio

    TPlayMechanismSpecialTransitionsValue
	TTransitionTime

	TPriority
	TPriorityDistanceOffset
	TMaxNumInstances
	TProbability

	TDialogueMode

    TAvailable0
    TAvailable1
    TAvailable2

    TPlaybackSpeed
	TLoop

    TAttenuationDistanceScaling
    TAttenuationID
)

var TranslateName = map[PropType]string{
	TVolume: "Volume",
	TLFE: "LFE",
	TPitch: "Pitch",
	TLPF: "LPF",
	THPF: "HPF",
	TMakeUpGain: "Make Up Gain",
	TInitialDelay: "Initial Delay",

	TDelayTime: "Delay Time",

	TUserAuxSendVolume0: "User Aux Send Volume 0",
	TUserAuxSendVolume1: "User Aux Send Volume 1",
	TUserAuxSendVolume2: "User Aux Send Volume 2",
	TUserAuxSendVolume3: "User Aux Send Volume 3",
	TUserAuxSendLPF0: "User Aux Send LPF 0",
	TUserAuxSendLPF1: "User Aux Send LPF 1",
	TUserAuxSendLPF2: "User Aux Send LPF 2",
	TUserAuxSendLPF3: "User Aux Send LPF 3",
	TUserAuxSendHPF0: "User Aux Send HPF 0",
	TUserAuxSendHPF1: "User Aux Send HPF 1",
	TUserAuxSendHPF2: "User Aux Send HPF 2",
	TUserAuxSendHPF3: "User Aux Send HPF 3",

	TGameAuxSendVolume: "Game Aux Send Volume",
	TGameAuxSendLPF: "Game Aux Send LPF",
	TGameAuxSendHPF: "Game Aux Send HPF",

	TBusVolume: "Bus Volume",
	TOutputBusVolume: "Output Bus Volume",
	TOutputBusHPF: "Output Bus HPF",
	TOutputBusLPF: "Output Bus LPF",
	TReflectionBusVolume: "Reflection Bus Volume",

	THDRBusThreshold: "HDR Bus Threshold",
	THDRBusRatio: "HDR Bus Ratio",
	THDRBusReleaseTime: "HDR Bus Release Time",
	THDRActiveRange: "HDR Active Range",
	THDRBusGameParam: "HDR Bus Game Param",
	THDRBusGameParamMin: "HDR Bus Game Param Min",
	THDRBusGameParamMax: "HDR Bus Game Param Max",

	TMidiTrackingRootNote: "Midi Tracking Root Note",
	TMidiPlayOnNoteType: "Midi Play On Note Type",
	TMidiTransposition: "Midi Transposition",
	TMidiVelocityOffset: "Midi Velocity Offset",
	TMidiKeyRangeMin: "Midi Key Range Min",
	TMidiKeyRangeMax: "Midi Key Range Max",
	TMidiVelocityRangeMin: "Midi Velocity Range Min",
	TMidiVelocityRangeMax: "Midi Velocity Range Max",
	TMidiChannelMask: "Midi Channel Mask",
	TMidiTempoSource: "Midi Tempo Source",
	TMidiTargetNode: "Midi Target Node",

	TPANLR: "PAN LR",
	TPANFR: "PAN FR",
	TPositioningPanX2D: "Positioning Pan X 2D",
	TPositioningPanY2D: "Positioning Pan Y 2D",
	TPositioningPanZ2D: "Positioning Pan Z 2D",
	TPositioningPanX3D: "Positioning Pan X 3D",
	TPositioningPanY3D: "Positioning Pan Y 3D",
	TPositioningPanZ3D: "Positioning Pan Z 3D",
	TPositioningCenterPercent: "Positioning Center Percent",
	TPositioningTypeBlend: "Positioning Type Blend",
	TPositioningEnableAttenuation: "Positioning Enable Attenuation",
	TPositioningConeAttenuationOnOff: "Positioning Cone Attenuation On Off",
	TPositioningConeAttenuation: "Positioning Cone Attenuation",
	TPositioningConeLPF: "Positioning Cone LPF",
	TPositioningConeHPF: "Positioning Cone HPF",

	TBypassFX: "Bypass FX",
	TBypassAllFX: "Bypass All FX",
	TBypassAllMetadata: "Bypass All Metadata",

	TMuteRatio: "Mute Ratio",

	TPlayMechanismSpecialTransitionsValue: "Play Mechanism Special Transitions Value",

	TTransitionTime: "Transition Time",

	TPriority: "Priority",
	TPriorityDistanceOffset: "Priority Distance Offset",
	TMaxNumInstances: "Max Num Instances",
	TProbability: "Probability",

	TDialogueMode: "Dialogue Mode",


	TAvailable0: "Available 0",
	TAvailable1: "Available 1",
	TAvailable2: "Available 2",

	TPlaybackSpeed: "Playback Speed",

	TLoop: "Loop",

	TAttenuationDistanceScaling: "Attenuation Distance Scaling",
	TAttenuationID: "Attenuation ID",
}

type TranslatePair struct {
	T PropType
	O uint8
}

var TranslationV154 []TranslatePair = []TranslatePair{
	{TVolume,0},
	{TPitch,1},
	{TLPF,2},
	{THPF,3},
	{TMakeUpGain,5},
	{TInitialDelay,34},

	{TDelayTime,58},

	{TUserAuxSendVolume0,8},
	{TUserAuxSendVolume1,9},
	{TUserAuxSendVolume2,10},
	{TUserAuxSendVolume3,11},
	{TUserAuxSendLPF0,16},
	{TUserAuxSendLPF1,17},
	{TUserAuxSendLPF2,18},
	{TUserAuxSendLPF3,19},
	{TUserAuxSendHPF0,20},
	{TUserAuxSendHPF1,21},
	{TUserAuxSendHPF2,22},
	{TUserAuxSendHPF3,23},

	{TGameAuxSendVolume,12},
	{TGameAuxSendLPF,24},
	{TGameAuxSendHPF,25},

	{TBusVolume,4},
	{TOutputBusVolume,13},
	{TOutputBusHPF,14},
	{TOutputBusLPF,15},
	{TReflectionBusVolume,26},

	{THDRBusThreshold,27},
    {THDRBusRatio,28},
    {THDRBusReleaseTime,29},
	{THDRActiveRange,30},
    {THDRBusGameParam,62},
    {THDRBusGameParamMin,63},
    {THDRBusGameParamMax,64},

	{TMidiTrackingRootNote,65},
    {TMidiPlayOnNoteType,66},
	{TMidiTransposition,31},
	{TMidiVelocityOffset,32},
    {TMidiKeyRangeMin,67},
    {TMidiKeyRangeMax,68},
    {TMidiVelocityRangeMin,69},
    {TMidiVelocityRangeMax,70},
    {TMidiChannelMask,71},
    {TMidiTempoSource,72},
    {TMidiTargetNode,73},

    {TPositioningPanX2D,35},
    {TPositioningPanY2D,36},
    {TPositioningPanZ2D,37},
    {TPositioningPanX3D,38},
    {TPositioningPanY3D,39},
    {TPositioningPanZ3D,40},
    {TPositioningCenterPercent,41},
    {TPositioningTypeBlend,42},
    {TPositioningEnableAttenuation,43},
    {TPositioningConeAttenuationOnOff,44},
    {TPositioningConeAttenuation,45},
    {TPositioningConeLPF,46},
    {TPositioningConeHPF,47},

	{TBypassFX,48},
    {TBypassAllFX,49},
    {TBypassAllMetadata,54},
    {TMuteRatio,7},

	{TPlayMechanismSpecialTransitionsValue,55},
	{TTransitionTime,59},

	{TPriority,6},
	{TPriorityDistanceOffset,57},
	{TMaxNumInstances,53},
	{TProbability,60},

	{TDialogueMode,61},

	{TAvailable0,50},
    {TAvailable1,51},
    {TAvailable2,52},

	{TPlaybackSpeed,33},
	{TLoop,74},

    {TAttenuationDistanceScaling,56},
	{TAttenuationID,75},
}

var TranslationV128 []TranslatePair = []TranslatePair{
	{TVolume,0x00},
	{TPitch,0x02},
	{TLPF,0x03},
	{THPF,0x04},
	{TMakeUpGain,0x06},
	{TInitialDelay,0x3B},

	{TDelayTime,0x0F},

	{TUserAuxSendVolume0,0x13},
	{TUserAuxSendVolume1,0x14},
	{TUserAuxSendVolume2,0x15},
	{TUserAuxSendVolume3,0x16},
	{TUserAuxSendLPF0,0x3C},
	{TUserAuxSendLPF1,0x3D},
	{TUserAuxSendLPF2,0x3E},
	{TUserAuxSendLPF3,0x3F},
	{TUserAuxSendHPF0,0x40},
	{TUserAuxSendHPF1,0x41},
	{TUserAuxSendHPF2,0x42},
	{TUserAuxSendHPF3,0x43},

	{TGameAuxSendVolume,0x17},
	{TGameAuxSendLPF,0x44},
	{TGameAuxSendHPF,0x45},

	{TBusVolume,0x05},
	{TOutputBusVolume,0x18},
	{TOutputBusHPF,0x19},
	{TOutputBusLPF,0x1A},
	{TReflectionBusVolume,0x48},

	{THDRBusThreshold,0x1B},
    {THDRBusRatio,0X1C},
    {THDRBusReleaseTime,0x1D},
	{THDRActiveRange,0x21},
    {THDRBusGameParam,0x1E},
    {THDRBusGameParamMin,0x1F},
    {THDRBusGameParamMax,0x20},

    {TMidiTrackingRootNote,0x2D},
    {TMidiPlayOnNoteType,0x2E},
	{TMidiTransposition,0x2F},
	{TMidiVelocityOffset,0x30},
    {TMidiKeyRangeMin,0x31},
    {TMidiKeyRangeMax,0x32},
    {TMidiVelocityRangeMin,0x33},
    {TMidiVelocityRangeMax,0x34},
    {TMidiChannelMask,0x35},
    {TMidiTempoSource,0x37},
    {TMidiTargetNode,0x38},

    {TPositioningTypeBlend,0x47},

	{TMuteRatio,0x0B},

	{TTransitionTime,0x10},

	{TPriority,0x07},
	{TPriorityDistanceOffset,0x08},
	{TProbability,0x11},

	{TDialogueMode,0x12},

	{TPlaybackSpeed,0x36},
	{TLoop,0x3A},

	{TAttenuationID,0x46},
}

var ForwardTranslationV154 map[PropType]uint8 = nil
var InverseTranslationV154 map[uint8]PropType = nil
var ForwardTranslationV128 map[PropType]uint8 = nil
var InverseTranslationV128 map[uint8]PropType = nil

func InitTranslation() {
	ForwardTranslationV154 = make(map[PropType]uint8, len(TranslationV154))
	InverseTranslationV154 = make(map[uint8]PropType, len(TranslationV154))
	for _, t := range TranslationV154 {
		ForwardTranslationV154[t.T] = t.O
		InverseTranslationV154[t.O] = t.T
	}
	ForwardTranslationV128 = make(map[PropType]uint8, len(TranslationV128))
	InverseTranslationV128 = make(map[uint8]PropType, len(TranslationV128))
	for _, t := range TranslationV128 {
		ForwardTranslationV128[t.T] = t.O
		InverseTranslationV128[t.O] = t.T
	}
}

func ForwardTranslateProp(p PropType, v int) uint8 {
	if v < 150 {
		tp, in := ForwardTranslationV128[p]
		if !in {
			panic(fmt.Sprintf("Failed to translate version 128 property ID %d", p))
		}
		return tp
	}
	if v >= 150 && v < 154 {
		panic("Translation is not implemented between version 150 (inclusive) and version 154 (exclusive)")
	}
	if v >= 154 {
		tp, in := ForwardTranslationV154[p]
		if !in {
			panic(fmt.Sprintf("Failed to translate version 154 property ID %d", p))
		}
		return tp
	}
	panic(fmt.Sprintf("Forward Translation is not implemented for version %d", v))
}

func InverseTranslateProp(p uint8, v int) PropType {
	if v < 150 {
		tp, in := InverseTranslationV128[p]
		if !in {
			panic(fmt.Sprintf("Failed to inverse translate version 128 property ID %d", p))
		}
		return tp
	}
	if v >= 150 && v < 154 {
		panic("Inverse translation is not implemented between version 150 (inclusive) and version 154 (exclusive)")
	}
	if v >= 154 {
		tp, in := InverseTranslationV154[p]
		if !in {
			panic(fmt.Sprintf("Failed to inverse translate version 154 property ID %d", p))
		}
		return tp
	}
	panic(fmt.Sprintf("Inverse Translation is not implemented for version %d", v))
}

func PropLabel(p PropType) string {
	name, in := TranslateName[p]
	if !in {
		return fmt.Sprintf("Unknown %d", p)
	}
	return name
}

var BasePropType []PropType = []PropType{
	TVolume,
	TPitch,
	TLPF,
	THPF,
	TMakeUpGain,
	TGameAuxSendVolume,
	TInitialDelay,
}

var BaseRangePropType []PropType = []PropType {
	TVolume,
	TPitch,
	TLPF,
	THPF,
	TMakeUpGain,
	TInitialDelay,
}

var UserAuxSendVolumePropType []PropType = []PropType {
	TUserAuxSendVolume0,
	TUserAuxSendVolume1,
	TUserAuxSendVolume2,
	TUserAuxSendVolume3,
}
