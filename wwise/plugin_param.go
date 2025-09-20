package wwise

import (
	"fmt"
)

// --- struct definition --- //

type PluginParam struct {
	Size   u32

	// No interface require like the old design. The reason behind this is that 
	// concrete data structure is only required when users need to view them 
	// or modify them. Transformation is cheap.
	Data []byte
}

type PluginParamComponent struct {
	PluginParam  map[u32]PluginParam
}

// --- allocation / freeing --- //

func AllocPluginParamComponent(size u32) (c *PluginParamComponent) {
	c = &PluginParamComponent{}
	if size <= 0 {
		c.PluginParam = make(map[u32]PluginParam)
	} else {
		c.PluginParam = make(map[u32]PluginParam, size)
	}
	return c
}

func AllocPluginParam(size u32) *PluginParam {
	return &PluginParam{
		Size: size,
		Data: make([]byte, size, size),
	}
}

// --- assertion --- //

// No side effect
func AssertPluginParm(p PluginParam) error {
	if p.Size != u32(len(p.Data)) {
		return fmt.Errorf("Plugin parameter size counter (%d) does not equal to size of actual data size (%d)", p.Size, len(p.Data))
	}
	return nil
}

// --- sizing --- //

// No side effect
func SizeOfPluginParam(p PluginParam) u32 {
	return Size32 + u32(len(p.Data))
}

// --- encoding --- //

func EncodePluginParam(e *HircEncoderCtx, p PluginParam) (err error) {
	size := SizeOfPluginParam(p)

	curr := e.Encoder.Count

	if err := e.Primitive(p.Size); err != nil {
		return fmt.Errorf("Failed to encode plugin parameter size: %w", err)
	}

	if err := e.Bytes(p.Data); err != nil {
		return fmt.Errorf("Failed to encode plugin parameter data bytes: %w", err)
	}

	return e.Expect(curr, size)
}

// --- component assertion wrapper --- //

// No side effect
func (c *PluginParamComponent) AssertPluginParamById(internalId u32) error {
	p, in := c.PluginParam[internalId]
	if !in {
		return fmt.Errorf("Failed to locate plugin parameter")
	}
	return AssertPluginParm(p)
}

// --- component sizing wrapper --- //

// No side effect
func (c *PluginParamComponent) SizeOfPluginParamById(internalId u32) u32 {
	p, in := c.PluginParam[internalId]
	if !in {
		panic("Failed to locate plugin parameter")
	}
	return SizeOfPluginParam(p)
}

// --- component getter and setter --- //

// No side effect
func (c *PluginParamComponent) HasPluginParam(internalId u32) (in bool) {
	_, in = c.PluginParam[internalId]
	return in
}

// No side effect
func (c *PluginParamComponent) GetPluginParam(internalId u32) (p PluginParam) {
	p, in := c.PluginParam[internalId]
	if !in {
		panic("Failed to locate plugin param")
	}
	return p
}

// Has side effect
func (c *PluginParamComponent) AddPluginParam(internalId u32, p PluginParam) {
	if _, in := c.PluginParam[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.PluginParam[internalId] = p
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetPluginParam(internalId u32) PluginParam {
	return h.PluginParamComponent.GetPluginParam(internalId)
}
