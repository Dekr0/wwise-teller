package wio

import (
	"bytes"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	reader := bytes.NewReader(data)

	data[2] = 16
	data[7] = 255

	d, _ := io.ReadAll(reader)
	t.Log(d)

	d[0] = 96

	r := NewInPlaceReader(data, ByteOrder)
	if r.Cap() != uint(len(data)) {
		t.Fatalf("T1: Expect %d but receive %d", len(data), r.Cap())
	}
	data[2] = 16
	if r.Len() != uint(len(data)) {
		t.Fatalf("T2: Expect %d but receive %d", len(data), r.Len())
	}
	t.Log(r.U8Unsafe())
	if r.Len() != uint(len(data)) - 1 {
		t.Fatalf("T3: Expect %d but receive %d", len(data) - 1, r.Len())
	}

	t.Log(r.U8Unsafe())
	if a := r.U8Unsafe(); a != 16 {
		t.Fatalf("T4: Expect 16 but received %d", a)
	}
	if r.Len() != uint(len(data)) - 3 {
		t.Fatalf("T5: Expect %d but receive %d", len(data) - 3, r.Len())
	}

	sr, err := r.NewInPlaceReader(r.Len())
	if err != nil {
		t.Fatalf("T6: %v", err)
	}
	data[7] = 255
	if sr.Len() != 5 {
		t.Fatal("T7: Expect 5 but receive", sr.Len())
	}
	for sr.Len() > 0 {
		t.Log(sr.U8Unsafe())
	}

	if err := sr.RelSeek(-5); err != nil {
		t.Fatal(err)
	}

	if err := sr.AbsSeek(5); err == nil {
		t.Fatal("Expecting error")
	}
	
	if err := sr.AbsSeek(4); err != nil {
		t.Fatal(err)
	}

	t.Log(sr.U8Unsafe())
	if sr.curr != 5 {
		t.Fatalf("Expecting 5 but received %d", sr.curr)
	}

	if err := sr.RelSeek(-3); err != nil {
		t.Fatal(err)
	}
	if a := sr.U8Unsafe(); a != 6 {
		t.Fatalf("Expecting 6 but received %d", a)
	}
	if a := sr.U8Unsafe(); a != 7 {
		t.Fatalf("Expecting 7 but received %d", a)
	}
	t.Log(sr.U8Unsafe())
}
