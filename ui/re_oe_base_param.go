package ui

import (
	"encoding/binary"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBaseParam(t *BankTab, o wwise.HircObj) {
	imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyB))
	if imgui.TreeNodeExStr("Base Parameter") {
		hid, err := o.HircID()
		if err != nil {
			panic(err)
		}
		b := o.BaseParameter()
		renderChangeParentQuery(t, b, hid, o.HircType() != wwise.HircTypeSound)
		imgui.SameLine()
		renderChangeParentListing(t)
		renderByBitVec(b)
		renderAuxParam(t, o)
		renderBaseProp(&b.PropBundle)
		renderBaseRangeProp(&b.RangePropBundle)
		renderAdvSetting(b, &b.AdvanceSetting)
		imgui.TreePop()
	}
}

func renderChangeParentQuery(t *BankTab, b *wwise.BaseParameter, hid uint32, disable bool) {
	size := imgui.NewVec2(imgui.ContentRegionAvail().X * 0.40, 160)
	imgui.BeginChildStrV("ChangeParentQuery", size, 0, 0)

	var filter func() = nil
	imgui.Text("Filtered parent by ID")
	if imgui.InputTextWithHint("##Filtered parent by ID", "", &t.parentIdQuery, 0, nil) {
		filter = t.filterParent
	}

	preview := wwise.HircTypeName[t.parentTypeQuery]
	imgui.Text("Filtered by hierarchy type")
	if imgui.BeginComboV("##Filtered by hierarchy type", preview, 0) {
		for _, ht := range wwise.ContainerHircType {
			selected := t.parentTypeQuery == ht
			preview = wwise.HircTypeName[ht]
			if imgui.SelectableBoolPtr(preview, &selected) {
				t.parentTypeQuery = ht
				filter = t.filterParent
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}

	if filter != nil {
		t.filterParent()
	}

	preview = strconv.FormatUint(uint64(b.DirectParentId), 10)
	imgui.BeginDisabledV(disable)
	imgui.Text("Direct Parent ID")
	if imgui.BeginComboV("##Direct Parent ID", preview, 0) {
		var changeParent func() = nil

		for _, p := range t.filteredParent {
			id, err := p.HircID()
			if err != nil {
				continue
			}

			selected := b.DirectParentId == id
			preview := strconv.FormatUint(uint64(id), 10)
			if imgui.SelectableBoolPtr(preview, &selected) {
				changeParent = bindChangeRoot(t, hid, id, b.DirectParentId)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()

		if changeParent != nil {
			changeParent()
		}
	}
	imgui.EndDisabled()
	imgui.EndChild()
}

func bindChangeRoot(t *BankTab, hid, np, op uint32) func() {
	return func() {
		t.changeRoot(hid, np, op)
	}
}

func bindRemoveRoot(t *BankTab, hid, op uint32) func() {
	return func() {
		t.removeRoot(hid, op)
	}
}

func renderChangeParentListing(t *BankTab) {
	size := imgui.NewVec2(0, 160)
	imgui.BeginChildStrV("ChangeParentListing", size, 0, 0)

	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("ChangeParentTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(t.filteredParent)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				o := t.filteredParent[n]

				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)

				id, err := o.HircID()
				if err != nil {
					panic(err)
				}
				imgui.Text(strconv.FormatUint(uint64(id), 10))

				imgui.TableSetColumnIndex(1)
				imgui.Text(wwise.HircTypeName[o.HircType()])
			}
		}
		imgui.EndTable()
	}
	imgui.EndChild()
}

func renderByBitVec(o *wwise.BaseParameter) {
	if imgui.TreeNodeExStr("Override (Category 1)") {
		size := imgui.NewVec2(0, 136)
		imgui.BeginChildStrV("Playback Priority", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

		imgui.BeginDisabledV(o.DirectParentId == 0)
		priorityOverrideParent := o.PriorityOverrideParent()
		if imgui.Checkbox("Priority Override Parent", &priorityOverrideParent) {
			o.SetPriorityOverrideParent(priorityOverrideParent)
		}
		imgui.EndDisabled()

		if o.DirectParentId == 0 {
			if imgui.Button("Add Playback Priorty Property") {
				o.PropBundle.AddPriority()
			}
		}

		if o.PriorityOverrideParent() || o.DirectParentId == 0 {
			i, p := o.PropBundle.Priority()
			if i != -1 {
				var val float32
				binary.Decode(p.V, wio.ByteOrder, &val)
				imgui.SetNextItemWidth(80)
				if imgui.InputFloat("Priority", &val) {
					if val >= 0.0 && val <= 100.0 {
						o.PropBundle.SetPropByIdxF32(i, val)
					}
				}
			}
		}

		imgui.BeginDisabledV(o.DirectParentId != 0 && !o.PriorityOverrideParent())
		priorityApplyDistFactor := o.PriorityApplyDistFactor()
		if imgui.Checkbox("Priority Apply Dist Factor", &priorityApplyDistFactor) {
			o.SetPriorityApplyDistFactor(priorityApplyDistFactor)
		}
		imgui.EndDisabled()

		if o.PriorityApplyDistFactor() {
			i, p := o.PropBundle.PriorityApplyDistFactor()
			if i != -1 {
				var val float32
				binary.Decode(p.V, wio.ByteOrder, &val)
				intFloat := int32(val)
				imgui.SetNextItemWidth(80)
				if imgui.InputInt("Offset priority by", &intFloat) {
					if intFloat >= -100 && intFloat <= 100 {
						o.PropBundle.SetPropByIdxF32(i, float32(intFloat))
					}
				}
				imgui.SameLine()
				imgui.Text("at max distance")
			}
		}
		imgui.EndChild()

		imgui.BeginDisabledV(o.DirectParentId == 0)
		overrideMidiEventsBehavior := o.OverrideMidiEventsBehavior()
		if imgui.Checkbox("Override Midi Events Behavior", &overrideMidiEventsBehavior) {
			o.SetOverrideMidiEventsBehavior(overrideMidiEventsBehavior)
		}
		imgui.EndDisabled()

		imgui.BeginDisabledV(o.DirectParentId == 0)
		overrideMidiNoteTracking := o.OverrideMidiNoteTracking()
		if imgui.Checkbox("Override Midi Note Tracking", &overrideMidiNoteTracking) {
			o.SetOverrideMidiNoteTracking(overrideMidiNoteTracking)
		}
		imgui.EndDisabled()

		imgui.BeginDisabledV(o.DirectParentId == 0)
		enableMidiNoteTracking := o.EnableMidiNoteTracking()
		if imgui.Checkbox("Enable Midi Note Tracking", &enableMidiNoteTracking) {
			o.SetEnableMidiNoteTracking(enableMidiNoteTracking)
		}
		imgui.EndDisabled()

		imgui.BeginDisabledV(o.DirectParentId == 0)
		midiBreakLoopOnNoteOff := o.MidiBreakLoopOnNoteOff()
		if imgui.Checkbox("MIDI Break Loop On Note Off", &midiBreakLoopOnNoteOff) {
			o.SetMidiBreakLoopOnNoteOff(midiBreakLoopOnNoteOff)
		}
		imgui.EndDisabled()
		imgui.TreePop()
	}
}

func renderAdvSetting(b *wwise.BaseParameter, a *wwise.AdvanceSetting) {
	if imgui.TreeNodeExStr("Advance Setting") {
		size := imgui.NewVec2(0, 128)
		imgui.BeginChildStrV("PlaybackLimit", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Playback Limit")

		imgui.BeginDisabledV(b.DirectParentId == 0)
		ignoreParentMaxLimit := a.IgnoreParentMaxNumInst()
		if imgui.Checkbox("Ignore Parent", &ignoreParentMaxLimit) {
			a.SetIgnoreParentMaxNumInst(ignoreParentMaxLimit)
		}
		imgui.EndDisabled()

		// To set behavior of playback limiting, hierarchy object need to 
		// enable "Ignore Parent"
		imgui.BeginDisabledV(b.DirectParentId != 0 && !a.IgnoreParentMaxNumInst())
		imgui.Text("Limit sound instances to:")
		imgui.SameLine()
		maxNumInstance := int32(a.MaxNumInstance)
		imgui.SetNextItemWidth(96)
		// Zero is no limiting
		if imgui.InputInt("##MaxNumInstance", &maxNumInstance) {
			if maxNumInstance >= 0 && maxNumInstance <= 1000 {
				a.MaxNumInstance = uint16(maxNumInstance)
			}
		}

		imgui.Text("When limit is reached:")
		imgui.SameLine()
		var preview string
		if !a.UseVirtualBehavior() {
			preview = "Kill voice"
		} else {
			preview = "Use virtual voice setting"
		}
		imgui.SetNextItemWidth(200)
		if imgui.BeginCombo("##ReachPlaybackLimitBehavior", preview) {
			selected := !a.UseVirtualBehavior()
			if imgui.SelectableBoolPtr("Kill voice", &selected) {
				a.SetUseVirtualBehavior(false)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
			selected = a.UseVirtualBehavior()
			if imgui.SelectableBoolPtr("Use virtual voice setting", &selected) {
				a.SetUseVirtualBehavior(true)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
			imgui.EndCombo()
		}
		imgui.SameLine()
		imgui.Text("for lowest priroty")

		imgui.Text("When priority is equal:")
		imgui.SameLine()
		if !a.KillNewest() {
			preview = "Discard oldest instance"
		} else {
			preview = "Discard newest instance"
		}
		imgui.SetNextItemWidth(192)
		if imgui.BeginCombo("##PriorityEqualBehavior", preview) {
			selected := !a.KillNewest()
			if imgui.SelectableBoolPtr("Discard oldest instance", &selected) {
				a.SetKillNewest(false)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
			selected = a.KillNewest()
			if imgui.SelectableBoolPtr("Discard newest instance", &selected) {
				a.SetKillNewest(true)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
			imgui.EndCombo()
		}
		imgui.EndDisabled()
		imgui.EndChild()

		// Virtual Voice
		size.Y = 112
		imgui.BeginChildStrV("VirtualVovice", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Virtual Voice")
		imgui.BeginDisabledV(b.DirectParentId == 0)
		overrideParentVVoice := a.OverrideParentVVoice()
		if imgui.Checkbox("Override Parent", &overrideParentVVoice) {
			a.SetVVoicesOptOverrideParent(overrideParentVVoice)
		}
		imgui.EndDisabled()
		// Hierarchy object need to enable "Override Parent Virtual Voice" in 
		// order to change virtual voice setting
		imgui.BeginDisabledV(b.DirectParentId != 0 && !a.OverrideParentVVoice())
		belowThreSholdBehavior := int32(a.BelowThresholdBehavior)
		if imgui.ComboStrarr(
			"Virtual Voice Behavior",
			&belowThreSholdBehavior,
			wwise.BelowThresholdBehaviorString,
			wwise.BelowThresholdBehaviorCount,
		) {
			a.BelowThresholdBehavior = uint8(belowThreSholdBehavior)
		}
		imgui.BeginDisabledV(a.VirtualQueueBehaviorDisable())
		virtualQueueBehavior := int32(a.VirtualQueueBehavior)
		if imgui.ComboStrarr(
			"On return to physical voice",
			&virtualQueueBehavior,
			wwise.VirtualQueueBehaviorString,
			wwise.VirtualQueueBehaviorCount,
		) {
			a.VirtualQueueBehavior = uint8(virtualQueueBehavior)
		}
		imgui.EndDisabled()
		imgui.EndDisabled()
		imgui.EndChild()
		// End of Virtual Voice

		// HDR Setting
		size.X = 0
		size.Y = 80
		imgui.BeginChildStrV("HDRSetting", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.BeginDisabledV(b.DirectParentId == 0)
		overrideHDREnvelope := a.OverrideHDREnvelope()
		if imgui.Checkbox("Override Parent", &overrideHDREnvelope) {
			a.SetOverrideHDREnvelope(overrideHDREnvelope)
		}
		imgui.EndDisabled()

		imgui.BeginDisabledV(b.DirectParentId != 0 && !a.OverrideHDREnvelope())
		enabledEnvelope := a.EnableEnvelope()
		if imgui.Checkbox("Enable Envelope", &enabledEnvelope) {
			b.SetEnableEnvelope(enabledEnvelope)
		}
		imgui.BeginDisabledV(!a.EnableEnvelope())
		i, prop := b.PropBundle.HDRActiveRange()
		if i != -1 {
			var val float32
			binary.Decode(prop.V, wio.ByteOrder, &val)
			if imgui.InputFloat("HDR Active Range", &val) {
				if val >= 0 && val <= 24 {
					b.PropBundle.SetPropByIdxF32(i, val)	
				}
			}
		}
		imgui.EndDisabled()
		imgui.EndDisabled()

		imgui.EndChild()
		// End of HDR Setting
		imgui.TreePop()
	}
}
