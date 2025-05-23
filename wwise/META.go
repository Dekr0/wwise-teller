package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type META struct {
	I uint8
	T []byte
	B []byte
}

func NewMETA(I uint8, T []byte, Data []byte) *META {
	return &META{I, T, Data}
}

func (e *META) Encode(ctx context.Context) ([]byte, error) {
	encoded := e.T
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(e.B)))
	if err != nil {
		panic(err)
	}
	return encoded, nil
}

func (e *META) Tag() []byte {
	return e.T
}

func (e *META) Idx() uint8 {
	return e.I
}
