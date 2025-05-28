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

const (
	PropTypeVolume uint8 = 0 
	PropTypeLFE = 1 
	PropTypePitch = 2 
	PropTypeLPF = 3 
	PropTypeHPF = 4 
	PropTypeBusVolume = 5 
	PropTypeMakeUpGain = 6 
	PropTypePriority = 7 
	PropTypePriorityDistanceOffset = 8 
	PropTypeFeedbackVolumeUnused = 9 
	PropTypeFeedbackLPFUnused = 10 
	PropTypeMuteRatio = 11 
	PropTypePANLR = 12 
	PropTypePANFR = 13 
	PropTypeCenterPCT = 14 
	PropTypeDelayTime = 15 
	PropTypeTransitionTime = 16 
	PropTypeProbability = 17 
	PropTypeDialogueMode = 18 
	PropTypeUserAuxSendVolume0 = 19 
	PropTypeUserAuxSendVolume1 = 20 
	PropTypeUserAuxSendVolume2 = 21 
	PropTypeUserAuxSendVolume3 = 22 
	PropTypeGameAuxSendVolume = 23 
	PropTypeOutputBusVolume = 24 
	PropTypeOutputBusHPF = 25 
	PropTypeOutputBusLPF = 26 
	PropTypeHDRBusThreshold = 27 
	PropTypeHDRBusRatio = 28 
	PropTypeHDRBusReleaseTime = 29 
	PropTypeHDRBusGameParam = 30 
	PropTypeHDRBusGameParamMin = 31 
	PropTypeHDRBusGameParamMax = 32 
	PropTypeHDRActiveRange = 33 
	PropTypeLoopStart = 34 
	PropTypeLoopEnd = 35 
	PropTypeTrimInTime = 36 
	PropTypeTrimOutTime = 37 
	PropTypeFadeInTime = 38 
	PropTypeFadeOutTime = 39 
	PropTypeFadeInCurve = 40 
	PropTypeFadeOutCurve = 41 
	PropTypeLoopCrossfadeDuration = 42 
	PropTypeCrossfadeUpCurve = 43 
	PropTypeCrossfadeDownCurve = 44 
	PropTypeMidiTrackingRootNote = 45 
	PropTypeMidiPlayOnNoteType = 46 
	PropTypeMidiTransposition = 47 
	PropTypeMidiVelocityOffset = 48 
	PropTypeMidiKeyRangeMin = 49 
	PropTypeMidiKeyRangeMax = 50 
	PropTypeMidiVelocityRangeMin = 51 
	PropTypeMidiVelocityRangeMax = 52 
	PropTypeMidiChannelMask = 53 
	PropTypePlaybackSpeed = 54 
	PropTypeMidiTempoSource = 55 
	PropTypeMidiTargetNode = 56 
	PropTypeAttachedPluginFXID = 57 
	PropTypeLoop = 58 
	PropTypeInitialDelay = 59 
	PropTypeUserAuxSendLPF0 = 60 
	PropTypeUserAuxSendLPF1 = 61 
	PropTypeUserAuxSendLPF2 = 62 
	PropTypeUserAuxSendLPF3 = 63 
	PropTypeUserAuxSendHPF0 = 64 
	PropTypeUserAuxSendHPF1 = 65 
	PropTypeUserAuxSendHPF2 = 66 
	PropTypeUserAuxSendHPF3 = 67 
	PropTypeGameAuxSendLPF = 68 
	PropTypeGameAuxSendHPF = 69 
	PropTypeAttenuationID = 70 
	PropTypePositioningTypeBlend = 71 
	PropTypeReflectionBusVolume = 72 
	PropTypePAN_UD = 73 
)

var BasePropType []uint8 = []uint8{
	PropTypeVolume,
	PropTypePitch,
	PropTypeLPF,
	PropTypeHPF,
	PropTypeMakeUpGain,
	PropTypeGameAuxSendVolume,
	PropTypeInitialDelay,
}

var BaseRangePropType []uint8 = []uint8 {
	PropTypeVolume,
	PropTypePitch,
	PropTypeLPF,
	PropTypeHPF,
	PropTypeMakeUpGain,
	PropTypeInitialDelay,
}

var UserAuxSendVolumePropType []uint8 = []uint8 {
	PropTypeUserAuxSendVolume0,
	PropTypeUserAuxSendVolume1,
	PropTypeUserAuxSendVolume2,
	PropTypeUserAuxSendVolume3,
}
