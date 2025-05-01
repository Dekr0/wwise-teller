package ui

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
)

type FileExplorer struct {
	fs         *FileSystem
	callback   func([]string)
	storage    imgui.SelectionBasicStorage
}

func NewFileExplorer(callback func([]string), initialDir string) (
	*FileExplorer, error,
) {
	fs, err := newFileSystem(initialDir, false, []string{".bnk", ".st_bnk"})
	if err != nil {
		return nil, err
	}
	return &FileExplorer{
		fs      : fs,
		callback: callback,
		storage : *imgui.NewSelectionBasicStorage(),
	}, nil
}

func (f *FileExplorer) Pwd() string {
	return f.fs.Pwd
}

func (f *FileExplorer) ResetSelection() {
	f.storage.Clear()
}

func (f *FileExplorer) Filter() {
	f.fs.filter()
	f.ResetSelection()
}

func (f *FileExplorer) IsFocusDir(n int) bool {
	if n >= 0 && n < len(f.fs.filtered) {
		return f.fs.filtered[n].entry.IsDir()
	}
	return false
}

func (f *FileExplorer) CdFocus(n int) error {
	if n >= 0 && n < len(f.fs.filtered) {
		if err := f.fs.cd(f.fs.filtered[n].entry.Name()); err != nil {
			return err
		}
		f.ResetSelection()
	}
	return nil
}

func (f *FileExplorer) OpenFocus(n int) {
	if n >= 0 && n < len(f.fs.filtered) {
		path := filepath.Join(f.fs.Pwd, f.fs.filtered[n].entry.Name())
		f.callback([]string{path})
	}
}

func (f *FileExplorer) OpenSelective() {
	paths := []string{}
	for i, d := range f.fs.filtered {
		if f.storage.Contains(imgui.ID(i)) && !d.entry.IsDir() {
			paths = append(paths, filepath.Join(f.Pwd(), d.entry.Name()))
		}
	}
	if len(paths) > 0 {
		f.callback(paths)
		f.ResetSelection()
	}
}

func (f *FileExplorer) CdParent() error {
	if err := f.fs.cdParent(); err != nil {
		return err
	}
	f.ResetSelection()
	return nil
}
