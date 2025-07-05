package ui

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Dekr0/wwise-teller/waapi"
)

func createPlayerNoCache(cache *sync.Map, sid uint32, wemData []byte) {
	ctx, cancel := context.WithCancel(context.Background())
	callback := func(ctx context.Context) {
		tmpWAVPath, err := waapi.WEMToWAVEByte(ctx, wemData)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to transform audio source %d data to WAV", sid), "error", err)
		}
		cache.Store(sid, tmpWAVPath)
		_, err = GlobalCtx.Manager.Player(tmpWAVPath)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to create audio player for audio source %d (No caching)", sid), "error", err)
		}
	}
	procMsg := fmt.Sprint("Creating audio player for audio source %d (No caching)", sid)
	doneMsg := fmt.Sprint("Done creating audio player for audio source %d (No caching)", sid)
	if err := GlobalCtx.Loop.QTask(ctx, cancel, procMsg, doneMsg, callback); err != nil {
		slog.Error(fmt.Sprintf("Failed to create background task to create audio player for audio source %d (No caching)", sid), "error", err)
	}
}

func createPlayerCache(tmpWAVPath string, sid uint32) {
	ctx, cancel := context.WithCancel(context.Background())
	callback := func(ctx context.Context) {
		_, err := GlobalCtx.Manager.Player(tmpWAVPath)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to create audio player for audio source %d (Cached)", sid), "error", err)
		}
	}
	procMsg := fmt.Sprint("Creating audio player for audio source %d (Cached)", sid)
	doneMsg := fmt.Sprint("Done creating audio player for audio source %d (Cached)", sid)
	if err := GlobalCtx.Loop.QTask(ctx, cancel, procMsg, doneMsg, callback); err != nil {
		slog.Error(fmt.Sprintf("Failed to create background task to create audio player for audio source %d (Cached)", sid), "error", err)
	}
}
