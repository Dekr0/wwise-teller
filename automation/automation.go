package automation

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dekr0/wwise-teller/wwise"
)

func RewireSoundsWithNewSourcesCSV(
	ctx context.Context,
	h *wwise.HIRC,
	mappingFile string,
) error {
	f, err := os.Open(mappingFile)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to obtain workspace directory: %w", err)
	}
	if strings.Compare(header[0], "workspace") != 0 {
		return fmt.Errorf("Expecting Row 1, Column 1 to be workspace")
	}
	workspace := filepath.Dir(mappingFile)
	if len(header) >= 2 {
		overwrite := header[1]
		// Relative path
		if !filepath.IsAbs(overwrite) {
			overwrite = filepath.Join(workspace, overwrite)
		}
		_, err := os.Lstat(overwrite)
		if err != nil {
			if os.IsNotExist(err) {
				slog.Info(fmt.Sprintf("%s does not exist. Creating this workspace", overwrite))
				if err := os.MkdirAll(overwrite, 0777); err != nil {
					slog.Error(fmt.Sprintf("Failed to create workspace %s", overwrite), "error", err)
					slog.Info(fmt.Sprintf("Fallback to workspace %s", workspace))
				} else {
					workspace = overwrite
				}
			} else {
				slog.Error(fmt.Sprintf("Failed to obtain information about workspace %s", workspace), "error", err)
				slog.Info(fmt.Sprintf("Fallback to workspace %s", workspace))
			}
		}
	}

	mapping := make(map[string][]*wwise.Sound, 32)
	rowNum := 2
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if len(row) < 2 {
			slog.Error(fmt.Sprintf("Expecting two columns, filename and # of targeting sound IDs, at row %d", rowNum))
			rowNum += 1
			continue
		}

		// Extension
		input := row[0]
		ext := filepath.Ext(input)
		if ext != ".wav" {
			slog.Error("Wave file is the only supported file format.")
			rowNum += 1
			continue
		}

		// Existence
		if !filepath.IsAbs(input) {
			input = filepath.Join(workspace, input)
		}
		_, err = os.Lstat(input)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to obtain info of wave file %s", input), "error", err)
			rowNum += 1
			continue
		}

		// Counter
		count, err := strconv.ParseUint(row[1], 10, 8)
		if err != nil {
			slog.Error("Failed to parse # of targeting sound IDs", "error", err)
			continue
		}
		if count > uint64(len(row) - 2) {
			slog.Error(fmt.Sprintf("Expecting %d of targeting sound IDs but receiving %d", count, len(row)- 2))
			continue
		}

		if _, in := mapping[row[0]]; !in {
			// IDing
			columnNum := 3
			targets := make([]*wwise.Sound, 0, count)
			for _, i := range row[2:] {
				id, err := strconv.ParseUint(i, 10, 32)
				if err != nil {
					slog.Error(fmt.Sprintf("Failed to parse sound ID at column %d", columnNum), "error", err)
					columnNum += 1
					continue
				}
				if v, ok := h.ActorMixerHirc.Load(uint32(id)); !ok {
					slog.Error(fmt.Sprintf("ID %d does not have an associated sound object", id))
				} else {
					// no duplication check because it will map to the same sid 
					// at the end
					targets = append(targets, v.(*wwise.Sound))
				}
				columnNum += 1
			}
			mapping[row[0]] = targets
		} else {
			slog.Error(fmt.Sprintf("Duplicated input file %s at row %d, column 1", row[0], rowNum))
		}
		rowNum += 1
	}

	if len(mapping) == 0 {
		return fmt.Errorf("No mapping is provided. Aborted")
	}
	return nil
}
