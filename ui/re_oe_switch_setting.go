package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func RenderSwitchSetting(t *bank_explorer.BankTab, o *wwise.SwitchCntr) {
	if imgui.TreeNodeExStr("Switch Setting") {
		imgui.SetNextItemWidth(128)
		imgui.BeginDisabled()
		preview := ""
	if o.GroupType == 0 {
			preview = "Switch Groups"
		} else if o.GroupType == 1 {
			preview = "State Groups"
		}
		if imgui.BeginCombo("Switch Group Type", preview) {
			if imgui.SelectableBool("Switch Groups") {
				o.GroupType = 0
			}
			if imgui.SelectableBool("State Groups") {
				o.GroupType = 1
			}
			imgui.EndCombo()
		}
		imgui.EndDisabled()

		imgui.Text(fmt.Sprintf("Group ID: %d", o.GroupID))
	}
}
