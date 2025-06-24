package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/Dekr0/wwise-teller/db/id"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var DatabaseEnvNotSet error = errors.New("Enviroment variable IDATABASE is not set.")

// Create a database connection using the default environmental variable
func CreateDefaultConn() (*sql.DB, error) {
	p := os.Getenv("IDATABASE")
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
			slog.Error("Failed to rollback database transaction", "error", err)
		}
	}
	return id.New(db).WithTx(tx), closeDb, commit, rollback, nil
}
