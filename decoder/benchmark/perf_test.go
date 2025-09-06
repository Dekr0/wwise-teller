package perf_test

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func BenchmarkUnbufferReadLargest(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	f, err := os.Open(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 4096, 4096)
	start := time.Now().UnixMilli()
	for {
		_, err = f.Read(buf)
		if err != nil {
			if err == io.EOF {
				b.Log(time.Now().UnixMilli() - start)
				break
			}
			b.Fatal(err)
		}
	}
}

func BenchmarkBufferReadLargest(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	f, err := os.Open(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 4096, 4096)
	r := bufio.NewReaderSize(f, 4096)
	start := time.Now().UnixMilli()
	for {
		_, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				b.Log(time.Now().UnixMilli() - start)
				break
			}
			b.Fatal(err)
		}
	}
}
