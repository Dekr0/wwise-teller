package wwise

import (
	"context"
	"encoding/binary"

	"github.com/Dekr0/wwise-teller/wio"
)

type STID struct {
	I uint8
	T []byte
	B []byte
}

func NewSTID(I uint8, T []byte, b []byte) *STID {
	return &STID{I, T, b}
}

func (s *STID) Encode(ctx context.Context) ([]byte, error) {
	encoded := s.T 
	encoded, err := binary.Append(encoded, wio.ByteOrder, uint32(len(s.B)))
	if err != nil {
		panic(err)
	}
	return encoded, nil
}

func (s *STID) Tag() []byte {
	return s.T
}

func (s *STID) Idx() uint8 {
	return s.I
}
