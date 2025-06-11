// TODO
// - Clean up filter for hierarchy object listing
//   - Action
//   - Attenuation
//   - Event
//   - State
//   - ...
//
// - Tree View Keyboard navigation
// - Tree Render does not work when there's one not expanded root entry in the table
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

func renderBankExplorer(bnkMngr *BankManager, saveActive bool, iType int) (
	string, *BankTab, string, int,
) {
	closedPath := "" 

	imgui.BeginV("Bank Explorer", nil, imgui.WindowFlagsMenuBar)

	savedTab, savedName, iType := renderBankExplorerMenu(bnkMngr, iType)
	if imgui.BeginTabBarV("BankExplorerTabBar", DefaultTabFlags) {
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

func renderBankExplorerTab(path string, t *BankTab) {
	imgui.Text("Sound bank: " + path)
	if imgui.BeginTabBar("SubBankExplorerTabBar") {
		if imgui.BeginTabItem("Hierarchy Listing") {
			renderActorMixerHircTable(t)
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

func renderActorMixerHircTable(b *BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar("By ID", imgui.DataTypeU32, uintptr(utils.Ptr(&b.HircFilter.Id))) {
		b.FilterHircs()
	}

	imgui.BeginDisabledV(b.HircFilter.Type != wwise.HircTypeAll && b.HircFilter.Type != wwise.HircTypeSound)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar("By source ID", imgui.DataTypeU32, uintptr(utils.Ptr(&b.HircFilter.Sid))) {
		b.FilterHircs()
	}
	imgui.EndDisabled()

	typeFilter := int32(b.HircFilter.Type)
	imgui.SetNextItemWidth(256)
	if imgui.ComboStrarr("By hierarchy object type", &typeFilter, wwise.HircTypeName, int32(len(wwise.HircTypeName)),
	) {
		b.HircFilter.Type = wwise.HircType(typeFilter)
		b.FilterHircs()
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

		storage := b.LinearStorage

		msIO := imgui.BeginMultiSelectV(DefaultMultiSelectFlags, storage.Size(), int32(len(b.HircFilter.HircObjs)))
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

func renderActorMixerHircTree(t *BankTab)  {
	imgui.Begin("Actor Mixer Hierarchy")
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil {
		imgui.End()
		return
	}
	renderActorMixerHircTreeTable(t)
	imgui.End()
}

func renderActorMixerHircTreeTable(t *BankTab) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("TreeTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		// Clipper does not play well with Tree Node :(
		for _, root := range t.Bank.HIRC().ActorMixerRoots {
			renderActorMixerHircNode(t, &root)
		}
		imgui.EndTable()
	}
}

func renderActorMixerHircNode(t *BankTab, node *wwise.ActorMixerHircNode) {
	o := node.Obj

	var sid string
	selected := false
	id, err := o.HircID()
	if err != nil { panic("Panic Trap") }

	sid = strconv.FormatUint(uint64(id), 10)
	selected = t.LinearStorage.Contains(imgui.ID(id))

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)

	flags := imgui.TreeNodeFlagsSpanAllColumns | imgui.TreeNodeFlagsOpenOnDoubleClick
	if selected {
		flags |= imgui.TreeNodeFlagsSelected
	}

	open := imgui.TreeNodeExStrV(sid, flags)
	if imgui.IsItemClicked() {
		if !imgui.CurrentIO().KeyCtrl() {
			t.LinearStorage.Clear()
		}
		t.LinearStorage.SetItemSelected(imgui.ID(id), true)
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
	if open {
		for _, leaf := range node.Leafs {
			renderActorMixerHircNode(t, &leaf)
		}
		imgui.TreePop()
	}
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
