package wwise

var ForwardTranslationV128 = map[PropType]uint8{
    TVolume: 0,
    TLFE: 1,
    TPitch: 2,
    TLPF: 3,
    THPF: 4,
    TMakeUpGain: 6,
    TInitialDelay: 59,
    TDelayTime: 15,
    TUserAuxSendVolume0: 19,
    TUserAuxSendVolume1: 20,
    TUserAuxSendVolume2: 21,
    TUserAuxSendVolume3: 22,
    TUserAuxSendLPF0: 60,
    TUserAuxSendLPF1: 61,
    TUserAuxSendLPF2: 62,
    TUserAuxSendLPF3: 63,
    TUserAuxSendHPF0: 64,
    TUserAuxSendHPF1: 65,
    TUserAuxSendHPF2: 66,
    TUserAuxSendHPF3: 67,
    TGameAuxSendVolume: 23,
    TGameAuxSendLPF: 68,
    TGameAuxSendHPF: 69,
    TBusVolume: 5,
    TOutputBusVolume: 24,
    TOutputBusHPF: 25,
    TOutputBusLPF: 26,
    TReflectionBusVolume: 72,
    THDRBusThreshold: 27,
    THDRBusRatio: 28,
    THDRBusReleaseTime: 29,
    THDRActiveRange: 33,
    THDRBusGameParam: 30,
    THDRBusGameParamMin: 31,
    THDRBusGameParamMax: 32,
    TMidiTrackingRootNote: 45,
    TMidiPlayOnNoteType: 46,
    TMidiTransposition: 47,
    TMidiVelocityOffset: 48,
    TMidiKeyRangeMin: 49,
    TMidiKeyRangeMax: 50,
    TMidiVelocityRangeMin: 51,
    TMidiVelocityRangeMax: 52,
    TMidiChannelMask: 53,
    TMidiTempoSource: 55,
    TMidiTargetNode: 56,
    TPANLR: 12,
    TPANFR: 13,
    TPANUD: 73,
    TCenterPCT: 14,
    TPositioningTypeBlend: 71,
    TFeedbackVolumeUnused: 9,
    TFeedbackLPFUnused: 10,
    TAttachedPluginFXID: 57,
    TMuteRatio: 11,
    TTransitionTime: 16,
    TTrimInTime: 36,
    TTrimOutTime: 37,
    TFadeInTime: 38,
    TFadeOutTime: 39,
    TFadeInCurve: 40,
    TFadeOutCurve: 41,
    TLoopCrossfadeDuration: 42,
    TCrossfadeUpCurve: 43,
    TCrossfadeDownCurve: 44,
    TPriority: 7,
    TPriorityDistanceOffset: 8,
    TProbability: 17,
    TDialogueMode: 18,
    TPlaybackSpeed: 54,
    TLoop: 58,
    TLoopStart: 34,
    TLoopEnd: 35,
    TAttenuationID: 70,
}
var InverseTranslationV128 = map[uint8]PropType{
    0: TVolume,
    1: TLFE,
    2: TPitch,
    3: TLPF,
    4: THPF,
    6: TMakeUpGain,
    59: TInitialDelay,
    15: TDelayTime,
    19: TUserAuxSendVolume0,
    20: TUserAuxSendVolume1,
    21: TUserAuxSendVolume2,
    22: TUserAuxSendVolume3,
    60: TUserAuxSendLPF0,
    61: TUserAuxSendLPF1,
    62: TUserAuxSendLPF2,
    63: TUserAuxSendLPF3,
    64: TUserAuxSendHPF0,
    65: TUserAuxSendHPF1,
    66: TUserAuxSendHPF2,
    67: TUserAuxSendHPF3,
    23: TGameAuxSendVolume,
    68: TGameAuxSendLPF,
    69: TGameAuxSendHPF,
    5: TBusVolume,
    24: TOutputBusVolume,
    25: TOutputBusHPF,
    26: TOutputBusLPF,
    72: TReflectionBusVolume,
    27: THDRBusThreshold,
    28: THDRBusRatio,
    29: THDRBusReleaseTime,
    33: THDRActiveRange,
    30: THDRBusGameParam,
    31: THDRBusGameParamMin,
    32: THDRBusGameParamMax,
    45: TMidiTrackingRootNote,
    46: TMidiPlayOnNoteType,
    47: TMidiTransposition,
    48: TMidiVelocityOffset,
    49: TMidiKeyRangeMin,
    50: TMidiKeyRangeMax,
    51: TMidiVelocityRangeMin,
    52: TMidiVelocityRangeMax,
    53: TMidiChannelMask,
    55: TMidiTempoSource,
    56: TMidiTargetNode,
    12: TPANLR,
    13: TPANFR,
    73: TPANUD,
    14: TCenterPCT,
    71: TPositioningTypeBlend,
    9: TFeedbackVolumeUnused,
    10: TFeedbackLPFUnused,
    57: TAttachedPluginFXID,
    11: TMuteRatio,
    16: TTransitionTime,
    36: TTrimInTime,
    37: TTrimOutTime,
    38: TFadeInTime,
    39: TFadeOutTime,
    40: TFadeInCurve,
    41: TFadeOutCurve,
    42: TLoopCrossfadeDuration,
    43: TCrossfadeUpCurve,
    44: TCrossfadeDownCurve,
    7: TPriority,
    8: TPriorityDistanceOffset,
    17: TProbability,
    18: TDialogueMode,
    54: TPlaybackSpeed,
    58: TLoop,
    34: TLoopStart,
    35: TLoopEnd,
    70: TAttenuationID,
}

var ForwardTranslationV154 = map[PropType]uint8{
    TVolume: 0,
    TPitch: 1,
    TLPF: 2,
    THPF: 3,
    TMakeUpGain: 5,
    TInitialDelay: 34,
    TDelayTime: 58,
    TUserAuxSendVolume0: 8,
    TUserAuxSendVolume1: 9,
    TUserAuxSendVolume2: 10,
    TUserAuxSendVolume3: 11,
    TUserAuxSendLPF0: 16,
    TUserAuxSendLPF1: 17,
    TUserAuxSendLPF2: 18,
    TUserAuxSendLPF3: 19,
    TUserAuxSendHPF0: 20,
    TUserAuxSendHPF1: 21,
    TUserAuxSendHPF2: 22,
    TUserAuxSendHPF3: 23,
    TGameAuxSendVolume: 12,
    TGameAuxSendLPF: 24,
    TGameAuxSendHPF: 25,
    TBusVolume: 4,
    TOutputBusVolume: 13,
    TOutputBusHPF: 14,
    TOutputBusLPF: 15,
    TReflectionBusVolume: 26,
    THDRBusThreshold: 27,
    THDRBusRatio: 28,
    THDRBusReleaseTime: 29,
    THDRActiveRange: 30,
    THDRBusGameParam: 62,
    THDRBusGameParamMin: 63,
    THDRBusGameParamMax: 64,
    TMidiTrackingRootNote: 65,
    TMidiPlayOnNoteType: 66,
    TMidiTransposition: 31,
    TMidiVelocityOffset: 32,
    TMidiKeyRangeMin: 67,
    TMidiKeyRangeMax: 68,
    TMidiVelocityRangeMin: 69,
    TMidiVelocityRangeMax: 70,
    TMidiChannelMask: 71,
    TMidiTempoSource: 72,
    TMidiTargetNode: 73,
    TPositioningPanX2D: 35,
    TPositioningPanY2D: 36,
    TPositioningPanZ2D: 37,
    TPositioningPanX3D: 38,
    TPositioningPanY3D: 39,
    TPositioningPanZ3D: 40,
    TPositioningCenterPercent: 41,
    TPositioningTypeBlend: 42,
    TPositioningEnableAttenuation: 43,
    TPositioningConeAttenuationOnOff: 44,
    TPositioningConeAttenuation: 45,
    TPositioningConeLPF: 46,
    TPositioningConeHPF: 47,
    TBypassFX: 48,
    TBypassAllFX: 49,
    TBypassAllMetadata: 54,
    TMuteRatio: 7,
    TPlayMechanismSpecialTransitionsValue: 55,
    TTransitionTime: 59,
    TPriority: 6,
    TPriorityDistanceOffset: 57,
    TMaxNumInstances: 53,
    TProbability: 60,
    TDialogueMode: 61,
    TAvailable0: 50,
    TAvailable1: 51,
    TAvailable2: 52,
    TPlaybackSpeed: 33,
    TLoop: 74,
    TAttenuationID: 75,
}
var InverseTranslationV154 = map[uint8]PropType{
    0: TVolume,
    1: TPitch,
    2: TLPF,
    3: THPF,
    5: TMakeUpGain,
    34: TInitialDelay,
    58: TDelayTime,
    8: TUserAuxSendVolume0,
    9: TUserAuxSendVolume1,
    10: TUserAuxSendVolume2,
    11: TUserAuxSendVolume3,
    16: TUserAuxSendLPF0,
    17: TUserAuxSendLPF1,
    18: TUserAuxSendLPF2,
    19: TUserAuxSendLPF3,
    20: TUserAuxSendHPF0,
    21: TUserAuxSendHPF1,
    22: TUserAuxSendHPF2,
    23: TUserAuxSendHPF3,
    12: TGameAuxSendVolume,
    24: TGameAuxSendLPF,
    25: TGameAuxSendHPF,
    4: TBusVolume,
    13: TOutputBusVolume,
    14: TOutputBusHPF,
    15: TOutputBusLPF,
    26: TReflectionBusVolume,
    27: THDRBusThreshold,
    28: THDRBusRatio,
    29: THDRBusReleaseTime,
    30: THDRActiveRange,
    62: THDRBusGameParam,
    63: THDRBusGameParamMin,
    64: THDRBusGameParamMax,
    65: TMidiTrackingRootNote,
    66: TMidiPlayOnNoteType,
    31: TMidiTransposition,
    32: TMidiVelocityOffset,
    67: TMidiKeyRangeMin,
    68: TMidiKeyRangeMax,
    69: TMidiVelocityRangeMin,
    70: TMidiVelocityRangeMax,
    71: TMidiChannelMask,
    72: TMidiTempoSource,
    73: TMidiTargetNode,
    35: TPositioningPanX2D,
    36: TPositioningPanY2D,
    37: TPositioningPanZ2D,
    38: TPositioningPanX3D,
    39: TPositioningPanY3D,
    40: TPositioningPanZ3D,
    41: TPositioningCenterPercent,
    42: TPositioningTypeBlend,
    43: TPositioningEnableAttenuation,
    44: TPositioningConeAttenuationOnOff,
    45: TPositioningConeAttenuation,
    46: TPositioningConeLPF,
    47: TPositioningConeHPF,
    48: TBypassFX,
    49: TBypassAllFX,
    54: TBypassAllMetadata,
    7: TMuteRatio,
    55: TPlayMechanismSpecialTransitionsValue,
    59: TTransitionTime,
    6: TPriority,
    57: TPriorityDistanceOffset,
    53: TMaxNumInstances,
    60: TProbability,
    61: TDialogueMode,
    50: TAvailable0,
    51: TAvailable1,
    52: TAvailable2,
    33: TPlaybackSpeed,
    74: TLoop,
    75: TAttenuationID,
}

const (
    TVolume PropType = 0
    TLFE PropType = 1
    TPitch PropType = 2
    TLPF PropType = 3
    THPF PropType = 4
    TMakeUpGain PropType = 5
    TInitialDelay PropType = 6
    TDelayTime PropType = 7
    TUserAuxSendVolume0 PropType = 8
    TUserAuxSendVolume1 PropType = 9
    TUserAuxSendVolume2 PropType = 10
    TUserAuxSendVolume3 PropType = 11
    TUserAuxSendLPF0 PropType = 12
    TUserAuxSendLPF1 PropType = 13
    TUserAuxSendLPF2 PropType = 14
    TUserAuxSendLPF3 PropType = 15
    TUserAuxSendHPF0 PropType = 16
    TUserAuxSendHPF1 PropType = 17
    TUserAuxSendHPF2 PropType = 18
    TUserAuxSendHPF3 PropType = 19
    TGameAuxSendVolume PropType = 20
    TGameAuxSendLPF PropType = 21
    TGameAuxSendHPF PropType = 22
    TBusVolume PropType = 23
    TOutputBusVolume PropType = 24
    TOutputBusHPF PropType = 25
    TOutputBusLPF PropType = 26
    TReflectionBusVolume PropType = 27
    THDRBusThreshold PropType = 28
    THDRBusRatio PropType = 29
    THDRBusReleaseTime PropType = 30
    THDRActiveRange PropType = 31
    THDRBusGameParam PropType = 32
    THDRBusGameParamMin PropType = 33
    THDRBusGameParamMax PropType = 34
    TMidiTrackingRootNote PropType = 35
    TMidiPlayOnNoteType PropType = 36
    TMidiTransposition PropType = 37
    TMidiVelocityOffset PropType = 38
    TMidiKeyRangeMin PropType = 39
    TMidiKeyRangeMax PropType = 40
    TMidiVelocityRangeMin PropType = 41
    TMidiVelocityRangeMax PropType = 42
    TMidiChannelMask PropType = 43
    TMidiTempoSource PropType = 44
    TMidiTargetNode PropType = 45
    TPANLR PropType = 46
    TPANFR PropType = 47
    TPANUD PropType = 48
    TCenterPCT PropType = 49
    TPositioningPanX2D PropType = 50
    TPositioningPanY2D PropType = 51
    TPositioningPanZ2D PropType = 52
    TPositioningPanX3D PropType = 53
    TPositioningPanY3D PropType = 54
    TPositioningPanZ3D PropType = 55
    TPositioningCenterPercent PropType = 56
    TPositioningTypeBlend PropType = 57
    TPositioningEnableAttenuation PropType = 58
    TPositioningConeAttenuationOnOff PropType = 59
    TPositioningConeAttenuation PropType = 60
    TPositioningConeLPF PropType = 61
    TPositioningConeHPF PropType = 62
    TFeedbackVolumeUnused PropType = 63
    TFeedbackLPFUnused PropType = 64
    TBypassFX PropType = 65
    TBypassAllFX PropType = 66
    TBypassAllMetadata PropType = 67
    TAttachedPluginFXID PropType = 68
    TMuteRatio PropType = 69
    TPlayMechanismSpecialTransitionsValue PropType = 70
    TTransitionTime PropType = 71
    TTrimInTime PropType = 72
    TTrimOutTime PropType = 73
    TFadeInTime PropType = 74
    TFadeOutTime PropType = 75
    TFadeInCurve PropType = 76
    TFadeOutCurve PropType = 77
    TLoopCrossfadeDuration PropType = 78
    TCrossfadeUpCurve PropType = 79
    TCrossfadeDownCurve PropType = 80
    TPriority PropType = 81
    TPriorityDistanceOffset PropType = 82
    TMaxNumInstances PropType = 83
    TProbability PropType = 84
    TDialogueMode PropType = 85
    TAvailable0 PropType = 86
    TAvailable1 PropType = 87
    TAvailable2 PropType = 88
    TPlaybackSpeed PropType = 89
    TLoop PropType = 90
    TLoopStart PropType = 91
    TLoopEnd PropType = 92
    TAttenuationDistanceScaling PropType = 93
    TAttenuationID PropType = 94
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
    TPANUD: "PAN UD",
    TCenterPCT: "Center PCT",
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
    TFeedbackVolumeUnused: "Feedback Volume Unused",
    TFeedbackLPFUnused: "Feedback LPF Unused",
    TBypassFX: "Bypass FX",
    TBypassAllFX: "Bypass All FX",
    TBypassAllMetadata: "Bypass All Metadata",
    TAttachedPluginFXID: "Attached Plugin FX ID",
    TMuteRatio: "Mute Ratio",
    TPlayMechanismSpecialTransitionsValue: "Play Mechanism Special Transitions Value",
    TTransitionTime: "Transition Time",
    TTrimInTime: "Trim In Time",
    TTrimOutTime: "Trim Out Time",
    TFadeInTime: "Fade In Time",
    TFadeOutTime: "Fade Out Time",
    TFadeInCurve: "Fade In Curve",
    TFadeOutCurve: "Fade Out Curve",
    TLoopCrossfadeDuration: "Loop Crossfade Duration",
    TCrossfadeUpCurve: "Crossfade Up Curve",
    TCrossfadeDownCurve: "Crossfade Down Curve",
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
    TLoopStart: "Loop Start",
    TLoopEnd: "Loop End",
    TAttenuationDistanceScaling: "Attenuation Distance Scaling",
    TAttenuationID: "Attenuation ID",
}
