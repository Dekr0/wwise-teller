package dockmanager

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

type Layout uint8

const (
	ActorMixerObjEditorLayout    Layout = 0
	ActorMixerEventLayout        Layout = 1
	MasterMixerLayout            Layout = 2
	LayoutCount                  Layout = 3
)

var LayoutName []string = []string{
	"Editor Layout (Actor Mixer)",
	"Event Layout (Actor Mixer)",
	"Master Mixer Layout",
}

type DockWindowTag uint8

const (
	AttenuationsTag           DockWindowTag = 0
	ActorMixerHierarchyTag    DockWindowTag = 1
	BankExplorerTag           DockWindowTag = 2
	BusesTag                  DockWindowTag = 3
	DebugTag                  DockWindowTag = 4
	EventsTag                 DockWindowTag = 5
	FileExplorerTag           DockWindowTag = 6
	FXTag                     DockWindowTag = 7
	GameSyncTag               DockWindowTag = 8
	LogTag                    DockWindowTag = 9
	MasterMixerHierarchyTag   DockWindowTag = 10
	MusicHierarchyTag         DockWindowTag = 11
	ObjectEditorActorMixerTag DockWindowTag = 12
	ObjectEditorMusicTag      DockWindowTag = 13
)

var DockWindowNames []string = []string{
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
}

type DockManager struct {
	Opens         []bool
	Focused       int
	Layout        Layout
	Rebuild       bool
}

const DockSpaceFlags imgui.DockNodeFlags = 
	imgui.DockNodeFlagsNone |
	imgui.DockNodeFlags(imgui.DockNodeFlagsNoWindowMenuButton)

func NewDockManager() *DockManager {
	open := make([]bool, len(DockWindowNames))
	for i := range DockWindowNames {
		open[i] = false
	}

	return &DockManager{
		Opens: open,
		Layout: ActorMixerEventLayout,
		Rebuild: true,
	}
}

func (d *DockManager) HideAllDockingWindow() {
	for tag := range d.Opens {
		d.Opens[tag] = false
	}
}

func (d *DockManager) BuildDockSpace() imgui.ID {
	dockSpaceID := imgui.IDStr("MainDock")
	if !d.Rebuild {
		return dockSpaceID
	}

	d.HideAllDockingWindow()

	if d.Layout == ActorMixerEventLayout {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		explorerDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		editorDock := imgui.ID(0)
		eventDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &editorDock,
		)
		editorDock = imgui.InternalDockBuilderSplitNode(
			editorDock, imgui.DirDown, 0.50, nil, &eventDock,
		)

		opens := []DockWindowTag{
			FileExplorerTag,
			BankExplorerTag,
			ActorMixerHierarchyTag,
			MusicHierarchyTag,
			MasterMixerHierarchyTag,
			ObjectEditorActorMixerTag,
			EventsTag,
			GameSyncTag,
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow("File Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy", hierarchyDock)
		imgui.InternalDockBuilderDockWindow("Music Hierarchy", hierarchyDock)
		imgui.InternalDockBuilderDockWindow("Master Mixer Hierarchy", hierarchyDock)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", editorDock)
		imgui.InternalDockBuilderDockWindow("Events", eventDock)
		imgui.InternalDockBuilderDockWindow("Game Sync", eventDock)

		imgui.InternalDockBuilderFinish(eventDock)
		d.Rebuild = false
	} else if d.Layout == ActorMixerObjEditorLayout {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		explorerDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		editorDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &editorDock,
		)

		opens := []DockWindowTag{
			FileExplorerTag,
			BankExplorerTag,
			ActorMixerHierarchyTag,
			ObjectEditorActorMixerTag,
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow("File Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Actor Mixer Hierarchy", hierarchyDock)
		imgui.InternalDockBuilderDockWindow("Object Editor (Actor Mixer)", editorDock)

		imgui.InternalDockBuilderFinish(mainDock)
		d.Rebuild = false
	} else if d.Layout == MasterMixerLayout {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		explorerDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		busDock := imgui.ID(0)
		fxDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &busDock,
		)
		busDock = imgui.InternalDockBuilderSplitNode(
			busDock, imgui.DirDown, 0.50, nil, &fxDock,
		)

		opens := []DockWindowTag{
			FileExplorerTag,
			BankExplorerTag,
			MasterMixerHierarchyTag,
			BusesTag,
			FXTag,
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow("File Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Bank Explorer", explorerDock)
		imgui.InternalDockBuilderDockWindow("Master Mixer Hierarchy", hierarchyDock)
		imgui.InternalDockBuilderDockWindow("Buses", busDock)
		imgui.InternalDockBuilderDockWindow("FX", fxDock)

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
