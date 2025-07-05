package waapi

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/utils"
)

func TestCreateConversionList(t *testing.T) {
	in := []string{
		"tests/wavs/ak74/core_01.wav",
		"tests/wavs/ak74/core_02.wav",
		"tests/wavs/ak74/core_03.wav",
		"tests/wavs/ak74/core_04.wav",
	}

	wavsMapping := make(map[string]struct{}, len(in))
	var err error
	for i := range in {
		in[i], err = filepath.Abs(in[i])
		if err != nil {
			t.Fatal(err)
		}
		wavsMapping[in[i]] = struct{}{}
	}

	wemsMapping := make(map[string]struct{}, len(wavsMapping))
	ctx := context.Background()
	if err := utils.InitTmp(); err != nil {
		t.Fatal(err)
	}
	wsource, err := CreateConversionList(ctx, wavsMapping, wemsMapping, "Vorbis Quality High", true); 
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(wsource)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
	utils.CleanTmp()
}

func TestWwiseConversion(t *testing.T) {
	in := []string{
		"tests/wavs/ak74/core_01.wav",
		"tests/wavs/ak74/core_02.wav",
		"tests/wavs/ak74/core_03.wav",
		"tests/wavs/ak74/core_04.wav",
	}

	wavsMapping := make(map[string]struct{}, len(in))
	var err error
	for i := range in {
		in[i], err = filepath.Abs(in[i])
		if err != nil {
			t.Fatal(err)
		}
		wavsMapping[in[i]] = struct{}{}
	}

	wemsMapping := make(map[string]struct{}, len(wavsMapping))
	ctx := context.Background()
	if err := utils.InitTmp(); err != nil {
		t.Fatal(err)
	}
	wsource, err := CreateConversionList(ctx, wavsMapping, wemsMapping, "Vorbis Quality High", true); 
	if err != nil {
		t.Fatal(err)
	}
	proj, err := filepath.Abs("WwiseTeller/WwiseTeller.wproj") 
	if err != nil {
		t.Fatal(err)
	}
	if err := WwiseConversion(ctx, wsource, proj); err != nil {
		t.Fatal(err)
	}
	utils.CleanTmp()
}
