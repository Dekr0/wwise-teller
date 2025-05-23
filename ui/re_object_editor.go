package ui

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/AllenDang/cimgui-go/utils"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderObjEditor(t *bankTab) {
	imgui.Begin("Object Editor")

	if t == nil {
		imgui.End()
		return
	}
	if t.writeLock.Load() {
		imgui.End()
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV(
		"ObjectEditorTabBar",
		imgui.TabBarFlagsReorderable |
		imgui.TabBarFlagsAutoSelectNewTabs |
		imgui.TabBarFlagsTabListPopupButton | imgui.TabBarFlagsFittingPolicyScroll,
	) {
		s := []wwise.HircObj{}
		for _, h := range t.bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if t.storage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}

		for i, h := range s {
			renderHircTab(t, i, h)
		}

		imgui.EndTabBar()
	}

	imgui.End()
}

func renderHircTab(t *bankTab, i int, h wwise.HircObj) {
	var label string
	switch h.(type) {
	case *wwise.Unknown:
		label = fmt.Sprintf("Unknown Object %d", i)
	default:
		id, _ := h.HircID()
		label = fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
	}

	if imgui.BeginTabItem(label) {
		if t.activeHirc != h {
			t.cntrStorage.Clear()
			t.playListStorage.Clear()
		}
		t.activeHirc = h
		switch h.(type) {
		case *wwise.ActorMixer:
			renderActorMixer(t, h.(*wwise.ActorMixer))
		case *wwise.LayerCntr:
			renderLayerCntr(t, h.(*wwise.LayerCntr))
		case *wwise.RanSeqCntr:
			renderRanSeqCntr(t, h.(*wwise.RanSeqCntr))
		case *wwise.SwitchCntr:
			renderSwitchCntr(t, h.(*wwise.SwitchCntr))
		case *wwise.Sound:
			renderSound(t, h.(*wwise.Sound))
		case *wwise.Unknown:
			renderUnknown(h.(*wwise.Unknown))
		}
		imgui.EndTabItem()
	}
}

func renderActorMixer(t *bankTab, o *wwise.ActorMixer) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, o.Container)
}

func renderLayerCntr(t *bankTab, o *wwise.LayerCntr) {
	renderBaseParam(t, o)
}

func renderRanSeqCntr(t *bankTab, o *wwise.RanSeqCntr) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, o.Container)
	renderPlayListSetting(t, o)
}

func renderSwitchCntr(t *bankTab, o *wwise.SwitchCntr) {
	renderBaseParam(t, o)
}

func renderSound(t *bankTab, o *wwise.Sound) {
	renderBaseParam(t, o)
}

func renderUnknown(o *wwise.Unknown) {
	imgui.Text(
		fmt.Sprintf(
			"Support for hierarchy object type %s is still under construction.",
			wwise.HircTypeName[o.HircType()],
		),
	)
}

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

	tableFlags := imgui.TableFlagsResizable   |
			      imgui.TableFlagsReorderable |
		          imgui.TableFlagsRowBg       |
		          imgui.TableFlagsBordersH    |
				  imgui.TableFlagsBordersV    |
		          imgui.TableFlagsScrollY
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("ChangeParentTable", 2, tableFlags, outerSize, 0) {
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

func renderProp(p *wwise.PropBundle) {
	if imgui.TreeNodeExStr("Property") {
		if imgui.Button("Add Property") {
			if _, err := p.New(); err != nil {
				slog.Info("Failed to add new property", "error", err)
			}
		}
		renderPropTable(p)
		imgui.TreePop()
	}
}

func renderPropTable(p *wwise.PropBundle) {
	flags := imgui.TableFlagsResizable | 
			 imgui.TableFlagsRowBg     | 
			 imgui.TableFlagsBordersH  | 
	         imgui.TableFlagsBordersV
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PropTable", 4, flags, outerSize, 0) {
		var deleteProp func() = nil
		var changeProp func() = nil

		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property ID")
		imgui.TableSetupColumn("Property Value (decimal view)")
		imgui.TableSetupColumn("Property Value (integer view)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for i := range p.PropValues {
			v := p.PropValues[i]
			currP := v.P
			currV := slices.Clone(v.V) // Performance disaster overtime?

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)
			
			imgui.PushIDStr(fmt.Sprintf("DeleteProperty_%d", i))
			if imgui.Button("X") {
				deleteProp = bindDeleteProp(p, currP)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			stageP := int32(currP)
			imgui.PushIDStr(fmt.Sprintf("PropertySelection_%d", i))

			if imgui.ComboStrarr("", &stageP, wwise.PropLabel_140, int32(len(wwise.PropLabel_140))) {
				if _, found := p.HasPid(uint8(stageP)); !found {
					changeProp = bindChangeProp(p, v, stageP)
				}
			}

			imgui.PopID()

			var stageVF float32
			var stageVI int32

			_, err := binary.Decode(currV, wio.ByteOrder, &stageVF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currV, wio.ByteOrder, &stageVI)
			if err != nil {
				panic(err)
			}


			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("PropertyValueFloat_%d", i))

			if imgui.InputFloat("", &stageVF) {
				p.UpdatePropF32(currP, stageVF)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("PropertyValueInt_%d", i))

			if imgui.InputInt("", &stageVI) {
				p.UpdatePropI32(currP, stageVI)
			}

			imgui.PopID()
		}

		imgui.EndTable()

		if deleteProp != nil {
			deleteProp()
		}
		if changeProp != nil {
			changeProp()
		}
	}
}

func bindChangeProp(p *wwise.PropBundle, v *wwise.PropValue, pid int32) func() {
	return func() {
		v.P = uint8(pid)
		p.Sort()
	}
}

func bindDeleteProp(p *wwise.PropBundle, pid uint8) func() {
	return func() {
		p.Remove(pid)
	}
}

func renderRangeProp(r *wwise.RangePropBundle) {
	if imgui.TreeNodeExStr("Range Property") {
		if imgui.Button("Add Range Property") {
			if _, err := r.New(); err != nil {
				slog.Info("Failed to add new range property", "error", err)
			}
		}
		renderRangePropTable(r)
		imgui.TreePop()
	}
}

func renderRangePropTable(r *wwise.RangePropBundle) {
	flags := imgui.TableFlagsResizable | 
			 imgui.TableFlagsRowBg     | 
			 imgui.TableFlagsBordersH  | 
	         imgui.TableFlagsBordersV
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("RangePropTable", 6, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property ID")
		imgui.TableSetupColumn("Min (decimal view)")
		imgui.TableSetupColumn("Min (integer view)")
		imgui.TableSetupColumn("Max (decimal view)")
		imgui.TableSetupColumn("Max (integer view)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var deleteProp func() = nil
		var changeProp func() = nil

		for i := range r.RangeValues {
			v := r.RangeValues[i]
			currP := v.PId
			currMin := v.Min
			currMax := v.Max

			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("DelRangProp%d", i))

			if imgui.Button("X") {
				deleteProp = bindRemoveRangeProp(r, currP)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			stageP := int32(currP)

			imgui.PushIDStr(fmt.Sprintf("RangePropertySelection_%d", i))
			
			if imgui.ComboStrarr(
				"", &stageP, 
				wwise.PropLabel_140, int32(len(wwise.PropLabel_140)),
			) {
				if _, found := r.HasPid(uint8(stageP)); !found {
					changeProp = bindChangeRangeProp(r, v, uint8(stageP))
				}
			}

			imgui.PopID()

			var stageMinF float32
			var stageMaxF float32
			var stageMinI int32
			var stageMaxI int32

			_, err := binary.Decode(currMin, wio.ByteOrder, &stageMinF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMin, wio.ByteOrder, &stageMinI)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxI)
			if err != nil {
				panic(err)
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMinF32_%d", i))

			// Delay?
			if imgui.InputFloat("", &stageMinF) {
				r.UpdatePropF32(currP, stageMinF, stageMaxF)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMinI32_%d", i))

			// Delay?
			if imgui.InputInt("", &stageMinI) {
				r.UpdatePropI32(currP, stageMinI, stageMaxI)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(4)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxF32_%d", i))

			// Delay?
			if imgui.InputFloat("", &stageMaxF) {
				r.UpdatePropF32(currP, stageMinF, stageMaxF)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(5)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxI32_%d", i))

			// Delay?
			if imgui.InputInt("", &stageMaxI) {
				r.UpdatePropI32(currP, stageMinI, stageMaxI)
			}

			imgui.PopID()
		}
		imgui.EndTable()

		if deleteProp != nil {
			deleteProp()
		}
		if changeProp != nil {
			changeProp()
		}
	}
}

func bindRemoveRangeProp(r *wwise.RangePropBundle, p uint8) func() {
	return func() {
		r.Remove(p)
	}
}

func bindChangeRangeProp(
	r *wwise.RangePropBundle, v *wwise.RangeValue, p uint8,
) func() {
	return func() {
		v.PId = p
		r.Sort()
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

func renderRTPC(rtpc *wwise.RTPC) {
	imgui.SeparatorText("RTPC")

	if !imgui.TreeNodeExStr("RTPC") {
		return
	}
	for _, r := range rtpc.RTPCItems {
		if !imgui.TreeNodeStr(fmt.Sprintf("RTPC %d (Curve %d)", r.RTPCID, r.RTPCCurveID)) {
			continue
		}
		if implot.BeginPlotV(
			fmt.Sprintf("RTPC %d (Curve %d) Plot", r.RTPCID, r.RTPCCurveID),
			imgui.Vec2{X: -1, Y: 360}, 0,
		) {
			implot.PlotLineFloatPtrInt(
				fmt.Sprintf("RTPC %d (Curve %d) Line", r.RTPCID, r.RTPCCurveID),
				utils.SliceToPtr(r.SamplePoints),
				int32(len(r.SamplePoints)),
			)
			implot.EndPlot()
		}
		imgui.TreePop()
	}
	imgui.TreePop()
}

func renderContainer(t *bankTab, id uint32, cntr *wwise.Container) {
	if imgui.TreeNodeExStr("Container") {
		if imgui.Button("Add New Children") {
		}

		flags := imgui.TableFlagsResizable | 
				 imgui.TableFlagsRowBg     | 
			     imgui.TableFlagsBordersH  | 
				 imgui.TableFlagsBordersV
		outerSize := imgui.NewVec2(0.0, 0.0)
		if imgui.BeginTableV("CntrTable", 2, flags, outerSize, 0) {
			imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumn("Children Hierarchy ID")
			imgui.TableHeadersRow()

			var deleteChild func() = nil
			for _, i := range cntr.Children {
				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)
				imgui.SetNextItemWidth(40)

				imgui.PushIDStr(fmt.Sprintf("DelChild%d", i))
				if imgui.Button("X") {
					deleteChild = bindRemoveRoot(t, i, id)
				}
				imgui.PopID()

				imgui.TableSetColumnIndex(1)
				imgui.SetNextItemWidth(-1)

				imgui.Text(strconv.FormatUint(uint64(i), 10))
			}

			imgui.EndTable()

			if deleteChild != nil {
				deleteChild()
			}
		}
		imgui.TreePop()
	}
}

func renderPlayListSetting(t *bankTab, r *wwise.RanSeqCntr) {
	if imgui.TreeNodeExStr("Play list setting") {
		imgui.PushItemWidth(160)

		renderPlayListValue(r.PlayListSetting)
		renderPlayListMode(r.PlayListSetting)
		renderPlayListMisc(r.PlayListSetting)

		imgui.PopItemWidth()

		renderPlayListTableSet(t, r)

		imgui.TreePop()
	}
}

func renderPlayListTableSet(t *bankTab, r *wwise.RanSeqCntr) {
	outerSize := imgui.NewVec2(0, 0)
	flags := imgui.TableFlagsNone 
	if imgui.BeginTableV("PLTransfer", 3, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthStretch, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthStretch, 0, 0)
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableNextRow()

		imgui.TableSetColumnIndex(0)
		renderPlayListPendingTable(t, r)

		imgui.TableSetColumnIndex(1)
		if imgui.Button(">>") {
			for i := range r.Container.Children {
				r.AddLeafToPlayList(i)
			}
			t.playListStorage.Clear()
		}

		imgui.TableSetColumnIndex(2)
		renderPlayListTable(t, r)

		imgui.EndTable()
	}
}

func renderPlayListValue(p *wwise.PlayListSetting) {
	if imgui.InputScalar("Loop Count", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopCount)),
	) {
	}

	if imgui.InputScalar("Loop Mod Min", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopModMin)),
	) {
	}

	if imgui.InputScalar("Loop Mod Max", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopModMax)),
	) {
	}

	if imgui.InputFloat("Transition Time", &p.TransitionTime) {
	}

	if imgui.InputFloat("Transition Time Mod Min", &p.TransitionTimeModMin) {
	}

	if imgui.InputFloat("Transition Time Mod Max", &p.TransitionTimeModMax) {
	}

	if imgui.InputScalar("Avoid Repeat Count", imgui.DataTypeU16, uintptr(utils.Ptr(&p.AvoidRepeatCount))) {
	}
}

func renderPlayListMode(p *wwise.PlayListSetting) {
	tMode := int32(p.TransitionMode)
	if imgui.ComboStrarr("TransitionMode", &tMode, wwise.TransitionModeString, int32(len(wwise.TransitionModeString))) {
		p.TransitionMode = uint8(tMode)
	}

	rMode := int32(p.RandomMode)
	if imgui.ComboStrarr("Random Mode", &rMode, wwise.RandomModeString, int32(len(wwise.RandomModeString))) {
		p.RandomMode = uint8(rMode)
	}

	mode := int32(p.Mode)
	if imgui.ComboStrarr("Mode", &mode, wwise.PlayListModeString, int32(len(wwise.PlayListModeString))) {
		p.Mode = uint8(mode)
	}
}

func renderPlayListMisc(p *wwise.PlayListSetting) {
	usingWeight := p.UsingWeight()
	if imgui.Checkbox("Using Weight", &usingWeight) {
		p.SetUsingWeight(usingWeight)
	}

	resetPlayListAtEachPlay := p.ResetPlayListAtEachPlay()
	if imgui.Checkbox("Reset Playlist At Each Play", &resetPlayListAtEachPlay) {
		p.SetResetPlayListAtEachPlay(resetPlayListAtEachPlay)
	}

	restartBackward := p.RestartBackward()
	if imgui.Checkbox("Restart Backward", &restartBackward) {
		p.SetRestartBackward(restartBackward)
	}

	continuous := p.Continuous()
	if imgui.Checkbox("Continuous", &continuous) {
		p.SetContinuous(continuous)
	}

	global := p.Global()
	if imgui.Checkbox("Global", &global) {
		p.SetGlobal(global)
	}
}

func renderPlayListPendingTable(t *bankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLPendingCell")
	flags := imgui.TableFlagsResizable | 
			 imgui.TableFlagsRowBg     | 
			 imgui.TableFlagsBordersH  | 
			 imgui.TableFlagsBordersV
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLPendTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Target ID")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var toPlayList func() = nil

		for i, child := range r.Container.Children {
			if slices.ContainsFunc(
				r.PlayListItems, func(p *wwise.PlayListItem) bool {
					return p.UniquePlayID == child
				},
			) {
				continue
			}

			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("ToPlayList%d", i))
			if imgui.Button(">") {
				toPlayList = bindToPlayList(t, r, i)
			}
			imgui.PopID()

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatUint(uint64(child), 10))
		}

		imgui.EndTable()

		if toPlayList != nil {
			toPlayList()
		}
	}
	imgui.EndChild()
}

func bindToPlayList(t *bankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.AddLeafToPlayList(i)
		t.playListStorage.Clear()
	}
}

func renderPlayListTable(t *bankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLCell")
	flags := imgui.TableFlagsResizable | 
			 imgui.TableFlagsRowBg     | 
	         imgui.TableFlagsBordersH  | 
	         imgui.TableFlagsBordersV
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLTable", 4, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Sequence")
		imgui.TableSetupColumn("Target ID")
		imgui.TableSetupColumn("Weight")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var move func() = nil
		var del func() = nil
		var delSel func() = nil

		flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d
		storageSize := t.playListStorage.Size()
		itemCount := int32(len(r.PlayListItems))
		msIO := imgui.BeginMultiSelectV(flags, storageSize, itemCount)
		t.playListStorage.ApplyRequests(msIO)

		for i, p := range r.PlayListItems {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)

			imgui.PushIDStr(fmt.Sprintf("DelPlayListItem%d", i))
			if imgui.Button("X") {
				del = bindPendPlayListItem(t, r, i)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			if c := renderPlayListItemOrderCombo(i, r); c != nil {
				move = c
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			selected := t.playListStorage.Contains(imgui.ID(i))
			label := strconv.FormatUint(uint64(p.UniquePlayID), 10)
			flags := imgui.SelectableFlagsSpanAllColumns | 
					 imgui.SelectableFlagsAllowOverlap
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(i))
			imgui.SelectableBoolPtrV(label, &selected, flags, imgui.NewVec2(0, 0))

			if c := renderPlayListTableCtxMenu(t, r); c != nil {
				delSel = c
			}

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(-1)
			imgui.PushIDStr(fmt.Sprintf("PLWeight%d", i))
			if imgui.InputScalar("", imgui.DataTypeU32, uintptr(utils.Ptr(&p.Weight))) {
			}
			imgui.PopID()
		}

		imgui.EndMultiSelect()
		t.playListStorage.ApplyRequests(msIO)

		if move != nil { 
			move()
		}
		if del != nil {
			del()
		}
		if delSel != nil {
			delSel()
		}

		imgui.EndTable()
	}
	imgui.EndChild()
}

func renderPlayListItemOrderCombo(i int, r *wwise.RanSeqCntr) func() {
	var move func() = nil

	preview := strconv.FormatUint(uint64(i), 10)
	label := fmt.Sprintf("PLSequence%d", i)
	if imgui.BeginCombo(label, preview) {
		for j := range r.PlayListItems {
			selected := i == j
			label := strconv.FormatUint(uint64(j), 10)
			if imgui.SelectableBoolPtr(label, &selected) {
				move = bindChangePlayListItemOrder(i, j, r)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}

	return move
}

func renderPlayListTableCtxMenu(t *bankTab, r *wwise.RanSeqCntr) func() {
	var delSel func() = nil

	if imgui.BeginPopupContextItem() {
		if imgui.Button("Delete") {
			delSel = bindPendSelectPlayListItem(t, r)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	return delSel
}

func bindChangePlayListItemOrder(i, j int, r *wwise.RanSeqCntr) func() {
	return func() {
		r.MovePlayListItem(i, j)
	}
}

func bindPendPlayListItem(t *bankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.RemoveLeafFromPlayList(i)
		t.playListStorage.Clear()
	}
}

func bindPendSelectPlayListItem(t *bankTab, r *wwise.RanSeqCntr) func() {
	return func() {
		mut := false
		tids := []uint32{}
		for i, p := range r.PlayListItems {
			if t.playListStorage.Contains(imgui.ID(i)) {
				tids = append(tids, p.UniquePlayID)
				mut = true
			}
		}
		if mut {
			t.playListStorage.Clear()
		}
		r.RemoveLeafsFromPlayList(tids)
	}
}
