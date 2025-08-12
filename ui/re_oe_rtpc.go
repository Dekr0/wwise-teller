package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/AllenDang/cimgui-go/utils"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderRTPC(hid uint32, r *wwise.RTPC, title string) {
	if imgui.TreeNodeExStr(title) {
		var removeRTPC func() = nil
		for i := range r.RTPCItems {
			ri := &r.RTPCItems[i]
			imgui.PushIDStr(fmt.Sprintf("%dRmRTPC%d", hid, i))
			if imgui.Button("X") {
				removeRTPC = bindRTPCRemove(r, i)
			}
			imgui.PopID()
			imgui.SameLine()
			if imgui.TreeNodeExStrStr(fmt.Sprintf("RTPC%d%d", ri.RTPCID, i), 0, fmt.Sprintf("RTPC %d", ri.RTPCID)) {
				imgui.BeginDisabledV(!ModifiyEverything)
				rtpcType := int32(ri.RTPCType)
				imgui.SetNextItemWidth(128)
				if imgui.ComboStrarr(
					"RTPC Type",
					&rtpcType,
					wwise.RTPCTypeName,
					wwise.RTPCTypeCount,
				) {
					ri.RTPCType = uint8(rtpcType)
				}
				imgui.EndDisabled()

				imgui.BeginDisabledV(!ModifiyEverything)
				accumType := int32(ri.RTPCAccum)
				imgui.SetNextItemWidth(96)
				if imgui.ComboStrarr(
					"Accumulation Type",
					&accumType,
					wwise.RTPCAccumTypeName,
					int32(wwise.RTPCAccumTypeCount),
				) {
					ri.RTPCAccum = wwise.RTPCAccumType(accumType)
				}
				imgui.EndDisabled()

				imgui.BeginDisabledV(!ModifiyEverything)
				paramID := int32(ri.ParamID.Value)
				imgui.SetNextItemWidth(320)
				if imgui.ComboStrarr(
					"RTPC Parameter",
					&paramID,
					wwise.RTPCParameterIDName,
					int32(wwise.RTPCParameterTypeCount),
				) {
				}
				imgui.EndDisabled()

				imgui.Text(fmt.Sprintf("RTPC Curve ID %d", ri.RTPCCurveID))

				imgui.BeginDisabledV(!ModifiyEverything)
				scaling := int32(ri.Scaling)
				imgui.SetNextItemWidth(96)
				if imgui.ComboStrarr(
					"Scaling",
					&scaling,
					wwise.CurveScalingTypeName,
					int32(wwise.CurveScalingTypeCount),
				) {
					ri.Scaling = wwise.CurveScalingType(scaling)
				}
				imgui.EndDisabled()

				renderRTPCGraph(
					ri.RTPCID,
					ri.RTPCCurveID,
					wwise.RTPCParameterType(ri.ParamID.Value),
					ri.RTPCGraphPointsX,
					ri.RTPCGraphPointsY,
					ri.RTPCGraphPointsInterp,
				)

				const plotFlags = implot.FlagsNoLegend | 
					              implot.FlagsNoTitle  | 
							      implot.FlagsNoFrame  |
					              implot.FlagsCanvasOnly
				const axisFlags = implot.AxisFlagsNoLabel | 
					              implot.AxisFlagsAutoFit
				
				if implot.BeginPlotV(fmt.Sprintf("%dRTPCGraph%dPlot", hid, i), imgui.Vec2{X: -1, Y: 128}, plotFlags) {
					implot.SetupAxesV("", "", axisFlags, axisFlags)
					implot.PlotScatterFloatPtrFloatPtr(
						fmt.Sprintf("%dRTPCGraph%dScatter", hid, i),
						utils.SliceToPtr(ri.RTPCGraphPointsX),
						utils.SliceToPtr(ri.RTPCGraphPointsY),
						int32(len(ri.RTPCGraphPointsX)),
					)
					implot.EndPlot()
				}

				imgui.TreePop()
			}
		}
		imgui.TreePop()
		if removeRTPC != nil {
			removeRTPC()
		}
	}
}

func bindRTPCRemove(r *wwise.RTPC, i int) func() {
	return func() { r.RemoveRTPCItem(i) }
}

func renderRTPCGraph(
	rtpcID uint32,
	curveID uint32,
	paramID wwise.RTPCParameterType,
	pointsX []float32,
	pointsY []float32,
	pointsInterp []uint32,
) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	size := imgui.NewVec2(0, 160)
	if imgui.BeginTableV(
		fmt.Sprintf("%d%dPoints", rtpcID, curveID),
		5, flags, size, 0,
	) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("X", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("Y", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Curve Interpolation")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var stackID string = fmt.Sprintf("%d%d", rtpcID, curveID)
		for i := range pointsX {
			x := &pointsX[i]
			y := &pointsY[i]
			interp := &pointsInterp[i]

			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)
			imgui.BeginDisabledV(i == 0 || i >= len(pointsX) - 1)
			imgui.PushIDStr(fmt.Sprintf("%sRmPoint%d", stackID, i))
			if imgui.Button("X") {}
			imgui.PopID()
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(40)
			imgui.BeginDisabledV(i >= len(pointsX) - 1)
			imgui.PushIDStr(fmt.Sprintf("%sAppendPoint%d", stackID, i))
			if imgui.Button("+") {}
			imgui.PopID()
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(128)
			if i == 0 || i >= len(pointsX) - 1 {
				imgui.BeginDisabledV(i == 0 || i >= len(pointsX) - 1)
				imgui.SliderFloat(
					fmt.Sprintf("##%sFromSlider%d", stackID, i),
					x,
					-8192.0,
					8192,
				)
				imgui.EndDisabled()
			} else {
				imgui.SliderFloat(
					fmt.Sprintf("##%sFromSlider%d", stackID, i),
					x,
					pointsX[i - 1] + 1e-4,
					pointsX[i + 1] - 1e-4,
				)
			}

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(128)
			imgui.BeginDisabledV(!ModifiyEverything)
			imgui.SliderFloat(
				fmt.Sprintf("##%sToSlider%d", stackID, i),
				y, -8192.0, 8192.0,
			)
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(4)
			imgui.SetNextItemWidth(-1)
			curve := int32(*interp)
			if imgui.ComboStrarr(
				fmt.Sprintf("##%sRTPCCurve%d", stackID, i),
				&curve, wwise.InterpCurveTypeName, int32(wwise.InterpCurveTypeCount),
			) {
				*interp = uint32(curve)
			}
		}
		imgui.EndTable()
	}
}
