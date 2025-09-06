package decoder_test

import (
	"encoding/binary"
	"io"
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

func TestRawReadAll(t *testing.T) {
	entries, err := os.ReadDir(SoundBanksDir)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 4096, 4096)
	for _, entry := range entries {
		f, err := os.Open(filepath.Join(SoundBanksDir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		for {
			_, err = f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}
		}
	}
}

func TestRawReadLargest(t *testing.T) {
	const bank = "content_audio_weapons_superearth.st_bnk"
	f, err := os.Open(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	buf := make([]byte, 4096, 4096)
	for {
		_, err = f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
	}
}
