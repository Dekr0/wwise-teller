package ui

import (
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

type Context struct {
	Loop      async.EventLoop
	ModalQ    ModalQ
	Config    config.Config
}

var GlobalCtx Context = Context{
	async.NewEventLoop(),
	NewModalQ(),
	config.Config{},
}
