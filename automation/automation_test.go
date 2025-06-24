package automation

import (
	"context"
	"database/sql"
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/db/id"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/waapi"

	_ "github.com/mattn/go-sqlite3"
)

var TestBankDir string = "../tests/bnks"
var TestCSVOkDir string = "./tests/csvs/ok"
var TestCSVCompleteDir string = "./tests/csvs/ok/complete"
var TestCSVFailDir string = "./tests/csvs/fail"

func TestParseCSVHeaderOk(t *testing.T) {
	tests, err := os.ReadDir(TestCSVOkDir)
	if err != nil {
		t.Fatal(err)
	}
	var p string = ""
	var b string = ""
	for _, test := range tests {
		p = filepath.Join(TestCSVOkDir, test.Name())
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		b = filepath.Base(p)
		reader := csv.NewReader(f)
		header := CSVHeader{Workspace: b, Output: b, Type: 0}
		if err = ParseCSVHeader(&header, reader); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseCSVHeaderFail(t *testing.T) {
	tests, err := os.ReadDir(TestCSVFailDir)
	if err != nil {
		t.Fatal(err)
	}
	var p string = ""
	var b string = ""
	for _, test := range tests {
		p = filepath.Join(TestCSVFailDir, test.Name())
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		b = filepath.Base(p)
		reader := csv.NewReader(f)
		header := CSVHeader{Workspace: b, Output: b, Type: 0}
		if err = ParseCSVHeader(&header, reader); err == nil {
			t.Fatalf("Expecting test case %s to fail", p)
		} else {
			t.Log(err)
		}
	}
}

type RewireSoundsWithNewSourcesCSVTest struct {
	CSV string
	Bank string
}

func TestRewireSoundsWithNewSourceCSV(t *testing.T) {
	waapi.InitTmp()
	tests := []RewireSoundsWithNewSourcesCSVTest{
		{
			"squad_ak74_ar19_01.csv",
			"wep_ar19_liberator.bnk",
		},
		{
			"squad_ak74_ar19_02.csv",
			"wep_ar19_liberator.bnk",
		},
		{
			"squad_ak74_ar19_03.csv",
			"wep_ar19_liberator.bnk",
		},
		{
			"squad_ak74_ar19_04.csv",
			"wep_ar19_liberator.bnk",
		},
		{
			"squad_ak74_ar19_05.csv",
			"wep_ar19_liberator.bnk",
		},
	}
	for _, test := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()
		bnk, err := parser.ParseBank(
			filepath.Join(TestBankDir, test.Bank),
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
		if err := RewireSoundsWithNewSourcesCSV(
			ctx,
			h, 
			d,
			filepath.Join(TestCSVCompleteDir, test.CSV),
			"VORBIS High Quality",
			"None",
			true,
		); err != nil {
			t.Fatal(err)
		}
	}
	waapi.CleanTmp()
}

func BenchmarkIDCollisionCheck(b *testing.B) {
	db, err := sql.Open("sqlite3", "../id_15314")
	if err != nil {
		b.Fatal(err)
	}
	q := id.New(db)
	ctx := context.Background()
	var in uint64
	for b.Loop() {
		_, err := TrySid(ctx, q)
		if err != nil {
			in += 1
		}
	}
	b.Logf("Collision: %d.\n", in)
}
