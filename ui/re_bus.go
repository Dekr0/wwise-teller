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
	"golang.design/x/clipboard"
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
		text := ""
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

				if imgui.BeginPopupContextItem() {
					cloneId := idA
					Disabled(!GlobalCtx.CopyEnable, func() {
						if imgui.SelectableBool("Copy ID") {
							clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(cloneId), 10)))
						}
						if imgui.SelectableBool("Expand Node in Hierarchy") {
							t.OpenBusHircNode(cloneId)
						}
					})
					imgui.EndPopup()
				}

				switch b := bus.(type) {
				case *wwise.Bus:
					if b.CanSetHDR == -1 {
						panic(fmt.Sprintf("HDR Availability isn't resolve for Bus %d", b.Id))
					}
					if b.IsHDRBus() {
						text = fmt.Sprintf("%s (HDR Bus)", wwise.HircTypeName[bus.HircType()])
					} else {
						text = wwise.HircTypeName[bus.HircType()]
					}
				case *wwise.AuxBus:
					text = wwise.HircTypeName[bus.HircType()]
				}
				imgui.TableSetColumnIndex(1)
				imgui.Text(text)
			}
		}
		imgui.EndTable()
	}
}

func renderMasterMixerHierarchy(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Master Mixer Hierarchy", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil {
		return 
	}
	renderMasterMixerHierarchyTreeTable(t)
}

func renderMasterMixerHierarchyTreeTable(t *be.BankTab) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	if imgui.BeginTableV("MasterMixerHierarchy", 2, flags, DefaultSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		for _, root := range t.Bank.HIRC().BusRoots {
			renderMasterMixerHierarchyNode(t, root)
		}
		imgui.EndTable()
	}
}

func renderMasterMixerHierarchyNode(t *be.BankTab, node *wwise.BusHircNode) {
	o := node.Obj

	var sid string
	selected := false
	id, err := o.HircID()
	if err != nil { panic(err) }
	sid = strconv.FormatUint(uint64(id), 10) 
	selected = t.BusViewer.ActiveBus == o

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	flags := DefaultTreeFlags
	if selected {
		flags |= imgui.TreeNodeFlagsSelected
	}
	imgui.SetNextItemOpen(node.Open)
	open := imgui.TreeNodeExStrV(sid, flags)
	if imgui.IsItemClicked() {
		t.BusViewer.ActiveBus = o
	}
	node.Open = open
	if imgui.BeginPopupContextItem() {
		Disabled(!GlobalCtx.CopyEnable, func() {
			if imgui.SelectableBool("Copy ID") {
				clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(id), 10)))
			}
		})
		imgui.EndPopup()
	}

	imgui.TableSetColumnIndex(1)
	st := ""
	switch bus := o.(type) {
	case *wwise.Bus:
		if bus.CanSetHDR == -1 {
			panic(fmt.Sprintf("HDR Availability isn't resolve for Bus %d", bus.Id))
		}
		if bus.IsHDRBus() {
			st = fmt.Sprintf("%s (HDR Bus)", wwise.HircTypeName[bus.HircType()])
		} else {
			st = wwise.HircTypeName[bus.HircType()]
		}
	case *wwise.AuxBus:
		st = wwise.HircTypeName[bus.HircType()]
	}
	imgui.Text(st)
	if open {
		for _, leafIdx := range node.LeafsIdx {
			if leaf, in := node.Leafs[leafIdx]; !in {
				panic(fmt.Sprintf("Failed to locate leaf using leaf index in Bus %d", id))
			} else {
				renderMasterMixerHierarchyNode(t, leaf)
			}
		}
		imgui.TreePop()
	}
}

func renderBusViewer(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Buses", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		return
	}
	if t.BusViewer.ActiveBus != nil {
		switch bus := t.BusViewer.ActiveBus.(type) {
		case *wwise.Bus:
			renderBus(t, bus)
		case *wwise.AuxBus:
			renderAuxBus(t, bus)
		default:
			panic("Non Master Mix Hierarchy object is in the bus viewer")
		}
	}
}

func renderBus(t *be.BankTab, b *wwise.Bus) {
	imgui.Text(fmt.Sprintf("Bus %d", b.Id))
	imgui.Text(fmt.Sprintf("Override Bus ID (Parent Bus) %d", b.OverrideBusId))
	if b.OverrideBusId != 0 {
		imgui.SameLine()
		if imgui.ArrowButton(fmt.Sprintf("GoToOverrideBus%d", b.OverrideBusId), imgui.DirRight) {
			t.SetActiveBus(b.OverrideBusId)
			imgui.SetWindowFocusStr("Buses")
		}
	}
	renderBusAuxParam(t, &b.AuxParam, &b.PropBundle)
	renderBusEarlyReflection(t, &b.AuxParam, &b.PropBundle)
	renderHDR(b)
	renderBusAdvanceSetting(b)
	renderAllProp(&b.PropBundle, nil)
	renderBusFxParam(t, &b.BusFxParam)
}

func renderAuxBus(t *be.BankTab, b *wwise.AuxBus) {
	imgui.Text(fmt.Sprintf("Auxiliary Bus %d", b.Id))
	imgui.Text(fmt.Sprintf("Override Bus ID (Parent Bus) %d", b.OverrideBusId))
	if b.OverrideBusId != 0 {
		imgui.SameLine()
		if imgui.ArrowButton(fmt.Sprintf("GoToOverrideBus%d", b.OverrideBusId), imgui.DirRight) {
			t.SetActiveBus(b.OverrideBusId)
			imgui.SetWindowFocusStr("Buses")
		}
	}
	renderBusAuxParam(t, &b.AuxParam, &b.PropBundle)
	renderBusEarlyReflection(t, &b.AuxParam, &b.PropBundle)
	renderAuxBusAdvanceSetting(b)
	renderAllProp(&b.PropBundle, nil)
	renderBusFxParam(t, &b.BusFxParam)
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

func renderBusAdvanceSetting(b *wwise.Bus) {
	if imgui.TreeNodeExStr("Advance Settings") {
		DefaultSize.Y = 128
		imgui.BeginChildStrV("BusAdvanceSetting", DefaultSize, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Playback Limit")
		{
			imgui.BeginDisabledV(b.OverrideBusId == 0)
			ignoreParent := b.IgnoreParentMaxNumInstance()
			if imgui.Checkbox("Ignore Parent", &ignoreParent) {
				b.SetIgnoreParentMaxNumInstance(ignoreParent)
			}	
			imgui.EndDisabled()
		}
		{
			imgui.BeginDisabledV(b.OverrideBusId != 0 && !b.IgnoreParentMaxNumInstance())
			{
				imgui.Text("Limit sound instances to:")
				imgui.SameLine()
				maxNumInstance := int32(b.MaxNumInstance)
				imgui.SetNextItemWidth(96)
				// Zero is no limiting
				if imgui.InputInt("##MaxNumInstance", &maxNumInstance) {
					if maxNumInstance >= 0 && maxNumInstance <= 1000 {
						b.MaxNumInstance = uint16(maxNumInstance)
					}
				}
			}
			var preview string
			{
				imgui.Text("When limit is reached:")
				imgui.SameLine()

				if !b.UseVirtualBehavior() {
					preview = "Kill voice"
				} else {
					preview = "Use virtual voice setting"
				}

				imgui.SetNextItemWidth(200)
				if imgui.BeginCombo("##ReachPlaybackLimitBehavior", preview) {
					selected := !b.UseVirtualBehavior()
					if imgui.SelectableBoolPtr("Kill voice", &selected) {
						b.SetUseVirtualBehavior(false)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					selected = b.UseVirtualBehavior()
					if imgui.SelectableBoolPtr("Use virtual voice setting", &selected) {
						b.SetUseVirtualBehavior(true)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					imgui.EndCombo()
				}
				imgui.SameLine()
				imgui.Text("for lowest priroty")
			}
			{
				imgui.Text("When priority is equal:")
				imgui.SameLine()
				if !b.KillNewest() {
					preview = "Discard oldest instance"
				} else {
					preview = "Discard newest instance"
				}
				imgui.SetNextItemWidth(192)
				if imgui.BeginCombo("##PriorityEqualBehavior", preview) {
					selected := !b.KillNewest()
					if imgui.SelectableBoolPtr("Discard oldest instance", &selected) {
						b.SetKillNewest(false)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					selected = b.KillNewest()
					if imgui.SelectableBoolPtr("Discard newest instance", &selected) {
						b.SetKillNewest(true)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					imgui.EndCombo()
				}
			}
			imgui.EndDisabled()
		}
		imgui.EndChild()
		imgui.TreePop()
		DefaultSize.Y = 0
	}
}

func renderHDR(b *wwise.Bus) {
	if imgui.TreeNodeExStr("HDR") {
		if b.CanSetHDR == -1 {
			panic(fmt.Sprintf("HDR Availability isn't resolve for Bus %d", b.Id))
		}
		{
			imgui.BeginDisabledV(b.CanSetHDR == 0)

			if b.CanSetHDR == 0 {
				imgui.Text("HDR is enabled in a parent bus.")
			} else {
				HDREnable := b.IsHDRBus()
				imgui.BeginDisabled()
				if imgui.Checkbox("Enable HDR", &HDREnable) {
					b.SetHDRBus(HDREnable)
				}
				imgui.EndDisabled()
			}
			{
				DefaultSize.Y = 128
				imgui.BeginChildStrV("HDRDynamics", DefaultSize, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
				imgui.SeparatorText("Dynamics")
				imgui.EndChild()
				DefaultSize.Y = 0
			}
			imgui.EndDisabled()
		}
		imgui.TreePop()
	}
}

func renderAuxBusAdvanceSetting(b *wwise.AuxBus) {
	if imgui.TreeNodeExStr("Advance Settings") {
		DefaultSize.Y = 128
		imgui.BeginChildStrV("BusAdvanceSetting", DefaultSize, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Playback Limit")
		{
			imgui.BeginDisabledV(b.OverrideBusId == 0)
			ignoreParent := b.IgnoreParentMaxNumInstance()
			if imgui.Checkbox("Ignore Parent", &ignoreParent) {
				b.SetIgnoreParentMaxNumInstance(ignoreParent)
			}	
			imgui.EndDisabled()
		}
		{
			imgui.BeginDisabledV(b.OverrideBusId != 0 && !b.IgnoreParentMaxNumInstance())
			{
				imgui.Text("Limit sound instances to:")
				imgui.SameLine()
				maxNumInstance := int32(b.MaxNumInstance)
				imgui.SetNextItemWidth(96)
				// Zero is no limiting
				if imgui.InputInt("##MaxNumInstance", &maxNumInstance) {
					if maxNumInstance >= 0 && maxNumInstance <= 1000 {
						b.MaxNumInstance = uint16(maxNumInstance)
					}
				}
			}
			var preview string
			{
				imgui.Text("When limit is reached:")
				imgui.SameLine()

				if !b.UseVirtualBehavior() {
					preview = "Kill voice"
				} else {
					preview = "Use virtual voice setting"
				}

				imgui.SetNextItemWidth(200)
				if imgui.BeginCombo("##ReachPlaybackLimitBehavior", preview) {
					selected := !b.UseVirtualBehavior()
					if imgui.SelectableBoolPtr("Kill voice", &selected) {
						b.SetUseVirtualBehavior(false)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					selected = b.UseVirtualBehavior()
					if imgui.SelectableBoolPtr("Use virtual voice setting", &selected) {
						b.SetUseVirtualBehavior(true)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					imgui.EndCombo()
				}
				imgui.SameLine()
				imgui.Text("for lowest priroty")
			}
			{
				imgui.Text("When priority is equal:")
				imgui.SameLine()
				if !b.KillNewest() {
					preview = "Discard oldest instance"
				} else {
					preview = "Discard newest instance"
				}
				imgui.SetNextItemWidth(192)
				if imgui.BeginCombo("##PriorityEqualBehavior", preview) {
					selected := !b.KillNewest()
					if imgui.SelectableBoolPtr("Discard oldest instance", &selected) {
						b.SetKillNewest(false)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					selected = b.KillNewest()
					if imgui.SelectableBoolPtr("Discard newest instance", &selected) {
						b.SetKillNewest(true)
					}
					if selected {
						imgui.SetItemDefaultFocus()
					}
					imgui.EndCombo()
				}
			}
			imgui.EndDisabled()
		}
		imgui.EndChild()
		DefaultSize.Y = 0
		imgui.TreePop()
	}
}

func renderBusFxParam(t *be.BankTab, b *wwise.BusFxParam) {
	if imgui.TreeNodeStr("FX") {
		{
			imgui.BeginDisabled()
			byPassFx := b.FxChunk.BitsFxByPass != 0
			imgui.Checkbox("By Passing FX", &byPassFx)
			imgui.EndDisabled()
		}

		if imgui.BeginTableV("FXChunkTable", 3, DefaultTableFlags, DefaultSize, 0) {
			imgui.TableSetupColumn("Unique FX Index")
			imgui.TableSetupColumn("FX ID")
			imgui.TableSetupColumn("Is Share Set")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()
			for i := range b.FxChunk.FxChunkItems {
				fi := &b.FxChunk.FxChunkItems[i]

				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)
				imgui.Text(strconv.FormatUint(uint64(fi.UniqueFxIndex), 10))

				imgui.TableSetColumnIndex(1)
				{
					imgui.Text(strconv.FormatUint(uint64(fi.FxId), 10))
					imgui.SameLine()

					imgui.BeginDisabled()
					imgui.ArrowButton(fmt.Sprintf("##SetFXID%d", i), imgui.DirDown)
					imgui.EndDisabled()

					imgui.SameLine()
					if imgui.ArrowButton(fmt.Sprintf("##GoToFXID%d", i), imgui.DirRight) {
						t.SetActiveFX(fi.FxId)
					}
				}

				imgui.TableSetColumnIndex(2)
				imgui.BeginDisabled()
				isSharedSet := fi.BitIsShareSet != 0
				imgui.Checkbox(fmt.Sprintf("##IsShareSet%d", i), &isSharedSet)
				imgui.EndDisabled()
			}

			imgui.EndTable()
		}

		imgui.Text(fmt.Sprintf("FX ID: %d", b.FxID_0))
		imgui.SameLine()

		imgui.BeginDisabled()
		imgui.ArrowButton("##SetFXID0%d", imgui.DirDown)
		imgui.EndDisabled()

		imgui.SameLine()
		if imgui.ArrowButton("##GoToFXID0", imgui.DirRight) {
			t.SetActiveFX(b.FxID_0)
		}

		imgui.BeginDisabled()
		isSharedSet := b.IsShareSet_0 != 0
		imgui.Checkbox("Is Share Set", &isSharedSet)
		imgui.EndDisabled()

		imgui.TreePop()
	}
}
