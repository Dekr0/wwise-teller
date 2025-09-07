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
	const bank = "content_audio_weapons_superearth.st_bnk"
	_, err := decoder.DecodeMem(
		t.Context(),
		filepath.Join(SoundBanksDir, bank),
		binary.LittleEndian,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
}
