package automation

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/wwise-teller/db/id"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var testBankDir string = "../tests/st_bnks"
var testProcessSpecDirOk string = "./tests/process/ok"
var testProcessSpecDirFail string = "./tests/process/fail"

func setDatabase() {
	p, _ := filepath.Abs("../id_15314")
	os.Setenv("IDATABASE", p)
}

func TestParseProcessSpec(t *testing.T) {

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
