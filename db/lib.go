package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/Dekr0/wwise-teller/db/id"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/cenkalti/backoff"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const DatabaseEnv = "IDATABASE"

var DatabaseEnvNotSet error = fmt.Errorf("Enviromental variable %s is not set.", DatabaseEnv)
var DatabaseEnvNotAbs error = fmt.Errorf("Enviromental variable %s is not in absolute path.", DatabaseEnv)
var WriteLock sync.Mutex

func CheckDatabaseEnv() error {
	p := os.Getenv(DatabaseEnv)
	if p == "" {
		return DatabaseEnvNotSet
	}
	if !filepath.IsAbs(p) {
		return DatabaseEnvNotAbs
	}
	db, err := CreateDefaultConn()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("%s is not a valid ID database because of database connection error (%w).", p, err)
	}
	return nil
}

// Create a database connection using the default environmental variable
func CreateDefaultConn() (*sql.DB, error) {
	p := os.Getenv(DatabaseEnv)
	if p == "" {
		return nil, DatabaseEnvNotSet
	}
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database %s: %w", p, err)
	}
	return db, nil
}

func CreateDefaultConnWithQuery() (*id.Queries, func(), error) {
	db, err := CreateDefaultConn()
	if err != nil {
		return nil, nil, err
	}
	closeDb := func() { 
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}
	return id.New(db), closeDb, nil
}

func CreateDefaultConnWithTxQuery(ctx context.Context) (*id.Queries, func(), func() error, func(), error) {
	db, err := CreateDefaultConn()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	closeDb := func() { 
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}
	commit := func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("Failed to commit database transaction: %w", err)
		}
		return nil
	}
	rollback := func() {
		if err := tx.Rollback(); err != nil {
			slog.Error("Failed to rollback database transaction. Please manually rollback database by using the backup database", "error", err)
		}
	}
	return id.New(db).WithTx(tx), closeDb, commit, rollback, nil
}

func TrySid(ctx context.Context, q *id.Queries) (uint32, error) {
	b := backoff.WithContext(backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 16), ctx)
	var sid uint32 = 0
	if err := backoff.Retry(func() error {
		var err error
		sid, err = utils.ShortID()
		if err != nil {
			slog.Error("Failed to generate 32 bit unsigned integer ID", "error", err)
			return err
		}
		count, err := q.SourceId(ctx, int64(sid))
		if err != nil {
			slog.Error("Failed to query source ID from database", "error", err)
			return err
		}
		if count > 0 {
			err := fmt.Errorf("Source ID %d already exists.", sid)
			slog.Error(err.Error())
			return err
		}
		return nil
	}, b); err != nil {
		return 0, err
	}
	if sid == 0 {
		return 0, errors.New("Source ID uses invalid value of 0.")
	}
	if err := q.InsertSource(ctx, int64(sid)); err != nil {
		return 0, err
	}
	return sid, nil
}

func TryHid(ctx context.Context, q *id.Queries) (uint32, error) {
	b := backoff.WithContext(backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 16), ctx)
	var hid uint32 = 0
	if err := backoff.Retry(func() error {
		var err error
		hid, err = utils.ShortID()
		if err != nil {
			slog.Error("Failed to generate 32 bit unsigned integer ID", "error", err)
			return err
		}
		count, err := q.HierarchyId(ctx, int64(hid))
		if err != nil {
			slog.Error("Failed to query source ID from database", "error", err)
			return err
		}
		if count > 0 {
			err := fmt.Errorf("Source ID %d already exists.", hid)
			slog.Error(err.Error())
			return err
		}
		return nil
	}, b); err != nil {
		return 0, err
	}
	if hid == 0 {
		return 0, errors.New("Source ID uses invalid value of 0.")
	}
	if err := q.InsertHierarchy(ctx, int64(hid)); err != nil {
		return 0, err
	}
	return hid, nil
}
