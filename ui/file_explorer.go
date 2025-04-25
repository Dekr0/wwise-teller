package ui

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type openFile struct {
	pwd string
	query string
	onOpen func(string)
	entries []os.DirEntry 
	filtered []os.DirEntry
}

func newOpenFileAtHome(onOpen func(string)) (*openFile, error) {
	home, err := utils.GetHome()
	if err != nil { return nil, err }
	return newOpenFile(home, onOpen)
}

func newOpenFile(p string, onOpen func(string)) (*openFile, error) {
	entries, err := utils.GetDirAndBank(p, false)
	if err != nil {
		return nil, err
	}
	filtered := make([]os.DirEntry, len(entries))
	for i, e := range entries {
		filtered[i] = e 
	}
	return &openFile{p, "", onOpen, entries, filtered}, nil
}

func (f *openFile) filter() {
	i := 0
	old := len(f.filtered)
	for _, d := range f.entries {
		if !fuzzy.Match(f.query, d.Name()) {
			continue
		}
		if i < len(f.filtered) {
			f.filtered[i] = d
		} else {
			f.filtered = append(f.filtered, d)
		}
		i += 1
	}
	if i < old {
		f.filtered = slices.Delete(f.filtered, i, old)
	}
}

func (f *openFile) cdParent() error {
	pwd := filepath.Dir(f.pwd)

	if runtime.GOOS == "windows" && pwd == "." {
		return nil
	}
	
	if pwd == f.pwd {
		return nil
	}

	newEntries, err := utils.GetDirAndBank(pwd, false)
	if err != nil {
		return err
	}

	f.pwd = pwd
	f.entries = newEntries

	f.filter()

	return nil
}

func (f *openFile) cd(basename string) error {
	pwd := filepath.Join(f.pwd, basename)

	newEntries, err := utils.GetDirAndBank(pwd, false)
	if err != nil {
		return err
	}

	f.pwd = pwd
	f.entries = newEntries

	f.filter()

	return nil
}

func showFileExplorer(openFile *openFile) string {
	bookmark := ""

	imgui.Begin("File Explorer")

	if imgui.BeginTabBar("FileExplorerTabBar") {
		bookmark = showFileExplorerTab(openFile)
		showWorkspaceExplorerTab()
		imgui.EndTabBar()
	}

	imgui.End()

	return bookmark
}

func showFileExplorerTab(openFile *openFile) string {
	bookmark := ""
	if !imgui.BeginTabItem("File Explorer") {
		return bookmark
	}

	if imgui.InputTextWithHint("Query", "", &openFile.query, 0, nil) {
		openFile.filter()
	}

	if imgui.BeginTableV("FileExplorerTable",
		1, imgui.TableFlagsRowBg | imgui.TableFlagsScrollY,
		imgui.NewVec2(0.0, 0.0), 0.0,
	) {
		imgui.TableSetupColumn("File name")
		imgui.TableHeadersRow()

		imgui.TableNextRow()
		imgui.TableSetColumnIndex(0)

		imgui.PushIDStr("FileExplorerTableCdParent")
		imgui.PushStyleVarVec2(
			imgui.StyleVarButtonTextAlign,
			imgui.Vec2{X: 0.0, Y: 0.5},
		)
		imgui.PushStyleColorVec4(
			imgui.ColButton,
			imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
		)

		if imgui.ButtonV("..", imgui.Vec2{X: -1, Y: 0}) {
			err := openFile.cdParent()
			if err != nil {
				slog.Error("Failed to go to parent directory", "error", err)
			}
		}

		imgui.PopStyleColor()
		imgui.PopStyleVar()
		imgui.PopID()

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(openFile.filtered)))

		clipper:
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				if showOpenFileEntry(openFile, openFile.filtered[n]) {
					break clipper
				}
				if imgui.BeginPopupContextItem() {
					if imgui.Button("Add to asset sourcing list") {
						bookmark = filepath.Join(openFile.pwd, openFile.filtered[n].Name())
						imgui.CloseCurrentPopup()
					}
					imgui.EndPopup()
				}
			}
		}

		imgui.EndTable()
	}
	imgui.EndTabItem()
	return bookmark
}

func showOpenFileEntry(openFile *openFile, entry os.DirEntry) bool {
	name := entry.Name()
	
	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	
	imgui.PushIDStr("FileExplorerTableCd" + entry.Name())
	imgui.PushStyleVarVec2(
		imgui.StyleVarButtonTextAlign,
		imgui.Vec2{X: 0.0, Y: 0.5},
	)
	imgui.PushStyleColorVec4(
		imgui.ColButton, imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
	)
	
	if imgui.ButtonV(name, imgui.Vec2{X: -1, Y: 0}) {
		name := entry.Name()
		if entry.IsDir() {
			err := openFile.cd(name)
			if err != nil {
				errMsg := fmt.Sprintf(
					"Failed to change to %s directory: %v", name, err,
					)
				slog.Error(errMsg)
			} else {
				openFile.query = ""
				openFile.filter()
			}
		} else {
			openFile.onOpen(filepath.Join(openFile.pwd, name))
		}
		imgui.PopStyleColor()
		imgui.PopStyleVar()
		imgui.PopID()
		return true
	}
	imgui.PopStyleColor()
	imgui.PopStyleVar()
	imgui.PopID()
	return false
}

func showWorkspaceExplorerTab() {
	if !imgui.BeginTabItem("Workspace") {
		return
	}
	imgui.TextDisabled("Workspace is currently under construction.")
	imgui.EndTabItem()
}

type saveFile struct {
	pwd string
	query string
	dest string
	onSave func(string)
	entries []os.DirEntry 
	filtered []os.DirEntry
}

func newSaveFileAtHome(onSave func(string)) (*saveFile, error) {
	home, err := utils.GetHome()
	if err != nil { return nil, err }
	return newSaveFile(home, onSave)
}

func newSaveFile(p string, onSave func(string)) (*saveFile, error) {
	entries, err := utils.GetDirAndBank(p, true)
	if err != nil {
		return nil, err
	}
	filtered := make([]os.DirEntry, len(entries))
	for i, e := range entries {
		filtered[i] = e 
	}
	return &saveFile{p, "", "", onSave, entries, filtered}, nil
}

func (s *saveFile) filter() {
	i := 0
	old := len(s.filtered)
	for _, d := range s.entries {
		if !fuzzy.Match(s.query, d.Name()) {
			continue
		}
		if i < len(s.filtered) {
			s.filtered[i] = d
		} else {
			s.filtered = append(s.filtered, d)
		}
		i += 1
	}
	if i < old {
		s.filtered = slices.Delete(s.filtered, i, old)
	}
}

func (s *saveFile) cdParent() error {
	pwd := filepath.Dir(s.pwd)

	if runtime.GOOS == "windows" && pwd == "." {
		return nil
	}
	
	if pwd == s.pwd {
		return nil
	}

	newEntries, err := utils.GetDirAndBank(pwd, true)
	if err != nil {
		return err
	}

	s.pwd = pwd
	s.entries = newEntries

	s.filter()

	return nil
}

func (s *saveFile) cd(basename string) error {
	pwd := filepath.Join(s.pwd, basename)

	newEntries, err := utils.GetDirAndBank(pwd, true)
	if err != nil {
		return err
	}

	s.pwd = pwd
	s.entries = newEntries

	s.filter()

	return nil
}

func showSaveFileModal(saveFile *saveFile, showSaveFile bool) {
	if showSaveFile {
		if !imgui.IsPopupOpenStr("Save Sound Bank") {
			imgui.OpenPopupStr("Save Sound Bank")
		}
	}

	if imgui.BeginPopupModal("Save Sound Bank") {
		if imgui.InputTextWithHint("Dest.", "", &saveFile.dest, 0, nil) {}
		imgui.SameLine()
		if imgui.Button("Save") {
			if saveFile.dest != "" {
				if saveFile.onSave != nil {
					saveFile.onSave(filepath.Join(saveFile.pwd, saveFile.dest))
					saveFile.onSave = nil
				}
				imgui.CloseCurrentPopup()
			}
		}
		if saveFile.dest == "" {
			imgui.Text("File name cannot be empty")
		}

		if imgui.InputTextWithHint("Query", "", &saveFile.query, 0, nil) {
			saveFile.filter()
		}
		imgui.SameLine()
		if imgui.Button("Cancel") {
			saveFile.onSave = nil
			imgui.CloseCurrentPopup()
		}

		if imgui.BeginTableV("FileExplorerTable",
			1, imgui.TableFlagsRowBg | imgui.TableFlagsScrollY,
			imgui.NewVec2(0.0, 0.0), 0.0,
		) {
			imgui.TableSetupColumn("File name")
			imgui.TableHeadersRow()

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)

			imgui.PushIDStr("FileExplorerTableCdParent")
			imgui.PushStyleVarVec2(
				imgui.StyleVarButtonTextAlign,
				imgui.Vec2{X: 0.0, Y: 0.5},
			)
			imgui.PushStyleColorVec4(
				imgui.ColButton,
				imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
			)

			if imgui.ButtonV("..", imgui.Vec2{X: -1, Y: 0}) {
				err := saveFile.cdParent()
				if err != nil {
					slog.Error(
						"Failed to go to parent directory", 
						"error", err,
					)
				}
			}
			imgui.PopStyleColor()
			imgui.PopStyleVar()
			imgui.PopID()

			clipper := imgui.NewListClipper()
			clipper.Begin(int32(len(saveFile.filtered)))
			clipper:
			for clipper.Step() {
				for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
					if showSaveFileEntry(saveFile, saveFile.filtered[n]) {
						break clipper
					}
				}
			}
			imgui.EndTable()
		}

		imgui.EndPopup()
	}
}

func showSaveFileEntry(saveFile *saveFile, entry os.DirEntry) bool {
	name := entry.Name()
	
	imgui.TableNextRow()
	imgui.TableSetColumnIndex(0)
	
	imgui.PushIDStr("FileExplorerTableCd" + entry.Name())
	imgui.PushStyleVarVec2(
		imgui.StyleVarButtonTextAlign,
		imgui.Vec2{X: 0.0, Y: 0.5},
	)
	imgui.PushStyleColorVec4(
		imgui.ColButton, imgui.Vec4{X: 0.0, Y: 0.0, Z: 0.0, W: 0.0},
	)
	
	if imgui.ButtonV(name, imgui.Vec2{X: -1, Y: 0}) {
		err := saveFile.cd(name)
		if err != nil {
			errMsg := fmt.Sprintf(
				"Failed to change to %s directory: %v", name, err,
				)
			slog.Error(errMsg)
		} else {
			saveFile.query = ""
			saveFile.filter()
		}
		imgui.PopStyleColor()
		imgui.PopStyleVar()
		imgui.PopID()
		return true
	}
	
	imgui.PopStyleColor()
	imgui.PopStyleVar()
	imgui.PopID()
	return false
}
