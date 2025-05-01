package ui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
)

func showBankExplorer(
	bnkMngr *BankManager,
	saveActive bool,
) (*bankTab, string, *bankTab, string) {
	var activeTab *bankTab = nil
	var closeTab string = "" 

	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)

	saveTab, saveName := showBankExplorerMenu(bnkMngr)

	if imgui.BeginTabBarV(
		"BankExplorerTabBar",
		imgui.TabBarFlagsReorderable | imgui.TabBarFlagsAutoSelectNewTabs | 
		imgui.TabBarFlagsTabListPopupButton | imgui.TabBarFlagsFittingPolicyScroll,
	) {
		bnkMngr.banks.Range(func(key any, value any) bool {
			base := filepath.Base(key.(string))
			open := true
			if imgui.BeginTabItemV(base, &open, 0) {
				showHierarchy(key.(string), value.(*bankTab))
				activeTab = value.(*bankTab)
				imgui.EndTabItem()
			}

			if !open {
				activeTab = nil
				closeTab = key.(string)
			}

			return true
		})
		imgui.EndTabBar()
	}
	imgui.End()

	if saveActive {
		saveTab = activeTab
	}

	return activeTab, closeTab, saveTab, saveName
}

func showBankExplorerMenu(bnkMngr *BankManager) (*bankTab, string) {
	var saveTab *bankTab = nil
	saveName := ""

	if !imgui.BeginMenuBar() {
		return saveTab, saveName
	}

	if imgui.BeginMenu("File") {
		if imgui.BeginMenuV("Save", !bnkMngr.writeLock.Load()) {
			bnkMngr.banks.Range(func(key, value any) bool {
				if imgui.MenuItemBool(key.(string)) {
					saveTab = value.(*bankTab)
					saveName = key.(string)
				}
				return true
			})
			imgui.EndMenu()
		}
		imgui.EndMenu()
	}
	imgui.EndMenuBar()

	return saveTab, saveName
}

func showHierarchy(path string, t *bankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	imgui.SetNextItemShortcut(DefaultSearchSC)
	if imgui.InputTextWithHint(
		"Filter by hierarchy object ID", "", &t.idQuery, 0, nil,
	) {
		t.filter()
	}

	if imgui.ComboStrarr(
		"Filter by hierarchy object type",
		&t.typeQuery,
		wwise.HircTypeName,
		int32(len(wwise.HircTypeName)),
	) {
		t.filter()
	}

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}
	if !imgui.BeginTableV(path, 2, 
		imgui.TableFlagsResizable | imgui.TableFlagsReorderable | 
		imgui.TableFlagsRowBg | 
		imgui.TableFlagsBordersH | imgui.TableFlagsBordersV |
		imgui.TableFlagsScrollY,
		imgui.Vec2{X: 0.0, Y: 0.0}, 0,
	) {
		return
	}
	imgui.TableSetupColumn("Hierarchy ID")
	imgui.TableSetupColumn("Hierarchy Type")
	imgui.TableSetupScrollFreeze(0, 1)
	imgui.TableHeadersRow()
	if focusTable {
		imgui.SetKeyboardFocusHere()
	}

	lSelStorage := t.lSelStorage
	hircObjs := t.filtered

	flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d
	msIO := imgui.BeginMultiSelectV(flags, lSelStorage.Size(), int32(len(hircObjs)))
	lSelStorage.ApplyRequests(msIO)

	clipper := imgui.NewListClipper()
	clipper.Begin(int32(len(hircObjs)))
	if msIO.RangeSrcItem() != 1 {
		// Ensure RangeSrc item is not clipped
		clipper.IncludeItemByIndex(int32(msIO.RangeSrcItem()))
	}
	for clipper.Step() {
		for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
			o := hircObjs[n]

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)

			var idS string
			id, err := o.HircID()
			if err != nil {
				idS = "-"
			} else {
				idS = strconv.FormatUint(uint64(id), 10)
			}

			selected := lSelStorage.Contains(imgui.ID(n))
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(n))
			if err != nil {
				imgui.PushIDStr(fmt.Sprintf("UnknownID%d", n))
			}
			imgui.SelectableBoolPtrV(
				idS, &selected, 
				imgui.SelectableFlagsSpanAllColumns | 
				imgui.SelectableFlagsAllowOverlap,
				imgui.Vec2{X: 0, Y: 0},
			)

			imgui.TableSetColumnIndex(1)
			imgui.Text(wwise.HircTypeName[o.HircType()])
			if err != nil {
				imgui.PopID()
			}
		}
	}

	msIO = imgui.EndMultiSelect()
	lSelStorage.ApplyRequests(msIO)

	imgui.EndTable()
}
