package decoder_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func TestDecodeBKHD(t *testing.T) {
	entries, err := os.ReadDir(SoundBanksDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if bank, err := decoder.Decode(
			t.Context(),
			filepath.Join(SoundBanksDir, entry.Name()),
			binary.LittleEndian,
			&decoder.DecoderOption{},
		); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("Version: %d; Id: %d", bank.BKHD.Version, bank.BKHD.Id)
		}
	}
}
