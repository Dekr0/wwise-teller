package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

type DockManager struct {
	focused     int
	DockWindows []string
}

func NewDockManager() *DockManager {
	return &DockManager{
		focused: 0,
		DockWindows: []string{
			"Bank Explorer (Linear View)",
			"Bank Explorer (Tree View)",
			"File Explorer",
			"Log",
			"Object Editor",
		},
	}
}

func (d *DockManager) FocusNext() {
	if d.focused + 1 > len(d.DockWindows) - 1 {
		d.focused = 0
	} else {
		d.focused += 1
	}
}

func (d *DockManager) FocusPrev() {
	if d.focused - 1 < 0 {
		d.focused = len(d.DockWindows) - 1
	} else {
		d.focused -= 1
	}
}

func (d *DockManager) Focus() string {
	return d.DockWindows[d.focused]
}

func buildDockSpace(dockSpaceID imgui.ID, dockSpaceFlags imgui.DockNodeFlags) {
	imgui.InternalDockBuilderRemoveNode(dockSpaceID)
	imgui.InternalDockBuilderAddNodeV(dockSpaceID, dockSpaceFlags)
	
	mainDock := dockSpaceID
	dock1 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirLeft, 0.45, nil, &mainDock,
	)
	dock2 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirRight, 0.75, nil, &mainDock,
	)
	dock3 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirDown, 0.35, nil, &dock2,
	)
	
	imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
	imgui.InternalDockBuilderDockWindow("Bank Explorer", dock2)
	imgui.InternalDockBuilderDockWindow("Hierarchy View", dock2)
	imgui.InternalDockBuilderDockWindow("Object Editor", dock3)
	imgui.InternalDockBuilderDockWindow("Log", dock3)
	imgui.InternalDockBuilderFinish(mainDock)
}

