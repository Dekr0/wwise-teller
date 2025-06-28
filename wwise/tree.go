package wwise

import (
	"fmt"
	"slices"
	"sort"
)

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

type BusHircNode struct {
	Obj      HircObj
	Root     uint32 // OverrideBusId
	LeafsIdx []uint32
	Leafs    map[uint32]*BusHircNode
}

func (h *HIRC) BuildTree() {
	// This should be able to allocated in a deterministic manner because the 
	// exact # can be accumulative during the parsing phase and post modification phase
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

	// Memory allocation most likely scatter all over the place
	busNodes := make(map[uint32]*BusHircNode)
	// This should be able to allocated in a deterministic manner because the 
	// exact # can be accumulative during the parsing phase and post modification phase
	h.BusRoots = []BusHircNode{}
	// First pass, obtain all nodes. Since bus hierarchy object does not keep 
	// track of leaf buses with a list of IDs, there's chances that a leaf node 
	// is allocated first but parent node is yet allocated.
	for _, o := range h.HircObjs {
		if !BusHircType(o) {
			continue
		}
		switch bus := o.(type) {
		case *Bus:
			node := &BusHircNode{
				bus, bus.OverrideBusId, []uint32{}, make(map[uint32]*BusHircNode),
			}
			if _, in := busNodes[bus.Id]; in {
				panic("Panic Trap")
			}
			busNodes[bus.Id] = node
			if bus.OverrideBusId == 0 {
				h.BusRoots = append(h.BusRoots, *node)
			}
		case *AuxBus:
			node := &BusHircNode{
				bus, bus.OverrideBusId, []uint32{}, make(map[uint32]*BusHircNode),
			}
			if _, in := busNodes[bus.Id]; in {
				panic("Panic Trap")
			}
			busNodes[bus.Id] = node
			if bus.OverrideBusId == 0 {
				h.BusRoots = append(h.BusRoots, *node)
			}
		}
	}
	// Second pass, construct tree
	for i, o := range h.HircObjs {
		if !BusHircType(o) {
			continue
		}

		id, err := o.HircID()
		if err != nil {
			panic("Panic Trap")
		}

		node, in := busNodes[id];
		if !in {
			panic("Panic Trap")
		}
		if node.Root == 0 {
			continue
		}

		root, in := busNodes[node.Root]; 
		if !in {
			panic("Panic Trap")
		} 

		if _, in := root.Leafs[uint32(i)]; in {
			panic("Panic Trap")
		}
		root.Leafs[uint32(i)] = node

		insertIdx, in := sort.Find(len(root.LeafsIdx), func(j int) int {
			if root.LeafsIdx[j] == uint32(i) {
				panic("Panic Trap")
			}
			if root.LeafsIdx[j] > uint32(i) {
				return 1
			}
			return -1
		})
		if !in {
			root.LeafsIdx = slices.Insert(root.LeafsIdx, insertIdx, uint32(i))
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
