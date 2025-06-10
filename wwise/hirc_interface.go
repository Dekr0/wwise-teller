package wwise

type HircType uint8

const (
	HircTypeAll                 HircType = 0x00
	HircTypeState               HircType = 0x01 // ???
	HircTypeSound               HircType = 0x02
	HircTypeAction              HircType = 0x03 // ???
	HircTypeEvent               HircType = 0x04 // ???
	HircTypeRanSeqCntr          HircType = 0x05
	HircTypeSwitchCntr          HircType = 0x06
	HircTypeActorMixer          HircType = 0x07
	HircTypeBus                 HircType = 0x08 // ???
	HircTypeLayerCntr           HircType = 0x09
	HircTypeMusicSegment        HircType = 0x0a
	HircTypeMusicTrack          HircType = 0x0b // ???
	HircTypeMusicSwitchCntr     HircType = 0x0c // ???
	HircTypeMusicRanSeqCntr     HircType = 0x0d // ???
	HircTypeAttenuation         HircType = 0x0e // ???
	HircTypeDialogueEvent       HircType = 0x0f // ???
	HircTypeFxShareSet          HircType = 0x10 // ???
	HircTypeFxCustom            HircType = 0x11 // ???
	HircTypeAuxBus              HircType = 0x12 // ???
	HircTypeLFOModulator        HircType = 0x13 // ???
	HircTypeEnvelopeModulator   HircType = 0x14 // ???
	HircTypeAudioDevice         HircType = 0x15 // ???
	HircTypeTimeModulator       HircType = 0x16 // ???
	HircTypeCount                    = 0x17
)

var KnownHircTypes []HircType = []HircType{
	0x00,
	HircTypeState,
	HircTypeSound,
	HircTypeAction,
	HircTypeEvent,
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeBus,
	HircTypeLayerCntr,
	HircTypeMusicSegment,
	HircTypeMusicTrack,
	HircTypeMusicSwitchCntr,
	HircTypeMusicRanSeqCntr,
	HircTypeAttenuation,
	HircTypeFxShareSet,
	HircTypeFxCustom,
	HircTypeAuxBus,
}

var ContainerHircTypes []HircType = []HircType{
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
	HircTypeMusicSegment,
	HircTypeMusicSwitchCntr,
	HircTypeMusicRanSeqCntr,
}

func ContainerHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeRanSeqCntr      ||
		   t == HircTypeSwitchCntr      ||
		   t == HircTypeActorMixer      ||
		   t == HircTypeLayerCntr       ||
		   t == HircTypeMusicSegment    ||
		   t == HircTypeMusicRanSeqCntr ||
		   t == HircTypeMusicSwitchCntr
}

func NonHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeState       ||
	       t == HircTypeAction      || 
		   t == HircTypeEvent       || 
		   t == HircTypeBus         ||
	       t == HircTypeAttenuation ||
		   t == HircTypeFxShareSet  ||
		   t == HircTypeFxCustom    ||
		   t == HircTypeAuxBus
}

var HircTypeName []string = []string{
	"",
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
	"Aux Bus",
	"LFO Modulator",
	"Envelope Modulator",
	"Audio Device",
	"Time Modulator",
}

type HircObj interface {
	Encode() []byte
	BaseParameter() *BaseParameter
	HircID() (uint32, error)
	HircType() HircType
	IsCntr() bool
	NumLeaf() int
	ParentID() uint32
	// Modify DirectParentId,
	// pre condition: o.DirectParentId == 0
	// post condition: o.DirectParentId == HircObj.HircID
	AddLeaf(o HircObj)
	// Modify DirectParentId,
	// pre condition: o.DirectParentId == HircObj.HircID
	// post condition: DirectParentId = 0
	RemoveLeaf(o HircObj)
}

const SizeOfHircObjHeader = 1 + 4

type HircObjHeader struct {
	Type HircType // U8x
	Size uint32   // U32
}
