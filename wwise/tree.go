package wwise

import "fmt"

type ActorMixerHircNode struct {
	Obj   HircObj
	Root  uint32
	Leafs []ActorMixerHircNode
}

// Don't abstract with ActorMixerHircNode
type MusicHircNode struct {
	Obj   HircObj
	Root  uint32
	Leafs []MusicHircNode
}

func (h *HIRC) BuildTree() {
	h.ActorMixerRoots = []ActorMixerHircNode{}
	h.MusicHircRoots = []MusicHircNode{}
	var o HircObj
	l := len(h.HircObjs)
	for i := range h.HircObjs {
		o = h.HircObjs[l - i - 1]
		if ActorMixerHircType(o) {
			b := o.BaseParameter()
			if b == nil {
				panic(fmt.Sprintf("%d should have base parameter", o.HircType()))
			}
			if b.DirectParentId == 0 {
				node := ActorMixerHircNode{o, 0, make([]ActorMixerHircNode, o.NumLeaf())}
				h.WalkActorMixerHircRoot(&node)
				h.ActorMixerRoots = append(h.ActorMixerRoots, node)
			}
		} else if MusicHircType(o) {
			b := o.BaseParameter()
			if b == nil {
				panic(fmt.Sprintf("%d should have base parameter", o.HircType()))
			}
			if b.DirectParentId == 0 {
				node := MusicHircNode{o, 0, make([]MusicHircNode, o.NumLeaf())}
				h.WalkMusicHirciRoot(&node)
				h.MusicHircRoots = append(h.MusicHircRoots, node)
			}
		}
	}
}

func (h *HIRC) WalkActorMixerHircRoot(node *ActorMixerHircNode) {
	if !ActorMixerHircType(node.Obj) {
		panic("Panic Trap")
	}
	id, err := node.Obj.HircID()
	if err != nil {
		panic("Panic trap")
	}
	leafs := node.Obj.Leafs()
	l := len(leafs)
	for i := range leafs {
		v, ok := h.ActorMixerHirc.Load(leafs[l - i - 1])
		if !ok {
			panic("Panic Trap")
		}
		obj := v.(HircObj)
		node.Leafs[i] = ActorMixerHircNode{obj, id, make([]ActorMixerHircNode, obj.NumLeaf())}
		h.WalkActorMixerHircRoot(&node.Leafs[i])
	}
}

func (h *HIRC) WalkMusicHirciRoot(node *MusicHircNode) {
	if !MusicHircType(node.Obj) {
		panic("Panic Trap")
	}
	id, err := node.Obj.HircID()
	if err != nil {
		panic("Panic trap")
	}
	leafs := node.Obj.Leafs()
	l := len(leafs)
	for i := range leafs {
		v, ok := h.MusicHirc.Load(leafs[l - i - 1])
		if !ok {
			panic("Panic Trap")
		}
		obj := v.(HircObj)
		node.Leafs[i] = MusicHircNode{obj, id, make([]MusicHircNode, obj.NumLeaf())}
		h.WalkMusicHirciRoot(&node.Leafs[i])
	}
}
