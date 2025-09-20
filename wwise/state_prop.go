package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type StateProp struct {
	NumStateProp   uio.V128
	PropId       []uio.V128
	AccumType    []AccumType
	InDb         []u8
}

type StatePropComponent struct {
	StateProp map[u32]StateProp
}

// --- allocating & freeing --- //

func AllocStateProp(numStateProp uio.V128) *StateProp {
	return &StateProp{
		NumStateProp: numStateProp,
		PropId: make([]uio.V128, numStateProp.V, numStateProp.V),
		AccumType: make([]AccumType, numStateProp.V, numStateProp.V),
		InDb: make([]u8, numStateProp.V, numStateProp.V),
	}
}

func AllocStatePropComponent(size u32) *StatePropComponent {
	if size <= 0 {
		return &StatePropComponent{
			StateProp: make(map[u32]StateProp),
		}
	}
	return &StatePropComponent{
		StateProp: make(map[u32]StateProp, size),
	}
}

// --- assertion --- //

func AssertStateProp(s StateProp) error {
	if s.NumStateProp.V != u64(len(s.PropId)) {
		return fmt.Errorf(
			"State property counter (%d) does not equal to # of state property ids (%d)",
			s.NumStateProp.V, len(s.PropId),
		)
	}
	if len(s.PropId) != len(s.AccumType) {
		return fmt.Errorf(
			"# of state property id (%d) does not equal to # of state accumluation type (%d)", 
			len(s.PropId), len(s.AccumType),
		)
	}
	if len(s.AccumType) != len(s.InDb) {
		return fmt.Errorf(
			"# of accumluation type (%d) does not equal to # of state decibel enabler (%d)", 
			len(s.AccumType), len(s.InDb),
		)
	}
	return nil
}

// --- sizing --- //

func SizeOfStateProp(s StateProp) (size u32) {
	size = u32(len(s.NumStateProp.B))
	size += Size8 * 2 * u32(len(s.PropId))
	for _, propId := range s.PropId {
		size += u32(len(propId.B))
	}
	return size
}

// --- encoding --- //

func EncodeStateProp(e *HircEncoderCtx, s StateProp) error {
	curr := e.Encoder.Count
	size := SizeOfStateProp(s)
	if err := e.Bytes(s.NumStateProp.B); err != nil {
		return fmt.Errorf("Failed to encode number of state property in LE128: %w", err)
	}
	for i, propId := range s.PropId {
		if err := e.Bytes(propId.B); err != nil {
			return fmt.Errorf("Failed to encode %d-th property id %d: %w", i, propId.V, err)
		}
		if err := e.Primitive(s.AccumType[i]); err != nil {
			return fmt.Errorf("Failed to encode %d-th accumulation type in state property: %w", i, err)
		}
		if err := e.Primitive(s.InDb[i]); err != nil {
			return fmt.Errorf("Failed to encode %d-th decibel enabled in state property: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

// --- component wrapper around assertion, encoding, and sizing --- //

func (c *StatePropComponent) AssertStatePropById(internalId u32) error {
	s, in := c.StateProp[internalId]
	if !in {
		return fmt.Errorf("Failed to locate state property")
	}
	return AssertStateProp(s)
}

func (c *StatePropComponent) SizeOfStatePropById(internalId u32) u32 {
	s, in := c.StateProp[internalId]
	if !in {
		panic("Failed to locate state property")
	}
	return SizeOfStateProp(s)
}

func (c *StatePropComponent) EncodeStatePropById(e *HircEncoderCtx, internalId u32) error {
	s, in := c.StateProp[internalId]
	if !in {
		return fmt.Errorf("Failed to locate state property")
	}
	return EncodeStateProp(e, s)
}

// --- component getter and setter --- //

func (c *StatePropComponent) HasStateProp(internalId u32) (in bool) {
	_, in = c.StateProp[internalId]
	return in
}

func (c *StatePropComponent) GetStateProp(internalId u32) (p StateProp) {
	p, in := c.StateProp[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate state prop."))
	}
	return p
}

func (c *StatePropComponent) AddStateProp(internalId u32, s StateProp) {
	if _, in := c.StateProp[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.StateProp[internalId] = s
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetStateProp(internalId u32) StateProp {
	return h.StatePropComponent.GetStateProp(internalId)
}
