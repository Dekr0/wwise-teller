// TODO
// - Tree View Keyboard navigation
package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderActorMixerHircTree(t *be.BankTab)  {
	imgui.Begin("Actor Mixer Hierarchy")
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil {
		imgui.End()
		return
	}
	renderActorMixerHircTreeTable(t)
	imgui.End()
}

func renderActorMixerHircTreeTable(t *be.BankTab) {
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

func renderActorMixerHircNode(t *be.BankTab, node *wwise.ActorMixerHircNode) {
	o := node.Obj

	var sid string
	selected := false
	id, err := o.HircID()
	if err != nil { panic("Panic Trap") }

	sid = strconv.FormatUint(uint64(id), 10)
	selected = t.ActorMixerViewer.LinearStorage.Contains(imgui.ID(id))

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)

	flags := imgui.TreeNodeFlagsSpanAllColumns | imgui.TreeNodeFlagsOpenOnDoubleClick
	if selected {
		flags |= imgui.TreeNodeFlagsSelected
	}

	open := imgui.TreeNodeExStrV(sid, flags)
	if imgui.IsItemClicked() {
		if !imgui.CurrentIO().KeyCtrl() {
			t.ActorMixerViewer.LinearStorage.Clear()
		}
		t.ActorMixerViewer.LinearStorage.SetItemSelected(imgui.ID(id), true)
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

func renderMusicHircTree(t *be.BankTab) {
	imgui.Begin("Music Hierarchy")
	imgui.End()
}
