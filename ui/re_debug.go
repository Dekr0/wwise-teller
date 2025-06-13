package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/ui/async"
)

func renderDebug(bnkMngr *BankManager, loop *async.EventLoop, modalQ *ModalQ) {
	imgui.Begin("Debug")
	imgui.PushTextWrapPos()

	numBnks := 0
	bnkMngr.Banks.Range(func(key, value any) bool {
		numBnks += 1
		return true
	})

	activeBankName := ""
	bnkMngr.Banks.Range(func(key, value any) bool {
		if value.(*BankTab) == bnkMngr.ActiveBank {
			activeBankName = key.(string)
			return false
		}
		return true
	})

	imgui.SeparatorText("Bank Manager")
	imgui.Text(fmt.Sprintf("# of bank tabs: %d", numBnks))
	imgui.Text(fmt.Sprintf("Active bank: %s", activeBankName))
	imgui.Text(fmt.Sprintf("Mounted Init.bnk? %v", bnkMngr.InitBank != nil))
	imgui.SeparatorText("Modal")
	imgui.Text(fmt.Sprintf("# of modals: %d", len(modalQ.modals)))
	imgui.SeparatorText("Event Loop")
	stat := loop.TaskStatus()
	imgui.Text(fmt.Sprintf("Sync task counter: %d", loop.SyncTaskCounter))
	imgui.Text(fmt.Sprintf("Async task counter: %d", loop.AsyncTaskCounter))
	imgui.Text(fmt.Sprintf("# of sync tasks: %d", stat.TotalNumSyncTask))
	imgui.Text(fmt.Sprintf("# of async tasks: %d", stat.TotalNumAsyncTask))
	imgui.Text(fmt.Sprintf("# of running async tasks: %d", stat.NumRunningAsyncTask))
	imgui.Text(fmt.Sprintf("# of pending async tasks: %d", stat.NumRunningAsyncTask))
	imgui.SeparatorText("Memory")
	imgui.PopTextWrapPos()
	imgui.End()
}
