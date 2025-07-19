package automation

import (
	"os"
	"path/filepath"
)

var testDir string = os.Getenv("TEST")
var testScriptsDir string = filepath.Join(testDir, "automations")
var testStBankDir string = filepath.Join(testDir, "default_st_bnks")
