package wio

import (
	bin "encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

func U8(r io.Reader, o bin.ByteOrder) (u8 uint8, err error) {
	buf := make([]uint8, 1, 1)
	_, err = r.Read(buf)
	u8 = buf[0]
	return u8, err
}

func MustU8(r io.Reader, o bin.ByteOrder) (u8 uint8) {
	var err error
	u8, err = U8(r, o)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a 8-bits integer: %s", err.Error()))
	}
	return u8
}

func I8(r io.Reader, o bin.ByteOrder) (int8, error) {
	i8, err := U8(r, o)
	return int8(i8), err
}

func MustI8(r io.Reader, o bin.ByteOrder) int8 {
	i8 := MustU8(r, o)
	return int8(i8)
}

func U16(r io.Reader, o bin.ByteOrder) (uint16, error) {
	data := make([]uint8, 2, 2)
	if _, err := r.Read(data); err != nil {
		return 0, err
	}
	return o.Uint16(data), nil
}

func MustU16(r io.Reader, o bin.ByteOrder) (u16 uint16) {
	var err error
	u16, err = U16(r, o)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a 16-bits integer: %s", err.Error()))
	}
	return u16
}

func I16(r io.Reader, o bin.ByteOrder) (int16, error) {
	i16, err := U16(r, o)
	return int16(i16), err
}

func MustI16(r io.Reader, o bin.ByteOrder) int16 {
	i16 := MustU16(r, o)
	return int16(i16)
}

func U32(r io.Reader, o bin.ByteOrder) (uint32, error) {
	data := make([]uint8, 4, 4)
	if _, err := r.Read(data); err != nil {
		return 0, err
	}
	return o.Uint32(data), nil
}

func MustU32(r io.Reader, o bin.ByteOrder) (u32 uint32) {
	var err error
	u32, err = U32(r, o)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a 32-bits integer: %s", err.Error()))
	}
	return u32
}

func I32(r io.Reader, o bin.ByteOrder) (int32, error) {
	i32, err := U32(r, o)
	return int32(i32), err
}

func MustI32(r io.Reader, o bin.ByteOrder) int32 {
	i32 := MustU32(r, o)
	return int32(i32)
}

func U64(r io.Reader, o bin.ByteOrder) (uint64, error) {
	data := make([]uint8, 8, 8)
	if _, err := r.Read(data); err != nil {
		return 0, err
	}
	return o.Uint64(data), nil
}

func MustU64(r io.Reader, o bin.ByteOrder) (u64 uint64) {
	var err error
	u64, err = U64(r, o)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a 64-bits integer: %s", err.Error()))
	}
	return u64
}

func I64(r io.Reader, o bin.ByteOrder) (int64, error) {
	i64, err := U64(r, o)
	return int64(i64), err
}

func MustI64(r io.Reader, o bin.ByteOrder) int64 {
	i64 := MustU64(r, o)
	return int64(i64)
}

func F32(r io.Reader, o bin.ByteOrder) (f32 float32, err error) {
	u32, err := U32(r, o)
	if err != nil {
		return 0.0, err
	}
	return math.Float32frombits(u32), nil
}

func MustF32(r io.Reader, o bin.ByteOrder) float32 {
	f32, err := F32(r, o)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a 32-bits float: %s", err.Error()))
	}
	return f32
}

func Stz(r io.Reader) (stz []byte, err error) {
	stz = []byte{}
	buf := make([]byte, 1, 1)
	i := 0
	for {
		_, err := r.Read(buf)
		if err != nil {
			return nil, err
		}
		stz = append(stz, buf...)
		i += 1
		if stz[len(stz) - 1] == 0 {
			break
		}
		if i >= 255 {
			return stz, errors.New("Zero-terminated string is larger than 255 bytes!")
		}
	}
	return stz, nil
}

func MustStz(r io.Reader) (stz []byte) {
	var err error
	stz, err = Stz(r)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a zero-terminated string: %s", err.Error()))
	}
	return stz
}

func Var(r io.Reader, o bin.ByteOrder, v *VarT) (err error) {
	if v == nil {
		return fmt.Errorf("The provided LE128 handle is is nil")
	}
	buf := make([]uint8, 1, 1)
	_, err = r.Read(buf)
	if err != nil {
		return err
	}

	v.B = make([]uint8, 0, 10)
	v.B, v.V = append(v.B, buf[0]), uint64(buf[0]) & 0b0111_1111

	i := 0
	for buf[0] & 0b1000_0000 > 0 && i < 10 {
		_, err = r.Read(buf)
		if err != nil {
			return err
		}
		v.B, v.V = append(v.B, buf[0]), (v.V << 7) | (uint64(buf[0]) & 0b0111_1111)
		i += 1
	}

	if i >= 10 {
		return fmt.Errorf("LE128 is larger than 10 bytes.")
	}

	return err
}

func MustVar(r io.Reader, o bin.ByteOrder, v *VarT) {
	err := Var(r, o, v)
	if err != nil {
		panic(fmt.Sprintf("Failed to read a LE128: %s", err.Error()))
	}
}
