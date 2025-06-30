package dockmanager

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

type Layout uint8

const (
	Layout01    Layout = 0
	Layout02    Layout = 1
	LayoutCount Layout = 2
)

type DockManager struct {
	Focused     int
	DockWindows []string
	Layout      Layout
	Rebuild     bool
}

const DockSpaceFlags imgui.DockNodeFlags = 
	imgui.DockNodeFlagsNone |
	imgui.DockNodeFlags(imgui.DockNodeFlagsNoWindowMenuButton)

func NewDockManager() *DockManager {
	return &DockManager{
		Focused: 0,
		DockWindows: []string{
			"Attenuations",
			"Actor Mixer Hierarchy",
			"Bank Explorer",
			"Buses",
			"Debug",
			"Events",
			"File Explorer",
			"FX",
			"Game Sync",
			"Log",
			"Master Mixer Hierarchy",
			"Music Hierarchy",
			"Object Editor (Actor Mixer)",
			"Object Editor (Music)",
		},
		Layout: Layout02,
		Rebuild: true,
	}
}

func (d *DockManager) FocusNext() {
	if d.Focused + 1 > len(d.DockWindows) - 1 {
		d.Focused = 0
	} else {
		d.Focused += 1
	}
}

func (d *DockManager) FocusPrev() {
	if d.Focused - 1 < 0 {
		d.Focused = len(d.DockWindows) - 1
	} else {
		d.Focused -= 1
	}
}

func (d *DockManager) Focus() string {
	return d.DockWindows[d.Focused]
}

func (d *DockManager) BuildDockSpace() imgui.ID {
	dockSpaceID := imgui.IDStr("MainDock")
	if !d.Rebuild {
		return dockSpaceID
	}

	if d.Layout == Layout01 {
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

		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Attenuations", dock3)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", dock2)
		imgui.InternalDockBuilderDockWindow("Buses", dock3)
		imgui.InternalDockBuilderDockWindow("Debug", dock3)
		imgui.InternalDockBuilderDockWindow("Events", dock3)
		imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("FX", dock3)
		imgui.InternalDockBuilderDockWindow("Game Sync", dock3)
		imgui.InternalDockBuilderDockWindow("Log", dock3)
		imgui.InternalDockBuilderDockWindow("Master Mixer Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Music Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", dock3)
		imgui.InternalDockBuilderDockWindow("Object Editor (Music)", dock3)
		imgui.InternalDockBuilderFinish(mainDock)
		d.Rebuild = false
	} else if d.Layout == Layout02 {
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

		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Attenuations", dock3)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("Buses", dock3)
		imgui.InternalDockBuilderDockWindow("Debug", dock3)
		imgui.InternalDockBuilderDockWindow("Events", dock3)
		imgui.InternalDockBuilderDockWindow("FX", dock3)
		imgui.InternalDockBuilderDockWindow("Game Sync", dock3)
		imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
		imgui.InternalDockBuilderDockWindow("Log", dock3)
		imgui.InternalDockBuilderDockWindow("Master Mixer Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Music Hierarchy", dock2)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", dock3)
		imgui.InternalDockBuilderDockWindow("Object Editor (Music)", dock3)
		imgui.InternalDockBuilderFinish(mainDock)
		d.Rebuild = false
	}
	return dockSpaceID
}

func (d *DockManager) SetFocus(s string) {
	imgui.SetWindowFocusStr(s)
}

func (d *DockManager) SetLayout(l Layout) {
	d.Layout = l
	d.Rebuild = true
}
