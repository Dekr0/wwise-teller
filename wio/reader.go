package wio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

var ByteOrder = binary.LittleEndian

var InvalidSeek error = errors.New("Invalid Seek")

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
	nread, err := io.ReadFull(r.r, section)
	if err != nil {
		return nil, err
	}
	r.p += uint64(nread)
	return &Reader{ 0, bytes.NewReader(section), r.o }, nil
}

func (r *Reader) ReadNUnsafe(s uint64, reserved uint64) []byte {
	b, err := r.ReadN(s, reserved)
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return b
}

func (r *Reader) ReadN(s uint64, reserved uint64) ([]byte, error) {
	b := make([]byte, s, s + reserved)
	nread, err := io.ReadFull(r.r, b)
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *Reader) I64() (int64, error) {
	var v int64
	err := binary.Read(r.r, r.o, &v)
	r.p += 8
	return v, err
}

func (r *Reader) F64Unsafe() (v float64) {
	v, err := r.F64()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *Reader) F64() (v float64, err error) {
	err = binary.Read(r.r, r.o, &v)
	r.p += 8
	return v, err
}

func (r *Reader) StzUnsafe() (s []byte) {
	s, err := r.Stz()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return s
}

// TODO: Refactor 
func (r *Reader) Stz() (stz []byte, err error) {
	stz = []byte{}
	buffer := make([]byte, 1, 1)
	i := 0
	for {
		_, err := io.ReadFull(r.r, buffer)
		if err != nil {
			return nil, err
		}
		stz = append(stz, buffer...)
		i += 1
		r.p += 1
		if stz[len(stz) - 1] == 0 {
			break
		}
		if i >= 255 {
			return nil, errors.New("Zero-terminated string is too long!")
		}
	}
	return stz, nil
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
		slog.Error("Error log before panicking", "error", err)
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
		slog.Error("Error log before panicking", "error", err)
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
