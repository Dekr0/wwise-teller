package wwise

var PropLabel_140 = []string{
  "Volume",
  "LFE",
  "Pitch",
  "LPF",
  "HPF",
  "Bus Volume",
  "Make Up Gain",
  "Priority",
  "Priority Distance Offset",
  "Feedback Volume (Unused)",
  "Feedback LPF (Unused)",
  "Mute Ratio",
  "PAN LR",
  "PAN FR",
  "Center PCT",
  "Delay Time",
  "Transition Time",
  "Probability",
  "Dialogue Mode",
  "User AuxSend Volume 0",
  "User AuxSend Volume 1",
  "User AuxSend Volume 2",
  "User AuxSend Volume 3",
  "Game Aux Send Volume",
  "Output Bus Volume",
  "Output Bus HPF",
  "Output Bus LPF",
  "HDR Bus Threshold",
  "HDR Bus Ratio",
  "HDR Bus Release Time",
  "HDR Bus Game Param",
  "HDR Bus Game Param Min",
  "HDR Bus Game Param Max",
  "HDR Active Range",
  "Loop Start",
  "Loop End",
  "Trim In Time",
  "Trim Out Time",
  "Fade In Time",
  "Fade Out Time",
  "Fade In Curve",
  "Fade Out Curve",
  "Loop Crossfade Duration",
  "Crossfade Up Curve",
  "Crossfade Down Curve",
  "Midi Tracking Root Note",
  "Midi Play On Note Type",
  "Midi Transposition",
  "Midi Velocity Offset",
  "Midi Key Range Min",
  "Midi Key Range Max",
  "Midi Velocity Range Min",
  "Midi Velocity Range Max",
  "Midi Channel Mask",
  "Playback Speed",
  "Midi Tempo Source",
  "Midi Target Node",
  "Attached Plugin FX ID",
  "Loop",
  "Initial Delay",
  "User Aux Send LPF 0",
  "User Aux Send LPF 1",
  "User Aux Send LPF 2",
  "User Aux Send LPF 3",
  "User Aux Send HPF 0",
  "User Aux Send HPF 1",
  "User Aux Send HPF 2",
  "User Aux Send HPF 3",
  "Game Aux Send LPF",
  "Game Aux Send HPF",
  "Attenuation ID",
  "Positioning Type Blend",
  "Reflection Bus Volume",
  "PAN_UD",
}

type PropType uint8

const (
	PropTypeVolume PropType = 0 
	PropTypeLFE PropType = 1 
	PropTypePitch PropType = 2 
	PropTypeLPF PropType = 3 
	PropTypeHPF PropType = 4 
	PropTypeBusVolume PropType = 5 
	PropTypeMakeUpGain PropType = 6 
	PropTypePriority PropType = 7 
	PropTypePriorityDistanceOffset PropType = 8 
	PropTypeFeedbackVolumeUnused PropType = 9 
	PropTypeFeedbackLPFUnused PropType = 10 
	PropTypeMuteRatio PropType = 11 
	PropTypePANLR PropType = 12 
	PropTypePANFR PropType = 13 
	PropTypeCenterPCT PropType = 14 
	PropTypeDelayTime PropType = 15 
	PropTypeTransitionTime PropType = 16 
	PropTypeProbability PropType = 17 
	PropTypeDialogueMode PropType = 18 
	PropTypeUserAuxSendVolume0 PropType = 19 
	PropTypeUserAuxSendVolume1 PropType = 20 
	PropTypeUserAuxSendVolume2 PropType = 21 
	PropTypeUserAuxSendVolume3 PropType = 22 
	PropTypeGameAuxSendVolume PropType = 23 
	PropTypeOutputBusVolume PropType = 24 
	PropTypeOutputBusHPF PropType = 25 
	PropTypeOutputBusLPF PropType = 26 
	PropTypeHDRBusThreshold PropType = 27 
	PropTypeHDRBusRatio PropType = 28 
	PropTypeHDRBusReleaseTime PropType = 29 
	PropTypeHDRBusGameParam PropType = 30 
	PropTypeHDRBusGameParamMin PropType = 31 
	PropTypeHDRBusGameParamMax PropType = 32 
	PropTypeHDRActiveRange PropType = 33 
	PropTypeLoopStart PropType = 34 
	PropTypeLoopEnd PropType = 35 
	PropTypeTrimInTime PropType = 36 
	PropTypeTrimOutTime PropType = 37 
	PropTypeFadeInTime PropType = 38 
	PropTypeFadeOutTime PropType = 39 
	PropTypeFadeInCurve PropType = 40 
	PropTypeFadeOutCurve PropType = 41 
	PropTypeLoopCrossfadeDuration PropType = 42 
	PropTypeCrossfadeUpCurve PropType = 43 
	PropTypeCrossfadeDownCurve PropType = 44 
	PropTypeMidiTrackingRootNote PropType = 45 
	PropTypeMidiPlayOnNoteType PropType = 46 
	PropTypeMidiTransposition PropType = 47 
	PropTypeMidiVelocityOffset PropType = 48 
	PropTypeMidiKeyRangeMin PropType = 49 
	PropTypeMidiKeyRangeMax PropType = 50 
	PropTypeMidiVelocityRangeMin PropType = 51 
	PropTypeMidiVelocityRangeMax PropType = 52 
	PropTypeMidiChannelMask PropType = 53 
	PropTypePlaybackSpeed PropType = 54 
	PropTypeMidiTempoSource PropType = 55 
	PropTypeMidiTargetNode PropType = 56 
	PropTypeAttachedPluginFXID PropType = 57 
	PropTypeLoop PropType = 58 
	PropTypeInitialDelay PropType = 59 
	PropTypeUserAuxSendLPF0 PropType = 60 
	PropTypeUserAuxSendLPF1 PropType = 61 
	PropTypeUserAuxSendLPF2 PropType = 62 
	PropTypeUserAuxSendLPF3 PropType = 63 
	PropTypeUserAuxSendHPF0 PropType = 64 
	PropTypeUserAuxSendHPF1 PropType = 65 
	PropTypeUserAuxSendHPF2 PropType = 66 
	PropTypeUserAuxSendHPF3 PropType = 67 
	PropTypeGameAuxSendLPF PropType = 68 
	PropTypeGameAuxSendHPF PropType = 69 
	PropTypeAttenuationID PropType = 70 
	PropTypePositioningTypeBlend PropType = 71 
	PropTypeReflectionBusVolume PropType = 72 
	PropTypePAN_UD PropType = 73 
)

var BasePropType []PropType = []PropType{
	PropTypeVolume,
	PropTypePitch,
	PropTypeLPF,
	PropTypeHPF,
	PropTypeMakeUpGain,
	PropTypeGameAuxSendVolume,
	PropTypeInitialDelay,
}

var BaseRangePropType []PropType = []PropType {
	PropTypeVolume,
	PropTypePitch,
	PropTypeLPF,
	PropTypeHPF,
	PropTypeMakeUpGain,
	PropTypeInitialDelay,
}

var UserAuxSendVolumePropType []PropType = []PropType {
	PropTypeUserAuxSendVolume0,
	PropTypeUserAuxSendVolume1,
	PropTypeUserAuxSendVolume2,
	PropTypeUserAuxSendVolume3,
}
