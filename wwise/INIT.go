package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type INIT struct {
	I uint8
	T []byte
	B []byte
}

func NewINIT(I uint8, T []byte, b []byte) *INIT {
	return &INIT{I, T, b}
}

func (i *INIT) Encode(ctx context.Context) ([]byte, error) {
	encoded := i.T
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(i.B)))
	if err != nil {
		panic(err)
	}
	return encoded, nil
}

func (i *INIT) Tag() []byte {
	return i.T
}

func (i *INIT) Idx() uint8 {
	return i.I
}
