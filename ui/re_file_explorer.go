package ui

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/utils"
)

func showFileExplorerWindow(fe *FileExplorer, modalQ *ModalQ) {
	imgui.BeginV("File Explorer", nil, 0)
	if imgui.BeginTabBar("FileExplorerTabBar") {
		showFileExplorerTab(fe, modalQ)
		imgui.EndTabBar()
	}
	imgui.End()
}

func showFileExplorerTab(fe *FileExplorer, modalQ *ModalQ) {
	focusTable := false

	if !imgui.BeginTabItem("File Explorer") {
		return
	}

	setFileExplorerShortcut(fe)

	showFileExplorerVol(fe)
	imgui.SameLine()
	imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF),)
	if imgui.InputTextWithHint("Query", "", &fe.fs.query, 0, nil) {
		fe.filter()
	}
	imgui.SameLine()
	imgui.SetNextItemShortcut(
		imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS),
	)
	if imgui.Button("Open") {
		fe.openSelective()
	}

	imgui.SetNextItemShortcut(imgui.KeyChord(ModCtrlShift) | imgui.KeyChord(imgui.KeyN))
	if imgui.Button("+") {
		onOK := func(name string) {
			if name == "" { return }
			path := filepath.Join(fe.pwd(), name)
			if err := os.MkdirAll(filepath.Join(fe.pwd(), name), os.ModePerm); err != nil {
				slog.Error(
					fmt.Sprintf("Failed to create directory %s", path),
					"error", err,
				)
			}
			if err := fe.refresh(); err != nil {
				slog.Error(
					fmt.Sprintf("Failed to refresh %s", fe.fs.pwd),
					"error", err,
				)
			}
		}
		pushSimpleTextModal(modalQ, "Make directory", onOK)
	}
	imgui.SameLine()
	imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyR))
	if imgui.Button("R") {
		if err := fe.refresh(); err != nil {
			slog.Error(
				fmt.Sprintf("Failed to refresh %s", fe.fs.pwd),
				"error", err,
			)
		}
	}
	imgui.SameLine()
	if imgui.ArrowButton("FileExplorerArrowButton", imgui.DirLeft) {
		if err := fe.cdParent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}
	imgui.SameLine()
	imgui.Text(fe.pwd())

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	showFileExplorerTabTable(fe, focusTable)

	imgui.EndTabItem()
}

func setFileExplorerShortcut(fe *FileExplorer) {
	if isLeftShortcut() {
		if err := fe.cdParent(); err != nil {
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

func showFileExplorerVol(fe *FileExplorer) {
	if runtime.GOOS != "windows" {
		return
	}
	if len(utils.Vols) == 0 {
		return
	}

	vol := fe.vol()
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
		if err := fe.switchVol(vol); err != nil {
			slog.Error("Failed to switch volume to " + vol, "error", err)
		}
	}
	imgui.PopItemWidth()
	imgui.PopID()
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

	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	imgui.SelectableBool("..")
	focused := imgui.IsItemFocused()
	doubleClicked := imgui.IsMouseDoubleClicked(0)
	righted := isRightShortcut()
	if focused && (doubleClicked || righted) {
		if err := fe.cdParent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}

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
				if fe.isFocusDir(int(n)) {
					if err := fe.cdFocus(int(n)); err != nil {
						slog.Error(
							"Failed to change current directory to selective directory",
							"error", err,
							)
					}
					break clipper
				} else {
					fe.openFocus(int(n))
				}
			}
		}
	}

	msIO = imgui.EndMultiSelect()
	storage.ApplyRequests(msIO)

	imgui.EndTable()
}
