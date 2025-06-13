package fs

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
)

type OpenFileDialog struct {
	Callback  func([]string)
	Fs        *FileSystem
	Storage   imgui.SelectionBasicStorage
}

func NewOpenFileDialog(
	callback   func([]string),
	dirOnly    bool,
	initialDir string, 
	exts       []string,
) (
	*OpenFileDialog, error,
) {
	fs, err := NewFileSystem(initialDir, dirOnly, exts)
	if err != nil {
		return nil, err
	}
	return &OpenFileDialog{
		Callback: callback,
		Fs      : fs,
		Storage : *imgui.NewSelectionBasicStorage(),
	}, nil
}

func (f *OpenFileDialog) CD(n int) error {
	if n >= 0 && n < len(f.Fs.Filtered) {
		if err := f.Fs.CD(f.Fs.Filtered[n].Entry.Name()); err != nil {
			return err
		}
		f.ResetSelection()
	}
	return nil
}

func (d *OpenFileDialog) Parent() error {
	if err := d.Fs.Parent(); err != nil {
		return err
	}
	d.ResetSelection()
	return nil
}

func (d *OpenFileDialog) Filter() {
	d.Fs.Filter()
	d.ResetSelection()
}

func (d *OpenFileDialog) FocusDir(n int) bool {
	if n >= 0 && n < len(d.Fs.Filtered) {
		return d.Fs.Filtered[n].Entry.IsDir()
	}
	return false
}

func (d *OpenFileDialog) Open() {
	paths := []string{}
	for i, e := range d.Fs.Filtered {
		if d.Storage.Contains(imgui.ID(i))  {
			if !d.Fs.DirOnly && e.Entry.IsDir() {
				continue
			}
			paths = append(paths, filepath.Join(d.Fs.Pwd, e.Entry.Name()))
		}
	}
	d.Callback(paths)
}

func (d *OpenFileDialog) ResetSelection() {
	d.Storage.Clear()
}

func (d *OpenFileDialog) SwitchVol(vol string) error {
	return d.Fs.SwitchVolume(vol)
}

func (d *OpenFileDialog) Vol() string {
	return d.Fs.Vol()
}
