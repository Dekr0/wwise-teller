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

	flags := imgui.TabBarFlagsReorderable        | 
	         imgui.TabBarFlagsAutoSelectNewTabs  | 
		     imgui.TabBarFlagsTabListPopupButton | 
	         imgui.TabBarFlagsFittingPolicyScroll
	if imgui.BeginTabBarV("BankExplorerTabBar", flags) {
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

	imgui.BeginDisabledV(b.typeQuery != 0 && b.typeQuery != int32(wwise.HircTypeSound))
	if imgui.InputTextWithHint("Filter by source ID", "", &b.sidQuery, 0, nil) {
		b.filter()
	}
	imgui.EndDisabled()

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

	tableFlags := imgui.TableFlagsResizable   |
			      imgui.TableFlagsReorderable |
		          imgui.TableFlagsRowBg       |
		          imgui.TableFlagsBordersH    |
				  imgui.TableFlagsBordersV    |
		          imgui.TableFlagsScrollY
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("LinearTable", 2, tableFlags, outerSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		storage := b.storage
		hircObjs := b.filtered

		flags := imgui.MultiSelectFlagsClearOnEscape | 
		         imgui.MultiSelectFlagsBoxSelect2d
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
	tableFlags := imgui.TableFlagsResizable   | 
				  imgui.TableFlagsReorderable | 
			      imgui.TableFlagsRowBg       | 
		          imgui.TableFlagsBordersH    | 
				  imgui.TableFlagsBordersV
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("TreeTable", 2, tableFlags, outerSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		hirc := t.bank.HIRC()
		hircObjs := hirc.HircObjs

		c := imgui.NewListClipper()
		c.Begin(int32(len(hircObjs)) - int32(hirc.ActionCount.Load()) - int32(hirc.EventCount.Load()))

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
}

func renderHircNode(
	c *imgui.ListClipper,
	drawIdx *int32,
	treeIdx *int,
	hircObjs []wwise.HircObj,
) bool {
	l := len(hircObjs)
	o := hircObjs[l - *treeIdx - 1]
	if o.HircType() == wwise.HircTypeAction || o.HircType() == wwise.HircTypeEvent {
		*treeIdx += 1
		return true
	}
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
		st := wwise.HircTypeName[o.HircType()]
		if o.HircType() == wwise.HircTypeSound {
			st = fmt.Sprintf(
				"%s (Audio Source %d)",
				st, o.(*wwise.Sound).BankSourceData.SourceID,
			)
		}
		imgui.Text(st)
		if open {
			for j := 0; j < o.NumLeaf(); {
				if !renderHircNode(c, drawIdx, treeIdx, hircObjs) {
					j += 1
				}
			}
			imgui.TreePop()
		} else {
			for j := 0; j < o.NumLeaf(); {
				if !clippedHircNode(treeIdx, hircObjs) {
					j += 1
				}
			}
		}
	} else if o.NumLeaf() > 0 {
		// clipped
		if imgui.StateStorage().Int(imgui.IDInt(int32(*treeIdx))) != 0 { // open?
			imgui.TreePushStr(sid)
			for j := 0; j < o.NumLeaf(); {
				if !renderHircNode(c, drawIdx, treeIdx, hircObjs) {
					j += 1
				}
			}
			imgui.TreePop()
		} else {
			for j := 0; j < o.NumLeaf(); {
				if !clippedHircNode(treeIdx, hircObjs) {
					j += 1
				}
			}
		}
	}
	return rootless
}

func clippedHircNode(treeIdx *int, hircObjs []wwise.HircObj) bool {
	l := len(hircObjs)
	o := hircObjs[l - *treeIdx - 1]
	*treeIdx += 1
	freeFloat := false
	if o.ParentID() == 0 {
		freeFloat = true
	}
	for j := 0; j < o.NumLeaf(); {
		if !clippedHircNode(treeIdx, hircObjs) {
			j += 1
		}
	}
	return freeFloat
}
