package ui

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/Dekr0/wwise-teller/utils"
)

func openSoundBankFunc(
	loop *async.EventLoop,
	bnkMngr *BankManager,
) func([]string) {
	return func(paths []string) {
		for _, path := range paths {
			dispatchOpenSoundBank(path, bnkMngr, loop)
		}
	}
}

func dispatchOpenSoundBank(
	path string,
	bnkMngr *BankManager,
	loop *async.EventLoop,
) {
	timeout, cancel := context.WithTimeout(
		context.Background(), time.Second * 30,
	)
	base := filepath.Base(path)
	onProcMsg := fmt.Sprintf("Loading sound bank %s", base)
	onDoneMsg := fmt.Sprintf("Loaded sound bank %s", base)
	if err := loop.QTask(timeout, cancel,
		onProcMsg, onDoneMsg,
		func (ctx context.Context) {
			slog.Info(onProcMsg)
			err := bnkMngr.openBank(ctx, path)
			if err != nil {
				slog.Error(
					fmt.Sprintf("Failed to load sound bank %s", base),
					"error", err,
				)
			} else {
				slog.Info(onDoneMsg)
			}
		},
	); err != nil {
		slog.Error("Failed to open sound bank", "error", err)
		cancel()
	}
}

func saveSoundBankFunc(
	loop *async.EventLoop,
	bnkMngr *BankManager,
	saveTab *bankTab,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 30,
		)

		onProcMsg := fmt.Sprintf("Saving sound bank to %s", path)
		onDoneMsg := fmt.Sprintf("Saved sound bank to %s", path)

		if err := loop.QTask(timeout, cancel,
			onProcMsg, onDoneMsg,
			func (ctx context.Context) {
				slog.Info(onProcMsg)

				bnkMngr.writeLock.Store(true)

				if data, err := saveTab.encode(ctx); err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", filepath.Base(path)),
						"error", err,
					)
				} else {
					if err := utils.SaveFileWithRetry(data, path); err != nil {
						slog.Error(
							fmt.Sprintf("Failed to save sound bank to %s", path),
							"error", err,
						)
					} else {
						slog.Info(onDoneMsg)
					}
				}

				bnkMngr.writeLock.Store(false)
			},
		); err != nil {
			slog.Error(fmt.Sprintf("Failed to save sound bank to %s", path),
				"error", err,
			)
		}
	}
}
