package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
)

func renderDebug(bnkMngr *be.BankManager, dockMngr *dockmanager.DockManager) {
	if !dockMngr.ShowDebug {
		return
	}

	imgui.BeginV("Debug", &dockMngr.ShowDebug, 0)

	imgui.PushTextWrapPos()

	numBnks := 0
	bnkMngr.Banks.Range(func(key, value any) bool {
		numBnks += 1
		return true
	})

	activeBankName := ""
	bnkMngr.Banks.Range(func(key, value any) bool {
		if value.(*be.BankTab) == bnkMngr.ActiveBank {
			activeBankName = key.(string)
			return false
		}
		return true
	})

	imgui.SeparatorText("Clipboard")
	imgui.Text(fmt.Sprintf("Clipboard enabled: %v", GlobalCtx.CopyEnable))

	imgui.SeparatorText("Bank Manager")
	imgui.Text(fmt.Sprintf("# of bank tabs: %d", numBnks))
	imgui.Text(fmt.Sprintf("Active bank: %s", activeBankName))
	
	nextBankName := ""
	bnkMngr.Banks.Range(func(key, value any) bool {
		if value.(*be.BankTab) == bnkMngr.SetNextBank {
			nextBankName = key.(string)
			return false
		}
		return true
	})
	imgui.Text(fmt.Sprintf("Next Bank nil? %v", bnkMngr.SetNextBank == nil))
	imgui.Text(fmt.Sprintf("Next Bank: %s", nextBankName))

	mountedBnk := "None"
	if bnkMngr.InitBank != nil {
		bnkMngr.Banks.Range(func(key, value any) bool {
			if value.(*be.BankTab) == bnkMngr.InitBank {
				mountedBnk = key.(string)
				return false
			}
			return true
		})
	}
	imgui.Text(fmt.Sprintf("Mounted Init.bnk: %s", mountedBnk))
	imgui.SeparatorText("Modal")
	imgui.Text(fmt.Sprintf("# of modals: %d", len(GlobalCtx.ModalQ.Modals)))
	imgui.SeparatorText("Event Loop")
	stat := GlobalCtx.Loop.TaskStatus()
	imgui.Text(fmt.Sprintf("Sync task counter: %d", GlobalCtx.Loop.SyncTaskCounter))
	imgui.Text(fmt.Sprintf("Async task counter: %d", GlobalCtx.Loop.AsyncTaskCounter))
	imgui.Text(fmt.Sprintf("# of sync tasks: %d", stat.TotalNumSyncTask))
	imgui.Text(fmt.Sprintf("# of async tasks: %d", stat.TotalNumAsyncTask))
	imgui.Text(fmt.Sprintf("# of running async tasks: %d", stat.NumRunningAsyncTask))
	imgui.Text(fmt.Sprintf("# of pending async tasks: %d", stat.NumRunningAsyncTask))
	imgui.SeparatorText("Memory")
	imgui.PopTextWrapPos()
	imgui.End()
}
