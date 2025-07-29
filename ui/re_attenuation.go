package ui

import (
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderAttenuationTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	filterState := &t.AttenuationViewer.Filter

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
	) {
		t.FilterAttenuations()
	}

	imgui.SeparatorText("")

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	if imgui.BeginTableV("AttenuationTable", 1, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("Attenuation ID")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Attenuations)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				attenuation := filterState.Attenuations[n]
				
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				if t.AttenuationViewer.ActiveAttenuation == nil {
					t.AttenuationViewer.ActiveAttenuation = attenuation
				}
				idA := attenuation.Id
				idB := t.AttenuationViewer.ActiveAttenuation.Id
				selected := idA == idB
				if imgui.SelectableBoolPtrV(
					strconv.FormatUint(uint64(idA), 10),
					&selected,
					DefaultSelectableFlags,
					DefaultSize,
				) {
					t.AttenuationViewer.ActiveAttenuation = attenuation
				}
			}
		}
		imgui.EndTable()
	}
}

func renderAttenuationTableCtx(t *be.BankTab, o wwise.HircObj, id uint32) {
	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(id), 10)))
		}
	})
}
