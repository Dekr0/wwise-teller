package ui

import (
	"encoding/binary"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func RenderPositioningParam(t *bank_explorer.BankTab, b *wwise.BaseParameter, v int) {
	if imgui.TreeNodeExStr("Positioning Parameter") {
		prop := &b.PropBundle
		pos := &b.PositioningParam

		parentless := b.DirectParentId == 0
		overrideParent := pos.OverrideParent()

		imgui.BeginDisabledV(parentless)
		if imgui.Checkbox("Override Parent", &overrideParent) {
			pos.SetOverrideParent(overrideParent)
		}
		imgui.EndDisabled()

		listenerRelativeRouting := pos.ListenerRelativeRouting()
		imgui.BeginDisabledV(!parentless && !overrideParent) // (
		if imgui.Checkbox("Listener Relative Routing", &listenerRelativeRouting) {
			pos.SetListenerRelativeRouting(listenerRelativeRouting)
		}

		attenuation := pos.Attenuation()
		imgui.BeginDisabledV(!listenerRelativeRouting)
		if imgui.Checkbox("Attenuation", &attenuation) {
			pos.EnableAttenuation(attenuation)
		}
		if attenuation {
			idx, p := prop.Prop(wwise.TAttenuationID, v)
			var attenuationID uint32
			if idx >= 0 {
				binary.Decode(p.V, wio.ByteOrder, &attenuationID)
				imgui.SetNextItemWidth(128)
				if imgui.BeginCombo("##AttenuationID", strconv.FormatUint(uint64(attenuationID), 10)) {
					imgui.EndCombo()
				}
			}
			if idx >= 0 {
				imgui.SameLine()
				if imgui.ArrowButton("##GoToAttenuationID", imgui.DirRight) {
					t.SetActiveAttenuation(attenuationID)
					DockMngr.SetLayout(dockmanager.AttenuationLayout)
				}
			}
		}
		imgui.EndDisabled()

		imgui.EndDisabled()

		imgui.TreePop()
	}
}
