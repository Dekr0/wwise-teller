package main

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/ui"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	logger := slog.New(log.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := ui.Run(); err != nil {
		logger.Error("Failed to launch GUI", "error", err)
	}
}
