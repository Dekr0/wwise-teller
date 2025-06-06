// All parser used for decoding HIRC chunks expect an io.Reader that operates
// on in memory buffer. All parser will have a side effect of which it will
// advance the cursor position of the accepted io.Reader.
// All hierarchy project parser only consume all data excluding hierarchy object
// header data (hierarchy object type [u8] and hierarchy data size [u32])

package parser

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync/atomic"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const MaxNumParseRoutine = 6

type ParserResult struct {
	i   uint32
	obj wwise.HircObj
}

type Parser func(uint32, *wio.Reader) wwise.HircObj

func ParseHIRC(ctx context.Context, r *wio.Reader, I uint8, T []byte, size uint32) (
	*wwise.HIRC, error,
) {
	assert.Equal(0, r.Pos(), "Parser for HIRC does not start at byte 0.")

	numHircItem := r.U32Unsafe()

	hirc := wwise.NewHIRC(I, T, numHircItem)

	/* sync signal */
	sem := make(chan struct{}, MaxNumParseRoutine)
	i := uint32(0)
	parsed := atomic.Uint32{}

	slog.Debug("Start scanning through all hierarchy object, and scheduling parser",
		"numHircItem", numHircItem,
	)

	for parsed.Load() < numHircItem {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if i >= numHircItem {
			continue
		}
		eHircType := r.U8Unsafe()
		dwSectionSize := r.U32Unsafe()
		if SkipHircObjType(wwise.HircType(eHircType)) {
			unknown := wwise.NewUnknown(
				wwise.HircType(eHircType),
				dwSectionSize,
				r.ReadNUnsafe(uint64(dwSectionSize), 4),
			)
			hirc.HircObjs[i] = unknown

			i += 1
			parsed.Add(1)
			slog.Debug("Skipped hierarchy object",
				"index", i,
				"hircType", eHircType,
				"dwSectionSize", dwSectionSize,
				"readerPosition", r.Pos(),
			)
			continue
		}
		select {
		case sem <- struct{}{}:
			switch wwise.HircType(eHircType) {
			case wwise.HircTypeState:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseState,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeSound:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseSound,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeAction:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseAction,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeEvent:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseEvent,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeRanSeqCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseRanSeqCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeSwitchCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseSwitchCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeActorMixer:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseActorMixer,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeLayerCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseLayerCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicSegment:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicSegment,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicTrack:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicTrack,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicSwitchCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicSwitchCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicRanSeqCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicRanSeqCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeAttenuation:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseAttenuation,
					hirc,
					sem,
					&parsed,
				)
			default:
				panic("Assertion Trap")
			}
			i += 1
			slog.Debug(
				fmt.Sprintf("Scheduled %s parser", wwise.HircTypeName[eHircType]),
				"index", i,
				"hircType", eHircType,
				"dwSectionSize", dwSectionSize,
				"readerPosition", r.Pos(),
			)
		default:
			var obj wwise.HircObj
			switch wwise.HircType(eHircType) {
			case wwise.HircTypeState:
				obj = ParseState(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeSound:
				obj = ParseSound(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeAction:
				obj = ParseAction(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeEvent:
				obj = ParseEvent(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeRanSeqCntr:
				obj = ParseRanSeqCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeSwitchCntr:
				obj = ParseSwitchCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeActorMixer:
				obj = ParseActorMixer(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeLayerCntr:
				obj = ParseLayerCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicSegment:
				obj = ParseMusicSegment(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicTrack:
				obj = ParseMusicTrack(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicSwitchCntr:
				obj = ParseMusicSwitchCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicRanSeqCntr:
				obj = ParseMusicRanSeqCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeAttenuation:
				obj = ParseAttenuation(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			default:
				panic("Assertion Trap")
			}
			AddHircObj(hirc, uint32(i), obj)
			i += 1
			parsed.Add(1)
		}
	}

	assert.Equal(
		size,
		uint32(r.Pos()),
		"There are data that is not consumed after parsing all HIRC blob",
	)

	return hirc, nil
}

// Side effect: It will modify HIRC. Specifically, HIRC.HircObjs and maps for
// different types of hierarchy objects.
func AddHircObj(h *wwise.HIRC, i uint32, obj wwise.HircObj) {
	t := obj.HircType()
	id, err := obj.HircID()
	if err != nil {
		panic(err)
	}
	switch t {
	case wwise.HircTypeState:
		if _, in := h.States.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate state object %d", id))
		}
		h.StateCount.Add(1)
	case wwise.HircTypeSound:
		if _, in := h.Sounds.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate sound object %d", id))
		}
	case wwise.HircTypeAction:
		if _, in := h.Actions.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate action object %d", id))
		}
		h.ActionCount.Add(1)
	case wwise.HircTypeEvent:
		if _, in := h.Events.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate event object %d", id))
		}
		h.EventCount.Add(1)
	case wwise.HircTypeRanSeqCntr:
		if _, in := h.RanSeqCntrs.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate random / sequence container object %d", id))
		}
	case wwise.HircTypeSwitchCntr:
		if _, in := h.SwitchCntrs.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate switch container object %d", id))
		}
	case wwise.HircTypeActorMixer:
		if _, in := h.ActorMixers.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate actor mixer object %d", id))
		}
	case wwise.HircTypeLayerCntr:
		if _, in := h.LayerCntrs.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate layer container object %d", id))
		}
	case wwise.HircTypeMusicSegment:
		if _, in := h.MusicSegments.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate music segment object %d", id))
		}
	case wwise.HircTypeMusicTrack:
		if _, in := h.MusicTracks.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate music track object %d", id))
		}
	case wwise.HircTypeMusicSwitchCntr:
		if _, in := h.MusicSwitchCntr.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate music switch container object %d", id))
		}
	case wwise.HircTypeMusicRanSeqCntr:
		if _, in := h.MusicRanSeqCntr.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate music random sequence container object %d", id))
		}
	case wwise.HircTypeAttenuation:
		if _, in := h.Attenuations.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate attenuation object %d", id))
		}
		h.AttenuationCount.Add(1)
	default:
		panic("Assertion Trap")
	}
	h.HircObjs[i] = obj
	if _, in := h.HircObjsMap.LoadOrStore(id, obj); in {
		panic(fmt.Sprintf("Duplicate hierarchy object %d", id))
	}
	slog.Debug(fmt.Sprintf("Collected %s parser", wwise.HircTypeName[obj.HircType()]))
}

func SkipHircObjType(t wwise.HircType) bool {
	_, find := sort.Find(len(wwise.KnownHircTypes), func(i int) int {
		if t < wwise.KnownHircTypes[i] {
			return -1
		}
		if t == wwise.KnownHircTypes[i] {
			return 0
		}
		return 1
	})
	return !find
}

func ParserRoutine[T wwise.HircObj](
	size uint32,
	i uint32,
	r *wio.Reader,
	f func(uint32, *wio.Reader) T,
	h *wwise.HIRC,
	sem chan struct{},
	parsed *atomic.Uint32,
) {
	AddHircObj(h, i, f(size, r))
	parsed.Add(1)
	<-sem
}

func ParseState(size uint32, r *wio.Reader) *wwise.State {
	assert.Equal(0, r.Pos(), "State parser position doesn't start at position 0.")
	begin := r.Pos()

	state := wwise.State{
		StateID: r.U32Unsafe(),
		StateProps: make([]struct{PID uint16; Val float32}, r.U16Unsafe()),
	}
	for i := range state.StateProps {
		state.StateProps[i].PID, state.StateProps[i].Val = r.U16Unsafe(), r.F32Unsafe()
	}

	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return &state
}

func ParseActorMixer(size uint32, r *wio.Reader) *wwise.ActorMixer {
	assert.Equal(0, r.Pos(), "Actor mixer parser position doesn't start at position 0.")
	begin := r.Pos()
	a := &wwise.ActorMixer{}
	a.Id = r.U32Unsafe()
	a.BaseParam = ParseBaseParam(r)
	a.Container = ParseContainer(r)
	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return a
}

func ParseEvent(size uint32, r *wio.Reader) *wwise.Event {
	assert.Equal(0, r.Pos(), "Layer container parser position doesn't start at position 0.")
	begin := r.Pos()
	e := wwise.Event{}
	e.Id = r.U32Unsafe()
	e.ActionIDs = make([]uint32, r.U8Unsafe())
	for i := range e.ActionIDs {
		e.ActionIDs[i] = r.U32Unsafe()
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &e
}

func ParseLayerCntr(size uint32, r *wio.Reader) *wwise.LayerCntr {
	assert.Equal(0, r.Pos(), "Layer container parser position doesn't start at position 0.")
	begin := r.Pos()
	l := &wwise.LayerCntr{
		Id:                     r.U32Unsafe(),
		BaseParam:              ParseBaseParam(r),
		Container:              ParseContainer(r),
		Layers:                 ParseLayers(r, make([]wwise.Layer, r.U32Unsafe())),
		IsContinuousValidation: r.U8Unsafe(),
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return l
}

func ParseLayers(r *wio.Reader, layers []wwise.Layer) []wwise.Layer {
	for i := range layers {
		l := &layers[i]
		l.Id = r.U32Unsafe()
		ParseRTPC(r, &l.InitialRTPC)
		l.RTPCId = r.U32Unsafe()
		l.RTPCType = r.U8Unsafe()
		l.LayerRTPCs = make([]wwise.LayerRTPC, r.U32Unsafe())
		for j := range l.LayerRTPCs {
			lr := &l.LayerRTPCs[j]
			lr.AssociatedChildID = r.U32Unsafe()
			lr.RTPCGraphPoints = make([]wwise.RTPCGraphPoint, r.U32Unsafe())
			for k := range lr.RTPCGraphPoints {
				lr.RTPCGraphPoints[k].From = r.F32Unsafe()
				lr.RTPCGraphPoints[k].To = r.F32Unsafe()
				lr.RTPCGraphPoints[k].Interp = r.U32Unsafe()
			}
		}
	}
	return layers
}

func ParseMusicSegment(size uint32, r *wio.Reader) *wwise.MusicSegment {
	assert.Equal(0, r.Pos(), "Music Segment parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicSegment{}

	m.Id = r.U32Unsafe()
	m.OverrideFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r)
	m.Children = *ParseContainer(r)
	ParseMeterInfo(r, &m.MeterInfo)
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

func ParseMeterInfo(r *wio.Reader, m *wwise.MeterInfo) {
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

func ParseMusicTrack(size uint32, r *wio.Reader) *wwise.MusicTrack {
	assert.Equal(0, r.Pos(), "Music Segment parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicTrack{
		Id:            r.U32Unsafe(),
		OverrideFlags: r.U8Unsafe(),
		Sources:       make([]wwise.BankSourceData, r.U32Unsafe()),
	}
	for i := range m.Sources {
		m.Sources[i] = ParseBankSourceData(r)
	}

	m.PlayListItems = make([]wwise.MusicTrackPlayListItem, r.U32Unsafe())
	ParseMusicTrackPlayList(r, m.PlayListItems)
	if len(m.PlayListItems) > 0 {
		m.NumSubTrack = r.U32Unsafe()
	}

	m.ClipAutomations = make([]wwise.ClipAutomation, r.U32Unsafe())
	ParseClipAutomation(r, m.ClipAutomations)

	m.BaseParam = *ParseBaseParam(r)

	m.TrackType = r.U8Unsafe()
	if m.UseSwitchAndTransition() {
		ParseMusicTrackSwitchParam(r, &m.SwitchParam)
		ParseMusicTrackTransitionParam(r, &m.TransitionParam)
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

func ParseMusicTrackPlayList(r *wio.Reader, p []wwise.MusicTrackPlayListItem) {
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

func ParseClipAutomation(r *wio.Reader, cs []wwise.ClipAutomation) {
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

func ParseMusicTrackSwitchParam(r *wio.Reader, s *wwise.MusicTrackSwitchParam) {
	s.GroupType = r.U8Unsafe()
	s.GroupID = r.U32Unsafe()
	s.DefaultSwitch = r.U32Unsafe()
	// NumSwitchAssoc uint32
	s.SwitchAssociates = make([]uint32, r.U32Unsafe())
	for i := range s.SwitchAssociates {
		s.SwitchAssociates[i] = r.U32Unsafe()
	}
}

func ParseMusicTrackTransitionParam(r *wio.Reader, t *wwise.MusicTrackTransitionParam) {
	t.SrcTransitionTime = r.I32Unsafe()
	t.SrcFadeCurve = r.U32Unsafe()
	t.SrcFadeOffset = r.I32Unsafe()
	t.SyncType = r.U32Unsafe()
	t.CueFilterHash = r.U32Unsafe()
	t.DestTransitionTime = r.I32Unsafe()
	t.DestFadeCurve = r.U32Unsafe()
	t.DestFadeOffset = r.I32Unsafe()
}

func ParseMusicRanSeqCntr(size uint32, r *wio.Reader) *wwise.MusicRanSeqCntr {
	assert.Equal(0, r.Pos(), "Muisc random / sequence container parser position doesn't start at position 0.")
	begin := r.Pos()

	m := wwise.MusicRanSeqCntr{}
	m.Id = r.U32Unsafe()
	m.OverwriteFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r)
	m.Children = *ParseContainer(r)
	ParseMeterInfo(r, &m.MeterInfo)
	m.Stingers = make([]wwise.Stinger, r.U32Unsafe())
	ParseStingers(r, m.Stingers)
	m.TransitionRules = make([]wwise.MusicTransitionRule, r.U32Unsafe())
	ParseTransitionRules(r, m.TransitionRules)
	totalNumPlayListNodes := r.U32Unsafe()
	ParseMusicPlayListNodes(r, &m.PlayListNode, &totalNumPlayListNodes)
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

func ParseTransitionRules(r *wio.Reader, rules []wwise.MusicTransitionRule) {
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
	r *wio.Reader, p *wwise.MusicPlayListNode, totalNumPlayListNodes *uint32,
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
		ParseMusicPlayListNodes(r, &p.PlayListLeafs[i], totalNumPlayListNodes)
	}
}

func ParseMusicSwitchCntr(size uint32, r *wio.Reader) (m *wwise.MusicSwitchCntr) {
	assert.Equal(0, r.Pos(), "Music switch container parser position doesn't start at position 0.")
	begin := r.Pos()

	m = &wwise.MusicSwitchCntr{}
	m.Id = r.U32Unsafe()
	m.OverwriteFlags = r.U8Unsafe()
	m.BaseParam = *ParseBaseParam(r)
	m.Children = *ParseContainer(r)
	ParseMeterInfo(r, &m.MeterInfo)
	m.Stingers = make([]wwise.Stinger, r.U32Unsafe())
	ParseStingers(r, m.Stingers)
	m.TransitionRules = make([]wwise.MusicTransitionRule, r.U32Unsafe())
	ParseTransitionRules(r, m.TransitionRules)
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

func ParseRanSeqCntr(size uint32, r *wio.Reader) *wwise.RanSeqCntr {
	assert.Equal(0, r.Pos(), "Random / Sequence container parser position doesn't start at position 0.")
	begin := r.Pos()
	rs := &wwise.RanSeqCntr{
		Id:              r.U32Unsafe(),
		BaseParam:       ParseBaseParam(r),
		PlayListSetting: ParsePlayListSetting(r),
		Container:       ParseContainer(r),
		PlayListItems:   make([]*wwise.PlayListItem, r.U16Unsafe()),
	}
	for i := range rs.PlayListItems {
		rs.PlayListItems[i] = ParsePlayListItem(r)
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return rs
}

func ParsePlayListItem(r *wio.Reader) *wwise.PlayListItem {
	return &wwise.PlayListItem{
		UniquePlayID: r.U32Unsafe(),
		Weight:       r.I32Unsafe(),
	}
}

func ParsePlayListSetting(r *wio.Reader) *wwise.PlayListSetting {
	return &wwise.PlayListSetting{
		LoopCount:            r.U16Unsafe(),
		LoopModMin:           r.U16Unsafe(),
		LoopModMax:           r.U16Unsafe(),
		TransitionTime:       r.F32Unsafe(),
		TransitionTimeModMin: r.F32Unsafe(),
		TransitionTimeModMax: r.F32Unsafe(),
		AvoidRepeatCount:     r.U16Unsafe(),
		TransitionMode:       r.U8Unsafe(),
		RandomMode:           r.U8Unsafe(),
		Mode:                 r.U8Unsafe(),
		PlayListBitVector:    r.U8Unsafe(),
	}
}

func ParseSound(size uint32, r *wio.Reader) *wwise.Sound {
	assert.Equal(0, r.Pos(), "Sound parser position doesn't start 0.")
	begin := r.Pos()
	s := &wwise.Sound{}
	s.Id = r.U32Unsafe()
	s.BankSourceData = ParseBankSourceData(r)
	s.BaseParam = ParseBaseParam(r)
	end := r.Pos()
	if begin >= end {
		panic("reader consumes zero byte")
	}
	assert.Equal(uint64(size), end-begin,
		"the amount of bytes reader consume doesn't equal size in the hierarchy header",
	)
	return s
}

func ParseBankSourceData(r *wio.Reader) wwise.BankSourceData {
	b := wwise.BankSourceData{
		PluginID:          r.U32Unsafe(),
		StreamType:        r.U8Unsafe(),
		SourceID:          r.U32Unsafe(),
		InMemoryMediaSize: r.U32Unsafe(),
		SourceBits:        r.U8Unsafe(),
		PluginParam:       nil,
	}
	if !b.HasParam() {
		return b
	}
	if b.PluginID == 0 {
		return b
	}
	b.PluginParam = &wwise.PluginParam{
		PluginParamSize: r.U32Unsafe(), PluginParamData: []byte{},
	}
	if b.PluginParam.PluginParamSize <= 0 {
		return b
	}
	begin := r.Pos()
	b.PluginParam.PluginParamData = r.ReadNUnsafe(
		uint64(b.PluginParam.PluginParamSize), 0,
	)
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(b.PluginParam.PluginParamSize, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal size in "+
			"source plugin parameter header",
	)
	return b
}

func ParseSwitchCntr(size uint32, r *wio.Reader) *wwise.SwitchCntr {
	assert.Equal(0, r.Pos(), "Switch container parser position doesn't start at 0.")
	begin := r.Pos()
	s := &wwise.SwitchCntr{
		Id:                     r.U32Unsafe(),
		BaseParam:              ParseBaseParam(r),
		GroupType:              r.U8Unsafe(),
		GroupID:                r.U32Unsafe(),
		DefaultSwitch:          r.U32Unsafe(),
		IsContinuousValidation: r.U8Unsafe(),
		Container:              ParseContainer(r),
		SwitchGroups:           make([]*wwise.SwitchGroupItem, r.U32Unsafe()),
	}
	for i := range s.SwitchGroups {
		item := &wwise.SwitchGroupItem{
			SwitchID: r.U32Unsafe(),
			NodeList: make([]uint32, r.U32Unsafe()),
		}
		nodeList := item.NodeList
		for j := range nodeList {
			nodeList[j] = r.U32Unsafe()
		}
		s.SwitchGroups[i] = item
	}
	NumSwitchParam := r.U32Unsafe()
	s.SwitchParams = make([]*wwise.SwitchParam, NumSwitchParam)
	for i := range s.SwitchParams {
		s.SwitchParams[i] = &wwise.SwitchParam{
			NodeId:            r.U32Unsafe(),
			PlayBackBitVector: r.U8Unsafe(),
			ModeBitVector:     r.U8Unsafe(),
			FadeOutTime:       r.I32Unsafe(),
			FadeInTime:        r.I32Unsafe(),
		}
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return s
}

func ParseBaseParam(r *wio.Reader) *wwise.BaseParameter {
	b := wwise.BaseParameter{}
	b.BitIsOverrideParentFx = r.U8Unsafe()
	ParseFxChunk(r, &b.FxChunk)
	ParseFxChunkMetadata(r, &b.FxChunkMetadata)
	b.BitOverrideAttachmentParams = r.U8Unsafe()
	b.OverrideBusId = r.U32Unsafe()
	b.DirectParentId = r.U32Unsafe()
	b.ByBitVectorA = r.U8Unsafe()
	ParsePropBundle(r, &b.PropBundle)
	ParseRangePropBundle(r, &b.RangePropBundle)
	ParsePositioningParam(r, &b.PositioningParam)
	ParseAuxParam(r, &b.AuxParam)
	ParseAdvanceSetting(r, &b.AdvanceSetting)
	ParseStateProp(r, &b.StateProp)
	ParseStateGroup(r, &b.StateGroup)
	ParseRTPC(r, &b.RTPC)
	return &b
}

func ParseFxChunk(r *wio.Reader, f *wwise.FxChunk) {
	UniqueNumFx := r.U8Unsafe()
	if UniqueNumFx <= 0 {
		f.BitsFxByPass = 0
		f.FxChunkItems = make([]wwise.FxChunkItem, 0)
		return
	}
	f.BitsFxByPass = r.U8Unsafe()
	f.FxChunkItems = make([]wwise.FxChunkItem, UniqueNumFx)
	for i := range f.FxChunkItems {
		f.FxChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxChunkItems[i].FxId = r.U32Unsafe()
		f.FxChunkItems[i].BitIsShareSet = r.U8Unsafe()
		f.FxChunkItems[i].BitIsRendered = r.U8Unsafe()
	}
}

func ParseFxChunkMetadata(r *wio.Reader, f *wwise.FxChunkMetadata)  {
	f.BitIsOverrideParentMetadata = r.U8Unsafe()
	UniqueNumFxMetadata := r.U8Unsafe()
	f.FxMetaDataChunkItems = make([]wwise.FxChunkMetadataItem, UniqueNumFxMetadata)
	for i := range f.FxMetaDataChunkItems {
		f.FxMetaDataChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxMetaDataChunkItems[i].FxId = r.U32Unsafe()
		f.FxMetaDataChunkItems[i].BitIsShareSet = r.U8Unsafe()
	}
}

func ParsePropBundle(r *wio.Reader, p *wwise.PropBundle) {
	CProps := r.U8Unsafe()
	p.PropValues = make([]wwise.PropValue, CProps)
	for i := range CProps {
		p.PropValues[i].P = wwise.PropType(r.U8Unsafe())
	}
	for i := range CProps {
		p.PropValues[i].V = r.ReadNUnsafe(4, 0)
	}
}

func ParseRangePropBundle(r *wio.Reader, p *wwise.RangePropBundle) {
	CProps := r.U8Unsafe()
	p.RangeValues = make([]wwise.RangeValue, CProps)
	for i := range p.RangeValues {
		p.RangeValues[i].P = wwise.PropType(r.U8Unsafe())
	}
	for i := range p.RangeValues {
		p.RangeValues[i].Min = r.ReadNUnsafe(4, 0)
		p.RangeValues[i].Max = r.ReadNUnsafe(4, 0)
	}
}

func ParsePositioningParam(r *wio.Reader, p *wwise.PositioningParam) {
	p.BitsPositioning = r.U8Unsafe()
	if !p.HasPositioningAnd3D() {
		return
	}
	p.Bits3D = r.U8Unsafe()
	if !p.HasAutomation() {
		return
	}
	p.PathMode = r.U8Unsafe()
	p.TransitionTime = r.I32Unsafe()
	NumPositionVertices := r.U32Unsafe()
	p.PositionVertices = make([]wwise.PositionVertex, NumPositionVertices)
	for i := range p.PositionVertices {
		p.PositionVertices[i].X = r.F32Unsafe()
		p.PositionVertices[i].Y = r.F32Unsafe()
		p.PositionVertices[i].Z = r.F32Unsafe()
		p.PositionVertices[i].Duration = r.I32Unsafe()
	}
	NumPositionPlayListItem := r.U32Unsafe()
	p.PositionPlayListItems = make([]wwise.PositionPlayListItem, NumPositionPlayListItem)
	for i := range p.PositionPlayListItems {
		p.PositionPlayListItems[i].UniqueVerticesOffset = r.U32Unsafe()
		p.PositionPlayListItems[i].INumVertices =  r.U32Unsafe()
	}
	p.Ak3DAutomationParams = make([]wwise.Ak3DAutomationParam, NumPositionPlayListItem)
	for i := range p.Ak3DAutomationParams {
		p.Ak3DAutomationParams[i].XRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].YRange = r.F32Unsafe()
		p.Ak3DAutomationParams[i].ZRange = r.F32Unsafe()
	}
}

func ParseAuxParam(r *wio.Reader, a *wwise.AuxParam) {
	a.AuxBitVector = r.U8Unsafe()
	if a.HasAux() {
		a.AuxIds[0] = r.U32Unsafe()
		a.AuxIds[1] = r.U32Unsafe()
		a.AuxIds[2] = r.U32Unsafe()
		a.AuxIds[3] = r.U32Unsafe()
		a.RestoreAuxIds[0] = a.AuxIds[0]
		a.RestoreAuxIds[1] = a.AuxIds[1]
		a.RestoreAuxIds[2] = a.AuxIds[2]
		a.RestoreAuxIds[3] = a.AuxIds[3]
	}
	a.ReflectionAuxBus = r.U32Unsafe()
	a.RestoreReflectionAuxBus = a.ReflectionAuxBus
}

func ParseAdvanceSetting(r *wio.Reader, a *wwise.AdvanceSetting) {
	a.AdvanceSettingBitVector = r.U8Unsafe()
	a.VirtualQueueBehavior = r.U8Unsafe()
	a.MaxNumInstance = r.U16Unsafe()
	a.BelowThresholdBehavior =  r.U8Unsafe()
	a.HDRBitVector = r.U8Unsafe()
}

func ParseStateProp(r *wio.Reader, s *wwise.StateProp) {
	NumStateProps := r.U8Unsafe()
	s.StatePropItems = make([]wwise.StatePropItem, NumStateProps)
	for i := range s.StatePropItems {
		s.StatePropItems[i].PropertyId = r.U8Unsafe()
		s.StatePropItems[i].AccumType = r.U8Unsafe()
		s.StatePropItems[i].InDb = r.U8Unsafe()
	}
}

func ParseStateGroup(r *wio.Reader, s *wwise.StateGroup) {
	NumStateGroups := r.U8Unsafe()
	s.StateGroupItems = make([]wwise.StateGroupItem, NumStateGroups)
	for i := range s.StateGroupItems {
		item := &s.StateGroupItems[i]
		item.StateGroupID = r.U32Unsafe()
		item.StateSyncType = r.U8Unsafe()
		NumStates := r.U8Unsafe()
		item.States = make([]wwise.StateGroupItemState, NumStates)
		for i := range item.States {
			item.States[i].StateID = r.U32Unsafe()
			item.States[i].StateInstanceID = r.U32Unsafe()
		}
	}
}

func ParseRTPC(r *wio.Reader, rtpc *wwise.RTPC) {
	NumRTPC := r.U16Unsafe()
	rtpc.RTPCItems = make([]wwise.RTPCItem, NumRTPC, NumRTPC)
	for i := range rtpc.RTPCItems {
		item := &rtpc.RTPCItems[i]
		item.RTPCID = r.U32Unsafe()
		item.RTPCType = r.U8Unsafe()
		item.RTPCAccum = r.U8Unsafe()
		item.ParamID = r.U8Unsafe()
		item.RTPCCurveID = r.U32Unsafe()
		item.Scaling = r.U8Unsafe()
		NumRTPCGraphPoints := r.U16Unsafe()
		item.RTPCGraphPoints = make([]wwise.RTPCGraphPoint, NumRTPCGraphPoints, NumRTPCGraphPoints)
		ParseRTPCGraphPoints(r, item.RTPCGraphPoints)
	}
}

func ParseContainer(r *wio.Reader) *wwise.Container {
	c := wwise.NewCntrChildren()
	NumChild := r.U32Unsafe()
	c.Children = make([]uint32, NumChild)
	for i := range c.Children {
		c.Children[i] = r.U32Unsafe()
	}
	return c
}

func ParseAttenuation(size uint32, r *wio.Reader) *wwise.Attenuation {
	assert.Equal(0, r.Pos(), "Attenuation parser position doesn't start at position 0.")
	begin := r.Pos()

	a := wwise.Attenuation{
		Id: r.U32Unsafe(),
		IsHeightSpreadEnabled: r.U8Unsafe(),
		IsConeEnabled: r.U8Unsafe(),
		Curves: [7]int8(make([]int8, 7)),
	}

	if a.ConeEnabled() {
		a.InsideDegrees = r.F32Unsafe()
		a.OutsideDegrees = r.F32Unsafe()
		a.OutsideVolume = r.F32Unsafe()
		a.LoPass = r.F32Unsafe()
		a.HiPass = r.F32Unsafe()
	}

	for i := range a.Curves {
		a.Curves[i] = r.I8Unsafe()
	}

	a.AttenuationConversionTables = make([]wwise.AttenuationConversionTable, r.U8Unsafe())
	ParseAttenuationConversionTables(r, a.AttenuationConversionTables)
	ParseRTPC(r, &a.RTPC)

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &a
}

func ParseAttenuationConversionTables(r *wio.Reader, t []wwise.AttenuationConversionTable) {
	for i := range t {
		t[i].EnumScaling = r.U8Unsafe()
		t[i].RTPCGraphPoints = make([]wwise.RTPCGraphPoint, r.U16Unsafe())
		ParseRTPCGraphPoints(r, t[i].RTPCGraphPoints)
	}
}

func ParseRTPCGraphPoints(r *wio.Reader, pts []wwise.RTPCGraphPoint) {
	for i := range pts {
		pts[i].From = r.F32Unsafe()
		pts[i].To = r.F32Unsafe()
		pts[i].Interp = r.U32Unsafe()
	}
}
