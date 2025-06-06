package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderRTPC(hid uint32, r *wwise.RTPC) {
	if imgui.TreeNodeExStr("RTPC") {
		for i := range r.RTPCItems {
			ri := &r.RTPCItems[i]
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
					wwise.RTPCAccumTypeCount,
				) {
					ri.RTPCAccum = uint8(accumType)
				}
				imgui.EndDisabled()

				imgui.BeginDisabledV(!ModifiyEverything)
				paramID := int32(ri.ParamID)
				imgui.SetNextItemWidth(320)
				if imgui.ComboStrarr(
					"RTPC Parameter",
					&paramID,
					wwise.RTPCParameterIDName,
					wwise.RTPCParameterIDCount,
				) {
					ri.ParamID = uint8(paramID)
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
					wwise.CurveScalingTypeCount,
				) {
					ri.Scaling = uint8(scaling)
				}
				imgui.EndDisabled()

				renderRTPCGraph(ri.RTPCID, ri.RTPCCurveID, ri.ParamID, ri.RTPCGraphPoints)

				imgui.TreePop()
			}
		}
		imgui.TreePop()
	}
}

func renderRTPCGraph(
	rtpcID uint32,
	curveID uint32,
	paramID uint8,
	points []wwise.RTPCGraphPoint,
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
		for i := range points {
			p := &points[i]

			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)
			imgui.BeginDisabledV(i == 0 || i >= len(points) - 1)
			imgui.PushIDStr(fmt.Sprintf("%sRmPoint%d", stackID, i))
			if imgui.Button("X") {}
			imgui.PopID()
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(40)
			imgui.BeginDisabledV(i >= len(points) - 1)
			imgui.PushIDStr(fmt.Sprintf("%sAppendPoint%d", stackID, i))
			if imgui.Button("+") {}
			imgui.PopID()
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(128)
			imgui.BeginDisabledV(i == 0 || i >= len(points) - 1)
			if i == 0 || i >= len(points) - 1 {
				imgui.SliderFloat(
					fmt.Sprintf("##%sFromSlider%d", stackID, i),
					&p.From,
					-8192.0,
					8192,
				)
			} else {
				imgui.SliderFloat(
					fmt.Sprintf("##%sFromSlider%d", stackID, i),
					&p.From,
					points[i - 1].From + 1e-4,
					points[i + 1].From - 1e-4,
				)
			}
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(128)
			imgui.BeginDisabledV(!ModifiyEverything)
			imgui.SliderFloat(
				fmt.Sprintf("##%sToSlider%d", stackID, i),
				&p.To, -8192.0, 8192.0,
			)
			imgui.EndDisabled()

			imgui.TableSetColumnIndex(4)
			imgui.SetNextItemWidth(-1)
			curve := int32(p.Interp)
			if imgui.ComboStrarr(
				fmt.Sprintf("##%sRTPCCurve%d", stackID, i),
				&curve, wwise.InterpCurveTypeName, int32(wwise.InterpCurveTypeConst),
			) {
				p.Interp = uint32(curve)
			}
		}

		imgui.EndTable()
	}
}
