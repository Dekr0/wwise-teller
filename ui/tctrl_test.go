package ui

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/waapi"
)

var testBanksDir = os.Getenv("TESTS")

func TestCreatePlayerNoCache(t *testing.T) {
	utils.InitTmp()
	defer utils.CleanTmp()
	waapi.InitWEMCache()
	defer waapi.CleanWEMCache()

	bank := "content_audio_wep_ar19_liberator.st_bnk"

	bnk, err := parser.ParseBank(filepath.Join(testBanksDir, bank), context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}

	var cacheMap sync.Map
	data := bnk.DATA()

	sid := uint32(279367945)
	wemData, in := data.AudiosMap[sid]
	if !in {
		t.Fatalf("No audio source has ID %d.", sid)
	}

	task := createPlayerNoCacheTask(&cacheMap, sid, wemData)
	task(t.Context())

	v, in := cacheMap.Load(sid)
	if !in {
		t.Fatalf("Data of audio source %d is not transformed into WAVE and cached", sid)
	}
	_, err = GlobalCtx.Manager.Player(v.(string))
	if err != nil {
		t.Fatal(err)
	}
}

