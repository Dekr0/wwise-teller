package main

import (
	"log/slog"
	"runtime"

	"github.com/Dekr0/wwise-teller/ui"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := ui.Run(); err != nil {
		slog.Error("Failed to launch GUI", "error", err)
	}
}
