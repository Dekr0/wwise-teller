package wwise

import (
	"fmt"
	"slices"
)

type HierarchyHeader struct {
	Type HircType
	Size u32
}

type HierarchyNode struct {
	Type   HircType
	Id     u32
}

type Container struct {
	Ids []u32
}

type Hierarchy struct {
	NextInternalId     u32 
	InternalIds      []u32
	Nodes            map[u32]*HierarchyNode
	DirectParentId   map[u32]u32
	Container        map[u32]*Container
	EncodedNodes     map[u32][]byte
}

// --- allocating --- //

func AllocHierarchy(numHirc u32) *Hierarchy {
	return &Hierarchy{
		NextInternalId: 0,
		InternalIds: make([]u32, 0, numHirc),
		DirectParentId: make(map[u32]u32, numHirc),
		Nodes: make(map[u32]*HierarchyNode, numHirc),
		Container: make(map[u32]*Container, numHirc),
		EncodedNodes: make(map[u32][]byte),
	}
}

func AllocContainer(size u32) *Container {
	return &Container{ make([]u32, size, size) }
}

// --- sizing --- //

func SizeOfContainer(c *Container) u32 {
	return Size32 + Size32 * u32(len(c.Ids))
}

// --- encoding --- //

func EncodeContainer(e *HircEncoderCtx, c *Container) error {
	size := SizeOfContainer(c)
	curr := e.Count()
	if err := e.Primitive(u32(len(c.Ids))); err != nil {
		return fmt.Errorf("Failed to encode # of ids in container: %w", err)
	}
	if err := e.Primitive(c.Ids); err != nil {
		return fmt.Errorf("Failed to encode ids in container: %w", err)
	}
	return e.Expect(curr, size)
}

// --- getter and setter --- //

func (h *Hierarchy) HasHierarchyNode(internalId u32) (in bool) {
	_, in = h.Nodes[internalId]
	return in
}

func (h *Hierarchy) GetHierarchyNode(internalId u32) (n *HierarchyNode) {
	n, in := h.Nodes[internalId]
	if !in {
		panic("Failed to locate hierarchy node")
	}
	return n
}

// Has side effect
func (h *Hierarchy) AddHierarchyNode(id u32, t HircType) (internalId u32) {
	internalId = h.NextInternalId
	if _, in := h.Nodes[id]; in {
		panic(MonotonicIdCollision)
	}
	// Get rid off this once I figure out the tree traversal algorithm 
	if slices.Contains(h.InternalIds, internalId) {
		panic(MonotonicIdCollision)
	}
	h.InternalIds = append(h.InternalIds, internalId)
	h.Nodes[internalId] = &HierarchyNode{ t, id, }
	h.NextInternalId++
	return internalId
}

func (h *Hierarchy) GetDirectParentId(internalId u32) u32 {
	if directParentId, in := h.DirectParentId[internalId]; !in {
		panic("Failed to locate direct parent id")
	} else {
		return directParentId
	}
}

func (h *Hierarchy) AddDirectParentId(internalId u32, directParentId u32) {
	if _, in := h.DirectParentId[internalId]; in {
		panic(MonotonicIdCollision)
	}
	h.DirectParentId[internalId] = directParentId
}

func (h *Hierarchy) GetContainer(internalId u32) *Container {
	if container, in := h.Container[internalId]; !in {
		panic("Failed to locate container")
	} else {
		return container
	}
}

func (h *Hierarchy) AddContainer(internalId u32, container *Container) {
	if container == nil {
		panic("Container is nil")
	}
	if _, in := h.Container[internalId]; in {
		panic(MonotonicIdCollision)
	}
	h.Container[internalId] = container
}

func (h *Hierarchy) GetEncodedHierarchyNode(internalId u32) []byte {
	if encoded, in := h.EncodedNodes[internalId]; !in {
		panic("Failed to locate encoded nodes")
	} else {
		return encoded
	}
}

// Has side effect
func (h *Hierarchy) AddEncodedHierarchyNode(id u32, t HircType, encoded []byte) {
	if encoded == nil {
		panic("Encoded hierarchy data is nil")
	}
	internalId := h.NextInternalId
	if _, in := h.Nodes[internalId]; in {
		panic(MonotonicIdCollision)
	}
	// Get rid off this once I figure out the tree traversal algorithm 
	if slices.Contains(h.InternalIds, internalId) {
		panic(MonotonicIdCollision)
	}
	h.InternalIds = append(h.InternalIds, internalId)
	h.Nodes[internalId] = &HierarchyNode{ t, id, }
	h.NextInternalId++
	if _, in := h.EncodedNodes[internalId]; in {
		panic(MonotonicIdCollision)
	}
	h.EncodedNodes[internalId] = encoded
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetHierarchyNode(internalId u32) *HierarchyNode {
	return h.Hierarchy.GetHierarchyNode(internalId)
}

func (h *HIRC) GetDirectParentId(internalId u32) u32 {
	return h.Hierarchy.GetDirectParentId(internalId)
}

func (h *HIRC) GetEncodedHierarchyNode(internalId u32) []byte {
	return h.Hierarchy.GetEncodedHierarchyNode(internalId)
}

func (h *HIRC) AddContainer(internalId u32, container *Container) {
	h.Hierarchy.AddContainer(internalId, container)
}

func (h *HIRC) GetContainer(internalId u32) *Container {
	return h.Hierarchy.GetContainer(internalId)
}

func (h *HIRC) EncodeEncodedHierarchy(e *HircEncoderCtx, t HircType, internalId u32) error {
	return EncodeEncodedHierarchy(
		e, t, h.GetEncodedHierarchyNode(internalId),
	)
}

// --- core procedure --- //

func EncodeEncodedHierarchy(e *HircEncoderCtx, t HircType, encoded []byte) (err error) {
	header := HierarchyHeader{ t, u32(len(encoded)) }
	if err = e.Struct(header, SizeOfHierarchyHeader); err != nil {
		return fmt.Errorf("Failed to encode hierarchy header: %w", err)
	}
	curr := e.Encoder.Count
	if err = e.Bytes(encoded); err != nil {
		return fmt.Errorf("Failed to write encoded chunk of: %w", err)
	}
	return e.Expect(curr, u32(len(encoded)))
}
