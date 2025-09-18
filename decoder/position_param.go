package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodePositionParam(r io.Reader, o order) (p *wwise.PositionParam) {
	p = wwise.AllocPositionParam() 
	p.SettingVector = uio.U8P(r, o)
	if !wwise.PositionOverrideParentAndListenerRelativeRounting(p) {
		return p
	}
	p.SpatialSettingVector = uio.U8P(r, o)
	if !wwise.PositionHas3D(p) {
		return p
	}
	p.PathMode = uio.U8P(r, o)
	p.TransitionTime = uio.I32P(r, o)

	n := uio.U32P(r, o)
	for range n {
		p.PositionVertices.PositionVerticesX = append(p.PositionVertices.PositionVerticesX, uio.F32P(r, o))
		p.PositionVertices.PositionVerticesY = append(p.PositionVertices.PositionVerticesY, uio.F32P(r, o))
		p.PositionVertices.PositionVerticesZ = append(p.PositionVertices.PositionVerticesZ, uio.F32P(r, o))
		p.PositionVertices.PositionVerticesDuration = append(p.PositionVertices.PositionVerticesDuration, uio.I32P(r, o))
	}

	n = uio.U32P(r, o)
	for range n {
		p.PositionPlayList.VerticesOffset = append(p.PositionPlayList.VerticesOffset, uio.U32P(r, o))
		p.PositionPlayList.NumVertices = append(p.PositionPlayList.NumVertices, uio.U32P(r, o))
	}

	for range n {
		p.SpatialAutomation.RangeX = append(p.SpatialAutomation.RangeX, uio.F32P(r,o))
		p.SpatialAutomation.RangeY = append(p.SpatialAutomation.RangeY, uio.F32P(r,o))
		p.SpatialAutomation.RangeZ = append(p.SpatialAutomation.RangeZ, uio.F32P(r,o))
	}

	return p
}
