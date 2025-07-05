package aio

import (
	"os"
	"testing"

)

var testWavesDir string = os.Getenv("TESTS")

func TestPlayerOpenResume(t *testing.T) {
	entries, err := os.ReadDir(testWavesDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
	}
}
