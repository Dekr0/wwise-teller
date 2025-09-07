package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder/benchmark/baseline"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")
var BigFile string = os.Getenv("BIGFILE")

func benchmarkReadOnceBigFile(b *testing.B) {
	for b.Loop() {
		baseline.ReadOnce(BigFile)
	}
}

func benchmarkBufferReadBigFile(b *testing.B) {
	const recvSize = 4096
	const bufSize = 4096
	for b.Loop() {
		baseline.BufferRead(recvSize, bufSize, BigFile)
	}
}

func BenchmarkReadOnceReadLargest(b *testing.B) {
	const name = "content_audio_weapons_superearth.st_bnk"
	path := filepath.Join(SoundBanksDir, name)
	for b.Loop() {
		baseline.ReadOnce(path)
	}
}

func BenchmarkUnbufferReadLargest(b *testing.B) {
	const recvSize = 4096

	const name = "content_audio_weapons_superearth.st_bnk"
	path := filepath.Join(SoundBanksDir, name)

	for b.Loop() {
		baseline.UnbufferRead(recvSize, path)
	}
}

func BenchmarkBufferReadLargest(b *testing.B) {
	const recvSize = 4096
	const bufSize = 4096

	const name = "content_audio_weapons_superearth.st_bnk"
	path := filepath.Join(SoundBanksDir, name)

	for b.Loop() {
		baseline.BufferRead(recvSize, bufSize, path)
	}
}
