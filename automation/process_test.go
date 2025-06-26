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
	"github.com/Dekr0/wwise-teller/wwise"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var testBankDir string = "../tests/st_bnks"
var testCSVOkDir string = "./tests/csvs/ok"
var testCSVCompleteDir string = "./tests/csvs/ok/complete"
var testCSVPartialDir string = "./tests/csvs/ok/partial"
var testCSVFailDir string = "./tests/csvs/fail"

func setDatabase() {
	p, _ := filepath.Abs("../id_15314")
	os.Setenv("IDATABASE", p)
}

func TestParseRewireHeaderOk(t *testing.T) {
	tests, err := os.ReadDir(testCSVOkDir)
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
		p = filepath.Join(testCSVOkDir, test.Name())
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
	tests, err := os.ReadDir(testCSVFailDir)
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
		p = filepath.Join(testCSVFailDir, test.Name())
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
			filepath.Join(testCSVCompleteDir, test.CSV),
			false,
		); err != nil {
			if err == wwise.NoHIRC {
				t.Fatalf("%s does not have HIRC chunk", test.Bank)
			}
			t.Fatal(err)
		}

		// bnkData, err := bnk.Encode(ctx)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// meta := bnk.META()
		// if meta == nil {
		// 	t.Fatalf("%s does not have META chunk", test.Bank)
		// }

		// if err := helldivers.GenHelldiversPatchStable(bnkData, meta.B, "."); err != nil {
		// 	t.Fatal(err)
		// }
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
			filepath.Join(testCSVPartialDir, test.CSV),
			false,
		); err != nil {
			if err == wwise.NoHIRC {
				t.Fatalf("%s does not have HIRC chunk", test.Bank)
			}
			t.Fatal(err)
		}

		// bnkData, err := bnk.Encode(ctx)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// meta := bnk.META()
		// if meta == nil {
		// 	t.Fatalf("%s does not have META chunk", test.Bank)
		// }

		// if err := helldivers.GenHelldiversPatchStable(bnkData, meta.B, "."); err != nil {
		// 	t.Fatal(err)
		// }
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
