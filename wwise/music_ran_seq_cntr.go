package wwise

import "github.com/Dekr0/wwise-teller/wio"

type MusicRanSeqCntr struct {
	Id                  uint32
	OverwriteFlags      uint8
	BaseParam           BaseParameter
	Children            Container
	MeterInfo           MeterInfo
	// NumStingers      uint32
	Stingers            []Stinger
	// NumRules         uint32
	TransitionRules     []MusicTransitionRule
	// NumPlayListItems uint32
	PlayListNode        MusicPlayListNode
}

func (h *MusicRanSeqCntr) Encode(v int) []byte {
	dataSize := h.DataSize(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeMusicRanSeqCntr))
	w.Append(dataSize)
	w.Append(h.Id)
	w.AppendByte(h.OverwriteFlags)
	w.AppendBytes(h.BaseParam.Encode(v))
	w.AppendBytes(h.Children.Encode(v))
	w.Append(h.MeterInfo)
	w.Append(uint32(len(h.Stingers)))
	for _, s := range h.Stingers {
		w.Append(s)
	}
	w.Append(uint32(len(h.TransitionRules)))
	for _, t := range h.TransitionRules {
		w.AppendBytes(t.Encode(v))
	}
	w.Append(h.PlayListNode.NumNodes())
	w.AppendBytes(h.PlayListNode.Encode(v))
	return w.BytesAssert(int(size))
}

func (h *MusicRanSeqCntr) DataSize(v int) uint32 {
	size := 4 + 1 +
		h.BaseParam.Size(v) + 
		h.Children.Size(v) + 
		SizeOfMeterInfo + 
		4 + uint32(len(h.Stingers)) * SizeOfStinger
	size += 4
	for _, t := range h.TransitionRules {
		size += t.Size(v)
	}
	size += 4  + h.PlayListNode.Size(v)
	return size
}

func (h *MusicRanSeqCntr) BaseParameter() *BaseParameter { return &h.BaseParam }

func (h *MusicRanSeqCntr) HircType() HircType { return HircTypeMusicRanSeqCntr }

func (h *MusicRanSeqCntr) HircID() (uint32, error) { return h.Id, nil }

func (h *MusicRanSeqCntr) IsCntr() bool { return true }

func (h *MusicRanSeqCntr) NumLeaf() int { return len(h.Children.Children) }

func (h *MusicRanSeqCntr) ParentID() uint32 { return h.BaseParam.DirectParentId }

func (h *MusicRanSeqCntr) AddLeaf(o HircObj) {}

func (h *MusicRanSeqCntr) RemoveLeaf(o HircObj) {}

func (h *MusicRanSeqCntr) Leafs() []uint32 { return h.Children.Children }

var JumpToTypeName []string = []string{
  	"StartOfPlaylist",
  	"SpecificItem",
  	"LastPlayedSegment",
  	"NextSegment",
}

var EntryTypeName []string = []string{
   "EntryMarker",
   "SameTime",
   "RandomMarker",
   "RandomUserMarker",
   "LastExitTime",
}

const SizeOfTransitionObj = 4 + 12 + 12 + 2
const SizeOfTransitionRulePair = 21 + 26
type MusicTransitionRule struct {
	// NumSrc              uint32
	SrcIDs                 []uint32
	// NumDest             uint32
	DestIDs                []uint32

	TransitionSourceRule   struct {
		TransitionTime int32
		FadeCurve      uint32
		FadeOffset     uint32
		SyncType       uint32
		CueFilterHash  uint32
		PlayPostExit   uint8
	}

	TransitionDestRule     struct {
		TransitionTime         int32
		FadeCurve              uint32
		FadeOffset             uint32
		CueFilterHash          uint32
		JumpToID               uint32
		JumpToType             uint16
		EntryType              uint16
		PlayPreEntry           uint8
		DestMatchSourceCueName uint8
	}

	AllocTransitionObjFlag uint8

	TransitionObj          struct {
		SegmentID   uint32
		FadeInParam struct {
			TransitionTime int32
			FadeCurve      uint32
			FadeOffset     int32
		}
		FadeOutParam struct {
			TransitionTime int32
			FadeCurve      uint32
			FadeOffset     int32
		}
		PlayPreEntry uint8
		PlayPostExit uint8
	}
}

func (m *MusicTransitionRule) Encode(v int) []byte {
	size := m.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(uint32(len(m.SrcIDs)))
	for _, s := range m.SrcIDs {
		w.Append(s)
	}
	w.Append(uint32(len(m.DestIDs)))
	for _, d := range m.DestIDs {
		w.Append(d)
	}
	w.Append(m.TransitionSourceRule)
	w.Append(m.TransitionDestRule)
	w.Append(m.AllocTransitionObjFlag)
	if m.HasTransitionObj() {
		w.Append(m.TransitionObj)
	}
	return w.BytesAssert(int(size))
}

func (m *MusicTransitionRule) Size(int) uint32 {
	size := 4 + uint32(len(m.SrcIDs)) * 4 +
			4 + uint32(len(m.DestIDs)) * 4 +
			SizeOfTransitionRulePair + 1 
	if m.HasTransitionObj() {
		size += SizeOfTransitionObj
	}
	return size
}

func (m *MusicTransitionRule) HasTransitionObj() bool {
	return m.AllocTransitionObjFlag != 0
}

const SizeOfPlayListNode = 4 + 4 + 4 + 4 + 2 + 2 + 2 + 4 + 2 + 1 + 1
type MusicPlayListNode struct {
	SegmentID        uint32
	PlayListItemID   uint32
	// NumChildren   uint32
	RSType           uint32
	Loop             int16
	LoopMin          int16
	LoopMax          int16
	Weight           uint32
	AvoidRepeatCount uint16
	UsingWeight      uint8
	Shuffle          uint8
	PlayListLeafs    []MusicPlayListNode
}

func (p *MusicPlayListNode) Encode(v int) []byte {
	size := p.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(p.SegmentID)
	w.Append(p.PlayListItemID)
	w.Append(uint32(len(p.PlayListLeafs)))
	w.Append(p.RSType)
	w.Append(p.Loop)
	w.Append(p.LoopMin)
	w.Append(p.LoopMax)
	w.Append(p.Weight)
	w.Append(p.AvoidRepeatCount)
	w.AppendByte(p.UsingWeight)
	w.AppendByte(p.Shuffle)
	for _, p := range p.PlayListLeafs {
		w.AppendBytes(p.Encode(v))
	}
	return w.BytesAssert(int(size))
}

func (p *MusicPlayListNode) Size(v int) uint32 {
	size := uint32(SizeOfPlayListNode)
	for _, l := range p.PlayListLeafs {
		size += l.Size(v)
	}
	return size 
}

func (p *MusicPlayListNode) NumNodes() uint32 {
	n := uint32(1)
	for _, l := range p.PlayListLeafs {
		n += l.NumNodes()
	}
	return n
}
