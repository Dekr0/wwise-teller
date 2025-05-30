package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderObjEditor(t *bankTab) {
	imgui.Begin("Object Editor")

	if t == nil {
		imgui.End()
		return
	}
	if t.writeLock.Load() {
		imgui.End()
		return
	}

	useViUp()
	useViDown()

	if imgui.BeginTabBarV(
		"ObjectEditorTabBar",
		imgui.TabBarFlagsReorderable |
		imgui.TabBarFlagsAutoSelectNewTabs |
		imgui.TabBarFlagsTabListPopupButton | imgui.TabBarFlagsFittingPolicyScroll,
	) {
		s := []wwise.HircObj{}
		for _, h := range t.bank.HIRC().HircObjs {
			id, err := h.HircID()
			if err != nil {
				continue
			}
			if t.storage.Contains(imgui.ID(id)) {
				s = append(s, h)
			}
		}

		for i, h := range s {
			renderHircTab(t, i, h)
		}

		imgui.EndTabBar()
	}

	imgui.End()
}

func renderHircTab(t *bankTab, i int, h wwise.HircObj) {
	var label string
	switch h.(type) {
	case *wwise.Unknown:
		label = fmt.Sprintf("Unknown Object %d", i)
	default:
		id, _ := h.HircID()
		label = fmt.Sprintf("%s %d", wwise.HircTypeName[h.HircType()], id)
	}

	if imgui.BeginTabItem(label) {
		if t.activeHirc != h {
			t.cntrStorage.Clear()
			t.playListStorage.Clear()
		}
		t.activeHirc = h
		switch h.(type) {
		case *wwise.ActorMixer:
			renderActorMixer(t, h.(*wwise.ActorMixer))
		case *wwise.LayerCntr:
			renderLayerCntr(t, h.(*wwise.LayerCntr))
		case *wwise.MusicTrack:
			renderMusicTrack(t, h.(*wwise.MusicTrack))
		case *wwise.RanSeqCntr:
			renderRanSeqCntr(t, h.(*wwise.RanSeqCntr))
		case *wwise.SwitchCntr:
			renderSwitchCntr(t, h.(*wwise.SwitchCntr))
		case *wwise.Sound:
			renderSound(t, h.(*wwise.Sound))
		case *wwise.Unknown:
			renderUnknown(h.(*wwise.Unknown))
		}
		imgui.EndTabItem()
	}
}

func renderActorMixer(t *bankTab, o *wwise.ActorMixer) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, o.Container)
}

func renderLayerCntr(t *bankTab, o *wwise.LayerCntr) {
	renderBaseParam(t, o)
}

func renderRanSeqCntr(t *bankTab, o *wwise.RanSeqCntr) {
	renderBaseParam(t, o)
	renderContainer(t, o.Id, o.Container)
	renderRanSeqPlayList(t, o)
}

func renderSwitchCntr(t *bankTab, o *wwise.SwitchCntr) {
	renderBaseParam(t, o)
}

func renderSound(t *bankTab, o *wwise.Sound) {
	renderBankSourceData(t, o)
	renderBaseParam(t, o)
}

func renderUnknown(o *wwise.Unknown) {
	imgui.Text(
		fmt.Sprintf(
			"Support for hierarchy object type %s is still under construction.",
			wwise.HircTypeName[o.HircType()],
		),
	)
}

func renderContainer(t *bankTab, id uint32, cntr *wwise.Container) {
	if imgui.TreeNodeExStr("Container") {
		if imgui.Button("Add New Children") {
		}

		const flags = DefaultTableFlags
		outerSize := imgui.NewVec2(0.0, 0.0)
		if imgui.BeginTableV("CntrTable", 2, flags, outerSize, 0) {
			imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumn("Children Hierarchy ID")
			imgui.TableHeadersRow()

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
			}

			imgui.EndTable()

			if deleteChild != nil {
				deleteChild()
			}
		}
		imgui.TreePop()
	}
}
