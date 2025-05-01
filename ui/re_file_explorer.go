package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
)

func showFileExplorerWindow(fe *FileExplorer) {
	imgui.Begin("File Explorer")
	if imgui.BeginTabBar("FileExplorerTabBar") {
		showFileExplorerTab(fe)
		imgui.EndTabBar()
	}
	imgui.End()
}

func showFileExplorerTab(fe *FileExplorer) {
	focusTable := false

	if !imgui.BeginTabItem("File Explorer") {
		return
	}

	if isLeftShortcut() {
		if err := fe.CdParent(); err != nil {
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
	if imgui.InputTextWithHint("Query", "", &fe.fs.query, 0, nil) {
		fe.Filter()
	}

	imgui.SameLine()

	imgui.SetNextItemShortcut(
		imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS),
	)
	if imgui.Button("Open") {
		fe.OpenSelective()
	}

	if imgui.ArrowButton("FileExplorerArrowButton", imgui.DirLeft) {
		if err := fe.CdParent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}

	imgui.SameLine()

	imgui.Text(fe.Pwd())

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	showFileExplorerTabTable(fe, focusTable)

	imgui.EndTabItem()
}

func showFileExplorerTabTable(fe *FileExplorer, focusTable bool) {
	if !imgui.BeginTableV("FileExplorerTable",
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

	filtered := fe.fs.filtered
	storage := fe.storage
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
				if fe.IsFocusDir(int(n)) {
					if err := fe.CdFocus(int(n)); err != nil {
						slog.Error(
							"Failed to change current directory to selective directory",
							"error", err,
							)
					}
					break clipper
				} else {
					fe.OpenFocus(int(n))
				}
			}
		}
	}

	msIO = imgui.EndMultiSelect()
	storage.ApplyRequests(msIO)

	imgui.EndTable()
}
