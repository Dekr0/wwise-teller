package ui

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankExplorerL(bnkMngr *BankManager, saveActive bool, iType int) (
	string, *BankTab, string, int,
) {
	closedPath := "" 

	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)

	savedTab, savedName, iType := renderBankExplorerMenu(bnkMngr, iType)

	flags := imgui.TabBarFlagsReorderable        | 
	         imgui.TabBarFlagsAutoSelectNewTabs  | 
		     imgui.TabBarFlagsTabListPopupButton | 
	         imgui.TabBarFlagsFittingPolicyScroll
	if imgui.BeginTabBarV("BankExplorerTabBar", flags) {
		bnkMngr.Banks.Range(func(key any, value any) bool {
			open := true
			path := key.(string)
			tab := value.(*BankTab)

			imgui.PushIDStr(path)
			if imgui.BeginTabItemV(filepath.Base(path), &open, 0) {
				renderHircLTable(path, value.(*BankTab))
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
				closedPath = path
			}

			return true
		})
		imgui.EndTabBar()
	}
	imgui.End()


	if saveActive {
		savedTab = bnkMngr.ActiveBank
		savedName = bnkMngr.ActivePath
	}

	return closedPath, savedTab, savedName, iType
}

func renderBankExplorerMenu(bnkMngr *BankManager, itype int) (*BankTab, string, int) {
	var saveTab *BankTab = nil
	saveName := ""

	if !imgui.BeginMenuBar() {
		return saveTab, saveName, itype
	}

	if imgui.BeginMenu("File") {
		if imgui.BeginMenuV("Save", !bnkMngr.WriteLock.Load()) {
			bnkMngr.Banks.Range(func(key, value any) bool {
				if imgui.MenuItemBool(key.(string)) {
					saveTab = value.(*BankTab)
					saveName = key.(string)
					itype = -1
				}
				return true
			})
			imgui.EndMenu()
		}
		if imgui.BeginMenu("Integration") {
			if imgui.BeginMenuV("Helldivers 2", !bnkMngr.WriteLock.Load()) {
				bnkMngr.Banks.Range(func(key, value any) bool {
					if imgui.MenuItemBool(key.(string)) {
						saveTab = value.(*BankTab)
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

func renderHircLTable(path string, b *BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.Text("Sound bank: " + path)

	imgui.SeparatorText("Filter")

	imgui.SetNextItemShortcut(DefaultSearchSC)
	if imgui.InputScalar("Filter by hierarchy object ID", imgui.DataTypeU32, uintptr(utils.Ptr(&b.HircFilter.Id))) {
		b.FilterHircs()
	}

	imgui.BeginDisabledV(b.HircFilter.Type != wwise.HircTypeAll && b.HircFilter.Type != wwise.HircTypeSound)
	if imgui.InputScalar("Filter by source ID", imgui.DataTypeU32, uintptr(utils.Ptr(&b.HircFilter.Sid))) {
		b.FilterHircs()
	}
	imgui.EndDisabled()

	typeFilter := int32(b.HircFilter.Type)
	if imgui.ComboStrarr("Filter by hierarchy object type", &typeFilter, wwise.HircTypeName, int32(len(wwise.HircTypeName)),
	) {
		b.HircFilter.Type = wwise.HircType(typeFilter)
		b.FilterHircs()
	}

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	flags := DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("LinearTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		storage := b.LinearStorage

		flags := imgui.MultiSelectFlagsClearOnEscape | 
		         imgui.MultiSelectFlagsBoxSelect2d
		msIO := imgui.BeginMultiSelectV(flags, storage.Size(), int32(len(b.HircFilter.HircObjs)))
		storage.ApplyRequests(msIO)

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(b.HircFilter.HircObjs)))
		if msIO.RangeSrcItem() != 1 {
			// Ensure RangeSrc item is not clipped
			clipper.IncludeItemByIndex(int32(msIO.RangeSrcItem()))
		}
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				o := b.HircFilter.HircObjs[n]

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

func renderHircTree(t *BankTab)  {
	imgui.Begin("Hierarchy View")
	if t == nil || t.Bank == nil || t.Bank.HIRC() != nil {
		imgui.End()
		return
	}
	renderHircTTable(t)
	imgui.End()
}

func renderHircTTable(t *BankTab) {
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

		hirc := t.Bank.HIRC()
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
