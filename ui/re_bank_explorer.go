package ui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankExplorerL(bnkMngr *BankManager, saveActive bool, iType int) (
	*bankTab, string, *bankTab, string, int,
) {
	var activeTab *bankTab = nil
	activeName := ""
	closedTab := "" 

	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)

	savedTab, savedName, iType := renderBankExplorerMenu(bnkMngr, iType)

	if imgui.BeginTabBarV(
		"BankExplorerTabBar",
		imgui.TabBarFlagsReorderable | imgui.TabBarFlagsAutoSelectNewTabs | 
		imgui.TabBarFlagsTabListPopupButton | imgui.TabBarFlagsFittingPolicyScroll,
	) {
		bnkMngr.banks.Range(func(key any, value any) bool {
			base := filepath.Base(key.(string))
			open := true
			if imgui.BeginTabItemV(base, &open, 0) {
				renderHircLTable(value.(*bankTab))
				activeTab = value.(*bankTab)
				activeName = key.(string)
				imgui.EndTabItem()
			}

			if !open {
				activeTab = nil
				closedTab = key.(string)
			}

			return true
		})
		imgui.EndTabBar()
	}
	imgui.End()

	if saveActive {
		savedTab = activeTab
		savedName = activeName
	}

	return activeTab, closedTab, savedTab, savedName, iType
}

func renderBankExplorerMenu(bnkMngr *BankManager, itype int) (*bankTab, string, int) {
	var saveTab *bankTab = nil
	saveName := ""

	if !imgui.BeginMenuBar() {
		return saveTab, saveName, itype
	}

	if imgui.BeginMenu("File") {
		if imgui.BeginMenuV("Save", !bnkMngr.writeLock.Load()) {
			bnkMngr.banks.Range(func(key, value any) bool {
				if imgui.MenuItemBool(key.(string)) {
					saveTab = value.(*bankTab)
					saveName = key.(string)
					itype = -1
				}
				return true
			})
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Integration") {
			if imgui.BeginMenuV("Helldivers 2", !bnkMngr.writeLock.Load()) {
				bnkMngr.banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						saveTab = value.(*bankTab)
						saveName = key.(string)
						itype = int(helldivers.IntegrationTypeHelldivers2)
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

	return saveTab, saveName, itype
}

func renderHircLTable(b *bankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	imgui.SetNextItemShortcut(DefaultSearchSC)
	if imgui.InputTextWithHint(
		"Filter by hierarchy object ID", "", &b.idQuery, 0, nil,
	) {
		b.filter()
	}

	if imgui.ComboStrarr(
		"Filter by hierarchy object type",
		&b.typeQuery,
		wwise.HircTypeName,
		int32(len(wwise.HircTypeName)),
	) {
		b.filter()
	}

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}
	if !imgui.BeginTableV("LinearTable", 2, 
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

	storage := b.storage
	hircObjs := b.filtered

	flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d
	msIO := imgui.BeginMultiSelectV(flags, storage.Size(), int32(len(hircObjs)))
	storage.ApplyRequests(msIO)

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

			selected := storage.Contains(imgui.ID(n))
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
	storage.ApplyRequests(msIO)

	imgui.EndTable()
}

func renderHircTree(t *bankTab)  {
	imgui.Begin("Hierarchy View")
	if t == nil {
		imgui.End()
		return
	}
	renderHircTTable(t)
	imgui.End()
}

func renderHircTTable(t *bankTab) {
	if !imgui.BeginTableV("TreeTable", 2, 
		imgui.TableFlagsResizable | imgui.TableFlagsReorderable | 
		imgui.TableFlagsRowBg | 
		imgui.TableFlagsBordersH | imgui.TableFlagsBordersV,
		imgui.Vec2{X: 0.0, Y: 0.0}, 0,
	) {
		return
	}
	imgui.TableSetupColumn("Hierarchy ID")
	imgui.TableSetupColumn("Hierarchy Type")
	imgui.TableSetupScrollFreeze(0, 1)
	imgui.TableHeadersRow()

	hircObjs := t.bank.HIRC().HircObjs

	c := imgui.NewListClipper()
	c.Begin(int32(len(hircObjs)))

	treeIdx := 0 
	drawIdx := int32(0)
	for c.Step() {
		for drawIdx < c.DisplayEnd() && treeIdx < len(hircObjs) {
			renderHircNode(c, &drawIdx, &treeIdx, hircObjs)
		}
	}
	for treeIdx < len(hircObjs) {
		renderHircNode(c, &drawIdx, &treeIdx, hircObjs)
	}

	imgui.EndTable()
}

func renderHircNode(
	c *imgui.ListClipper,
	drawIdx *int32,
	treeIdx *int,
	hircObjs []wwise.HircObj,
) bool {
	o := hircObjs[*treeIdx]
	visible := *drawIdx >= c.DisplayStart() && *drawIdx < c.DisplayEnd()
	*drawIdx += 1
	*treeIdx += 1

	var sid string
	id, err := o.HircID()
	if err != nil {
		sid = fmt.Sprintf("Unknown Id (Tree Index %d)", *treeIdx)
	} else {
		sid = strconv.FormatUint(uint64(id), 10)
	}

	rootless := false
	if o.ParentID() == 0 {
		rootless = true
	}

	if visible {
		imgui.SetNextItemStorageID(imgui.IDInt(int32(*treeIdx)))
		imgui.TableNextRow()
		imgui.TableSetColumnIndex(0)
		open := imgui.TreeNodeExStrV(sid, imgui.TreeNodeFlagsSpanAllColumns)
		imgui.TableSetColumnIndex(1)
		imgui.Text(wwise.HircTypeName[o.HircType()])
		if open {
			for j := 0; j < o.NumChild(); {
				if !renderHircNode(c, drawIdx, treeIdx, hircObjs) {
					j += 1
				}
			}
			imgui.TreePop()
		} else {
			for j := 0; j < o.NumChild(); {
				if !clippedHircNode(treeIdx, hircObjs) {
					j += 1
				}
			}
		}
	} else if o.NumChild() > 0 {
		// clipped
		if imgui.StateStorage().Int(imgui.IDInt(int32(*treeIdx))) != 0 { // open?
			imgui.TreePushStr(sid)
			for j := 0; j < o.NumChild(); {
				if !renderHircNode(c, drawIdx, treeIdx, hircObjs) {
					j += 1
				}
			}
			imgui.TreePop()
		} else {
			for j := 0; j < o.NumChild(); {
				if !clippedHircNode(treeIdx, hircObjs) {
					j += 1
				}
			}
		}
	}
	return rootless
}

func clippedHircNode(treeIdx *int, hircObjs []wwise.HircObj) bool {
	o := hircObjs[*treeIdx]
	*treeIdx += 1
	freeFloat := false
	if o.ParentID() == 0 {
		freeFloat = true
	}
	for j := 0; j < o.NumChild(); {
		if !clippedHircNode(treeIdx, hircObjs) {
			j += 1
		}
	}
	return freeFloat
}
