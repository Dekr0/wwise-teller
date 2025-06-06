package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type PLAT struct {
	I uint8
	T []byte
	B []byte
}

func NewPLAT(I uint8, T []byte, b []byte) *PLAT {
	return &PLAT{I, T, b}
}

func (p *PLAT) Encode(ctx context.Context) ([]byte, error) {
	encoded := p.T
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(p.B)))
	if err != nil {
		panic(err)
	}
	encoded = append(encoded, p.B...)
	return encoded, nil
}

func (p *PLAT) Tag() []byte {
	return p.T
}

func (p *PLAT) Idx() uint8 {
	return p.I
}
