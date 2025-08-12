package wio_test

import (
	bin "encoding/binary"
	"bytes"
	"slices"
	"testing"

	"github.com/Dekr0/wwise-teller/wio"
)

func TestVar(t *testing.T) {
	va := wio.VarT{B: make([]byte, 0, 10), V: 0}
	va.B = append(va.B, 0)

	if err := va.Set(65535); err != nil {
		t.Fatal(err)
	}

	b := bytes.NewBuffer(slices.Clone(va.B))
	vb := wio.VarT{}
	if err := wio.Var(b, bin.LittleEndian, &vb); err != nil {
		t.Fatal(err)
	}
	
	if vb.V != va.V {
		t.Fatalf("Expect vb = %d but receive %d", va.V, vb.V)
	}

	if !bytes.Equal(vb.B, va.B) {
		t.Fatalf("Expect bytes content of vb equals to va\nva: %v\nvb: %v", va.B, vb.B)
	}
}
