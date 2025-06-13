package ui

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderAuxParam(t *be.BankTab, init *wwise.Bank, o wwise.HircObj) {
	if imgui.TreeNodeExStr("User-Defined Auxiliary Send") {
		b := o.BaseParameter()
		a := &o.BaseParameter().AuxParam
		p := &o.BaseParameter().PropBundle

		parentID := o.ParentID()

		// Hierarchy object without parent bypass the requirement of enabling override 
		// user-defined auxiliary send.
		imgui.BeginDisabledV(parentID == 0)
		overrideAuxSend := a.OverrideAuxSends()
		if imgui.Checkbox("Override User-Defined Auxiliary Send", &overrideAuxSend) {
			b.SetOverrideAuxSends(overrideAuxSend)
		}
		imgui.EndDisabled()

		if parentID != 0 && a.OverrideAuxSends() {
			renderUserAuxSendVolumeTable(p)
		} else if parentID == 0 {
			// Hierarchy object without parent can change auxiliary bus setting freely.
			renderUserAuxSendVolumeTable(p)
		}

		renderAuxBusIDTable(t, init, parentID, a)

		overrideReflectionAuxBus := a.OverrideReflectionAuxBus()
		// Hierarchy object without parent bypass the requirement of enabling
		// override early reflection bus
		imgui.BeginDisabledV(parentID == 0)
		if imgui.Checkbox("Override early reflection bus", &overrideReflectionAuxBus) {
			b.SetOverrideReflectionAuxBus(overrideReflectionAuxBus)
		}
		imgui.EndDisabled()

		imgui.BeginDisabledV(parentID != 0 && !a.OverrideReflectionAuxBus())
		if imgui.Button("Add reflection auxiliary bus volume") {
			p.AddReflectionBusVolume()
		}
		imgui.EndDisabled()

		i, prop := p.ReflectionBusVolume()
		if i != -1 {
			var val float32
			binary.Decode(prop.V, wio.ByteOrder, &val)
			if imgui.SliderFloat("##ReflectionBusVolumeSlider", &val, -96.0, 12.0) {
				p.SetPropByIdxF32(i, val)
			}
			imgui.SameLine()
			if imgui.InputFloat("##ReflectionBusVolumeInputV", &val) {
				if val >= -96.0 && val <= 12.0 {
					p.SetPropByIdxF32(i, val)
				}
			}
		}

		imgui.TreePop()
	}
}

func renderUserAuxSendVolumeTable(p *wwise.PropBundle) {
	if imgui.Button("Add User-Defined Auxiliary Volume") {
		p.AddUserAuxSendVolume()
	}
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("UserAuxSendVolumeTable", 4, flags, outerSize, 0) {
		var removeUserAuxSendVolumeProp func() = nil
		var changeUserAuxSendVolumeProp func() = nil

		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Property")
		imgui.TableSetupColumn("Value (Slider)")
		imgui.TableSetupColumn("Value (InputV)")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		j := -1
		for i := range p.PropValues {
			pid := p.PropValues[i].P
			if !slices.Contains(wwise.UserAuxSendVolumePropType, pid) {
				continue
			}
			imgui.TableNextRow()

			j += 1
			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("RmUserAuxSendVolume%d", i))
			if imgui.Button("X") {
				removeUserAuxSendVolumeProp = bindRemoveProp(p, pid)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)
			preview := wwise.PropLabel_140[pid]
			if imgui.BeginCombo(fmt.Sprintf("##ChangeUserAuxSendVolume%d", i), preview) {
				for _, t := range wwise.UserAuxSendVolumePropType {
					selected := pid == t

					label := wwise.PropLabel_140[t]
					if imgui.SelectableBoolPtr(label, &selected) {
						changeUserAuxSendVolumeProp = bindChangeUserAuxSendVolumeProp(p, i, t)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
				}
				imgui.EndCombo()
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)
			var val float32
			binary.Decode(p.PropValues[i].V, wio.ByteOrder, &val)
			switch pid {
			case wwise.PropTypeUserAuxSendVolume0:
				fallthrough
			case wwise.PropTypeUserAuxSendVolume1:
				fallthrough
			case wwise.PropTypeUserAuxSendVolume2:
				fallthrough
			case wwise.PropTypeUserAuxSendVolume3:
				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				if imgui.SliderFloat(fmt.Sprintf("##UserAuxSendVolume%dSlider", j), &val, -96.0, 12.0) {
					p.SetPropByIdxF32(i, val)
				}
				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				if imgui.InputFloat(fmt.Sprintf("##UserAuxSendVolume%dInputV", j), &val) {
					if val >= -96.0 && val <= 12.0 {
						p.SetPropByIdxF32(i, val)
					}
				}
			default:
				panic("Panic trap")
			}
		}
		imgui.EndTable()
		if removeUserAuxSendVolumeProp != nil {
			removeUserAuxSendVolumeProp()
		}
		if changeUserAuxSendVolumeProp != nil {
			changeUserAuxSendVolumeProp()
		}
	}
}

func bindChangeUserAuxSendVolumeProp(p *wwise.PropBundle, idx int, nextPid wwise.PropType) func() {
	return func() { p.ChangeUserAuxSendVolumeProp(idx, nextPid) }
}

func renderAuxBusIDTable(t *be.BankTab, init *wwise.Bank, parentID uint32, a *wwise.AuxParam) {
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("AuxiliaryBusIDTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Auxiliary Bus")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		for i, aid := range a.AuxIds {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.Text(fmt.Sprintf("%d", i))

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)
			imgui.BeginDisabledV(init == nil || (parentID != 0 && !a.OverrideAuxSends()))
			label := fmt.Sprintf("##AuxiliaryBus%d", i)
			preview := strconv.FormatUint(uint64(aid), 10) 
			if imgui.BeginCombo(label, preview) {
				imgui.EndCombo()
			}
			imgui.EndDisabled()
		}
		imgui.EndTable()
	}
}
