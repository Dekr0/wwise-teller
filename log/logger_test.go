package log

import (
	"log/slog"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	l := slog.New(NewHandler(os.Stdout, nil))

	l.Info("Test 1", "message", "1")
	l.Warn("Test 2", "message", "2")
	l.Error("Test 3", "message", "3")
}
