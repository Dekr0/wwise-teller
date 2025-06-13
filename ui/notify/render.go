package notify

import (
	"fmt"
	"slices"
	"github.com/AllenDang/cimgui-go/imgui"
)

func RenderNotify(nQ *NotifyQ) {
	if len(nQ.Queue) <= 0 {
		return
	}
	size := imgui.MainViewport().Size()
	y := imgui.CalcTextSize(nQ.Queue[0].Msg).Y * float32(len(nQ.Queue))
	x := float32(0.0)
	for _, n := range nQ.Queue {
		x = max(imgui.CalcTextSize(n.Msg).X, x)
	}
	size.X -= size.X * 0.01 + x + 32.0
	size.Y -= size.Y * 0.03 + y + float32(len(nQ.Queue)) * size.Y * 0.01
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
	for i < len(nQ.Queue) {
		select {
		case <- nQ.Queue[i].Timer.C:
			nQ.Queue = slices.Delete(nQ.Queue, i, i + 1) 
		default:
			imgui.PushIDStr(fmt.Sprintf("RemoveNotfiy_%d", i))
			if imgui.Button("X") {
				nQ.Queue = slices.Delete(nQ.Queue, i, i + 1)
				imgui.PopID()
				break
			}
			imgui.PopID()
			imgui.SameLine()
			imgui.Text(nQ.Queue[i].Msg)
		}
		i += 1
	}
	imgui.End()
}

