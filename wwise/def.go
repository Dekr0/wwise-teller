package wwise

const Size8  = 1
const Size16 = 2
const Size32 = 4
const Size64 = 8

const SizeOfHierarchyId = Size32
const SizeOfActionId    = Size32

const SizeOfHierarchyHeader = Size8 + Size32

const SizeOfHierarchyNumCounter = Size32

type LerpType = u8

const (
	LerpLog3      LerpType = 0
	LerpSine      LerpType = 1
	LerpLog1      LerpType = 2
	LerpInvSCurve LerpType = 3
	LerpLinear    LerpType = 4
	LerpSCurve    LerpType = 5
	LerpExp1      LerpType = 6
	LerpnvSine    LerpType = 7
	LerpExp3      LerpType = 8
	LerpConst     LerpType = 9
	LerpCount     LerpType = 10
)

type HircType = uint8

const (
	HircTypeState               HircType = 0x01
	HircTypeSound               HircType = 0x02
	HircTypeAction              HircType = 0x03
	HircTypeEvent               HircType = 0x04
	HircTypeRanSeqCntr          HircType = 0x05
	HircTypeSwitchCntr          HircType = 0x06
	HircTypeActorMixer          HircType = 0x07
	HircTypeBus                 HircType = 0x08
	HircTypeLayerCntr           HircType = 0x09
	HircTypeMusicSegment        HircType = 0x0a
	HircTypeMusicTrack          HircType = 0x0b
	HircTypeMusicSwitchCntr     HircType = 0x0c
	HircTypeMusicRanSeqCntr     HircType = 0x0d
	HircTypeAttenuation         HircType = 0x0e
	HircTypeDialogueEvent       HircType = 0x0f
	HircTypeFxShareSet          HircType = 0x10
	HircTypeFxCustom            HircType = 0x11
	HircTypeAuxBus              HircType = 0x12
	HircTypeLFOModulator        HircType = 0x13
	HircTypeEnvelopeModulator   HircType = 0x14
	HircTypeAudioDevice         HircType = 0x15
	HircTypeTimeModulator       HircType = 0x16
)

var HircTypeName [HircTypeTimeModulator]string = [HircTypeTimeModulator]string{
	"State",
	"Sound",
	"Action",
	"Event",
	"Random / Sequence Container",
	"Switch Container",
	"Actor Mixer",
	"Bus",
	"Layer Container",
	"Music Segment",
	"Music Track",
	"Music Switch Container",
	"Music Random / Sequence Container",
	"Attenuation",
	"Dialogue Event",
	"FX Share Set",
	"FX Custom",
	"Auxiliary Bus",
	"LFO Modulator",
	"Envelope Modulator",
	"Audio Device",
	"Time Modulator",
}

func GetHircTypeName(t HircType) string {
	return HircTypeName[t - 1]
}

type Tag = string
const (
	TagBKHD Tag = "BKHD"
	TagDIDX Tag = "DIDX"
	TagDATA Tag = "DATA"
	TagSTMG Tag = "STMG"
	TagHIRC Tag = "HIRC"
	TagFXPR Tag = "FXPR"
	TagENVS Tag = "ENVS"
	TagSTID Tag = "STID"
	TagINIT Tag = "INIT"
	TagPLAT Tag = "PLAT"
	TagMETA Tag = "META"
)
var KnownTag []string = []string{
	TagBKHD,
	TagDIDX,
	TagDATA,
	TagSTMG,
	TagHIRC,
	TagFXPR,
	TagENVS,
	TagSTID,
	TagINIT,
	TagPLAT,
	TagMETA,
}

type AccumType = u8
