package automation

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wwise"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var testRewireOkDir string = "./tests/csvs/ok"
var testRewireCompleteDir string = "./tests/csvs/ok/complete"
var testRewirePartialDir string = "./tests/csvs/ok/partial"
var testRewireFailDir string = "./tests/csvs/fail"

func TestParseRewireHeaderOk(t *testing.T) {
	tests, err := os.ReadDir(testRewireOkDir)
	if err != nil {
		t.Fatal(err)
	}
	var p string = ""
	var b string = ""
	for _, test := range tests {
		t.Logf("---Running test %s---", test.Name())
		if test.IsDir() {
			continue
		}
		p = filepath.Join(testRewireOkDir, test.Name())
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		b = filepath.Base(p)
		reader := csv.NewReader(f)
		header := CSVHeader{Workspace: b, Type: 0}
		if err = ParseRewireHeader(&header, reader); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseRewireHeaderFail(t *testing.T) {
	tests, err := os.ReadDir(testRewireFailDir)
	if err != nil {
		t.Fatal(err)
	}
	var p string = ""
	var b string = ""
	for _, test := range tests {
		t.Logf("---Running test %s---", test.Name())
		if test.IsDir() {
			continue
		}
		p = filepath.Join(testRewireFailDir, test.Name())
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		b = filepath.Base(p)
		reader := csv.NewReader(f)
		header := CSVHeader{Workspace: b, Type: 0}
		if err = ParseRewireHeader(&header, reader); err == nil {
			t.Fatalf("Expecting test case %s to fail", p)
		} else {
			t.Log(err)
		}
	}
}

type RewireWithNewSourcesTest struct {
	CSV string
	Bank string
}

func TestRewireWithNewSources(t *testing.T) {
	setDatabase()
	waapi.InitTmp()
	tests := []RewireWithNewSourcesTest{
		{
			"squad_ak74_ar19_01.csv",
			"wep_ar19_liberator.st_bnk",
		},
		{
			"squad_ak74_ar19_02.csv",
			"wep_ar19_liberator.st_bnk",
		},
	}
	for _, test := range tests {
		t.Logf("---Running test %s---", test.CSV)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		bnk, err := parser.ParseBank(
			filepath.Join(testBankDir, test.Bank),
			ctx,
		)
		if err != nil {
			t.Fatal(err)
		}

		h, d := bnk.HIRC(), bnk.DIDX()
		if h == nil {
			t.Fatal("HIRC chunk is nil")
		}
		if d == nil {
			t.Fatal("DIDX chunk is nil")
		}

		if err := RewireWithNewSources(
			ctx,
			bnk,
			filepath.Join(testRewireCompleteDir, test.CSV),
			false,
		); err != nil {
			if err == wwise.NoHIRC {
				t.Fatalf("%s does not have HIRC chunk", test.Bank)
			}
			t.Fatal(err)
		}
	}
	waapi.CleanTmp()
}

func TestRewireWithNewSourcesPartial(t *testing.T) {
	setDatabase()
	waapi.InitTmp()
	tests := []RewireWithNewSourcesTest{
		{
			"squad_ak74_ar19_01.csv",
			"wep_ar19_liberator.st_bnk",
		},
		{
			"squad_ak74_ar19_02.csv",
			"wep_ar19_liberator.st_bnk",
		},
	}
	for _, test := range tests {
		t.Logf("---Running test %s---", test.CSV)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		bnk, err := parser.ParseBank(
			filepath.Join(testBankDir, test.Bank),
			ctx,
		)
		if err != nil {
			t.Fatal(err)
		}

		h, d := bnk.HIRC(), bnk.DIDX()
		if h == nil {
			t.Fatal("HIRC chunk is nil")
		}
		if d == nil {
			t.Fatal("DIDX chunk is nil")
		}

		if err := RewireWithNewSources(
			ctx,
			bnk,
			filepath.Join(testRewirePartialDir, test.CSV),
			false,
		); err != nil {
			if err == wwise.NoHIRC {
				t.Fatalf("%s does not have HIRC chunk", test.Bank)
			}
			t.Fatal(err)
		}
	}
	waapi.CleanTmp()
}
