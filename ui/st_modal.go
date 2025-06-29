package ui

import (
	"slices"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
)

type Modal struct {
	done *bool
	flag imgui.WindowFlags
	name string
	loop func()
	onClose func()
}

type ModalQ struct {
	Modals []*Modal
}

func NewModalQ() ModalQ {
	return ModalQ{make([]*Modal, 0, 4)}
}

func (m *ModalQ) QModal(
	done *bool,
	flag imgui.WindowFlags,
	name string,
	loop func(),
	onClose func(),
) {
	if slices.ContainsFunc(m.Modals, func(modal *Modal) bool {
		return strings.Compare(modal.name, name) == 0
	}) {
		return
	}
	m.Modals = append(m.Modals, &Modal{done, flag, name, loop, onClose})
}
