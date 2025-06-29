package ui

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
)

func selectHD2ArchiveFunc() func([]string) {
	return func(paths []string) {
		pushExtractSoundBanksModal(paths)
	}
}

func extractHD2SoundBanksFunc(paths []string) func(string) {
	return func(dest string) {
		for _, path := range paths {
			dispatchHD2ExtractSoundBank(path, dest)
		}
	}
}

func dispatchHD2ExtractSoundBank(path string, dest string) {
	timeout, cancel := context.WithTimeout(
		context.Background(), time.Second * 8,
	)
	onProcMsg := fmt.Sprintf(
		"Extract sound banks from Helldivers 2 game archive %s", path,
	)
	onDoneMsg := fmt.Sprintf(
		"Extracted sound banks from Helldivers 2 game archive %s", path,
	)
	f := func(ctx context.Context) {
		if err := helldivers.ExtractSoundBankStable(path, dest, false); err != nil {
			slog.Error(
				fmt.Sprintf(
					"Failed to extract sound bank from game archive %s", path,
				),
				"error", err,
			)
		}
	}
	err := GlobalCtx.Loop.QTask(timeout, cancel, onProcMsg, onDoneMsg, f)
	if err != nil {
		slog.Error(
			fmt.Sprintf(
				"Failed to extract sound bank from Helldivers 2 game " + 
				"archive %s", path,
			),
			"error", err,
		)
	}
}
