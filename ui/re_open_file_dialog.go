package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
)

func openFileDialogFunc(
	callback   func([]string),
	dirOnly    bool,
	initialDir string, 
	exts       []string,
) (func(), *bool, error) {
	d, err := NewOpenFileDialog(callback, dirOnly, initialDir, exts)
	if err != nil {
		return nil, nil, err
	}
	done := false
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
		useViShiftUp()
		useViDown()
		useViShiftDown()

		imgui.SetNextItemShortcut(
			imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF),
		)
		if imgui.InputTextWithHint("Query", "", &d.fs.query, 0, nil) {
			d.Filter()
		}

		imgui.SameLine()

		align := imgui.CursorPos().X
		imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS))
		if imgui.Button("Open") {
			d.OpenSelective()
			done = true
			return
		}

		if imgui.ArrowButton("OpenFileDialogArrowButton", imgui.DirLeft) {
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
			focusTable = true
			imgui.SetKeyboardFocusHere()
		}

		openFileDialogTable(d, focusTable)
	}, &done, nil
}

func openFileDialogTable(d *OpenFileDialog, focusTable bool) {
	if !imgui.BeginTableV("OpenFileDialogTable",
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
	storage := d.storage
	msIO := imgui.BeginMultiSelectV(
		imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d,
		storage.Size(),
		int32(len(filtered)),
	)
	storage.ApplyRequests(msIO)

	clipper := imgui.NewListClipper()
	clipper.Begin(int32(len(filtered)))
	if msIO.RangeSrcItem() != 1 {
		clipper.IncludeItemByIndex(int32(msIO.RangeSrcItem()))
	}
	clipper:
	for clipper.Step() {
		for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
			entry := filtered[n].entry

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)

			label := entry.Name()
			if entry.IsDir() {
				label = fmt.Sprintf("%s/", label)
			}

			selected := storage.Contains(imgui.ID(n))
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(n))
			imgui.SelectableBoolPtrV(
				label,
				&selected,
				imgui.SelectableFlagsSpanAllColumns | 
				imgui.SelectableFlagsAllowOverlap,
				imgui.Vec2{X: 0, Y: 0.0},
			)

			focused := imgui.IsItemFocused()
			doubleClicked := imgui.IsMouseDoubleClicked(0)
			righted := isRightShortcut()
			if focused && (doubleClicked || righted) {
				if d.IsFocusDir(int(n)) {
					if err := d.CdFocus(int(n)); err != nil {
						slog.Error(
							"Failed to change current directory to selective directory",
							"error", err,
							)
					}
					break clipper
				}
			}
		}
	}

	msIO = imgui.EndMultiSelect()
	storage.ApplyRequests(msIO)

	imgui.EndTable()
}
