package db

import (
	"context"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestConcurrent(t *testing.T) {
	if err := InitDatabase(); err != nil {
		t.Fatal(err)
	}

	var w sync.WaitGroup

	bg := context.Background()

	w.Add(1)
	go func() {
		defer w.Done()
		sids := make([]uint32, 16, 16)
		closeConn, commit, rollback, err := AllocateSids(bg, sids)
		if err != nil {
			t.Log(err)
			return
		}
		defer closeConn()
		if err := commit(); err != nil {
			t.Log(err)
			rollback()
			return
		}
		t.Log(sids)
	}()

	w.Add(1)
	go func() {
		defer w.Done()
		hids := make([]uint32, 16, 16)
		closeConn, commit, rollback, err := AllocateHids(bg, hids)
		if err != nil {
			t.Log(err)
			return
		}
		defer closeConn()
		if err := commit(); err != nil {
			t.Log(err)
			rollback()
			return
		}
		t.Log(hids)
	}()

	w.Wait()
}
