package ui

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

type Modal struct {
	done *bool
	flag imgui.WindowFlags
	name string
	loop func()
	onClose func()
}

type ModalQ struct {
	modals []*Modal
}

func NewModalQ() *ModalQ {
	return &ModalQ{make([]*Modal, 0, 4)}
}

func (m *ModalQ) QModal(
	done *bool,
	flag imgui.WindowFlags,
	name string,
	loop func(),
	onClose func(),
) {
	if slices.ContainsFunc(m.modals, func(modal *Modal) bool {
		return strings.Compare(modal.name, name) == 0
	}) {
		return
	}
	m.modals = append(m.modals, &Modal{done, flag, name, loop, onClose})
}

func (m *ModalQ) ShowModal() {
	if len(m.modals) <= 0 {
		return
	}
	top := m.modals[len(m.modals)-1]
	if *top.done {
		imgui.CloseCurrentPopup()
		m.modals = m.modals[:len(m.modals)-1]
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

func pushConfigModalFunc(modalQ *ModalQ, conf *config.Config) {
	renderF, done := configModalFunc(modalQ, conf)
	modalQ.QModal(
		done,
		imgui.WindowFlagsAlwaysAutoResize,
		"Config",
		renderF, nil,
	)
}

func pushSetHomeModal(modalQ *ModalQ, conf *config.Config) {
	onSave := func(path string) {
		if err := conf.SetHome(path); err != nil {
			slog.Error(
				"Failed to set initial directory for file " +
				"explorer",
				"error", err,
			)
		}
	}
	renderF, done, err := saveFileDialogFunc(onSave, conf.Home)
	if err != nil {
		slog.Error(
			"Failed to create save file dialog for setting initial" + 
			" directory for file explorer",
			"error", err,
		)
	} else {
		modalQ.QModal(
			done,
			0,
			"Set starting directory for file explorer",
			renderF, nil,
		)
	}
}

func pushSetHelldivers2DataModal(modalQ *ModalQ, conf *config.Config) {
	onSave := func(path string) {
		if err := conf.SetHelldiversData(path); err != nil {
			slog.Error(
				"Failed to set Helldivers 2 data directory",
				"error", err,
			)
		}
	}
	renderF, done, err := saveFileDialogFunc(onSave, conf.HelldiversData)
	if err != nil {
		slog.Error(
			"Failed to create save file dialog for setting Helldivers " +
			"2 data directory",
			"error", err,
		)
	} else {
		modalQ.QModal(
			done,
			0,
			"Set Helldivers 2 data directory",
			renderF, nil,
		)
	}
}

func pushSaveSoundBankModal(
	modalQ *ModalQ,
	loop *async.EventLoop,
	conf *config.Config,
	bnkMngr *BankManager,
	saveTab *bankTab,
	saveName string,
) {
	onSave := saveSoundBankFunc(loop, bnkMngr, saveTab, saveName)
	renderF, done, err := saveFileDialogFunc(onSave, conf.DefaultSave)
	if err != nil {
		slog.Error(
			fmt.Sprintf("Failed create save file dialog for saving sound bank %s",
				saveName,
			),
			"error", err,
		)
	} else {
		modalQ.QModal(
			done,
			0,
			fmt.Sprintf("Save sound bank %s to ...", saveName),
			renderF, nil,
		)
	}
}

func pushHD2PatchModal(
	modalQ *ModalQ,
	loop *async.EventLoop,
	conf *config.Config,
	bnkMngr *BankManager,
	saveTab *bankTab,
	saveName string,
) {
	onSave := HD2PatchFunc(loop, bnkMngr, saveTab, saveName)
	if renderF, done, err := saveFileDialogFunc(onSave, conf.DefaultSave);
	   err != nil {
		slog.Error(
			fmt.Sprintf("Failed create save file dialog for saving sound " +
				"bank %s to HD2 patch", saveName,
			),
			"error", err,
		)
	} else {
		modalQ.QModal(
			done,
			0,
			fmt.Sprintf("Save sound bank %s to HD2 patch ...", saveName),
			renderF, nil,
		)
	}
}

func pushSelectGameArchiveModal(
	modalQ *ModalQ,
	loop *async.EventLoop,
	conf *config.Config,
) {
	onOpen := selectGameArchiveFunc(modalQ, loop, conf)
	renderF, done, err := openFileDialogFunc(
		onOpen, false, conf.HelldiversData, []string{},
	)
	if err != nil {
		slog.Error(
			"Failed to create open file dialog for opening " +
			"Helldivers 2 game archives",
			"error", err,
		)
	} else {
		modalQ.QModal(
			done, 
			0,
			"Select Helldivers 2 game archives",
			renderF, nil,
		)
	}
}

func pushExtractSoundBanksModal(
	modalQ *ModalQ,
	loop *async.EventLoop,
	conf *config.Config,
	paths []string,
) {
	onSave := extractSoundBanksFunc(loop, paths)
	renderF, done, err := saveFileDialogFunc(onSave, conf.DefaultSave)
	if err != nil {
		slog.Error(
			"Failed create save file dialog for saving extracted sound banks",
			"error", err,
		)
		return
	}
	modalQ.QModal(
		done,
		0,
		"Save extracted sound banks to ...",
		renderF, nil,
	)
}

