// TODO:
// Vim navigation
// Better styling?
// Search for action?
package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderEventsViewer(open *bool) {
	if !*open {
		return
	}
	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.EventsTag], open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	activeBank, valid := BnkMngr.ActiveBankV()
	if !valid || activeBank.SounBankLock.Load() {
		return
	}
	if activeBank.EventViewer.ActiveEvent != nil {
		size := imgui.NewVec2(0, 0)

		size.Y = imgui.ContentRegionAvail().Y * 0.5
		imgui.BeginChildStrV("EventActionTableChild", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		renderActionsTable(activeBank)
		imgui.EndChild()

		if activeBank.EventViewer.ActiveAction != nil {
			size.X, size.Y = 0, 0
			imgui.BeginChildStrV("EventActionParameter", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
			renderAction(activeBank, activeBank.EventViewer.ActiveAction)
			imgui.EndChild()
		}
	}
}

func renderActionsTable(t *be.BankTab) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	if imgui.BeginTableV("EventActionTable", 3, flags, DefaultSize, 0) {
		imgui.TableSetupColumn("Action ID")
		imgui.TableSetupColumn("Action Type")
		imgui.TableSetupColumn("Action Target")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var action *wwise.Action
		hirc := t.Bank.HIRC()
		flags := imgui.SelectableFlagsSpanAllColumns | 
			     imgui.SelectableFlagsAllowOverlap
		for _, actionID := range t.EventViewer.ActiveEvent.ActionIDs {
			imgui.TableNextRow()

			value, ok := hirc.Actions.Load(actionID)
			if !ok { panic(fmt.Sprintf("No Action object has ID %d.", actionID)) }

			action = value.(wwise.HircObj).(*wwise.Action)
			imgui.TableSetColumnIndex(0)
			if t.EventViewer.ActiveAction == nil {
				t.EventViewer.ActiveAction = action
			}
			selected := actionID == t.EventViewer.ActiveAction.Id
			label := strconv.FormatUint(uint64(actionID), 10)
			if imgui.SelectableBoolPtrV(label, &selected, flags, DefaultSize) {
				t.EventViewer.ActiveAction = action
			}

			if imgui.BeginPopupContextItem() {
				renderActionCtxMenu(t, action, actionID)
				imgui.EndPopup()
			}

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)
			if preview, in := wwise.ActionTypeName[action.Type()]; !in {
				imgui.Text(fmt.Sprintf("Unknown Action Type %d", action.Type()))
			} else {
				imgui.Text(preview)
			}

			imgui.TableSetColumnIndex(2)
			imgui.Text(strconv.FormatUint(uint64(action.IdExt), 10))
		}
		imgui.EndTable()
	}
}

func renderAction(t *be.BankTab, a *wwise.Action) {
	imgui.Text("Type")
	imgui.SameLine()
	imgui.SetNextItemWidth(160)
	if preview, in := wwise.ActionTypeName[a.Type()]; !in {
		imgui.Text(fmt.Sprintf("Unknown Action Type %d", a.Type()))
	} else {
		imgui.Text(preview)
	}

	imgui.Text(fmt.Sprintf("Target: %d", a.IdExt))
	imgui.SameLine()
	if imgui.ArrowButton(fmt.Sprintf("GoToTarget%d", a.IdExt), imgui.DirRight) {
		switch a.ActionParam.(type) {
		case *wwise.ActionActiveParam:
			t.SetActiveActorMixerHirc(a.IdExt)
			t.OpenActorMixerHircNode(a.IdExt)
			// This is only handle at the end of rendering loop
			// Thus, only Bank Explorer is focused
			imgui.SetWindowFocusStr("Actor Mixer Hierarchy")
			imgui.SetWindowFocusStr("Bank Explorer")
			t.Focus = be.BankTabActorMixer
		case *wwise.ActionPlayParam:
			t.SetActiveActorMixerHirc(a.IdExt)
			t.OpenActorMixerHircNode(a.IdExt)
			imgui.SetWindowFocusStr("Actor Mixer Hierarchy")
			imgui.SetWindowFocusStr("Bank Explorer")
			t.Focus = be.BankTabActorMixer
		}
	}

	renderAllProp(&a.PropBundle, &a.RangePropBundle, t.Version())
	renderActionParam(a)
	imgui.SameLine()
}

func renderActionCtxMenu(t *be.BankTab, action *wwise.Action, actionID uint32) {
	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(actionID), 10)))
		}
	})
	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy Action Target ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(action.IdExt), 10)))
		}
	})
}

func renderEventCtxMenu(t *be.BankTab, event *wwise.Event) {
	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(event.Id), 10)))
		}
	})
	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy Action IDs") {

		}
	})
}

func renderEventsTable(t *be.BankTab) {
	imgui.SeparatorText("Filter")
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar("By event ID", imgui.DataTypeU32, uintptr(utils.Ptr(&t.EventViewer.Filter.Id))) {
		t.FilterEvents()
	}
	imgui.SeparatorText("")

	if imgui.BeginTableV("EventsTable", 1, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("Event ID")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(t.EventViewer.Filter.Events)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				event := t.EventViewer.Filter.Events[n]
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)

				if t.EventViewer.ActiveEvent == nil {
					t.EventViewer.ActiveEvent = event
				}
				selected := event.Id == t.EventViewer.ActiveEvent.Id
				label := strconv.FormatUint(uint64(event.Id), 10)
				if imgui.SelectableBoolPtr(label, &selected) {
					t.EventViewer.ActiveEvent = event
					t.EventViewer.ActiveAction = nil
				}

				if imgui.BeginPopupContextItem() {
					renderEventCtxMenu(t, event)
					imgui.EndPopup()
				}
			}
		}
		imgui.EndTable()
	}
}

func renderActionParam(a *wwise.Action) {
	if imgui.TreeNodeStr("Action Parameter") {
		switch t := a.ActionParam.(type) {
		case *wwise.ActionNoParam:
			imgui.Text("This type of action does not have any action parameter.")
		case *wwise.ActionActiveParam:
			t.EnumFadeCurve = renderFadeInCurveCombo(t.EnumFadeCurve)
			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionPlayParam:
			t.EnumFadeCurve = renderFadeInCurveCombo(t.EnumFadeCurve)
			imgui.Text(fmt.Sprintf("Bank ID: %d", t.BankID))
		case *wwise.ActionSetValueParam:
			t.EnumFadeCurve = renderFadeInCurveCombo(t.EnumFadeCurve)
			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionSetStateParam:
			imgui.Text(fmt.Sprintf("State Group ID: %d", t.StateGroupID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToStateGroup%d", t.StateGroupID), imgui.DirRight)
			imgui.EndDisabled()

			imgui.Text(fmt.Sprintf("Target State ID: %d", t.TargetStateID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToTargetState%d", t.TargetStateID), imgui.DirRight)
			imgui.EndDisabled()
		case *wwise.ActionSetSwitchParam:
			imgui.Text(fmt.Sprintf("Switch Group ID: %d", t.SwitchGroupID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToSwitchGroup%d", t.SwitchGroupID), imgui.DirRight)
			imgui.EndDisabled()

			imgui.Text(fmt.Sprintf("Switch State ID: %d", t.SwitchStateID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToSwitchState%d", t.SwitchStateID), imgui.DirRight)
			imgui.EndDisabled()
		case *wwise.ActionSetRTPCParam:
			imgui.Text(fmt.Sprintf("RTPC ID: %d", t.RTPCID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToRTPC%d", t.RTPCID), imgui.DirRight)
			imgui.EndDisabled()
			
			imgui.BeginDisabled()
			imgui.SetNextItemWidth(96)
			imgui.InputFloat("Set RTPC Value To", &t.RTPCValue)
			imgui.EndDisabled()
		case *wwise.ActionSetFXParam:
			audioDeviceEle := t.AudioDeviceElement()
			imgui.BeginDisabled()
			imgui.Checkbox("Audio Device Element", &audioDeviceEle)
			imgui.EndDisabled()
			
			imgui.BeginDisabled()
			imgui.SetNextItemWidth(48)
			imgui.SliderScalar("Effect Slot Index", imgui.DataTypeU8, uintptr(utils.Ptr(&t.SlotIndex)), 0, 254)
			imgui.EndDisabled()

			imgui.Text(fmt.Sprintf("Target FX ID: %d", t.FXID))
			imgui.SameLine()
			imgui.BeginDisabled()
			imgui.ArrowButton(fmt.Sprintf("GoToTargetFX%d", t.FXID), imgui.DirRight)
			imgui.EndDisabled()

			shared := t.Shared()
			imgui.BeginDisabled()
			imgui.Checkbox("Audio Device Element", &shared)
			imgui.EndDisabled()

			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionByPassFXParam:
			bypass := t.ByPass()
			imgui.BeginDisabled()
			imgui.Checkbox("By Pass", &bypass)
			imgui.EndDisabled()
			
			imgui.BeginDisabled()
			imgui.SetNextItemWidth(48)
			imgui.SliderScalar("Bypass Effect Slot Index", imgui.DataTypeU8, uintptr(utils.Ptr(&t.ByFxSolt)), 0, 254)
			imgui.EndDisabled()

			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionSeekParam:
			relative := t.SeekRelativeDuration()
			if imgui.Checkbox("Relative Seek", &relative) {
				t.SetSeekRelativeDuration(relative)
			}
			imgui.PushItemWidth(96)
			if relative {
				imgui.SliderFloat("Seek (Percent)", &t.SeekValue, 0, 100)
			} else {
				imgui.SliderFloat("Seek (Time)", &t.SeekValue, 0, 10)
			}
			if relative {
				imgui.SliderFloat("Seek Offset Min (Percent)", &t.SeekValueMin, -100, 0)
				imgui.SliderFloat("Seek Offset Max (Percent)", &t.SeekValueMax, 0, 100)
			} else {
				imgui.SliderFloat("Seek Offset Min (Time)", &t.SeekValueMin, -10, 0)
				imgui.SliderFloat("Seek Offset Max (Time)", &t.SeekValueMax, 0, 10)
			}
			imgui.PopItemWidth()
			snap := t.IsSnapToNearestMark()
			if imgui.Checkbox("Seek To Nearest Marker", &snap) {
				t.SetSnapToNearestMark(snap)
			}
			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionReleaseParam:
			imgui.Text("This type of action does not have any action parameter.")
		case *wwise.ActionPlayEventParam:
			imgui.Text("This type of action does not have any action parameter.")
		}
		imgui.TreePop()
	}
}

func renderFadeInCurveCombo(EnumFadeCurve wwise.InterpCurveType) wwise.InterpCurveType {
	curveType := int32(EnumFadeCurve)
	imgui.Text("Fade-In Curve")
	imgui.SameLine()
	imgui.SetNextItemWidth(256)
	if imgui.ComboStrarr("##Fade-In Curve", &curveType, wwise.InterpCurveTypeName, int32(wwise.InterpCurveTypeCount)) {
		EnumFadeCurve = wwise.InterpCurveType(curveType)
	}
	return EnumFadeCurve
}

func renderActionExceptParamTable(params []wwise.ExceptParam) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	DefaultSize.Y = 128
	if imgui.BeginTableV("ActionExceptParamTable", 2, flags, DefaultSize, 0) {
		imgui.TableSetupColumn("Exception Target ID")
		imgui.TableSetupColumn("Is Bus")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for _, param := range params {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.Text(strconv.FormatUint(uint64(param.ID), 10))
			imgui.SameLine()
			imgui.BeginDisabled()
			if imgui.ArrowButton(fmt.Sprintf("GoToExceptionTarget%d", param.ID), imgui.DirRight) {

			}
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatBool(param.IsBus == 1))
		}

		imgui.EndTable()
	}
	DefaultSize.Y = 0
}
