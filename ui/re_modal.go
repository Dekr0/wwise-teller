package ui

import (
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/modal"
)

func renderModal(m *modal.ModalQ) {
	if len(m.Modals) <= 0 {
		return
	}
	top := m.Modals[len(m.Modals)-1]
	if *top.Done {
		imgui.CloseCurrentPopup()
		m.Modals = m.Modals[:len(m.Modals)-1]
		if top.OnClose != nil {
			top.OnClose()
		}
		return
	}
	if !imgui.IsPopupOpenStr(top.Name) {
		imgui.OpenPopupStr(top.Name)
		imgui.SetNextWindowSize(imgui.NewVec2(640, 640))
	}
	center := imgui.MainViewport().Center()
	imgui.SetNextWindowPosV(center, imgui.CondAppearing, imgui.NewVec2(0.5, 0.5))
	if imgui.BeginPopupModalV(top.Name, nil, top.Flag) {
		if imgui.Shortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyQ)) {
			*top.Done = true
		}
		top.Loop()
		imgui.EndPopup()
	}
}

func pushSetHomeModal() {
	onSave := func(path string) {
		if err := GCtx.Config.SetHome(path); err != nil {
			slog.Error(
				"Failed to set initial directory for file " +
				"explorer",
				"error", err,
			)
		}
	}
	renderF, done, err := saveFileDialogFunc(onSave, GCtx.Config.Home)
	if err != nil {
		slog.Error(
			"Failed to create save file dialog for setting initial" + 
			" directory for file explorer",
			"error", err,
		)
	} else {
		Modal(done, 0, "Set starting directory for file explorer", renderF, nil)
	}
}

func PushSimpleTextModal(title string, label string, confirm string, c func(string)) {
	done := false
	text := ""
	f := func() {
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF))
		imgui.InputTextWithHint(label, "", &text, 0, nil)
		imgui.SameLine()
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS))
		if imgui.Button(confirm) {
			done = true
			c(text)
		}
	}
	Modal(&done, imgui.WindowFlagsAlwaysAutoResize, title, f, nil)
}
