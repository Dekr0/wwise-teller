package io

import (
	bin "encoding/binary"
	"io"
)

// # Naming convention
// - If a read function contains a T, this function will perform tracking on # of
// many bytes read.
// - If a read function contains a P, this function will panic instead of
// returning error.

// Usage
// If a struct is deterministic (i.e. size is fixed; no variable-size field),
// use

func U8(r io.Reader, o order) (v u8, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func U8T(r io.Reader, o order, t *int) (v u8, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size8
	}
	return v, err
}

func U8P(r io.Reader, o order) (v u8) {
	v, err := U8(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func U8PT(r io.Reader, o order, t *int) (v u8) {
	v, err := U8T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func I8(r io.Reader, o order) (v i8, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func I8T(r io.Reader, o order, t *int) (v i8, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size8
	}
	return v, err
}

func I8P(r io.Reader, o order) (v i8) {
	v, err := I8(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func I8PT(r io.Reader, o order, t *int) (v i8) {
	v, err := I8T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func U16(r io.Reader, o order) (v u16, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func U16T(r io.Reader, o order, t *int) (v u16, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size16
	}
	return v, err
}

func U16P(r io.Reader, o order) (v u16) {
	v, err := U16(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func U16PT(r io.Reader, o order, t *int) (v u16) {
	v, err := U16T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func I16(r io.Reader, o order) (v i16, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func I16T(r io.Reader, o order, t *int) (v i16, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size16
	}
	return v, err
}

func I16P(r io.Reader, o order) (v i16) {
	v, err := I16(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func I16PT(r io.Reader, o order, t *int) (v i16) {
	v, err := I16T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func U32(r io.Reader, o order) (v u32, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func U32T(r io.Reader, o order, t *int) (v u32, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size32
	}
	return v, err
}

func U32P(r io.Reader, o order) (v u32) {
	v, err := U32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func U32PT(r io.Reader, o order, t *int) (v u32) {
	v, err := U32T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func I32(r io.Reader, o order) (v i32, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func I32T(r io.Reader, o order, t *int) (v i32, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size32
	}
	return v, err
}

func I32P(r io.Reader, o order) (v i32) {
	v, err := I32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func I32PT(r io.Reader, o order, t *int) (v i32) {
	v, err := I32T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func F32(r io.Reader, o order) (v f32, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func F32T(r io.Reader, o order, t *int) (v f32, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size32
	}
	return v, err
}

func F32P(r io.Reader, o order) (v f32) {
	v, err := F32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func F32PT(r io.Reader, o order, t *int) (v f32) {
	v, err := F32T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func U64(r io.Reader, o order) (v u64, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func U64T(r io.Reader, o order, t *int) (v u64, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size64
	}
	return v, err
}

func U64P(r io.Reader, o order) (v u64) {
	v, err := U64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func U64PT(r io.Reader, o order, t *int) (v u64) {
	v, err := U64T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func I64(r io.Reader, o order) (v i64, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func I64T(r io.Reader, o order, t *int) (v i64, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size64
	}
	return v, err
}

func I64P(r io.Reader, o order) (v i64) {
	v, err := I64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func I64PT(r io.Reader, o order, t *int) (v i64) {
	v, err := I64T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

func F64(r io.Reader, o order) (v f64, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func F64T(r io.Reader, o order, t *int) (v f64, err error) {
	err = bin.Read(r, o, &v)
	if (err != nil) {
		return v, err
	}
	if (t != nil) {
		*t += Size64
	}
	return v, err
}

func F64P(r io.Reader, o order) (v f64) {
	v, err := F64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func F64PT(r io.Reader, o order, t *int) (v f64) {
	v, err := F64T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}

// All struct read functions need to allocate ahead of time.
func Struct[T any](r io.Reader, o order, data T) error {
	return bin.Read(r, o, data)
}

func StructP[T any](r io.Reader, o order, data T) {
	err := Struct(r, o, data)
	if err != nil {
		panic(err)
	}
}

func StructT[T any](r io.Reader, o order, data T, size int, t *int) (err error) {
	err = bin.Read(r, o, data)
	if err != nil {
		return err
	}
	if t != nil {
		*t += size
	}
	return err
}

func StructPT[T any](r io.Reader, o order, data T, size int, t *int) {
	err := StructT(r, o, data, size, t)
	if err != nil {
		panic(err)
	}
}

// All Stz read functions need to allocate ahead of time.
func Stz(r io.Reader, o order, buf []byte) ([]byte, error) {
	var err error
	var i u8 = 0
	var b byte = 0
	for {
		err = bin.Read(r, o, &b)
		if err != nil {
			return nil, err
		}

		buf = append(buf, b)

		i += 1

		if buf[len(buf) - 1] == 0 {
			break
		}

		if i >= MaxStzSize {
			return nil, ExceedStzSize
		}
	}
	return buf, err
}

func StzT(r io.Reader, o order, buf []byte, t *int) ([]byte, error) {
	var err error
	var i u8 = 0
	var b byte = 0
	for {
		err = bin.Read(r, o, &b)
		if err != nil {
			return nil, err
		}

		buf = append(buf, b)

		i += 1

		if buf[len(buf) - 1] == 0 {
			break
		}

		if i >= MaxStzSize {
			return nil, ExceedStzSize
		}
	}

	if t != nil {
		*t += int(i)
	}

	return buf, err
}

func StzP(r io.Reader, o order, buf []byte) []byte {
	var err error
	buf, err = Stz(r, o, buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func StzPT(r io.Reader, o order, buf []byte, t *int) []byte {
	var err error
	buf, err = StzT(r, o, buf, t)
	if err != nil {
		panic(err)
	}
	return buf
}

// All V128 read functions need to allocate ahead of time.
func VV128(r io.Reader, o order) (v *V128, err error) {
	var cur u8 = 0
	err = bin.Read(r, o, &cur)
	if err != nil {
		return nil, err
	}

	v = &V128{ make([]byte, 0, 8), 0 }

	v.V = u64(cur) & V128Mask
	v.B = append(v.B, cur)

	i := 0

	for cur & V128Mask > 0 && i < V128MaxSize {
		err = bin.Read(r, o, &cur)
		if err != nil {
			return nil, err
		}

		v.B = append(v.B, cur)
		v.V = (v.V << 7) | (u64(cur) & V128Mask)

		i += 1
	}

	if i >= 10 {
		return nil, ExceedV128Size
	}

	return v, nil
}

func VV128T(r io.Reader, o order, t *int) (v *V128, err error) {
	var cur u8 = 0
	err = bin.Read(r, o, &cur)
	if err != nil {
		return nil, err
	}

	v = &V128{ make([]byte, 0, 16), 0 }

	v.V = u64(cur) & V128Mask
	v.B = append(v.B, cur)

	i := 0

	for cur & V128Mask > 0 && i < V128MaxSize {
		err = bin.Read(r, o, &cur)
		if err != nil {
			return nil, err
		}

		v.B = append(v.B, cur)
		v.V = (v.V << 7) | (u64(cur) & V128Mask)

		i += 1
	}

	if i >= 10 {
		return nil, ExceedV128Size
	}

	if t != nil {
		*t += i
	}

	return v, nil 
}

func VV128P(r io.Reader, o order) *V128 {
	v, err := VV128(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func VV128PT(r io.Reader, o order, t *int) *V128 {
	v, err := VV128T(r, o, t)
	if err != nil {
		panic(err)
	}
	return v
}
