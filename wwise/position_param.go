package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type PositionParam struct {
	// bit 0 override parent
	//  - when a hierarchy is a root level hierarchy (direct parent id is zero), 
	//  bit 0 is automatically set to 1
	// bit 1 listener relative routing 
	//  - when a hierarchy is a root level hierarchy (direct parent id is zero), 
	//  bit 1 is automatically set to 1
	// bit 2 to bit 4 - panner type
	// bit 5 to bit 8 - 3D Position Type
	SettingVector        u8
	// bit 0 to bit 2 - Spatialization mode
	// bit 3 - Enable Attenuation
	// bit 4 - Hold Emitter Position And Orientation
	// bit 5 - Hold Listener Orientation
	// bit 6 - Enable diffraction
	SpatialSettingVector u8
	PathMode             PathMode
	TransitionTime       i32
	PositionVertices     PositionVertices       
	PositionPlayList     PositionPlayList
	SpatialAutomation    SpatialAutomation
}

type PositionVertices struct {
	PositionVerticesX        []f32
	PositionVerticesY        []f32
	PositionVerticesZ        []f32
	PositionVerticesDuration []i32
}

type PositionVertexS struct {
	X        f32
	Y        f32
	Z        f32
	Duration i32
}

type PositionPlayList struct {
	VerticesOffset []u32
	NumVertices    []u32
}

type PositionPlayListS struct {
	VerticesOffset u32
	NumVertices    u32
}

type SpatialAutomation struct {
	RangeX []f32
	RangeY []f32
	RangeZ []f32
}

type SpatialAutomationS struct {
	RangeX f32
	RangeY f32
	RangeZ f32
}

type PositionParamComponent struct {
	PositionParam map[u32]*PositionParam
}

// --- allocation and freeing --- //

func AllocPositionParam() *PositionParam {
	return &PositionParam{
		PositionVertices: PositionVertices{
			make([]f32, 0, 4),
			make([]f32, 0, 4),
			make([]f32, 0, 4),
			make([]i32, 0, 4),
		},
		PositionPlayList: PositionPlayList{
			make([]u32, 0, 4),
			make([]u32, 0, 4),
		},
		SpatialAutomation: SpatialAutomation{
			make([]f32, 0, 4),
			make([]f32, 0, 4),
			make([]f32, 0, 4),
		},
	}
}

func AllocPositionParamComponent(size u32) *PositionParamComponent {
	if size <= 0 {
		return &PositionParamComponent{make(map[u32]*PositionParam)}
	}
	return &PositionParamComponent{make(map[u32]*PositionParam, size)}
}

// --- assertion --- //

func AssertPositionParam(p *PositionParam, isRoot bool) error {
	if isRoot {
		if !PositionOverrideParent(p) {
			return fmt.Errorf("Position override parent is not set for a root level hierarchy")
		}
		if !PositionListenerRelativeRouting(p) {
			return fmt.Errorf("Position listener relative routing is not set for a root level hierarchy")
		}
	}

	if !PositionOverrideParent(p) {
		if PositionListenerRelativeRouting(p) {
			return fmt.Errorf("Position override parent is not set but listener relative routing is enabled.")
		}
		if ((p.SettingVector >> 5) & 3) != 0 {
			return fmt.Errorf("Position override parent is not set but 3D automation is enabled")
		}
		if p.SpatialSettingVector != 0 {
			return fmt.Errorf("Position override parent is not set but spatial setting vector is non zero")
		}
		if p.PathMode != 0 {
			return fmt.Errorf("Position override parent is not set but path mode value is non zero")
		}
		if p.TransitionTime != 0 {
			return fmt.Errorf("Position override parent is not set but transition time value is non zero")
		}
		if len(p.PositionVertices.PositionVerticesX) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesY) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesZ) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesDuration) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position vertices.")
		}
		if len(p.PositionPlayList.VerticesOffset) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position playlist.")
		}
		if len(p.PositionPlayList.NumVertices) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in position playlist.")
		}
		if len(p.SpatialAutomation.RangeX) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeY) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeZ) > 0 {
			return fmt.Errorf("Position override parent is not set but there are elements in spatial 3D automation.")
		}
	}
	if !PositionListenerRelativeRouting(p) {
		if ((p.SettingVector >> 5) & 3) != 0 {
			return fmt.Errorf("Position listen relative routing is not set but 3D automation is enabled.")
		}
		if p.SpatialSettingVector != 0 {
			return fmt.Errorf("Position listen relative routing is not set but spatial setting vector is non zero")
		}
		if p.PathMode != 0 {
			return fmt.Errorf("Position listen relative routing is not set but path mode value is non zero")
		}
		if p.TransitionTime != 0 {
			return fmt.Errorf("Position listen relative routing is not set but transition time value is non zero")
		}
		if len(p.PositionVertices.PositionVerticesX) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesY) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesZ) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesDuration) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position vertices.")
		}
		if len(p.PositionPlayList.VerticesOffset) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position playlist.")
		}
		if len(p.PositionPlayList.NumVertices) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in position playlist.")
		}
		if len(p.SpatialAutomation.RangeX) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeY) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeZ) > 0 {
			return fmt.Errorf("Position listen relative routing is not set but there are elements in spatial 3D automation.")
		}
	}
	if ((p.SettingVector >> 5) & 3) == 0 {
		if p.PathMode != 0 {
			return fmt.Errorf("3D Emitter is set but path mode value is non zero")
		}
		if p.TransitionTime != 0 {
			return fmt.Errorf("3D Emitter is set but transition time value is non zero")
		}
		if len(p.PositionVertices.PositionVerticesX) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesY) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesZ) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position vertices.")
		}
		if len(p.PositionVertices.PositionVerticesDuration) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position vertices.")
		}
		if len(p.PositionPlayList.VerticesOffset) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position playlist.")
		}
		if len(p.PositionPlayList.NumVertices) > 0 {
			return fmt.Errorf("3D Emitter is set but there are elements in position playlist.")
		}
		if len(p.SpatialAutomation.RangeX) > 0 {
			return fmt.Errorf("3D Emitter is set set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeY) > 0 {
			return fmt.Errorf("3D Emitter is set set but there are elements in spatial 3D automation.")
		}
		if len(p.SpatialAutomation.RangeZ) > 0 {
			return fmt.Errorf("3D Emitter is set set but there are elements in spatial 3D automation.")
		}
	}
	if len(p.PositionVertices.PositionVerticesX) != len(p.PositionVertices.PositionVerticesY) {
		return fmt.Errorf("# of position X vertices does not equal to # of position Y vertices")
	}
	if len(p.PositionVertices.PositionVerticesY) != len(p.PositionVertices.PositionVerticesZ) {
		return fmt.Errorf("# of position Y vertices does not equal to # of position Z vertices")
	}
	if len(p.PositionVertices.PositionVerticesZ) != len(p.PositionVertices.PositionVerticesDuration) {
		return fmt.Errorf("# of position Z vertices does not equal to # of position play duration for each vertice")
	}
	if len(p.PositionPlayList.VerticesOffset) != len(p.PositionPlayList.NumVertices) {
		return fmt.Errorf("# of vertices offset in position playlist does not equal to # of number vertices in position playlist")
	}
	if len(p.SpatialAutomation.RangeX) != len(p.SpatialAutomation.RangeY) {
		return fmt.Errorf("# of range X values does not equal to # of range Y values in spatial automation")
	}
	if len(p.SpatialAutomation.RangeY) != len(p.SpatialAutomation.RangeZ) {
		return fmt.Errorf("# of range Y values does not equal to # of range Z values in spatial automation")
	}
	return nil
}

// --- sizing --- //

func SizeOfPositionParam(p *PositionParam) (size u32) {
	size = Size8
	if !PositionOverrideParentAndListenerRelativeRounting(p) {
		return size
	}
	size += Size8
	if !PositionHas3D(p) {
		return size
	}
	size += Size8 + Size32
	size += Size32 + u32(len(p.PositionVertices.PositionVerticesX)) * SizeOfPositionVertice
	size += Size32 + u32(len(p.PositionPlayList.VerticesOffset)) * SizeOfPositionPlayList
	return size + u32(len(p.SpatialAutomation.RangeX)) * SizeOfSpatialAutomation
}

// --- encoding --- //

func EncodePositionParam(e *HircEncoderCtx, p *PositionParam) error {
	curr := e.Encoder.Count
	size := SizeOfPositionParam(p)
	if err := e.Primitive(p.SettingVector); err != nil {
		return fmt.Errorf("Failed to encode position parameter setting vector: %w", err)
	}
	if !PositionOverrideParentAndListenerRelativeRounting(p) {
		return e.Expect(curr, size)
	}
	if err := e.Primitive(p.SpatialSettingVector); err != nil {
		return fmt.Errorf("Failed to encode position parameter spatial setting vector: %w", err)
	}
	if !PositionHas3D(p) {
		return e.Expect(curr, size)
	}
	if err := e.Primitive(p.PathMode); err != nil {
		return fmt.Errorf("Failed to encode position parameter path mode: %w", err)
	}
	if err := e.Primitive(p.TransitionTime); err != nil {
		return fmt.Errorf("Failed to encode transition time %d: %w", p.TransitionTime, err)
	}
	if err := e.Primitive(u32(len(p.PositionVertices.PositionVerticesX))); err != nil {
		return fmt.Errorf("Failed to encode # of position vertices in position parameter: %w", err)
	}
	for i, x := range p.PositionVertices.PositionVerticesX {
		payload := PositionVertexS{
			x, 
			p.PositionVertices.PositionVerticesY[i],
			p.PositionVertices.PositionVerticesZ[i],
			p.PositionVertices.PositionVerticesDuration[i],
		}
		if err := e.Struct(payload, SizeOfPositionVertice); err != nil {
			return fmt.Errorf("Failed to encode %d-th of position vertice in position parameter: %w", i, err)
		}
	}
	if err := e.Primitive(u32(len(p.PositionPlayList.VerticesOffset))); err != nil {
		return fmt.Errorf("Failed to encode # of position playlist items in position parameter: %w", err)
	}
	for i, offset := range p.PositionPlayList.VerticesOffset {
		payload := PositionPlayListS{
			offset, p.PositionPlayList.NumVertices[i],
		}
		if err := e.Struct(payload, SizeOfPositionPlayList); err != nil {
			return fmt.Errorf("Failed to encode %d-th of position playlist item in position parameter: %w", i, err)
		}
	}
	for i, x := range p.SpatialAutomation.RangeX {
		payload := SpatialAutomationS{
			RangeX: x,
			RangeY: p.SpatialAutomation.RangeY[i],
			RangeZ: p.SpatialAutomation.RangeZ[i],
		}
		if err := e.Struct(payload, SizeOfSpatialAutomation); err != nil {
			return fmt.Errorf("Failed to encode %d-th of spatial item in position parameter: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

// --- component assertion wrapper --- //

func (c *PositionParamComponent) AssertPositionParamById(internalId u32, isRoot bool) error {
	p, in := c.PositionParam[internalId]
	if !in {
		return fmt.Errorf("Failed to locate position parameter.")
	}
	return AssertPositionParam(p, isRoot)
}

// --- component sizing wrapper --- //

func (c *PositionParamComponent) SizeOfPositionParamById(internalId u32) u32 {
	p, in := c.PositionParam[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate position parameter"))
	}
	return SizeOfPositionParam(p)
}

// --- component encoding wrapper --- //

func (c *PositionParamComponent) EncodePositionParamById(e *HircEncoderCtx, internalId u32) error {
	p, in := c.PositionParam[internalId]
	if !in {
		return fmt.Errorf("Failed to locate position parameter")
	}
	return EncodePositionParam(e, p)
}

// --- component getter and setter --- //

// Has no side effect
func (c *PositionParamComponent) HasPositionParam(internalId u32) (in bool) {
	_, in = c.PositionParam[internalId]
	return in
}

func (c *PositionParamComponent) GetPositionParam(internalId u32) (p *PositionParam) {
	p, in := c.PositionParam[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate position parameter."))
	}
	return p
}

// Has side effect
func (c *PositionParamComponent) AddPositionParam(internalId u32, p *PositionParam) {
	if _, in := c.PositionParam[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.PositionParam[internalId] = p
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetPositionParam(internalId u32) *PositionParam {
	return h.PositionParamComponent.GetPositionParam(internalId)
}

// --- core procedure --- //

// Has no side effect
func PositionOverrideParent(p *PositionParam) bool {
	return uio.Bit(p.SettingVector, 0)
}

// Has no side effect
func PositionListenerRelativeRouting(p *PositionParam) bool {
	if !PositionOverrideParent(p) {
		return false
	}
	return uio.Bit(p.SettingVector, 1)
}

func PositionOverrideParentAndListenerRelativeRounting(p *PositionParam) bool {
	return PositionOverrideParent(p) && PositionListenerRelativeRouting(p)
}

func Position3DPositionType(p *PositionParam) u8 {
	if !PositionOverrideParentAndListenerRelativeRounting(p) {
		panic("Attempt to obtain 3D position type when position override parent and listener relative routing is not enable.")
	}
	return (p.SettingVector >> 5) & 3
}

// Has no side effect
func PositionHas3D(p *PositionParam) bool {
	if !PositionOverrideParentAndListenerRelativeRounting(p) {
		return false
	}
	return PositionOverrideParentAndListenerRelativeRounting(p) && Position3DPositionType(p) != 0
}

// --- definition and constant --- //

type PathMode = u8

const SizeOfPositionVertice = 4 * Size32
const SizeOfPositionPlayList = 2 * Size32
const SizeOfSpatialAutomation = 3 * Size32
