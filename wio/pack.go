package wio

import (
	bin "encoding/binary"
	"fmt"
	"io"
)

func MustWrite(w io.Writer, o bin.ByteOrder, data any, name string) {
	if err := bin.Write(w, o, data); err != nil {
		panic(fmt.Sprintf("Failed to write %s", name))
	}
}

func MustWriteTrack(w io.Writer, o bin.ByteOrder, data any, name string, p *uint32, size uint32) {
	if p != nil {
		*p += size
	}
	MustWrite(w, o, data, name)
}

func AssertFit(expect uint32, recv uint32) {
	if expect != recv {
		panic(fmt.Sprintf("Expecting %d bytes are all fitted but only %d bytes are used", expect, recv))
	}
}
