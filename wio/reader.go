package wio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var ByteOrder = binary.LittleEndian

var InvalidSeek error = errors.New("Invalid Seek")

// TYPE_SID: ShortID (uint32_t)
// TYPE_TID: target ShortID (uint32_t) same thing but for easier understanding of output
// TYPE_UNI: union (float / int32_t)
// TYPE_D64: double
// TYPE_F32: float
// TYPE_4CC: FourCC
// TYPE_S64: int64_t
// TYPE_U64: uint64_t
// TYPE_S32: int32_t
// TYPE_U32: uint32_t
// TYPE_S16: int16_t
// TYPE_U16: uint16_t
// TYPE_S8 : int8_t
// TYPE_U8 : uint8_t
// TYPE_VAR: variable size #u8/u16/u32
// TYPE_GAP: byte gap
// TYPE_STR: string
// TYPE_STZ: string (null-terminated)

// A helper struct that wraps an io.ReadSeeker. It provides a set of short 
// hand functions to read from this io.ReadSeeker, and produce commonly seen 
// data types.
// This helper struct is not designed for concurrent read with zero write and 
// zero copy. It's designed for an io.ReadSeeker that operates on bytes which 
// are not in the memory.
// If an io.ReadSeeker is operates on bytes that are in memory, such as 
// bytes.Reader, the following function will return copies instead of slices 
// that point to the same memory region:
// - Reader.ReadNUnsafe
// - Reader.ReadN, 
// - Reader.ReadAllUnsafe, 
// - Reader.ReadAll, 
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

// TODO: Return InPlaceReader instead of Reader. Only use this function when 
// the io.ReadSeeker of this Reader is operating bytes that are not in memory.
func (r *Reader) NewBufferReaderUnsafe(s uint64) *Reader {
	n, err := r.NewBufferReader(s)
	if err != nil { panic(err) }
	return n
}

// TODO: Return InPlaceReader instead of Reader. Only use this function when 
// the io.ReadSeeker of a Reader is operating bytes that are not in memory.
func (r *Reader) NewBufferReader(s uint64) (*Reader, error) {
	section := make([]byte, s)
	nread, err := r.r.Read(section)
	if err != nil {
		return nil, err
	}
	r.p += uint64(nread)
	return &Reader{ 0, bytes.NewReader(section), r.o }, nil
}

func (r *Reader) ReadNUnsafe(s uint64, reserved uint64) []byte {
	b, err := r.ReadN(s, reserved)
	if err != nil {
		panic(b)
	}
	return b
}

func (r *Reader) ReadN(s uint64, reserved uint64) ([]byte, error) {
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

func (r *Reader) ReadAllUnsafe() []byte {
	b, err := r.ReadAll()
	if err != nil { panic(err) }
	return b
}

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

// Since all read operations from bytes.Reader will create a copy for portion of 
// a slice, and in a concurrent read setting with zero write, it's ideal to 
// spawn readers that operates on some portion of a memory region without 
// copying, InPlaceReader fills up this gap.
// Example
//   Reader 1
//   |
//   | -----------------------
//   |     Reader 2          |
//   |     |                 |
//   |     |---------        |
//	 |     |        |        |
//   V     V        V        V
// [ 0, 1, 2, 3, 4, 5, 6, 7, 8 ]

type InPlaceReader struct {
	curr uint
	Buff []byte // Escape hatch for accessing this
	o binary.ByteOrder
}

// NOTES: Make sure pass in a slice of a memory region instead a copy of a 
// memory region
func NewInPlaceReader(buff []byte, o binary.ByteOrder) *InPlaceReader {
	return &InPlaceReader{curr: 0, Buff: buff, o: o}
}

func (r *InPlaceReader) Cap() uint {
	return uint(len(r.Buff))
}

func (r *InPlaceReader) Len() uint {
	return r.Cap() - r.curr
}

func (r *InPlaceReader) Tell() uint {
	return r.curr
}

func (r *InPlaceReader) AbsSeekUnsafe(j uint) {
	if err := r.AbsSeek(j); err != nil {
		panic(err)
	}
}

func (r *InPlaceReader) AbsSeek(j uint) error {
	if j >= 0 && j < r.Cap() {
		r.curr = j
		return nil
	}
	return InvalidSeek
}

func (r *InPlaceReader) RelSeekUnsafe(j int) {
	if err := r.RelSeek(j); err != nil {
		panic(err)
	}
}

func (r *InPlaceReader) RelSeek(j int) error {
	if j < 0 {
		flip := uint(-j)
		if flip > r.curr {
			return InvalidSeek
		}
		r.curr -= flip 
	} else {
		j := uint(j)
		if j + r.curr >= r.Cap() {
			return InvalidSeek
		}
		r.curr += j
	}
	return nil
}

func (r *InPlaceReader) NewInPlaceReader(s uint) (*InPlaceReader, error) {
	if s > r.Len() {
		return nil, io.ErrShortBuffer
	}
	nr := NewInPlaceReader(r.Buff[r.curr:r.curr + s], r.o)
	r.curr += s
	return nr, nil
}

func (r *InPlaceReader) NewInPlaceReaderOffset(offset uint, s uint) {}

func (r *InPlaceReader) NewInPlaceReaderUnsafe(s uint) (*InPlaceReader) {
	nr, err := r.NewInPlaceReader(s)
	if err != nil {
		panic(err)
	}
	return nr
}

func (r *InPlaceReader) U8Unsafe() uint8 {
	v, err := r.U8()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) U8() (uint8, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v uint8
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 1], r.o, &v)
	if err == nil {
		r.curr += 1
	} 
	return v, err
}


func (r *InPlaceReader) I8Unsafe() int8 {
	v, err := r.I8()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) I8() (int8, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v int8
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 1], r.o, &v)
	if err == nil {
		r.curr += 1
	} 
	return v, err
}

func (r *InPlaceReader) U16Unsafe() uint16 {
	v, err := r.U16()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) U16() (uint16, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v uint16
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 2], r.o, &v)
	if err == nil {
		r.curr += 2
	} 
	return v, err
}

func (r *InPlaceReader) I16Unsafe() int16 {
	v, err := r.I16()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) I16() (int16, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v int16
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 2], r.o, &v)
	if err == nil {
		r.curr += 2
	} 
	return v, err
}

func (r *InPlaceReader) U32Unsafe() uint32 {
	v, err := r.U32()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) U32() (uint32, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v uint32
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 4], r.o, &v)
	if err == nil {
		r.curr += 4
	} 
	return v, err
}

func (r *InPlaceReader) I32Unsafe() int32 {
	v, err := r.I32()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) I32() (int32, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v int32
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 4], r.o, &v)
	if err == nil {
		r.curr += 4
	} 
	return v, err
}

func (r *InPlaceReader) F32Unsafe() float32 {
	v, err := r.F32()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) F32() (float32, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v float32
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 4], r.o, &v)
	if err == nil {
		r.curr += 4
	} 
	return v, err
}

func (r *InPlaceReader) FourCC() ([]byte, error) {
	if r.Len() <= 0 {
		return nil, io.ErrShortBuffer
	}
	
	b := bytes.Clone(r.Buff[r.curr:r.curr + 4])
	r.curr += 4
	return b, nil
}

func (r *InPlaceReader) U64Unsafe() uint64 {
	v, err := r.U64()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) U64() (uint64, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v uint64
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 8], r.o, &v)
	if err == nil {
		r.curr += 8
	} 
	return v, err
}

func (r *InPlaceReader) I64Unsafe() int64 {
	v, err := r.I64()
	if err != nil { panic(err) }
	return v
}

func (r *InPlaceReader) I64() (int64, error) {
	if r.Len() <= 0 {
		return 0, io.ErrShortBuffer
	}
	var v int64
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 8], r.o, &v)
	if err == nil {
		r.curr += 8
	} 
	return v, err
}
