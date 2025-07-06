package ui

import (
	"github.com/Dekr0/wwise-teller/aio"
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
)

type Context struct {
	Loop           async.EventLoop
	ModalQ         ModalQ
	Config         config.Config
	CopyEnable     bool
	PlayersManager aio.PlayersManager
}

var GlobalCtx Context = Context{
	async.NewEventLoop(),
	NewModalQ(),
	config.Config{},
	false,
	aio.PlayersManager{
		Max: aio.DefaultMaxLRUWEMPlayerManagerSize,
		Players: make([]*aio.Player, 0, 8),
	},
}
