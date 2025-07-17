package wwise

import "github.com/Dekr0/wwise-teller/wio"

type MusicSegment struct {
	HircObj

	Id             uint32
	OverrideFlags uint8
	BaseParam      BaseParameter
	Children       Container
	MeterInfo      MeterInfo
	// NumStingers uint32
	Stingers       []Stinger
	Duration       float64
	// NumMarkers  uint32
	Markers        []MusicSegmentMarker
}

func (h *MusicSegment) OverwriteParentMIDITempo() bool {
	return wio.GetBit(h.OverrideFlags, 1)
}

func (h *MusicSegment) SetOverwriteParentMIDITempo(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 1, set)
}

func (h *MusicSegment) OverwriteParentMIDITarget() bool {
	return wio.GetBit(h.OverrideFlags, 2)
}

func (h *MusicSegment) SetOverwriteParentMIDITarget(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 2, set)
}

func (h *MusicSegment) MidiTargetTypeBus() bool {
	return wio.GetBit(h.OverrideFlags, 3)
}

func (h *MusicSegment) SetMidiTargetTypeBus(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 3, set)
}

func (h *MusicSegment) Encode(v int) []byte {
	dataSize := h.Size(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicSegment))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendByte(h.OverrideFlags)
	w.AppendBytes(h.BaseParam.Encode(v))
	w.AppendBytes(h.Children.Encode(v))
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

func (h *MusicSegment) Size(v int) uint32 {
	dataSize := 4 + 1 + h.BaseParam.Size(v) + h.Children.Size(v) + SizeOfMeterInfo + 4 + SizeOfStinger*uint32(len(h.Stingers))
	dataSize += 8 + 4
	for _, m := range h.Markers {
		dataSize += uint32(m.Size())
	}
	return dataSize
}

func (h *MusicSegment) BaseParameter() *BaseParameter { return &h.BaseParam }

func (h *MusicSegment) HircType() HircType { return HircTypeMusicSegment }

func (h *MusicSegment) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicSegment) IsCntr() bool { return true }

func (h *MusicSegment) NumLeaf() int { return len(h.Children.Children) }

func (h *MusicSegment) ParentID() uint32 { return h.BaseParam.DirectParentId }

func (h *MusicSegment) AddLeaf(o HircObj) {}

func (h *MusicSegment) RemoveLeaf(o HircObj) {}

func (h *MusicSegment) Leafs() []uint32 { return h.Children.Children }

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
