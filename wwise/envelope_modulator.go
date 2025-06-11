package wwise

import "github.com/Dekr0/wwise-teller/wio"

type ModulatorPropType uint8

const (
    ModulatorPropTypeScope                ModulatorPropType = 0x0
    ModulatorPropTypeEnvelopeStopPlayback ModulatorPropType = 0x1
    ModulatorPropTypeLFODepth             ModulatorPropType = 0x2
    ModulatorPropTypeLFOAttack            ModulatorPropType = 0x3
    ModulatorPropTypeLFOFrequency         ModulatorPropType = 0x4
    ModulatorPropTypeLFOWaveform          ModulatorPropType = 0x5
    ModulatorPropTypeLFOSmoothing         ModulatorPropType = 0x6
    ModulatorPropTypeLFOPWM               ModulatorPropType = 0x7
    ModulatorPropTypeLFOInitialPhase      ModulatorPropType = 0x8
    ModulatorPropTypeEnvelopeAttackTime   ModulatorPropType = 0x9
    ModulatorPropTypeEnvelopeAttackCurve  ModulatorPropType = 0xA
    ModulatorPropTypeEnvelopeDecayTime    ModulatorPropType = 0xB
    ModulatorPropTypeEnvelopeSustainLevel ModulatorPropType = 0xC
    ModulatorPropTypeEnvelopeSustainTime  ModulatorPropType = 0xD
    ModulatorPropTypeEnvelopeReleaseTime  ModulatorPropType = 0xE
    ModulatorPropTypeEnvelopeTriggerOn    ModulatorPropType = 0xF
    ModulatorPropTypeTimeDuration         ModulatorPropType = 0x10
    ModulatorPropTypeTimeLoops            ModulatorPropType = 0x11
    ModulatorPropTypeTimePlaybackRate     ModulatorPropType = 0x12
    ModulatorPropTypeTimeInitialDelay     ModulatorPropType = 0x13
	ModulatorPropTypeCount                ModulatorPropType = 0x14
)

var ModulatorPropTypeName []string = []string{
  "Scope",
  "Envelope Stop Playback",
  "LFO Depth",
  "LFO Attack",
  "LFO Frequency",
  "LFO Waveform",
  "LFO Smoothing",
  "LFO PWM",
  "LFO Initial Phase",
  "Envelope Attack Time",
  "Envelope Attack Curve",
  "Envelope Decay Time",
  "Envelope Sustain Level",
  "Envelope Sustain Time",
  "Envelope Release Time",
  "Envelope Trigger On",
  "Time Duration", // 132
  "Time Loops", // 132
  "Time Playback Rate", // 132
  "Time Initial Delay", // 132
}

type ModulatorRTPCParamIDType uint8

const (
	ModulatorRTPCParamIDTypeLfoDepth             ModulatorRTPCParamIDType = 0x0
	ModulatorRTPCParamIDTypeLfoAttack            ModulatorRTPCParamIDType = 0x1
	ModulatorRTPCParamIDTypeLfoFrequency         ModulatorRTPCParamIDType = 0x2
	ModulatorRTPCParamIDTypeLfoWaveform          ModulatorRTPCParamIDType = 0x3
	ModulatorRTPCParamIDTypeLfoSmoothing         ModulatorRTPCParamIDType = 0x4
	ModulatorRTPCParamIDTypeLfoPWM               ModulatorRTPCParamIDType = 0x5
	ModulatorRTPCParamIDTypeLfoInitialPhase      ModulatorRTPCParamIDType = 0x6
	ModulatorRTPCParamIDTypeLfoRetrigger         ModulatorRTPCParamIDType = 0x7
	ModulatorRTPCParamIDTypeEnvelopeAttackTime   ModulatorRTPCParamIDType = 0x8
	ModulatorRTPCParamIDTypeEnvelopeAttackCurve  ModulatorRTPCParamIDType = 0x9
	ModulatorRTPCParamIDTypeEnvelopeDecayTime    ModulatorRTPCParamIDType = 0xA
	ModulatorRTPCParamIDTypeEnvelopeSustainLevel ModulatorRTPCParamIDType = 0xB
	ModulatorRTPCParamIDTypeEnvelopeSustainTime  ModulatorRTPCParamIDType = 0xC
	ModulatorRTPCParamIDTypeEnvelopeReleaseTime  ModulatorRTPCParamIDType = 0xD
	ModulatorRTPCParamIDTypeTimePlaybackSpeed    ModulatorRTPCParamIDType = 0xE
	ModulatorRTPCParamIDTypeTimeInitialDelay     ModulatorRTPCParamIDType = 0xF
)

var ModulatorRTPCParamID []string = []string{
    "Modulator LFO Depth",
    "Modulator LFO Attack",
    "Modulator LFO Frequency",
    "Modulator LFO Waveform",
    "Modulator LFO Smoothing",
    "Modulator LFO PWM",
    "Modulator LFO Initial Phase",
    "Modulator LFO Retrigger",
    "Modulator Envelope AttackTime",
    "Modulator Envelope AttackCurve",
    "Modulator Envelope DecayTime",
    "Modulator Envelope SustainLevel",
    "Modulator Envelope SustainTime",
    "Modulator Envelope Release Time",
    "Modulator Time Playback Speed", // #132~~
    "Modulator Time Initial Delay", // #132~~
}

type EnvelopeModulator struct {
	HircObj

	Id              uint32
	PropBundle      PropBundle
	RangePropBundle RangePropBundle
	RTPC            RTPC
}

func (e *EnvelopeModulator) Encode() []byte {
	dataSize := e.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.Append(HircTypeEnvelopeModulator)
	w.Append(dataSize)
	w.Append(e.Id)
	w.AppendBytes(e.PropBundle.Encode())
	w.AppendBytes(e.RangePropBundle.Encode())
	w.AppendBytes(e.RTPC.Encode())
	return w.BytesAssert(int(dataSize))
}

func (e *EnvelopeModulator) DataSize() uint32 {
	return 4 + e.PropBundle.Size() + e.RangePropBundle.Size() + e.RTPC.Size()
}

func (h *EnvelopeModulator) BaseParameter() *BaseParameter { return nil }

func (h *EnvelopeModulator) HircType() HircType { return HircTypeEnvelopeModulator }

func (h *EnvelopeModulator) HircID() (uint32, error) { return h.Id, nil }

func (h *EnvelopeModulator) IsCntr() bool { return false }

func (h *EnvelopeModulator) NumLeaf() int { return 0 }

func (h *EnvelopeModulator) ParentID() uint32 { return 0 }

func (h *EnvelopeModulator) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *EnvelopeModulator) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (h *EnvelopeModulator) Leafs() []uint32 { return []uint32{} }
