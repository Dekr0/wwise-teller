package wwise

import (
	"fmt"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

// At decoding phase, the array of auxiliary bus IDs  must have four elements 
// regardless whether if it uses user-definded auxiliary send. If it uses it, 
// decode all the auxiliary bus IDs. Otherwise, it will zero out the entire 
// array.
// At encoding phase, if all auxiliary bus IDs are zero, that means it doesn't 
// use user-defined auxiliary send. It needs to set the corresponding bit to 
// zero.
type AuxParam struct {
	AuxBitVector            uint8     // U8x
	AuxIds                  [4]uint32 // 4 * tid
	RestoreAuxIds           [4]uint32
	ReflectionAuxBus        uint32    // tid
	RestoreReflectionAuxBus uint32
}

func NewAuxParam() *AuxParam {
	return &AuxParam{0, [4]uint32{0, 0, 0, 0}, [4]uint32{0, 0, 0, 0}, 0, 0}
}

func (a *AuxParam) OverrideAuxSends() bool {
	return wio.GetBit(a.AuxBitVector, 2)
}

func (a *AuxParam) SetOverrideAuxSends(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 2, set)
	if !a.OverrideAuxSends() {
		for i, aid := range a.RestoreAuxIds {
			a.AuxIds[i] = aid
		}
	}
}

func (a *AuxParam) HasAux() bool {
	return wio.GetBit(a.AuxBitVector, 3)
}

func (a *AuxParam) SetAux(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 3, set)
}

func (a *AuxParam) OverrideReflectionAuxBus() bool {
	return wio.GetBit(a.AuxBitVector, 4)
}

func (a *AuxParam) SetOverrideReflectionAuxBus(set bool) {
	a.AuxBitVector = wio.SetBit(a.AuxBitVector, 4, set)
	if a.OverrideReflectionAuxBus() {
		a.ReflectionAuxBus = a.RestoreReflectionAuxBus
	}
}

func (a *AuxParam) Encode() []byte {
	n := 0
	for _, a := range a.AuxIds {
		if a == 0 {
			n += 1
		}
	}
	if n == len(a.AuxIds) {
		a.SetAux(false)
	}

	a.assert()

	if !a.HasAux() {
		size := a.Size()
		w := wio.NewWriter(uint64(size))
		w.AppendByte(a.AuxBitVector)
		w.Append(a.ReflectionAuxBus)
		return w.BytesAssert(int(size))
	}

	size := a.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendByte(a.AuxBitVector)
	for _, id := range a.AuxIds { w.Append(id) }
	w.Append(a.ReflectionAuxBus)

	return w.BytesAssert(int(size))
}

func (a *AuxParam) Size() uint32 {
	if !a.HasAux() {
		return 5
	}
	return uint32(1 + 4 * 4 + 4)
}

func (a *AuxParam) assert() {
	if !a.HasAux() {
		for _, a := range a.AuxIds {
			msg := fmt.Sprintf(
				"Aux bit vector indicate user-defined auxiliary send is not " + 
				"used but auxiliary bus %d has non-zero ID", a,
			)
			assert.True(a <= 0, msg)
		}
		return
	}
}

