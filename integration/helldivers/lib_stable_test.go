package helldivers

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
)

func TestExtractSoundBanksStable(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelInfo.Level())

	data := "/mnt/d/Program Files/Steam/steamapps/common/Helldivers 2/data"
	target := filepath.Join(data, "22749a294788af66")
	if err := ExtractSoundBankStable(target, ".", false); err != nil {
		t.Fatal(err)
	}

	bnk, err := parser.ParseBank("./content_audio_vehicle_combat_walker.st_bnk", context.Background())
	if err != nil {
		t.Fatal(err)
	}

	b, err := bnk.Encode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	GenHelldiversPatchStable(b, bnk.META().B, ".")
}

func TestExtractSoundBanksPatchStable(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelInfo.Level())

	if err := ExtractSoundBankStable("./9ba626afa44a3aa3.patch_0", "..", false); err != nil {
		t.Fatal(err)
	}

	_, err := parser.ParseBank("../content_audio_vehicle_combat_walker.st_bnk", context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
