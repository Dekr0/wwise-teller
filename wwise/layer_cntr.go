package wwise

import (
	"log/slog"

	"github.com/Dekr0/wwise-teller/wio"
)

type LayerCntr struct {
	HircObj
	Id uint32
	BaseParam *BaseParameter
	Container *Container

	// NumLayers uint32 // u32
	
	Layers []Layer
	IsContinuousValidation uint8 // U8x
}

func (l *LayerCntr) Encode() []byte {
	dataSize := l.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeLayerCntr))
	w.Append(dataSize)
	w.Append(l.Id)
	w.AppendBytes(l.BaseParam.Encode())
	w.AppendBytes(l.Container.Encode())
	w.Append(uint32(len(l.Layers)))
	for _, i := range l.Layers {
		w.AppendBytes(i.Encode())
	}
	w.AppendByte(l.IsContinuousValidation)
	return w.BytesAssert(int(size))
}

func (l *LayerCntr) DataSize() uint32 {
	size := 4 + l.BaseParam.Size() + l.Container.Size() + 4
	for _, i := range l.Layers {
		size += i.Size()
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

type Layer struct {
	Id uint32 // tid
	InitialRTPC RTPC
	RTPCId uint32 // tid
	RTPCType uint8 // U8x

	// NumAssoc uint32 // u32

	LayerRTPCs []LayerRTPC
}

func (l *Layer) Encode() []byte {
	size := l.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(l.Id)
	w.AppendBytes(l.InitialRTPC.Encode())
	w.Append(l.RTPCId)
	w.AppendByte(l.RTPCType)
	w.Append(uint32(len(l.LayerRTPCs)))
	for _, i := range l.LayerRTPCs {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (l *Layer) Size() uint32 {
	size := uint32(4 + l.InitialRTPC.Size() + 4 + 1 + 4)
	for _, i := range l.LayerRTPCs {
		size += i.Size()
	}
	return size
}

type LayerRTPC struct {
	AssociatedChildID uint32 // tid

	// NumRTPCGraphPoints / CurveSize uint32 // u32

	RTPCGraphPoints []RTPCGraphPoint
}

func (l *LayerRTPC) Encode() []byte {
	size := l.Size()
	w := wio.NewWriter(uint64(size))
	w.Append(l.AssociatedChildID)
	w.Append(uint32(len(l.RTPCGraphPoints)))
	for _, i := range l.RTPCGraphPoints {
		w.Append(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (l *LayerRTPC) Size() uint32 {
	return uint32(4 + 4 + len(l.RTPCGraphPoints) * SizeOfRTPCGraphPoint)
}
