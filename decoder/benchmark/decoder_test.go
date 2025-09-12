package benchmark

import (
	"encoding/binary"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func BenchmarkDecodeBaseline(b *testing.B) {
	slog.SetDefault(slog.New(slog.DiscardHandler))

	const bank = "content_audio_weapons_superearth.st_bnk"

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}
	hircDecodeOpt.NumRoutine = 0

	for b.Loop() {
		_, err := decoder.AllocDecode(
			b.Context(),
			filepath.Join(SoundBanksDir, bank),
			binary.LittleEndian,
			&bankDecodeOpt,
			&hircDecodeOpt,
		)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
	}
}
