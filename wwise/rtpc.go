package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

type RTPC struct {
	Id      []u32
	Type    []RTPCType
	Accum   []RTPCAccum
	ParamId []uio.V128
	CurveId []u32
	Scaling []Scaling
	Graph   []RTPCGraph
}

type RTPCS struct {
	Id       u32
	Type     RTPCType
	Accum    RTPCAccum
	ParamId *uio.V128
	CurveId  u32
	Scaling  Scaling
}

type RTPCGraph struct {
	PointX    []f32
	PointY    []f32
	PointLerp []u32
}

type RTPCGraphFrame struct {
	PointX    f32
	PointY    f32
	PointLerp u32
}

// 1 RPTC -> 1 RTPC Graph
type RTPCComponent struct {
	BaseRTPC         map[u32]*RTPC
	LayerInitialRTPC map[u32]*RTPC
	LayerRTPC        map[u32]*RTPC
}

// --- allocating and freeing --- //

func AllocateRTPCComponent(baseSize u32, numLayer u32) (r *RTPCComponent) {
	r = &RTPCComponent{}
	if baseSize <= 0 {
		r.BaseRTPC = make(map[u32]*RTPC)
	} else {
		r.BaseRTPC = make(map[u32]*RTPC, baseSize)
	}
	if numLayer <= 0 {
		r.LayerInitialRTPC = make(map[u32]*RTPC)
		r.LayerRTPC = make(map[u32]*RTPC)
	} else {
		r.LayerInitialRTPC = make(map[u32]*RTPC, numLayer)
		r.LayerRTPC = make(map[u32]*RTPC, numLayer)
	}
	return r
}

func AllocateRTPC(size u32) *RTPC {
	return &RTPC{
		Id: make([]u32, size, size),
		Type: make([]RTPCType, size, size),
		Accum: make([]RTPCAccum, size, size),
		ParamId: make([]uio.V128, size, size),
		CurveId: make([]u32, size, size),
		Scaling: make([]Scaling, size, size),
		Graph: make([]RTPCGraph, size, size),
	}
}

func AllocateRTPCGraph(size u32) *RTPCGraph {
	return &RTPCGraph{
		PointX: make([]f32, size, size),
		PointY: make([]f32, size, size),
		PointLerp: make([]u32, size, size),
	}
}

// --- assertion --- //

func AssertRTPC(r *RTPC) error {
	if len(r.Id) != len(r.Type) {
		return fmt.Errorf("# of RTPC id (%d) does not equal to # of RTPC type (%d)",
			len(r.Id), len(r.Type),
		)
	}
	if len(r.Type) != len(r.Accum) {
		return fmt.Errorf("# of RTPC Type (%d) does not equal to # of RTPC Accum (%d)",
			len(r.Type), len(r.Accum),
		)
	}
	if len(r.Accum) != len(r.ParamId)  {
		return fmt.Errorf("# of RTPC Accum (%d) does not equal to # of RTPC parameter id (%d)",
			len(r.Accum), len(r.ParamId),
		)
	}
	if len(r.ParamId) != len(r.CurveId) {
		return fmt.Errorf("# of RTPC parameter id (%d) does not equal to # of RTPC curve id (%d)",
			len(r.ParamId), len(r.CurveId),
		)
	}
	if len(r.CurveId) != len(r.Scaling) {
		return fmt.Errorf("# of RTPC curve id (%d) does not equal to # of RTPC scaling (%d)",
			len(r.CurveId), len(r.Scaling),
		)
	}
	if len(r.Scaling) != len(r.Graph) {
		return fmt.Errorf("# of RTPC scaling value (%d) does not equal to # of RTPC graph (%d)",
			len(r.Scaling), len(r.Graph),
		)
	}
	return AssertRTPCGraphs(r.Graph)
}

func AssertRTPCGraphs(rs []RTPCGraph) error {
	for _, r := range rs {
		if err := AssertRTPCGraph(&r); err != nil {
			return err
		}
	}
	return nil
}

func AssertRTPCGraph(r *RTPCGraph) error {
	if len(r.PointX) != len(r.PointY) {
		return fmt.Errorf("# of X values (%d) does not equal to # of Y values (%d) in RTPC graph",
			len(r.PointX), len(r.PointY),
		)
	}
	if len(r.PointY) != len(r.PointLerp) {
		return fmt.Errorf("# of Y values (%d) does not equal to # of line interpret type (%d) in RTPC graph",
			len(r.PointY), len(r.PointLerp),
		)
	}
	return nil
}

// --- sizing --- //

func SizeOfRTPC(r *RTPC) (size u32) {	
	size = Size16 + SizeOfRTPCCurvePrimitive * u32(len(r.Id))
	for i, paramId := range r.ParamId {
		size += u32(len(paramId.B)) + SizeOfRTPCGraph(&r.Graph[i])
	}
	return size
}

func SizeOfRTPCCurve(paramId *uio.V128) u32 {
	return SizeOfRTPCCurvePrimitive + u32(len(paramId.B))
}

func SizeOfRTPCGraph(r *RTPCGraph) u32 {
	return SizeOfRTPCGraphFrame * u32(len(r.PointX))
}

// --- encoding --- //

func EncodeRTPC(e *HircEncoderCtx, r *RTPC) error {
	curr := e.Count()
	size := SizeOfRTPC(r)
	if err := e.Primitive(u16(len(r.Id))); err != nil {
		return fmt.Errorf("Failed to encode # of RTPC curve: %w", err)
	}
	for i, id := range r.Id {
		if err := EncodeRTPCCurve(e, 
			&RTPCS{
				id,
				r.Type[i],
				r.Accum[i],
				&r.ParamId[i],
				r.CurveId[i],
				r.Scaling[i],
			},
			&r.Graph[i],
		); err != nil {
			return fmt.Errorf("Failed encode %d-th RTPC curve: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

func EncodeRTPCCurve(e *HircEncoderCtx, r *RTPCS, g *RTPCGraph) error {
	curr := e.Count()
	size := SizeOfRTPCCurve(r.ParamId) + SizeOfRTPCGraph(g)
	if err := e.Struct(
		RTPCPartialLow{Id: r.Id, Type: r.Type, Accum: r.Accum},
		SizeOfRTPCCurvePrimitiveLow,
	); err != nil {
		return fmt.Errorf("Failed to encode lower portion of RTPC curve data: %w", err)
	}
	if err := e.Bytes(r.ParamId.B); err != nil {
		return fmt.Errorf("Failed to encode RTPC parameter id %d: %w", r.ParamId.V, err)
	}
	if err := e.Struct(
		RTPCPartialHigh{
			CurveId: r.CurveId,
			Scaling: r.Scaling, 
			NumRTPCFrames: u16(len(g.PointX)),
		},
		SizeOfRTPCCurvePrimitiveHigh,
	); err != nil {
		return fmt.Errorf("Failed to encode high portion of RTPC curve data: %w", err)
	}
	if err := EncodeRTPCGraph(e, g); err != nil {
		return fmt.Errorf("Failed to encode RTPC graph frames: %w", err)
	}
	return e.Expect(curr, size)
}

func EncodeRTPCGraph(e *HircEncoderCtx, g *RTPCGraph) error {
	curr := e.Count()
	size := SizeOfRTPCGraph(g)
	for i, x := range g.PointX {
		payload := RTPCGraphFrame{
			PointX: x,
			PointY: g.PointY[i],
			PointLerp: g.PointLerp[i],
		}
		if err := e.Struct(payload, SizeOfRTPCGraphFrame); err != nil {
			return fmt.Errorf("Failed to encode %d-th RTPC Graph frame: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

// --- component getter and setter --- //

func (c *RTPCComponent) HasBaseRTPC(internalId u32) (in bool) {
	_, in = c.BaseRTPC[internalId]
	return in
}

func (c *RTPCComponent) GetBaseRTPC(internalId u32) (r *RTPC) {
	r, in := c.BaseRTPC[internalId]
	if !in {
		panic("Failed to locate base RTPC")
	}
	return r
}

func (c *RTPCComponent) AddBaseRTPC(internalId u32, r *RTPC) {
	if r == nil {
		panic("Base parameter RTPC is nil")
	}
	if _, in := c.BaseRTPC[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.BaseRTPC[internalId] = r
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetBaseRTPC(internalId u32) (r *RTPC) {
	return h.RTPCComponent.GetBaseRTPC(internalId)
}

func (h *HIRC) AssertBaseRTPCById(internalId u32) error {
	r, in := h.RTPCComponent.BaseRTPC[internalId]
	if !in {
		return fmt.Errorf("Failed to locate base RTPC curves")
	}
	if err := AssertRTPC(r); err != nil {
		return err
	}
	return nil
}

// --- constant and definition --- //

const SizeOfRTPCCurvePrimitiveLow  = Size32 + Size8 * 2
const SizeOfRTPCCurvePrimitiveHigh = Size32 + Size8 + Size16
const SizeOfRTPCCurvePrimitive     = SizeOfRTPCCurvePrimitiveLow + SizeOfRTPCCurvePrimitiveHigh
const SizeOfRTPCGraphFrame         = Size32 * 3

type RTPCType  = u8
type RTPCAccum = u8
type Scaling   = u8

type RTPCPartialLow struct {
	Id    u32
	Type  RTPCType
	Accum RTPCAccum
}

type RTPCPartialHigh struct {
	CurveId       u32
	Scaling       Scaling
	NumRTPCFrames u16
}
