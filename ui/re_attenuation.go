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

func renderAttenuationViewer(open *bool) {
	if !*open {
		return
	}

	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.AttenuationsTag], open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open{
		return
	}

	activeBank, valid := BnkMngr.ActiveBankV()
	if !valid || activeBank.SounBankLock.Load() {
		return
	}
	if activeBank.AttenuationViewer.ActiveAttenuation != nil {
		renderAttenuation(activeBank.AttenuationViewer.ActiveAttenuation)
	}
}

func renderAttenuation(a *wwise.Attenuation) {
	imgui.Text(fmt.Sprintf("Attenuation ID %d", a.Id))

	heightSpreadEnabled := a.HeightSpreadEnabled()
	if imgui.Checkbox("Enable Height Spread", &heightSpreadEnabled) {
		a.SetHeightSpreadEnabled(heightSpreadEnabled)
	}

	if imgui.TreeNodeExStr("Cone Attenuation") {
		imgui.BeginDisabled()
		coneEnabled := a.ConeEnabled()
		if imgui.Checkbox("Cone Used", &coneEnabled) {}
		imgui.PushItemWidth(160.0)
		imgui.SliderFloat("Cone Max Attenuation", &a.OutsideVolume, -360.0, 360.0)
		imgui.SliderFloat("Cone Inner Angle Degree", &a.InsideDegrees, -360.0, 360.0)
		imgui.SliderFloat("Cone Outer Angle Degree", &a.OutsideDegrees, -360.0, 360.0)
		imgui.SliderFloat("Cone LPF", &a.LoPass, -8192.0, 8192.0)
		imgui.SliderFloat("Cone HPF", &a.HiPass, -8192.0, 8192.0)
		imgui.PopItemWidth()
		imgui.EndDisabled()
		imgui.TreePop()
	}
}
