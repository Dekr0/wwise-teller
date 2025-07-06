package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/async"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
)

func renderMainMenuBar(
	dockMngr *dockmanager.DockManager,
	cmdPaletteMngr *CmdPaletteMngr,
) {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Layout") {
			for i := range dockmanager.LayoutCount {
				selected := dockMngr.Layout == i
				if imgui.MenuItemBoolV(dockmanager.LayoutName[i], fmt.Sprintf("F%d", i + 1), selected, true) {
					dockMngr.SetLayout(i)
				}
			}
			imgui.EndMenu()
		}

		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)
		if imgui.Button("Config") {
			pushConfigModalFunc()
		}
		imgui.PopStyleColor()

		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)
		imgui.SetNextItemShortcutV(
			imgui.KeyChord(imgui.ModCtrl) | 
			imgui.KeyChord(imgui.ModShift) | 
			imgui.KeyChord(imgui.KeyP),
			imgui.InputFlagsRouteGlobal,
		)
		if imgui.Button("Command") {
			pushCommandPaletteModal(cmdPaletteMngr)
		}
		imgui.PopStyleColor()

		if imgui.BeginMenu("Integration") {
			if imgui.BeginMenu("Helldivers 2") {
				if imgui.MenuItemBool("Extract sound banks from game archives") {
					pushSelectGameArchiveModal()
				}
				imgui.EndMenu()
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu("Views") {
			for tag, open := range dockMngr.Opens {
				if imgui.MenuItemBoolV(dockmanager.DockWindowNames[tag], "", open, true) {
					dockMngr.Opens[tag] = !open
				}
			}
			imgui.EndMenu()
		}

		imgui.EndMenuBar()
	}
}

func renderStatusBar(asyncTasks []*async.Task) {
	renderTaskPopup := false

	viewport := imgui.MainViewport()

	statusBarFlags := imgui.WindowFlagsMenuBar | imgui.WindowFlagsNoScrollbar
	statusBarFlags |= imgui.WindowFlagsNoSavedSettings

	if imgui.InternalBeginViewportSideBar("StatusBar", 
		viewport, imgui.DirDown, imgui.FrameHeight(), statusBarFlags,
	) {
		if imgui.BeginMenuBar() {
			imgui.PushStyleColorVec4(
				imgui.ColButton,
				imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
			)
			if imgui.Button("Task status") {
				renderTaskPopup = true
			}
			imgui.PopStyleColor()

			if renderTaskPopup {
				if !imgui.IsPopupOpenStr("Tasks") {
					imgui.OpenPopupStr("Tasks")
				}
			}
			renderTasks(asyncTasks)
			imgui.SameLine()
			imgui.Text(fmt.Sprintf("%f FPS", imgui.CurrentIO().Framerate()))

			imgui.EndMenuBar()
		}
		imgui.End()
	}
}
