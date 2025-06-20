package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderFxTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	filterState := &t.FxViewer.Filter

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
	) {
		t.FilterFxS()
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(128)
	preview := wwise.HircTypeName[filterState.Type]
	if imgui.BeginCombo("By Type", preview) {
		var filter func() = nil
		for _, _type := range wwise.FxHircTypes {
			selected := filterState.Type == _type
			preview = wwise.HircTypeName[_type]
			if imgui.SelectableBoolPtr(preview, &selected) {
				filterState.Type = _type
				filter = t.FilterFxS
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

	if imgui.BeginTableV("FxSTable", 3, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Type")
		imgui.TableSetupColumn("FX Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Fxs)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				fx := filterState.Fxs[n]
				
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				if t.FxViewer.ActiveFx == nil {
					t.FxViewer.ActiveFx = fx
				}
				idA, err := fx.HircID()
				if err != nil { panic(err) }
				idB, err := t.FxViewer.ActiveFx.HircID()
				if err != nil { panic(err) }
				selected := idA == idB
				if imgui.SelectableBoolPtrV(
					strconv.FormatUint(uint64(idA), 10),
					&selected,
					DefaultSelectableFlags,
					DefaultSize,
				) {
					t.FxViewer.ActiveFx = fx
				}

				imgui.TableSetColumnIndex(1)
				imgui.Text(wwise.HircTypeName[fx.HircType()])

				imgui.TableSetColumnIndex(2)
				switch sfx := fx.(type) {
				case *wwise.FxCustom:
					name, in := wwise.PluginNameLUT[int32(sfx.PluginTypeId)]
					if !in {
						imgui.Text(fmt.Sprintf("Plugin ID %d", sfx.PluginTypeId))
					} else {
						imgui.Text(name)
					}
				case *wwise.FxShareSet:
					name, in := wwise.PluginNameLUT[int32(sfx.PluginTypeId)]
					if !in {
						imgui.Text(fmt.Sprintf("Plugin ID %d", sfx.PluginTypeId))
					} else {
						imgui.Text(name)
					}
				}
			}
		}
		imgui.EndTable()
	}
}
