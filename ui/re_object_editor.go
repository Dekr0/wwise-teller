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
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderObjEditorActorMixer(m *be.BankManager, t *be.BankTab, init *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Object Editor (Actor Mixer)", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.SounBankLock.Load() {
		return
	}

	// Display loading screen for write lock

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		for _, h := range t.Bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if !wwise.ActorMixerHircType(h) {
				continue
			}
			if t.ActorMixerViewer.Selected(id) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderActorMixerTab(m, t, init, h)
		}
		imgui.EndTabBar()
	}
}

func renderActorMixerTab(m *be.BankManager, t *be.BankTab, init *be.BankTab, h wwise.HircObj) {
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
			renderActorMixer(m, t, init, ah)
		case *wwise.LayerCntr:
			renderLayerCntr(m, t, init, ah)
		case *wwise.RanSeqCntr:
			renderRanSeqCntr(m, t, init, ah)
		case *wwise.SwitchCntr:
			renderSwitchCntr(m, t, init, ah)
		case *wwise.Sound:
			renderSound(m, t, init, ah)
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

func renderObjEditorMusic(m *be.BankManager, t *be.BankTab, init *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Object Editor (Music)", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.SounBankLock.Load() {
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		for _, h := range t.Bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if t.MusicHircViewer.LinearStorage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderMusicTab(m, t, init, h)
		}
		imgui.EndTabBar()
	}
}

func renderMusicTab(m *be.BankManager, t *be.BankTab, init *be.BankTab, h wwise.HircObj) {
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
			renderMusicTrack(m, t, init, mh)
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

func renderActorMixer(m *be.BankManager, t *be.BankTab, init *be.BankTab, o *wwise.ActorMixer) {
	renderBaseParam(m, t, init, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
}

func renderLayerCntr(m *be.BankManager, t *be.BankTab, init *be.BankTab, o *wwise.LayerCntr) {
	renderBaseParam(m, t, init, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
	renderLayer(o)
}

func renderRanSeqCntr(m *be.BankManager, t *be.BankTab, init *be.BankTab, o *wwise.RanSeqCntr) {
	renderBaseParam(m, t, init, o)
	renderContainer(t, o.Id, &o.Container, wwise.ActorMixerHircType(o))
	renderRanSeqPlayList(t, o)
}

func renderSwitchCntr(m *be.BankManager, t *be.BankTab, init *be.BankTab, o *wwise.SwitchCntr) {
	renderBaseParam(m, t, init, o)
}

func renderSound(m *be.BankManager, t *be.BankTab, init *be.BankTab, o *wwise.Sound) {
	renderBankSourceData(t, o)
	renderBaseParam(m, t, init, o)
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
		imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
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
