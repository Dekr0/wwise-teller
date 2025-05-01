package ui

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/ui/async"
)

func selectGameArchiveFunc(
	modalQ *ModalQ,
	loop *async.EventLoop,
	conf *config.Config,
) func([]string) {
	return func(paths []string) {
		pushExtractSoundBanksModal(modalQ, loop, conf, paths)
	}
}

func extractSoundBanksFunc(loop *async.EventLoop, paths []string) func(string) {
	return func(dest string) {
		for _, path := range paths {
			dispatchExtractSoundBank(loop, path, dest)
		}
	}
}

func dispatchExtractSoundBank(loop *async.EventLoop, path string, dest string) {
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
		if err := helldivers.ExtractSoundBank(ctx, path, dest); err != nil {
			slog.Error(
				fmt.Sprintf(
					"Failed to extract sound bank from game archive %s", path,
				),
				"error", err,
			)
		}
	}
	err := loop.QTask(timeout, cancel, onProcMsg, onDoneMsg, f)
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
