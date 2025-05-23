package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

func renderMainMenuBar(
	dockMngr       *DockManager,
	conf           *config.Config,
	cmdPaletteMngr *CmdPaletteMngr,
	modalQ         *ModalQ,
	loop           *async.EventLoop,
) {
	if imgui.BeginMenuBar() {
		if imgui.BeginMenu("Layout") {
			if imgui.MenuItemBool("Layout 1") {
				dockMngr.layout = Layout01
				dockMngr.rebuild = true
			}
			if imgui.MenuItemBool("Layout 2") {
				dockMngr.layout = Layout02
				dockMngr.rebuild = true
			}
			imgui.EndMenu()
		}

		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)
		if imgui.Button("Config") {
			pushConfigModalFunc(modalQ, conf)
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
			pushCommandPaletteModal(modalQ, cmdPaletteMngr)
		}
		imgui.PopStyleColor()

		if imgui.BeginMenu("Integration") {
			if imgui.BeginMenu("Helldivers 2") {
				if imgui.MenuItemBool("Extract sound banks from game archives") {
					pushSelectGameArchiveModal(modalQ, loop, conf)
				}
				imgui.EndMenu()
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

