package aio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gopxl/beep/v2/speaker"
)

var testWavesDir string = filepath.Join(os.Getenv("TEST"))

func TestPlayOpen(t *testing.T) {
	player := Player{}
	entries, err := os.ReadDir(testWavesDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err = player.Open(filepath.Join(testWavesDir, entry.Name())); err != nil {
			t.Fatal(err)
		}
		speaker.Play(player.Streamer)
		return
	}
}
