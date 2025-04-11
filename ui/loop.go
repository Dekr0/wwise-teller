package ui

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/imguizmo"
)

var activeBackend backend.Backend[sdlbackend.SDLWindowFlags]

var state stateStore 

func Run() error {
	err := setup()
	if err != nil {
		return err
	}
	activeBackend.Run(loop)
	return nil
}

func setup() error {
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	home, err = filepath.Abs(home)
	if err != nil {
		return err
	}
	fe, err := newFileExplorer(home)
	if err != nil {
		return err
	}
	state = stateStore{context.Background(), fe, newEventBus()}

	activeBackend, err = backend.CreateBackend(sdlbackend.NewSDLBackend())
	if err != nil {
		return err
	}

	activeBackend.SetBgColor(imgui.NewVec4(0.0, 0.0, 0.0, 1.0))

	activeBackend.CreateWindow("Wwise Teller", 1280, 720)

	imgui.CurrentIO().SetConfigFlags(imgui.ConfigFlagsDockingEnable)

	return nil
}

func loop() {
	state.eventBus.runAll()
	imgui.ClearSizeCallbackPool()
	imguizmo.BeginFrame()
	showFileExplorer()
}

func showFileExplorer() {
	imgui.Begin("File Explorer")

	imgui.Text(state.fileExplorer.pwd)

	if imgui.BeginTableV( "file-explorer-table",
		1,
		imgui.TableFlagsRowBg,
		imgui.NewVec2(0.0, 0.0), 0.0,
	) {
		imgui.TableSetupColumn("Name")
		imgui.TableHeadersRow()

		imgui.TableNextRow()
		imgui.TableSetColumnIndex(0)

		imgui.PushIDStr("file-explorer-table-cd-back")
		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)

		if imgui.Button("..") {
			state.eventBus.enqueue(
				"onCdParentClick",
				func () {
					err := state.fileExplorer.cdParent()
					if err != nil {
						slog.Error("Failed to go to parent directory", "error", err)
					}
				},
				state.ctx,
			)
		}

		imgui.PopStyleColor()
		imgui.PopID()

		for _, entry := range state.fileExplorer.entries {
			name := entry.Name()

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)

			imgui.PushIDStr("file-explorer-table-" + entry.Name() + "-entry")
			imgui.PushStyleColorVec4(
				imgui.ColButton, imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
			)

			if imgui.Button(name) {
				state.eventBus.enqueue(
					"onDirEntryClick",
					func() {
						if filepath.Ext(name) == ".bnk" {
							slog.Info(fmt.Sprintf("Loading sound bank %s", name))
							return
						}
						err := state.fileExplorer.cd(name)
						if err != nil {
							errMsg := fmt.Sprintf(
								"Failed to change to %s directory: %v", name, err,
							)
							slog.Error(errMsg)
						}
					},
					state.ctx,
				)
			}

			imgui.PopStyleColor()
			imgui.PopID()
		}
		imgui.EndTable()
	}

	imgui.End()
}
