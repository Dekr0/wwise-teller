package automation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Dekr0/wwise-teller/db/id"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/cenkalti/backoff"
)

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
