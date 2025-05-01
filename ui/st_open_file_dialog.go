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

func (d *OpenFileDialog) ResetSelection() {
	d.storage.Clear()
}

func (d *OpenFileDialog) Filter() {
	d.fs.filter()
	d.ResetSelection()
}

func (d *OpenFileDialog) IsFocusDir(n int) bool {
	if n >= 0 && n < len(d.fs.filtered) {
		return d.fs.filtered[n].entry.IsDir()
	}
	return false
}

func (f *OpenFileDialog) CdFocus(n int) error {
	if n >= 0 && n < len(f.fs.filtered) {
		if err := f.fs.cd(f.fs.filtered[n].entry.Name()); err != nil {
			return err
		}
		f.ResetSelection()
	}
	return nil
}

func (d *OpenFileDialog) CdParent() error {
	if err := d.fs.cdParent(); err != nil {
		return err
	}
	d.ResetSelection()
	return nil
}

func (d *OpenFileDialog) OpenSelective() {
	paths := []string{}
	for i, e := range d.fs.filtered {
		if d.storage.Contains(imgui.ID(i))  {
			if !d.fs.dirOnly && e.entry.IsDir() {
				continue
			}
			paths = append(paths, filepath.Join(d.fs.Pwd, e.entry.Name()))
		}
	}
	d.callback(paths)
}
