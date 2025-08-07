package ui

import (
	"container/ring"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/imguizmo"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/Dekr0/wwise-teller/aio"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/ui/async"
	uctx "github.com/Dekr0/wwise-teller/ui/context"
	"github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/ui/fs"
	glog "github.com/Dekr0/wwise-teller/ui/log"
	"github.com/Dekr0/wwise-teller/ui/notify"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/waapi"
	"golang.design/x/clipboard"
)

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

	utils.InitTmp()
	waapi.InitWEMCache()

	uctx.Init()

	err = clipboard.Init()
	if err != nil {
		slog.Error("Failed to initialized clipboard. Copying is disabled.")
		GCtx.CopyEnable = false
	}
	GCtx.CopyEnable = true

	err = aio.InitBeep()
	if err != nil {
		slog.Error("Failed to initialized audio player. WEM playback is disabled.")
	}

	err = utils.ScanMountPoint()
	if err != nil {
		return err
	}

	err = config.Load(Config)
	if err != nil {
		return err
	}
	slog.Info("Loaded configuration file")

	fileExplorer, err := fs.NewFileExplorer(fileExplorerCallback, GCtx.Config.Home)
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
		if err := GCtx.Config.Save(); err != nil {
			slog.Error("Failed to save configuration file", "error", err)
		}
		logF.Close()
	})
	slog.Info("Created rendering backend and registered all necessary hooks")

	if err := setupImGUI(); err != nil {
		return err
	}
	slog.Info("Setup ImGUI context and configuration")

	c := CmdPaletteMngr{}
	NewCmdPaletteMngrP(&c, DockMngr)

	backend.Run(createLoop(fileExplorer, &c, nQ, gLog))

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
		for _, onDone := range GCtx.Loop.Update() {
			nQ.Q(onDone, time.Second * 8)
		}
	}
}

// The loop try to follow the following rule:
// 1. Prioritize return computed state over persistence state with computation
// 2. If point 1 cannot be done, a render function only accept the state it needs 
// to use
func createLoop(
	fileExplorer *fs.FileExplorer,
	cmdPaletteMngr *CmdPaletteMngr,
	nQ *notify.NotifyQ,
	gLog *glog.GuiLog,
) func() {
	return func() {
		imgui.ClearSizeCallbackPool()
		imguizmo.BeginFrame()

		if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyF1), imgui.InputFlagsRouteGlobal) {
			DockMngr.SetLayout(dockmanager.ActorMixerObjEditorLayout)
		}
		if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyF2), imgui.InputFlagsRouteGlobal) {
			DockMngr.SetLayout(dockmanager.ActorMixerEventLayout)
		}
		if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyF3), imgui.InputFlagsRouteGlobal) {
			DockMngr.SetLayout(dockmanager.MasterMixerLayout)
		}
		if imgui.ShortcutNilV(imgui.KeyChord(imgui.KeyF4), imgui.InputFlagsRouteGlobal) {
			DockMngr.SetLayout(dockmanager.AttenuationLayout)
		}

		viewport := imgui.MainViewport()

		renderStatusBar(GCtx.Loop.AsyncTasks)

		imgui.SetNextWindowPos(viewport.WorkPos())
		imgui.SetNextWindowSize(viewport.WorkSize())
		imgui.SetNextWindowViewport(viewport.ID())

		imgui.BeginV("MainDock", nil, MainDockFlags)

		dockSpaceID := DockMngr.BuildDockSpace()
		imgui.DockSpaceV(
			dockSpaceID,
			DefaultSize,
			dockmanager.DockSpaceFlags,
			imgui.NewEmptyWindowClass(),
		)

		renderMainMenuBar(DockMngr, cmdPaletteMngr)
		renderModal(&GCtx.ModalQ)

		glog.RenderLog(gLog, &DockMngr.Opens[dockmanager.LogTag])

		renderDebug(&DockMngr.Opens[dockmanager.DebugTag])
		renderFileExplorer(fileExplorer, &DockMngr.Opens[dockmanager.FileExplorerTag])
		renderBankExplorer()
		renderActorMixerHircTree(&DockMngr.Opens[dockmanager.ActorMixerHierarchyTag])
		renderMusicHircTree(&DockMngr.Opens[dockmanager.MusicHierarchyTag])
		renderMasterMixerHierarchy(&DockMngr.Opens[dockmanager.MasterMixerHierarchyTag])
		renderObjEditorActorMixer(&DockMngr.Opens[dockmanager.ObjectEditorActorMixerTag])
		renderObjEditorMusic(&DockMngr.Opens[dockmanager.ObjectEditorMusicTag])
		renderBusViewer(&DockMngr.Opens[dockmanager.BusesTag])
		renderFXViewer(&DockMngr.Opens[dockmanager.FXTag])
		renderEventsViewer(&DockMngr.Opens[dockmanager.EventsTag])
		renderAttenuationViewer(&DockMngr.Opens[dockmanager.AttenuationsTag])
		RenderTransportControl(&DockMngr.Opens[dockmanager.TransportControlTag])
		// processor.RenderProcessorEditor(&GCtx.Editor, &DockMngr.Opens[dockmanager.ProcessorEditorTag])

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
