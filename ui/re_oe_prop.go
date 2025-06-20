package ui

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBaseProp(p *wwise.PropBundle) {
	if imgui.TreeNodeExStr("Property") {
		if imgui.Button("Add Base Property") {
			p.AddBaseProp()
		}
		renderBasePropTable(p)
		imgui.TreePop()
	}
}

func renderBasePropTable(p *wwise.PropBundle) {
	const flags = DefaultTableFlags
	if imgui.BeginTableV("PropTable", 4, flags, DefaultSize, 0) {
		var removeBaseProp func() = nil
		var changeBaseProp func() = nil

		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property")
		imgui.TableSetupColumn("Value (Slider)")
		imgui.TableSetupColumn("Value (InputV)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for i := range p.PropValues {
			pid := p.PropValues[i].P
			if !slices.Contains(wwise.BasePropType, pid) {
				continue
			}
			imgui.TableNextRow()
			
			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("RmBaseProp%d", i))
			if imgui.Button("X") {
				removeBaseProp = bindRemoveProp(p, pid)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)
			preview := wwise.PropLabel_140[pid]
			if imgui.BeginCombo(fmt.Sprintf("##ChangeBaseProp%d", i), preview) {
				for _, t := range wwise.BasePropType {
					selected := pid == t

					label := wwise.PropLabel_140[t]
					if imgui.SelectableBoolPtr(label, &selected) {
						changeBaseProp = bindChangeBaseProp(p, i, t)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
				}
				imgui.EndCombo()
			}

			var val float32
			binary.Decode(p.PropValues[i].V, wio.ByteOrder, &val)
			switch pid {
			case wwise.PropTypeVolume:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##VolumeSlider", &val, -96.0, 12.0) {
					p.SetPropByIdxF32(i, val)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##VolumeInputV", &val) {
					if val >= -96.0 && val <= 12.0 {
						p.SetPropByIdxF32(i, val)
					}
				}
			case wwise.PropTypePitch:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				intFloat := int32(val)
				if imgui.SliderInt("##PitchSlider", &intFloat, -2400, 2400) {
					p.SetPropByIdxF32(i, float32(intFloat))
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputInt("##PitchInputV", &intFloat) {
					if intFloat >= -2400 && intFloat <= 2400 {
						p.SetPropByIdxF32(i, float32(intFloat))
					}
				}
			case wwise.PropTypeLPF:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				intFloat := int32(val)
				if imgui.SliderInt("##LPFSlider", &intFloat, 0, 100) {
					p.SetPropByIdxF32(i, float32(intFloat))
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputInt("##LPFInputV", &intFloat) {
					if intFloat >= 0 && intFloat <= 100 {
						p.SetPropByIdxF32(i, float32(intFloat))
					}
				}
			case wwise.PropTypeHPF:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				intFloat := int32(val)
				if imgui.SliderInt("##HPFSlider", &intFloat, 0, 100) {
					p.SetPropByIdxF32(i, float32(intFloat))
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputInt("##HPFInputV", &intFloat) {
					if intFloat >= 0 && intFloat <= 100 {
						p.SetPropByIdxF32(i, float32(intFloat))
					}
				}
			case wwise.PropTypeMakeUpGain:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##MakeUpGainSlider", &val, -96.0, 96.0) {
					p.SetPropByIdxF32(i, val)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##MakeUpGainInputV", &val) {
					if val >= -96.0 && val <= 96.0 {
						p.SetPropByIdxF32(i, val)
					}
				}
			case wwise.PropTypeGameAuxSendVolume:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##GameAuxSendVolumeSlider", &val, -96.0, 12.0) {
					p.SetPropByIdxF32(i, val)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##GameAuxSendVolumeInputV", &val) {
					if val >= -96.0 && val <= 12.0 {
						p.SetPropByIdxF32(i, val)
					}
				}
			case wwise.PropTypeInitialDelay:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##InitialDelaySlider", &val, 0, 60.0) {
					p.SetPropByIdxF32(i, val)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##InitialDelayInputV", &val) {
					if val >= 0 && val <= 60.0 {
						p.SetPropByIdxF32(i, val)
					}
				}
			default:
				panic("Panic Trap")
			}
		}

		imgui.EndTable()

		if removeBaseProp != nil {
			removeBaseProp()
		}
		if changeBaseProp != nil {
			changeBaseProp()
		}
	}
}

func bindChangeBaseProp(p *wwise.PropBundle, idx int, nextPid wwise.PropType) func() {
	return func() { p.ChangeBaseProp(idx, nextPid) }
}

func bindRemoveProp(p *wwise.PropBundle, pid wwise.PropType) func() {
	return func() {
		p.Remove(pid)
	}
}

func renderBaseRangeProp(r *wwise.RangePropBundle) {
	if imgui.TreeNodeExStr("Base Range Property (Randomizer)") {
		if imgui.Button("Add Base Range Property") {
			r.AddBaseProp()
		}
		renderBaseRangePropTable(r)
		imgui.TreePop()
	}
}

func renderBaseRangePropTable(r *wwise.RangePropBundle) {
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("RangePropTableSlider", 6, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property")
		imgui.TableSetupColumn("Min (Slider)")
		imgui.TableSetupColumn("Min (Inputv)")
		imgui.TableSetupColumn("Max (Slider)")
		imgui.TableSetupColumn("Max (Inputv)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var removeBaseRangeProp func() = nil
		var changeBaseRangeProp func() = nil

		for i := range r.RangeValues {
			pid := r.RangeValues[i].P
			if !slices.Contains(wwise.BaseRangePropType, pid) {
				continue
			}
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("RmBaseRangeProp%d", i))
			if imgui.Button("X") {
				removeBaseRangeProp = bindRemoveBaseRangeProp(r, pid)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)
			preview := wwise.PropLabel_140[pid]
			if imgui.BeginCombo(fmt.Sprintf("##ChangeBaseRangeProp%d", i), preview) {
				for _, t := range wwise.BaseRangePropType {
					selected := pid == t
					label := wwise.PropLabel_140[t]
					if imgui.SelectableBoolPtr(label, &selected) {
						changeBaseRangeProp = bindChangeBaseRangeProp(r, i, t)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
				}
				imgui.EndCombo()
			}

			var valMin float32
			var valMax float32
			binary.Decode(r.RangeValues[i].Min, wio.ByteOrder, &valMin)
			binary.Decode(r.RangeValues[i].Max, wio.ByteOrder, &valMax)
			switch pid {
			case wwise.PropTypeVolume:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##VolumeMinSlider", &valMin, -108.0, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##VolumeMinInputV", &valMin) {
					if valMin >= -108.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##VolumeMaxSlider", &valMax, 0, 108.0) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##VolumeMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 108.0 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			case wwise.PropTypePitch:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##PitchMinSlider", &valMin, -4800, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##PitchMinInputV", &valMin) {
					if valMin >= -4800.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##PitchMaxSlider", &valMax, 0, 4800) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##PitchMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 4800 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			case wwise.PropTypeLPF:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##LPFMinSlider", &valMin, -100.0, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##LPFMinInputV", &valMin) {
					if valMin >= -100.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##LPFMaxSlider", &valMax, 0, 100.0) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##LPFMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 100.0 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			case wwise.PropTypeHPF:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##HPFMinSlider", &valMin, -100.0, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##HPFMinInputV", &valMin) {
					if valMin >= -100.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##HPFMaxSlider", &valMax, 0, 100.0) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##HPFMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 100.0 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			case wwise.PropTypeMakeUpGain:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##MakeUpGainMinSlider", &valMin, -192.0, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##MakeUpGainMinInputV", &valMin) {
					if valMin >= -192.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##MakeUpGainMaxSlider", &valMax, 0, 192.0) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##MakeUpGainMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 192.0 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			case wwise.PropTypeInitialDelay:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##InitialDelayMinSlider", &valMin, -60.0, 0) {
					r.SetPropMinByIdxF32(i, valMin)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##InitialDelayMinInputV", &valMin) {
					if valMin >= -60.0 && valMin <= 0 {
						r.SetPropMinByIdxF32(i, valMin)
					}
				}
				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat("##InitialDelayMaxSlider", &valMax, 0, 60.0) {
					r.SetPropMaxByIdxF32(i, valMax)
				}
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat("##InitialDelayMaxInputV", &valMax) {
					if valMax >= 0.0 && valMax <= 60.0 {
						r.SetPropMaxByIdxF32(i, valMax)
					}
				}
			}
		}
		imgui.EndTable()

		if removeBaseRangeProp != nil {
			removeBaseRangeProp()
		}
		if changeBaseRangeProp != nil {
			changeBaseRangeProp()
		}
	}
}

func bindRemoveBaseRangeProp(r *wwise.RangePropBundle, p wwise.PropType) func() {
	return func() {
		r.Remove(p)
	}
}

func bindChangeBaseRangeProp(r *wwise.RangePropBundle, idx int, nextPid wwise.PropType) func() {
	return func() {
		r.ChangeBaseProp(idx, nextPid)
	}
}

func renderAllProp(p *wwise.PropBundle, r *wwise.RangePropBundle) {
	if p != nil && imgui.TreeNodeExStr("All Property (Read-Only)") {
		if imgui.BeginTableV("AllPropReadOnly", 2, DefaultTableFlags, DefaultSize, 0) {
			imgui.TableSetupColumn("Property")
			imgui.TableSetupColumn("Value")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()
			var attenuationID uint32
			var value float32
			for _, p := range p.PropValues {
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				imgui.Text(wwise.PropLabel_140[p.P])
				imgui.TableSetColumnIndex(1)
				if p.P == wwise.PropType(wwise.PropTypeAttenuationID) {
					binary.Decode(p.V, wio.ByteOrder, &attenuationID)
					imgui.Text(strconv.FormatUint(uint64(attenuationID), 10))
				} else {
					binary.Decode(p.V, wio.ByteOrder, &value)
					imgui.Text(strconv.FormatFloat(float64(value), 'f', 4, 32))
				}
			}
			imgui.EndTable()
		}
		imgui.TreePop()
	}
	if r != nil && imgui.TreeNodeExStr("All Range Property (Read-Only)") {
		if imgui.BeginTableV("AllRangePropReadOnly", 3, DefaultTableFlags, DefaultSize, 0) {
			imgui.TableSetupColumn("Property")
			imgui.TableSetupColumn("Min")
			imgui.TableSetupColumn("Max")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()
			var min float32
			var max float32
			for _, r := range r.RangeValues {
				imgui.TableNextRow()
				binary.Decode(r.Min, wio.ByteOrder, &min)
				binary.Decode(r.Max, wio.ByteOrder, &max)
				imgui.TableSetColumnIndex(0)
				imgui.Text(wwise.PropLabel_140[r.P])
				imgui.TableSetColumnIndex(1)
				imgui.Text(strconv.FormatFloat(float64(min), 'f', 4, 32))
				imgui.TableSetColumnIndex(2)
				imgui.Text(strconv.FormatFloat(float64(max), 'f', 4, 32))
			}
			imgui.EndTable()
		}
		imgui.TreePop()
	}
}
