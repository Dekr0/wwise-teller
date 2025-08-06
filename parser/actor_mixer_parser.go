package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseActorMixer(size uint32, r *wio.Reader, v int) *wwise.ActorMixer {
	assert.Equal(0, r.Pos(), "Actor mixer parser position doesn't start at position 0.")
	begin := r.Pos()
	a := wwise.ActorMixer{}
	a.Id = r.U32Unsafe()
	a.BaseParam = ParseBaseParam(r, v)
	ParseContainer(r, &a.Container, v)
	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return &a
}

func ParseLayerCntr(size uint32, r *wio.Reader, v int) *wwise.LayerCntr {
	assert.Equal(0, r.Pos(), "Layer container parser position doesn't start at position 0.")
	begin := r.Pos()
	l := wwise.LayerCntr{
		Id:        r.U32Unsafe(),
		BaseParam: ParseBaseParam(r, v),
	}
	ParseContainer(r, &l.Container, v)
	l.Layers = ParseLayers(r, make([]wwise.Layer, r.U32Unsafe()), v)
	l.IsContinuousValidation = r.U8Unsafe()
	
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &l
}

func ParseLayers(r *wio.Reader, layers []wwise.Layer, v int) []wwise.Layer {
	for i := range layers {
		l := &layers[i]
		l.Id = r.U32Unsafe()
		ParseRTPC(r, &l.InitialRTPC, v)
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

func ParseRanSeqCntr(size uint32, r *wio.Reader, v int) *wwise.RanSeqCntr {
	assert.Equal(0, r.Pos(), "Random / Sequence container parser position doesn't start at position 0.")
	begin := r.Pos()

	rs := wwise.RanSeqCntr{
		Id:               r.U32Unsafe(),
		BaseParam:       *ParseBaseParam(r, v),
	}
	ParsePlayListSetting(r, &rs.PlayListSetting)
	ParseContainer(r, &rs.Container, v)
	rs.PlayListItems = make([]wwise.PlayListItem, r.U16Unsafe())
	for i := range rs.PlayListItems {
		rs.PlayListItems[i].UniquePlayID = r.U32Unsafe()
		rs.PlayListItems[i].Weight = r.I32Unsafe()
	}

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &rs
}

func ParsePlayListSetting(r *wio.Reader, p *wwise.PlayListSetting) {
	p.LoopCount = r.U16Unsafe()
	p.LoopModMin = r.U16Unsafe()
	p.LoopModMax = r.U16Unsafe()
	p.TransitionTime = r.F32Unsafe()
	p.TransitionTimeModMin = r.F32Unsafe()
	p.TransitionTimeModMax = r.F32Unsafe()
	p.AvoidRepeatCount = r.U16Unsafe()
	p.TransitionMode = r.U8Unsafe()
	p.RandomMode = r.U8Unsafe()
	p.Mode = r.U8Unsafe()
	p.PlayListBitVector = r.U8Unsafe()
}

func ParseSound(size uint32, r *wio.Reader, v int) *wwise.Sound {
	assert.Equal(0, r.Pos(), "Sound parser position doesn't start 0.")
	begin := r.Pos()
	s := &wwise.Sound{}
	s.Id = r.U32Unsafe()
	s.BankSourceData = ParseBankSourceData(r, v)
	s.BaseParam = ParseBaseParam(r, v)
	end := r.Pos()
	if begin >= end {
		panic("reader consumes zero byte")
	}
	assert.Equal(uint64(size), end-begin,
		"the amount of bytes reader consume doesn't equal size in the hierarchy header",
	)
	return s
}

func ParseSwitchCntr(size uint32, r *wio.Reader, v int) *wwise.SwitchCntr {
	assert.Equal(0, r.Pos(), "Switch container parser position doesn't start at 0.")
	begin := r.Pos()
	s := wwise.SwitchCntr{
		Id:                     r.U32Unsafe(),
		BaseParam:              ParseBaseParam(r, v),
		GroupType:              r.U8Unsafe(),
		GroupID:                r.U32Unsafe(),
		DefaultSwitch:          r.U32Unsafe(),
		IsContinuousValidation: r.U8Unsafe(),
	}
	ParseContainer(r, &s.Container, v)
	s.SwitchGroups = make([]wwise.SwitchGroupItem, r.U32Unsafe())
	for i := range s.SwitchGroups {
		s.SwitchGroups[i].SwitchID = r.U32Unsafe()
		s.SwitchGroups[i].NodeList = make([]uint32, r.U32Unsafe())
		for j := range s.SwitchGroups[i].NodeList {
			s.SwitchGroups[i].NodeList[j] = r.U32Unsafe()
		}
	}
	NumSwitchParam := r.U32Unsafe()
	s.SwitchParams = make([]wwise.SwitchParam, NumSwitchParam)
	for i := range s.SwitchParams {
		s.SwitchParams[i].NodeId = r.U32Unsafe()
		if v <= 150 {
			s.SwitchParams[i].PlayBackBitVector = r.U8Unsafe()
			s.SwitchParams[i].ModeBitVector = r.U8Unsafe()
		} else {
			s.SwitchParams[i].PlayBackBitVector = r.U8Unsafe()
		}
		s.SwitchParams[i].FadeOutTime = r.I32Unsafe()
		s.SwitchParams[i].FadeInTime = r.I32Unsafe()
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader consume zero byte.")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to size in hierarchy header",
	)
	return &s
}

func ParseContainer(r *wio.Reader, c *wwise.Container, v int) {
	NumChild := r.U32Unsafe()
	c.Children = make([]uint32, NumChild)
	for i := range c.Children {
		c.Children[i] = r.U32Unsafe()
	}
}

func ParseAttenuation(size uint32, r *wio.Reader, v int) *wwise.Attenuation {
	assert.Equal(0, r.Pos(), "Attenuation parser position doesn't start at position 0.")
	begin := r.Pos()

	a := wwise.Attenuation{
		Id: r.U32Unsafe(),
		IsHeightSpreadEnabled: r.U8Unsafe(),
		IsConeEnabled: r.U8Unsafe(),
	}
	if v <= 141 {
		a.Curves = make([]int8, 7, 7)
	} else {
		a.Curves = make([]int8, 19, 19)
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
	ParseAttenuationConversionTables(r, a.AttenuationConversionTables, v)
	ParseRTPC(r, &a.RTPC, v)

	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &a
}

func ParseAttenuationConversionTables(r *wio.Reader, t []wwise.AttenuationConversionTable, v int) {
	for i := range t {
		t[i].EnumScaling = r.U8Unsafe()
		NumRTPCGraphPoints := r.U16Unsafe()
		t[i].RTPCGraphPointsX = make([]float32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		t[i].RTPCGraphPointsY = make([]float32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		t[i].RTPCGraphPointsInterp = make([]uint32, NumRTPCGraphPoints, NumRTPCGraphPoints)
		ParseRTPCGraphPoints(
			r,
			t[i].RTPCGraphPointsX,
			t[i].RTPCGraphPointsY,
			t[i].RTPCGraphPointsInterp,
			v,
		)
	}
}
