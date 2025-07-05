package automation

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/db"
	"github.com/Dekr0/wwise-teller/db/id"
	"github.com/Dekr0/wwise-teller/utils"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var testStBankDir = filepath.Join(os.Getenv("TESTS"), "default_st_bnks")
var testProcessDir = filepath.Join(os.Getenv("TESTS"), "process")
var testProcessOkDir = filepath.Join(testProcessDir, "ok")

func TestProcess(t *testing.T) {
	if err := utils.InitTmp(); err != nil {
		t.Fatal(err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	entries, err := os.ReadDir(testProcessOkDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		Process(context.Background(), filepath.Join(testProcessOkDir, entry.Name()))
	}
	utils.CleanTmp()
}

func BenchmarkIDCollisionCheck(b *testing.B) {
	conn, err := sql.Open("sqlite3", "../id_15314")
	if err != nil {
		b.Fatal(err)
	}
	q := id.New(conn)
	ctx := context.Background()
	var in uint64
	for b.Loop() {
		_, err := db.TrySid(ctx, q)
		if err != nil {
			in += 1
		}
	}
	b.Logf("Collision: %d.\n", in)
}
