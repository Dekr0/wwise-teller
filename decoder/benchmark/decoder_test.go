package benchmark

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func BenchmarkDecodeBaseline(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	for b.Loop() {
		_, err := decoder.Decode(
			b.Context(),
			filepath.Join(SoundBanksDir, bank),
			binary.LittleEndian,
			nil,
			nil,
		)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
	}
}
