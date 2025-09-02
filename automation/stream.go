package automation

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/Dekr0/wwise-teller/wwise"
)

// This only address problem related to prefetch stream.
func ToStreamTypeBnk(bnk *wwise.Bank, script string) (err error) {
	if bnk.DIDX() == nil {
		return wwise.NoDIDX
	}
	if bnk.DATA() == nil {
		return wwise.NoDATA
	}
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}
	di := bnk.DIDX()
	da := bnk.DATA()

	f, err := os.Open(script)
	if err != nil {
		return err
	}

	reader := csv.NewReader(f)
	rowNum := 0
	var row []string
	for {
		row, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		soundId, err := strconv.ParseUint(row[0], 10, 32)
		if err != nil {
			slog.Error(fmt.Sprintf("Invalid sound id %s at row %d", row[0], rowNum))
			rowNum += 1
			continue
		}
		v, ok := h.ActorMixerHirc.Load(uint32(soundId))
		if !ok {
			slog.Error(fmt.Sprintf("Sound id %d does not have an associated sound object at row %d", soundId, rowNum))
			rowNum += 1
			continue
		}
		sound := v.(*wwise.Sound)
		data := &sound.BankSourceData
		if data.StreamType == wwise.SourceTypeDATA {
			rowNum += 1
			continue
		}
		data.StreamType = wwise.SourceTypeDATA
		sid := data.SourceID
		if _, in := di.MediaIndexsMap[sid]; !in {
			slog.Error(fmt.Sprintf("Source ID %d does not have media index entry at row %d", sid, rowNum))
			rowNum += 1
			continue
		}
		if _, in := da.AudiosMap[sid]; !in {
			slog.Error(fmt.Sprintf("Source ID %d does not have audio source data at row %d", sid, rowNum))
		}
	}
	return nil
}
