package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type STMG struct {
	I uint8
	T []byte
	B []byte
}

func NewSTMG(I uint8, T []byte, b []byte) *STMG {
	return &STMG{I, T, b}
}

func (s *STMG) Encode(ctx context.Context) ([]byte, error) {
	encoded := s.T 
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(s.B)))
	if err != nil {
		panic(err)
	}
	encoded = append(encoded, s.B...)
	return encoded, nil
}

func (s *STMG) Tag() []byte {
	return s.T
}

func (s *STMG) Idx() uint8 {
	return s.I
}
