package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

func commandPaletteModal(cmdMngr *CmdPaletteMngr) (func(), *bool) {
	done := false
	cmdMngr.Selected = 0 
	return func() {
		if useViEnter() {
			done = true
			cmdMngr.Filtered[cmdMngr.Selected].Cmd.Action()
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
		if imgui.InputTextWithHint("Command", "", &cmdMngr.Query, 0, nil) {
			cmdMngr.Filter()
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
		for i, cmd := range cmdMngr.Filtered {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			if imgui.SelectableBoolV(
				cmd.Cmd.Name,
				i == cmdMngr.Selected,
				0,
				imgui.NewVec2(0, 0),
			) {
				cmdMngr.Selected = i
			}
			focused := imgui.IsItemFocused()
			if focused {
				cmdMngr.Selected = int(i)
			}
		}
		imgui.EndTable()
	}, &done
}

func pushCommandPaletteModal(cmdPaletteMngr *CmdPaletteMngr) {
	renderF, done := commandPaletteModal(cmdPaletteMngr)
	GCtx.ModalQ.QModal(done, 0, "Command Palette", renderF, nil)
}
