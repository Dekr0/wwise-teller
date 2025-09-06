package benchmark

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func BenchmarkDecodeBaselineStreaming(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	_, err := decoder.DecodeStreaming(
		b.Context(),
		filepath.Join(SoundBanksDir, bank),
		binary.LittleEndian,
		nil,
	)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkDecodeBaselineMem(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	_, err := decoder.DecodeMem(
		b.Context(),
		filepath.Join(SoundBanksDir, bank),
		binary.LittleEndian,
		nil,
	)
	if err != nil {
		b.Fatal(err)
	}
}
