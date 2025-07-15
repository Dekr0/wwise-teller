package scripts

import (
	"context"
	"os"
	"testing"
)

func TestExtractHD2(t *testing.T) {
	if err := ExtractHD2(context.Background(), os.Getenv("DATA"), "./output"); err != nil {
		t.Fatal(err)
	}
}
