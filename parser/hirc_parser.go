package parser

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
	"github.com/Dekr0/wwise-teller/wwise"
)

const MAX_NUM_ROUTINE = 8

type parseResult struct {
	index   uint32
	hircObj wwise.HircObj
}

func parseHIRC(
	ctx context.Context,
	chunkSize uint32,
	r *reader.BankReader,
) (*wwise.HIRC, error) {
	assert.AssertEqual(0, r.Tell(), "Parser for HIRC does not start at byte 0.")

	numHircItem, err := r.U32()
	if err != nil {
		return nil, err
	}

	hirc := wwise.NewHIRC(chunkSize, numHircItem)

	cErr := make(chan error, numHircItem)
	sem := make(chan struct{}, numHircItem)
	cParseResult := make(chan *parseResult, numHircItem)
	activeParser := 0

	slog.Info("Start scanning through all hierarchy object, and scheduling parser",
		"numHircItem", numHircItem,
	)

	for i := range numHircItem {
		eHircType, err := r.U8()
		if err != nil {
			return nil, err
		}
		dwSectionSize, err := r.U32()
		if err != nil {
			return nil, err
		}
		switch eHircType {
		/*
			case 0x02:
				initErr := dispatchParserRoutine(
					dwSectionSize, r, cHirc, cErr, sem, parseSound,
				)
				if initErr != nil {
					return nil, err
				}
			case 0x05:
				initErr := dispatchParserRoutine(
					dwSectionSize, r, cHirc, cErr, sem, parseRanSeqCntr,
				)
				if initErr != nil {
					return nil, err
				}
			case 0x06:
				initErr := dispatchParserRoutine(
					dwSectionSize, r, cHirc, cErr, sem, parseSwitchCntr,
				)
				if initErr != nil {
					return nil, err
				}
		*/
		case 0x07:
			slog.Info(
				"Scheduling actor mixer parser...", 
				"index", i, 
				"dwSectionSize", dwSectionSize,
				"consumeSize", r.Tell(),
			)
			initErr := dispatchParserRoutine(
				dwSectionSize, i, r, cParseResult, cErr, sem, parseActorMixer,
			)
			if initErr != nil {
				return nil, err
			}
			activeParser += 1
		/*
			case 0x09:
				initErr := dispatchParserRoutine(
					dwSectionSize, r, cHirc, cErr, sem, parseLayerCntr,
				)
				if initErr != nil {
					return nil, err
				}
		*/
		default:
			blob, err := r.ReadFull(uint64(dwSectionSize), 4)
			if err != nil {
				return nil, err
			}
			unknown := &wwise.Unknown{
				Header: &wwise.HircObjHeader{
					HircType: eHircType,
					HircSize: dwSectionSize,
				},
				Blob: blob,
			}
			hirc.HircObjs[i] = unknown
			slog.Info(
				"Skip parsing hierarchy object", 
				"index", i,
				"eHircType", eHircType,
				"dwSectionSize", dwSectionSize,
				"consumeSize", r.Tell(),
			)
		}
	}

	slog.Info("Scanned hierarchy objects, and scheduled parser",
		"numHircItem", numHircItem,
		"chunkSize", chunkSize,
		"consumeSize", r.Tell(),
	)

	assert.AssertEqual(
		chunkSize,
		uint32(r.Tell()),
		"There are data that is not consumed after parsing all HIRC blob",
	)

	for activeParser > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case parseResult := <-cParseResult:
			slog.Info("Collected parser result", "eHircType", parseResult.hircObj.GetHircType())
			hirc.HircObjs[parseResult.index] = parseResult.hircObj
			switch parseResult.hircObj.GetHircType() {
			case 0x07:
				hircId, _ := parseResult.hircObj.GetHircID()
				if _, in := hirc.ActorMixers[hircId]; in {
					panic(fmt.Sprintf("Duplicate actor mixer object %d", hircId))
				}
				hirc.ActorMixers[hircId] = parseResult.hircObj.(*wwise.ActorMixer)
				slog.Info("Sorted parser result to actor mixer", "hircId", hircId)
				activeParser 
			}
		default: 
			slog.Info("Sleep for 5 seconds")
			time.Sleep(time.Second * 5) 
		}
	}

	return hirc, nil
}

func dispatchParserRoutine[T wwise.HircObj](
	dwSectionSize uint32,
	index uint32,
	parentReader *reader.BankReader,
	cParseResult chan *parseResult,
	cErr chan error,
	sem chan struct{},
	parserFunc func(uint32, *reader.BankReader) T,
) error {
	sectionReader, readErr := parentReader.NewSectionBankReader(uint64(dwSectionSize))
	if readErr != nil {
		return readErr
	}
	sem <- struct{}{}
	go func() {
		defer func() { <-sem }()
		hircObj := parserFunc(dwSectionSize, sectionReader)
		cParseResult <- &parseResult{index, hircObj}
	}()
	return nil
}

func parseActorMixer(dwSectionSize uint32, p *reader.BankReader) *wwise.ActorMixer {
	begin := p.Tell()
	actorMixer := &wwise.ActorMixer{}
	actorMixer.HircId = p.U32Unsafe()
	actorMixer.BaseParam = parseBaseParam(p)
	actorMixer.Children = parseChildren(p)
	end := p.Tell()
	if begin >= end {
		panic("Reader does not consume any byte at all!")
	}
	assert.AssertEqual(
		uint64(dwSectionSize),
		end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return actorMixer
}

func parseLayerCntr(dwSectionSize uint32, p *reader.BankReader) (*wwise.LayerCntr, error) {
	return nil, nil
}

func parseRanSeqCntr(dwSectionSize uint32, p *reader.BankReader) (*wwise.RanSeqCntr, error) {
	return nil, nil
}

func parseSound(dwSectionSize uint32, p *reader.BankReader) (*wwise.Sound, error) {
	return nil, nil
}

func parseSwitchCntr(dwSectionSize uint32, p *reader.BankReader) (*wwise.SwitchCntr, error) {
	return nil, nil
}

func parseBaseParam(r *reader.BankReader) *wwise.BaseParameter {
	bp := wwise.BaseParameter{}
	bp.BitIsOverrideParentFx = r.U8Unsafe()
	bp.FxChunk = parseFxChunk(r)
	bp.FxChunkMetadata = parseFxChunkMetadata(r)
	bp.BitOverrideAttachmentParams = r.U8Unsafe()
	bp.OverrideBusId = r.U32Unsafe()
	bp.DirectParentId = r.U32Unsafe()
	bp.ByBitVectorA = r.U8Unsafe()
	bp.PropBundle = parsePropBundle(r)
	bp.RangePropBundle = parseRangePropBundle(r)
	bp.PositioningParam = parsePositioningParam(r)
	bp.AuxParam = parseAuxParam(r)
	bp.AdvanceSetting = parseAdvanceSetting(r)
	bp.StateProp = parseStateProp(r)
	bp.StateGroup = parseStateGroup(r)
	bp.RTPC = parseRTPC(r)
	return &bp
}

func parseFxChunk(r *reader.BankReader) *wwise.FxChunk {
	f := wwise.NewFxChunk()
	f.UniqueNumFx = r.U8Unsafe()
	if f.UniqueNumFx <= 0 {
		f.BitsFxByPass = 0
		f.FxChunkItems = make([]*wwise.FxChunkItem, 0)
		return f
	}
	f.BitsFxByPass = r.U8Unsafe()
	f.FxChunkItems = make([]*wwise.FxChunkItem, f.UniqueNumFx)
	for i := range f.FxChunkItems {
		fxChunkItem := &wwise.FxChunkItem{
			UniqueFxIndex: r.U8Unsafe(),
			FxId: r.U32Unsafe(),
			BitIsShareSet: r.U8Unsafe(),
			BitIsRendered: r.U8Unsafe(),
		}
		f.FxChunkItems[i] = fxChunkItem
	}
	return f
}

func parseFxChunkMetadata(r *reader.BankReader) *wwise.FxChunkMetadata {
	f := wwise.NewFxChunkMetadata()
	f.BitIsOverrideParentMetadata = r.U8Unsafe()
	f.UniqueNumFxMetadata = r.U8Unsafe()
	if f.UniqueNumFxMetadata <= 0 {
		f.FxMetaDataChunkItems = make([]*wwise.FxChunkMetadataItem, 0)
		return f
	}
	f.FxMetaDataChunkItems = make([]*wwise.FxChunkMetadataItem, f.UniqueNumFxMetadata)
	for i := range f.FxMetaDataChunkItems {
		f.FxMetaDataChunkItems[i].UniqueFxIndex = r.U8Unsafe()
		f.FxMetaDataChunkItems[i].FxId = r.U32Unsafe()
		f.FxMetaDataChunkItems[i].BitIsShareSet = r.U8Unsafe()
	}
	return f
}

func parsePropBundle(r *reader.BankReader) *wwise.PropBundle {
	p := wwise.NewPropBundle()
	p.CProps = r.U8Unsafe()
	p.PIds = make([]uint8, p.CProps)
	for i := range p.PIds {
		p.PIds[i] = r.U8Unsafe()
	}
	p.PValues = make([][]byte, p.CProps)
	for i := range p.PValues {
		p.PValues[i] = r.ReadFullUnsafe(4, 0)
	}
	return p
}

func parseRangePropBundle(r *reader.BankReader) *wwise.RangePropBundle {
	p := wwise.NewRangePropBundle()
	p.CProps = r.U8Unsafe()
	p.PIds = make([]uint8, p.CProps)
	for i := range p.PIds {
		p.PIds[i] = r.U8Unsafe()
	}
	p.RangeValues = make([]*wwise.RangeValue, p.CProps)
	for i := range p.RangeValues {
		p.RangeValues[i] = &wwise.RangeValue{
			Min: r.ReadFullUnsafe(4, 0),
			Max: r.ReadFullUnsafe(4, 0),
		}
	}
	return p
}

func parsePositioningParam(r *reader.BankReader) *wwise.PositioningParam {
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
	p.NumPositionVertices = r.U32Unsafe()
	p.PositionVertices = make([]*wwise.PositionVertex, p.NumPositionVertices)
	for i := range p.PositionVertices {
		p.PositionVertices[i] = &wwise.PositionVertex{
			X:        r.F32Unsafe(),
			Y:        r.F32Unsafe(),
			Z:        r.F32Unsafe(),
			Duration: r.I32Unsafe(),
		}
	}
	p.NumPositionPlayListItem = r.U32Unsafe()
	p.PositionPlayListItems = make([]*wwise.PositionPlayListItem, p.NumPositionPlayListItem)
	for i := range p.PositionPlayListItems {
		p.PositionPlayListItems[i] = &wwise.PositionPlayListItem{
			UniqueVerticesOffset: r.U32Unsafe(),
			INumVertices:         r.U32Unsafe(),
		}
	}
	p.Ak3DAutomationParams = make([]*wwise.Ak3DAutomationParam, p.NumPositionPlayListItem)
	for i := range p.Ak3DAutomationParams {
		p.Ak3DAutomationParams[i] = &wwise.Ak3DAutomationParam{
			XRange: r.F32Unsafe(),
			YRange: r.F32Unsafe(),
			ZRange: r.F32Unsafe(),
		}
	}
	return p
}

func parseAuxParam(r *reader.BankReader) *wwise.AuxParam {
	a := wwise.NewAuxParam()
	a.AuxBitVector = r.U8Unsafe()
	if a.HasAux() {
		a.AuxIds = make([]uint32, 4, 4)
		a.AuxIds[0] = r.U32Unsafe()
		a.AuxIds[1] = r.U32Unsafe()
		a.AuxIds[2] = r.U32Unsafe()
		a.AuxIds[3] = r.U32Unsafe()
	}
	a.ReflectionAuxBus = r.U32Unsafe()
	return a
}

func parseAdvanceSetting(r *reader.BankReader) *wwise.AdvanceSetting {
	return &wwise.AdvanceSetting{
		AdvanceSettingBitVector: r.U8Unsafe(),
		VirtualQueueBehavior:    r.U8Unsafe(),
		MaxNumInstance:          r.U16Unsafe(),
		BelowThresholdBehavior:  r.U8Unsafe(),
		HDRBitVector:            r.U8Unsafe(),
	}
}

func parseStateProp(r *reader.BankReader) *wwise.StateProp {
	sp := wwise.NewStateProp()
	sp.NumStateProps = r.U8Unsafe()
	sp.StatePropItems = make([]*wwise.StatePropItem, sp.NumStateProps)
	for i := range sp.StatePropItems {
		sp.StatePropItems[i] = &wwise.StatePropItem{
			PropertyId: r.U8Unsafe(),
			AccumType:  r.U8Unsafe(),
			InDb:       r.U8Unsafe(),
		}
	}
	return sp
}

func parseStateGroup(r *reader.BankReader) *wwise.StateGroup {
	sg := wwise.NewStateGroup()
	sg.NumStateGroups = r.U8Unsafe()
	sg.StateGroupItems = make([]*wwise.StateGroupItem, sg.NumStateGroups)
	for i := range sg.StateGroupItems {
		sgi := wwise.NewStateGroupItem()
		sgi.StateGroupID = r.U32Unsafe()
		sgi.StateSyncType = r.U8Unsafe()
		sgi.NumStates = r.U8Unsafe()
		sgi.States = make([]*wwise.StateGroupItemState, sgi.NumStates)
		for i := range sgi.States {
			sgi.States[i] = &wwise.StateGroupItemState{
				StateID:         r.U32Unsafe(),
				StateInstanceID: r.U32Unsafe(),
			}
		}
		sg.StateGroupItems[i] = sgi
	}
	return sg
}

func parseRTPC(r *reader.BankReader) *wwise.RTPC {
	rtpc := wwise.NewRTPC()
	rtpc.NumRTPC = r.U16Unsafe()
	rtpc.RTPCItems = make([]*wwise.RTPCItem, rtpc.NumRTPC, rtpc.NumRTPC)
	for i := range rtpc.RTPCItems {
		ri := wwise.NewRTPCItem()
		ri.RTPCID = r.U32Unsafe()
		ri.RTPCType = r.U8Unsafe()
		ri.RTPCAccum = r.U8Unsafe()
		ri.ParamID = r.U8Unsafe()
		ri.RTPCCurveID = r.U32Unsafe()
		ri.Scaling = r.U8Unsafe()
		ri.NumRTPCGraphPoints = r.U16Unsafe()
		ri.RTPCGraphPoints = make([]*wwise.RTPCGraphPoint, ri.NumRTPCGraphPoints, ri.NumRTPCGraphPoints)
		for j := range ri.NumRTPCGraphPoints {
			rs := &wwise.RTPCGraphPoint{
				From: r.F32Unsafe(), To: r.F32Unsafe(), Interp: r.U32Unsafe(),
			}
			ri.RTPCGraphPoints[j] = rs
		}
		rtpc.RTPCItems[i] = ri
	}
	return rtpc
}

func parseChildren(r *reader.BankReader) *wwise.CntrChildren {
	cntrChildren := wwise.NewCntrChildren()
	cntrChildren.NumChild = r.U32Unsafe()
	cntrChildren.Children = make([]uint32, cntrChildren.NumChild)
	for i := range cntrChildren.Children {
		cntrChildren.Children[i] = r.U32Unsafe()
	}
	return cntrChildren
}
