package wwise

import (
	"github.com/Dekr0/wwise-teller/wio"
)

const (
	VirtualQueueBehaviorFromBeginning = 0
	VirtualQueueBehaviorPlayFromElapsedTime = 1
	VirtualQueueBehaviorResume = 2
	VirtualQueueBehaviorCount = 3
)
var VirtualQueueBehaviorString []string = []string{
	"From Beginning", "Play From Elapsed Time", "Resume",
}

const (
	BelowThresholdBehaviorContinueToPlay = 0
	BelowThresholdBehaviorKillVoice = 1
	BelowThresholdBehaviorSendToVirtualVoice = 2
	BelowThresholdBehaviorKillIfOneShotElseVirtual = 3
	BelowThresholdBehaviorCount = 4
)
var BelowThresholdBehaviorString []string = []string{
	"Continue To Play", "Kill Voice", "Send To Virtual Voice", "Kill if finite else virtual",
}

const SizeOfAdvanceSetting = 6
type AdvanceSetting struct {
	AdvanceSettingBitVector uint8 // U8x
	VirtualQueueBehavior uint8 // U8x
	MaxNumInstance uint16 // u16
	BelowThresholdBehavior uint8 // U8x
	HDRBitVector uint8 // U8x
}

func (a *AdvanceSetting) Clone() AdvanceSetting {
	return AdvanceSetting{
		a.AdvanceSettingBitVector,
		a.VirtualQueueBehavior,
		a.MaxNumInstance,
		a.BelowThresholdBehavior,
		a.HDRBitVector,
	}
}

func (a *AdvanceSetting) KillNewest() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 0)
}

func (a *AdvanceSetting) SetKillNewest(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 0, set)
}

func (a *AdvanceSetting) UseVirtualBehavior() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 1)
}

func (a *AdvanceSetting) SetUseVirtualBehavior(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 1, set)
	if !a.UseVirtualBehavior() {
		a.BelowThresholdBehavior = BelowThresholdBehaviorContinueToPlay
		a.VirtualQueueBehavior = VirtualQueueBehaviorPlayFromElapsedTime
	} else {
		if a.BelowThresholdBehavior == BelowThresholdBehaviorContinueToPlay {
			a.BelowThresholdBehavior = BelowThresholdBehaviorSendToVirtualVoice
		}
		if a.VirtualQueueBehavior == VirtualQueueBehaviorPlayFromElapsedTime {
			a.VirtualQueueBehavior = VirtualQueueBehaviorResume
		}
	}
}

func (a *AdvanceSetting) VirtualQueueBehaviorDisable() bool {
	return a.BelowThresholdBehavior == BelowThresholdBehaviorContinueToPlay || 
		   a.BelowThresholdBehavior == BelowThresholdBehaviorKillVoice
}

func (a *AdvanceSetting) IgnoreParentMaxNumInst() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 3)
}

func (a *AdvanceSetting) SetIgnoreParentMaxNumInst(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 3, set)
}

func (a *AdvanceSetting) OverrideParentVVoice() bool {
	return wio.GetBit(a.AdvanceSettingBitVector, 4)
}

func (a *AdvanceSetting) SetVVoicesOptOverrideParent(set bool) {
	a.AdvanceSettingBitVector = wio.SetBit(a.AdvanceSettingBitVector, 4, set)
	if !a.OverrideParentVVoice() {
		if a.IgnoreParentMaxNumInst() && a.UseVirtualBehavior() {
			a.BelowThresholdBehavior = BelowThresholdBehaviorSendToVirtualVoice
			a.VirtualQueueBehavior = VirtualQueueBehaviorResume
		} else {
			a.BelowThresholdBehavior = BelowThresholdBehaviorContinueToPlay
			a.VirtualQueueBehavior = VirtualQueueBehaviorPlayFromElapsedTime
		}
	}
}

func (a *AdvanceSetting) OverrideHDREnvelope() bool {
	return wio.GetBit(a.HDRBitVector, 0)
}

func (a *AdvanceSetting) SetOverrideHDREnvelope(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 0, set)
}

func (a *AdvanceSetting) OverrideAnalysis() bool {
	return wio.GetBit(a.HDRBitVector, 1)
}

func (a *AdvanceSetting) SetOverrideAnalysis(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 1, set)
}

func (a *AdvanceSetting) NormalizeLoudness() bool {
	return wio.GetBit(a.HDRBitVector, 2)
}

func (a *AdvanceSetting) SetNormalizeLoudness(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 2, set)
}

func (a *AdvanceSetting) EnableEnvelope() bool {
	return wio.GetBit(a.HDRBitVector, 3)
}

func (a *AdvanceSetting) SetEnableEnvelope(set bool) {
	a.HDRBitVector = wio.SetBit(a.HDRBitVector, 3, set)
}

