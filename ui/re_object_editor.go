// TODO
// - Add New Children
package ui

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderObjEditorActorMixer(open *bool) {
	if !*open {
		return
	}
	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.ObjectEditorActorMixerTag], open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}

	activeBank, valid := BnkMngr.ActiveBankV()
	if !valid || activeBank.SounBankLock.Load() {
		return
	}
	// Display loading screen for write lock

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		for _, h := range activeBank.Bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if !wwise.ActorMixerHircType(h) {
				continue
			}
			if activeBank.ActorMixerViewer.Selected(id) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderActorMixerTab(activeBank, h)
		}
		imgui.EndTabBar()
	}
}

func renderActorMixerTab(t *be.BankTab, h wwise.HircObj) {
	viewer := &t.ActorMixerViewer
	id, err := h.HircID()
	if err != nil {
		panic(err)
	}
	label := fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
	open := true
	if imgui.BeginTabItemV(label, &open, imgui.TabItemFlagsNone) {
		if viewer.ActiveHirc != h {
			viewer.CntrStorage.Clear()
			viewer.RanSeqPlaylistStorage.Clear()
		}
		viewer.ActiveHirc = h
		switch ah := h.(type) {
		case *wwise.ActorMixer:
			renderActorMixer(t, ah)
		case *wwise.LayerCntr:
			renderLayerCntr(t, ah)
		case *wwise.RanSeqCntr:
			renderRanSeqCntr(t, ah)
		case *wwise.SwitchCntr:
			renderSwitchCntr(t, ah)
		case *wwise.Sound:
			renderSound(t, ah)
		default:
			panic("Panic Trap")
		}
		imgui.EndTabItem()
	}
	if !open {
		viewer.SetSelected(id, false)
		viewer.CntrStorage.Clear()
		viewer.RanSeqPlaylistStorage.Clear()
	}
}

func renderObjEditorMusic(open *bool) {
	if !*open {
		return
	}
	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.ObjectEditorMusicTag], open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}

	activeBank, valid := BnkMngr.ActiveBankV()
	if !valid || activeBank.SounBankLock.Load() {
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		hirc := activeBank.Bank.HIRC()
		for _, h := range hirc.HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if activeBank.MusicHircViewer.LinearStorage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderMusicTab(activeBank, h)
		}
		imgui.EndTabBar()
	}
}

func renderMusicTab(t *be.BankTab, h wwise.HircObj) {
	viewer := &t.MusicHircViewer
	id, err := h.HircID()
	if err != nil {
		panic(err)
	}
	label := fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
	open := true
	if imgui.BeginTabItemV(label, &open, imgui.TabItemFlagsNone) {
		if viewer.ActiveMusicHirc != h {
			viewer.CntrStorage.Clear()
		}
		viewer.ActiveMusicHirc = h
		switch mh := h.(type) {
		case *wwise.MusicTrack:
			renderMusicTrack(t, mh)
		case *wwise.MusicSegment:
			renderMusicSegment(t, mh)
		case *wwise.MusicRanSeqCntr:
			renderMusicRanSeqCntr(t, mh)
		case *wwise.MusicSwitchCntr:
			renderMusicSwitchCntr(t, mh)
		default:
			panic("Panic Trap")
		}
		imgui.EndTabItem()
	}
	if !open {
		viewer.LinearStorage.SetItemSelected(imgui.ID(id), false)
		viewer.CntrStorage.Clear()
	}
}

func renderActorMixer(t *be.BankTab, o *wwise.ActorMixer) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
}

func renderLayerCntr(t *be.BankTab, o *wwise.LayerCntr) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
	renderLayer(o)
}

func renderRanSeqCntr(t *be.BankTab, o *wwise.RanSeqCntr) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
	renderRanSeqPlayList(t, o)
}

func renderSwitchCntr(t *be.BankTab, o *wwise.SwitchCntr) {
	renderBaseParam(t, o)
}

func renderSound(t *be.BankTab, o *wwise.Sound) {
	renderBankSourceData(t, o)
	renderBaseParam(t, o)
}

func renderMusicSegment(t *be.BankTab, o *wwise.MusicSegment) {
	imgui.Text("Under construction")
}

func renderMusicSwitchCntr(t *be.BankTab, o *wwise.MusicSwitchCntr) {
	imgui.Text("Under construction")
}

func renderMusicRanSeqCntr(t *be.BankTab, o *wwise.MusicRanSeqCntr) {
	imgui.Text("Under construction")
}

func renderContainer(t *be.BankTab, id uint32, cntr *wwise.Container, actorMixer bool) {
	if imgui.TreeNodeExStr("Container") {
		imgui.BeginDisabled()
		imgui.Button("Add New Children")
		imgui.EndDisabled()

		imgui.SameLine()
		imgui.BeginDisabledV(!GCtx.CopyEnable)
		if imgui.Button("Copy IDs") {
			var builder strings.Builder
			var err error
			for _, child := range cntr.Children {
				if _, err = builder.WriteString(strconv.FormatUint(uint64(child), 10)+"\n"); err != nil {
					slog.Error("Failed to copy children IDs", "error", err)
					break
				}
			}
			clipboard.Write(clipboard.FmtText, []byte(builder.String()))
		}
		imgui.EndDisabled()

		const flags = DefaultTableFlags
		outerSize := imgui.NewVec2(0.0, 0.0)
		if imgui.BeginTableV("CntrTable", 5, flags, outerSize, 0) {
			imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumn("ID")
			imgui.TableSetupColumn("Type")
			imgui.TableSetupColumn("Source ID")
			imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()

			hirc := t.Bank.HIRC()
			var deleteChild func() = nil
			for _, i := range cntr.Children {
				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)
				imgui.SetNextItemWidth(40)

				imgui.PushIDStr(fmt.Sprintf("DelChild%d", i))
				if imgui.Button("X") {
					deleteChild = bindRemoveRoot(t, i, id)
				}
				imgui.PopID()

				imgui.TableSetColumnIndex(1)
				imgui.SetNextItemWidth(-1)

				imgui.Text(strconv.FormatUint(uint64(i), 10))

				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(-1)
				value, ok := hirc.ActorMixerHirc.Load(i)
				if !ok {
					imgui.Text("-")
				} else {
					obj := value.(wwise.HircObj)
					imgui.Text(wwise.HircTypeName[obj.HircType()])
					if obj.HircType() == wwise.HircTypeSound {
						imgui.TableSetColumnIndex(3)
						imgui.Text(strconv.FormatUint(uint64(obj.(*wwise.Sound).BankSourceData.SourceID), 10))
					}
				}

				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(56)
				imgui.BeginDisabledV(!ok)
				if imgui.ArrowButton("CntrGoTo" + strconv.FormatUint(uint64(i), 10), imgui.DirRight) {
					if actorMixer {
						t.ActorMixerViewer.SetSelected(i, true)
					} else {
						t.MusicHircViewer.LinearStorage.SetItemSelected(imgui.ID(i), true)
					}
				}
				imgui.EndDisabled()
			}

			imgui.EndTable()

			if deleteChild != nil {
				deleteChild()
			}
		}
		imgui.TreePop()
	}
}
