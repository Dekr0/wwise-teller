package aio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/utils"
)

var testWEMsDir string = filepath.Join(os.Getenv("TESTS"), "wems")

func TestGetDurationWEMFile(t *testing.T) {
	utils.InitTmp()
	defer utils.CleanTmp()
	entries, err := os.ReadDir(testWEMsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
	}
}
