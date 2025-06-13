package ui

import (
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/fs"
	"github.com/Dekr0/wwise-teller/utils"
)

func saveFileDialogFunc(
	callback   func(string),
	initialDir string,
) (func(), *bool, error) {
	done := false
	d, err := fs.NewSaveFileDialog(callback, initialDir)
	if err != nil {
		return nil, nil, err
	}
	return func() {
		focusTable := false

		saveFileDialogShortcut(d)
		saveFileDialogVol(d)
		imgui.SameLine()
		imgui.SetNextItemShortcut(
			imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF),
		)
		if imgui.InputTextWithHint("Query", "", &d.Fs.Query, 0, nil) {
			d.Filter()
		}
		imgui.SameLine()
		align := imgui.CursorPos().X
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS))
		if imgui.Button("Save") {
			d.Save()
			done = true
			return
		}

		if imgui.ArrowButton("SaveFileDialogArrowButton", imgui.DirLeft) {
			if err := d.Parent(); err != nil {
				slog.Error(
					"Failed to change current directory to parent directory",
					"error", err,
				)
			}
		}
		imgui.SameLine()
		imgui.Text(d.Fs.Pwd)
		imgui.SameLine()
		imgui.SetCursorPosX(align)
		if imgui.Button("Cancel") {
			done = true
			return
		}

		if imgui.Shortcut(UnFocusQuerySC) {
			imgui.SetKeyboardFocusHere()
			focusTable = true
		}

		saveFileDialogTable(d, focusTable)
	}, &done, nil
}

func saveFileDialogShortcut(d *fs.SaveFileDialog) {
	if isLeftShortcut() {
		if err := d.Parent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}
	useViUp()
	useViDown()
}

func saveFileDialogVol(d *fs.SaveFileDialog) {
	if runtime.GOOS != "windows" {
		return
	}
	if len(utils.Vols) == 0 {
		return
	}

	vol := d.Vol()
	idx := int32(slices.IndexFunc(utils.Vols, func(v string) bool {
		return strings.Compare(v, vol) == 0
	}))
	if idx == -1 {
		idx = 0
	}

	imgui.PushIDStr("FileExplorerVol")
	imgui.PushItemWidth(imgui.CalcTextSize("C:\\").X + 24.0)
	if imgui.ComboStrarr("", &idx, utils.Vols, int32(len(utils.Vols))) {
		vol := utils.Vols[idx]
		if err := d.SwitchVol(vol); err != nil {
			slog.Error("Failed to switch volume to " + vol, "error", err)
		}
	}
	imgui.PopItemWidth()
	imgui.PopID()
}

func saveFileDialogTable(d *fs.SaveFileDialog, focusTable bool) {
	if !imgui.BeginTableV("SaveFileDialogTable",
		1, imgui.TableFlagsRowBg | imgui.TableFlagsScrollY,
		imgui.NewVec2(0.0, 0.0), 0.0,
	) {
		return
	}
	imgui.TableSetupColumn("File / Directory name")
	imgui.TableHeadersRow()

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)

	if focusTable {
		imgui.SetKeyboardFocusHere()
	}
	imgui.SelectableBool(".")

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	imgui.SelectableBool("..")
	focused := imgui.IsItemFocused()
	doubleClicked := imgui.IsMouseDoubleClicked(0)
	righted := isRightShortcut()
	if focused && (doubleClicked || righted) {
		if err := d.Parent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}

	filtered := d.Fs.Filtered
	clipper := imgui.NewListClipper()
	clipper.Begin(int32(len(filtered)))
	clipper:
	for clipper.Step() {
		for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
			entry := filtered[n].Entry

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			label := fmt.Sprintf("%s/", entry.Name())
			if imgui.SelectableBoolV(
				label, n == int32(d.Selected), 
				imgui.SelectableFlagsSpanAllColumns | 
				imgui.SelectableFlagsAllowOverlap,
				imgui.NewVec2(0, 0),
			) {
				d.Selected = int(n)
			}

			focused := imgui.IsItemFocused()
			if focused {
				d.Selected = int(n)
			}

			doubleClicked := imgui.IsMouseDoubleClicked(0)
			righted := isRightShortcut() 
			if focused && (doubleClicked || righted) {
				if err := d.CD(); err != nil {
					slog.Error(
						"Failed to change current directory selected directory",
						"error", err,
					)
				}
				break clipper
			}
		}
	}

	imgui.EndTable()
}
