package wwise

import (
	"log/slog"

	"github.com/Dekr0/wwise-teller/wio"
)

type LayerCntr struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container  Container

	// NumLayers uint32 // u32
	
	Layers []Layer
	IsContinuousValidation uint8 // U8x
}

func (l *LayerCntr) Encode(v int) []byte {
	dataSize := l.DataSize(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeLayerCntr))
	w.Append(dataSize)
	w.Append(l.Id)
	w.AppendBytes(l.BaseParam.Encode(v))
	w.AppendBytes(l.Container.Encode(v))
	w.Append(uint32(len(l.Layers)))
	for _, i := range l.Layers {
		w.AppendBytes(i.Encode(v))
	}
	w.AppendByte(l.IsContinuousValidation)
	return w.BytesAssert(int(size))
}

func (l *LayerCntr) DataSize(v int) uint32 {
	size := 4 + l.BaseParam.Size(v) + l.Container.Size(v) + 4
	for _, i := range l.Layers {
		size += i.Size(v)
	}
	return size + 1
}

func (l *LayerCntr) BaseParameter() *BaseParameter { return l.BaseParam }

func (l *LayerCntr) HircID() (uint32, error) { return l.Id, nil }

func (l *LayerCntr) HircType() HircType { return HircTypeLayerCntr }

func (l *LayerCntr) IsCntr() bool { return true }

func (l *LayerCntr) NumLeaf() int { return len(l.Container.Children) }

func (l *LayerCntr) ParentID() uint32 { return l.BaseParam.DirectParentId }

func (l *LayerCntr) AddLeaf(o HircObj) {
	slog.Warn("Adding new leaf is not implemented for layer container.")
}

func (l *LayerCntr) RemoveLeaf(o HircObj) {
	slog.Warn("Removing old leaf is not implemented for layer container.")
}

func (h *LayerCntr) Leafs() []uint32 { return h.Container.Children }

type Layer struct {
	Id uint32 // tid
	InitialRTPC RTPC
	RTPCId uint32 // tid
	RTPCType uint8 // U8x

	// NumAssoc uint32 // u32

	LayerRTPCs []LayerRTPC
}

func (l *Layer) Encode(v int) []byte {
	size := l.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(l.Id)
	w.AppendBytes(l.InitialRTPC.Encode(v))
	w.Append(l.RTPCId)
	w.AppendByte(l.RTPCType)
	w.Append(uint32(len(l.LayerRTPCs)))
	for _, i := range l.LayerRTPCs {
		w.AppendBytes(i.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (l *Layer) Size(v int) uint32 {
	size := uint32(4 + l.InitialRTPC.Size(v) + 4 + 1 + 4)
	for _, i := range l.LayerRTPCs {
		size += i.Size(v)
	}
	return size
}

type LayerRTPC struct {
	AssociatedChildID uint32 // tid

	// NumRTPCGraphPoints / CurveSize uint32 // u32

	RTPCGraphPoints []RTPCGraphPoint
}

func (l *LayerRTPC) Encode(v int) []byte {
	size := l.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(l.AssociatedChildID)
	w.Append(uint32(len(l.RTPCGraphPoints)))
	for _, i := range l.RTPCGraphPoints {
		w.Append(i.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (l *LayerRTPC) Size(v int) uint32 {
	return uint32(4 + 4 + len(l.RTPCGraphPoints) * SizeOfRTPCGraphPoint)
}
