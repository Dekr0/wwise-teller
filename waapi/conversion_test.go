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
	if _, err := WavToWem(
		context.Background(),
		in, out,
		"../AudioConversionProject/AudioConversionProject.wproj",
		"Vorbis Quality High",
	); err != nil {
		t.Fatal(err)
	}
}
