package ui

import (
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBaseParam(t *bankTab, o wwise.HircObj) {
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
		renderProp(b.PropBundle)
		renderRangeProp(b.RangePropBundle)
		renderAdvSetting(b)
		imgui.TreePop()
	}
}

func renderChangeParentQuery(t *bankTab, b *wwise.BaseParameter, hid uint32, disable bool) {
	size := imgui.NewVec2(imgui.ContentRegionAvail().X * 0.40, 256)
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

func bindChangeRoot(t *bankTab, hid, np, op uint32) func() {
	return func() {
		t.changeRoot(hid, np, op)
	}
}

func bindRemoveRoot(t *bankTab, hid, op uint32) func() {
	return func() {
		t.removeRoot(hid, op)
	}
}

func renderChangeParentListing(t *bankTab) {
	size := imgui.NewVec2(0, 256)
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
		priorityOverrideParent := o.PriorityOverrideParent()
		if imgui.Checkbox("Priority Override Parent", &priorityOverrideParent) {
			o.SetPriorityOverrideParent(priorityOverrideParent)
		}

		priorityApplyDistFactor := o.PriorityApplyDistFactor()
		if imgui.Checkbox("Priority Apply Dist Factor", &priorityApplyDistFactor) {
			o.SetPriorityApplyDistFactor(priorityApplyDistFactor)
		}

		overrideMidiEventsBehavior := o.OverrideMidiEventsBehavior()
		if imgui.Checkbox("Override Midi Events Behavior", &overrideMidiEventsBehavior) {
			o.SetOverrideMidiEventsBehavior(overrideMidiEventsBehavior)
		}

		overrideMidiNoteTracking := o.OverrideMidiNoteTracking()
		if imgui.Checkbox("Override Midi Note Tracking", &overrideMidiNoteTracking) {
			o.SetOverrideMidiNoteTracking(overrideMidiNoteTracking)
		}

		enableMidiNoteTracking := o.EnableMidiNoteTracking()
		if imgui.Checkbox("Enable Midi Note Tracking", &enableMidiNoteTracking) {
			o.SetEnableMidiNoteTracking(enableMidiNoteTracking)
		}

		midiBreakLoopOnNoteOff := o.MidiBreakLoopOnNoteOff()
		if imgui.Checkbox("MIDI Break Loop On Note Off", &midiBreakLoopOnNoteOff) {
			o.SetMidiBreakLoopOnNoteOff(midiBreakLoopOnNoteOff)
		}
		imgui.TreePop()
	}
}

func renderAdvSetting(o *wwise.BaseParameter) {
	if imgui.TreeNodeExStr("Advanced Setting") {
		killNewest := o.AdvanceSetting.KillNewest()
		if imgui.Checkbox("Kill Newest", &killNewest) {
			o.AdvanceSetting.SetKillNewest(killNewest)
		}

		useVirtualBehavior := o.AdvanceSetting.UseVirtualBehavior()
		if imgui.Checkbox("Use virtual behavior", &useVirtualBehavior) {
			o.AdvanceSetting.SetUseVirtualBehavior(useVirtualBehavior)
		}

		ignoreParentMaxNumInst := o.AdvanceSetting.IgnoreParentMaxNumInst()
		if imgui.Checkbox("Ignore parent max number instance", &ignoreParentMaxNumInst) {
			o.AdvanceSetting.SetIgnoreParentMaxNumInst(ignoreParentMaxNumInst)
		}

		isVVoicesOptOverrideParent := o.AdvanceSetting.IsVVoicesOptOverrideParent()
		if imgui.Checkbox("Is Virtual Voices Opt Override Parent", &isVVoicesOptOverrideParent) {
			o.AdvanceSetting.SetVVoicesOptOverrideParent(isVVoicesOptOverrideParent)
		}

		maxNumInstance := int32(o.AdvanceSetting.MaxNumInstance)
		imgui.SetNextItemWidth(128.0)
		if imgui.InputInt("Max number of instance", &maxNumInstance) {
			if maxNumInstance >= 0 && maxNumInstance <= 0xFFFF {
				o.AdvanceSetting.MaxNumInstance = uint16(maxNumInstance)
			}
		}

		imgui.TreePop()
	}
}
