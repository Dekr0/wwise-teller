// work in progress command struct for undo
// using pointer of the change state might be a bad idea to restore old state
// in the case of removing and creating new state.
package ui

import (
	"github.com/Dekr0/wwise-teller/wwise"
)

type removePropCmd struct {
	oldPid uint8
	oldValue []byte
	prop *wwise.PropBundle
}

func (r *removePropCmd) Do() {
	r.prop.Remove(r.oldPid)
}

func (r *removePropCmd) Undo() {
	r.prop.UpdatePropBytes(r.oldPid, r.oldValue)
}
