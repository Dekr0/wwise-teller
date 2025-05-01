package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

const ModCtrlShift =
	imgui.KeyChord(imgui.ModCtrl) |
	imgui.KeyChord(imgui.ModShift)
const DefaultNavPrevSC =
	ModCtrlShift |
	imgui.KeyChord(imgui.KeyJ)
const DefaultNavNextSC =
	ModCtrlShift |
	imgui.KeyChord(imgui.KeyK)
const DefaultSearchSC =
	imgui.KeyChord(imgui.ModCtrl) |
	imgui.KeyChord(imgui.KeyF)
const DefaultSaveAsSC =
	ModCtrlShift |
	imgui.KeyChord(imgui.KeyS)
const UnFocusQuerySC =
	ModCtrlShift |
	imgui.KeyChord(imgui.KeyF)

func isLeftShortcut() bool {
	return imgui.Shortcut(imgui.KeyChord(imgui.KeyLeftArrow)) || 
	       imgui.Shortcut(imgui.KeyChord(imgui.KeyH))
}

func isRightShortcut() bool {
	return imgui.Shortcut(imgui.KeyChord(imgui.KeyRightArrow)) || 
		   imgui.Shortcut(imgui.KeyChord(imgui.KeyL))
}

func isUpShortcut() bool {
	return imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyUpArrow), imgui.InputFlagsRepeat) || 
	       imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyK), imgui.InputFlagsRepeat)
}

func isDownShortcut() bool {
	return imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyDownArrow), imgui.InputFlagsRepeat) || 
	       imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyJ), imgui.InputFlagsRepeat)
}

// Miss ctrl Vi

func useViDown() {
	if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyJ), imgui.InputFlagsRepeat) {
		imgui.CurrentIO().AddKeyEvent(imgui.KeyDownArrow, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyDownArrow, false)
	}
}

func useViShiftDown() {
	if imgui.ShortcutNilV(imgui.KeyChord(imgui.ModShift) | imgui.KeyChord(imgui.KeyJ), imgui.InputFlagsRepeat) {
		imgui.CurrentIO().AddKeyEvent(imgui.KeyLeftShift, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyDownArrow, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyLeftShift, false)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyDownArrow, false)
	}
}

func useViUp() {
	if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyK), imgui.InputFlagsRepeat) {
		imgui.CurrentIO().AddKeyEvent(imgui.KeyUpArrow, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyUpArrow, false)
	}
}

func useViShiftUp() {
	if imgui.ShortcutNilV(imgui.KeyChord(imgui.ModShift) | imgui.KeyChord(imgui.KeyK), imgui.InputFlagsRepeat) {
		imgui.CurrentIO().AddKeyEvent(imgui.KeyLeftShift, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyUpArrow, true)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyLeftShift, false)
		imgui.CurrentIO().AddKeyEvent(imgui.KeyUpArrow, false)
	}
}

func useViEnter() bool {
	return imgui.Shortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS))
}

func clearQuery() {
	if imgui.ShortcutNilV(UnFocusQuerySC, 0) {
		imgui.SetKeyboardFocusHere()
	}
}
