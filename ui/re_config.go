package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
)

func configModalFunc() (func(), *bool) {
	done := false
	return func() {
		imgui.Checkbox("Disable All Guard Rails", &ModifiyEverything)
		imgui.Text(fmt.Sprintf("Home: %s", GCtx.Config.Home))
		imgui.SameLine()
		if imgui.ArrowButton("SetHome", imgui.DirRight) {
			pushSetHomeModal()
		}

		if imgui.Button("Apply") {
			if err := GCtx.Config.Save(); err != nil {
				slog.Error("Failed to save configuration")
			}
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			done = true
		}
	}, &done
}

func pushConfigModalFunc() {
	renderF, done := configModalFunc()
	Modal(done, imgui.WindowFlagsAlwaysAutoResize, "Config", renderF, nil)
}
