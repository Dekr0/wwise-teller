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
	"slices"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/interp"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const maxNumParseRoutine = 8

type parserResult struct {
	i uint32
	obj wwise.HircObj
}

type parser func(uint32, *wio.Reader) wwise.HircObj

func parseHIRC(ctx context.Context, size uint32, r *wio.Reader) (*wwise.HIRC, error) {
	assert.Equal(0, r.Pos(), "Parser for HIRC does not start at byte 0.")

	numHircItem := r.U32Unsafe()

	hirc := wwise.NewHIRC(size, numHircItem)

	/* sync signal */
	sem := make(chan struct{}, maxNumParseRoutine)
	parserResult := make(chan *parserResult, numHircItem)
	i := 0
	parsed := 0

	slog.Debug("Start scanning through all hierarchy object, and scheduling parser",
		"numHircItem", numHircItem,
	)

	for parsed < int(numHircItem) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res := <- parserResult:
			addHircObj(hirc, res.i, res.obj)
			parsed += 1
		default:
		}
		if i >= int(numHircItem) {
			continue
		}
		eHircType := r.U8Unsafe()
		dwSectionSize := r.U32Unsafe()
		if skipHircObjType(wwise.HircType(eHircType)) {
			unknown := wwise.NewUnknown(
				wwise.HircType(eHircType),
				dwSectionSize,
				r.ReadNUnsafe(uint64(dwSectionSize), 4),
			)
			hirc.HircObjs[i] = unknown
			i += 1
			parsed += 1
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
			case wwise.HircTypeSound:
				go parserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					parseSound,
					parserResult,
					sem,
				)
			case wwise.HircRanSeqCntr:
				go parserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					parseRanSeqCntr,
					parserResult,
					sem,
				)
			case wwise.HircSwitchCntr:
				go parserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					parseSwitchCntr,
					parserResult,
					sem,
				)
			case wwise.HircTypeActorMixer:
				go parserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					parseActorMixer,
					parserResult,
					sem,
				)
			case wwise.HircTypeLayerCntr:
				go parserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					parseLayerCntr,
					parserResult,
					sem,
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
			case wwise.HircTypeSound:
				obj = parseSound(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircRanSeqCntr:
				obj = parseRanSeqCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircSwitchCntr:
				obj = parseSwitchCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeActorMixer:
				obj = parseActorMixer(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeLayerCntr:
				obj = parseLayerCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			default:
				panic("Assertion Trap")
			}
			addHircObj(hirc, uint32(i), obj)
			i += 1
			parsed += 1
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
func addHircObj(h *wwise.HIRC, i uint32, obj wwise.HircObj) {
	t := obj.HircType()
	id, err := obj.HircID()
	if err != nil { panic(err) }
	switch t {
	case wwise.HircTypeSound:
		if _, in := h.Sounds[id]; in {
			panic(fmt.Sprintf("Duplicate sound object %d", id))
		}
		h.Sounds[id] = obj.(*wwise.Sound)
	case wwise.HircRanSeqCntr:
		if _, in := h.RanSeqCntrs[id]; in {
			panic(fmt.Sprintf("Duplicate random / sequence container object %d", id))
		}
		h.RanSeqCntrs[id] = obj.(*wwise.RanSeqCntr)
	case wwise.HircSwitchCntr:
		if _, in := h.SwitchCntrs[id]; in {
			panic(fmt.Sprintf("Duplicate switch container object %d", id))
		}
		h.SwitchCntrs[id] = obj.(*wwise.SwitchCntr)
	case wwise.HircTypeActorMixer:
		if _, in := h.ActorMixers[id]; in {
			panic(fmt.Sprintf("Duplicate actor mixer object %d", id))
		}
		h.ActorMixers[id] = obj.(*wwise.ActorMixer)
	case wwise.HircTypeLayerCntr:
		if _, in := h.LayerCntrs[id]; in {
			panic(fmt.Sprintf("Duplicate layer container object %d", id))
		}
		h.LayerCntrs[id] = obj.(*wwise.LayerCntr)
	default:
		panic("Assertion Trap")
	}
	h.HircObjs[i] = obj
	slog.Debug(fmt.Sprintf("Collected %s parser", wwise.HircTypeName[obj.HircType()]))
}

func skipHircObjType(t wwise.HircType) bool {
	_, find := sort.Find(len(wwise.KnownHircType), func(i int) int {
		if t < wwise.KnownHircType[i] {
			return -1
		}
		if t == wwise.KnownHircType[i] {
			return 0
		}
		return 1
	})
	return !find
}

// Side effect: Channel will receive a parser result. Semaphore channel will 
// release an item once parser finishes to allows other more parser routine to 
// start.
func parserRoutine[T wwise.HircObj](
	size uint32,
	i uint32,
	r *wio.Reader,
	f func(uint32, *wio.Reader) T,
	c chan *parserResult,
	sem chan struct{},
) {
	hircObj := f(size, r)
	c <- &parserResult{i, hircObj}
	<- sem
}

func parseActorMixer(size uint32, r *wio.Reader) *wwise.ActorMixer {
	assert.Equal(0, r.Pos(), "Actor mixer parser position doesn't start at position 0.")
	begin := r.Pos()
	a := &wwise.ActorMixer{}
	a.Id = r.U32Unsafe()
	a.BaseParam = parseBaseParam(r)
	a.Container = parseContainer(r)
	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return a
}

func parseLayerCntr(size uint32, r *wio.Reader) *wwise.LayerCntr {
	assert.Equal(0, r.Pos(), "Layer container parser position doesn't start at position 0.")
	begin := r.Pos()
	l := &wwise.LayerCntr{
		Id: r.U32Unsafe(),
		BaseParam: parseBaseParam(r),
		Container: parseContainer(r),
		Layers: parseLayers(r),
		IsContinuousValidation: r.U8Unsafe(),
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return l
}

func parseLayers(r *wio.Reader) []*wwise.Layer {
	layers := make([]*wwise.Layer, r.U32Unsafe())
	for i := range layers {
		l := &wwise.Layer{
			Id: r.U32Unsafe(),
			InitialRTPC: parseRTPC(r),
			RTPCId: r.U32Unsafe(),
			RTPCType: r.U8Unsafe(),
			LayerRTPCs: make([]*wwise.LayerRTPC, r.U32Unsafe()),
		}
		for j := range l.LayerRTPCs {
			lr := &wwise.LayerRTPC{
				AssociatedChildID: r.U32Unsafe(),
				RTPCGraphPoints: make([]*wwise.RTPCGraphPoint, r.U32Unsafe()),
			}
			for k := range lr.RTPCGraphPoints {
				lr.RTPCGraphPoints[k] = &wwise.RTPCGraphPoint{
					From: r.F32Unsafe(),
					To: r.F32Unsafe(),
					Interp: r.U32Unsafe(),
				}
			}
			l.LayerRTPCs[j] = lr
		}
		layers[i] = l
	}
	return layers
}

func parseRanSeqCntr(size uint32, r *wio.Reader) *wwise.RanSeqCntr {
	assert.Equal(0, r.Pos(), "Random / Sequence container parser position doesn't start at position 0.")
	begin := r.Pos()
	rs := &wwise.RanSeqCntr{
		Id: r.U32Unsafe(),
		BaseParam: parseBaseParam(r),
		PlayListSetting: parsePlayListSetting(r),
		Container: parseContainer(r),
		PlayListItems: make([]*wwise.PlayListItem, r.U16Unsafe()),
	}
	for i := range rs.PlayListItems {
		rs.PlayListItems[i] = parsePlayListItem(r)
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return rs
}

func parsePlayListItem(r *wio.Reader) *wwise.PlayListItem {
	return &wwise.PlayListItem{
		UniquePlayID: r.U32Unsafe(),
		Weight: r.I32Unsafe(),
	}
}

func parsePlayListSetting(r *wio.Reader) *wwise.PlayListSetting {
	return &wwise.PlayListSetting{
		LoopCount: r.U16Unsafe(),
		LoopModMin: r.U16Unsafe(),
		LoopModMax: r.U16Unsafe(),
		TransitionTime: r.F32Unsafe(),
		TransitionTimeModMin: r.F32Unsafe(),
		TransitionTimeModMax: r.F32Unsafe(),
		AvoidRepeatCount: r.U16Unsafe(),
		TransitionMode: r.U8Unsafe(),
		RandomMode: r.U8Unsafe(),
		Mode: r.U8Unsafe(),
		PlayListBitVector: r.U8Unsafe(),
	}
}

func parseSound(size uint32, r *wio.Reader) *wwise.Sound {
	assert.Equal(0, r.Pos(), "Sound parser position doesn't start 0.")
	begin := r.Pos()
	s := &wwise.Sound{}
	s.Id = r.U32Unsafe()
	s.BankSourceData = parseBankSourceData(r)
	s.BaseParam = parseBaseParam(r)
	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte")
	}
	assert.Equal(uint64(size), end - begin,
		"The amount of bytes reader consume doesn't equal size in the hierarchy header",
	)
	return s
}

func parseBankSourceData(r *wio.Reader) *wwise.BankSourceData {
	b := &wwise.BankSourceData{
		PluginID: r.U32Unsafe(),
		StreamType: r.U8Unsafe(),
		SourceID: r.U32Unsafe(),
		InMemoryMediaSize: r.U32Unsafe(),
		SourceBits: r.U8Unsafe(),
		PluginParam: nil,
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
	assert.Equal(b.PluginParam.PluginParamSize, uint32(end - begin),
		"The amount of bytes reader consume doesn't equal size in " + 
		"source plugin parameter header",
	)
	return b
}

func parseSwitchCntr(size uint32, r *wio.Reader) *wwise.SwitchCntr {
	assert.Equal(0, r.Pos(), "Switch container parser position doesn't start at 0.")
	begin := r.Pos()
	s := &wwise.SwitchCntr{
		Id: r.U32Unsafe(),
		BaseParam: parseBaseParam(r),
		GroupType: r.U8Unsafe(),
		GroupID: r.U32Unsafe(),
		DefaultSwitch: r.U32Unsafe(),
		IsContinuousValidation: r.U8Unsafe(),
		Container: parseContainer(r),
		SwitchGroups: make([]*wwise.SwitchGroupItem, r.U32Unsafe()),
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
			NodeId: r.U32Unsafe(),
			PlayBackBitVector: r.U8Unsafe(),
			ModeBitVector: r.U8Unsafe(),
			FadeOutTime: r.I32Unsafe(),
			FadeInTime: r.I32Unsafe(),
		}
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end - begin), 
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return s
}

func parseBaseParam(r *wio.Reader) *wwise.BaseParameter {
	b := wwise.BaseParameter{}
	b.BitIsOverrideParentFx = r.U8Unsafe()
	b.FxChunk = parseFxChunk(r)
	b.FxChunkMetadata = parseFxChunkMetadata(r)
	b.BitOverrideAttachmentParams = r.U8Unsafe()
	b.OverrideBusId = r.U32Unsafe()
	b.DirectParentId = r.U32Unsafe()
	b.ByBitVectorA = r.U8Unsafe()
	b.PropBundle = parsePropBundle(r)
	b.RangePropBundle = parseRangePropBundle(r)
	b.PositioningParam = parsePositioningParam(r)
	b.AuxParam = parseAuxParam(r)
	b.AdvanceSetting = parseAdvanceSetting(r)
	b.StateProp = parseStateProp(r)
	b.StateGroup = parseStateGroup(r)
	b.RTPC = parseRTPC(r)
	return &b
}

func parseFxChunk(r *wio.Reader) *wwise.FxChunk {
	f := wwise.NewFxChunk()
	UniqueNumFx := r.U8Unsafe()
	if UniqueNumFx <= 0 {
		f.BitsFxByPass = 0
		f.FxChunkItems = make([]*wwise.FxChunkItem, 0)
		return f
	}
	f.BitsFxByPass = r.U8Unsafe()
	f.FxChunkItems = make([]*wwise.FxChunkItem, UniqueNumFx)
	for i := range f.FxChunkItems {
		fxChunkItem := &wwise.FxChunkItem{
			UniqueFxIndex: r.U8Unsafe(),
			FxId:          r.U32Unsafe(),
			BitIsShareSet: r.U8Unsafe(),
			BitIsRendered: r.U8Unsafe(),
		}
		f.FxChunkItems[i] = fxChunkItem
	}
	return f
}

func parseFxChunkMetadata(r *wio.Reader) *wwise.FxChunkMetadata {
	f := wwise.NewFxChunkMetadata()
	f.BitIsOverrideParentMetadata = r.U8Unsafe()
	UniqueNumFxMetadata := r.U8Unsafe()
	f.FxMetaDataChunkItems = make([]*wwise.FxChunkMetadataItem, UniqueNumFxMetadata)
	for i := range f.FxMetaDataChunkItems {
		f.FxMetaDataChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxMetaDataChunkItems[i].FxId = r.U32Unsafe()
		f.FxMetaDataChunkItems[i].BitIsShareSet = r.U8Unsafe()
	}
	return f
}

func parsePropBundle(r *wio.Reader) *wwise.PropBundle {
	p := wwise.NewPropBundle()
	CProps := r.U8Unsafe()
	p.PropValues = make([]*wwise.PropValue, CProps)
	for i := range CProps {
		p.PropValues[i] = &wwise.PropValue{P:r.U8Unsafe()}
	}
	for i := range CProps {
		p.PropValues[i].V = r.ReadNUnsafe(4, 0)
	}
	return p
}

func parseRangePropBundle(r *wio.Reader) *wwise.RangePropBundle {
	p := wwise.NewRangePropBundle()
	CProps := r.U8Unsafe()
	p.RangeValues = make([]*wwise.RangeValue, CProps)
	for i := range p.RangeValues {
		p.RangeValues[i] = &wwise.RangeValue{PId:r.U8Unsafe()}
	}
	for i := range p.RangeValues {
		p.RangeValues[i].Min = r.ReadNUnsafe(4, 0)
		p.RangeValues[i].Max = r.ReadNUnsafe(4, 0)
	}
	return p
}

func parsePositioningParam(r *wio.Reader) *wwise.PositioningParam {
	p := wwise.NewPositioningParam()
	p.BitsPositioning = r.U8Unsafe()
	if !p.HasPositioningAnd3D() {
		return p
	}
	p.Bits3D = r.U8Unsafe()
	if !p.HasAutomation() {
		return p
	}
	p.PathMode = r.U8Unsafe()
	p.TransitionTime = r.I32Unsafe()
	NumPositionVertices := r.U32Unsafe()
	p.PositionVertices = make([]*wwise.PositionVertex, NumPositionVertices)
	for i := range p.PositionVertices {
		p.PositionVertices[i] = &wwise.PositionVertex{
			X:        r.F32Unsafe(),
			Y:        r.F32Unsafe(),
			Z:        r.F32Unsafe(),
			Duration: r.I32Unsafe(),
		}
	}
	NumPositionPlayListItem := r.U32Unsafe()
	p.PositionPlayListItems = make([]*wwise.PositionPlayListItem, NumPositionPlayListItem)
	for i := range p.PositionPlayListItems {
		p.PositionPlayListItems[i] = &wwise.PositionPlayListItem{
			UniqueVerticesOffset: r.U32Unsafe(),
			INumVertices:         r.U32Unsafe(),
		}
	}
	p.Ak3DAutomationParams = make([]*wwise.Ak3DAutomationParam, NumPositionPlayListItem)
	for i := range p.Ak3DAutomationParams {
		p.Ak3DAutomationParams[i] = &wwise.Ak3DAutomationParam{
			XRange: r.F32Unsafe(),
			YRange: r.F32Unsafe(),
			ZRange: r.F32Unsafe(),
		}
	}
	return p
}

func parseAuxParam(r *wio.Reader) *wwise.AuxParam {
	a := wwise.NewAuxParam()
	a.AuxBitVector = r.U8Unsafe()
	if a.HasAux() {
		a.AuxIds = make([]uint32, 4, 4)
		a.AuxIds[0] = r.U32Unsafe()
		a.AuxIds[1] = r.U32Unsafe()
		a.AuxIds[2] = r.U32Unsafe()
		a.AuxIds[3] = r.U32Unsafe()
		a.RestoreAuxIds = slices.Clone(a.AuxIds)
	}
	a.ReflectionAuxBus = r.U32Unsafe()
	return a
}

func parseAdvanceSetting(r *wio.Reader) *wwise.AdvanceSetting {
	return &wwise.AdvanceSetting{
		AdvanceSettingBitVector: r.U8Unsafe(),
		VirtualQueueBehavior:    r.U8Unsafe(),
		MaxNumInstance:          r.U16Unsafe(),
		BelowThresholdBehavior:  r.U8Unsafe(),
		HDRBitVector:            r.U8Unsafe(),
	}
}

func parseStateProp(r *wio.Reader) *wwise.StateProp {
	s := wwise.NewStateProp()
	NumStateProps := r.U8Unsafe()
	s.StatePropItems = make([]*wwise.StatePropItem, NumStateProps)
	for i := range s.StatePropItems {
		s.StatePropItems[i] = &wwise.StatePropItem{
			PropertyId: r.U8Unsafe(),
			AccumType:  r.U8Unsafe(),
			InDb:       r.U8Unsafe(),
		}
	}
	return s
}

func parseStateGroup(r *wio.Reader) *wwise.StateGroup {
	s := wwise.NewStateGroup()
	NumStateGroups := r.U8Unsafe()
	s.StateGroupItems = make([]*wwise.StateGroupItem, NumStateGroups)
	for i := range s.StateGroupItems {
		item := wwise.NewStateGroupItem()
		item.StateGroupID = r.U32Unsafe()
		item.StateSyncType = r.U8Unsafe()
		NumStates := r.U8Unsafe()
		item.States = make([]*wwise.StateGroupItemState, NumStates)
		for i := range item.States {
			item.States[i] = &wwise.StateGroupItemState{
				StateID:         r.U32Unsafe(),
				StateInstanceID: r.U32Unsafe(),
			}
		}
		s.StateGroupItems[i] = item
	}
	return s
}

func parseRTPC(r *wio.Reader) *wwise.RTPC {
	rtpc := wwise.NewRTPC()
	NumRTPC := r.U16Unsafe()
	rtpc.RTPCItems = make([]*wwise.RTPCItem, NumRTPC, NumRTPC)
	for i := range rtpc.RTPCItems {
		item := wwise.NewRTPCItem()
		item.RTPCID = r.U32Unsafe()
		item.RTPCType = r.U8Unsafe()
		item.RTPCAccum = r.U8Unsafe()
		item.ParamID = r.U8Unsafe()
		item.RTPCCurveID = r.U32Unsafe()
		item.Scaling = r.U8Unsafe()
		NumRTPCGraphPoints := r.U16Unsafe()
		item.RTPCGraphPoints = make([]*wwise.RTPCGraphPoint, NumRTPCGraphPoints, NumRTPCGraphPoints)
		for j := range NumRTPCGraphPoints {
			r := &wwise.RTPCGraphPoint{
				From: r.F32Unsafe(), To: r.F32Unsafe(), Interp: r.U32Unsafe(),
			}
			item.RTPCGraphPoints[j] = r
		}
		// item.SamplePoints = computeRTPCSamplePoints(item.RTPCGraphPoints)
		rtpc.RTPCItems[i] = item
	}
	return rtpc
}

func computeRTPCSamplePoints(rpts []*wwise.RTPCGraphPoint) []float32 {
	spts := make([]float32, 0, len(rpts) * interp.NumSamples)
	for i, p1 := range rpts {
		if i >= len(rpts) - 1 {
			break
		}
		p2 := rpts[i + 1]
		switch p1.Interp {
		case 1:
			spts = append(spts, interp.SampleLog3(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 2:
			spts = append(spts, interp.SampleSine(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 3:
			spts = append(spts, interp.SampleLog1(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 4:
			spts = append(spts, interp.SampleInvertSCurve(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 5:
			spts = append(spts, interp.SampleLinear(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 6:
			spts = append(spts, interp.SampleSCurve(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 7:
			spts = append(spts, interp.SampleExp1(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 8:
			spts = append(spts, interp.SampleConst(p2.To)...)
		case 9:
			spts = append(spts, interp.SampleExp3(float64(p1.From), float64(p1.To), float64(p2.From), float64(p2.To))...)
		case 10:
			spts = append(spts, interp.SampleConst(p2.To)...)
		}
	}
	return spts
}

func parseContainer(r *wio.Reader) *wwise.Container {
	c := wwise.NewCntrChildren()
	NumChild := r.U32Unsafe()
	c.Children = make([]uint32, NumChild)
	for i := range c.Children {
		c.Children[i] = r.U32Unsafe()
	}
	return c
}
