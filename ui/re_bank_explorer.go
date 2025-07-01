package ui

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderBankExplorer(bnkMngr *be.BankManager) {
	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)
	renderBankExplorerMenu(bnkMngr)
	if imgui.BeginTabBarV("BankExplorerTabBar", DefaultTabFlags) {
		paths := []string{}
		bnkMngr.Banks.Range(func(key any, value any) bool {
			open := true
			path := key.(string)
			tab := value.(*be.BankTab)

			imgui.PushIDStr(filepath.Base(path))
			selected := imgui.TabItemFlagsNone
			if bnkMngr.SetNextBank != nil && bnkMngr.SetNextBank == tab {
				selected = imgui.TabItemFlagsSetSelected
				bnkMngr.SetNextBank = nil
			}
			if imgui.BeginTabItemV(filepath.Base(path), &open, selected) {
				renderBankExplorerTab(path, value.(*be.BankTab))
				bnkMngr.ActiveBank = tab
				bnkMngr.ActivePath = path
				imgui.EndTabItem()
			}
			imgui.PopID()

			if !open {
				if bnkMngr.ActiveBank == tab {
					bnkMngr.ActiveBank = nil
					bnkMngr.ActivePath = ""
				}
				if bnkMngr.InitBank == tab {
					bnkMngr.InitBank = nil
				}
				if bnkMngr.SetNextBank == tab {
					bnkMngr.SetNextBank = nil
				}
				paths = append(paths, path)
			}

			return true
		})
		imgui.EndTabBar()
		for _, path := range paths {
			bnkMngr.CloseBank(path)
		}
	}
	imgui.End()
}

func renderBankExplorerTab(path string, t *be.BankTab) {
	imgui.PushTextWrapPos()
	imgui.Text("Sound bank: " + path)
	imgui.PopTextWrapPos()
	if imgui.BeginTabBarV("SubBankExplorerTabBar", DefaultTabFlags) {
		selected := imgui.TabItemFlagsNone
		if t.Focus == be.BankTabActorMixer {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Actor Mixer", nil, selected) {
			renderActorMixerHircTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabMusic {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Music", nil, selected) {
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabAttenuation {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Attenuation", nil, selected) {
			renderAttenuationTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabBuses {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Buses", nil, selected) {
			renderBusTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabFX {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Fx", nil, selected) {
			renderFxTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabModulator {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Modulators", nil, selected) {
			renderModulatorTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabEvents {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Events", nil, selected) {
			renderEventsTable(t)
			imgui.EndTabItem()
		}

		selected = imgui.TabItemFlagsNone
		if t.Focus == be.BankTabGameSync {
			selected = imgui.TabItemFlagsSetSelected
			t.Focus = be.BankTabNone
		}
		if imgui.BeginTabItemV("Game Sync", nil, selected) {
			selected = imgui.TabItemFlagsNone
			imgui.EndTabItem()
		}
		imgui.EndTabBar()
	}
}

func renderBankExplorerMenu(bnkMngr *be.BankManager) {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.BeginMenuV("Save", !bnkMngr.WriteLock.Load()) {
				bnkMngr.Banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						pushSaveSoundBankModal(bnkMngr, value.(*be.BankTab), key.(string))
						return false
					}
					return true
				})
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Project") {
				imgui.BeginDisabledV(bnkMngr.ActiveBank == nil)
				if imgui.MenuItemBool("Set Selected Bank As Init.bnk") {
					bnkMngr.InitBank = bnkMngr.ActiveBank
				}
				imgui.EndDisabled()
				if imgui.BeginMenu("Set Init.bnk Using Existed Banks") {
					bnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							bnkMngr.InitBank = value.(*be.BankTab)
							return false
						}
						return true
					})
					imgui.EndMenu()
				}
				if imgui.MenuItemBool("Unmount Init.bnk") {
					bnkMngr.InitBank = nil
				}
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Integration") {
				if imgui.BeginMenuV("Helldivers 2", !bnkMngr.WriteLock.Load()) {
					bnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							pushHD2PatchModal(bnkMngr, value.(*be.BankTab), key.(string))
							return false
						}
						return true
					})
					imgui.EndMenu()
				}
				imgui.EndMenu()
			}

			imgui.EndMenu()
		}
		imgui.EndMenuBar()
	}
}

func renderActorMixerHircTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	filterState := &t.ActorMixerViewer.HircFilter

	imgui.SeparatorText("Filter")

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
	) {
		t.FilterActorMixerHircs()
	}

	imgui.BeginDisabledV(filterState.Type != wwise.HircTypeAll && filterState.Type != wwise.HircTypeSound)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By source ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Sid)),
	) {
		t.FilterActorMixerHircs()
	}
	imgui.EndDisabled()

	imgui.SetNextItemWidth(256)
	preview := wwise.HircTypeName[filterState.Type]
	if imgui.BeginCombo("By Type", preview) {
		var filter func() = nil
		for _, _type := range wwise.ActorMixerHircTypes {
			selected := filterState.Type == _type
			preview = wwise.HircTypeName[_type]
			if imgui.SelectableBoolPtr(preview, &selected) {
				filterState.Type = _type
				filter = t.FilterActorMixerHircs
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

	if imgui.BeginTableV("LinearTable", 2, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		storage := t.ActorMixerViewer.LinearStorage

		msIO := imgui.BeginMultiSelectV(
			DefaultMultiSelectFlags,
			storage.Size(),
			int32(len(filterState.Hircs)),
		)
		storage.ApplyRequests(msIO)

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Hircs)))
		if msIO.RangeSrcItem() != 1 {
			// Ensure RangeSrc item is not clipped
			clipper.IncludeItemByIndex(int32(msIO.RangeSrcItem()))
		}
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				o := filterState.Hircs[n]

				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)

				var idS string
				var selected bool = false
				id, err := o.HircID()
				if err != nil {
					idS = "-"
				} else {
					idS = strconv.FormatUint(uint64(id), 10)
					selected = storage.Contains(imgui.ID(id))
					imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(id))
				}

				if err != nil {
					imgui.PushIDStr(fmt.Sprintf("UnknownID%d", n))
				}
				flags := imgui.SelectableFlagsSpanAllColumns | 
					     imgui.SelectableFlagsAllowOverlap
				size := imgui.NewVec2(0, 0)
				imgui.SelectableBoolPtrV(idS, &selected, flags, size)

				if imgui.BeginPopupContextItem() {
					renderActorMixerHircTableCtx(t, o, id)
					imgui.EndPopup()
				}

				imgui.TableSetColumnIndex(1)
				st := wwise.HircTypeName[o.HircType()]
				if o.HircType() == wwise.HircTypeSound {
					st = fmt.Sprintf(
						"%s (Audio Source %d)",
						st, o.(*wwise.Sound).BankSourceData.SourceID,
					)
				}
				imgui.Text(st)

				if err != nil {
					imgui.PopID()
				}
			}
		}
		msIO = imgui.EndMultiSelect()
		storage.ApplyRequests(msIO)
		imgui.EndTable()
	}
}

func renderActorMixerHircTableCtx(t *be.BankTab, o wwise.HircObj, id uint32) {
	if imgui.SelectableBool("Expand in Actor Mixer Hierarchy View") {
		t.OpenActorMixerHircNode(id)
		return
	}

	imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
	if imgui.SelectableBool("Copy ID") {
		clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(id), 10)))
		imgui.EndDisabled()
		return
	}
	imgui.EndDisabled()

	switch sound := o.(type) {
	case *wwise.Sound:
		imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
		if imgui.SelectableBool("Copy Source ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10)))
			imgui.EndDisabled()
			return
		}
		imgui.EndDisabled()
	}

	leafs := o.Leafs()
	if len(leafs) <= 0 {
		return
	}

	imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
	if imgui.SelectableBool("Copy Leafs' IDs") {
		var builder strings.Builder
		for _, l := range leafs {
			if _, err := builder.WriteString(fmt.Sprintf("%d\n", l)); err != nil {
				slog.Error(fmt.Sprintf("Failed to copy leafs' IDs of hierarchy object %d", id), "error", err)
				return
			}
		}
		clipboard.Write(clipboard.FmtText, []byte(builder.String()))
		imgui.EndDisabled()
		return
	}
	imgui.EndDisabled()
	
	imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
	if imgui.SelectableBool("Copy Leafs' Source IDs") {
		h := t.Bank.HIRC()
		if h == nil {
			panic("Context menu is shown in actor mixer hierarchy table but no HIRC chunk is found.")
		}
		var builder strings.Builder
		for _, l := range leafs {
			v, in := h.ActorMixerHirc.Load(l)
			if !in {
				panic(fmt.Sprintf("No actor mixer hierarchy object has ID %d", l))
			}
			switch sound := v.(wwise.HircObj).(type) {
			case *wwise.Sound:
				if _, err := builder.WriteString(fmt.Sprintf("%d\n", sound.BankSourceData.SourceID)); err != nil {
					slog.Error(fmt.Sprintf("Failed to copy leafs' Source IDs of hierarchy object %d", id), "error", err)
					return
				}
			}
		}
		clipboard.Write(clipboard.FmtText, []byte(builder.String()))
		imgui.EndDisabled()
		return
	}
	imgui.EndDisabled()
}
