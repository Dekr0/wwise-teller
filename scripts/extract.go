// TODO parallelization
package scripts

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
)

func ExtractHD2(ctx context.Context, data string, output string) error {
	f, err := os.Open(data)
	if err != nil {
		return err
	}

	stat, err := os.Lstat(output)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(output, 0777); err != nil {
				return fmt.Errorf("Failed to create output directory %s: %w", output, err)
			}
		} else {
			return fmt.Errorf("Failed to obtain information about output directory %s: %w", output, err)
		}
	} else {
		if !stat.IsDir() {
			return fmt.Errorf("%s is not a directory.", output)
		}
	}

	for {
		entries, err := f.ReadDir(1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		select {
		case <- ctx.Done():
			return ctx.Err()
		default:
		}
		for _, entry := range entries {
			if entry.IsDir() { continue }

			ext := filepath.Ext(entry.Name())
			if strings.EqualFold(ext, ".stream") { continue }
			if strings.EqualFold(ext, ".gpu_resources") { continue }
			if strings.EqualFold(ext, ".ini") { continue }
			if strings.EqualFold(ext, ".data") { continue }
			if strings.Contains(ext, "patch") { continue }
			
			select {
			default:
				helldivers.ExtractSoundBankStable(filepath.Join(data, entry.Name()), output, false)
			}
		}
	}
	return nil
}
