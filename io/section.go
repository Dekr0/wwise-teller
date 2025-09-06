package io

import (
	bin "encoding/binary"
	"io"
)

func SectionU8(r *io.SectionReader, o order) (v u8, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionU8P(r *io.SectionReader, o order) u8 {
	v, err := SectionU8(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionI8(r *io.SectionReader, o order) (v i8, err error) {
	uv, err := SectionU8(r, o)
	if err != nil {
		return v, err
	}
	v = i8(uv)
	return v, nil 
}

func SectionI8P(r *io.SectionReader, o order) i8 {
	v, err := SectionI8(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionU16(r *io.SectionReader, o order) (v u16, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionU16P(r *io.SectionReader, o order) u16 {
	v, err := SectionU16(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionI16(r *io.SectionReader, o order) (v i16, err error) {
	uv, err := SectionI16(r, o)
	if err != nil {
		return v, err
	}
	v = i16(uv)
	return v, nil 
}

func SectionI16P(r *io.SectionReader, o order) i16 {
	v, err := SectionI16(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionU32(r *io.SectionReader, o order) (v u32, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionU32P(r *io.SectionReader, o order) u32 {
	v, err := SectionU32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionI32(r *io.SectionReader, o order) (v i32, err error) {
	uv, err := SectionI32(r, o)
	if err != nil {
		return v, err
	}
	v = i32(uv)
	return v, nil 
}

func SectionI32P(r *io.SectionReader, o order) i32 {
	v, err := SectionI32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionF32(r *io.SectionReader, o order) (v f32, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionF32P(r *io.SectionReader, o order) f32 {
	v, err := SectionF32(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionU64(r *io.SectionReader, o order) (v u64, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionU64P(r *io.SectionReader, o order) u64 {
	v, err := SectionU64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionI64(r *io.SectionReader, o order) (v i64, err error) {
	uv, err := SectionI64(r, o)
	if err != nil {
		return v, err
	}
	v = i64(uv)
	return v, nil 
}

func SectionI64P(r *io.SectionReader, o order) i64 {
	v, err := SectionI64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

func SectionF64(r *io.SectionReader, o order) (v f64, err error) {
	err = bin.Read(r, o, &v)
	return v, err
}

func SectionF64P(r *io.SectionReader, o order) f64 {
	v, err := SectionF64(r, o)
	if err != nil {
		panic(err)
	}
	return v
}

// All struct read functions need to allocate ahead of time.
func SectionStruct[T any](r *io.SectionReader, o order, data T) error {
	return bin.Read(r, o, data)
}

func SectionStructP[T any](r *io.SectionReader, o order, data T) {
	err := SectionStruct(r, o, data)
	if err != nil {
		panic(err)
	}
}
