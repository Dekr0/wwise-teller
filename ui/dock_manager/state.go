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
	AttenuationsTag           DockWindowTag = iota
	ActorMixerHierarchyTag
	BankExplorerTag
	BusesTag
	DebugTag
	EventsTag
	FileExplorerTag
	FXTag
	GameSyncTag
	LogTag
	MasterMixerHierarchyTag
	MusicHierarchyTag
	ObjectEditorActorMixerTag
	ObjectEditorMusicTag
	NotificationTag
	TransportControlTag
	ProcessorEditorTag           
	DockWindowTagCount
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
	"Notifications",
	"Transport Control",
	"Processor Editor",
}

type DockManager struct {
	Opens         []bool
	Focused       int
	Layout        Layout
	Rebuild       bool
	Init          bool
}

const DockSpaceFlags imgui.DockNodeFlags = 
	imgui.DockNodeFlagsNone |
	imgui.DockNodeFlags(imgui.DockNodeFlagsNoWindowMenuButton)

func NewDockManagerP(d *DockManager) {
	open := make([]bool, DockWindowTagCount)
	for i := range DockWindowNames {
		open[i] = false
	}
	d.Opens = open
	d.Layout = ActorMixerEventLayout
	d.Rebuild = true
	d.Init = true
}

func NewDockManager() DockManager {
	open := make([]bool, DockWindowTagCount)
	for i := range DockWindowNames {
		open[i] = false
	}

	return DockManager{
		Opens: open,
		Layout: ActorMixerEventLayout,
		Rebuild: true,
		Init: true,
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
		notificationDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		editorDock := imgui.ID(0)
		transportDock := imgui.ID(0)
		eventDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		explorerDock = imgui.InternalDockBuilderSplitNode(
			explorerDock, imgui.DirUp, 0.85, nil, &notificationDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &editorDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirUp, 0.65, nil, &transportDock,
		)
		editorDock = imgui.InternalDockBuilderSplitNode(
			editorDock, imgui.DirDown, 0.50, nil, &eventDock,
		)

		opens := []DockWindowTag{
			BankExplorerTag,
			ActorMixerHierarchyTag,
			MusicHierarchyTag,
			MasterMixerHierarchyTag,
			ObjectEditorActorMixerTag,
			EventsTag,
			GameSyncTag,
			TransportControlTag,
		}
		if d.Init {
			opens = append(opens, FileExplorerTag)
			d.Init = false
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow(DockWindowNames[FileExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[BankExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[NotificationTag], notificationDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ActorMixerHierarchyTag], hierarchyDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[MusicHierarchyTag], hierarchyDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[MasterMixerHierarchyTag], hierarchyDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ObjectEditorActorMixerTag], editorDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ProcessorEditorTag], editorDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[TransportControlTag], transportDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[EventsTag], eventDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[GameSyncTag], eventDock)

		imgui.InternalDockBuilderFinish(eventDock)
		d.Rebuild = false
	} else if d.Layout == ActorMixerObjEditorLayout {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		explorerDock := imgui.ID(0)
		notificationDock := imgui.ID(0)
		transportDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		editorDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		explorerDock = imgui.InternalDockBuilderSplitNode(
			explorerDock, imgui.DirUp, 0.85, nil, &notificationDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &editorDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirUp, 0.65, nil, &transportDock,
		)

		opens := []DockWindowTag{
			BankExplorerTag,
			ActorMixerHierarchyTag,
			ObjectEditorActorMixerTag,
			TransportControlTag,
		}
		if d.Init {
			opens = append(opens, FileExplorerTag)
			d.Init = false
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow(DockWindowNames[FileExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[BankExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[NotificationTag], notificationDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ActorMixerHierarchyTag], hierarchyDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[TransportControlTag], transportDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ObjectEditorActorMixerTag], editorDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[ProcessorEditorTag], editorDock)

		imgui.InternalDockBuilderFinish(mainDock)
		d.Rebuild = false
	} else if d.Layout == MasterMixerLayout {
		imgui.InternalDockBuilderRemoveNode(dockSpaceID)
		imgui.InternalDockBuilderAddNodeV(dockSpaceID, DockSpaceFlags)

		mainDock := dockSpaceID
		explorerDock := imgui.ID(0)
		notificationDock := imgui.ID(0)
		hierarchyDock := imgui.ID(0)
		busDock := imgui.ID(0)
		fxDock := imgui.ID(0)

		explorerDock = imgui.InternalDockBuilderSplitNode(
			mainDock, imgui.DirLeft, 0.30, nil, &hierarchyDock,
		)
		explorerDock = imgui.InternalDockBuilderSplitNode(
			explorerDock, imgui.DirUp, 0.85, nil, &notificationDock,
		)
		hierarchyDock = imgui.InternalDockBuilderSplitNode(
			hierarchyDock, imgui.DirLeft, 0.50, nil, &busDock,
		)
		busDock = imgui.InternalDockBuilderSplitNode(
			busDock, imgui.DirDown, 0.50, nil, &fxDock,
		)

		opens := []DockWindowTag{
			BankExplorerTag,
			MasterMixerHierarchyTag,
			BusesTag,
			FXTag,
		}
		if d.Init {
			opens = append(opens, FileExplorerTag)
			d.Init = false
		}
		for _, tag := range opens {
			d.Opens[tag] = true
		}

		imgui.InternalDockBuilderDockWindow(DockWindowNames[FileExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[BankExplorerTag], explorerDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[NotificationTag], notificationDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[MasterMixerHierarchyTag], hierarchyDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[BusesTag], busDock)
		imgui.InternalDockBuilderDockWindow(DockWindowNames[FXTag], fxDock)

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
