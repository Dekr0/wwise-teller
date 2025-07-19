package automation

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
)

func TestReplace(t *testing.T) {
	if err := utils.InitTmp(); err != nil {
		t.Fatal(err)
	}
	bnk, err := parser.ParseBank(
		filepath.Join(testStBankDir, "content_audio_wep_ar61_marauder.st_bnk"),
		context.Background(),
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = bnk.Encode(context.Background(), true, false)
	if err != nil {
		t.Fatal(err)
	}
	utils.CleanTmp()
}
