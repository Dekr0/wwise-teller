package wwise

import "github.com/Dekr0/wwise-teller/wio"

type Attenuation struct {
	HircObj
	Id                            uint32
	IsHeightSpreadEnabled         uint8
	IsConeEnabled                 uint8
	InsideDegrees                 float32
	OutsideDegrees                float32
	OutsideVolume                 float32
	LoPass                        float32
	HiPass                        float32
	Curves                        []int8 // [7]uint8 <= 141; [19]uint8 > 141
	// NumCurves                  uint8
	AttenuationConversionTables   []AttenuationConversionTable
	RTPC                          RTPC
}

func (h *Attenuation) Encode(v int) []byte {
	dataSize := h.Size(v)
	size := dataSize + SizeOfHircObjHeader
	w := wio.NewWriter(uint64(size))
	w.Append(uint8(HircTypeAttenuation))
	w.Append(dataSize)
	w.Append(h.Id)
	w.Append(h.IsHeightSpreadEnabled)
	w.Append(h.IsConeEnabled)
	if h.IsConeEnabled & 1 != 0 {
		w.Append(h.InsideDegrees)
		w.Append(h.OutsideDegrees)
		w.Append(h.OutsideVolume)
		w.Append(h.LoPass)
		w.Append(h.HiPass)
	}
	for _, curve := range h.Curves {
		w.Append(curve)
	}
	w.Append(uint8(len(h.AttenuationConversionTables)))
	for _, t := range h.AttenuationConversionTables {
		w.AppendBytes(t.Encode(v))
	}
	w.AppendBytes(h.RTPC.Encode(v))
	return w.BytesAssert(int(size))
}

func (h *Attenuation) ConeEnabled() bool {
	return h.IsConeEnabled & 1 != 0
}

func (h *Attenuation) BaseParameter() *BaseParameter { return nil }

func (h *Attenuation) HircType() HircType { return HircTypeAttenuation }

func (h *Attenuation) HircID() (uint32, error) { return h.Id, nil }

func (h *Attenuation) IsCntr() bool { return false }

func (h *Attenuation) NumLeaf() int { return 0 }

func (h *Attenuation) ParentID() uint32 { return 0 }

func (h *Attenuation) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *Attenuation) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (h *Attenuation) Leafs() []uint32 { return []uint32{} }

func (h *Attenuation) Size(v int) uint32 {
	var size uint32 = 0
	if v <= 141 {
		size = 14
	} else {
		size = 26
	}
	if h.IsConeEnabled & 1 != 0 {
		size += 20
	}
	for _, r := range h.AttenuationConversionTables {
		size += r.Size(v)
	}
	size += h.RTPC.Size(v)
	return size
}

type AttenuationConversionTable struct {
	EnumScaling     uint8
	// Size         uint16
	RTPCGraphPoints []RTPCGraphPoint
}

func (a *AttenuationConversionTable) Size(v int) uint32 {
	return 1 + 2 + uint32(len(a.RTPCGraphPoints)) * SizeOfRTPCGraphPoint
}

func (a *AttenuationConversionTable) Encode(v int) []byte {
	size := a.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(a.EnumScaling)
	w.Append(uint16(len(a.RTPCGraphPoints)))
	for _, p := range a.RTPCGraphPoints {
		w.AppendBytes(p.Encode(v))
	}
	return w.BytesAssert(int(size))
}
