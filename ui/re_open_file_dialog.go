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

func openFileDialogFunc(
	callback   func([]string),
	dirOnly    bool,
	initialDir string, 
	exts       []string,
) (func(), *bool, error) {
	d, err := fs.NewOpenFileDialog(callback, dirOnly, initialDir, exts)
	if err != nil {
		return nil, nil, err
	}
	done := false
	return func() {
		focusTable := false

		openFileDialogShortcut(d)
		openFileDialogVol(d)
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
		if imgui.Button("Open") {
			d.Open()
			done = true
			return
		}

		if imgui.ArrowButton("OpenFileDialogArrowButton", imgui.DirLeft) {
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
			focusTable = true
			imgui.SetKeyboardFocusHere()
		}

		openFileDialogTable(d, focusTable)
	}, &done, nil
}

func openFileDialogShortcut(d *fs.OpenFileDialog) {
	if isLeftShortcut() {
		if err := d.Parent(); err != nil {
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
}

func openFileDialogVol(d *fs.OpenFileDialog) {
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

func openFileDialogTable(d *fs.OpenFileDialog, focusTable bool) {
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
	storage := d.Storage
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
			entry := filtered[n].Entry

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
				if d.FocusDir(int(n)) {
					if err := d.CD(int(n)); err != nil {
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
