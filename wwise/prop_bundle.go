package wwise

import (
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const SizeOfPropValue = 4

var ExceedMaxNumOfProperty = errors.New("Exceed max number of property.")

type PropValue struct {
	P uint8
	V []byte
}
type PropBundle struct {
	// CProps uint8 // u8i
	// PIds []uint8 // CProps * u8i
	// PValues [][]byte // CProps * (Union[tid, uni / float32])
	PropValues []*PropValue
}

func NewPropBundle() *PropBundle {
	return &PropBundle{[]*PropValue{}}
}

func (p *PropBundle) Encode() []byte {
	size := p.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(p.PropValues)))
	for _, i := range p.PropValues {
		w.AppendByte(i.P)
	}
	for _, i := range p.PropValues {
		w.AppendBytes(i.V)
	}
	return w.BytesAssert(int(size))
}

func (p *PropBundle) Size() uint32 {
	return uint32(1 + len(p.PropValues) + SizeOfPropValue * len(p.PropValues))
}

func (p *PropBundle) HasPid(pId uint8) (int, bool) {
	return sort.Find(len(p.PropValues), func(i int) int {
		if pId < p.PropValues[i].P {
			return -1
		} else if pId == p.PropValues[i].P {
			return 0
		} else {
			return 1
		}
	})
}

func (p *PropBundle) UpdatePropBytes(pId uint8, b []byte) {
	if len(b) != 4 {
		panic("Inserting a property value using a byte slice with length " + 
			"less than 4.")
	}
	i, found := p.HasPid(pId)
	if !found {
		p.PropValues = slices.Insert(p.PropValues, i, &PropValue{pId, b})
	} else {
		p.PropValues[i].V = b
	}
}

func (p *PropBundle) UpdatePropI32(pId uint8, v int32) {
	i, found := p.HasPid(pId)
	w := wio.NewWriter(4)
	w.Append(v)
	b := w.BytesAssert(4)
	if !found {
		p.PropValues = slices.Insert(p.PropValues, i, &PropValue{pId, b})
	} else {
		p.PropValues[i].V = b
	}
}

func (p *PropBundle) UpdatePropF32(pId uint8, v float32) {
	i, found := p.HasPid(pId)
	w := wio.NewWriter(4)
	w.Append(v)
	b := w.BytesAssert(4)
	if !found {
		p.PropValues = slices.Insert(p.PropValues, i, &PropValue{pId, b})
	} else {
		p.PropValues[i].V = b
	}
}

func (p *PropBundle) Sort() {
	slices.SortFunc(p.PropValues, 
		func(a *PropValue, b *PropValue) int {
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

// Add a new property. The new property ID will look for the one that is not in 
// used.
func (p *PropBundle) New() (uint8, error) {
	if len(p.PropValues) == len(PropLabel_140) {
		return 0, ExceedMaxNumOfProperty
	}
	// Mid point
	if len(p.PropValues) == 0 {
		p.PropValues = append(p.PropValues, &PropValue{
			uint8(len(PropLabel_140) / 2), []byte{0, 0, 0, 0},
		})
		return p.PropValues[0].P, nil
	}
	PL := p.PropValues[0].P
	right := p.PropValues[len(p.PropValues) - 1].P
	PR := uint8(len(PropLabel_140)) - right - 1
	if PL > 0 || PR > 0 {
		if PL >= PR {
			p.PropValues = append(
				[]*PropValue{{PL - 1, []byte{0, 0, 0, 0}}}, 
				p.PropValues...
				)
			return PL - 1, nil
		}
		p.PropValues = append(
			p.PropValues, 
			&PropValue{right + 1, []byte{0, 0, 0, 0}}, 
		)
		return right + 1, nil 
	}
	for i := range len(PropLabel_140) {
		if !slices.ContainsFunc(p.PropValues, func(p *PropValue) bool {
			return p.P == uint8(i)
		}) {
			p.UpdatePropBytes(uint8(i), []byte{0, 0, 0, 0})
			return uint8(i), nil
		}
	}
	panic("Dead code path")
}

func (p *PropBundle) Remove(pId uint8) error {
	i, found := p.HasPid(pId)
	if !found {
		return fmt.Errorf("Failed to find property ID %d", pId)
	}
	p.PropValues = slices.Delete(p.PropValues, i, i + 1)
	return nil
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
	RangeValues []*RangeValue // CProps * sizeof(RangeValue)
}

func NewRangePropBundle() *RangePropBundle {
	return &RangePropBundle{[]*RangeValue{}}
}

func (r *RangePropBundle) Encode() []byte {
	size := r.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(len(r.RangeValues)))
	for _, i := range r.RangeValues {
		w.AppendByte(i.PId)
	}
	for _, i := range r.RangeValues {
		w.AppendBytes(i.Encode())
	}
	return w.BytesAssert(int(size))
}

func (r *RangePropBundle) Size() uint32 {
	return uint32(1 + len(r.RangeValues) + SizeOfRangeValue * len(r.RangeValues))
}

func (r *RangePropBundle) HasPid(pID uint8) (int, bool) {
	return sort.Find(len(r.RangeValues), func(i int) int {
		if pID < r.RangeValues[i].PId {
			return -1
		} else if pID == r.RangeValues[i].PId {
			return 0
		} else {
			return 1
		}
	})
}

func (r *RangePropBundle) UpdatePropBytes(pId uint8, min []byte, max []byte) {
	if len(min) != 4 || len(max) != 4 {
		panic("Inserting a property value using a byte slice with length " + 
			"less than 4.")
	}
	i, found := r.HasPid(pId)
	if !found {
		r.RangeValues = slices.Insert(r.RangeValues, i, &RangeValue{pId, min, max})
	} else {
		r.RangeValues[i].Min = min
		r.RangeValues[i].Max = max
	}
}

func (r *RangePropBundle) UpdatePropF32(pId uint8, min float32, max float32) {
	i, found := r.HasPid(pId)
	w := wio.NewWriter(4)
	w.Append(min)
	minB := w.BytesAssert(4)
	w = wio.NewWriter(4)
	w.Append(min)
	maxB := w.BytesAssert(4)
	if !found {
		r.RangeValues = slices.Insert(r.RangeValues, i, &RangeValue{pId, minB, maxB})
	} else {
		r.RangeValues[i].Min = minB
		r.RangeValues[i].Max = maxB
	}
}

func (r *RangePropBundle) UpdatePropI32(pId uint8, min int32, max int32) {
	i, found := r.HasPid(pId)
	w := wio.NewWriter(4)
	w.Append(min)
	minB := w.BytesAssert(4)
	w = wio.NewWriter(4)
	w.Append(min)
	maxB := w.BytesAssert(4)
	if !found {
		r.RangeValues = slices.Insert(r.RangeValues, i, &RangeValue{pId, minB, maxB})
	} else {
		r.RangeValues[i].Min = minB
		r.RangeValues[i].Max = maxB
	}
}

func (r *RangePropBundle) New() (uint8, error) {
	if len(r.RangeValues) == len(PropLabel_140) {
		return 0, ExceedMaxNumOfProperty
	}
	// Mid point
	if len(r.RangeValues) == 0 {
		r.RangeValues = append(r.RangeValues, &RangeValue{
			uint8(len(PropLabel_140) / 2),
			[]byte{0, 0, 0, 0},
			[]byte{0, 0, 0, 0},
		})
		return r.RangeValues[0].PId, nil
	}
	PL := r.RangeValues[0].PId
	right := r.RangeValues[len(r.RangeValues) - 1].PId
	PR := uint8(len(PropLabel_140)) - right - 1
	if PL > 0 || PR > 0 {
		if PL >= PR {
			r.RangeValues = append(
				[]*RangeValue{{PL - 1, []byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}}}, 
				r.RangeValues...
			)
			return PL - 1, nil
		}
		r.RangeValues= append(
			r.RangeValues, 
			&RangeValue{right + 1, []byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}}, 
		)
		return right + 1, nil 
	}
	for i := range len(PropLabel_140) {
		if !slices.ContainsFunc(r.RangeValues, func(r *RangeValue) bool {
			return r.PId == uint8(i)
		}) {
			r.UpdatePropBytes(uint8(i), []byte{0, 0, 0, 0}, []byte{0, 0, 0, 0})
			return uint8(i), nil
		}
	}
	panic("Dead code path")
}

func (r *RangePropBundle) Remove(pId uint8) error {
	i, found := r.HasPid(pId)
	if !found {
		return fmt.Errorf("Failed to find property ID %d", pId)
	}
	r.RangeValues = slices.Delete(r.RangeValues, i, i + 1)
	return nil
}

func (r *RangePropBundle) Sort() {
	slices.SortFunc(r.RangeValues, 
		func(a *RangeValue, b *RangeValue) int {
			if a.PId < b.PId {
				return -1
			}
			if a.PId > b.PId {
				return 1
			}
			return 0
		},
	)
}

const SizeOfRangeValue = 8
type RangeValue struct {
	PId uint8
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

