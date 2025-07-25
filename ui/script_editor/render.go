package script_editor

import (
	"github.com/AllenDang/cimgui-go/imgui"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
)

func renderScriptEditor(open *bool) {
	if !*open {
		return
	}

	imgui.Begin(dockmanager.DockWindowNames[dockmanager.ScriptEditorTag])
	defer imgui.End()

	if !*open {
		return
	}
}
