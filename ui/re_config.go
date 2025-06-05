package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/config"
)

func configModalFunc(modalQ *ModalQ,conf *config.Config) (func(), *bool) {
	done := false
	return func() {
		imgui.Checkbox("Disable All Guard Rails", &ModifiyEverything)
		imgui.Text(fmt.Sprintf("Home: %s", conf.Home))
		imgui.SameLine()
		if imgui.ArrowButton("SetHome", imgui.DirRight) {
			pushSetHomeModal(modalQ, conf)
		}
		imgui.Text(fmt.Sprintf("Helldivers 2 Data Directory: %s", conf.HelldiversData))
		imgui.SameLine()
		if imgui.ArrowButton("SetHelldivers2Data", imgui.DirRight) {
			pushSetHelldivers2DataModal(modalQ, conf)
		}
		if imgui.Button("Save") {
			if err := conf.Save(); err != nil {
				slog.Error("Failed to save configuration")
			}
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			done = true
		}
	}, &done
}
