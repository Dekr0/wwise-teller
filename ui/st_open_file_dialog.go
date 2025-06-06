package ui

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
)

type OpenFileDialog struct {
	callback  func([]string)
	fs        *FileSystem
	storage   imgui.SelectionBasicStorage
}

func NewOpenFileDialog(
	callback   func([]string),
	dirOnly    bool,
	initialDir string, 
	exts       []string,
) (
	*OpenFileDialog, error,
) {
	fs, err := newFileSystem(initialDir, dirOnly, exts)
	if err != nil {
		return nil, err
	}
	return &OpenFileDialog{
		callback: callback,
		fs      : fs,
		storage : *imgui.NewSelectionBasicStorage(),
	}, nil
}

func (f *OpenFileDialog) cdFocus(n int) error {
	if n >= 0 && n < len(f.fs.filtered) {
		if err := f.fs.cd(f.fs.filtered[n].entry.Name()); err != nil {
			return err
		}
		f.resetSelection()
	}
	return nil
}

func (d *OpenFileDialog) cdParent() error {
	if err := d.fs.cdParent(); err != nil {
		return err
	}
	d.resetSelection()
	return nil
}

func (d *OpenFileDialog) filter() {
	d.fs.filter()
	d.resetSelection()
}

func (d *OpenFileDialog) isFocusDir(n int) bool {
	if n >= 0 && n < len(d.fs.filtered) {
		return d.fs.filtered[n].entry.IsDir()
	}
	return false
}

func (d *OpenFileDialog) openSelective() {
	paths := []string{}
	for i, e := range d.fs.filtered {
		if d.storage.Contains(imgui.ID(i))  {
			if !d.fs.dirOnly && e.entry.IsDir() {
				continue
			}
			paths = append(paths, filepath.Join(d.fs.pwd, e.entry.Name()))
		}
	}
	d.callback(paths)
}

func (d *OpenFileDialog) resetSelection() {
	d.storage.Clear()
}

func (d *OpenFileDialog) switchVol(vol string) error {
	return d.fs.switchVolume(vol)
}

func (d *OpenFileDialog) vol() string {
	return d.fs.vol()
}
