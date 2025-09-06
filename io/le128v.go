package io

const V128Mask    = 0b0111_1111
const V128MaxSize = 10

type V128 struct {
	B []byte
	V   u64
}

func V128Alloc(v *V128) {
	if v == nil {
		panic("Passing a nil point of V128")
	}
	v.B = make([]byte, 0, V128MaxSize)
	v.V = 0
}
