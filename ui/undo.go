package ui

import "slices"

type undoNode struct {
	name string
	do   func()
	undo func()
}

type undoList struct {
	curr uint
	nodes []*undoNode
}

func newUndoList() *undoList {
	nodes := make([]*undoNode, 0, 32)
	nodes = append(nodes, &undoNode{"Initial State", nil, nil})
	return &undoList{0, nodes}
}

func (u *undoList) Do(name string, do func(), undo func()) {
	if u.curr == uint(len(u.nodes) - 1) {
		u.nodes = append(u.nodes, &undoNode{name, do, undo})
	} else {
		u.nodes = slices.Delete(u.nodes, int(u.curr) + 1, len(u.nodes))
		u.nodes = append(u.nodes, &undoNode{name, do, undo})
	}
	do()
	u.curr += 1
}

func (u *undoList) Undo() {
	if u.curr <= 0 {
		return
	}
	u.nodes[u.curr].undo()
	u.curr -= 1
}

func (u *undoList) Redo() {
	if u.curr >= uint(len(u.nodes) - 1) {
		return
	}
	u.curr += 1
	u.nodes[u.curr].do()
}
