package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type FXPR struct {
	I uint8
	T []byte
	b []byte
}

func NewFXPR(I uint8, T []byte, b []byte) *FXPR {
	return &FXPR{I, T, b}
}

func (f *FXPR) Encode(ctx context.Context) ([]byte, error) {
	encoded := f.T
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(f.b)))
	if err != nil {
		panic(err)
	}
	return encoded, nil
}

func (f *FXPR) Tag() []byte {
	return f.T
}

func (f *FXPR) Idx() uint8 {
	return f.I
}
