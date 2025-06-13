package main

import (
	"log/slog"

	"github.com/Dekr0/wwise-teller/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		slog.Error("Failed to launch GUI", "error", err)
	}
}
