package reader

import (
	"fmt"
	"encoding/binary"
)

type FixedSizeBlobWriter struct { blob []byte }

func NewFixedSizeBlobWriter(cap uint64) *FixedSizeBlobWriter {
	return &FixedSizeBlobWriter{make([]byte, 0, cap)}
}

/*
Append will panic upon exceeding maximum capacity since FixedSizeBlobWriter is 
initialized with assumption of capacity is set upfront.
*/
func (b *FixedSizeBlobWriter) AppendByte(data byte) {
	if cap(b.blob) < len(b.blob) + 1 {
		panic(fmt.Sprintf("Exceed maximum capacity"))
	}
	b.blob = append(b.blob, data)
}

/*
Append will panic upon exceeding maximum capacity since FixedSizeBlobWriter is 
initialized with assumption of capacity is set upfront.
*/
func (b *FixedSizeBlobWriter) AppendBytes(data []byte) {
	if cap(b.blob) < len(b.blob) + len(data) {
		panic(fmt.Sprintf("Exceed maximum capacity"))
	}
	b.blob = append(b.blob, data...)
}

/*
Append will panic upon error since FixedSizeBlobWriter is operate on in-memory 
slice.
*/
func (b *FixedSizeBlobWriter) Append(data any) {
	var err error
	b.blob, err = binary.Append(b.blob, BYTE_ORDER, data)
	if err != nil {
		panic(err)
	}
}

func (b *FixedSizeBlobWriter) GetBlob() []byte {
	return b.blob
}

/*
Before flushing out the blob, perform a panic check to ensure the entire buffer 
is used.
*/
func (b *FixedSizeBlobWriter) Flush(expect int) []byte {
	if len(b.blob) != expect {
		panic("The buffer is not fully used.")
	}
	return b.blob
}

func (b *FixedSizeBlobWriter) Len() int {
	return len(b.blob)
}

