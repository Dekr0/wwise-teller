package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestFNV(t *testing.T) {
	u, err := uuid.NewUUID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := u.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	h, err := FNV30(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(h)
}

