package wwise

import "github.com/Dekr0/wwise-teller/wio"

type StatePropType uint16

type StatePropValue struct {
	P   StatePropType
	V []byte
}
type StatePropBundle struct {
	StatePropValues []StatePropValue
}

func (p *StatePropBundle) Encode(v int) []byte {
	size := p.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(uint16(len(p.StatePropValues)))
	for _, i := range p.StatePropValues {
		w.Append(i.P)
	}
	for _, i := range p.StatePropValues {
		w.AppendBytes(i.V)
	}
	return w.BytesAssert(int(size))
}

func (s *StatePropBundle) Size(int) uint32 {
	return 2 + uint32(len(s.StatePropValues)) * 6
}
