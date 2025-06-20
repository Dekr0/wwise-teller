package ui

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderAuxParam(m *be.BankManager, init *be.BankTab, o wwise.HircObj) {
	if imgui.TreeNodeExStr("User-Defined Auxiliary Send") {
		b := o.BaseParameter()
		a := &o.BaseParameter().AuxParam
		p := &o.BaseParameter().PropBundle

		parentID := o.ParentID()
		root := parentID == 0

		// Hierarchy object without parent bypass the requirement of enabling 
		// override user-defined auxiliary send.
		imgui.BeginDisabledV(root)
		overrideAuxSend := a.OverrideAuxSends()
		if imgui.Checkbox("Override User-Defined Auxiliary Send", &overrideAuxSend) {
			b.SetOverrideAuxSends(overrideAuxSend)
		}
		imgui.EndDisabled()
		renderUserAuxSendTable(m, init, root, p, a)

		imgui.TreePop()
	}
}

func renderUserAuxSendTable(
	m *be.BankManager,
	init *be.BankTab,
	root bool,
	p *wwise.PropBundle,
	a *wwise.AuxParam,
) {
	DefaultSize.X = 440
	if imgui.BeginTableV("UserAuxSendTable", 7, DefaultTableFlags, DefaultSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 8, 0)
		imgui.TableSetupColumnV("Aux Bus", imgui.TableColumnFlagsWidthFixed, 128, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 20, 0)
		imgui.TableSetupColumnV("Fader", imgui.TableColumnFlagsWidthFixed, 96, 0)
		imgui.TableSetupColumnV("Input", imgui.TableColumnFlagsWidthFixed, 96, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 16, 0)
		imgui.TableSetupColumn("")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		i := wwise.PropTypeUserAuxSendVolume0
		for j, aid := range a.AuxIds {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(-1)
			imgui.Text(fmt.Sprintf("%d", j))

			imgui.TableSetColumnIndex(1)
			{
				imgui.BeginDisabledV(init == nil || (!root && !a.OverrideAuxSends()))

				imgui.PushIDStr(fmt.Sprintf("BusUserAuxSend%d", j))
				if imgui.Button("x") {
					a.AuxIds[j] = 0
				}
				imgui.PopID()
				imgui.SameLine()

				imgui.SetNextItemWidth(96)
				preview := strconv.FormatUint(uint64(aid), 10)
				imgui.Text(preview)

				imgui.SameLine()
				popup := fmt.Sprintf("SetUserAuxSendBus%d", j)
				if imgui.ArrowButton(fmt.Sprintf("SetUserAuxSendBusBtn%d", j), imgui.DirDown) {
					imgui.OpenPopupStr(popup)
				}

				if imgui.BeginPopup(popup) {
					filterState := &init.BusViewer.Filter
					imgui.Text("Aux Bus")
					imgui.SameLine()
					imgui.SetNextItemWidth(96)
					if imgui.BeginCombo(fmt.Sprintf("##UserAuxSendBus%dCombo", j), preview) {
						for _, b := range filterState.Buses {
							if b.HircType() != wwise.HircTypeAuxBus {
								continue
							}
							selected := false

							optionID, err := b.HircID()
							if err != nil { panic(err) }

							if init.BusViewer.ActiveBus != nil {
								activeId, err := init.BusViewer.ActiveBus.HircID()
								if err != nil { panic(err) }
								selected = optionID == activeId
							}

							if imgui.SelectableBoolPtr(strconv.FormatUint(uint64(optionID), 10), &selected) {
								a.AuxIds[j] = optionID
							}
							if selected {
								imgui.SetItemDefaultFocus()
							}
						}
						imgui.EndCombo()
					}

					imgui.Text("Search ")
					imgui.SameLine()
					imgui.SetNextItemWidth(96)
					if imgui.InputScalar(
						"##ID",
						imgui.DataTypeU32,
						uintptr(utils.Ptr(&filterState.Id)),
						) {
						init.FilterBuses()
					}
					imgui.EndPopup()
				}
				imgui.EndDisabled()
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)
			{
				imgui.BeginDisabledV(init == nil || aid == 0)
				if imgui.ArrowButton(fmt.Sprintf("##GoToUserAuxSend%d", i), imgui.DirRight) {
					init.SetActiveBus(aid)
					init.Focus = be.BankTabBuses
					m.SetNextBank = init
					imgui.SetWindowFocusStr("Buses")
				}
				imgui.EndDisabled()
			}

			val := float32(0.0)
			idx, pv := p.Prop(i)
			in := idx != -1
			if in {
				binary.Decode(pv.V, wio.ByteOrder, &val)
			}

			{
				imgui.BeginDisabledV(!root && !a.OverrideAuxSends())

				{
					imgui.BeginDisabledV(!in)
					imgui.TableSetColumnIndex(3)
					imgui.SetNextItemWidth(-1)
					if imgui.SliderFloat(fmt.Sprintf("###UserAuxSendFader%d", j), &val, -96.0, 12.0) {
						p.SetPropByIdxF32(idx, val)
					}
					imgui.TableSetColumnIndex(4)
					imgui.SetNextItemWidth(-1)
					if imgui.InputFloat(fmt.Sprintf("###UserAuxSendInput%d", j), &val) {
						if val >= -96.0 && val <= 12.0 {
							p.SetPropByIdxF32(idx, val)
						}
					}
					imgui.EndDisabled()
				}

				{
					imgui.BeginDisabledV(in)
					imgui.TableSetColumnIndex(5)
					imgui.SetNextItemWidth(16)
					imgui.PushIDStr(fmt.Sprintf("AddUserAuxSend%d", j))
					if imgui.Button("+") {
						p.Add(i)
					}
					imgui.PopID()
					imgui.EndDisabled()
				}

				{
					imgui.BeginDisabledV(!in)
					imgui.TableSetColumnIndex(6)
					imgui.SetNextItemWidth(16)
					imgui.PushIDStr(fmt.Sprintf("RmUserAuxSend%d", j))
					if imgui.Button("X") {
						p.Remove(i)
					}
					imgui.PopID()
					imgui.EndDisabled()
				}

				imgui.EndDisabled()
			}
			i += 1
		}
		imgui.EndTable()
	}
	DefaultSize.X = 0
}

func renderEarlyReflection(m *be.BankManager, init *be.BankTab, o wwise.HircObj) {
	if imgui.TreeNodeExStr("Early Reflections") {
		b := o.BaseParameter()
		a := &o.BaseParameter().AuxParam
		p := &o.BaseParameter().PropBundle

		parentID := o.ParentID()
		root := parentID == 0

		overrideReflectionAuxBus := a.OverrideReflectionAuxBus()
		// Hierarchy object without parent bypass the requirement of enabling
		// override early reflection bus
		imgui.BeginDisabledV(root)
		if imgui.Checkbox("Override early reflection bus", &overrideReflectionAuxBus) {
			b.SetOverrideReflectionAuxBus(overrideReflectionAuxBus)
		}
		imgui.EndDisabled()

		imgui.Text("Auxiliary Bus")
		{
			imgui.BeginDisabledV(!root && !a.OverrideReflectionAuxBus())

			imgui.PushIDStr("ZeroBusReflectionAuxSend")
			if imgui.Button("x") {
				a.ReflectionAuxBus = 0
			}
			imgui.PopID()
			imgui.SameLine()

			imgui.SetNextItemWidth(96)
			preview := strconv.FormatUint(uint64(a.ReflectionAuxBus), 10) 
			imgui.Text(preview)

			imgui.SameLine()
			const popup = "SetReflectionAuxSend"
			imgui.BeginDisabledV(init == nil)
			if imgui.ArrowButton("SetReflectionAuxSendBtn", imgui.DirDown) {
				imgui.OpenPopupStr(popup)
			}
			imgui.EndDisabled()

			if imgui.BeginPopup(popup) {
				filterState := &init.BusViewer.Filter
				imgui.Text("Aux Bus")
				imgui.SameLine()
				imgui.SetNextItemWidth(96)
				if imgui.BeginCombo("##ReflectionAuxSendBusCombo", preview) {
					for _, b := range filterState.Buses {
						if b.HircType() != wwise.HircTypeAuxBus {
							continue
						}
						selected := false

						optionID, err := b.HircID()
						if err != nil { panic(err) }

						if init.BusViewer.ActiveBus != nil {
							activeId, err := init.BusViewer.ActiveBus.HircID()
							if err != nil { panic(err) }
							selected = optionID == activeId
						}
						if imgui.SelectableBoolPtr(strconv.FormatUint(uint64(optionID), 10), &selected) {
							a.ReflectionAuxBus = optionID
						}
						if selected {
							imgui.SetItemDefaultFocus()
						}
					}
					imgui.EndCombo()
				}

				imgui.Text("Search ")
				imgui.SameLine()
				imgui.SetNextItemWidth(96)
				if imgui.InputScalar(
					"##ID",
					imgui.DataTypeU32,
					uintptr(utils.Ptr(&filterState.Id)),
				) {
					init.FilterBuses()
				}
				imgui.EndPopup()
			}
			imgui.EndDisabled()
		}

		{
			imgui.SameLine()
			imgui.BeginDisabledV(init == nil || a.ReflectionAuxBus == 0)
			if imgui.ArrowButton("##GoToReflectionAuxBus", imgui.DirRight) {
				init.SetActiveBus(a.ReflectionAuxBus)
				init.Focus = be.BankTabBuses
				m.SetNextBank = init
				imgui.SetWindowFocusStr("Buses")
			}
			imgui.EndDisabled()
		}

		val := float32(0)
		idx, pv := p.Prop(wwise.PropTypeReflectionBusVolume)
		in := idx != -1
		if in {
			binary.Decode(pv.V, wio.ByteOrder, &val)
		}

		{
			imgui.BeginDisabledV(!root && !a.OverrideReflectionAuxBus()) 
			{
				imgui.BeginDisabledV(!in)
				imgui.SetNextItemWidth(96)
				imgui.SameLine()
				if imgui.SliderFloat("##ReflectionAuxBusFader", &val, -96.0, 12.0) {
					p.SetPropByIdxF32(idx, val)
				}
				imgui.SameLine()
				imgui.SetNextItemWidth(96)
				if imgui.InputFloat("##ReflectionAuxBusInput", &val, ) {
					if val >= -96.0 && val <= 12.0 {
						p.SetPropByIdxF32(idx, val)
					}
				}
				imgui.EndDisabled()
			}
			{
				imgui.BeginDisabledV(in)
				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(16)
				imgui.PushIDStr("AddReflectionAuxBusSend")
				imgui.SameLine()
				if imgui.Button("+") {
					p.Add(wwise.PropTypeReflectionBusVolume)
				}
				imgui.PopID()
				imgui.EndDisabled()
			}
			{
				imgui.BeginDisabledV(!in)
				imgui.TableSetColumnIndex(6)
				imgui.SetNextItemWidth(16)
				imgui.PushIDStr("RmReflectionAuxBusSend")
				imgui.SameLine()
				if imgui.Button("X") {
					p.Remove(wwise.PropTypeReflectionBusVolume)
				}
				imgui.PopID()
				imgui.EndDisabled()
			}
			imgui.EndDisabled()
		}

		imgui.TreePop()
	}
}
