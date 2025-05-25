package wio

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"
)

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
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return nr
}

func (r *InPlaceReader) U8Unsafe() uint8 {
	v, err := r.U8()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) U8() (uint8, error) {
	if r.Len() < 1 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) I8() (int8, error) {
	if r.Len() < 1 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) U16() (uint16, error) {
	if r.Len() < 2 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) I16() (int16, error) {
	if r.Len() < 2 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) U32() (uint32, error) {
	if r.Len() < 4 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) I32() (int32, error) {
	if r.Len() < 4 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) F32() (float32, error) {
	if r.Len() < 4 {
		return 0, io.ErrShortBuffer
	}
	var v float32
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 4], r.o, &v)
	if err == nil {
		r.curr += 4
	} 
	return v, err
}

func (r *InPlaceReader) FourCCNoCopyUnsafe() ([]byte) {
	b, err := r.FourCCNoCopy()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return b
}

func (r *InPlaceReader) FourCCNoCopy() ([]byte, error) {
	if r.Len() < 4 {
		return nil, io.ErrShortBuffer
	}
	b := r.Buff[r.curr:r.curr + 4]
	r.curr += 4
	return b, nil
}

func (r *InPlaceReader) FourCCUnsafe() []byte {
	b, err := r.FourCC()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return b
}

func (r *InPlaceReader) FourCC() ([]byte, error) {
	if r.Len() < 4 {
		return nil, io.ErrShortBuffer
	}
	
	b := bytes.Clone(r.Buff[r.curr:r.curr + 4])
	r.curr += 4
	return b, nil
}

func (r *InPlaceReader) U64Unsafe() uint64 {
	v, err := r.U64()
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) U64() (uint64, error) {
	if r.Len() < 8 {
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
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return v
}

func (r *InPlaceReader) I64() (int64, error) {
	if r.Len() < 8 {
		return 0, io.ErrShortBuffer
	}
	var v int64
	
	_, err := binary.Decode(r.Buff[r.curr:r.curr + 8], r.o, &v)
	if err == nil {
		r.curr += 8
	} 
	return v, err
}

func (r *InPlaceReader) ReadNoCopyUnsafe(n uint) []byte {
	b, err := r.ReadNoCopy(n)
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return b
}

func (r *InPlaceReader) ReadNoCopy(n uint) ([]byte, error) {
	if r.Len() < n {
		return nil, io.ErrShortBuffer
	}
	b := r.Buff[r.curr:r.curr + n]
	r.curr += n
	return b, nil
}

func (r *InPlaceReader) ReadUnsafe(n uint) []byte {
	b, err := r.Read(n)
	if err != nil {
		slog.Error("Error log before panicking", "error", err)
		panic(err)
	}
	return b
}

func (r *InPlaceReader) Read(n uint) ([]byte, error) {
	if r.Len() < n {
		return nil, io.ErrShortBuffer
	}
	
	b := bytes.Clone(r.Buff[r.curr:r.curr + n])
	r.curr += n
	return b, nil
}

func (r *InPlaceReader) ReadAllNoCopy() ([]byte, error) {
	if r.Len() == 0 {
		return nil,io.ErrUnexpectedEOF
	}
	b := r.Buff[r.curr:]
	r.curr += uint(len(r.Buff))
	return b, nil
}
