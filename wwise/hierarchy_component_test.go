package wwise

import (
	"encoding/binary"
	"testing"
)

func TestFxChunkEncode(t *testing.T) {
	f := FxChunkItem{1, 12345, 0, 0}
	b := make([]byte, SizeOfFxChunk, SizeOfFxChunk)
	_, err := binary.Encode(b, binary.LittleEndian, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFxChunkMetadataEncode(t *testing.T) {
	f := FxChunkMetadataItem{1, 12345, 0}
	b := make([]byte, SizeOfFxChunkMetadata, SizeOfFxChunkMetadata)
	_, err := binary.Encode(b, binary.LittleEndian, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPropBundleAddFull(t *testing.T) {
	p := NewPropBundle()
	for range len(PropLabel_140) {
		if _, err := p.New(); err != nil {
			t.Fail()
		}
	}
	if _, err := p.New(); err != nil {
		t.Fail()
		return
	}
	if len(p.PropValues) != len(PropLabel_140) {
		t.Fail()
		return
	}
}

func TestPropBundleNewExist(t *testing.T) {
	p := &PropBundle{
		[]PropValue{
			{0, []byte{0, 0, 0, 0}},
			{6, []byte{0, 0, 0, 0}},
			{73, []byte{0, 0, 0, 0}},
		},
	}
	if _, err := p.New(); err != nil {
		t.Fail()
		return
	}
	if _, found := p.HasPid(1); !found {
		t.Fail()
		return
	}
}

func TestPropBundleNewWithUpdate(t *testing.T) {
	p := &PropBundle{}
	if _, err := p.New(); err != nil {
		t.Fail()
		return
	}
	if _, err := p.New(); err != nil {
		t.Fail()
		return
	}
	p.UpdatePropF32(6, -96.0)
	if _, found := p.HasPid(6); !found {
		t.Fail()
		return
	}
}
