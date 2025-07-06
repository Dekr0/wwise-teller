package notify

import (
	"fmt"
	"slices"
	"github.com/AllenDang/cimgui-go/imgui"
)

func RenderNotify(nQ *NotifyQ) {
	imgui.BeginV("Notifications", 
		nil,
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoResize |
		imgui.WindowFlagsAlwaysAutoResize,
	)
	defer imgui.End()
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
}

