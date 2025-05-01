package ui

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/AllenDang/cimgui-go/utils"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func showObjectEditor(activeTab *bankTab) {
	imgui.Begin("Object Editor")
	if activeTab == nil {
		imgui.End()
		return
	}

	if activeTab.writeLock.Load() {
		imgui.End()
		return
	}

	if !imgui.BeginTabBarV("ObjectEditorTabBar",
		imgui.TabBarFlagsReorderable |
		imgui.TabBarFlagsAutoSelectNewTabs |
		imgui.TabBarFlagsTabListPopupButton | imgui.TabBarFlagsFittingPolicyScroll,
	) {
		imgui.End()
		return
	}

	s := []wwise.HircObj{}
	for i, h := range activeTab.filtered {
		if activeTab.lSelStorage.Contains(imgui.ID(i)) {
			s = append(s, h)
		}
	}

	for i, h := range s {
		var label string
		switch h.(type) {
		case *wwise.Unknown:
			label = fmt.Sprintf("Unknown Object %d", i)
		default:
			id, _ := h.HircID()
			label = fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
		}
		if imgui.BeginTabItem(label) {
			switch h.(type) {
			case *wwise.ActorMixer:
				showActorMixerProperties(h.(*wwise.ActorMixer))
			case *wwise.RanSeqCntr:
				showRanSeqCntrProperties(h.(*wwise.RanSeqCntr))
			case *wwise.Sound:
				showSoundProperties(h.(*wwise.Sound))
			case *wwise.Unknown:
				showUnknowProperties(h.(*wwise.Unknown))
			}
			imgui.EndTabItem()
		}
	}
	imgui.EndTabBar()
	imgui.End()
}

func showActorMixerProperties(o *wwise.ActorMixer) {
	showBaseParameter(o.BaseParam)
}

func showRanSeqCntrProperties(o *wwise.RanSeqCntr) {
	showBaseParameter(o.BaseParam)
}

func showSoundProperties(o *wwise.Sound) {
	showBaseParameter(o.BaseParam)
}

func showUnknowProperties(o *wwise.Unknown) {
	imgui.Text(
		fmt.Sprintf(
			"Support for hierarchy object type %s is still under construction.",
			wwise.HircTypeName[o.HircType()],
		),
	)
}

func showBaseParameter(o *wwise.BaseParameter) {
	if !imgui.TreeNodeExStr("Base Parameter") {
		return
	}

	// imgui.CurrentStyle().SetIndentSpacing(0.0)

	// Bus ID override and Direct ID override should be done through a separate 
	// UI since there are many ID. Filter and combo are required
	showByBitVector(o)
	showProperty(o.PropBundle)
	showRangeProperty(o.RangePropBundle)
	showAdvanceSetting(o)
	// showRTPC(o.RTPC)

	imgui.TreePop()
}

func showByBitVector(o *wwise.BaseParameter) {
	if !imgui.TreeNodeExStr("Override (Category 1)") {
		return
	}
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

func showProperty(p *wwise.PropBundle) {
	if !imgui.TreeNodeExStr("Property") {
		return
	}

	if imgui.Button("Add Property") {
		if _, err := p.New(); err != nil {
			slog.Info("Failed to add new property", "error", err)
		}
	}

	if !imgui.BeginTableV(
		"PropertyTable", 4,
		imgui.TableFlagsResizable | imgui.TableFlagsRowBg | 
		imgui.TableFlagsBordersH | imgui.TableFlagsBordersV,
		imgui.Vec2{X: 0.0, Y: 0.0}, 0,
	) {
		return
	}

	imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
	imgui.TableSetupColumn("Property ID")
	imgui.TableSetupColumn("Property Value (decimal view)")
	imgui.TableSetupColumn("Property Value (integer view)")
	imgui.TableHeadersRow()

	for i := range p.PropValues {
		v := p.PropValues[i]
		currP := v.P
		currV := slices.Clone(v.V)

		imgui.TableNextRow()
		imgui.TableSetColumnIndex(0)
		imgui.SetNextItemWidth(0)
		imgui.PushIDStr(fmt.Sprintf("DeleteProperty_%d", i))
		if imgui.Button("X") {
			p.Remove(currP)
			imgui.PopID()
			break
		}
		imgui.PopID()

		imgui.TableSetColumnIndex(1)
		imgui.SetNextItemWidth(-1)
		stageP := int32(currP)
		imgui.PushIDStr(fmt.Sprintf("PropertySelection_%d", i))
		if imgui.ComboStrarr(
			"", &stageP, 
			wwise.PropLabel_140, int32(len(wwise.PropLabel_140)),
		) {
			if _, found := p.HasPid(uint8(stageP)); !found {
				v.P = uint8(stageP)
				p.Sort()
				imgui.PopID()
				break
			}
		}
		imgui.PopID()

		var stageVF float32
		var stageVI int32
		_, err := binary.Decode(currV, wio.ByteOrder, &stageVF)
		if err != nil { panic(err) }
		_, err = binary.Decode(currV, wio.ByteOrder, &stageVI)
		if err != nil { panic(err) }

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
	imgui.TreePop()
}

func showRangeProperty(r *wwise.RangePropBundle) {
	if !imgui.TreeNodeExStr("Range Property") {
		return
	}

	if imgui.Button("Add Range Property") {
		if _, err := r.New(); err != nil {
			slog.Info("Failed to add new range property", "error", err)
		}
	}

	if !imgui.BeginTableV(
		"RangePropertyTable", 6,
		imgui.TableFlagsResizable | imgui.TableFlagsRowBg | 
		imgui.TableFlagsBordersH | imgui.TableFlagsBordersV,
		imgui.Vec2{X: 0.0, Y: 0.0}, 0,
	) {
		return
	}

	imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
	imgui.TableSetupColumn("Property ID")
	imgui.TableSetupColumn("Min (decimal view)")
	imgui.TableSetupColumn("Min (integer view)")
	imgui.TableSetupColumn("Max (decimal view)")
	imgui.TableSetupColumn("Max (integer view)")
	imgui.TableHeadersRow()

	for i := range r.RangeValues {
		v := r.RangeValues[i]
		currP := v.PId
		currMin := v.Min
		currMax := v.Max

		imgui.TableNextRow()

		imgui.TableSetColumnIndex(0)
		imgui.PushIDStr(fmt.Sprintf("DeleteRangeProperty_%d", i))
		if imgui.Button("X") {
			r.Remove(currP)
			imgui.PopID()
			break
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
				v.PId = uint8(stageP)
				r.Sort()
				imgui.PopID()
				break
			}
		}
		imgui.PopID()

		var stageMinF float32
		var stageMaxF float32
		var stageMinI int32
		var stageMaxI int32
		_, err := binary.Decode(currMin, wio.ByteOrder, &stageMinF)
		if err != nil { panic(err) }
		_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxF)
		if err != nil { panic(err) }
		_, err = binary.Decode(currMin, wio.ByteOrder, &stageMinI)
		if err != nil { panic(err) }
		_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxI)
		if err != nil { panic(err) }

		imgui.TableSetColumnIndex(2)
		imgui.SetNextItemWidth(-1)
		imgui.PushIDStr(fmt.Sprintf("RangePropertyMinF32_%d", i))
		if imgui.InputFloat("", &stageMinF) {
			r.UpdatePropF32(currP, stageMinF, stageMaxF)
		}
		imgui.PopID()

		imgui.TableSetColumnIndex(3)
		imgui.SetNextItemWidth(-1)
		imgui.PushIDStr(fmt.Sprintf("RangePropertyMinI32_%d", i))
		if imgui.InputInt("", &stageMinI) {
			r.UpdatePropI32(currP, stageMinI, stageMaxI)
		}
		imgui.PopID()

		imgui.TableSetColumnIndex(4)
		imgui.SetNextItemWidth(-1)
		imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxF32_%d", i))
		if imgui.InputFloat("", &stageMaxF) {
			r.UpdatePropF32(currP, stageMinF, stageMaxF)
		}
		imgui.PopID()

		imgui.TableSetColumnIndex(5)
		imgui.SetNextItemWidth(-1)
		imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxI32_%d", i))
		if imgui.InputInt("", &stageMaxI) {
			r.UpdatePropI32(currP, stageMinI, stageMaxI)
		}
		imgui.PopID()
	}
	imgui.EndTable()
	imgui.TreePop()
}

func showAdvanceSetting(o *wwise.BaseParameter) {
	if !imgui.TreeNodeStr("Advanced Setting") {
		return
	}
	
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

func showRTPC(rtpc *wwise.RTPC) {
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

