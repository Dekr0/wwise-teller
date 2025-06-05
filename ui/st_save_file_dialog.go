// TODO:
// - Modify the destination file name
package ui

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
	callback func(string)
	selected int
	fs       *FileSystem
}

func NewSaveFileDialog(callback func(string), initialDir string) (
	*SaveFileDialog, error,
) {
	fs, err := newFileSystem(initialDir, true, nil)
	if err != nil {
		return nil, err
	}
	return &SaveFileDialog{callback: callback, selected: 0, fs: fs}, nil
}

func (d *SaveFileDialog) cdSelected() error {
	if d.selected >= 0 && d.selected < len(d.fs.filtered) {
		if err := d.fs.cd(d.fs.filtered[d.selected].entry.Name()); err != nil {
			return err
		}
		d.selected = 0
	}
	return nil
}

func (d *SaveFileDialog) cdParent() error {
	if err := d.fs.cdParent(); err != nil {
		return err
	}
	d.selected = 0
	return nil
}

func (d *SaveFileDialog) filter() {
	d.fs.filter()
	d.selected = 0
}

func (d *SaveFileDialog) save() {
	if len(d.fs.filtered) <= 0 || d.selected > len(d.fs.filtered) {
		d.callback(d.fs.pwd)
	} else {
		d.callback(filepath.Join(d.fs.pwd, d.fs.filtered[d.selected].entry.Name()))
	}
}

func (d *SaveFileDialog) setNext(delta int) {
	if d.selected + int(delta) >= 0 && d.selected + int(delta) < len(d.fs.filtered) {
		d.selected += int(delta)
	}
}

func (d *SaveFileDialog) switchVol(vol string) error {
	return d.fs.switchVolume(vol)
}

func (d *SaveFileDialog) vol() string {
	return d.fs.vol()
}
