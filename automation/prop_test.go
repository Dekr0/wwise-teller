package automation

import (
	"os"
	"path/filepath"
	"testing"
)

var testPropOk = "./tests/jsons/ok"

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
