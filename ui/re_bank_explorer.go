package ui

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	wutils "github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderBankExplorer() {
	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.BankExplorerTag], nil, imgui.WindowFlagsMenuBar)
	renderBankExplorerMenu()
	if imgui.BeginTabBarV("BankExplorerTabBar", DefaultTabFlags) {
		paths := []string{}
		BnkMngr.Banks.Range(func(key any, value any) bool {
			open := true
			path := key.(string)
			tab := value.(*be.BankTab)

			imgui.PushIDStr(path)
			selected := imgui.TabItemFlagsNone
			if BnkMngr.SetNextBank != nil && BnkMngr.SetNextBank == tab {
				selected = imgui.TabItemFlagsSetSelected
				BnkMngr.SetNextBank = nil
			}
			if imgui.BeginTabItemV(filepath.Base(path), &open, selected) {
				renderBankExplorerTab(path, value.(*be.BankTab))
				BnkMngr.ActiveBank = tab
				BnkMngr.ActivePath = path
				imgui.EndTabItem()
			}
			imgui.PopID()

			if !open {
				if BnkMngr.ActiveBank == tab {
					BnkMngr.ActiveBank = nil
					BnkMngr.ActivePath = ""
				}
				if BnkMngr.InitBank == tab {
					BnkMngr.InitBank = nil
				}
				if BnkMngr.SetNextBank == tab {
					BnkMngr.SetNextBank = nil
				}
				paths = append(paths, path)
			}

			return true
		})
		imgui.EndTabBar()
		for _, path := range paths {
			BnkMngr.CloseBank(path)
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

func renderBankExplorerMenu() {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.BeginMenuV("Save (With Metadata)", !BnkMngr.WriteLock.Load()) {
				BnkMngr.Banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						pushSaveSoundBankModal(value.(*be.BankTab), key.(string), false)
						return false
					}
					return true
				})
				imgui.EndMenu()
			}
			if imgui.BeginMenuV("Save (Without Metadata)", !BnkMngr.WriteLock.Load()) {
				BnkMngr.Banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						pushSaveSoundBankModal(value.(*be.BankTab), key.(string), true)
						return false
					}
					return true
				})
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Project") {
				Disabled(BnkMngr.ActiveBank == nil, func() {
					if imgui.MenuItemBool("Set Selected Bank As Init.bnk") {
						BnkMngr.InitBank = BnkMngr.ActiveBank
					}
				})
				if imgui.BeginMenu("Set Init.bnk Using Existed Banks") {
					BnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							BnkMngr.InitBank = value.(*be.BankTab)
							return false
						}
						return true
					})
					imgui.EndMenu()
				}
				if imgui.MenuItemBool("Unmount Init.bnk") {
					BnkMngr.InitBank = nil
				}
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Integration") {
				if imgui.BeginMenuV("Helldivers 2", !BnkMngr.WriteLock.Load()) {
					BnkMngr.Banks.Range(func(key, value any) bool {
						if imgui.MenuItemBool(key.(string)) {
							pushHD2PatchModal(value.(*be.BankTab), key.(string))
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

// Save sound bank modal
func pushSaveSoundBankModal(saveTab *be.BankTab, saveName string, excludeMeta bool) {
	onSave := onSaveSoundBankModal(saveTab, saveName, excludeMeta)
	renderF, done, err := saveFileDialogFunc(onSave, GCtx.Config.Home)
	if err != nil {
		msg := fmt.Sprintf("Failed create save file dialog for saving sound bank %s", saveName)
		slog.Error(msg, "error", err)
	} else {
		Modal(done, 0, fmt.Sprintf("Save sound bank %s to ...", saveName), renderF, nil)
	}
}

func onSaveSoundBankModal(
	saveTab *be.BankTab,
	saveName string,
	excludeMeta bool,
) func(string) {
return func(path string) {
	onProcMsg := fmt.Sprintf("Saving sound bank %s to %s", saveName, path)
	onDoneMsg := fmt.Sprintf("Saved sound bank %s to %s", saveName, path)
	f := saveSoundBank(path, saveTab, saveName, excludeMeta)
	BG(time.Second * 8, onProcMsg, onDoneMsg, f)
}}

func saveSoundBank(
	path string,
	saveTab *be.BankTab,
	saveName string,
	excludeMeta bool,
) func(context.Context) {
return func(ctx context.Context) {
	BnkMngr.WriteLock.Store(true)
	data, err := saveTab.Encode(ctx, excludeMeta)
	if err != nil {
		msg := fmt.Sprintf("Failed to encode sound bank %s", saveName)
		slog.Error(msg, "error", err)
		BnkMngr.WriteLock.Store(false)
		return
	}
	dest := filepath.Join(path, filepath.Base(saveName))
	if err := wutils.SaveFileWithRetry(data, dest); err != nil {
		msg := fmt.Sprintf("Failed to save sound bank %s to %s", saveName, path)
		slog.Error(msg, "error", err)
	}
	BnkMngr.WriteLock.Store(false)
}}
// End of save sound bank modal


// Start of push HD2 Patch Modal
func pushHD2PatchModal(saveTab *be.BankTab, saveName string) {
	onSave := func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 4,
		)
		
		onProcMsg := fmt.Sprintf("Saving sound bank %s to HD2 patch %s",
			saveName, path)
		onDoneMsg := fmt.Sprintf("Saved sound bank %s to HD2 patch %s",
			saveName, path)

		if err := GCtx.Loop.QTask(timeout, cancel, onProcMsg, onDoneMsg, 
			func(ctx context.Context) {
				slog.Info(onProcMsg)
				BnkMngr.WriteLock.Store(true)
				defer BnkMngr.WriteLock.Store(false)
				bnkData, err := saveTab.Encode(ctx, true)
				if err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", saveName),
						"error", err,
					)
					return
				}
				meta := saveTab.Bank.META()
				if meta == nil {
					slog.Error(
						fmt.Sprintf("Sound bank %s is missing integration data.",
							saveName),
					)
					return
				}
				if err := helldivers.GenHelldiversPatchStable(
					bnkData, meta.B, path,
				); err != nil {
					slog.Error(fmt.Sprintf("Failed to write HD2 patch to %s", path))
				} else {
					slog.Info(onDoneMsg)
				}
			},
		);
		   err != nil {
			slog.Error(fmt.Sprintf("Failed to save HD2 patch to %s", path),
				"error", err,
			)
		}
	}

	if renderF, done, err := saveFileDialogFunc(onSave, GCtx.Config.Home);
	   err != nil {
		msg := fmt.Sprintf("Failed create save file dialog for saving sound bank %s to HD2 patch", saveName)
		slog.Error(msg, "error", err)
	} else {
		Modal(done, 0, fmt.Sprintf("Save sound bank %s to HD2 patch ...", saveName), renderF, nil)
	}
}

// End of push HD2 Patch Modal

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

	Disabled(filterState.Type != wwise.HircTypeAll && filterState.Type != wwise.HircTypeSound, func() {
		imgui.SetNextItemWidth(96)
		if imgui.InputScalar(
			"By source ID",
			imgui.DataTypeU32,
			uintptr(utils.Ptr(&filterState.Sid)),
		) {
			t.FilterActorMixerHircs()
		}
	})

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

				id, err := o.HircID()
				if err != nil {
					panic(fmt.Sprintf("Actor mixer hierarchy object should return an id without error"))
				}

				selected := t.ActorMixerViewer.Selected(id)
				imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(t.ActorMixerViewer.GetSelectionHash(id)))

				flags := imgui.SelectableFlagsSpanAllColumns | 
					     imgui.SelectableFlagsAllowOverlap
				size := imgui.NewVec2(0, 0)
				imgui.SelectableBoolPtrV(strconv.FormatUint(uint64(id), 10), &selected, flags, size)

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

	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(id), 10)))
		}
	})

	switch sound := o.(type) {
	case *wwise.Sound:
		Disabled(!GCtx.CopyEnable, func() {
			if imgui.SelectableBool("Copy Source ID") {
				clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10)))
			}
		})
	}

	leafs := o.Leafs()
	if len(leafs) <= 0 {
		return
	}

	Disabled(!GCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy Leafs' IDs") {
			var builder strings.Builder
			for _, l := range leafs {
				if _, err := builder.WriteString(fmt.Sprintf("%d\n", l)); err != nil {
					slog.Error(fmt.Sprintf("Failed to copy leafs' IDs of hierarchy object %d", id), "error", err)
					return
				}
			}
			clipboard.Write(clipboard.FmtText, []byte(builder.String()))
		}
	})
	
	Disabled(!GCtx.CopyEnable, func() {
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
		}
	})

	if imgui.SelectableBool("Search For Events and Actions") {
		if t.SearchNearestEventAction(id) {
			t.Focus = be.BankTabEvents
			imgui.SetWindowFocusStr("Events")
		}
	}
}
