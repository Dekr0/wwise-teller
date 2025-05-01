package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
)

func saveFileDialogFunc(
	callback   func(string),
	initialDir string,
) (func(), *bool, error) {
	done := false
	d, err := NewSaveFileDialog(callback, initialDir)
	if err != nil {
		return nil, nil, err
	}
	return func() {
		focusTable := false

		if isLeftShortcut() {
			if err := d.CdParent(); err != nil {
				slog.Error(
					"Failed to change current directory to parent directory",
					"error", err,
				)
			}
		}
		useViUp()
		useViDown()

		imgui.SetNextItemShortcut(
			imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF),
		)
		if imgui.InputTextWithHint("Query", "", &d.fs.query, 0, nil) {
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
			if err := d.CdParent(); err != nil {
				slog.Error(
					"Failed to change current directory to parent directory",
					"error", err,
				)
			}
		}

		imgui.SameLine()

		imgui.Text(d.fs.Pwd)

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

func saveFileDialogTable(d *SaveFileDialog, focusTable bool) {
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

	filtered := d.fs.filtered
	clipper := imgui.NewListClipper()
	clipper.Begin(int32(len(filtered)))
	clipper:
	for clipper.Step() {
		for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
			entry := filtered[n].entry

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			label := fmt.Sprintf("%s/", entry.Name())
			if imgui.SelectableBoolV(
				label, n == int32(d.selected), 
				imgui.SelectableFlagsSpanAllColumns | 
				imgui.SelectableFlagsAllowOverlap,
				imgui.NewVec2(0, 0),
			) {
				d.selected = int(n)
			}

			focused := imgui.IsItemFocused()
			if focused {
				d.selected = int(n)
			}

			doubleClicked := imgui.IsMouseDoubleClicked(0)
			righted := isRightShortcut() 
			if focused && (doubleClicked || righted) {
				if err := d.CdSelected(); err != nil {
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
