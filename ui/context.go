package ui

import (
	"log/slog"

	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

type Context struct {
	conf   *config.Config
	loop   *async.EventLoop
	modalQ *ModalQ
	nQ     *NotifyQ
}

func newContext() (*Context, error) {
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}
	slog.Info("Loaded configuration file.")
	loop := async.NewEventLoop()
	slog.Info("Created event loop.")
	modalQ := NewModalQ()
	slog.Info("Created modal queue.")
	nQ := &NotifyQ{make([]*notfiy, 0, 16)}
	slog.Info("Created notification queue.")
	return &Context{
		conf, loop, modalQ, nQ,
	}, nil
}
