package ui

import (
	"context"
	"time"
	uctx "github.com/Dekr0/wwise-teller/ui/context"
	"github.com/AllenDang/cimgui-go/imgui"
)

var ModifiyEverything = false
var GCtx     = &uctx.GCtx
var Config   = &GCtx.Config
var BnkMngr  = &GCtx.BankMngr
var DockMngr = &GCtx.DockMngr

func BG(
	timeout time.Duration,
	onProcMsg string,
	onDoneMsg string,
	f func(context.Context),
) {
	ctx, cancel := context.WithTimeout(GCtx.Ctx, timeout)
	if err := GCtx.Loop.QTask(ctx, cancel, onProcMsg, onDoneMsg, f); err != nil {
		cancel()
	}
}

func Modal(
	done *bool,
	flag imgui.WindowFlags,
	name string,
	loop func(),
	onClose func(),
) {
	GCtx.ModalQ.QModal(done, flag, name, loop, onClose)
}
