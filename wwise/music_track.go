package wwise

import (
	"slices"

	wio "github.com/Dekr0/wwise-teller/wio"
)

var TrackTypeName []string = []string{
	"Normal",
	"Random",
	"Sequence",
	"Switch",
}

type MusicTrack struct {
	HircObj

	Id              uint32
	OverrideFlags   uint8
	// NumSources
	Sources         []BankSourceData
	// NumPlayListItem
	PlayListItems   []MusicTrackPlayListItem
	NumSubTrack     uint32
	// NumClipAutomationItem uint32
	ClipAutomations []ClipAutomation
	BaseParam       BaseParameter

	TrackType       uint8
	SwitchParam     MusicTrackSwitchParam
	TransitionParam MusicTrackTransitionParam

	LookAheadTime   int32
}

func (h *MusicTrack) OverrideParentMIDITempo() bool {
	return wio.GetBit(h.OverrideFlags, 1)
}

func (h *MusicTrack) SetOverrideParentMIDITempo(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 1, set)
}

func (h *MusicTrack) OverrideParentMIDITarget() bool {
	return wio.GetBit(h.OverrideFlags, 2)
}

func (h *MusicTrack) SetOverrideParentMIDITarget(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 2, set)
}

func (h *MusicTrack) MidiTargetTypeBus() bool {
	return wio.GetBit(h.OverrideFlags, 3)
}

func (h *MusicTrack) SetMidiTargetTypeBus(set bool) {
	h.OverrideFlags = wio.SetBit(h.OverrideFlags, 3, set)
}

func (h *MusicTrack) AddNewAutomation() {
	h.ClipAutomations = append(h.ClipAutomations, ClipAutomation{})
}

func (h *MusicTrack) RemoveAutomation(i int) {
	h.ClipAutomations = slices.Delete(h.ClipAutomations, i, i + 1)
}

func (h *MusicTrack) Encode() []byte {
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicTrack))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendByte(h.OverrideFlags)
	w.Append(uint32(len(h.Sources)))
	for _, s := range h.Sources {
		w.AppendBytes(s.Encode())
	}
	w.Append(uint32(len(h.PlayListItems)))
	for _, p := range h.PlayListItems {
		w.Append(p)
	}
	if len(h.PlayListItems) > 0 {
		w.Append(h.NumSubTrack)
	}
	w.Append(uint32(len(h.ClipAutomations)))
	for _, c := range h.ClipAutomations {
		w.AppendBytes(c.Encode())
	}
	w.AppendBytes(h.BaseParam.Encode())
	w.AppendByte(h.TrackType)
	if h.UseSwitchAndTransition() {
		w.AppendBytes(h.SwitchParam.Encode())
		w.Append(h.TransitionParam)
	}
	w.Append(h.LookAheadTime)
	return w.BytesAssert(int(size))
}

func (h *MusicTrack) DataSize() uint32 {
	dataSize := uint32(4 + 1 + 4)
	for _, s := range h.Sources {
		dataSize += s.Size()
	}
	dataSize += 4 + uint32(len(h.PlayListItems)) * SizeOfMusicTrackPlayListItem
	if len(h.PlayListItems) > 0 {
		dataSize += 4
	}
	dataSize += 4
	for _, c := range h.ClipAutomations {
		dataSize += uint32(c.Size())
	}
	dataSize += h.BaseParam.Size() + 1 
	if h.UseSwitchAndTransition() {
		dataSize += uint32(h.SwitchParam.Size()) + SizeOfMusicTrackTransitionParam 
	}
	dataSize += 4
	return dataSize
}

func (h *MusicTrack) BaseParameter() *BaseParameter { return &h.BaseParam }

func (h *MusicTrack) HircType() HircType { return HircTypeMusicTrack }

func (h *MusicTrack) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicTrack) IsCntr() bool { return false }

func (h *MusicTrack) NumLeaf() int { return 0 }

func (h *MusicTrack) ParentID() uint32 { return h.BaseParam.DirectParentId }

func (h *MusicTrack) AddLeaf(o HircObj) { panic("") }

func (h *MusicTrack) RemoveLeaf(o HircObj) { panic("") }

func (h *MusicTrack) UseSwitchAndTransition() bool {
	return h.TrackType == 0x03
}

const SizeOfMusicTrackPlayListItem = 44
type MusicTrackPlayListItem struct {
	TrackID         uint32
	SourceID        uint32
	EventID         uint32
	PlayAt          float64
	BeginTrimOffset float64
	EndTrimOffset   float64
	SrcDuration     float64
}

var ClipAutomationTypeName []string = []string{
  	"Volume",
  	"LPF",
  	"HPF",
  	"FadeIn",
  	"FadeOut",
}


type ClipAutomation struct {
	ClipIndex       uint32
	AutoType        uint32
	// NumPoints       uint32
	RTPCGraphPoints []RTPCGraphPoint
}

func (c *ClipAutomation) AddRTPCGraphPoint() {
	c.RTPCGraphPoints = append(c.RTPCGraphPoints, RTPCGraphPoint{})
}

func (c *ClipAutomation) RemoveRTPCGraphPoint(i int) {
	c.RTPCGraphPoints = slices.Delete(c.RTPCGraphPoints, i, i + 1)
}

func (c *ClipAutomation) Encode() []byte {
	dataSize := c.Size()
	w := wio.NewWriter(uint64(dataSize))
	w.Append(c.ClipIndex)
	w.Append(c.AutoType)
	w.Append(uint32(len(c.RTPCGraphPoints)))
	for _, r := range c.RTPCGraphPoints {
		w.AppendBytes(r.Encode())
	}
	return w.BytesAssert(int(dataSize))
}

func (c *ClipAutomation) Size() uint16 {
	return 4 + 4 + 4 + uint16(len(c.RTPCGraphPoints)) * SizeOfRTPCGraphPoint
}

var MusicSwitchGroupTypeName []string = []string{
	"Switch",
	"State",
}

type MusicTrackSwitchParam struct {
	GroupType        uint8
	GroupID          uint32
	DefaultSwitch    uint32
	// NumSwitchAssoc uint32
	SwitchAssociates []uint32
}

func (m *MusicTrackSwitchParam) Encode() []byte {
	dataSize := m.Size()
	w := wio.NewWriter(uint64(dataSize))
	w.Append(m.GroupType)
	w.Append(m.GroupID)
	w.Append(m.DefaultSwitch)
	w.Append(uint32(len(m.SwitchAssociates)))
	for _, s := range m.SwitchAssociates {
		w.Append(s)
	}
	return w.BytesAssert(int(dataSize))
}

func (m *MusicTrackSwitchParam) Size() uint16 {
	return 1 + 4 + 4 + 4 + uint16(len(m.SwitchAssociates)) * 4
}

var SyncTypeName []string = []string{
  	"Immediate",
  	"NextGrid",
  	"NextBar",
  	"NextBeat",
  	"NextMarker",
  	"NextUserMarker",
  	"EntryMarker",
  	"ExitMarker",
  	"ExitNever",
  	"LastExitPosition",
}
const NumSyncType = 10
const SizeOfMusicTrackTransitionParam = 32
type MusicTrackTransitionParam struct {
	SrcTransitionTime  int32
	SrcFadeCurve       uint32
	SrcFadeOffset      int32
	SyncType           uint32
	CueFilterHash      uint32	
	DestTransitionTime int32
	DestFadeCurve      uint32
	DestFadeOffset     int32
}
