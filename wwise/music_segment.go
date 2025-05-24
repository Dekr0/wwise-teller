package wwise

import "github.com/Dekr0/wwise-teller/wio"

type MusicSegment struct {
	HircObj

	Id             uint32
	OverwriteFlags uint8
	BaseParam      BaseParameter
	Children       Container
	MeterInfo      MeterInfo
	// NumStingers uint32
	Stingers []Stinger
	Duration float64
	// NumMarkers  uint32
	Markers []MusicSegmentMarker
}

func (h *MusicSegment) Encode() []byte {
	dataSize := h.Size()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicSegment))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendByte(h.OverwriteFlags)
	w.AppendBytes(h.BaseParam.Encode())
	w.AppendBytes(h.Children.Encode())
	w.Append(h.MeterInfo)
	w.Append(uint32(len(h.Stingers)))
	for _, s := range h.Stingers {
		w.Append(s)
	}
	w.Append(h.Duration)
	w.Append(uint32(len(h.Markers)))
	for _, m := range h.Markers {
		w.AppendBytes(m.Encode())
	}
	return w.BytesAssert(int(size))
}

func (h *MusicSegment) Size() uint32 {
	dataSize := 4 + 1 + h.BaseParam.Size() + h.Children.Size() + SizeOfMeterInfo + 4 + SizeOfStinger*uint32(len(h.Stingers))
	dataSize += 8 + 4
	for _, m := range h.Markers {
		dataSize += uint32(m.Size())
	}
	return dataSize
}

func (h *MusicSegment) BaseParameter() *BaseParameter { return nil }

func (h *MusicSegment) HircType() HircType { return HircTypeMusicSegment }

func (h *MusicSegment) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicSegment) IsCntr() bool { return true }

func (h *MusicSegment) NumLeaf() int { return len(h.Children.Children) }

func (h *MusicSegment) ParentID() uint32 { return h.BaseParam.DirectParentId }

func (h *MusicSegment) AddLeaf(o HircObj) { panic("") }

func (h *MusicSegment) RemoveLeaf(o HircObj) { panic("") }

const SizeOfMeterInfo = 23

type MeterInfo struct {
	GridPeriod         float64
	GridOffset         float64
	Tempo              float32
	TimeSigNumBeatsBar uint8
	TimeSigBeatVal     uint8
	MeterInfoFlag      uint8
}

const SizeOfStinger = 24

type Stinger struct {
	TriggerID           uint32
	SegmentID           uint32
	SyncPlayAt          uint32
	CueFilterHash       uint32
	DontRepeatTime      int32
	NumSegmentLookAhead uint32
}

type MusicSegmentMarker struct {
	ID         uint32
	Position   float64
	MarkerName []byte
}

func (m *MusicSegmentMarker) Encode() []byte {
	dataSize := m.Size()
	w := wio.NewWriter(uint64(dataSize))
	w.Append(m.ID)
	w.Append(m.Position)
	w.AppendBytes(m.MarkerName)
	return w.BytesAssert(int(dataSize))
}

func (m *MusicSegmentMarker) Size() uint8 {
	return 4 + 8 + uint8(len(m.MarkerName))
}
