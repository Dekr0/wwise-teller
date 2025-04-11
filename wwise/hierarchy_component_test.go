package wwise

import (
	"encoding/binary"
	"testing"
)

func TestFxChunkEncode(t *testing.T) {
	f := FxChunkItem{1, 12345, 0, 0}
	b := make([]byte, FX_CHUNK_ITEM_FIELD_SIZE, FX_CHUNK_ITEM_FIELD_SIZE)
	_, err := binary.Encode(b, binary.LittleEndian, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFxChunkMetadataEncode(t *testing.T) {
	f := FxChunkMetadataItem{1, 12345, 0}
	b := make([]byte, FX_CHUNK_METADATA_FIELD_SIZE, FX_CHUNK_METADATA_FIELD_SIZE)
	_, err := binary.Encode(b, binary.LittleEndian, f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPropBundleEncode(t *testing.T) {
	pValues := make([][]byte, 3, 3)
	pValues[0] = CreatePropValue(float32(2.4))
	pValues[1] = CreatePropValue(float32(95.0))
	pValues[2] = CreatePropValue(float32(0.08))
	p := PropBundle{ 3, []uint8{6, 7, 59}, pValues, }
	p.Encode()
}

func TestRangePropBundleEncode(t *testing.T) {
	rangeValues := make([]*RangeValue, 3, 3)
	rangeValues[0] = CreateRangeValue(float32(-2.4), float32(2.4))
	rangeValues[1] = CreateRangeValue(float32(90.0), float32(95.0))
	rangeValues[2] = CreateRangeValue(float32(0.016), float32(0.032))
	r := RangePropBundle{ 3, []uint8{6, 7, 59}, rangeValues, }
	r.Encode()
}
