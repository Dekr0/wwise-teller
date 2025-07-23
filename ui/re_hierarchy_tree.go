// TODO
// - Tree View Keyboard navigation
package ui

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderActorMixerHircTree(t *be.BankTab, open *bool) {
	if !*open {
		return 
	}
	imgui.BeginV("Actor Mixer Hierarchy", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil {
		return
	}
	renderActorMixerHircTreeTable(t)
}

func renderActorMixerHircTreeTable(t *be.BankTab) {
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	if imgui.BeginTableV("ActorMixerHierarchy", 2, flags, DefaultSize, 0) {
		imgui.TableSetupColumn("Hierarchy ID")
		imgui.TableSetupColumn("Hierarchy Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		// Clipper does not play well with Tree Node :(
		root := t.Bank.HIRC().ActorMixerRoots
		for i := range root {
			renderActorMixerHircNode(t, root[i])
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
	selected = t.ActorMixerViewer.Selected(id)

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)

	flags := DefaultTreeFlags
	if selected {
		flags |= imgui.TreeNodeFlagsSelected
	}

	imgui.SetNextItemOpen(node.Open)
	open := imgui.TreeNodeExStrV(sid, flags)
	node.Open = open 
	if imgui.IsItemClickedV(imgui.MouseButtonLeft) {
		if !imgui.CurrentIO().KeyCtrl() {
			t.ActorMixerViewer.Clear()
		}
		t.ActorMixerViewer.SetSelected(id, true)
	}

	if imgui.BeginPopupContextItem() {
		renderActorMixerHircCtx(t, node, o, id)
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
	if open {
		for i := range node.Leafs {
			renderActorMixerHircNode(t, node.Leafs[i])
		}
		imgui.TreePop()
	}
}

func renderActorMixerHircCtx(
	t *be.BankTab,
	node *wwise.ActorMixerHircNode,
	o wwise.HircObj,
	id uint32,
) {
	Disabled(!GlobalCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy ID") {
			clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(id), 10)))
		}
	})

	switch sound := node.Obj.(type) {
	case *wwise.Sound:
		Disabled(!GlobalCtx.CopyEnable, func() {
			if imgui.SelectableBool("Copy Source ID") {
				clipboard.Write(clipboard.FmtText, []byte(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10)))
			}
		})
	}

	if len(node.Leafs) <= 0 {
		return
	}

	Disabled(!GlobalCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy Leafs' IDs") {
			l := len(node.Leafs)
			var builder strings.Builder
			for i := range node.Leafs {
				id, err := node.Leafs[l - i - 1].Obj.HircID()
				if err != nil {
					panic(err)
				}
				if _, err := builder.WriteString(fmt.Sprintf("%d\n", id)); err != nil {
					slog.Error(fmt.Sprintf("Failed to copy leafs IDs of actor mixer hierarchy object %d", id), "error", err)
				}
			}
			clipboard.Write(clipboard.FmtText, []byte(builder.String()))
		}
	})

	Disabled(!GlobalCtx.CopyEnable, func() {
		if imgui.SelectableBool("Copy Leafs' Source IDs") {
			l := len(node.Leafs)
			var builder strings.Builder
			for i := range node.Leafs {
				switch sound := node.Leafs[l - i - 1].Obj.(type) {
				case *wwise.Sound:
					if _, err := builder.WriteString(fmt.Sprintf("%d\n", sound.BankSourceData.SourceID)); err != nil {
						slog.Error(fmt.Sprintf("Failed to copy leafs IDs of actor mixer hierarchy object %d", id), "error", err)
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

func renderMusicHircTree(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Music Hierarchy", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	imgui.Text("Under Construction")
}
