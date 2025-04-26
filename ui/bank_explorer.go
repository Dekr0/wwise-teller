package ui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
)

func showBankExplorer(bnkMngr *bankManager) (*bankTab, string, *bankTab, string) {
	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)

	var activeTab *bankTab = nil
	var closeTab string = "" 
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

			// Active tab can be closed
			if !open {
				activeTab = nil
				closeTab = key.(string)
			}

			return true
		})

		imgui.EndTabBar()
	}

	imgui.End()

	return activeTab, closeTab, saveTab, saveName
}

func showBankExplorerMenu(bnkMngr *bankManager) (*bankTab, string) {
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

// Display hierarchy objects in either linear fashion or a tree fashion.
// Return a list of hierarchy object being selected in this **active bank tab**
func showHierarchy(path string, t *bankTab) {
	// if imgui.Button("Undo") {
	// 	t.undoList.Undo()
	// }
	// imgui.SameLine()
	// if imgui.Button("Redo") {
	// 	t.undoList.Redo()
	// }

	imgui.SeparatorText("Filter")
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

	lSelStorage := t.lSelStorage
	hircObjs := t.filtered

	// BoxSelect1d ensure a whole row is selected in the table
	flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect1d
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
