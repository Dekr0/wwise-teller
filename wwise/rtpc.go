package wwise

import (
	"sync"

	uio "github.com/Dekr0/unwise/io"
)

type RTPC struct {
	Id      []u32
	Type    []u8
	Accum   []u8
	ParamId []uio.V128
	CurveId []u32
	Scaling []u8
}

type RTPCGraph struct {
	PointX    []f32
	PointY    []f32
	PointLerp []u32
}

type RTPCComponent struct {
	dMu sync.Mutex

	ActorMixerRTPC map[u32]*RTPC
}

type RTPCGraphComponent struct {
	dMu sync.Mutex

	ActorMixerRTPC map[u32][]RTPCGraph
}
