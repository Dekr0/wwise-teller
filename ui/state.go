package ui

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/log"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func createOnOpenCallback(
	loop *async.EventLoop,
	bnkMngr *bankManager,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 30,
		)
		base := filepath.Base(path)
		onProc := fmt.Sprintf("Loading sound bank %s", base)
		onDone := fmt.Sprintf("Loaded sound bank %s", base)
		if err := loop.QTask(timeout, cancel,
			onProc, onDone,
			func (ctx context.Context) {
				slog.Info(onProc)
				err := bnkMngr.openBank(ctx, path)
				if err != nil {
					slog.Error(
						fmt.Sprintf("Failed to load sound bank %s", base),
						"error", err,
					)
				} else {
					slog.Info(onDone)
				}
			},
		); err != nil {
			slog.Error("Failed to open sound bank", "error", err)
			cancel()
		}
	}
}

func createOnSaveCallback(
	loop *async.EventLoop,
	bnkMngr *bankManager,
	saveTab *bankTab,
	saveName string,
) func(string) {
	return func(path string) {
		timeout, cancel := context.WithTimeout(
			context.Background(), time.Second * 30,
		)

		onProc := fmt.Sprintf("Saving sound bank to %s", path)
		onDone := fmt.Sprintf("Saved sound bank to %s", path)

		if err := loop.QTask(timeout, cancel,
			onProc, onDone,
			func (ctx context.Context) {
				slog.Info(onProc)

				bnkMngr.mu.Lock()
				bnkMngr.writeLock = true
				bnkMngr.mu.Unlock()

				if data, err := saveTab.encode(ctx); err != nil {
					slog.Error(
						fmt.Sprintf("Failed to encode sound bank %s", filepath.Base(path)),
						"error", err,
					)
				} else {
					if err := utils.SaveFileWithRetry(data, path); err != nil {
						slog.Error(
							fmt.Sprintf("Failed to save sound bank to %s", path),
							"error", err,
						)
					} else {
						slog.Info(onDone)
					}
				}

				bnkMngr.mu.Lock()
				bnkMngr.writeLock = false
				bnkMngr.mu.Unlock()
			},
		); err != nil {
			slog.Error(fmt.Sprintf("Failed to save sound bank to %s", path),
				"error", err,
			)
		}
	}
}

type bankManager struct {
	banks sync.Map
	writeLock bool
	mu sync.Mutex
}

type bankTab struct {
	bank *wwise.Bank
	idQuery string
	typeQuery int32
	filtered []wwise.HircObj
	lSelStorage imgui.SelectionBasicStorage

	// Sync
	writeLock bool
	mu sync.Mutex
}

func (b *bankTab) filter() {
	i := 0
	old := len(b.filtered)
	for _, d := range b.bank.HIRC.HircObjs {
		if b.typeQuery != int32(d.HircType()) {
			continue
		}

		id, err := d.HircID()

		if err != nil {
			if i < len(b.filtered) {
				b.filtered[i] = d
			} else {
				b.filtered = append(b.filtered, d)
			}
			i += 1
			continue
		}

		if !fuzzy.Match(b.idQuery, strconv.FormatUint(uint64(id), 10)) {
			continue
		}
		if i < len(b.filtered) {
			b.filtered[i] = d
		} else {
			b.filtered = append(b.filtered, d)
		}
		i += 1
	}
	if i < old {
		b.filtered = slices.Delete(b.filtered, i, old)
	}
}

func (b *bankTab) encode(ctx context.Context) ([]byte, error) {
	b.mu.Lock()
	b.writeLock = true
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		b.writeLock = false
		b.mu.Unlock()
	}()

	type result struct {
		data []byte
		err  error
	}

	c := make(chan *result)
	go func() {
		data, err := b.bank.Encode(ctx)
		time.Sleep(time.Second * 16)
		c <- &result{data, err}
	}()

	select {
	case <- ctx.Done():
		<- c
		return nil, ctx.Err()
	case r := <- c:
		return r.data, r.err
	}
}

func (b *bankManager) openBank(ctx context.Context, path string) error {
	type result struct {
		bank *wwise.Bank
		err error
	}

	c := make(chan *result, 1)

	if _, in := b.banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}
	
	go func() {
		bank, err := parser.ParseBank(path, ctx)
		time.Sleep(time.Second * 15)
		c <- &result{bank, err}
	}()

	var bank *wwise.Bank
	select {
	case res := <- c:
		if res.bank == nil {
			return res.err
		}
		bank = res.bank
	case <- ctx.Done():
		return ctx.Err()
	}

	if _, in := b.banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}

	filtered := make([]wwise.HircObj, len(bank.HIRC.HircObjs))
	for i, o := range bank.HIRC.HircObjs {
		filtered[i] = o
	}

	b.banks.Store(path, &bankTab{
		writeLock: false,
		bank: bank,
		idQuery: "",
		typeQuery: 0,
		filtered: filtered,
		lSelStorage: *imgui.NewSelectionBasicStorage(),
	})

	return nil
}

func (b *bankManager) closeBank(del string) {
	b.banks.Delete(del)
}

type notfiy struct {
	message string
	timer *time.Timer
}

type notifyQ struct {
	queue []*notfiy
}

type guiLog struct {
	log *log.InMemoryLog
	debug bool
	info bool
	warn bool
	error bool
}
