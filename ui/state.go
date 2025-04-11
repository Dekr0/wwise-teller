package ui

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

type stateStore struct {
	ctx context.Context
	fileExplorer *fileExplorer
	eventBus *eventBus
}

type fileExplorer struct {
	pwd string
	entries []os.DirEntry /* cache os.ReadDir every time current directory is changed */
}

func newFileExplorer(p string) (*fileExplorer, error) {
	entries, err := getDirAndBank(p)
	if err != nil {
		return nil, err
	}
	return &fileExplorer{ p, entries }, nil
}

func getDirAndBank(p string) ([]os.DirEntry, error) {
	fd, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	newEntries := make([]os.DirEntry, 0, 1024)
	bound := 128
	for bound > 0 {
		entries, err := fd.ReadDir(1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				newEntries = append(newEntries, entry)
				continue
			}
			if filepath.Ext(entry.Name()) == ".bnk" {
				newEntries = append(newEntries, entry)
				continue
			}
		}

		bound -= 1
	}

	if bound <= 0 {
		return nil, fmt.Errorf(
			"Failed to read files from %s: upper bound is reached", p,
		)
	}
	
	return newEntries, nil
}

func (f *fileExplorer) cdParent() error {
	pwd := filepath.Dir(f.pwd)

	if runtime.GOOS == "windows" && pwd == "." {
		return nil
	}
	
	if pwd == f.pwd {
		return nil
	}

	newEntries, err := getDirAndBank(pwd)
	if err != nil {
		return err
	}

	f.pwd = pwd
	f.entries = newEntries

	return nil
}

func (f *fileExplorer) cd(basename string) error {
	pwd := filepath.Join(f.pwd, basename)

	newEntries, err := getDirAndBank(pwd)
	if err != nil {
		return err
	}

	f.pwd = pwd
	f.entries = newEntries

	return nil
}

type eventBus struct {
	eventHandler []func()
}

func newEventBus() *eventBus {
	return &eventBus{ make([]func(), 0, 32) }
}

func (b *eventBus) enqueue(name string, callback func(), ctx context.Context) {
	b.eventHandler = append(b.eventHandler, func() {
		if err := ctx.Err(); err != nil {
			slog.Warn("Operation was canceled", "handlerName", name, "reason", err)
			return
		}
		callback()
	})
}

func (b *eventBus) runAll() {
	for _, h := range b.eventHandler {
		h()
	}
	b.eventHandler = slices.Delete(b.eventHandler, 0, len(b.eventHandler))
}
