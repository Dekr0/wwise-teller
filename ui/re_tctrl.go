package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func RenderTransportControl(t *be.BankTab) {
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.WriteLock.Load() {
		return
	}

	imgui.Begin("Transport Control")
	defer imgui.End()

	didx := t.Bank.DIDX()
	if didx == nil {
		imgui.Text("This sound bank does not DIDX chunk.")
		return
	}

	data := t.Bank.DATA()
	if data == nil {
		imgui.Text("This sound bank does not have DATA chunk.")
		return
	}

	var sound *wwise.Sound
	switch ah := t.ActorMixerViewer.ActiveHirc.(type) {
	case *wwise.Sound:
		sound = ah
	}
	if sound == nil {
		return
	}

	sid := sound.BankSourceData.SourceID
	v, ok := t.WEMToWaveCache.Load(sid)
	if !ok {
		wemData, in := data.AudiosMap[sid]
		if !in { // TODO, FX based audio source
			return
		}
		createPlayerNoCache(&t.WEMToWaveCache, sid, wemData)
		return
	}
	tmpWAVPath := v.(string)
	ok = GlobalCtx.Manager.HasPlayer(tmpWAVPath)
	if !ok {
		createPlayerCache(tmpWAVPath, sid)
		return
	}
}
