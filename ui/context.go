package ui

import (
	"github.com/Dekr0/wwise-teller/aio"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

type Context struct {
	Loop       async.EventLoop
	ModalQ     ModalQ
	Config     config.Config
	CopyEnable bool
	Manager    aio.LRUWEMPlayersManger
}

var GlobalCtx Context = Context{
	async.NewEventLoop(),
	NewModalQ(),
	config.Config{},
	false,
	aio.LRUWEMPlayersManger{},
}
