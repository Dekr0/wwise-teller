package automation

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

var testProp = filepath.Join(testScriptsDir, "prop")
var testPropOk = filepath.Join(testProp, "ok")
var testPropInt = filepath.Join(testProcessOkDir, "prop_int")

func TestPropOk(t *testing.T) {
	tests, err := os.ReadDir(testPropOk)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		if test.IsDir() {
			continue
		}
		t.Logf("---Running %s test---", test.Name())
		if spec, err := ParsePropModifierSpec(filepath.Join(testPropOk, test.Name())); err != nil {
			t.Fatal(err)
		} else {
			t.Log(spec)
		}
	}
}

func TestPropInt(t *testing.T) {
	tests, err := os.ReadDir(testPropInt)
	if err != nil {
		t.Fatal(err)
	}
	bg := context.Background()
	for _, test := range tests {
		if test.IsDir() {
			continue
		}
		t.Logf("---Running %s test---", test.Name())
		Process(bg, filepath.Join(testPropInt, test.Name()))
	}
}
