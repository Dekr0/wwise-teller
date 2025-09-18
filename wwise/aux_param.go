package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type AuxParam struct {
	// bit 2 Override User Aux Sends
	// bit 3 Has Aux
	// bit 4 Override Reflection Aux Bus
	SettingVector    u8
	AuxIds        [4]u32
	ReflectionAux    u32
}

type AuxParamComponent struct {
	AuxParam map[u32]*AuxParam
}

// --- allocating --- //

func AllocAuxParam() (a *AuxParam) {
	return &AuxParam{
		SettingVector: 0,
		AuxIds: [4]u32{ 0, 0, 0, 0 },
		ReflectionAux: 0,
	}
}

func AllocAuxParamComponent(size u32) *AuxParamComponent {
	if size <= 0 {
		return &AuxParamComponent{
			make(map[u32]*AuxParam),
		}
	}
	return &AuxParamComponent{make(map[u32]*AuxParam, size)}
}

// --- assertion --- //

func AssertAuxParam(a *AuxParam) error {
	if HasAux(a) {
		numVoidAux := 0
		for _, a := range a.AuxIds {
			if a == 0 {
				numVoidAux++
			}
		}
		if numVoidAux >= len(a.AuxIds) {
			return fmt.Errorf("User-defined auxiliary send is enabled but all auxiliary sends are 0")
		}
	} else {
		for i, a := range a.AuxIds {
			if a != 0 {
				return fmt.Errorf("User-defined auxiliary send is not enabled but %d-th auxiliary send is none zero", i)
			}
		}
	}
	if !OverrideReflectionAux(a) && a.ReflectionAux != 0 {
		return fmt.Errorf("Override reflection auxiliary send is not enabled but reflection aux send is none zero")
	}
	return nil
}

// --- size --- //

func SizeOfAuxParam(a *AuxParam) (size u32) {
	size = Size8 + Size32
	if HasAux(a) {
		numVoidAux := 0
		for _, a := range a.AuxIds {
			if a == 0 {
				numVoidAux++
			}
		}
		if numVoidAux < len(a.AuxIds) {
			size += 4 * Size32
		}
	}
	return size
}

// --- encoding --- //

func EncodeAuxParam(e *HircEncoderCtx, a *AuxParam) error {
	curr := e.Encoder.Count
	size := SizeOfAuxParam(a)
	if err := e.Primitive(a.SettingVector); err != nil {
		return fmt.Errorf("Failed to encode auxiliary parameter setting vector: %w", err)
	}
	if HasAux(a) {
		for i, a := range a.AuxIds {
			if err := e.Primitive(a); err != nil {
				return fmt.Errorf("Failed to encode %d-th user-defined auxiliary (%d): %w", i, a, err)
			}
		}
	}
	if err := e.Primitive(a.ReflectionAux); err != nil {
		return fmt.Errorf("Failed to encode relfection auxiliary send %d: %w",a.ReflectionAux, err)
	}
	return e.Expect(curr, size)
}

// --- Component assertion wrapper --- //

func (c *AuxParamComponent) AssertAuxParamById(internalId u32) error {
	a, in := c.AuxParam[internalId]
	if !in {
		return fmt.Errorf("Failed to locate auxiliary parameter")
	}
	return AssertAuxParam(a)
}

// --- component sizing wrapper --- //

func (c *AuxParamComponent) SizeOfAuxParamById(internalId u32) (size u32) {
	a, in := c.AuxParam[internalId]
	if !in {
		panic("Failed to locate auxiliary parameter")
	}
	return SizeOfAuxParam(a)
}

// --- component encoding wrapper --- //

func (c *AuxParamComponent) EncodeAuxParamById(e *HircEncoderCtx, internalId u32) error {
	a, in := c.AuxParam[internalId]
	if !in {
		return fmt.Errorf("Failed to locate auxiliary parameter")
	}
	return EncodeAuxParam(e, a)
}

// --- component getter and setter --- //

func (c *AuxParamComponent) HasAuxParam(internalId u32) (in bool) {
	_, in = c.AuxParam[internalId]
	return in
}

func (c *AuxParamComponent) GetAuxParam(internalId u32) (a *AuxParam) {
	a, in := c.AuxParam[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate auxiliary parameter."))
	}
	return a
}

func (c *AuxParamComponent) AddAuxParam(internalId u32, a *AuxParam) {
	if a == nil {
		panic("Auxiliary parameter is nil")
	}
	if _, in := c.AuxParam[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.AuxParam[internalId] = a
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetAuxParam(internalId u32) (a *AuxParam) {
	return h.AuxParamComponent.GetAuxParam(internalId)
}

// --- core procedure --- //

func HasAux(a *AuxParam) bool {
	return uio.Bit(a.SettingVector, 3)
}

func OverrideReflectionAux(a *AuxParam) bool {
	return uio.Bit(a.SettingVector, 4)
}

func SetAux(a *AuxParam, set bool) {
	a.SettingVector = uio.SetBit(a.SettingVector, 3, set)
}
