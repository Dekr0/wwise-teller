package io

import (
	bin "encoding/binary"
	"fmt"
	"io"
)

type EncoderCtx struct {
	Writer io.Writer
	Order  order
	Count  u32
}

func SizeOf(data any) u32 {
	switch data := data.(type) {
	case bool, i8, u8, *bool, *i8, *u8:
		return Size8
	case []bool:
		return u32(len(data))
	case []i8:
		return u32(len(data))
	case []u8:
		return u32(len(data))
	case i16, u16, *i16, *u16:
		return Size16
	case []i16:
		return Size16 * u32(len(data))
	case []u16:
		return Size16 * u32(len(data))
	case i32, u32, *i32, *u32:
		return Size32
	case []i32:
		return Size32 * u32(len(data))
	case []u32:
		return Size32 * u32(len(data))
	case i64, u64, *i64, *u64:
		return Size64
	case []i64:
		return Size64 * u32(len(data))
	case []u64:
		return Size64 * u32(len(data))
	case f32, *f32:
		return Size32
	case f64, *f64:
		return Size64
	case []f32:
		return Size32 * u32(len(data))
	case []f64:
		return Size64 * u32(len(data))
	}
	panic("Unsupported primitive type")
}

func Encode(e *EncoderCtx, data any) (err error) {
	size := SizeOf(data)
	if err := bin.Write(e.Writer, e.Order, data); err != nil {
		return err
	}
	e.Count += size
	return err
}

func EncodeBytes(e *EncoderCtx, data []byte) (err error) {
	n, err := e.Writer.Write(data)
	if err != nil {
		return err
	}
	l := len(data)
	if n != l {
		return fmt.Errorf("Failed to write all %d bytes", l)
	}
	e.Count += u32(l)
	return err
}

func EncodeStruct(e *EncoderCtx, data any, size u32) (err error) {
	if err := bin.Write(e.Writer, e.Order, data); err != nil {
		return err
	}
	e.Count += size
	return err
}

func AssertEncodeLimit(e *EncoderCtx, prev u32, expect u32) (err error) {
	recv := e.Count - prev
	if recv != expect {
		return fmt.Errorf("Expect encoded data has a size of %d but receive %d", expect, recv)
	}
	return nil
}
