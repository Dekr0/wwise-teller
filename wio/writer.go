package wio

import (
	"encoding/binary"
	"fmt"
	"io"
)

type FixedWriter struct { b []byte }

func NewWriter(cap uint64) *FixedWriter {
	return &FixedWriter{make([]byte, 0, cap)}
}

// Panic upon exceeding maximum capacity given the assumption of which capacity 
// is fixed.
func (w *FixedWriter) AppendByte(b byte) {
	if cap(w.b) < len(w.b) + 1 {
		panic(fmt.Sprintf("Exceed maximum capacity"))
	}
	w.b = append(w.b, b)
}

// Panic upon exceeding maximum capacity given the assumption of which capacity 
// is fixed.
func (w *FixedWriter) AppendBytes(b []byte) {
	if cap(w.b) < len(w.b) + len(b) {
		panic(fmt.Sprintf("Exceed maximum capacity"))
	}
	w.b = append(w.b, b...)
}

// Panic upon exceeding maximum capacity given the assumption of which capacity 
// is fixed.
func (w *FixedWriter) Append(v any) {
	var err error
	w.b, err = binary.Append(w.b, ByteOrder, v)
	if err != nil { panic(err) }
}

func (w *FixedWriter) Bytes() []byte {
	return w.b
}

// Perform assertion to make sure that a buffer is completely used before 
// returning the result.
func (w *FixedWriter) BytesAssert(expect int) []byte {
	if len(w.b) != expect {
		panic("The buffer size mismatch with expected size.")
	}
	return w.b
}

func (w *FixedWriter) Len() int {
	return len(w.b)
}

func GetBit(v uint8, pos int) bool {
	return (v >> pos) & 1 > 0
}

func SetBit(v uint8, pos int, set bool) uint8 {
	if !set {
		return v & (^(1 << pos))
	}
	return v | (1 << pos)
}

type BinaryWriteHelper struct {
	nwrite uint
	w      io.Writer
}

func NewBinaryWriteHelper(w io.Writer) *BinaryWriteHelper {
	return &BinaryWriteHelper{0, w}
}

func (w *BinaryWriteHelper) U8(v uint8) error {
	w.nwrite += 1
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) I8(v int8) error {
	w.nwrite += 1
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) U16(v uint16) error {
	w.nwrite += 2
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) I16(v int16) error {
	w.nwrite += 2
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) U32(v uint32) error {
	w.nwrite += 4
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) I32(v int32) error {
	w.nwrite += 4
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) F32(v float32) error {
	w.nwrite += 4
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) U64(v uint64) error {
	w.nwrite += 8
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) I64(v int64) error {
	w.nwrite += 8
	return binary.Write(w.w, ByteOrder, v)
}

func (w *BinaryWriteHelper) Bytes(p []byte) error {
	w.nwrite += uint(len(p))
	_, err := w.w.Write(p)
	return err
}

func (w *BinaryWriteHelper) Tell() uint { return w.nwrite }
