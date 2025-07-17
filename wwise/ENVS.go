package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type ENVS struct {
	I uint8
	T []byte
	B []byte
}

func (e *ENVS) Encode(ctx context.Context, v int) ([]byte, error) {
	encoded := e.T
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(e.B)))
	if err != nil {
		panic(err)
	}
	encoded = append(encoded, e.B...)
	return encoded, nil
}

func NewENVS(I uint8, T []byte, b []byte) *ENVS {
	return &ENVS{I, T, b}
}

func (e *ENVS) Tag() []byte {
	return e.T
}

func (e *ENVS) Idx() uint8 {
	return e.I
}
