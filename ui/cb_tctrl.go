package ui

import (
	"context"
	"fmt"
	"log/slog"

	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/waapi"
)

func createPlayerNoCacheTask(bnkTab *be.BankTab, sid uint32, wemData []byte) func(context.Context) {
	return func(ctx context.Context) {
		defer bnkTab.WEMExportLock.Store(false)
		waveFile, err := waapi.ExportWEMByte(ctx, wemData, true)
		if err != nil {
			bnkTab.UpdateErrorAudioSource(sid)
			slog.Error(fmt.Sprintf("Failed to export audio source %d", sid), "error", err)
			return
		}
		bnkTab.WEMExportCache.Store(sid, waveFile)
		err = GlobalCtx.PlayersManager.NewPlayer(waveFile)
		if err != nil {
			bnkTab.UpdateErrorAudioSource(sid)
			slog.Error(fmt.Sprintf("Failed to initialize audio player for audio source %d", sid), "error", err)
			return
		}
	}
}

func createPlayerNoCache(bnkTab *be.BankTab, sid uint32, wemData []byte) {
	ctx, cancel := context.WithCancel(context.Background())
	callback := createPlayerNoCacheTask(bnkTab, sid, wemData)
	procMsg := fmt.Sprintf("Initializing audio player for audio source %d", sid)
	doneMsg := fmt.Sprintf("Initialized audio player for audio source %d", sid)
	if err := GlobalCtx.Loop.QTask(ctx, cancel, procMsg, doneMsg, callback); err != nil {
		slog.Error(fmt.Sprintf("Failed to create background task to initialize audio player for audio source %d", sid), "error", err)
	} else {
		bnkTab.WEMExportLock.Store(true)
	}
}

func createPlayerCacheTask(tmpWAVPath string, sid uint32) func(context.Context) {
	return func(ctx context.Context) {
		defer GlobalCtx.PlayersManager.CreateLock.Store(false)
		err := GlobalCtx.PlayersManager.NewPlayer(tmpWAVPath)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to initialize audio player for audio source %d", sid), "error", err)
		}
	}
}

func createPlayerCache(tmpWAVPath string, sid uint32) {
	ctx, cancel := context.WithCancel(context.Background())
	callback := createPlayerCacheTask(tmpWAVPath, sid)
	procMsg := fmt.Sprintf("Initializing a new audio player for audio source %d", sid)
	doneMsg := fmt.Sprintf("Initialized a new audio player for audio source %d", sid)
	if err := GlobalCtx.Loop.QTask(ctx, cancel, procMsg, doneMsg, callback); err != nil {
		slog.Error(fmt.Sprintf("Failed to create background task to initialize audio player for audio source %d", sid), "error", err)
	} else {
		GlobalCtx.PlayersManager.CreateLock.Store(true)
	}
}
