package ui

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
)

type FileExplorer struct {
	fs       *FileSystem
	callback func([]string)
	storage  imgui.SelectionBasicStorage
}

func newFileExplorer(callback func([]string), initialDir string) (
	*FileExplorer, error,
) {
	fs, err := newFileSystem(initialDir, false, []string{".bnk", ".st_bnk"})
	if err != nil {
		return nil, err
	}
	return &FileExplorer{
		fs:       fs,
		callback: callback,
		storage:  *imgui.NewSelectionBasicStorage(),
	}, nil
}

func (f *FileExplorer) cdFocus(n int) error {
	if n >= 0 && n < len(f.fs.filtered) {
		if err := f.fs.cd(f.fs.filtered[n].entry.Name()); err != nil {
			return err
		}
		f.resetSelection()
	}
	return nil
}

func (f *FileExplorer) cdParent() error {
	if err := f.fs.cdParent(); err != nil {
		return err
	}
	f.resetSelection()
	return nil
}

func (f *FileExplorer) filter() {
	f.fs.filter()
	f.resetSelection()
}

func (f *FileExplorer) isFocusDir(n int) bool {
	if n >= 0 && n < len(f.fs.filtered) {
		return f.fs.filtered[n].entry.IsDir()
	}
	return false
}

func (f *FileExplorer) openFocus(n int) {
	if n >= 0 && n < len(f.fs.filtered) {
		path := filepath.Join(f.fs.pwd, f.fs.filtered[n].entry.Name())
		f.callback([]string{path})
	}
}

func (f *FileExplorer) openSelective() {
	paths := []string{}
	for i, d := range f.fs.filtered {
		if f.storage.Contains(imgui.ID(i)) && !d.entry.IsDir() {
			paths = append(paths, filepath.Join(f.pwd(), d.entry.Name()))
		}
	}
	if len(paths) > 0 {
		f.callback(paths)
		f.resetSelection()
	}
}

func (f *FileExplorer) pwd() string {
	return f.fs.pwd
}

func (f *FileExplorer) resetSelection() {
	f.storage.Clear()
}

func (f *FileExplorer) switchVol(vol string) error {
	return f.fs.switchVolume(vol)
}

func (f *FileExplorer) vol() string {
	return f.fs.vol()
}
