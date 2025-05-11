package ui

import (
	"fmt"
	"slices"
	"github.com/AllenDang/cimgui-go/imgui"
)

func renderNotfiy(nQ *NotifyQ) {
	if len(nQ.queue) <= 0 {
		return
	}
	size := imgui.MainViewport().Size()
	y := imgui.CalcTextSize(nQ.queue[0].message).Y * float32(len(nQ.queue))
	x := float32(0.0)
	for _, n := range nQ.queue {
		x = max(imgui.CalcTextSize(n.message).X, x)
	}
	size.X -= size.X * 0.01 + x + 32.0
	size.Y -= size.Y * 0.03 + y + float32(len(nQ.queue)) * size.Y * 0.01
	imgui.SetNextWindowPos(size)
	imgui.BeginV("Notify", 
		nil,
		imgui.WindowFlagsNoDecoration |
		imgui.WindowFlagsNoTitleBar |
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoResize |
		imgui.WindowFlagsAlwaysAutoResize,
	)
	i := 0
	for i < len(nQ.queue) {
		select {
		case <- nQ.queue[i].timer.C:
			nQ.queue = slices.Delete(nQ.queue, i, i + 1) 
		default:
			imgui.PushIDStr(fmt.Sprintf("RemoveNotfiy_%d", i))
			if imgui.Button("X") {
				nQ.queue = slices.Delete(nQ.queue, i, i + 1)
				imgui.PopID()
				break
			}
			imgui.PopID()
			imgui.SameLine()
			imgui.Text(nQ.queue[i].message)
		}
		i += 1
	}
	imgui.End()
}

