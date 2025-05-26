package wwise

import (
	"errors"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

type PositioningParam struct {
	BitsPositioning uint8 // U8x
	Bits3D uint8 // U8x
	PathMode uint8 // U8x
	TransitionTime int32 // s32
	// NumPositionVertices uint32 // u32
	PositionVertices []PositionVertex // NumPositionVertices * sizeof(PositionVertex)
	// NumPositionPlayListItem uint32 // u32
	PositionPlayListItems []PositionPlayListItem // NumPositionPlayListItem * sizeof(PositionPlayListItem)
	Ak3DAutomationParams []Ak3DAutomationParam // NumPositionPlayListItem * sizeof(Ak3DAutomationParams)
}

func NewPositioningParam() *PositioningParam {
	return &PositioningParam{
		0, 0, 0, 0, 
		[]PositionVertex{}, 
		[]PositionPlayListItem{},
		[]Ak3DAutomationParam{},
	}
}

func (p *PositioningParam) HasPositioning() bool {
	return (p.BitsPositioning >> 0) & 1 != 0
}

func (p *PositioningParam) Has3D() bool {
	if !p.HasPositioning() {
		return false
	}
	return (p.BitsPositioning >> 1) & 1 != 0
}

func (p *PositioningParam) HasPositioningAnd3D() bool {
	return p.HasPositioning() && p.Has3D()
}

func (p *PositioningParam) Type3DPosition() (uint8, error) {
	if !p.HasPositioning() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter is not enable.",
		)
	}
	if !p.Has3D() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter does not enable 3D setting",
		)
	}
	return (p.BitsPositioning >> 5) & 3, nil
}

func (p *PositioningParam) HasAutomation() bool {
	if !p.HasPositioningAnd3D() {
		return false
	}
	_3DPositioningType, err := p.Type3DPosition()
	assert.Nil(err, "Error of Get3DPositionType")
	return p.HasPositioningAnd3D() && _3DPositioningType != 0
}

func (p *PositioningParam) Encode() []byte {
	p.assert()

	if !p.HasPositioning() || !p.HasPositioningAnd3D() {
		size := p.Size()
		w := wio.NewWriter(uint64(size))
		w.AppendByte(p.BitsPositioning)
		return w.BytesAssert(int(size))
	}

	if !p.HasAutomation() {
		size := p.Size()
		w := wio.NewWriter(uint64(size))
		w.AppendByte(p.BitsPositioning)
		w.AppendByte(p.Bits3D)
		return w.BytesAssert(int(size))
	}

	size := p.Size()
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

func (p *PositioningParam) Size() uint32 {
	if !p.HasPositioning() || !p.HasPositioningAnd3D() {
		return 1
	}
	if !p.HasAutomation() {
		return 2
	}
	return uint32(1 + 1 + 1 + 4 + 4 + len(p.PositionVertices) * SizeOfPositionVertex + 4 + len(p.PositionPlayListItems) * SizeOfPositionPlayListItem + len(p.PositionPlayListItems) * SizeOfAk3DAutomationParam)
}

/* Will Panic */
func (p *PositioningParam) assert() {
	/* TODO, document assertion */
	if !p.HasPositioning() {
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
	if !p.HasPositioningAnd3D() {
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
	if !p.HasAutomation() {
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

