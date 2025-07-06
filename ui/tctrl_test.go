package ui

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/waapi"
)

var testBanksDir = os.Getenv("TESTS")

func TestCreatePlayerNoCache(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	utils.InitTmp()
	defer utils.CleanTmp()
	waapi.InitWEMCache()
	defer waapi.CleanWEMCache()

	bank := "content_audio_wep_ar19_liberator.st_bnk"

	bnk, err := parser.ParseBank(filepath.Join(testBanksDir, bank), context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}

	tab := bank_explorer.BankTab{}

	data := bnk.DATA()

	sid := uint32(279367945)
	wemData, in := data.AudiosMap[sid]
	if !in {
		t.Fatalf("No audio source has ID %d.", sid)
	}

	task := createPlayerNoCacheTask(&tab, sid, wemData)
	task(t.Context())

	_, in = tab.WEMExportCache.Load(sid)
	if !in {
		t.Fatalf("Data of audio source %d is not transformed into WAVE and cached", sid)
	}
}

