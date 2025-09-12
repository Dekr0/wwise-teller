package wwise

import "sync"

// --- struct definition --- //

type PositionParam struct {
	SettingVector        u8
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
	X f32
	Y f32
	Z f32
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
	dMu sync.Mutex

	ActorMixerPositionParam map[u32]*PositionParam
}

// --- definition and constant --- //

type PathMode = u8
