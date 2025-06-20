package wwise

import "github.com/Dekr0/wwise-teller/wio"

const SizeOfPlayListSetting = 24

const (
	TransitionModeDisable = 0
	TransitionModeCrossFadeAmp = 1
	TransitionModeCrossFadePower = 2
	TransitionModeDelay = 3
	TransitionModeSampleAccurate = 4
	TransitionModeTriggerRate = 5
	TransitionModeCount = 6
)
var TransitionModeString []string = []string{
	"Disabled",
	"Cross Fade Amp",
	"Cross Fade Power",
	"Delay",
	"Sample Accurate",
	"Trigger Rate",
}

const (
	RandomModeNormal = 0
	RandomModeShuffle = 1
)
var RandomModeString []string = []string{"Normal", "Shuffle"}

const (
	ModeRandom = 0
	ModeSequence = 1
	ModeCount = 2
)
var PlayListModeString []string = []string{"Random", "Sequence"}

type PlayListSetting struct {
	LoopCount uint16 // u16
	LoopModMin uint16 // u16
	LoopModMax uint16 // u16
	TransitionTime float32 // f32
	TransitionTimeModMin float32 // f32
	TransitionTimeModMax float32 // f32
	AvoidRepeatCount uint16 // u16
	TransitionMode uint8 // U8x
	RandomMode uint8 // U8x
	Mode uint8 // U8x

	// _bIsUsingWeight
	// bResetPlayListAtEachPlay
	// bIsRestartBackward
	// bIsContinuous
	// bIsGlobal
	PlayListBitVector uint8 // U8x
}

func (p *PlayListSetting) Clone() PlayListSetting {
	return PlayListSetting{
		p.LoopCount,
		p.LoopModMin,
		p.LoopModMax,
		p.TransitionTime,
		p.TransitionTimeModMin,
		p.TransitionTimeModMax,
		p.AvoidRepeatCount,
		p.TransitionMode,
		p.RandomMode,
		p.Mode,
		p.PlayListBitVector,
	}
}

func (p *PlayListSetting) Random() bool {
	return p.Mode == ModeRandom
}

func (p *PlayListSetting) UseRandom() {
	p.Mode = ModeRandom
}

func (p *PlayListSetting) Sequence() bool {
	return p.Mode == ModeSequence
}

func (p *PlayListSetting) UseSequence() {
	p.Mode = ModeSequence
}

func (p *PlayListSetting) UseInfiniteLoop() {
	p.LoopCount = 0
	p.LoopModMin = 0
	p.LoopModMax = 0
	p.SetResetPlayListAtEachPlay(true)
}

func (p *PlayListSetting) RandomModeNormal() bool {
	return p.RandomMode == RandomModeNormal
}

func  (p *PlayListSetting) UseRandomModeNormal() {
	p.RandomMode = RandomModeNormal
}

func (p *PlayListSetting) RandomModeShuffle() bool {
	return p.RandomMode == RandomModeShuffle
}

func  (p *PlayListSetting) UseRandomModeShuffle() {
	p.RandomMode = RandomModeShuffle
}

func (p *PlayListSetting) UsingWeight() bool {
	return wio.GetBit(p.PlayListBitVector, 0)
}

func (p *PlayListSetting) SetUsingWeight(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 0, set)
}

func (p *PlayListSetting) ResetPlayListAtEachPlay() bool {
	return wio.GetBit(p.PlayListBitVector, 1)
}

func (p *PlayListSetting) SetResetPlayListAtEachPlay(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 1, set)
}

func (p *PlayListSetting) RestartBackward() bool {
	return wio.GetBit(p.PlayListBitVector, 2)
}

// true is restart backward
// false is normal restart
func (p *PlayListSetting) SetRestartBackward(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 2, set)
}

func (p *PlayListSetting) Continuous() bool {
	return wio.GetBit(p.PlayListBitVector, 3)
}

func (p *PlayListSetting) SetContinuous(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 3, set)
}

func (p *PlayListSetting) Global() bool {
	return wio.GetBit(p.PlayListBitVector, 4)
}

func (p *PlayListSetting) SetGlobal(set bool) {
	p.PlayListBitVector = wio.SetBit(p.PlayListBitVector, 4, set)
}

func (p *PlayListSetting) Infinite() bool {
	return p.LoopCount == 0 && p.LoopModMin == 0  && p.LoopModMax == 0
}

const SizeOfPlayListItem = 8
type PlayListItem struct {
	UniquePlayID uint32 // tid
	Weight int32 // s32
}

