package helldivers

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestExtractSoundBanksStable(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelInfo.Level())

	data := "/mnt/d/Program Files/Steam/steamapps/common/Helldivers 2/data"
	target := filepath.Join(data, "c76c97b3dfb67c5c")
	if err := ExtractSoundBankStable(target, ".", false); err != nil {
		t.Fatal(err)
	}
}

func TestGenPatchStable(t *testing.T) {
	ctx := context.Background()
	bnksName := []string{
		"./8342649675791264267.st_bnk",
	}
	bnks := []wwise.Bank{}
	for _, bnkName := range bnksName {
		bnk, err:= parser.ParseBank(bnkName, ctx, false)
		if err != nil {
			t.Fatal(err)
		}
		bnks = append(bnks, *bnk)
	}
	bnksData := [][]byte{}
	metas := [][]byte{}
	for _, bnk := range bnks {
		bnkData, err := bnk.Encode(ctx, false)
		if err != nil {
			t.Fatal(err)
		}
		bnksData = append(bnksData, bnkData)
		meta := bnk.META()
		if err != nil {
			t.Fatal(err)
		}
		metas = append(metas, meta.B)
	}
	if err := GenHelldiversPatchStableMulti(bnksData, metas, "."); err != nil {
		t.Fatal(err)
	}
}

func TestExtractSoundBanksPatchStable(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelInfo.Level())

	if err := ExtractSoundBankStable("./9ba626afa44a3aa3.patch_0", "..", false); err != nil {
		t.Fatal(err)
	}

	_, err := parser.ParseBank("../content_audio_vehicle_combat_walker.st_bnk", context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
}
