package Context

import (
	"github.com/Dekr0/wwise-teller/config"
	"github.com/Dekr0/wwise-teller/ui/async"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/ui/modal"
	"github.com/Dekr0/wwise-teller/ui/processor"
)

type Context struct {
	Loop       async.EventLoop
	ModalQ     modal.ModalQ
	Config     config.Config
	DockMngr   dockmanager.DockManager
	Editor     processor.ProcessorEditor
	CopyEnable bool
}

var GCtx Context = Context{}

func Init() {
	GCtx.Loop = async.NewEventLoop()
	GCtx.ModalQ = modal.NewModalQ()
	GCtx.Config = config.Config{}
	GCtx.Editor = processor.New()
	GCtx.CopyEnable = false
	dockmanager.NewDockManagerP(&GCtx.DockMngr)
}
