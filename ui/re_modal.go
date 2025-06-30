package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
)

func (m *ModalQ) renderModal() {
	if len(m.Modals) <= 0 {
		return
	}
	top := m.Modals[len(m.Modals)-1]
	if *top.done {
		imgui.CloseCurrentPopup()
		m.Modals = m.Modals[:len(m.Modals)-1]
		if top.onClose != nil {
			top.onClose()
		}
		return
	}
	if !imgui.IsPopupOpenStr(top.name) {
		imgui.OpenPopupStr(top.name)
		imgui.SetNextWindowSize(imgui.NewVec2(640, 640))
	}
	center := imgui.MainViewport().Center()
	imgui.SetNextWindowPosV(center, imgui.CondAppearing, imgui.NewVec2(0.5, 0.5))
	if imgui.BeginPopupModalV(top.name, nil, top.flag) {
		if imgui.Shortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyQ)) {
			*top.done = true
		}
		top.loop()
		imgui.EndPopup()
	}
}

func pushConfigModalFunc() {
	renderF, done := configModalFunc()
	GlobalCtx.ModalQ.QModal(
		done,
		imgui.WindowFlagsAlwaysAutoResize,
		"Config",
		renderF, nil,
	)
}

func pushSetHomeModal() {
	onSave := func(path string) {
		if err := GlobalCtx.Config.SetHome(path); err != nil {
			slog.Error(
				"Failed to set initial directory for file " +
				"explorer",
				"error", err,
			)
		}
	}
	renderF, done, err := saveFileDialogFunc(onSave, GlobalCtx.Config.Home)
	if err != nil {
		slog.Error(
			"Failed to create save file dialog for setting initial" + 
			" directory for file explorer",
			"error", err,
		)
	} else {
		GlobalCtx.ModalQ.QModal(
			done,
			0,
			"Set starting directory for file explorer",
			renderF, nil,
		)
	}
}

func pushSaveSoundBankModal(
	bnkMngr *be.BankManager,
	saveTab *be.BankTab,
	saveName string,
) {
	onSave := saveSoundBankFunc(bnkMngr, saveTab, saveName)
	renderF, done, err := saveFileDialogFunc(onSave, GlobalCtx.Config.Home)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Failed create save file dialog for saving sound bank %s",
				saveName,
			),
			"error", err,
		)
	} else {
		GlobalCtx.ModalQ.QModal(
			done,
			0,
			fmt.Sprintf("Save sound bank %s to ...", saveName),
			renderF, nil,
		)
	}
}

func pushHD2PatchModal(
	bnkMngr *be.BankManager,
	saveTab *be.BankTab,
	saveName string,
) {
	onSave := HD2PatchFunc(bnkMngr, saveTab, saveName)
	if renderF, done, err := saveFileDialogFunc(onSave, GlobalCtx.Config.Home);
	   err != nil {
		slog.Error(
			fmt.Sprintf("Failed create save file dialog for saving sound " +
				"bank %s to HD2 patch", saveName,
			),
			"error", err,
		)
	} else {
		GlobalCtx.ModalQ.QModal(
			done,
			0,
			fmt.Sprintf("Save sound bank %s to HD2 patch ...", saveName),
			renderF, nil,
		)
	}
}

func pushSelectGameArchiveModal() {
	onOpen := selectHD2ArchiveFunc()
	initialDir, err := helldivers.GetHelldivers2Data()
	if err != nil {
		initialDir = GlobalCtx.Config.Home
	}
	renderF, done, err := openFileDialogFunc(
		onOpen, false, initialDir, []string{},
	)
	if err != nil {
		slog.Error(
			"Failed to create open file dialog for opening " +
			"Helldivers 2 game archives",
			"error", err,
		)
	} else {
		GlobalCtx.ModalQ.QModal(
			done, 
			0,
			"Select Helldivers 2 game archives",
			renderF, nil,
		)
	}
}

func pushExtractSoundBanksModal(paths []string) {
	onSave := extractHD2SoundBanksFunc(paths)
	renderF, done, err := saveFileDialogFunc(onSave, GlobalCtx.Config.Home)
	if err != nil {
		slog.Error(
			"Failed create save file dialog for saving extracted sound banks",
			"error", err,
		)
		return
	}
	GlobalCtx.ModalQ.QModal(
		done,
		0,
		"Save extracted sound banks to ...",
		renderF, nil,
	)
}

func pushCommandPaletteModal(cmdPaletteMngr *CmdPaletteMngr) {
	renderF, done := commandPaletteModal(cmdPaletteMngr)
	GlobalCtx.ModalQ.QModal(done, 0, "Command Palette", renderF, nil)
}

func pushSimpleTextModal(title string, c func(string)) {
	renderF, done := simpleTextModal(c)
	GlobalCtx.ModalQ.QModal(done, imgui.WindowFlagsAlwaysAutoResize, title, renderF, nil)
}

func simpleTextModal(c func(string)) (func(), *bool) {
	done := false
	text := ""
	return func() {
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF))
		imgui.InputTextWithHint("Directory Name", "", &text, 0, nil)
		imgui.SameLine()
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS))
		if imgui.Button("Create") {
			done = true
			c(text)
		}
	}, &done
}
