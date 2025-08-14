package wwise

import (
	"errors"
	"slices"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type PositioningParam struct {
	BitsPositioning                  uint8 // U8x
	FallbackBitsPositioning          uint8 // U8x
	// When 3D spatialization (bit 0 to bit 2) is none -> remove positioning type blend
	Bits3D                           uint8 // U8x
	FallbackBits3D                   uint8 // U8x
	PathMode                         uint8 // U8x
	FallbackPathMode                 uint8 // U8x
	TransitionTime                   int32 // s32
	FallbackTransitionTime           int32 // s32
	// NumPositionVertices          uint32 // u32
	PositionVertices               []PositionVertex // NumPositionVertices * sizeof(PositionVertex)
	FallbackPositionVertices       []PositionVertex // NumPositionVertices * sizeof(PositionVertex)
	// NumPositionPlayListItem       uint32 // u32
	PositionPlayListItems          []PositionPlayListItem // NumPositionPlayListItem * sizeof(PositionPlayListItem)
	FallbackPositionPlayListItems  []PositionPlayListItem // NumPositionPlayListItem * sizeof(PositionPlayListItem)
	Ak3DAutomationParams           []Ak3DAutomationParam  // NumPositionPlayListItem * sizeof(Ak3DAutomationParams)
	FallbackAk3DAutomationParams   []Ak3DAutomationParam  // NumPositionPlayListItem * sizeof(Ak3DAutomationParams)
}

func (p *PositioningParam) Clone() PositioningParam {
	cp := PositioningParam{}
	cp.BitsPositioning = p.BitsPositioning
	cp.FallbackBitsPositioning = p.FallbackBitsPositioning
	cp.Bits3D = p.Bits3D
	cp.FallbackBits3D = p.FallbackBits3D
	cp.PathMode = p.PathMode
	cp.FallbackPathMode = p.FallbackPathMode
	cp.TransitionTime = p.TransitionTime
	cp.FallbackTransitionTime = p.FallbackTransitionTime
	cp.PositionVertices = slices.Clone(p.PositionVertices)
	cp.FallbackPositionVertices = slices.Clone(p.FallbackPositionVertices)
	cp.PositionPlayListItems = slices.Clone(p.PositionPlayListItems)
	cp.FallbackPositionPlayListItems = slices.Clone(p.FallbackPositionPlayListItems)
	cp.Ak3DAutomationParams = slices.Clone(p.Ak3DAutomationParams)
	cp.FallbackAk3DAutomationParams = slices.Clone(p.FallbackAk3DAutomationParams)
	return cp
}

func (p *PositioningParam) OverrideParent() bool {
	return wio.GetBit(p.BitsPositioning, 0)
}

// Reset everything to default -> remove all properties related to positioning 
// -> remove attached attenuation
func (p *PositioningParam) SetOverrideParent(set bool) {
	if !set {
		p.BitsPositioning = p.FallbackBitsPositioning
		p.Bits3D = p.FallbackBits3D
		p.PathMode = p.FallbackPathMode
		p.TransitionTime = p.FallbackTransitionTime
		p.PositionVertices = slices.Clone(p.FallbackPositionVertices)
		p.PositionPlayListItems = slices.Clone(p.FallbackPositionPlayListItems)
		p.Ak3DAutomationParams = slices.Clone(p.FallbackAk3DAutomationParams)
	} else {
		p.BitsPositioning = wio.SetBit(p.BitsPositioning, 0, set)
	}
}

func (p *PositioningParam) ListenerRelativeRouting() bool {
	if !p.OverrideParent() {
		return false
	}
	return wio.GetBit(p.BitsPositioning, 1)
}

// Remove all property related to listener relative routing; remove attached 
// attenuation; fallback 3D positioning.
func (p *PositioningParam) SetListenerRelativeRouting(set bool) {
	if !p.OverrideParent() {
		return
	}
	p.BitsPositioning = wio.SetBit(p.BitsPositioning, 1, set)
	if !set {
		p.Bits3D = p.FallbackBits3D
	}
}

func (p *PositioningParam) OverrideParentAndHasListenerRelativeRouting() bool {
	return p.OverrideParent() && p.ListenerRelativeRouting()
}

// bit 2 and bit 3 are for speaker panning
func (p *PositioningParam) Get3DPositionType() (uint8, error) {
	if !p.OverrideParent() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter is not enable.",
		)
	}
	if !p.ListenerRelativeRouting() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter does not enable 3D setting",
		)
	}
	return (p.BitsPositioning >> 5) & 3, nil
}

func (p *PositioningParam) Attenuation() bool {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return false
	}
	return wio.GetBit(p.Bits3D, 3)
}

// Remove attached attenuation ID
func (p *PositioningParam) EnableAttenuation(enable bool) {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return
	}
	p.Bits3D = wio.SetBit(p.Bits3D, 3, enable)
}

func (p *PositioningParam) Has3DAutomation() bool {
	if !p.OverrideParentAndHasListenerRelativeRouting() {
		return false
	}
	_3DPositioningType, err := p.Get3DPositionType()
	assert.Nil(err, "Error of Get3DPositionType")
	return p.OverrideParentAndHasListenerRelativeRouting() && _3DPositioningType != 0
}

func (p *PositioningParam) HoldEmitterPositionAndOrientation() bool {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return false
	}
	return wio.GetBit(p.Bits3D, 4)
}

func (p *PositioningParam) SetHoldEmitterPositionAndOrientation(set bool) {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return
	}
	p.Bits3D = wio.SetBit(p.Bits3D, 4, set)
}

func (p *PositioningParam) Diffraction() bool {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return false
	}
	return wio.GetBit(p.Bits3D, 6)
}

func (p *PositioningParam) EnableDiffraction(set bool) {
	if !p.OverrideParent() || !p.ListenerRelativeRouting() {
		return
	}
}

func (p *PositioningParam) Encode(v int) []byte {
	p.assert()

	if !p.OverrideParent() || !p.OverrideParentAndHasListenerRelativeRouting() {
		size := p.Size(v)
		w := wio.NewWriter(uint64(size))
		w.AppendByte(p.BitsPositioning)
		return w.BytesAssert(int(size))
	}

	if !p.Has3DAutomation() {
		size := p.Size(v)
		w := wio.NewWriter(uint64(size))
		w.AppendByte(p.BitsPositioning)
		w.AppendByte(p.Bits3D)
		return w.BytesAssert(int(size))
	}

	size := p.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(p.BitsPositioning)
	w.Append(p.Bits3D)
	w.Append(p.PathMode)
	w.Append(p.TransitionTime)
	w.Append(uint32(len(p.PositionVertices)))
	for _, v := range p.PositionVertices { w.Append(v) }
	w.Append(uint32(len(p.PositionPlayListItems)))
	for _, i := range p.PositionPlayListItems { w.Append(i) }
	for _, p := range p.Ak3DAutomationParams { w.Append(p) }

	return w.BytesAssert(int(size))
}

func (p *PositioningParam) Size(v int) uint32 {
	if !p.OverrideParent() || !p.OverrideParentAndHasListenerRelativeRouting() {
		return 1
	}
	if !p.Has3DAutomation() {
		return 2
	}
	return uint32(1 + 1 + 1 + 4 + 4 + len(p.PositionVertices) * SizeOfPositionVertex + 4 + len(p.PositionPlayListItems) * SizeOfPositionPlayListItem + len(p.PositionPlayListItems) * SizeOfAk3DAutomationParam)
}

/* Will Panic */
func (p *PositioningParam) assert() {
	/* TODO, document assertion */
	if !p.OverrideParent() {
		assert.True(((p.BitsPositioning >> 1) & 1) == 0, "") // No 3D
		assert.True(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.True(p.Bits3D == 0, "")
		assert.True(p.PathMode == 0, "")
		assert.True(p.TransitionTime == 0, "")
		assert.True(len(p.PositionVertices) == 0, "")
		assert.True(len(p.PositionPlayListItems) == 0, "")
		assert.True(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	if !p.OverrideParentAndHasListenerRelativeRouting() {
		assert.True(((p.BitsPositioning >> 1) & 1) == 0, "") // No 3D
		assert.True(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.True(p.Bits3D == 0, "")
		assert.True(p.PathMode == 0, "")
		assert.True(p.TransitionTime == 0, "")
		assert.True(len(p.PositionVertices) == 0, "")
		assert.True(len(p.PositionPlayListItems) == 0, "")
		assert.True(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	if !p.Has3DAutomation() {
		assert.True(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.True(p.PathMode == 0, "")
		assert.True(p.TransitionTime == 0, "")
		assert.True(len(p.PositionVertices) == 0, "")
		assert.True(len(p.PositionPlayListItems) == 0, "")
		assert.True(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	assert.True(
		len(p.Ak3DAutomationParams) == len(p.PositionPlayListItems),
		"# of position play list item doesn't equal of # of 3D automation parameters",
	)
}

const SizeOfPositionVertex = 16 
type PositionVertex struct {
	X float32 // f32
	Y float32 // f32
	Z float32 // f32
	Duration int32 // s32
}

const SizeOfPositionPlayListItem = 8 
type PositionPlayListItem struct {
	UniqueVerticesOffset uint32 // U32
	INumVertices uint32 // u32
}

const SizeOfAk3DAutomationParam = 12
type Ak3DAutomationParam struct {
	XRange float32 // f32
	YRange float32 // f32
	ZRange float32 // f32
}

