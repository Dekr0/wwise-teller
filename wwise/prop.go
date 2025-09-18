package wwise

import (
	"fmt"
)

// --- struct definition --- //

type Prop struct {
	Ids  []u8
	Vals []f32
}

type PropS struct {
	Id  u8
	Val u32
}

type PropComponent struct {
	Prop map[u32]*Prop
}

type RProp struct {
	Ids  []u8
	Mins []f32
	Maxs []f32
}

type RPropS struct {
	Id    u8
	Range RPropRange
}

type RPropRange struct {
	Min f32
	Max f32
}

type RPropComponent struct {
	RProp map[u32]*RProp
}

// --- allocation / freeing --- //

func AllocProp(numProp u8) *Prop {
	return &Prop{
		Ids:  make([]u8, numProp, numProp),
		Vals: make([]f32, numProp, numProp),
	}
}

func AllocRProp(numRProp u8) *RProp {
	return &RProp{
		Ids:  make([]u8, numRProp, numRProp),
		Mins: make([]f32, numRProp, numRProp),
		Maxs: make([]f32, numRProp, numRProp),
	}
}

func AllocPropComponent(size u32) *PropComponent {
	if size <= 0 {
		return &PropComponent{make(map[u32]*Prop)}
	}
	return &PropComponent{make(map[u32]*Prop, size)}
}

func AllocRPropComponent(size u32) *RPropComponent {
	if size <= 0 {
		return &RPropComponent{make(map[u32]*RProp)}
	}
	return &RPropComponent{make(map[u32]*RProp, size)}
}

// --- assertion --- //

func AssertProp(p *Prop) error {
	if len(p.Ids) != len(p.Vals) {
		return fmt.Errorf("# of property ids (%d) does not equal to # of property value (%d)", len(p.Ids), len(p.Vals))
	}
	return nil
}

func AssertRProp(p *RProp) error {
	if len(p.Ids) != len(p.Mins) {
		return fmt.Errorf("# of property ids (%d) does not equal to # of min property value (%d)", len(p.Ids), len(p.Mins))
	}
	if len(p.Mins) != len(p.Maxs) {
		return fmt.Errorf("# of property ids (%d) does not equal to # of max property value (%d)", len(p.Ids), len(p.Maxs))
	}
	return nil
}

// --- sizing --- //

func SizeOfProp(p *Prop) u32 {
	return SizeOfPropCounter + (SizeOfPropId + SizeOfPropValue) * u32(len(p.Ids))
}

func SizeOfRProp(p *RProp) u32 {
	return SizeOfPropCounter + (SizeOfPropId + SizeOfRPropValue) * u32(len(p.Ids))
}

// --- encoding --- //

// No side effect 
func EncodeProp(e *HircEncoderCtx, p *Prop) error {
	size := SizeOfProp(p)
	curr := e.Encoder.Count
	if err := e.Primitive(u8(len(p.Ids))); err != nil {
		return fmt.Errorf("Failed to encode property size counter: %w", err)
	}
	if err := e.Primitive(p.Ids); err != nil {
		return fmt.Errorf("Failed to encode property ids: %w", err)
	}
	if err := e.Primitive(p.Vals); err != nil {
		return fmt.Errorf("Failed to encode property values: %w", err)
	}
	return e.Expect(curr, size)
}

// No side effect
func EncodeRProp(e *HircEncoderCtx, p *RProp) error {
	size := SizeOfRProp(p)
	curr := e.Encoder.Count
	if err := e.Primitive(u8(len(p.Ids))); err != nil {
		return fmt.Errorf("Failed to encode range-based property size counter: %w", err)
	}
	if err := e.Primitive(p.Ids); err != nil {
		return fmt.Errorf("Failed to encode range-based property ids: %w", err)
	}
	for i, id := range p.Ids {
		r := RPropRange{p.Mins[i], p.Maxs[i]}
		if err := e.Struct(r, SizeOfRPropValue); err != nil {
			return fmt.Errorf("Failed to encode range values of property id %d: %w", id, err)
		}
	}
	return e.Expect(curr, size)
}


// --- component assertion wrapper --- //

// No side effect
func (c *PropComponent) AssertPropById(internalId u32) error {
	if p, in := c.Prop[internalId]; !in {
		panic("Failed to locate Prop")
	} else {
		return AssertProp(p)
	}
}

// No side effect
func (c *RPropComponent) AssertRPropById(internalId u32) error {
	if r, in := c.RProp[internalId]; !in {
		panic("Failed to locate RProp")
	} else {
		return AssertRProp(r)
	}
}

// --- component sizing wrapper --- //

// No side effect
func (c *PropComponent) SizeOfPropById(internalId u32) u32 {
	if p, in := c.Prop[internalId]; !in {
		panic("Failed to locate Prop")
	} else {
		return SizeOfProp(p)
	}
}

// No side effect
func (c *RPropComponent) SizeOfRPropById(internalId u32) u32 {
	if p, in := c.RProp[internalId]; !in {
		panic("Failed to locate RProp")
	} else {
		return SizeOfRProp(p)
	}
}

// --- component getter and setter --- //

// Has no side effect
func (c *PropComponent) HasProp(internalId u32) (in bool) {
	_, in = c.Prop[internalId]
	return in
}

// Has no side effect
func (c *RPropComponent) HasRProp(internalId u32) (in bool) {
	_, in = c.RProp[internalId]
	return in
}

// Has no side effect
func (c *PropComponent) GetProp(internalId u32) (p *Prop) {
	p, in := c.Prop[internalId]
	if !in {
		panic("Failed to locate Prop")
	}
	return p
}

// Has no side effect
func (c *RPropComponent) GetRProp(internalId u32) (r *RProp) {
	r, in := c.RProp[internalId]
	if !in {
		panic("Failed to locate RProp")
	}
	return r
}

// Has side effect
func (c *PropComponent) AddProp(internalId u32, p *Prop) {
	if p == nil {
		panic("Propery is nil")
	}
	if _, in := c.Prop[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.Prop[internalId] = p
}

// Has side effect
func (c *RPropComponent) AddRProp(internalId u32, r *RProp) {
	if r == nil {
		panic("Range-based property is nil")
	}
	if _, in := c.RProp[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.RProp[internalId] = r

}

// --- HIRC component wrapper --- //

func (h *HIRC) GetProp(internalId u32) (p *Prop) {
	return h.PropComponent.GetProp(internalId)
}

func (h *HIRC) GetRProp(internalId u32) (p *RProp) {
	return h.RPropComponent.GetRProp(internalId)
}

// --- constant and type alias --- //

const SizeOfPropCounter = Size8
const SizeOfPropId      = Size8
const SizeOfPropValue   = Size32
const SizeOfRPropValue  = Size32 * 2
