package wm_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGenerateDockTag(t *testing.T) {
	names := []string{
		"Attenuations",
		"Actor Mixer Hierarchy",
		"Bank Explorer",
		"Debug",
		"File Explorer",
		"Log",
		"Master Mixer Hierarchy",
		"Music Hierarchy",
		"Notification",
		"Transport Control",
	}

	var b strings.Builder
	b.WriteString("package wm\n\nconst (\n")

	for i, name := range names {
		b.WriteString(fmt.Sprintf("    DockTag%s", strings.ReplaceAll(name, " ", "")))
		if i == 0 {
			b.WriteString(" = iota")
		}
		b.WriteByte('\n')
	}
	b.WriteString("    DockTagCount\n)\n")
	if err := os.WriteFile("enum.go", []byte(b.String()), 0666); err != nil {
		t.Fatal(err)
	}
}
