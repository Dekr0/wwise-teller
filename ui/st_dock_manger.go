package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

type Layout uint8

const (
	Layout01 Layout = 0
	Layout02 Layout = 1
)

type DockManager struct {
	focused     int
	dockWindows []string
	layout      Layout
	rebuild     bool
}

func NewDockManager() *DockManager {
	return &DockManager{
		focused: 0,
		dockWindows: []string{
			"Attenuations",
			"Actor Mixer Hierarchy Tree",
			"Bank Explorer",
			"Events",
			"File Explorer",
			"Game Sync",
			"Log",
			"Object Editor (Actor Mixer)",
			"Object Editor (Music)",
		},
		layout: Layout02,
		rebuild: true,
	}
}

func (d *DockManager) FocusNext() {
	if d.focused + 1 > len(d.dockWindows) - 1 {
		d.focused = 0
	} else {
		d.focused += 1
	}
}

func (d *DockManager) FocusPrev() {
	if d.focused - 1 < 0 {
		d.focused = len(d.dockWindows) - 1
	} else {
		d.focused -= 1
	}
}

func (d *DockManager) Focus() string {
	return d.dockWindows[d.focused]
}

func (d *DockManager) buildDockSpace() imgui.ID {
	dockSpaceID := imgui.IDStr("MainDock")
	if !d.rebuild {
		return dockSpaceID
	}

	if d.layout == Layout01 {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		dock1 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.45, nil, &mainDock,
		)
		dock2 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirRight, 0.75, nil, &mainDock,
		)
		dock3 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirDown, 0.45, nil, &dock2,
		)

		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy Tree", dock2)
		imgui.InternalDockBuilderDockWindow("Attenuations", dock3)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", dock2)
		imgui.InternalDockBuilderDockWindow("Events", dock3)
		imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("Game Sync", dock3)
		imgui.InternalDockBuilderDockWindow("Log", dock3)
		imgui.InternalDockBuilderDockWindow("Music Hierarchy Tree", dock2)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", dock3)
		imgui.InternalDockBuilderDockWindow("Object Editor (Music)", dock3)
		imgui.InternalDockBuilderFinish(mainDock)
		d.rebuild = false
	} else if d.layout == Layout02 {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		dock1 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &mainDock,
			)
		dock2 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirRight, 0.60, nil, &mainDock,
			)
		dock3 := imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirRight, 0.50, nil, &dock2,
		)

		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy Tree", dock2)
		imgui.InternalDockBuilderDockWindow("Attenuations", dock3)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("Events", dock3)
		imgui.InternalDockBuilderDockWindow("Game Sync", dock3)
		imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("Log", dock3)
		imgui.InternalDockBuilderDockWindow("Music Hierarchy Tree", dock2)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", dock3)
		imgui.InternalDockBuilderDockWindow("Object Editor (Music)", dock3)
		imgui.InternalDockBuilderFinish(mainDock)
		d.rebuild = false
	}
	return dockSpaceID
}
