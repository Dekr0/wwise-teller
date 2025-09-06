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
var BigFile string = os.Getenv("BIGFILE")

func benchmarkReadOnceBigFile(b *testing.B) {
	start := time.Now().UnixMilli()
	data, err := os.ReadFile(BigFile)
	if err != nil {
		b.Fatal(err)
	}
	b.Log(time.Now().UnixMilli() - start)
	for i := range data {
		data[i] = 0
	}
}

func BenchmarkBufferReadBigFile(b *testing.B) {
	f, err := os.Open(BigFile)
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

func benchmarkReadOnceReadLargest(b *testing.B) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	start := time.Now().UnixMilli()
	data, err := os.ReadFile(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		b.Fatal(err)
	}
	b.Log(time.Now().UnixMilli() - start)
	for i := range data {
		data[i] = 0
	}
}

func benchmarkUnbufferReadLargest(b *testing.B) {
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

func benchmarkBufferReadLargest(b *testing.B) {
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
