package ui

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
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
	saveName string,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 4,
		)

		onProcMsg := fmt.Sprintf("Saving sound bank %s to %s", saveName, path)
		onDoneMsg := fmt.Sprintf("Saved sound bank %s to %s", saveName, path)

		if err := loop.QTask(timeout, cancel,
			onProcMsg, onDoneMsg,
			func (ctx context.Context) {
				slog.Info(onProcMsg)

				bnkMngr.writeLock.Store(true)

				if data, err := saveTab.encode(ctx); err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", saveName),
						"error", err,
					)
				} else {
					if err := utils.SaveFileWithRetry(
						data, 
						filepath.Join(path, filepath.Base(saveName)),
					); err != nil {
						slog.Error(
							fmt.Sprintf("Failed to save sound bank %s to %s",
								saveName, path),
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

func HD2PatchFunc(
	loop *async.EventLoop,
	bnkMngr *BankManager,
	saveTab *bankTab,
	saveName string,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 4,
		)
		
		onProcMsg := fmt.Sprintf("Saving sound bank %s to HD2 patch %s",
			saveName, path)
		onDoneMsg := fmt.Sprintf("Saved sound bank %s to HD2 patch %s",
			saveName, path)

		if err := loop.QTask(timeout, cancel, onProcMsg, onDoneMsg, 
			func(ctx context.Context) {
				slog.Info(onProcMsg)
				bnkMngr.writeLock.Store(true)
				defer bnkMngr.writeLock.Store(false)
				bnkData, err := saveTab.encode(ctx)
				if err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", saveName),
						"error", err,
					)
					return
				}
				meta := saveTab.bank.META()
				if meta == nil {
					slog.Error(
						fmt.Sprintf("Sound bank %s is missing integration data.",
							saveName),
					)
					return
				}
				if err := helldivers.GenHelldiversPatch(
					ctx, bnkData, meta.Data, path,
				); err != nil {
					slog.Error(fmt.Sprintf("Failed to write HD2 patch to %s", path))
				} else {
					slog.Info(onDoneMsg)
				}
			},
		);
		   err != nil {
			slog.Error(fmt.Sprintf("Failed to save HD2 patch to %s", path),
				"error", err,
			)
		}
	}
}
