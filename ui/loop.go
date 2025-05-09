package ui

import (
	"container/ring"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/imguizmo"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/Dekr0/wwise-teller/utils"
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

const DockSpaceFlags imgui.DockNodeFlags = 
	imgui.DockNodeFlagsNone |
	imgui.DockNodeFlags(imgui.DockNodeFlagsNoWindowMenuButton)

func Run() error {
	// Begin of app state
	gLog := &GuiLog{
		log: &log.InMemoryLog{Logs: ring.New(log.DefaultSize)},
		debug: true,
		info: true,
		warn: true,
		error: true,
	}
	logF, err := os.OpenFile("wwise_teller.log", os.O_APPEND | os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(
		io.MultiWriter(gLog.log, logF),
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)))

	err = utils.ScanMountPoint()
	if err != nil {
		return err
	}

	conf, err := config.Load()
	if err != nil {
		return err
	}
	slog.Info("Loaded configuration file")

	loop := async.NewEventLoop()
	slog.Info("Created event loop")
	modalQ := NewModalQ()

	bnkMngr := &BankManager{writeLock: &atomic.Bool{}}
	bnkMngr.writeLock.Store(false)
	slog.Info("Created bank manager")

	dockMngr := NewDockManager()

	fileExplorer, err := newFileExplorer(
		openSoundBankFunc(loop, bnkMngr), conf.Home,
	)
	if err != nil {
		return err
	}
	slog.Info("Created file explorer backend for file exploring")

	reBuildDockSpace := true

	nQ := &NotifyQ{make([]*notfiy, 0, 16)}
	// End of app state

	backend, err := setupBackend()
	if err != nil {
		return err
	}
	backend.SetAfterRenderHook(createAfterRenderHook(loop, nQ))
	backend.SetBeforeDestroyContextHook(func() {
		if err := conf.Save(); err != nil {
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
		conf,
		loop,
		modalQ,
		dockMngr,
		fileExplorer,
		NewCmdPaletteMngr(&reBuildDockSpace, conf, loop, modalQ),
		bnkMngr, 
		nQ,
		gLog,
		&reBuildDockSpace,
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
	imgui.CurrentIO().SetConfigFlags(
		imgui.ConfigFlagsDockingEnable |
		imgui.ConfigFlagsNavEnableKeyboard,
	)
	return nil
}

func createAfterRenderHook(loop *async.EventLoop, nQ *NotifyQ) func() {
	return func() {
		for _, onDone := range loop.Update() {
			nQ.queue = append(nQ.queue, &notfiy{
				onDone, time.NewTimer(time.Second * 8),
			})
		}
	}
}

// The loop try to follow the following rule:
// 1. Prioritize return computed state over persistence state with computation
// 2. If point 1 cannot be done, a render function only accept the state it needs 
// to use
func createLoop(
	conf *config.Config,
	loop *async.EventLoop,
	modalQ *ModalQ,
	dockMngr *DockManager,
	fileExplorer *FileExplorer,
	cmdPaletteMngr *CmdPaletteMngr,
	bnkMngr *BankManager,
	nQ *NotifyQ,
	gLog *GuiLog,
	reBuildDockSpace *bool,
) func() {
	return func() {
		imgui.ClearSizeCallbackPool()
		imguizmo.BeginFrame()

		saveActive := false
		itype := -1
		viewport := imgui.MainViewport()

		if imgui.ShortcutNilV(DefaultSaveAsSC, imgui.InputFlagsRouteGlobal) {
			saveActive = true
			itype = -1
		}
		if imgui.ShortcutNilV(ModCtrlShift | imgui.KeyChord(imgui.KeyI), imgui.InputFlagsRouteGlobal) {
			saveActive = true
			itype = int(helldivers.IntegrationTypeHelldivers2)
		}
		if imgui.ShortcutNilV(DefaultNavPrevSC, imgui.InputFlagsRouteGlobal) {
			dockMngr.FocusPrev()
			imgui.SetWindowFocusStr(dockMngr.Focus())
		}
		if imgui.ShortcutNilV(DefaultNavNextSC, imgui.InputFlagsRouteGlobal) {
			dockMngr.FocusNext()
			imgui.SetWindowFocusStr(dockMngr.Focus())
		}

		showStatusBar(loop.AsyncTasks)

		imgui.SetNextWindowPos(viewport.WorkPos())
		imgui.SetNextWindowSize(viewport.WorkSize())
		imgui.SetNextWindowViewport(viewport.ID())

		imgui.BeginV("MainDock", nil, MainDockFlags)

		dockSpaceID := imgui.IDStr("MainDock")
		if *reBuildDockSpace {
			buildDockSpace(dockSpaceID, DockSpaceFlags)
			*reBuildDockSpace = false
		}
		imgui.DockSpaceV(
			dockSpaceID,
			imgui.NewVec2(0, 0),
			DockSpaceFlags,
			imgui.NewEmptyWindowClass(),
		)

		showMainMenuBar(reBuildDockSpace, conf, cmdPaletteMngr, modalQ, loop)

		modalQ.showModal()

		showFileExplorerWindow(fileExplorer, modalQ)

		activeTab, closeTab, saveTab, saveName, itype := showBankExplorer(
			bnkMngr, saveActive, itype,
		)
		if saveTab != nil {
			switch itype {
			case -1:
				pushSaveSoundBankModal(modalQ, loop, conf, bnkMngr, saveTab, saveName)
			case int(helldivers.IntegrationTypeHelldivers2):
				pushHD2PatchModal(modalQ, loop, conf, bnkMngr, saveTab, saveName)
			}
		}

		showObjectEditor(activeTab)

		showNotify(nQ)

		showLog(gLog)
		// showDevDebug(loop, modalQ)

		imgui.End()

		if closeTab != "" {
			bnkMngr.closeBank(closeTab)
		}
	}
}

func showTasks(asyncTasks []*async.Task) {
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

func showLog(gLog *GuiLog) {
	imgui.BeginV("Log", nil, imgui.WindowFlagsHorizontalScrollbar)
	gLog.log.Logs.Do(func(a any) {
		if a == nil {
			return
		}
		imgui.Text(a.(string))
	})
	imgui.End()
}

func showDevDebug(
	loop *async.EventLoop,
	modalQ *ModalQ,
) {
	imgui.Begin("Dev. Debug")
	imgui.Text(fmt.Sprintf("Number of modals: %d", len(modalQ.modals)))
	imgui.Text(fmt.Sprintf("Number of asynchronous tasks: %d", loop.NumTasks()))
	imgui.End()
}

