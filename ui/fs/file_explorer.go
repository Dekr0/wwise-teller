package fs

import (
	"path/filepath"

	"github.com/AllenDang/cimgui-go/imgui"
)

type FileExplorer struct {
	Fs       *FileSystem
	Callback func([]string)
	Storage  imgui.SelectionBasicStorage
}

func NewFileExplorer(callback func([]string), initialDir string) (
	*FileExplorer, error,
) {
	fs, err := NewFileSystem(initialDir, false, []string{".bnk", ".st_bnk"})
	if err != nil {
		return nil, err
	}
	return &FileExplorer{
		Fs:       fs,
		Callback: callback,
		Storage:  *imgui.NewSelectionBasicStorage(),
	}, nil
}

func (f *FileExplorer) CD(n int) error {
	if n >= 0 && n < len(f.Fs.Filtered) {
		if err := f.Fs.CD(f.Fs.Filtered[n].Entry.Name()); err != nil {
			return err
		}
		f.ResetSelection()
	}
	return nil
}

func (f *FileExplorer) Parent() error {
	if err := f.Fs.Parent(); err != nil {
		return err
	}
	f.ResetSelection()
	return nil
}

func (f *FileExplorer) Filter() {
	f.Fs.Filter()
	f.ResetSelection()
}

func (f *FileExplorer) IsFocusDir(n int) bool {
	if n >= 0 && n < len(f.Fs.Filtered) {
		return f.Fs.Filtered[n].Entry.IsDir()
	}
	return false
}

func (f *FileExplorer) OpenFocus(n int) {
	if n >= 0 && n < len(f.Fs.Filtered) {
		path := filepath.Join(f.Fs.Pwd, f.Fs.Filtered[n].Entry.Name())
		f.Callback([]string{path})
	}
}

func (f *FileExplorer) OpenSelective() {
	paths := []string{}
	for i, d := range f.Fs.Filtered {
		if f.Storage.Contains(imgui.ID(i)) && !d.Entry.IsDir() {
			paths = append(paths, filepath.Join(f.Pwd(), d.Entry.Name()))
		}
	}
	if len(paths) > 0 {
		f.Callback(paths)
		f.ResetSelection()
	}
}

func (f *FileExplorer) Pwd() string {
	return f.Fs.Pwd
}

func (f *FileExplorer) ResetSelection() {
	f.Storage.Clear()
}

func (f *FileExplorer) Refresh() error {
	return f.Fs.Refresh()
}

func (f *FileExplorer) SwitchVol(vol string) error {
	return f.Fs.SwitchVolume(vol)
}

func (f *FileExplorer) Vol() string {
	return f.Fs.Vol()
}
