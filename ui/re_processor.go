package ui

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
)

func renderProcessorEditor(open *bool) {
	if !*open {
		return
	}

	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.ProcessorEditorTag], open, imgui.WindowFlagsMenuBar)
	defer imgui.End()

	if !*open {
		return
	}

	for _, path := range GCtx.Editor.Path {
		if imgui.BeginTabBarV("EditorTabBar", DefaultTabFlags) {
			imgui.PushIDStr(path)
			open := true
			if imgui.BeginTabItemV(filepath.Base(path), &open, imgui.TabItemFlagsNone) {
			}
			imgui.PopID()
		}
	}
}

func renderProcessorEditorMenu() {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.BeginMenu("New Processor") {
				callback := func(name string) {
					GCtx.Editor.NewProcessor(name)
				}
				PushSimpleTextModal("New Processor", "Processor Name", "Create", callback)
			}
		}
		imgui.EndMenuBar()
	}
}
