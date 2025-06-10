package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

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
