package ui

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderProp(p *wwise.PropBundle) {
	if imgui.TreeNodeExStr("Property") {
		if imgui.Button("Add Property") {
			if _, err := p.New(); err != nil {
				slog.Info("Failed to add new property", "error", err)
			}
		}
		renderPropTable(p)
		imgui.TreePop()
	}
}

func renderPropTable(p *wwise.PropBundle) {
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PropTable", 4, flags, outerSize, 0) {
		var deleteProp func() = nil
		var changeProp func() = nil

		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property ID")
		imgui.TableSetupColumn("Property Value (decimal view)")
		imgui.TableSetupColumn("Property Value (integer view)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for i := range p.PropValues {
			v := &p.PropValues[i]
			currP := v.P
			currV := slices.Clone(v.V) // Performance disaster overtime?

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)
			
			imgui.PushIDStr(fmt.Sprintf("DeleteProperty_%d", i))
			if imgui.Button("X") {
				deleteProp = bindDeleteProp(p, currP)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			stageP := int32(currP)
			imgui.PushIDStr(fmt.Sprintf("PropertySelection_%d", i))

			if imgui.ComboStrarr("", &stageP, wwise.PropLabel_140, int32(len(wwise.PropLabel_140))) {
				if _, found := p.HasPid(uint8(stageP)); !found {
					changeProp = bindChangeProp(p, v, stageP)
				}
			}

			imgui.PopID()

			var stageVF float32
			var stageVI int32

			_, err := binary.Decode(currV, wio.ByteOrder, &stageVF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currV, wio.ByteOrder, &stageVI)
			if err != nil {
				panic(err)
			}


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

		if deleteProp != nil {
			deleteProp()
		}
		if changeProp != nil {
			changeProp()
		}
	}
}

func bindChangeProp(p *wwise.PropBundle, v *wwise.PropValue, pid int32) func() {
	return func() {
		v.P = uint8(pid)
		p.Sort()
	}
}

func bindDeleteProp(p *wwise.PropBundle, pid uint8) func() {
	return func() {
		p.Remove(pid)
	}
}

func renderRangeProp(r *wwise.RangePropBundle) {
	if imgui.TreeNodeExStr("Range Property") {
		if imgui.Button("Add Range Property") {
			if _, err := r.New(); err != nil {
				slog.Info("Failed to add new range property", "error", err)
			}
		}
		renderRangePropTable(r)
		imgui.TreePop()
	}
}

func renderRangePropTable(r *wwise.RangePropBundle) {
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("RangePropTable", 6, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property ID")
		imgui.TableSetupColumn("Min (decimal view)")
		imgui.TableSetupColumn("Min (integer view)")
		imgui.TableSetupColumn("Max (decimal view)")
		imgui.TableSetupColumn("Max (integer view)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var deleteProp func() = nil
		var changeProp func() = nil

		for i := range r.RangeValues {
			v := &r.RangeValues[i]
			currP := v.PId
			currMin := v.Min
			currMax := v.Max

			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("DelRangProp%d", i))

			if imgui.Button("X") {
				deleteProp = bindRemoveRangeProp(r, currP)
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
					changeProp = bindChangeRangeProp(r, v, uint8(stageP))
				}
			}

			imgui.PopID()

			var stageMinF float32
			var stageMaxF float32
			var stageMinI int32
			var stageMaxI int32

			_, err := binary.Decode(currMin, wio.ByteOrder, &stageMinF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxF)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMin, wio.ByteOrder, &stageMinI)
			if err != nil {
				panic(err)
			}
			_, err = binary.Decode(currMax, wio.ByteOrder, &stageMaxI)
			if err != nil {
				panic(err)
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMinF32_%d", i))

			// Delay?
			if imgui.InputFloat("", &stageMinF) {
				r.UpdatePropF32(currP, stageMinF, stageMaxF)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMinI32_%d", i))

			// Delay?
			if imgui.InputInt("", &stageMinI) {
				r.UpdatePropI32(currP, stageMinI, stageMaxI)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(4)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxF32_%d", i))

			// Delay?
			if imgui.InputFloat("", &stageMaxF) {
				r.UpdatePropF32(currP, stageMinF, stageMaxF)
			}

			imgui.PopID()

			imgui.TableSetColumnIndex(5)
			imgui.SetNextItemWidth(-1)

			imgui.PushIDStr(fmt.Sprintf("RangePropertyMaxI32_%d", i))

			// Delay?
			if imgui.InputInt("", &stageMaxI) {
				r.UpdatePropI32(currP, stageMinI, stageMaxI)
			}

			imgui.PopID()
		}
		imgui.EndTable()

		if deleteProp != nil {
			deleteProp()
		}
		if changeProp != nil {
			changeProp()
		}
	}
}

func bindRemoveRangeProp(r *wwise.RangePropBundle, p uint8) func() {
	return func() {
		r.Remove(p)
	}
}

func bindChangeRangeProp(
	r *wwise.RangePropBundle, v *wwise.RangeValue, p uint8,
) func() {
	return func() {
		v.PId = p
		r.Sort()
	}
}
