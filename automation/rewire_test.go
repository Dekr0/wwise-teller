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

var testRewire string = filepath.Join(os.Getenv("TESTS"), "rewire")
var testRewireNew string = filepath.Join(testRewire, "new")
var testRewireNewOkDir string = filepath.Join(testRewire, "ok")
var testRewireNewCompleteDir string = filepath.Join(testRewireNewOkDir, "complete")
var testRewireNewPartialDir string = filepath.Join(testRewireNewOkDir, "partial")
var testRewireNewFailDir string = filepath.Join(testRewire, "fail")

func TestParseRewireHeaderOk(t *testing.T) {
	tests, err := os.ReadDir(testRewireNewOkDir)
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
		p = filepath.Join(testRewireNewOkDir, test.Name())
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
	tests, err := os.ReadDir(testRewireNewFailDir)
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
		p = filepath.Join(testRewireNewFailDir, test.Name())
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
			filepath.Join(testStBankDir, test.Bank),
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
			filepath.Join(testRewireNewCompleteDir, test.CSV),
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
			filepath.Join(testStBankDir, test.Bank),
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
			filepath.Join(testRewireNewPartialDir, test.CSV),
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
