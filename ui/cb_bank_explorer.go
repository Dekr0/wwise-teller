package ui

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/utils"
)

func openSoundBankFunc(
	bnkMngr *be.BankManager,
) func([]string) {
	return func(paths []string) {
		for _, path := range paths {
			dispatchOpenSoundBank(path, bnkMngr)
		}
	}
}

func dispatchOpenSoundBank(
	path string,
	bnkMngr *be.BankManager,
) {
	timeout, cancel := context.WithTimeout(
		context.Background(), time.Second * 30,
	)
	base := filepath.Base(path)
	onProcMsg := fmt.Sprintf("Loading sound bank %s", base)
	onDoneMsg := fmt.Sprintf("Loaded sound bank %s", base)
	if err := GlobalCtx.Loop.QTask(timeout, cancel,
		onProcMsg, onDoneMsg,
		func (ctx context.Context) {
			slog.Info(onProcMsg)
			err := bnkMngr.OpenBank(ctx, path)
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
	bnkMngr *be.BankManager,
	saveTab *be.BankTab,
	saveName string,
	excludeMeta bool,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 4,
		)

		onProcMsg := fmt.Sprintf("Saving sound bank %s to %s", saveName, path)
		onDoneMsg := fmt.Sprintf("Saved sound bank %s to %s", saveName, path)

		if err := GlobalCtx.Loop.QTask(timeout, cancel,
			onProcMsg, onDoneMsg,
			func (ctx context.Context) {
				slog.Info(onProcMsg)

				bnkMngr.WriteLock.Store(true)

				if data, err := saveTab.Encode(ctx, excludeMeta); err != nil {
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
				bnkMngr.WriteLock.Store(false)
			},
		); err != nil {
			slog.Error(fmt.Sprintf("Failed to save sound bank to %s", path),
				"error", err,
			)
		}
	}
}

func HD2PatchFunc(
	bnkMngr *be.BankManager,
	saveTab *be.BankTab,
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

		if err := GlobalCtx.Loop.QTask(timeout, cancel, onProcMsg, onDoneMsg, 
			func(ctx context.Context) {
				slog.Info(onProcMsg)
				bnkMngr.WriteLock.Store(true)
				defer bnkMngr.WriteLock.Store(false)
				bnkData, err := saveTab.Encode(ctx, true)
				if err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", saveName),
						"error", err,
					)
					return
				}
				meta := saveTab.Bank.META()
				if meta == nil {
					slog.Error(
						fmt.Sprintf("Sound bank %s is missing integration data.",
							saveName),
					)
					return
				}
				if err := helldivers.GenHelldiversPatchStable(
					bnkData, meta.B, path,
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
