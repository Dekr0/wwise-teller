package ui

import (
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderModulatorTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	filterState := &t.ModulatorViewer.Filter

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
	) {
		t.FilterModulator()
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(150)
	preview := wwise.HircTypeName[filterState.Type]
	if imgui.BeginCombo("By Type", preview) {
		var filter func() = nil
		for _, _type := range wwise.ModulatorTypes {
			selected := filterState.Type == _type
			preview = wwise.HircTypeName[_type]
			if imgui.SelectableBoolPtr(preview, &selected) {
				filterState.Type = _type
				filter = t.FilterModulator
			}
		}
		imgui.EndCombo()
		if filter != nil {
			filter()
		}
	}
	imgui.SeparatorText("")

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	if imgui.BeginTableV("ModulatorsTable", 2, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Modulators)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				Modulator := filterState.Modulators[n]
				
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				if t.ModulatorViewer.ActiveModulator == nil {
					t.ModulatorViewer.ActiveModulator = Modulator
				}
				idA, err := Modulator.HircID()
				if err != nil { panic(err) }
				idB, err := t.ModulatorViewer.ActiveModulator.HircID()
				if err != nil { panic(err) }
				selected := idA == idB
				if imgui.SelectableBoolPtrV(
					strconv.FormatUint(uint64(idA), 10),
					&selected,
					DefaultSelectableFlags,
					DefaultSize,
				) {
					t.ModulatorViewer.ActiveModulator = Modulator
				}

				imgui.TableSetColumnIndex(1)
				imgui.Text(wwise.HircTypeName[Modulator.HircType()])
			}
		}
		imgui.EndTable()
	}
}
