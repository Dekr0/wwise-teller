package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseMusicSegment(size uint32, r *wio.Reader, v int) *wwise.MusicSegment {
	assert.Equal(0, r.Pos(), "Music Segment parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicSegment{}

	m.Id = r.U32Unsafe()
	m.OverrideFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r, v)
	ParseContainer(r, &m.Children, v)
	ParseMeterInfo(r, &m.MeterInfo, v)
	m.Stingers = make([]wwise.Stinger, r.U32Unsafe())
	ParseStingers(r, m.Stingers)

	m.Duration = r.F64Unsafe()
	m.Markers = make([]wwise.MusicSegmentMarker, r.U32Unsafe())
	ParseMusicSegmentMarkers(r, m.Markers)

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)

	return &m
}

func ParseMeterInfo(r *wio.Reader, m *wwise.MeterInfo, v int) {
	m.GridPeriod = r.F64Unsafe()
	m.GridOffset = r.F64Unsafe()
	m.Tempo = r.F32Unsafe()
	m.TimeSigNumBeatsBar = r.U8Unsafe()
	m.TimeSigBeatVal = r.U8Unsafe()
	m.MeterInfoFlag = r.U8Unsafe()
}

func ParseStingers(r *wio.Reader, stingers []wwise.Stinger) {
	for i := range stingers {
		stingers[i].TriggerID = r.U32Unsafe()
		stingers[i].SegmentID = r.U32Unsafe()
		stingers[i].SyncPlayAt = r.U32Unsafe()
		stingers[i].CueFilterHash = r.U32Unsafe()
		stingers[i].DontRepeatTime = r.I32Unsafe()
		stingers[i].NumSegmentLookAhead = r.U32Unsafe()
	}
}

func ParseMusicSegmentMarkers(r *wio.Reader, markers []wwise.MusicSegmentMarker) {
	for i := range markers {
		markers[i].ID = r.U32Unsafe()
		markers[i].Position = r.F64Unsafe()
		markers[i].MarkerName = r.StzUnsafe()
	}
}

func ParseMusicTrack(size uint32, r *wio.Reader, v int) *wwise.MusicTrack {
	assert.Equal(0, r.Pos(), "Music Segment parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicTrack{
		Id:            r.U32Unsafe(),
		OverrideFlags: r.U8Unsafe(),
		Sources:       make([]wwise.BankSourceData, r.U32Unsafe()),
	}
	for i := range m.Sources {
		m.Sources[i] = ParseBankSourceData(r, v)
	}

	m.PlayListItems = make([]wwise.MusicTrackPlayListItem, r.U32Unsafe())
	ParseMusicTrackPlayList(r, m.PlayListItems, v)
	if len(m.PlayListItems) > 0 {
		m.NumSubTrack = r.U32Unsafe()
	}

	m.ClipAutomations = make([]wwise.ClipAutomation, r.U32Unsafe())
	ParseClipAutomation(r, m.ClipAutomations, v)

	m.BaseParam = *ParseBaseParam(r, v)

	m.TrackType = r.U8Unsafe()
	if m.UseSwitchAndTransition() {
		ParseMusicTrackSwitchParam(r, &m.SwitchParam, v)
		ParseMusicTrackTransitionParam(r, &m.TransitionParam, v)
	}

	m.LookAheadTime = r.I32Unsafe()

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &m
}

func ParseMusicTrackPlayList(r *wio.Reader, p []wwise.MusicTrackPlayListItem, v int) {
	for i := range p {
		p[i].TrackID = r.U32Unsafe()
		p[i].SourceID = r.U32Unsafe()
		p[i].EventID = r.U32Unsafe()
		p[i].PlayAt = r.F64Unsafe()
		p[i].BeginTrimOffset = r.F64Unsafe()
		p[i].EndTrimOffset = r.F64Unsafe()
		p[i].SrcDuration = r.F64Unsafe()
	}
}

func ParseClipAutomation(r *wio.Reader, cs []wwise.ClipAutomation, v int) {
	for i := range cs {
		cs[i].ClipIndex = r.U32Unsafe()
		cs[i].AutoType = r.U32Unsafe()
		cs[i].RTPCGraphPoints = make([]wwise.RTPCGraphPoint, r.U32Unsafe())
		for j := range cs[i].RTPCGraphPoints {
			cs[i].RTPCGraphPoints[j].From = r.F32Unsafe()
			cs[i].RTPCGraphPoints[j].To = r.F32Unsafe()
			cs[i].RTPCGraphPoints[j].Interp = r.U32Unsafe()
		}
	}
}

func ParseMusicTrackSwitchParam(r *wio.Reader, s *wwise.MusicTrackSwitchParam, v int) {
	s.GroupType = r.U8Unsafe()
	s.GroupID = r.U32Unsafe()
	s.DefaultSwitch = r.U32Unsafe()
	// NumSwitchAssoc uint32
	s.SwitchAssociates = make([]uint32, r.U32Unsafe())
	for i := range s.SwitchAssociates {
		s.SwitchAssociates[i] = r.U32Unsafe()
	}
}

func ParseMusicTrackTransitionParam(r *wio.Reader, t *wwise.MusicTrackTransitionParam, v int) {
	t.SrcTransitionTime = r.I32Unsafe()
	t.SrcFadeCurve = r.U32Unsafe()
	t.SrcFadeOffset = r.I32Unsafe()
	t.SyncType = r.U32Unsafe()
	t.CueFilterHash = r.U32Unsafe()
	t.DestTransitionTime = r.I32Unsafe()
	t.DestFadeCurve = r.U32Unsafe()
	t.DestFadeOffset = r.I32Unsafe()
}

func ParseMusicRanSeqCntr(size uint32, r *wio.Reader, v int) *wwise.MusicRanSeqCntr {
	assert.Equal(0, r.Pos(), "Muisc random / sequence container parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicRanSeqCntr{}
	m.Id = r.U32Unsafe()
	m.OverwriteFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r, v)
	ParseContainer(r, &m.Children, v)
	ParseMeterInfo(r, &m.MeterInfo, v)
	m.Stingers = make([]wwise.Stinger, r.U32Unsafe())
	ParseStingers(r, m.Stingers)
	m.TransitionRules = make([]wwise.MusicTransitionRule, r.U32Unsafe())
	ParseTransitionRules(r, m.TransitionRules, v)
	totalNumPlayListNodes := r.U32Unsafe()
	ParseMusicPlayListNodes(r, &m.PlayListNode, &totalNumPlayListNodes, v)
	if totalNumPlayListNodes > 0 {
		panic("Number play list node checker is not equal to zero")
	}

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &m
}

func ParseTransitionRules(r *wio.Reader, rules []wwise.MusicTransitionRule, v int) {
	for i := range rules {
		rule := &rules[i]
		rule.SrcIDs = make([]uint32, r.U32Unsafe())
		for j := range rule.SrcIDs {
			rule.SrcIDs[j] = r.U32Unsafe()
		}
		rule.DestIDs = make([]uint32, r.U32Unsafe())
		for j := range rule.DestIDs {
			rule.DestIDs[j] = r.U32Unsafe()
		}

		rule.TransitionSourceRule.TransitionTime = r.I32Unsafe()
		rule.TransitionSourceRule.FadeCurve = r.U32Unsafe()
		rule.TransitionSourceRule.FadeOffset = r.U32Unsafe()
		rule.TransitionSourceRule.SyncType = r.U32Unsafe()
		rule.TransitionSourceRule.CueFilterHash = r.U32Unsafe()
		rule.TransitionSourceRule.PlayPostExit = r.U8Unsafe()

		rule.TransitionDestRule.TransitionTime = r.I32Unsafe()
		rule.TransitionDestRule.FadeCurve = r.U32Unsafe()
		rule.TransitionDestRule.FadeOffset = r.U32Unsafe()
		rule.TransitionDestRule.CueFilterHash = r.U32Unsafe()
		rule.TransitionDestRule.JumpToID = r.U32Unsafe()
		rule.TransitionDestRule.JumpToType = r.U16Unsafe()
		rule.TransitionDestRule.EntryType = r.U16Unsafe()
		rule.TransitionDestRule.PlayPreEntry = r.U8Unsafe()
		rule.TransitionDestRule.DestMatchSourceCueName = r.U8Unsafe()

		rule.AllocTransitionObjFlag = r.U8Unsafe()
		if rule.HasTransitionObj() {
			rule.TransitionObj.SegmentID = r.U32Unsafe()

			rule.TransitionObj.FadeInParam.TransitionTime = r.I32Unsafe()
			rule.TransitionObj.FadeInParam.FadeCurve = r.U32Unsafe()
			rule.TransitionObj.FadeInParam.FadeOffset = r.I32Unsafe()

			rule.TransitionObj.FadeOutParam.TransitionTime = r.I32Unsafe()
			rule.TransitionObj.FadeOutParam.FadeCurve = r.U32Unsafe()
			rule.TransitionObj.FadeOutParam.FadeOffset = r.I32Unsafe()

			rule.TransitionObj.PlayPreEntry = r.U8Unsafe()
			rule.TransitionObj.PlayPostExit = r.U8Unsafe()
		}
	}
}

func ParseMusicPlayListNodes(
	r *wio.Reader, p *wwise.MusicPlayListNode, totalNumPlayListNodes *uint32, v int,
) {
	*totalNumPlayListNodes -= 1
	p.SegmentID = r.U32Unsafe()
	p.PlayListItemID = r.U32Unsafe()
	p.PlayListLeafs = make([]wwise.MusicPlayListNode, r.U32Unsafe())
	p.RSType = r.U32Unsafe()
	p.Loop = r.I16Unsafe()
	p.LoopMin = r.I16Unsafe()
	p.LoopMax = r.I16Unsafe()
	p.Weight = r.U32Unsafe()
	p.AvoidRepeatCount = r.U16Unsafe()
	p.UsingWeight = r.U8Unsafe()
	p.Shuffle = r.U8Unsafe()
	for i := range p.PlayListLeafs {
		ParseMusicPlayListNodes(r, &p.PlayListLeafs[i], totalNumPlayListNodes, v)
	}
}

func ParseMusicSwitchCntr(size uint32, r *wio.Reader, v int) (m *wwise.MusicSwitchCntr) {
	assert.Equal(0, r.Pos(), "Music switch container parser position doesn't start at position 0.")
	begin := r.Pos()

	m = &wwise.MusicSwitchCntr{}
	m.Id = r.U32Unsafe()
	m.OverwriteFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r, v)
	ParseContainer(r, &m.Children, v)
	ParseMeterInfo(r, &m.MeterInfo, v)
	m.Stingers = make([]wwise.Stinger, r.U32Unsafe())
	ParseStingers(r, m.Stingers)
	m.TransitionRules = make([]wwise.MusicTransitionRule, r.U32Unsafe())
	ParseTransitionRules(r, m.TransitionRules, v)
	m.ContinuePlayBack = r.U8Unsafe()
	m.DecisionTreeData = r.ReadAllUnsafe()

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return m
}
