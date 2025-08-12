package wio

import (
	"fmt"
	"slices"
)

type VarT struct {
	B []byte
	V uint64
}

func (v *VarT) Set(val uint64) error {
	if val == 0 {
		v.B = slices.Delete(v.B, 0, len(v.B))
		v.B, v.V = append(v.B, 0), 0
		return nil
	}

	i, calc := 0, val
	for i < 10 {
		calc >>= 7
		if calc == 0 {
			break
		}
		i += 1
	}
	if i >= 10 {
		return fmt.Errorf("%d is too large for 128LE", val)
	}

	v.B, v.V, calc = slices.Delete(v.B, 0, len(v.B)), val, val

	curr := uint64(0)
	for range i + 1 {
		curr = calc & 0b0111_1111
		calc >>= 7
		v.B = append(v.B, uint8(curr))
	}
	
	slices.Reverse(v.B)
	for i := range v.B {
		if i < len(v.B) - 1 {
			v.B[i] |= 0b1000_0000
		}
	}
	return nil
}
