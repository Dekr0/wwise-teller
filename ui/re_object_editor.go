// TODO
// - Add New Children
package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderObjEditorActorMixer(t *BankTab, init *wwise.Bank) {
	imgui.Begin("Object Editor (Actor Mixer)")

	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		imgui.End()
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		for _, h := range t.Bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil { panic(err) }
			if t.ActorMixerViewer.LinearStorage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderActorMixerTab(t, init, h)
		}
		imgui.EndTabBar()
	}

	imgui.End()
}

func renderActorMixerTab(t *BankTab, init *wwise.Bank, h wwise.HircObj) {
	viewer := &t.ActorMixerViewer
	id, err := h.HircID()
	if err != nil {
		panic(err)
	}
	label := fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
	open := true
	if imgui.BeginTabItemV(label, &open, imgui.TabItemFlagsNone) {
		if viewer.ActiveActorMixerHirc != h {
			viewer.CntrStorage.Clear()
			viewer.RanSeqPlaylistStorage.Clear()
		}
		viewer.ActiveActorMixerHirc = h
		switch ah := h.(type) {
		case *wwise.ActorMixer:
			renderActorMixer(t, init, ah)
		case *wwise.LayerCntr:
			renderLayerCntr(t, init, ah)
		case *wwise.RanSeqCntr:
			renderRanSeqCntr(t, init, ah)
		case *wwise.SwitchCntr:
			renderSwitchCntr(t, init, ah)
		case *wwise.Sound:
			renderSound(t, init, ah)
		default:
			panic("Panic Trap")
		}
		imgui.EndTabItem()
	}
	if !open {
		viewer.LinearStorage.SetItemSelected(imgui.ID(id), false)
		viewer.CntrStorage.Clear()
		viewer.RanSeqPlaylistStorage.Clear()
	}
}

func renderObjEditorMusic(t *BankTab, init *wwise.Bank) {
	imgui.Begin("Object Editor (Music)")

	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		imgui.End()
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV("ObjectEditorTabBar", DefaultTabFlags) {
		s := []wwise.HircObj{}
		for _, h := range t.Bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil { panic(err) }
			if t.MusicHircViewer.LinearStorage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}
		for _, h := range s {
			renderMusicTab(t, init, h)
		}
		imgui.EndTabBar()
	}
	imgui.End()
}

func renderMusicTab(t *BankTab, init *wwise.Bank, h wwise.HircObj) {
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
			renderMusicTrack(t, init, mh)
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

func renderActorMixer(t *BankTab, init *wwise.Bank, o *wwise.ActorMixer) {
	renderBaseParam(t, init, o)
	renderContainer(t, o.Id, o.Container, wwise.ActorMixerHircType(o))
}

func renderLayerCntr(t *BankTab, init *wwise.Bank, o *wwise.LayerCntr) {
	renderBaseParam(t, init, o)
}

func renderRanSeqCntr(t *BankTab, init *wwise.Bank, o *wwise.RanSeqCntr) {
	renderBaseParam(t, init, o)
	renderContainer(t, o.Id, o.Container, wwise.ActorMixerHircType(o))
	renderRanSeqPlayList(t, o)
}

func renderSwitchCntr(t *BankTab, init *wwise.Bank, o *wwise.SwitchCntr) {
	renderBaseParam(t, init, o)
}

func renderSound(t *BankTab, init *wwise.Bank, o *wwise.Sound) {
	renderBankSourceData(t, o)
	renderBaseParam(t, init, o)
}

func renderMusicSegment(t *BankTab, o *wwise.MusicSegment) {
	imgui.Text("Under construction")
}

func renderMusicSwitchCntr(t *BankTab, o *wwise.MusicSwitchCntr) {
	imgui.Text("Under construction")
}

func renderMusicRanSeqCntr(t *BankTab, o *wwise.MusicRanSeqCntr) {
	imgui.Text("Under construction")
}

func renderContainer(t *BankTab, id uint32, cntr *wwise.Container, actorMixer bool) {
	if imgui.TreeNodeExStr("Container") {
		imgui.BeginDisabled()
		imgui.Button("Add New Children")
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
						t.ActorMixerViewer.LinearStorage.SetItemSelected(imgui.ID(i), true)
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
