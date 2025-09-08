package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder/benchmark/baseline"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")
var BigFile string = os.Getenv("BIGFILE")

func BenchmarkReadOnceBigFile(b *testing.B) {
	for b.Loop() {
		baseline.ReadOnce(BigFile)
	}
}

func BenchmarkBufferReadBigFile(b *testing.B) {
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

func BenchmarkWriteOnceLargest(b *testing.B) {
	os.Remove("tmp")

	const name = "content_audio_weapons_superearth.st_bnk"
	path := filepath.Join(SoundBanksDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		f, err := baseline.WriteOnce(data)

		if err != nil {
			b.Fatal(err)
		}

		{
			b.StopTimer()
			if f != nil {
				if err = f.Close(); err != nil {
					b.Fatal(err)
				}
			}
			if err = os.Remove("tmp"); err != nil {
				b.Fatal(err)
			}
			b.StartTimer()
		}
	}
}

func BenchmarkManyWriteLargest(b *testing.B) {
	os.Remove("tmp")

	const name = "content_audio_weapons_superearth.st_bnk"
	const writeSize = 4096

	path := filepath.Join(SoundBanksDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for b.Loop() {
		f, err := baseline.ManyWrite(data, writeSize)

		if err != nil {
			b.Fatal(err)
		}

		{
			b.StopTimer()
			if f != nil {
				if err = f.Close(); err != nil {
					b.Fatal(err)
				}
			}
			if err = os.Remove("tmp"); err != nil {
				b.Fatal(err)
			}
			b.StartTimer()
		}
	}
}

func TestManyWriteLargest(t *testing.T) {
	const name = "content_audio_weapons_superearth.st_bnk"
	const writeSize = 4096

	path := filepath.Join(SoundBanksDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	f, err := baseline.ManyWrite(data, writeSize)

	if err != nil {
		t.Fatal(err)
	}

	if f != nil {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	}
	if err = os.Remove("tmp"); err != nil {
		t.Fatal(err)
	}
}
