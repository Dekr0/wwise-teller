package Context

import (
	"context"
	"sync/atomic"

	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
	"github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/ui/modal"
	"github.com/Dekr0/wwise-teller/ui/processor"
)

type Context struct {
	Ctx        context.Context
	Loop       async.EventLoop
	ModalQ     modal.ModalQ
	Config     config.Config
	BankMngr   bank_explorer.BankManager
	DockMngr   dockmanager.DockManager
	Editor     processor.ProcessorEditor
	CopyEnable bool
}

var GCtx Context = Context{}

func Init() {
	GCtx.Ctx = context.Background()
	GCtx.Loop = async.NewEventLoop()
	GCtx.ModalQ = modal.NewModalQ()
	GCtx.Config = config.Config{}
	GCtx.Editor = processor.New()
	GCtx.CopyEnable = false
	GCtx.BankMngr = bank_explorer.BankManager{WriteLock: atomic.Bool{}}
	GCtx.BankMngr.WriteLock.Store(false)
	dockmanager.NewDockManagerP(&GCtx.DockMngr)
}
