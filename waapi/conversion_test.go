package waapi

import (
	"context"
	"testing"
)

func TestWaveToWem(t *testing.T) {
	in := []string{
		"../tests/wavs/ak74/core_01.wav",
		"../tests/wavs/ak74/core_02.wav",
		"../tests/wavs/ak74/core_03.wav",
		"../tests/wavs/ak74/core_04.wav",
	}
	out := make([]string, len(in))
	ctx := context.Background()
	wsource, err := CreateConversionList( ctx, in, out, "Vorbis Quality High"); 
	if err != nil {
		t.Fatal(err)
	}
	if err := WwiseConversion(
		ctx,
		wsource,
		"../AudioConversionProject/AudioConversionProject.wproj",
	); err != nil {
		t.Fatal(err)
	}
}
