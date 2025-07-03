package automation

import (
	"context"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/waapi"
)

func TestReplace(t *testing.T) {
	if err := waapi.InitTmp(); err != nil {
		t.Fatal(err)
	}
	bnk, err := parser.ParseBank(
		"../tests/default_st_bnks/content_audio_wep_ar61_marauder.st_bnk",
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = bnk.Encode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	waapi.CleanTmp()
}
