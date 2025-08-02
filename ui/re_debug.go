package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
)

func renderDebug(open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Debug", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}

	imgui.PushTextWrapPos()

	numBnks := 0
	BnkMngr.Banks.Range(func(key, value any) bool {
		numBnks += 1
		return true
	})

	activeBankName := ""
	BnkMngr.Banks.Range(func(key, value any) bool {
		if value.(*be.BankTab) == BnkMngr.ActiveBank {
			activeBankName = key.(string)
			return false
		}
		return true
	})

	imgui.SeparatorText("Clipboard")
	imgui.Text(fmt.Sprintf("Clipboard enabled: %v", GCtx.CopyEnable))

	imgui.SeparatorText("Bank Manager")
	imgui.Text(fmt.Sprintf("# of bank tabs: %d", numBnks))
	imgui.Text(fmt.Sprintf("Active bank: %s", activeBankName))
	
	nextBankName := ""
	BnkMngr.Banks.Range(func(key, value any) bool {
		if value.(*be.BankTab) == BnkMngr.SetNextBank {
			nextBankName = key.(string)
			return false
		}
		return true
	})
	imgui.Text(fmt.Sprintf("Next Bank nil? %v", BnkMngr.SetNextBank == nil))
	imgui.Text(fmt.Sprintf("Next Bank: %s", nextBankName))

	mountedBnk := "None"
	if BnkMngr.InitBank != nil {
		BnkMngr.Banks.Range(func(key, value any) bool {
			if value.(*be.BankTab) == BnkMngr.InitBank {
				mountedBnk = key.(string)
				return false
			}
			return true
		})
	}
	imgui.Text(fmt.Sprintf("Mounted Init.bnk: %s", mountedBnk))
	imgui.SeparatorText("Modal")
	imgui.Text(fmt.Sprintf("# of modals: %d", len(GCtx.ModalQ.Modals)))
	imgui.SeparatorText("Event Loop")
	stat := GCtx.Loop.TaskStatus()
	imgui.Text(fmt.Sprintf("Sync task counter: %d", GCtx.Loop.SyncTaskCounter))
	imgui.Text(fmt.Sprintf("Async task counter: %d", GCtx.Loop.AsyncTaskCounter))
	imgui.Text(fmt.Sprintf("# of sync tasks: %d", stat.TotalNumSyncTask))
	imgui.Text(fmt.Sprintf("# of async tasks: %d", stat.TotalNumAsyncTask))
	imgui.Text(fmt.Sprintf("# of running async tasks: %d", stat.NumRunningAsyncTask))
	imgui.Text(fmt.Sprintf("# of pending async tasks: %d", stat.NumRunningAsyncTask))
	if BnkMngr.ActiveBank != nil {
		imgui.SeparatorText("Active Sound Bank Session")
		BnkMngr.ActiveBank.Session.Mutex.Lock()
		imgui.Text(fmt.Sprintf("# of Streamers: %d", len(BnkMngr.ActiveBank.Session.Streamers)))
		for _, streamer := range BnkMngr.ActiveBank.Session.Streamers {
			imgui.Text(fmt.Sprintf("%d - Is Nil: %v", streamer.Id(), streamer.UWStreamer() == nil))
		}
		BnkMngr.ActiveBank.Session.Mutex.Unlock()
	}
	imgui.SeparatorText("Memory")
	imgui.PopTextWrapPos()
}
