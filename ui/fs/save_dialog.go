package fs

import (
	"path/filepath"
)

// Single select only
// Operate only on directory
// Confirm (Ctrl - s) and cancel
// Change directory behavior:
//   - Double click or enter
// On save (Ctrl - s):
//   - If nothing is selected, use current directory on confirm
//   - If one is selected, use selected on confirm
type SaveFileDialog struct {
	Callback func(string)
	Selected int
	Fs       *FileSystem
}

func NewSaveFileDialog(callback func(string), initialDir string) (
	*SaveFileDialog, error,
) {
	fs, err := NewFileSystem(initialDir, true, nil)
	if err != nil {
		return nil, err
	}
	return &SaveFileDialog{Callback: callback, Selected: 0, Fs: fs}, nil
}

func (d *SaveFileDialog) CD() error {
	if d.Selected >= 0 && d.Selected < len(d.Fs.Filtered) {
		if err := d.Fs.CD(d.Fs.Filtered[d.Selected].Entry.Name()); err != nil {
			return err
		}
		d.Selected = 0
	}
	return nil
}

func (d *SaveFileDialog) Parent() error {
	if err := d.Fs.Parent(); err != nil {
		return err
	}
	d.Selected = 0
	return nil
}

func (d *SaveFileDialog) Filter() {
	d.Fs.Filter()
	d.Selected = 0
}

func (d *SaveFileDialog) Save() {
	if len(d.Fs.Filtered) <= 0 || d.Selected > len(d.Fs.Filtered) {
		d.Callback(d.Fs.Pwd)
	} else {
		d.Callback(filepath.Join(d.Fs.Pwd, d.Fs.Filtered[d.Selected].Entry.Name()))
	}
}

func (d *SaveFileDialog) SetNext(delta int) {
	if d.Selected + int(delta) >= 0 && d.Selected + int(delta) < len(d.Fs.Filtered) {
		d.Selected += int(delta)
	}
}

func (d *SaveFileDialog) SwitchVol(vol string) error {
	return d.Fs.SwitchVolume(vol)
}

func (d *SaveFileDialog) Vol() string {
	return d.Fs.Vol()
}
