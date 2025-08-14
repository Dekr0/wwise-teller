package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/implot"
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

	imgui.SeparatorText("Cone Attenuation")
	coneEnabled := a.ConeEnabled()
	if imgui.Checkbox("Cone Used", &coneEnabled) {
		a.SetConeEnabled(coneEnabled)
	}
	imgui.BeginDisabledV(!coneEnabled)
	imgui.PushItemWidth(160.0)
	imgui.SliderFloat("Cone Max Attenuation", &a.OutsideVolume, -360.0, 360.0)
	imgui.SliderFloat("Cone Inner Angle Degree", &a.InsideDegrees, -360.0, 360.0)
	imgui.SliderFloat("Cone Outer Angle Degree", &a.OutsideDegrees, -360.0, 360.0)
	imgui.SliderFloat("Cone LPF", &a.LoPass, -8192.0, 8192.0)
	imgui.SliderFloat("Cone HPF", &a.HiPass, -8192.0, 8192.0)
	imgui.PopItemWidth()
	imgui.EndDisabled()

	imgui.SeparatorText("Attenuation Settings")
	RenderAttenuationSettingsGE141(a)

	imgui.SeparatorText("Attenuation Conversion Table")
	if imgui.BeginTabBarV("AttenuationConversionTableBar", DefaultTabFlags) {
		const plotFlags = implot.FlagsNoLegend | 
			              implot.FlagsNoTitle  | 
					      implot.FlagsNoFrame  |
			              implot.FlagsCanvasOnly
		const axisFlags = implot.AxisFlagsNoLabel | 
			              implot.AxisFlagsAutoFit

		for i, t := range a.AttenuationConversionTables {
			if imgui.BeginTabItemV(fmt.Sprintf("Table %d", i), nil, 0) {
				const flags = DefaultTableFlags | imgui.TableFlagsScrollY
				size := imgui.NewVec2(-1, 256)
				if imgui.BeginTableV(fmt.Sprintf("Table%dPoints", i), 5, flags, size, 0) {
					imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
					imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
					imgui.TableSetupColumnV("X", imgui.TableColumnFlagsWidthFixed, 0, 0)
					imgui.TableSetupColumnV("Y", imgui.TableColumnFlagsWidthFixed, 0, 0)
					imgui.TableSetupColumn("Curve Interpolation")
					imgui.TableSetupScrollFreeze(0, 1)
					imgui.TableHeadersRow()

					numPoints := len(t.RTPCGraphPointsX)
					for i := range t.RTPCGraphPointsX {
						disable := i == 0 || i >= numPoints - 1
						disableAdd := i >= numPoints - 1
						x := &t.RTPCGraphPointsX[i]
						y := &t.RTPCGraphPointsY[i]
						interp := &t.RTPCGraphPointsInterp[i]

						imgui.TableNextRow()

						imgui.TableSetColumnIndex(0)
						imgui.SetNextItemWidth(40)
						imgui.BeginDisabledV(disable)
						imgui.PushIDStr(fmt.Sprintf("##RmPoint%d", i))
						if imgui.Button("X") {}
						imgui.PopID()
						imgui.EndDisabled()

						imgui.TableSetColumnIndex(1)
						imgui.SetNextItemWidth(40)
						imgui.BeginDisabledV(disableAdd)
						imgui.PushIDStr(fmt.Sprintf("##AppendPoint%d", i))
						if imgui.Button("+") {}
						imgui.PopID()
						imgui.EndDisabled()

						imgui.TableSetColumnIndex(2)
						imgui.SetNextItemWidth(128)
						if disable {
							imgui.BeginDisabledV(disable)
							imgui.SliderFloat(
								fmt.Sprintf("##FromSlider%d", i),
								x,
								-8192.0,
								8192,
							)
							imgui.EndDisabled()
						} else {
							imgui.SliderFloat(
								fmt.Sprintf("##FromSlider%d", i),
								x,
								t.RTPCGraphPointsX[i - 1] + 1e-4,
								t.RTPCGraphPointsX[i + 1] - 1e-4,
							)
						}

						imgui.TableSetColumnIndex(3)
						imgui.SetNextItemWidth(128)
						imgui.BeginDisabledV(!ModifiyEverything)
						imgui.SliderFloat(
							fmt.Sprintf("##ToSlider%d", i),
							y, -8192.0, 8192.0,
						)
						imgui.EndDisabled()

						imgui.TableSetColumnIndex(4)
						imgui.SetNextItemWidth(-1)
						curve := int32(*interp)
						imgui.BeginDisabledV(disable)
						if imgui.ComboStrarr(
							fmt.Sprintf("##RTPCCurve%d", i),
							&curve, wwise.InterpCurveTypeName, int32(wwise.InterpCurveTypeCount),
						) {
							*interp = uint32(curve)
						}
						imgui.EndDisabled()
					}
					imgui.EndTable()
				}

				vec4 := imgui.NewVec4(0, 0, 0, -1)
				if implot.BeginPlotV(
					fmt.Sprintf("AttenuationConversionTable%dPlot", i),
					size,
					plotFlags,
				) {
					implot.SetupAxesV("", "", axisFlags, axisFlags)
					implot.SetNextMarkerStyleV(implot.MarkerCircle, -1, vec4, -1, vec4)
					implot.PlotLineFloatPtrFloatPtr(
						fmt.Sprintf("AttenuationConversionTableLinear%d", i),
						utils.SliceToPtr(t.RTPCGraphPointsX),
						utils.SliceToPtr(t.RTPCGraphPointsY),
						int32(len(t.RTPCGraphPointsX)),
					)
					implot.EndPlot()
				}
				imgui.EndTabItem()
			}
		}
		imgui.EndTabBar()
	}

	renderRTPC(a.Id, &a.RTPC, "Attenuation RTPC Initial")
}

func RenderAttenuationSettingsGE141(a *wwise.Attenuation)  {
	if imgui.BeginTableV("Attenuation Settings Table", 3, DefaultTableFlags, imgui.NewVec2(-1, 320), 0) {
		imgui.TableSetupColumn("Driver")
		imgui.TableSetupColumn("Property")
		imgui.TableSetupColumn("Curve")
		imgui.TableHeadersRow()

		i := 0
		for _, name := range wwise.AttenuationDistancePropertyG141 {
			imgui.TableNextRow()
			imgui.PushIDStr("Distance" + name)
			imgui.TableSetColumnIndex(0)
			imgui.Text("Distance")
			imgui.TableSetColumnIndex(1)
			imgui.Text(name)
			imgui.TableSetColumnIndex(2)
			imgui.Text(fmt.Sprintf("Table %d", a.Curves[i]))
			imgui.PopID()
			i += 1
		}

		for _, name := range wwise.AttenuationObstructionPropertyG141 {
			imgui.TableNextRow()
			imgui.PushIDStr("Obstruction" + name)
			imgui.TableSetColumnIndex(0)
			imgui.Text("Obstruction")
			imgui.TableSetColumnIndex(1)
			imgui.Text(name)
			imgui.TableSetColumnIndex(2)
			imgui.Text(fmt.Sprintf("Table %d", a.Curves[i]))
			imgui.PopID()
			i += 1
		}

		for _, name := range wwise.AttenuationObstructionPropertyG141 {
			imgui.TableNextRow()
			imgui.PushIDStr("Occulsion" + name)
			imgui.TableSetColumnIndex(0)
			imgui.Text("Occulsion")
			imgui.TableSetColumnIndex(1)
			imgui.Text(name)
			imgui.TableSetColumnIndex(2)
			imgui.Text(fmt.Sprintf("Table %d", a.Curves[i]))
			imgui.PopID()
			i += 1
		}

		for _, name := range wwise.AttenuationObstructionPropertyG141 {
			imgui.TableNextRow()
			imgui.PushIDStr("Diffraction" + name)
			imgui.TableSetColumnIndex(0)
			imgui.Text("Diffraction")
			imgui.TableSetColumnIndex(1)
			imgui.Text(name)
			imgui.TableSetColumnIndex(2)
			imgui.Text(fmt.Sprintf("Table %d", a.Curves[i]))
			imgui.PopID()
			i += 1
		}

		for _, name := range wwise.AttenuationObstructionPropertyG141 {
			imgui.TableNextRow()
			imgui.PushIDStr("Transmission" + name)
			imgui.TableSetColumnIndex(0)
			imgui.Text("Transmission")
			imgui.TableSetColumnIndex(1)
			imgui.Text(name)
			imgui.TableSetColumnIndex(2)
			imgui.Text(fmt.Sprintf("Table %d", a.Curves[i]))
			imgui.PopID()
			i += 1
		}

		imgui.EndTable()
	}
}
