package log

import "github.com/AllenDang/cimgui-go/imgui"

func RenderLog(gLog *GuiLog) {
	imgui.Begin("Log")
	gLog.Log.Logs.Do(func(a any) {
		if a == nil {
			return
		}
		imgui.PushTextWrapPos()
		imgui.Text(a.(string))
		imgui.PopTextWrapPos()
	})
	imgui.End()
}
