package wwise

type HircType uint8

const (
	HircTypeAll                 HircType = 0x00
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
	HircTypeCount               HircType = 0x17
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
	HircTypeLFOModulator,
	HircTypeEnvelopeModulator,
	HircTypeTimeModulator,
}

var MusicHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeMusicTrack,
	HircTypeMusicSegment,
	HircTypeMusicSwitchCntr,
	HircTypeMusicRanSeqCntr,
}

var ActorMixerHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeSound,
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
	HircTypeDialogueEvent,
}

func ActorMixerHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeSound           ||
	       t == HircTypeRanSeqCntr      ||
		   t == HircTypeSwitchCntr      ||
		   t == HircTypeActorMixer      ||
		   t == HircTypeLayerCntr       ||
		   t == HircTypeDialogueEvent
}

var ContainerActorMixerHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeRanSeqCntr,
	HircTypeSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
}

func ContainerActorMixerHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeRanSeqCntr      ||
		   t == HircTypeSwitchCntr      ||
		   t == HircTypeActorMixer      ||
		   t == HircTypeLayerCntr
}

func MusicHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeMusicTrack      ||
		   t == HircTypeMusicSegment    ||
	       t == HircTypeMusicRanSeqCntr ||
	       t == HircTypeMusicSwitchCntr
}

var ContainerMusicHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeMusicSegment,
	HircTypeMusicSwitchCntr,
	HircTypeMusicRanSeqCntr,
}

func ContainerMusicHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeMusicSegment    ||
		   t == HircTypeMusicRanSeqCntr ||
		   t == HircTypeMusicSwitchCntr
}

var FxHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeFxShareSet,
	HircTypeFxCustom,
}

func FxHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeFxShareSet ||
		   t == HircTypeFxCustom
}

var ModulatorTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeLFOModulator,
	HircTypeEnvelopeModulator,
	HircTypeTimeModulator,
}

func ModulatorType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeLFOModulator      ||
	       t == HircTypeEnvelopeModulator ||
	       t == HircTypeTimeModulator
}

var BusHircTypes []HircType = []HircType{
	HircTypeAll,
	HircTypeBus,
	HircTypeAuxBus,
}

func BusHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeBus ||  t == HircTypeAuxBus
}

func NonHircType(o HircObj) bool {
	t := o.HircType()
	return t == HircTypeState  ||
	       t == HircTypeAction || 
		   t == HircTypeEvent  
}

var HircTypeName []string = []string{
	"All",
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
	Leafs() []uint32
}

const SizeOfHircObjHeader = 1 + 4

type HircObjHeader struct {
	Type HircType // U8x
	Size uint32   // U32
}
