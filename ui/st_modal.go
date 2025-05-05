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
	modals []*Modal
}

func NewModalQ() *ModalQ {
	return &ModalQ{make([]*Modal, 0, 4)}
}

func (m *ModalQ) QModal(
	done *bool,
	flag imgui.WindowFlags,
	name string,
	loop func(),
	onClose func(),
) {
	if slices.ContainsFunc(m.modals, func(modal *Modal) bool {
		return strings.Compare(modal.name, name) == 0
	}) {
		return
	}
	m.modals = append(m.modals, &Modal{done, flag, name, loop, onClose})
}
