package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

func commandPaletteModal(cmdMngr *CmdPaletteMngr) (func(), *bool) {
	done := false
	cmdMngr.selected = 0 
	return func() {
		if useViEnter() {
			done = true
			cmdMngr.filtered[cmdMngr.selected].cmd.action()
			return
		}
		if isUpShortcut() {
			cmdMngr.SetNext(-1)
		}
		if isDownShortcut() {
			cmdMngr.SetNext(1)
		}
		imgui.SetNextItemShortcutV(
			imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF), 0,
		)
		if imgui.InputTextWithHint("Command", "", &cmdMngr.query, 0, nil) {
			cmdMngr.filter()
		}
		clearQuery()
		if !imgui.BeginTableV("CmdPaletteTable",
			1, imgui.TableFlagsRowBg | imgui.TableFlagsScrollY,
			imgui.NewVec2(0.0, 0.0), 0.0,
		) {
			return
		}
		imgui.TableSetupColumn("Command")
		imgui.TableHeadersRow()
		for i, cmd := range cmdMngr.filtered {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			if imgui.SelectableBoolV(
				cmd.cmd.name,
				i == cmdMngr.selected,
				0,
				imgui.NewVec2(0, 0),
			) {
				cmdMngr.selected = i
			}
			focused := imgui.IsItemFocused()
			if focused {
				cmdMngr.selected = int(i)
			}
		}
		imgui.EndTable()
	}, &done
}
