package wwise

import (
	"fmt"
	"slices"
	"sort"
)

type ActorMixerHircNode struct {
	Obj   HircObj
	Root *ActorMixerHircNode
	Open  bool // Helper state for tree node rendering
	Leafs []*ActorMixerHircNode
}

// Don't abstract with ActorMixerHircNode
type MusicHircNode struct {
	Obj   HircObj
	Root  uint32
	Leafs []MusicHircNode
}

type BusHircNode struct {
	Obj      HircObj
	Root    *BusHircNode // OverrideBusId
	Open     bool
	LeafsIdx []uint32
	Leafs    map[uint32]*BusHircNode
}

// TODO Optimize tree rebuilding
func (h *HIRC) BuildTree() {
	// This should be able to allocated in a deterministic manner because the 
	// exact # can be accumulative during the parsing phase and post modification phase
	h.ActorMixerRoots = []*ActorMixerHircNode{}
	h.ActorMixerHircNodesMap = map[uint32]*ActorMixerHircNode{}
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
				node := &ActorMixerHircNode{o, nil, false, make([]*ActorMixerHircNode, o.NumLeaf())}
				id, err := o.HircID()
				if err != nil {
					panic(err)
				}
				if _, in := h.ActorMixerHircNodesMap[id]; in {
					panic(fmt.Sprintf("Duplicate actor mixer hierarchy node %d", id))
				}
				h.ActorMixerHircNodesMap[id] = node
				h.WalkActorMixerHircRoot(node)
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
	h.BusHircNodesMap = make(map[uint32]*BusHircNode)
	// This should be able to allocated in a deterministic manner because the 
	// exact # can be accumulative during the parsing phase and post modification phase
	h.BusRoots = []*BusHircNode{}
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
				bus, nil, false, []uint32{}, make(map[uint32]*BusHircNode),
			}
			if _, in := h.BusHircNodesMap[bus.Id]; in {
				panic(fmt.Sprintf("Duplicate bus hierarchy node %d", bus.Id))
			}
			h.BusHircNodesMap[bus.Id] = node
			if bus.OverrideBusId == 0 {
				h.BusRoots = append(h.BusRoots, node)
			}
		case *AuxBus:
			node := &BusHircNode{
				bus, nil, false, []uint32{}, make(map[uint32]*BusHircNode),
			}
			if _, in := h.BusHircNodesMap[bus.Id]; in {
				panic(fmt.Sprintf("Duplicate aux bus hierarchy node %d", bus.Id))
			}
			h.BusHircNodesMap[bus.Id] = node
			if bus.OverrideBusId == 0 {
				h.BusRoots = append(h.BusRoots, node)
			}
		default:
			panic("Non bus hierarchy object is being processed during the bus hierarchy tree construction.")
		}
	}
	// Second pass, construct tree
	for i, o := range h.HircObjs {
		if !BusHircType(o) {
			continue
		}
		overrideBusId := uint32(0)
		switch bus := o.(type) {
		case *Bus:
			overrideBusId = bus.OverrideBusId
		case *AuxBus:
			overrideBusId = bus.OverrideBusId
		default:
			panic("Non bus hierarchy object is being processed during the bus hierarchy tree construction.")
		}
		if overrideBusId == 0 {
			continue
		}

		id, err := o.HircID()
		if err != nil {
			panic(err)
		}
		node, in := h.BusHircNodesMap[id];
		if !in {
			panic(fmt.Sprintf("No bus hierarchy node has ID %d.", id))
		}

		root, in := h.BusHircNodesMap[overrideBusId]; 
		if !in {
			panic(fmt.Sprintf("No root bus hierarchy node has ID %d.", overrideBusId))
		} 
		node.Root = root

		if _, in := root.Leafs[uint32(i)]; in {
			panic(fmt.Sprintf("Duplicate bus hierarchy object index %d", i))
		}
		root.Leafs[uint32(i)] = node

		insertIdx, in := sort.Find(len(root.LeafsIdx), func(j int) int {
			if root.LeafsIdx[j] == uint32(i) {
				panic(fmt.Sprintf("Duplicate bus hierarchy object index %d", i))
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

func (h *HIRC) OpenActorMixerHircNode(id uint32) {
	node, in := h.ActorMixerHircNodesMap[id]
	if !in {
		panic(fmt.Sprintf("No actor mixer hierarchy node has ID %d", id))
	}
	node.Open = true
	parent := node.Root
	for parent != nil {
		parent.Open = true
		parent = parent.Root
	}
}

func (h *HIRC) OpenBusHircNode(id uint32) {
	node, in := h.BusHircNodesMap[id]
	if !in {
		panic(fmt.Sprintf("No bus hierarchy node has ID %d", id))
	}
	node.Open = true
	root := node.Root
	for root != nil {
		root.Open = true
		root = root.Root
	}
}

// Resolve which bus can use enable HDR
func (h *HIRC) HDRAvailability() {
	for i := range h.HircObjs {
		o := h.HircObjs[i]
		switch bus := o.(type) {
		case *Bus:
			bus.CanSetHDR = h.CanSetHDR(bus, bus.OverrideBusId)
		default:
		}
	}
}

func (h *HIRC) CanSetHDR(b *Bus, parentID uint32) int8 {
	if parentID == 0 {
		b.CanSetHDR = 1
		return b.CanSetHDR
	}
	v, in := h.Buses.Load(parentID)
	if !in {
		panic("Panic Trap")
	}
	parentBus := v.(*Bus)
	if parentBus.CanSetHDR != -1 {
		b.CanSetHDR = parentBus.CanSetHDR
		return b.CanSetHDR
	}
	b.CanSetHDR = h.CanSetHDR(parentBus, parentBus.OverrideBusId)
	return b.CanSetHDR
}

func (h *HIRC) WalkActorMixerHircRoot(node *ActorMixerHircNode) {
	if !ActorMixerHircType(node.Obj) {
		panic("Panic Trap")
	}
	leafs := node.Obj.Leafs()
	l := len(leafs)
	for i := range leafs {
		id := leafs[l - i - 1]
		v, ok := h.ActorMixerHirc.Load(id)
		if !ok {
			panic("Panic Trap")
		}
		obj := v.(HircObj)
		node.Leafs[i] = &ActorMixerHircNode{
			obj, node, false, make([]*ActorMixerHircNode, obj.NumLeaf()),
		}
		if _, in := h.ActorMixerHircNodesMap[id]; in {
			panic(fmt.Sprintf("Duplicate actor mixer hierarchy object %d", id))
		}
		h.ActorMixerHircNodesMap[id] = node.Leafs[i]
		h.WalkActorMixerHircRoot(node.Leafs[i])
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
