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
	"github.com/Dekr0/wwise-teller/ui/fs"
)

func renderFileExplorer(fe *fs.FileExplorer) {
	imgui.BeginV("File Explorer", nil, 0)
	if imgui.BeginTabBar("FileExplorerTabBar") {
		renderFileExplorerTab(fe)
		imgui.EndTabBar()
	}
	imgui.End()
}

func renderFileExplorerTab(fe *fs.FileExplorer) {
	focusTable := false

	if !imgui.BeginTabItem("File Explorer") {
		return
	}

	setFileExplorerShortcut(fe)

	renderFileExplorerVol(fe)
	imgui.SameLine()
	imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyF),)
	imgui.SetNextItemWidth(160)
	if imgui.InputTextWithHint("Query", "", &fe.Fs.Query, 0, nil) {
		fe.Filter()
	}
	imgui.SameLine()
	imgui.SetNextItemShortcut(
		imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyS),
	)
	if imgui.Button("Open") {
		fe.OpenSelective()
	}

	imgui.SetNextItemShortcut(imgui.KeyChord(ModCtrlShift) | imgui.KeyChord(imgui.KeyN))
	if imgui.Button("+") {
		onOK := func(name string) {
			if name == "" { return }
			path := filepath.Join(fe.Pwd(), name)
			if err := os.MkdirAll(filepath.Join(fe.Pwd(), name), os.ModePerm); err != nil {
				slog.Error(
					fmt.Sprintf("Failed to create directory %s", path),
					"error", err,
				)
			}
			if err := fe.Refresh(); err != nil {
				slog.Error(
					fmt.Sprintf("Failed to refresh %s", fe.Fs.Pwd),
					"error", err,
				)
			}
		}
		pushSimpleTextModal("Make directory", onOK)
	}
	imgui.SameLine()
	imgui.SetNextItemShortcut(imgui.KeyChord(imgui.ModCtrl) | imgui.KeyChord(imgui.KeyR))
	if imgui.Button("R") {
		if err := fe.Refresh(); err != nil {
			slog.Error(
				fmt.Sprintf("Failed to refresh %s", fe.Fs.Pwd),
				"error", err,
			)
		}
	}
	imgui.SameLine()
	if imgui.ArrowButton("FileExplorerArrowButton", imgui.DirLeft) {
		if err := fe.Parent(); err != nil {
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

	renderFileExplorerTable(fe, focusTable)

	imgui.EndTabItem()
}

func setFileExplorerShortcut(fe *fs.FileExplorer) {
	if isLeftShortcut() {
		if err := fe.Parent(); err != nil {
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

func renderFileExplorerVol(fe *fs.FileExplorer) {
	if runtime.GOOS != "windows" {
		return
	}
	if len(utils.Vols) == 0 {
		return
	}

	vol := fe.Vol()
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
		if err := fe.SwitchVol(vol); err != nil {
			slog.Error("Failed to switch volume to " + vol, "error", err)
		}
	}
	imgui.PopItemWidth()
	imgui.PopID()
}

func renderFileExplorerTable(fe *fs.FileExplorer, focusTable bool) {
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
		if err := fe.Parent(); err != nil {
			slog.Error(
				"Failed to change current directory to parent directory",
				"error", err,
			)
		}
	}

	var cdFocus func() = nil
	var openFocus func() = nil

	filtered := fe.Fs.Filtered
	storage := fe.Storage
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
				if fe.IsFocusDir(int(n)) {
					cdFocus = bindCdFoucs(fe, n)
				} else {
					openFocus = bindOpenFoucs(fe, n)
				}
			}
		}
	}
	msIO = imgui.EndMultiSelect()
	storage.ApplyRequests(msIO)
	imgui.EndTable()
	if cdFocus != nil {
		cdFocus()
	}
	if openFocus != nil {
		openFocus()
	}
}

func bindCdFoucs(fe *fs.FileExplorer, n int32) func() {
	return func() {
		if err := fe.CD(int(n)); err != nil {
			slog.Error(
				"Failed to change current directory to selective directory",
				"error", err,
				)
		}
	}
}

func bindOpenFoucs(fe *fs.FileExplorer, n int32) func() {
	return func() {
		fe.OpenFocus(int(n))
	}
}
