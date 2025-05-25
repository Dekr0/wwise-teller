package helldivers

import (
	"context"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
)

func TestSoundBankExtracting(t *testing.T) {
	target := "/mnt/D/Program Files/Steam/steamapps/common/Helldivers 2/data/22749a294788af66"
	if err := ExtractSoundBank(nil, target, ".", false); err != nil {
		t.Fatal(err)
	}
}

func TestSoundBankPatching(t *testing.T) {
	target := "../../tests/st_bnk/content_audio_vehicle_combat_walker.st_bnk"
	bnk, err := parser.ParseBank(target, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	bnkData, err := bnk.Encode(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	GenHelldiversPatch(context.Background(), bnkData, bnk.META().B, ".")
}

func TestExtractSoundBankFromPatch(t *testing.T) {
	target := "../../tests/patch/9ba626afa44a3aa3.patch_0"
	if err := ExtractSoundBank(nil, target, ".", false); err != nil {
		t.Fatal(err)
	}
}

func TestParseRepatchedSoundBank(t *testing.T) {
	target := "../../tests/content_audio_vehicle_combat_walker.st_bnk"
	parser.ParseBank(target, context.Background())
}
