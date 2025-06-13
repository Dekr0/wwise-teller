package ui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankExplorer(bnkMngr *BankManager, conf *config.Config, modalQ *ModalQ, loop *async.EventLoop) {
	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)
	renderBankExplorerMenu(bnkMngr, conf, modalQ, loop)
	if imgui.BeginTabBarV("BankExplorerTabBar", DefaultTabFlags) {
		paths := []string{}
		bnkMngr.Banks.Range(func(key any, value any) bool {
			open := true
			path := key.(string)
			tab := value.(*BankTab)

			imgui.PushIDStr(path)
			if imgui.BeginTabItemV(filepath.Base(path), &open, 0) {
				renderBankExplorerTab(path, value.(*BankTab))
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

func renderBankExplorerTab(path string, t *BankTab) {
	imgui.Text("Sound bank: " + path)
	if imgui.BeginTabBar("SubBankExplorerTabBar") {
		if imgui.BeginTabItem("Actor Mixer Hierarchy") {
			renderActorMixerHircTable(t)
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Interative Music Hierarchy") {
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Attenuation") {
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Events") {
			renderEventsTable(t)
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Game Sync") {
			imgui.EndTabItem()
		}
		imgui.EndTabBar()
	}
}

func renderBankExplorerMenu(
	bnkMngr *BankManager,
	conf *config.Config,
	modalQ *ModalQ,
	loop *async.EventLoop,
) {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.BeginMenuV("Save", !bnkMngr.WriteLock.Load()) {
				bnkMngr.Banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						pushSaveSoundBankModal(modalQ, loop, conf, bnkMngr, value.(*BankTab), key.(string))
						return false
					}
					return true
				})
				imgui.EndMenu()
			}
			if imgui.BeginMenu("Project") {
				if imgui.MenuItemBool("Set Selected Bank As Init.bnk") {
					bnkMngr.InitBank = bnkMngr.ActiveBank.Bank
					deleteKey := ""
					bnkMngr.Banks.Range(func(key, value any) bool {
						if value.(*BankTab) == bnkMngr.ActiveBank {
							deleteKey = key.(string)
							return false
						}
						return true
					})
					bnkMngr.ActiveBank = nil
					bnkMngr.CloseBank(deleteKey)
				}
				if imgui.BeginMenu("Set Init.bnk Using Existed Banks") {
					deleteKey := ""
					bnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							bnkMngr.InitBank = value.(*BankTab).Bank
							if bnkMngr.ActiveBank == value.(*BankTab) {
								bnkMngr.ActiveBank = nil
							}
							deleteKey = key.(string)
							return false
						}
						return true
					})
					bnkMngr.CloseBank(deleteKey)
				}
				if imgui.BeginMenu("Unmount Init.bnk") {
					bnkMngr.InitBank = nil
				}
			}
			if imgui.BeginMenu("Integration") {
				if imgui.BeginMenuV("Helldivers 2", !bnkMngr.WriteLock.Load()) {
					bnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							pushHD2PatchModal(modalQ, loop, conf, bnkMngr, value.(*BankTab), key.(string))
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

func renderActorMixerHircTable(t *BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	filterState := &t.ActorMixerViewer.ActorMixerHircFilter

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

	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("LinearTable", 2, flags, outerSize, 0) {
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
			int32(len(filterState.ActorMixerHircs)),
		)
		storage.ApplyRequests(msIO)

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.ActorMixerHircs)))
		if msIO.RangeSrcItem() != 1 {
			// Ensure RangeSrc item is not clipped
			clipper.IncludeItemByIndex(int32(msIO.RangeSrcItem()))
		}
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				o := filterState.ActorMixerHircs[n]

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
