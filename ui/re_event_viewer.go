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
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderEventsViewer(t *BankTab) {
	imgui.Begin("Events")
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		imgui.End()
		return
	}

	if t.EventViewer.ActiveEvent != nil {
		size := imgui.NewVec2(0, 0)

		size.Y = imgui.ContentRegionAvail().Y * 0.5
		imgui.BeginChildStrV("EventActionTableChild", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

		size.X, size.Y = 0, 0
		const flags = DefaultTableFlags | imgui.TableFlagsScrollY
		if imgui.BeginTableV("EventActionTable", 3, flags, size, 0) {
			imgui.TableSetupColumn("Action ID")
			imgui.TableSetupColumn("Action Type")
			imgui.TableSetupColumn("Action Target")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()

			var action *wwise.Action
			hirc := t.Bank.HIRC()
			flags := imgui.SelectableFlagsSpanAllColumns | 
				     imgui.SelectableFlagsAllowOverlap
			for i, actionID := range t.EventViewer.ActiveEvent.ActionIDs {
				imgui.TableNextRow()

				value, ok := hirc.Actions.Load(actionID)
				if !ok {
					imgui.TableSetColumnIndex(0)
					imgui.Text(fmt.Sprintf("Action %d (Null)", actionID))
				} else {
					action = value.(wwise.HircObj).(*wwise.Action)
					imgui.TableSetColumnIndex(0)
					if t.EventViewer.ActiveAction == nil {
						t.EventViewer.ActiveAction = action
					}
					selected := actionID == t.EventViewer.ActiveAction.Id
					label := strconv.FormatUint(uint64(actionID), 10)
					if imgui.SelectableBoolPtrV(label, &selected, flags, size) {
						t.EventViewer.ActiveAction = action
					}

					imgui.TableSetColumnIndex(1)
					imgui.SetNextItemWidth(-1)
					t := int32(action.Type())
					label = fmt.Sprintf("##ActionType%d", i)
					imgui.ComboStrarr(label, &t, wwise.ActionTypeName[1:], wwise.ActionTypeCount)

					imgui.TableSetColumnIndex(2)
					imgui.Text(strconv.FormatUint(uint64(action.IdExt), 10))
				}
			}
			imgui.EndTable()
		}
		imgui.EndChild()

		if t.EventViewer.ActiveAction != nil {
			size.X, size.Y = 0, 0
			imgui.BeginChildStrV("EventActionParameter", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

			actionType := int32(t.EventViewer.ActiveAction.Type())
			imgui.Text("Type")
			imgui.SameLine()
			imgui.SetNextItemWidth(160)
			imgui.ComboStrarr("##ActionType", &actionType, wwise.ActionTypeName[1:], wwise.ActionTypeCount)

			preview := strconv.FormatUint(uint64(t.EventViewer.ActiveAction.IdExt), 10)
			imgui.Text("Target")
			imgui.SameLine()
			imgui.SetNextItemWidth(96)
			if imgui.BeginCombo("##ActionTarget", preview) {
				imgui.EndCombo()
			}
			imgui.SameLine()
			id := "GoTo" + preview
			if imgui.ArrowButton(id, imgui.DirRight) {
				t.LinearStorage.Clear()
				t.LinearStorage.SetItemSelected(imgui.ID(t.EventViewer.ActiveAction.IdExt), true)
			}

			// renderActionProp(t.ActiveAction)
			// renderActionRangeProp(t.ActiveAction)
			// renderActionParam(t.ActiveAction)

			imgui.SameLine()

			imgui.EndChild()
		}
	}
	imgui.End()
}

func renderEventsTable(t *BankTab) {
	imgui.SeparatorText("Filter")
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar("By event ID", imgui.DataTypeU32, uintptr(utils.Ptr(&t.EventViewer.EventFilter.Id))) {
		t.FilterEvents()
	}
	imgui.SeparatorText("")

	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	size := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("EventsTable", 1, flags, size, 0) {
		imgui.TableSetupColumn("Event ID")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(t.EventViewer.EventFilter.Events)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				event := t.EventViewer.EventFilter.Events[n]
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)

				size.X, size.Y = 0, 0
				if t.EventViewer.ActiveEvent == nil {
					t.EventViewer.ActiveEvent = event
				}
				selected := event.Id == t.EventViewer.ActiveEvent.Id
				label := strconv.FormatUint(uint64(event.Id), 10)
				if imgui.SelectableBoolPtr(label, &selected) {
					t.EventViewer.ActiveEvent = event
					t.EventViewer.ActiveAction = nil
				}
			}
		}
		imgui.EndTable()
	}
}

func renderActionProp(a *wwise.Action) {
	if imgui.TreeNodeStr("Property") {
		imgui.TreePop()
	}
}

func renderActionRangeProp(a *wwise.Action) {
	if imgui.TreeNodeStr("Randomizer Property") {
		imgui.TreePop()
	}
}

func renderActionParam(a *wwise.Action) {
	if imgui.TreeNodeStr("Action Parameter") {
		switch t := a.ActionParam.(type) {
		case *wwise.ActionNoParam:
		case *wwise.ActionActiveParam:
			t.EnumFadeCurve = renderFadeInCurveCombo(t.EnumFadeCurve)
			renderActionExceptParamTable(t.ExceptParams)
		case *wwise.ActionPlayParam:
			t.EnumFadeCurve = renderFadeInCurveCombo(t.EnumFadeCurve)
		}
		imgui.TreePop()
	}
}

func renderFadeInCurveCombo(EnumFadeCurve uint8) uint8 {
	curveType := int32(EnumFadeCurve)
	imgui.Text("Fade-In Curve")
	imgui.SetNextItemWidth(128)
	if imgui.ComboStrarr("##Fade-In Curve", &curveType, wwise.InterpCurveTypeName, int32(wwise.InterpCurveTypeCount)) {
		EnumFadeCurve = uint8(curveType)
	}
	return EnumFadeCurve
}

func renderActionExceptParamTable(params []wwise.ExceptParam) {
	flags := DefaultTableFlags
	size := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("ActionExceptParamTable", 2, flags, size, 0) {
		imgui.TableSetupColumn("Exception Target ID")
		imgui.TableSetupColumn("Is Bus")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for _, param := range params {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.Text(strconv.FormatUint(uint64(param.ID), 10))

			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatBool(param.IsBus == 1))
		}

		imgui.EndTable()
	}
}
