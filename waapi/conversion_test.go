package waapi

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateConversionList(t *testing.T) {
	in := []string{
		"../tests/wavs/ak74/core_01.wav",
		"../tests/wavs/ak74/core_02.wav",
		"../tests/wavs/ak74/core_03.wav",
		"../tests/wavs/ak74/core_04.wav",
	}
	var err error
	for i := range in {
		in[i], err = filepath.Abs(in[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	out := make([]string, len(in))
	ctx := context.Background()
	if err := InitTmp(); err != nil {
		t.Fatal(err)
	}
	wsource, err := CreateConversionList(ctx, in, out, "Vorbis Quality High", true); 
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(wsource)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
	CleanTmp()
}

func TestWwiseConversion(t *testing.T) {
	in := []string{
		"../tests/wavs/ak74/core_01.wav",
		"../tests/wavs/ak74/core_02.wav",
		"../tests/wavs/ak74/core_03.wav",
		"../tests/wavs/ak74/core_04.wav",
	}
	var err error
	for i := range in {
		in[i], err = filepath.Abs(in[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	out := make([]string, len(in))
	ctx := context.Background()
	if err := InitTmp(); err != nil {
		t.Fatal(err)
	}
	wsource, err := CreateConversionList(ctx, in, out, "Vorbis Quality High", true); 
	if err != nil {
		t.Fatal(err)
	}
	proj, err := filepath.Abs("WwiseTeller") 
	if err != nil {
		t.Fatal(err)
	}
	if err := WwiseConversion(ctx, wsource, proj); err != nil {
		t.Fatal(err)
	}
	CleanTmp()
}
