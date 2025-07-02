package log

import "github.com/AllenDang/cimgui-go/imgui"

func RenderLog(gLog *GuiLog, showLog *bool) {
	if !*showLog {
		return
	}
	imgui.BeginV("Log", showLog, 0)
	defer imgui.End()
	if !*showLog {
		return
	}
	gLog.Log.Logs.Do(func(a any) {
		if a == nil {
			return
		}
		imgui.PushTextWrapPos()
		imgui.Text(a.(string))
		imgui.PopTextWrapPos()
	})
}
