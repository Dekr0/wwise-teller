package ui

import (
	"container/ring"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/imguizmo"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/ui/async"
)

func Run() error {
	// Begin of app state
	gLog := &guiLog{
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

	conf, err := config.Load()
	if err != nil {
		return err
	}
	slog.Info("Loaded configuration file")

	loop := async.NewEventLoop()
	slog.Info("Created event loop")

	bnkMngr := &bankManager{writeLock: false}
	slog.Info("Created bank manager")

	openFile, err := newOpenFile(conf.Home, createOnOpenCallback(loop, bnkMngr))
	if err != nil {
		return err
	}
	slog.Info("Created file explorer backend for file exploring")

	saveFile, err := newSaveFile(conf.DefaultSave, nil)
	if err != nil {
		return err
	}
	slog.Info("Created file explorer backend for sound bank export")

	reBuildDockSpace := true

	nQ := &notifyQ{make([]*notfiy, 0, 16)}
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
		openFile,
		saveFile,
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
	imgui.CurrentIO().SetConfigFlags(imgui.ConfigFlagsDockingEnable)
	return nil
}

func createAfterRenderHook(loop *async.EventLoop, nQ *notifyQ) func() {
	return func() {
		for _, onDone := range loop.Update() {
			nQ.queue = append(nQ.queue, &notfiy{
				onDone, time.NewTimer(time.Second * 8),
			})
		}
	}
}

func createLoop(
	conf *config.Config,
	loop *async.EventLoop,
	openFile *openFile,
	saveFile *saveFile,
	bnkMngr *bankManager,
	nQ *notifyQ,
	gLog *guiLog,
	reBuildDockSpace *bool,
) func() {
	return func() {
		imgui.ClearSizeCallbackPool()
		imguizmo.BeginFrame()

		showStatusBar(loop.AsyncTasks)
		viewport := imgui.MainViewport()
		imgui.SetNextWindowPos(viewport.WorkPos())
		imgui.SetNextWindowSize(viewport.WorkSize())
		imgui.SetNextWindowViewport(viewport.ID())

		windowFlags := imgui.WindowFlagsNoDocking
		windowFlags |= imgui.WindowFlagsNoTitleBar | imgui.WindowFlagsNoCollapse
		windowFlags |= imgui.WindowFlagsNoResize | imgui.WindowFlagsNoMove
		windowFlags |= imgui.WindowFlagsNoBringToFrontOnFocus 
		windowFlags |= imgui.WindowFlagsNoNavFocus | imgui.WindowFlagsMenuBar

		imgui.BeginV("MainDock", nil, windowFlags)
		dockSpaceFlags := imgui.DockNodeFlagsNone 
		dockSpaceFlags |= imgui.DockNodeFlags(imgui.DockNodeFlagsNoWindowMenuButton)

		if imgui.BeginMenuBar() {
			if imgui.BeginMenu("Layout") {
				if imgui.MenuItemBool("Reset") {
					*reBuildDockSpace = true
				}
				imgui.EndMenu()
			}

			imgui.EndMenuBar()
		}

		dockSpaceID := imgui.IDStr("MainDock")
		if *reBuildDockSpace {
			buildDockSpace(dockSpaceID, dockSpaceFlags)
			*reBuildDockSpace = false
		}
		imgui.DockSpaceV(
			dockSpaceID,
			imgui.NewVec2(0, 0),
			dockSpaceFlags,
			imgui.NewEmptyWindowClass(),
		)

		showLog(gLog)

		showFileExplorer(openFile)
		
		activeTab, closeTab, saveTab, saveName := showBankExplorer(bnkMngr)
		if saveTab != nil {
			saveFile.onSave = createOnSaveCallback(
				loop, bnkMngr, saveTab, saveName,
			)
			saveFile.dest = filepath.Base(saveName) 
		}
		showSaveFileModal(saveFile, saveTab != nil)

		showObjectEditor(activeTab)
		showNotify(nQ)

		imgui.End()

		if closeTab != "" {
			bnkMngr.closeBank(closeTab)
		}
	}
}

func showStatusBar(asyncTasks []*async.AsyncTask) {
	showTaskPopup := false

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
				showTaskPopup = true
			}
			imgui.PopStyleColor()

			if showTaskPopup {
				if !imgui.IsPopupOpenStr("Tasks") {
					imgui.OpenPopupStr("Tasks")
				}
			}

			showTasks(asyncTasks)

			imgui.EndMenuBar()
		}
		imgui.End()
	}
}

func showTasks(asyncTasks []*async.AsyncTask) {
	if !imgui.BeginPopupV("Tasks", imgui.WindowFlagsAlwaysAutoResize) {
		return
	}

	for i, a := range asyncTasks {
		if a == nil { continue }

		imgui.PushIDStr(fmt.Sprintf("XTask_%", i))
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
			imgui.SetTooltip(a.OnProc)
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

func showNotify(nQ *notifyQ) {
	if len(nQ.queue) <= 0 {
		return
	}

	size := imgui.MainViewport().Size()
	y := imgui.CalcTextSize(nQ.queue[0].message).Y * float32(len(nQ.queue))
	x := float32(0.0)
	for _, n := range nQ.queue {
		x = max(imgui.CalcTextSize(n.message).X, x)
	}
	size.X -= size.X * 0.01 + x + 32.0
	size.Y -= size.Y * 0.03 + y + float32(len(nQ.queue)) * size.Y * 0.01
	imgui.SetNextWindowPos(size)

	imgui.BeginV("Notify", 
		nil,
		imgui.WindowFlagsNoDecoration |
		imgui.WindowFlagsNoTitleBar |
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoResize |
		imgui.WindowFlagsAlwaysAutoResize,
	)
	i := 0
	for i < len(nQ.queue) {
		select {
		case <- nQ.queue[i].timer.C:
			nQ.queue = slices.Delete(nQ.queue, i, i + 1) 
		default:
			imgui.PushIDStr(fmt.Sprintf("RemoveNotfiy", i))
			if imgui.Button("X") {
				nQ.queue = slices.Delete(nQ.queue, i, i + 1)
				imgui.PopID()
				break
			}
			imgui.PopID()
			imgui.SameLine()
			imgui.Text(nQ.queue[i].message)
		}
		i += 1
	}
	imgui.End()
}

func showLog(gLog *guiLog) {
	imgui.Begin("Log")
	gLog.log.Logs.Do(func(a any) {
		if a == nil {
			return
		}
		imgui.Text(a.(string))
	})
	imgui.End()
}

func buildDockSpace(dockSpaceID imgui.ID, dockSpaceFlags imgui.DockNodeFlags) {
	imgui.InternalDockBuilderRemoveNode(dockSpaceID)
	imgui.InternalDockBuilderAddNodeV(dockSpaceID, dockSpaceFlags)
	
	mainDock := dockSpaceID
	
	dock1 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirLeft, 0.25, nil, &mainDock,
	)
	
	dock2 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirRight, 0.75, nil, &mainDock,
	)
	
	dock3 := imgui.InternalDockBuilderSplitNode(
		mainDock, imgui.DirDown, 0.35, nil, &dock2,
	)
	
	imgui.InternalDockBuilderDockWindow("File Explorer", dock1)
	imgui.InternalDockBuilderDockWindow("Bank Explorer", dock2)
	imgui.InternalDockBuilderDockWindow("Object Editor", dock3)
	
	imgui.InternalDockBuilderFinish(mainDock)
}
