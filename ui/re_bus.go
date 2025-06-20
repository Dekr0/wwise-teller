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

func renderBusTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	filterState := &t.BusViewer.Filter

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
		) {
		t.FilterBuses()
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(80)
	preview := wwise.HircTypeName[filterState.Type]
	if imgui.BeginCombo("By Type", preview) {
		var filter func() = nil
		for _, _type := range wwise.BusHircTypes {
			selected := filterState.Type == _type
			preview = wwise.HircTypeName[_type]
			if imgui.SelectableBoolPtr(preview, &selected) {
				filterState.Type = _type
				filter = t.FilterBuses
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

	if imgui.BeginTableV("BusesTable", 2, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}
		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Buses)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				bus := filterState.Buses[n]
				
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				if t.BusViewer.ActiveBus == nil {
					t.BusViewer.ActiveBus = bus
				}
				idA, err := bus.HircID()
				if err != nil { panic(err) }
				idB, err := t.BusViewer.ActiveBus.HircID()
				if err != nil { panic(err) }
				selected := idA == idB
				if imgui.SelectableBoolPtrV(
					strconv.FormatUint(uint64(idA), 10),
					&selected,
					DefaultSelectableFlags,
					DefaultSize,
				) {
					t.BusViewer.ActiveBus = bus
				}

				imgui.TableSetColumnIndex(1)
				imgui.Text(wwise.HircTypeName[bus.HircType()])
			}
		}
		imgui.EndTable()
	}
}

func renderBusViewer(t *be.BankTab) {
	imgui.Begin("Buses")
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		imgui.End()
		return
	}

	if t.BusViewer.ActiveBus != nil {
		switch bus := t.BusViewer.ActiveBus.(type) {
		case *wwise.Bus:
			renderBus(t, bus)
		case *wwise.AuxBus:
			renderAuxBus(t, bus)
		default:
			panic("Panic trap")
		}
	}

	imgui.End()
}

func renderBus(t *be.BankTab, b *wwise.Bus) {
	imgui.Text(fmt.Sprintf("Bus %d", b.Id))
	renderBusAuxParam(t, &b.AuxParam, &b.PropBundle)
	renderBusEarlyReflection(t, &b.AuxParam, &b.PropBundle)
	renderAllProp(&b.PropBundle, nil)
}

func renderAuxBus(t *be.BankTab, b *wwise.AuxBus) {
	imgui.Text(fmt.Sprintf("Auxiliary Bus %d", b.Id))
	renderBusAuxParam(t, &b.AuxParam, &b.PropBundle)
	renderBusEarlyReflection(t, &b.AuxParam, &b.PropBundle)
	renderAllProp(&b.PropBundle, nil)
}

func renderBusAuxParam(t *be.BankTab, a *wwise.AuxParam, p *wwise.PropBundle) {
	if imgui.TreeNodeExStr("User-Defined Auxiliary Send") {
		imgui.BeginDisabled()
		overrideAuxSend := a.OverrideAuxSends()
		imgui.Checkbox("Override User-Defined Auxiliary Send", &overrideAuxSend)
		imgui.EndDisabled()
		renderBusUserAuxSendTable(t, p, a)
		imgui.TreePop()
	}
}

func renderBusUserAuxSendTable(t *be.BankTab, p *wwise.PropBundle, a *wwise.AuxParam) {
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

			imgui.PushIDStr(fmt.Sprintf("ZeroBusUserAuxSend%d", j))
			if imgui.Button("x") {
				a.AuxIds[j] = 0
			}
			imgui.PopID()
			imgui.SameLine()

			imgui.SetNextItemWidth(96)
			preview := strconv.FormatUint(uint64(aid), 10) 
			imgui.Text(preview)
			imgui.SameLine()

			popup := fmt.Sprintf("SetBusUserAuxSendBus%d", j)
			if imgui.ArrowButton(fmt.Sprintf("SetBusUserAuxSendBtn%d", j), imgui.DirDown) {
				imgui.OpenPopupStr(popup)
			}

			if imgui.BeginPopup(popup) {
				filterState := &t.BusViewer.Filter
				imgui.Text("Aux Bus")
				imgui.SameLine()

				imgui.SetNextItemWidth(96)
				if imgui.BeginCombo(fmt.Sprintf("##BusUserAuxSendBus%dCombo", j), preview) {
					for _, b := range filterState.Buses {
						if b.HircType() != wwise.HircTypeAuxBus {
							continue
						}

						selected := false
						optionID, err := b.HircID()
						if err != nil { panic(err) }
						if t.BusViewer.ActiveBus != nil {
							activeId, err := t.BusViewer.ActiveBus.HircID()
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
					t.FilterBuses()
				}
				imgui.EndPopup()
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)
			if imgui.ArrowButton(fmt.Sprintf("##GoToUserAuxSend%d", i), imgui.DirRight) {
				t.SetActiveBus(aid)
				t.Focus = be.BankTabBuses
				imgui.SetWindowFocusStr("Buses")
			}

			val := float32(0.0)
			idx, pv := p.Prop(i)
			in := idx != -1
			if in {
				binary.Decode(pv.V, wio.ByteOrder, &val)
			}

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

			imgui.BeginDisabledV(in)
			imgui.TableSetColumnIndex(5)
			imgui.SetNextItemWidth(16)
			imgui.PushIDStr(fmt.Sprintf("AddUserAuxSend%d", j))
			if imgui.Button("+") {
				p.Add(i)
			}
			imgui.PopID()
			imgui.EndDisabled()

			imgui.BeginDisabledV(!in)
			imgui.TableSetColumnIndex(6)
			imgui.SetNextItemWidth(16)
			imgui.PushIDStr(fmt.Sprintf("RmUserAuxSend%d", j))
			if imgui.Button("X") {
				p.Remove(i)
			}
			imgui.PopID()
			imgui.EndDisabled()

			i += 1
		}
		imgui.EndTable()
	}
	DefaultSize.X = 0
}

func renderBusEarlyReflection(t *be.BankTab, a *wwise.AuxParam, p *wwise.PropBundle) {
	if imgui.TreeNodeExStr("Early Reflections") {
		overrideReflectionAuxBus := a.OverrideReflectionAuxBus()
		imgui.BeginDisabled()
		imgui.Checkbox("Override early reflection bus", &overrideReflectionAuxBus)
		imgui.EndDisabled()

		imgui.Text("Auxiliary Bus")
		{
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
			if imgui.ArrowButton("SetReflectionAuxSendBtn", imgui.DirDown) {
				imgui.OpenPopupStr(popup)
			}

			if imgui.BeginPopup(popup) {
				filterState := &t.BusViewer.Filter
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

						if t.BusViewer.ActiveBus != nil {
							activeId, err := t.BusViewer.ActiveBus.HircID()
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
					t.FilterBuses()
				}
				imgui.EndPopup()
			}
		}

		imgui.SameLine()
		imgui.BeginDisabledV(a.ReflectionAuxBus == 0)
		if imgui.ArrowButton("##GoToReflectionAuxBus", imgui.DirRight) {
			t.SetActiveBus(a.ReflectionAuxBus)
			t.Focus = be.BankTabBuses
			imgui.SetWindowFocusStr("Buses")
		}
		imgui.EndDisabled()

		val := float32(0)
		idx, pv := p.Prop(wwise.PropTypeReflectionBusVolume)
		in := idx != -1
		if in {
			binary.Decode(pv.V, wio.ByteOrder, &val)
		}

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

		imgui.TreePop()
	}
}
