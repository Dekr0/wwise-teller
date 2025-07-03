package ui

import "github.com/AllenDang/cimgui-go/imgui"

func Disabled(disabled bool, render func()) {
	imgui.BeginDisabledV(disabled)
	render()
	imgui.EndDisabled()
}
