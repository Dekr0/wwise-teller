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

var DatabaseEnvNotSet   error = fmt.Errorf("Enviromental variable %s is not set.", DatabaseEnv)
var DatabaseEnvNotAbs   error = fmt.Errorf("Enviromental variable %s is not in absolute path.", DatabaseEnv)
var DatabaseInitRequire error = errors.New("Wwise sound bank ID database is yet opened and initialized")
var WriteLock sync.Mutex

func Ping() error {
	if WwiseIdDB == nil {
		return DatabaseInitRequire
	}
	err := WwiseIdDB.Ping()
	if err != nil {
		return fmt.Errorf("Failed to verify Connectivity of Wwise sound bank ID database: %w", err) 
	}
	return nil
}

var WwiseIdDB *sql.DB
// Open the database using the default environemntal variable
func InitDatabase() (err error) {
	p := os.Getenv(DatabaseEnv)
	if p == "" {
		return DatabaseEnvNotSet
	}
	if !filepath.IsAbs(p) {
		return DatabaseEnvNotAbs
	}
	WwiseIdDB, err = sql.Open("sqlite3", p)
	if err != nil {
		return fmt.Errorf("Failed to open database %s: %w", p, err)
	}
	slog.Info("Opened Wwise sound bank ID database")
	err = WwiseIdDB.Ping()
	if err != nil {
		return fmt.Errorf("Failed to verify Wwise sound bank ID database: %w", err) 
	}
	return nil
}

func CloseDatabase() {
	if WwiseIdDB != nil {
		WwiseIdDB.Close()
	}
}

// Create a database connection using the default environmental variable
func createConn(ctx context.Context) (conn *sql.Conn, err error) {
	if err = Ping(); err != nil {
		return conn, err
	}
	conn, err = WwiseIdDB.Conn(ctx)
	return conn, fmt.Errorf("Failed to connect to Wwise sound bank ID database: %w", err)
}

func createConnWithQuery(ctx context.Context) (*id.Queries, func(), error) {
	conn, err := createConn(ctx)
	if err != nil {
		return nil, nil, err
	}
	closeConn := func() { 
		if err := conn.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}
	return id.New(conn), closeConn, nil
}

func CreateConnWithTxQuery(ctx context.Context) (*id.Queries, func(), func() error, func(), error) {
	conn, err := createConn(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Failed begin Wwise sound bank ID database transaction: %w", err)
	}
	closeConn := func() { 
		if err := conn.Close(); err != nil {
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
	return id.New(conn).WithTx(tx), closeConn, commit, rollback, nil
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
		if err := q.InsertSource(ctx, int64(sid)); err != nil {
			return fmt.Errorf("Failed insert the new allocated source ID %d into Wwise sound bank database: %w", sid, err)
		}
		return nil
	}, b); err != nil {
		return 0, fmt.Errorf("Failed to allocate a new source ID after exhausting all retry: %w", err)
	}
	if sid == 0 {
		return 0, errors.New("Source ID uses invalid value of 0.")
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
		if err := q.InsertHierarchy(ctx, int64(hid)); err != nil {
			return fmt.Errorf("Failed insert the new allocated hierarchy ID %d into Wwise sound bank database: %w", hid, err)
		}
		return nil
	}, b); err != nil {
		return 0, fmt.Errorf("Failed to allocate a new hierarchy ID after exhausting all retry: %w", err)
	}
	if hid == 0 {
		return 0, errors.New("Source ID uses invalid value of 0.")
	}
	return hid, nil
}

func AllocateSids(ctx context.Context, ids []uint32) (
	closeConn func(), commit func() error, rollback func(), err error,
) {
	if len(ids) <= 0 {
		return closeConn, commit, rollback, fmt.Errorf("Empty source IDs array is provided")
	}
	var q *id.Queries
	q, closeConn, commit, rollback, err = CreateConnWithTxQuery(ctx)
	if err != nil {
		return closeConn, commit, rollback, err
	}
	for i := range ids {
		sid, err := TrySid(ctx, q)
		if err != nil {
			rollback()
			rollback = nil
			closeConn()
			closeConn = nil
			commit = nil
			return closeConn, commit, rollback, err
		}
		ids[i] = sid
	}
	return closeConn, commit, rollback, nil
}

func AllocateHids(ctx context.Context, ids []uint32) (
	closeConn func(), commit func() error, rollback func(), err error,
) {
	if len(ids) <= 0 {
		return closeConn, commit, rollback, fmt.Errorf("Empty hierarchy IDs array is provided")
	}
	var q *id.Queries
	q, closeConn, commit, rollback, err = CreateConnWithTxQuery(ctx)
	if err != nil {
		return closeConn, commit, rollback, err
	}
	for i := range ids {
		hid, err := TryHid(ctx, q)
		if err != nil {
			rollback()
			rollback = nil
			closeConn()
			closeConn = nil
			commit = nil
			return closeConn, commit, rollback, err
		}
		ids[i] = hid
	}
	return closeConn, commit, rollback, nil
}

func AllocateSounds(ctx context.Context, sids []uint32, hids []uint32) (
	closeConn func(), commit func() error, rollback func(), err error,
) {
	if len(sids) != len(hids) {
		return closeConn, commit, rollback, fmt.Errorf("Length of source IDs array and lenght of hierarchy IDs mismatch")
	}
	if len(sids) == 0 {
		return closeConn, commit, rollback, fmt.Errorf("Empty hierarchy IDs array is provided")
	}
	var q *id.Queries
	q, closeConn, commit, rollback, err = CreateConnWithTxQuery(ctx)
	if err != nil {
		return closeConn, commit, rollback, err
	}
	for i := range sids {
		sid, err := TrySid(ctx, q)
		if err != nil {
			rollback()
			rollback = nil
			closeConn()
			closeConn = nil
			commit = nil
			return closeConn, commit, rollback, err
		}
		sids[i] = sid
		hid, err := TryHid(ctx, q)
		if err != nil {
			rollback()
			rollback = nil
			closeConn()
			closeConn = nil
			commit = nil
			return closeConn, commit, rollback, err
		}
		hids[i] = hid
	}
	return closeConn, commit, rollback, nil
}
