// TODO: Rethink how to construct this for concurrency. This is very cooked.
package ui

import (
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/aio"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func RenderTransportControl(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("Transport Control", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}

	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.SounBankLock.Load() {
		return
	}

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

	if v, ok := t.ErrorAudioSources.Load(sid); ok && v.(int) >= be.MaxExportRetrys {
		imgui.Text(fmt.Sprintf("Export is suspended for audio source %d.", sid))
		return
	}

	v, ok := t.WEMExportCache.Load(sid)
	if !ok {
		if t.WEMExportLock.Load() {
			imgui.ProgressBarV(
				float32(-1.0 * imgui.Time()),
				imgui.NewVec2(0, 0),
				"A background audio source export task is running. Please wait...",
			)
			return
		}
		wemData, in := data.AudiosMap[sid]
		// TODO, FX based audio source
		if !in { 
			return
		}
		createPlayerNoCache(t, sid, wemData)
		return
	}

	waveFile := v.(string)
	mngr := &GlobalCtx.PlayersManager
	i := mngr.HasPlayer(waveFile)
	if i == -1 {
		if mngr.CreateLock.Load() {
			imgui.ProgressBarV(
				float32(-1.0 * imgui.Time()),
				imgui.NewVec2(0, 0),
				"A background audio player create task is running. Please wait...",
			)
			return
		}
		createPlayerCache(waveFile, sid)
		return
	}
	if i == -2 {
		imgui.Text(fmt.Sprintf("Audio player initialization is suspended for audio source %d", sid))
		return
	}

	Disabled(!aio.BeepEnable, func() {
		if imgui.Button("Play") {
			mngr.SetActivePlayer(waveFile)
			if err := mngr.PlayActivePlayer(); err != nil {
				slog.Error( "Failed to play active audio player", "error", err)
			}
		}
		imgui.SameLine()
		if imgui.Button("Pause") {
			mngr.SetActivePlayer(waveFile)
			mngr.PauseActivePlayer()
		}
		imgui.SameLine()
		if imgui.Button("Resume") {
			mngr.SetActivePlayer(waveFile)
			mngr.ResumeActivePlayer()
		}
		imgui.SameLine()
		if imgui.Button("Loop") {
			mngr.SetActivePlayer(waveFile)
			if err := mngr.LoopActivePlayer(); err != nil {
				slog.Error("Failed to loop audio cue", "error", err)
			}
		}
	})

	pos, l, format, in := mngr.PlayerInfo(waveFile)
	if in {
		fpos := float32(format.SampleRate.D(pos).Seconds())
		fl := float32(format.SampleRate.D(l).Seconds())
		imgui.Text(fmt.Sprintf("%f / %f", fpos, fl))
		imgui.BeginDisabled()
		imgui.SetNextItemWidth(-1)
		imgui.SliderFloat("##Timestamp", &fpos, 0.0, fl)
		imgui.EndDisabled()
		imgui.Text(fmt.Sprintf("# of Channels: %d", format.NumChannels))
	}
}
