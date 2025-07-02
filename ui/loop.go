package ui

import (
	"container/ring"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/imguizmo"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/ui/async"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/ui/fs"
	glog "github.com/Dekr0/wwise-teller/ui/log"
	"github.com/Dekr0/wwise-teller/ui/notify"
	"github.com/Dekr0/wwise-teller/utils"
	"golang.design/x/clipboard"
)

var ModifiyEverything = false

const MainDockFlags imgui.WindowFlags = 
	imgui.WindowFlagsNoDocking |
	imgui.WindowFlagsNoTitleBar | 
	imgui.WindowFlagsNoCollapse |
	imgui.WindowFlagsNoResize |
	imgui.WindowFlagsNoMove |
	imgui.WindowFlagsNoBringToFrontOnFocus |
	imgui.WindowFlagsNoNavFocus |
	imgui.WindowFlagsMenuBar

func Run() error {
	runtime.LockOSThread()

	err := clipboard.Init()
	if err != nil {
		slog.Error("Failed to initialized clipboard. Copying is disabled")
		GlobalCtx.CopyEnable = false
	}
	GlobalCtx.CopyEnable = true

	// Begin of app state
	gLog := &glog.GuiLog{
		Log: log.InMemoryLog{Logs: ring.New(log.DefaultSize)},
		Debug: true,
		Info: true,
		Warn: true,
		Error: true,
	}
	logF, err := os.OpenFile("wwise_teller.log", os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(&gLog.Log, logF),
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)))

	err = utils.ScanMountPoint()
	if err != nil {
		return err
	}

	err = config.Load(&GlobalCtx.Config)
	if err != nil {
		return err
	}
	slog.Info("Loaded configuration file")

	bnkMngr := &be.BankManager{WriteLock: atomic.Bool{}}
	bnkMngr.WriteLock.Store(false)
	slog.Info("Created bank manager")

	dockMngr := dockmanager.NewDockManager()

	fileExplorer, err := fs.NewFileExplorer(
		openSoundBankFunc(bnkMngr), GlobalCtx.Config.Home,
	)
	if err != nil {
		return err
	}
	slog.Info("Created file explorer backend for file exploring")

	nQ := &notify.NotifyQ{Queue: make([]notify.Notfiy, 0, 16)}
	// End of app state

	backend, err := setupBackend()
	if err != nil {
		return err
	}
	backend.SetAfterRenderHook(createAfterRenderHook(nQ))
	backend.SetBeforeDestroyContextHook(func() {
		if err := GlobalCtx.Config.Save(); err != nil {
			slog.Error("Failed to save configuration file", "error", err)
		}
		logF.Close()
	})
	slog.Info("Created rendering backend and registered all necessary hooks")

	if err := setupImGUI(); err != nil {
		return err
	}
	slog.Info("Setup ImGUI context and configuration")

	backend.Run(createLoop(
		dockMngr,
		fileExplorer,
		NewCmdPaletteMngr(dockMngr),
		bnkMngr, 
		nQ,
		gLog,
	))

	return nil
}

func setupBackend() (backend.Backend[sdlbackend.SDLWindowFlags], error) {
	backend, err := backend.CreateBackend(sdlbackend.NewSDLBackend())
	if err != nil {
		return nil, err
	}

	backend.SetBgColor(imgui.NewVec4(0.0, 0.0, 0.0, 1.0))
	backend.CreateWindow("Wwise Teller", 1280, 720)

	return backend, nil
}

func setupImGUI() error {
	imgui.CreateContext()
	implot.CreateContext()
	imgui.CurrentIO().SetConfigFlags(
		imgui.ConfigFlagsDockingEnable |
		imgui.ConfigFlagsViewportsEnable |
		imgui.ConfigFlagsNavEnableKeyboard,
	)
	return nil
}

func createAfterRenderHook(nQ *notify.NotifyQ) func() {
	return func() {
		for _, onDone := range GlobalCtx.Loop.Update() {
			nQ.Q(onDone, time.Second * 8)
		}
	}
}

// The loop try to follow the following rule:
// 1. Prioritize return computed state over persistence state with computation
// 2. If point 1 cannot be done, a render function only accept the state it needs 
// to use
func createLoop(
	dockMngr *dockmanager.DockManager,
	fileExplorer *fs.FileExplorer,
	cmdPaletteMngr *CmdPaletteMngr,
	bnkMngr *be.BankManager,
	nQ *notify.NotifyQ,
	gLog *glog.GuiLog,
) func() {
	return func() {
		imgui.ClearSizeCallbackPool()
		imguizmo.BeginFrame()

		viewport := imgui.MainViewport()

		if imgui.Shortcut(imgui.KeyChord(imgui.KeyF1)) {
			dockMngr.SetLayout(dockmanager.ActorMixerObjEditorLayout)
		}
		if imgui.Shortcut(imgui.KeyChord(imgui.KeyF2)) {
			dockMngr.SetLayout(dockmanager.ActorMixerEventLayout)
		}
		if imgui.Shortcut(imgui.KeyChord(imgui.KeyF3)) {
			dockMngr.SetLayout(dockmanager.MasterMixerLayout)
		}

		renderStatusBar(GlobalCtx.Loop.AsyncTasks)

		imgui.SetNextWindowPos(viewport.WorkPos())
		imgui.SetNextWindowSize(viewport.WorkSize())
		imgui.SetNextWindowViewport(viewport.ID())

		imgui.BeginV("MainDock", nil, MainDockFlags)

		dockSpaceID := dockMngr.BuildDockSpace()
		imgui.DockSpaceV(
			dockSpaceID,
			DefaultSize,
			dockmanager.DockSpaceFlags,
			imgui.NewEmptyWindowClass(),
		)

		renderMainMenuBar(dockMngr, cmdPaletteMngr)
		GlobalCtx.ModalQ.renderModal()

		glog.RenderLog(gLog, &dockMngr.Opens[dockmanager.LogTag])

		renderDebug(bnkMngr, &dockMngr.Opens[dockmanager.DebugTag])
		renderFileExplorer(fileExplorer, &dockMngr.Opens[dockmanager.FileExplorerTag])
		renderBankExplorer(bnkMngr)
		renderActorMixerHircTree(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.ActorMixerHierarchyTag])
		renderMusicHircTree(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.MusicHierarchyTag])
		renderMasterMixerHierarchy(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.MasterMixerHierarchyTag])
		renderObjEditorActorMixer(bnkMngr, bnkMngr.ActiveBank, bnkMngr.InitBank, &dockMngr.Opens[dockmanager.ObjectEditorActorMixerTag])
		renderObjEditorMusic(bnkMngr, bnkMngr.ActiveBank, bnkMngr.InitBank, &dockMngr.Opens[dockmanager.ObjectEditorMusicTag])
		renderBusViewer(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.BusesTag])
		renderFXViewer(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.FXTag])
		renderEventsViewer(bnkMngr.ActiveBank, &dockMngr.Opens[dockmanager.EventsTag])
		notify.RenderNotify(nQ)
		imgui.End()
	}
}

func renderTasks(asyncTasks []*async.Task) {
	if !imgui.BeginPopupV("Tasks", imgui.WindowFlagsAlwaysAutoResize) {
		return
	}
	for i, a := range asyncTasks {
		if a == nil { continue }

		imgui.PushIDStr(fmt.Sprintf("XTask_%d", i))
		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)
		if imgui.Button("X") { a.Cancel() }
		imgui.PopStyleColor()
		imgui.PopID()

		imgui.SameLine()

		imgui.Text(fmt.Sprintf("Task %d", i))
		if imgui.IsItemHoveredV(imgui.HoveredFlagsDelayNone) {
			imgui.SetTooltip(a.OnProcMsg)
		}

		imgui.SameLine()

		imgui.ProgressBarV(
			float32(-1.0 * imgui.Time()),
			imgui.NewVec2(0, 0),
			"Processing",
		)
	}
	imgui.EndPopup()
}
