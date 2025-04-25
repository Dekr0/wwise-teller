package wio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var ByteOrder = binary.LittleEndian

/*
TYPE_SID: ShortID (uint32_t)
TYPE_TID: target ShortID (uint32_t) same thing but for easier understanding of output
TYPE_UNI: union (float / int32_t)
TYPE_D64: double
TYPE_F32: float
TYPE_4CC: FourCC
TYPE_S64: int64_t
TYPE_U64: uint64_t
TYPE_S32: int32_t
TYPE_U32: uint32_t
TYPE_S16: int16_t
TYPE_U16: uint16_t
TYPE_S8 : int8_t
TYPE_U8 : uint8_t
TYPE_VAR: variable size #u8/u16/u32
TYPE_GAP: byte gap
TYPE_STR: string
TYPE_STZ: string (null-terminated)
*/

/*
A wrapper around io.ReadSeeker for reading binary data in different section of
a sound bank, or the entire section of a sound bank.
It maintains byte order consistency and the absolute location (from origin, 0
byte).
*/
type Reader struct {
	p uint64
	r io.ReadSeeker
	o binary.ByteOrder
}

func NewReader(r io.ReadSeeker, o binary.ByteOrder) *Reader {
	return &Reader{0, r, o}
}

func (r *Reader) ByteOrder() binary.ByteOrder {
	return r.o
}

func (r *Reader) Pos() uint64 {
	return r.p
}

func (r *Reader) SeekStart(p uint64) error {
	n, err := r.r.Seek(int64(p), io.SeekStart)
	if err != nil {
		return err
	}
	r.p = uint64(n)
	return nil
}

func (r *Reader) SeekCurrent(p int64) error {
	n, err := r.r.Seek(p, io.SeekCurrent)
	if err != nil {
		return err
	}
	r.p = uint64(n)
	return nil
}


// Read a given amount of bytes and create a new reader that operates on those 
// bytes.
func (r *Reader) NewBufferReader(s uint64) (*Reader, error) {
	section := make([]byte, s)
	nread, err := r.r.Read(section)
	if err != nil {
		return nil, err
	}
	r.p += uint64(nread)
	return &Reader{ 0, bytes.NewReader(section), r.o }, nil
}

func (r *Reader) NewBufferReaderUnsafe(s uint64) *Reader {
	n, err := r.NewBufferReader(s)
	if err != nil { panic(err) }
	return n
}

func (r *Reader) ReadFullUnsafe(s uint64, reserved uint64) []byte {
	b, err := r.ReadFull(s, reserved)
	if err != nil {
		panic(b)
	}
	return b
}


// Read a give amount of bytes at the current position. Reserve some bytes if 
// necessary to prevent overhead from resizing.
// Use this when reader is operating on data that is not yet in memory. Otherwise
// , it will large amount of copying.
func (r *Reader) ReadFull(s uint64, reserved uint64) ([]byte, error) {
	b := make([]byte, s, s + reserved)
	nread, err := r.r.Read(b)
	if err != nil {
		return nil, err
	}
	if uint64(nread) > s {
		err := fmt.Errorf(
			"Reader.ReadFull: io.Reader read more than the specified size." +
			" Expected read size %d, actual read size %d.", s, nread,
		)
		return nil, err
	}
	r.p += s
	return b, nil
}

// Read data from the current position until an error occurs, or until EOF.
func (r *Reader) ReadAll() ([]byte, error) {
	b, err := io.ReadAll(r.r)
	if err != nil { return nil, err }
	r.p += uint64(len(b))
	return b, nil
}

func (r *Reader) U8Unsafe() uint8 {
	v, err := r.U8()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) U8() (uint8, error) {
	var v uint8
	err := binary.Read(r.r, r.o, &v)
	r.p += 1
	return v, err
}


func (r *Reader) I8Unsafe() int8 {
	v, err := r.I8()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) I8() (int8, error) {
	var v int8
	err := binary.Read(r.r, r.o, &v)
	r.p += 1
	return v, err
}

func (r *Reader) U16Unsafe() uint16 {
	v, err := r.U16()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) U16() (uint16, error) {
	var v uint16
	err := binary.Read(r.r, r.o, &v)
	r.p += 2
	return v, err
}

func (r *Reader) I16Unsafe() int16 {
	v, err := r.I16()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) I16() (int16, error) {
	var v int16
	err := binary.Read(r.r, r.o, &v)
	r.p += 2
	return v, err
}

func (r *Reader) U32Unsafe() uint32 {
	v, err := r.U32()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) U32() (uint32, error) {
	var v uint32
	err := binary.Read(r.r, r.o, &v)
	r.p += 4
	return v, err
}

func (r *Reader) I32Unsafe() int32 {
	v, err := r.I32()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) I32() (int32, error) {
	var v int32
	err := binary.Read(r.r, r.o, &v)
	r.p += 4
	return v, err
}

func (r *Reader) F32Unsafe() float32 {
	v, err := r.F32()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) F32() (float32, error) {
	var v float32
	err := binary.Read(r.r, r.o, &v)
	r.p += 4
	return v, err
}

func (r *Reader) FourCC() ([]byte, error) {
	buf := make([]byte, 4, 4)
	err := binary.Read(r.r, r.o, &buf)
	r.p += 4
	return buf, err
}

func (r *Reader) U64Unsafe() uint64 {
	v, err := r.U64()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) U64() (uint64, error) {
	var v uint64
	err := binary.Read(r.r, r.o, &v)
	r.p += 8
	return v, err
}

func (r *Reader) I64Unsafe() int64 {
	v, err := r.I64()
	if err != nil { panic(err) }
	return v
}

func (r *Reader) I64() (int64, error) {
	var v int64
	err := binary.Read(r.r, r.o, &v)
	r.p += 8
	return v, err
}
