package wwise

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wio"
)

// Refactor the getter to a single get prop value using PropTypeEnum

const SizeOfPropValue = 4

type PropValue struct {
	P   uint8
	V []byte
}
type PropBundle struct {
	// CProps uint8 // u8i
	// PIds []uint8 // CProps * u8i
	// PValues [][]byte // CProps * (Union[tid, uni / float32])

	Modulator    bool
	PropValues []PropValue
}

func (p *PropBundle) Clone() PropBundle {
	np := PropBundle{PropValues: make([]PropValue, len(p.PropValues))}
	for i, pv := range p.PropValues {
		np.PropValues[i].P = pv.P
		np.PropValues[i].V = bytes.Clone(pv.V)
	}
	return np
}

func (p *PropBundle) Encode(v int) []byte {
	size := p.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(p.PropValues)))
	for _, i := range p.PropValues {
		w.Append(i.P)
	}
	for _, i := range p.PropValues {
		w.AppendBytes(i.V)
	}
	return w.BytesAssert(int(size))
}

func (p *PropBundle) Size(v int) uint32 {
	return uint32(1 + len(p.PropValues) + SizeOfPropValue * len(p.PropValues))
}

func (p *PropBundle) HasPidRaw(pid uint8) (int, bool) {
	return sort.Find(len(p.PropValues), func(i int) int {
		if pid < p.PropValues[i].P {
			return -1
		} else if pid == p.PropValues[i].P {
			return 0
		} else {
			return 1
		}
	})
}

func (p *PropBundle) HasPid(pid PropType, v int) (int, bool) {
	return p.HasPidRaw(ForwardTranslateProp(pid, v))
}

func (p *PropBundle) Prop(pid PropType, v int) (int, *PropValue) {
	tp := ForwardTranslateProp(pid, v)
	if idx, in := p.HasPidRaw(tp); !in {
		return -1, nil
	} else {
		return idx, &p.PropValues[idx]
	}
}

func (p *PropBundle) Add(pid PropType, v int) {
	tp := ForwardTranslateProp(pid, v)
	if idx, in := p.HasPidRaw(tp); !in {
		p.PropValues = slices.Insert(p.PropValues, idx, PropValue{tp, []byte{0, 0, 0, 0}})
	}
}

func (p *PropBundle) AddWithVal(pid PropType, val [4]byte, v int) {
	tp := ForwardTranslateProp(pid, v)
	if idx, in := p.HasPidRaw(tp); !in {
		p.PropValues = slices.Insert(p.PropValues, idx, PropValue{tp, bytes.Clone(val[:])})
	}
}

func (p *PropBundle) SetPropByIdxF32(idx int, v float32) {
	if len(p.PropValues) <= 0 || idx >= len(p.PropValues) {
		return
	}
	binary.Encode(p.PropValues[idx].V, wio.ByteOrder, v)
}

func (p *PropBundle) Sort() {
	slices.SortFunc(p.PropValues, 
		func(a PropValue, b PropValue) int {
			if a.P < b.P {
				return -1
			}
			if a.P > b.P {
				return 1
			}
			return 0
		},
	)
}

func (p *PropBundle) AddBaseProp(v int) {
	for _, t := range BasePropType {
		tp := ForwardTranslateProp(t, v)
		if i, in := p.HasPidRaw(tp); !in {
			p.PropValues = slices.Insert(p.PropValues, i, PropValue{tp, []byte{0, 0, 0, 0}})
			return
		}
	}
}

func (p *PropBundle) ChangeBaseProp(idx int, nextPid PropType, v int) {
	tp := ForwardTranslateProp(nextPid, v)
	if !slices.Contains(BasePropType, nextPid) {
		return
	}
	if slices.ContainsFunc(p.PropValues, func(p PropValue) bool {
		return p.P == tp
	}) {
		return
	}
	p.PropValues[idx].P = tp
	for i := range 4 {
		p.PropValues[idx].V[i] = 0
	}
	p.Sort()
}

var BasePropChecker map[PropType]func(float32) error = map[PropType]func(float32) error {
	TVolume: utils.FloatInBound(-96.0, 12.0, PropLabel(TVolume)),
	TPitch: utils.FloatInBound(-2400, 2400, PropLabel(TPitch)),
	TLPF: utils.FloatInBound(0, 100, PropLabel(TLPF)),
	THPF: utils.FloatInBound(0, 100, PropLabel(THPF)),
	TMakeUpGain: utils.FloatInBound(-96.0, 96.0, PropLabel(TMakeUpGain)),
	TGameAuxSendVolume: utils.FloatInBound(-96.0, 12.0, PropLabel(TGameAuxSendVolume)),
	TInitialDelay: utils.FloatInBound(0.0, 60.0, PropLabel(TInitialDelay)),
}

var BaseRangePropChecker map[PropType]func(float32, bool) error = map[PropType]func(float32, bool) error {
	TVolume: utils.FloatRangeInBound(-108, 0, 108, PropLabel(TVolume)),
	TPitch: utils.FloatRangeInBound(-4800, 0, 4800, PropLabel(TPitch)),
	TLPF: utils.FloatRangeInBound(-100, 0, 100, PropLabel(TLPF)),
	THPF: utils.FloatRangeInBound(-100, 0, 100, PropLabel(THPF)),
	TMakeUpGain: utils.FloatRangeInBound(-192, 0, 192, PropLabel(TMakeUpGain)),
	TInitialDelay: utils.FloatRangeInBound(-60, 0, 60, PropLabel(TInitialDelay)),
}

func CheckBasePropVal(p PropType, v float32) error {
	if !slices.Contains(BasePropType, p) {
		return fmt.Errorf("Invaild base property ID %d", p)
	}
	c, in := BasePropChecker[p]
	if !in {
		panic("Panic Trap")
	}
	return c(v)
}

func CheckBaseRangeProp(p PropType, minV float32, maxV float32) error {
	if !slices.Contains(BasePropType, p) {
		return fmt.Errorf("Invaild base property ID %d", p)
	}
	c, in := BaseRangePropChecker[p]
	if !in {
		panic("Panic Trap")
	}
	if err := c(minV, false); err != nil { return err }
	if err := c(maxV, true); err != nil { return err }
	return nil
}

func (p *PropBundle) RemoveAllUserAuxSendVolumeProp(v int) {
	tps := make([]uint8, len(UserAuxSendVolumePropType), len(UserAuxSendVolumePropType))
	for i, at := range UserAuxSendVolumePropType {
		tps[i] = ForwardTranslateProp(at, v)
	}
	p.PropValues = slices.DeleteFunc(p.PropValues, func(p PropValue) bool {
		return slices.Contains(tps, p.P)
	})
}

func (p *PropBundle) AddReflectionBusVolume(v int) {
	tp := ForwardTranslateProp(TReflectionBusVolume, v)
	if i, in := p.HasPidRaw(tp); !in {
		p.PropValues = slices.Insert(p.PropValues, i, PropValue{tp, []byte{0, 0, 0, 0}})
	}
}

func (p *PropBundle) ReflectionBusVolume(v int) (int, *PropValue) {
	tp := ForwardTranslateProp(TReflectionBusVolume, v)
	if i, in := p.HasPidRaw(tp); !in {
		return -1, nil
	} else {
		return i, &p.PropValues[i]
	}
}

func (p *PropBundle) AddHDRActiveRange(v int) {
	tp := ForwardTranslateProp(THDRActiveRange, v)
	if i, in := p.HasPidRaw(tp); !in {
		p.PropValues = slices.Insert(p.PropValues, i, PropValue{tp, []byte{0, 0, 0, 0}})
	}
}

func (p *PropBundle) HDRActiveRange(v int) (int, *PropValue) {
	tp := ForwardTranslateProp(THDRActiveRange, v)
	if i, in := p.HasPidRaw(tp); !in {
		return -1, nil
	} else {
		return i, &p.PropValues[i]
	}
}

func (p *PropBundle) Remove(pid PropType, v int) {
	tp := ForwardTranslateProp(pid, v)
	i, found := p.HasPidRaw(tp)
	if !found {
		return
	}
	p.PropValues = slices.Delete(p.PropValues, i, i + 1)
}

func (p *PropBundle) DisplayProp() {
	var f float32
	for i := range p.PropValues {
		_, err := binary.Decode(p.PropValues[i].V, wio.ByteOrder, &f)
		if err != nil { panic(err) }
		fmt.Println(p.PropValues[i].P, f)
	}
}
type RangePropBundle struct {
	// CProps uint8 // u8i
	// PIds []uint8 // CProps * u8i
	Modulator     bool
	RangeValues []RangeValue // CProps * sizeof(RangeValue)
}

func (r *RangePropBundle) Clone() RangePropBundle {
	rp := RangePropBundle{RangeValues: make([]RangeValue, len(r.RangeValues))}
	for i, rv := range r.RangeValues {
		rp.RangeValues[i].P = rv.P
		rp.RangeValues[i].Min = bytes.Clone(rv.Min)
		rp.RangeValues[i].Max = bytes.Clone(rv.Max)
	}
	return rp
}

func (r *RangePropBundle) Encode(v int) []byte {
	size := r.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(r.RangeValues)))
	for _, i := range r.RangeValues {
		w.Append(i.P)
	}
	for _, i := range r.RangeValues {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (r *RangePropBundle) Size(v int) uint32 {
	return uint32(1 + len(r.RangeValues) + SizeOfRangeValue * len(r.RangeValues))
}

func (r *RangePropBundle) HasPidRaw(pid uint8) (int, bool) {
	return sort.Find(len(r.RangeValues), func(i int) int {
		if pid < r.RangeValues[i].P {
			return -1
		} else if pid == r.RangeValues[i].P {
			return 0
		} else {
			return 1
		}
	})
}

func (r *RangePropBundle) HasPid(pid PropType, v int) (int, bool) {
	return r.HasPidRaw(ForwardTranslateProp(pid, v))
}

func (r *RangePropBundle) Add(pid PropType, v int) {
	tp := ForwardTranslateProp(pid, v)
	if idx, in := r.HasPidRaw(tp); !in {
		r.RangeValues = slices.Insert(r.RangeValues, idx,
			RangeValue{tp, []byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}},
		)
	}
}

func (r *RangePropBundle) AddWithVal(pid PropType, lower [4]byte, upper [4]byte, v int) {
	tp := ForwardTranslateProp(pid, v)
	if idx, in := r.HasPidRaw(tp); !in {
		r.RangeValues = slices.Insert(r.RangeValues, idx,
			RangeValue{tp, lower[:], upper[:]},
		)
	}
}

func (r *RangePropBundle) Remove(pid PropType, v int) error {
	tp := ForwardTranslateProp(pid, v)
	i, found := r.HasPidRaw(tp)
	if !found {
		return fmt.Errorf("Failed to find property ID %d", pid)
	}
	r.RangeValues = slices.Delete(r.RangeValues, i, i + 1)
	return nil
}

func (r *RangePropBundle) Sort() {
	slices.SortFunc(r.RangeValues, 
		func(a RangeValue, b RangeValue) int {
			if a.P < b.P {
				return -1
			}
			if a.P > b.P {
				return 1
			}
			return 0
		},
	)
}

func (r *RangePropBundle) AddBaseProp(v int) {
	for _, t := range BaseRangePropType {
		tp := ForwardTranslateProp(t, v)
		if i, in := r.HasPidRaw(tp); !in {
			r.RangeValues = slices.Insert(r.RangeValues, i, 
				RangeValue{tp, []byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}},
			)
			break
		}
	}
}

func (r *RangePropBundle) ChangeBaseProp(idx int, nextPid PropType, v int) {
	if !slices.Contains(BasePropType, nextPid) {
		return
	}
	tp := ForwardTranslateProp(nextPid, v)
	if slices.ContainsFunc(r.RangeValues, func(r RangeValue) bool {
		return r.P == tp
	}) {
		return
	}
	r.RangeValues[idx].P = tp 
	for i := range 4 {
		r.RangeValues[idx].Min[i] = 0
		r.RangeValues[idx].Max[i] = 0
	}
	r.Sort()

}

func (r *RangePropBundle) SetPropMinByIdxF32(idx int, v float32) {
	if len(r.RangeValues) <= 0 || idx >= len(r.RangeValues) {
		return
	}
	binary.Encode(r.RangeValues[idx].Min, wio.ByteOrder, v)
}

func (r *RangePropBundle) SetPropMaxByIdxF32(idx int, v float32) {
	if len(r.RangeValues) <= 0 || idx >= len(r.RangeValues) {
		return
	}
	binary.Encode(r.RangeValues[idx].Max, wio.ByteOrder, v)
}

const SizeOfRangeValue = 8
type RangeValue struct {
	P     uint8
	Min []byte // Union[tid, uni / float32]
	Max []byte // Union[tid, uni / float32]
}

func (r *RangeValue) Encode() []byte {
	b := make([]byte, 0, SizeOfRangeValue)
	b = append(b, r.Min...) 
	b = append(b, r.Max...)
	assert.Equal(
		SizeOfRangeValue, len(b),
		"Encoded data of RangeValue has incorrect size",
	)
	return b
}
