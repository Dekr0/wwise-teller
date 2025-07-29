// TODO: Rethink how to construct this for concurrency again. It's still very
// cooked.
package ui

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/implot"
	"github.com/AllenDang/cimgui-go/utils"
	"github.com/Dekr0/wwise-teller/aio"
	"github.com/Dekr0/wwise-teller/ui/audio"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	dockmanager "github.com/Dekr0/wwise-teller/ui/dock_manager"
	"github.com/Dekr0/wwise-teller/waapi"
	"github.com/Dekr0/wwise-teller/wwise"
)

func RenderTransportControl(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV(dockmanager.DockWindowNames[dockmanager.TransportControlTag], open, imgui.WindowFlagsNone)
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

	switch ah := t.ActorMixerViewer.ActiveHirc.(type) {
	case *wwise.Sound:
		renderPlaySound(t, data, ah)
	}
}

func renderPlaySound(bnkTab *be.BankTab, data *wwise.DATA, sound *wwise.Sound) {
	sid := sound.BankSourceData.SourceID
	if v, ok := bnkTab.ErrorAudioSources.Load(sid); ok && v.(int) >= be.MaxWEMExportRetrys {
		imgui.Text(fmt.Sprintf("Auto WEM export is suspended for audio source %d.", sid))
		return
	}

	v, ok := bnkTab.WEMExportCache.Load(sid)
	if !ok {
		if bnkTab.BusyWEMExport() {
			const msg = "A background audio source export task is running. Please wait..."
			imgui.ProgressBarV(float32(-1.0 * imgui.Time()), imgui.NewVec2(0, 0), msg)
			return
		}

		wemData, in := data.AudiosMap[sid]
		if !in { // FX based audio source is not covered yet
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		callback := func(ctx context.Context) {
			defer bnkTab.UnlockWEMExport()
			waveFile, err := waapi.ExportWEMByte(ctx, wemData, true)
			if err != nil {
				bnkTab.UpdateErrorAudioSource(sid)
				slog.Error(fmt.Sprintf("Failed to export audio source %d", sid), "error", err)
				return
			}
			bnkTab.WEMExportCache.Store(sid, waveFile)
		}
		proc := fmt.Sprintf("Exporting audio source %d", sid)
		done := fmt.Sprintf("Exported audio source %d", sid)
		if err := GCtx.Loop.QTask(ctx, cancel, proc, done, callback); err != nil {
			slog.Error(fmt.Sprintf("Failed to create background task to export audio source %d", sid), "error", err)
		} else {
			bnkTab.LockWEMExport()
		}
		return
	}

	soundId := sound.Id
	waveFile := v.(string)
	if v, ok := bnkTab.ErrorStreamers.Load(soundId); ok && v.(int) >= be.MaxInitStreamerRetrys {
		imgui.Text(fmt.Sprintf("Auto sound streamer initialization is suspended for sound %d", soundId))
		return
	}

	streamer, in := bnkTab.Session.Streamer(soundId)
	if !in {
		if bnkTab.Session.Busy() {
			const msg = "A background streamer initialization is running. Please wait..."
			imgui.ProgressBarV(float32(-1.0 * imgui.Time()), imgui.NewVec2(0, 0), msg)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		callback := func(ctx context.Context) {
			defer bnkTab.Session.Unlock()
			err := bnkTab.Session.NewSoundStreamerFile(ctx, soundId, waveFile)
			if err != nil {
				bnkTab.UpdateErrorStreamers(soundId)
				slog.Error(fmt.Sprintf("Failed to initialize sound streamer for sound %d", soundId), "error", err)
			}
		}
		proc := fmt.Sprintf("Initializing new sound streamer for sound %d", soundId)
		done := fmt.Sprintf("Initialized new sound streamer for sound %d", soundId)
		if err := GCtx.Loop.QTask(ctx, cancel, proc, done, callback); err != nil {
			slog.Error(fmt.Sprintf("Failed to create background task to initialize sound streamer for sound %d", soundId), "error", err)
		} else {
			bnkTab.Session.Lock()
		}
		return
	}

	var soundStreamer *audio.SoundStreamer
	switch t := streamer.(type) {
	case *audio.SoundStreamer:
		soundStreamer = t
	default: 
		panic("This should be a sound streamer!")
	}

	Disabled(!aio.BeepEnable, func() {
		if imgui.Button("Play") {
			if err := bnkTab.Session.Play(soundId); err != nil {
				slog.Error(fmt.Sprintf("Failed to play sound streamer for sound %d", soundId), "error", err)
			}
		}
		imgui.SameLine()
		if imgui.Button("Pause") {
			bnkTab.Session.Pause(soundId)
		}
		imgui.SameLine()
		if imgui.Button("Resume") {
			if err := bnkTab.Session.Resume(soundId); err != nil {
				slog.Error(fmt.Sprintf("Failed to resume sound streamer for sound %d", soundId), "error", err)
			}
		}
		imgui.SameLine()
		if imgui.Button("Loop") {
			bnkTab.Session.Pause(soundId)
			if err := soundStreamer.Loop(); err != nil {
				slog.Error(fmt.Sprintf("Failed to loop sound streamer for sound %d", soundId), "error", err)
			} else {
				if err := bnkTab.Session.Resume(soundId); err != nil {
					slog.Error(fmt.Sprintf("Failed to resume sound streamer for sound %d", soundId), "error", err)
				}
			}
		}
	})

	const plotFlags = implot.FlagsNoLegend | 
		              implot.FlagsNoTitle  | 
				      implot.FlagsNoFrame  |
		              implot.FlagsCanvasOnly
	const axisFlags = implot.AxisFlagsNoLabel      | 
			          implot.AxisFlagsNoTickMarks  |
		              implot.AxisFlagsNoTickLabels |
		              implot.AxisFlagsAutoFit
	const tableFlags = imgui.TableFlagsResizable |
		               imgui.TableFlagsBordersH  |
				       imgui.TableFlagsBordersV
	size := imgui.NewVec2(-1, 72)

	pos := 2000 * float64(soundStreamer.Position()) / float64(soundStreamer.Format.SampleRate)

	if imgui.BeginTableV("WaveformTable", 2, tableFlags, DefaultSize, 0) {
		imgui.TableSetupColumnV("Channel", imgui.TableColumnFlagsWidthFixed, 48, 0)
		imgui.TableSetupColumn("")
		imgui.TableHeadersRow()

		var channelNames []string = nil
		for _, channelConf := range waapi.ChannelLUT {
			if len(channelConf) == len(soundStreamer.PCMData) {
				channelNames = channelConf
			}
		}

		var text = ""
		for i, channelData := range soundStreamer.PCMData {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(48)
			if channelNames != nil {
				text = channelNames[i]
			} else {
				text = fmt.Sprintf("Channel %d", i)
			}
			imgui.Text(text)

			imgui.TableSetColumnIndex(1)
			implot.PushStyleVarVec2(implot.StyleVarPlotPadding, DefaultSize)
			if implot.BeginPlotV(fmt.Sprintf("Channel %d", i + 1), size, plotFlags) {
				implot.SetupAxesV("", "", axisFlags, axisFlags)
				implot.PlotLineS64PtrInt(
					fmt.Sprintf("Channel %d Sample Point", i + 1),
					utils.SliceToPtr(channelData),
					int32(len(channelData)),
					)
				implot.DragLineX(0, &pos, imgui.NewVec4(1, 0, 0, 1))
				if i + 1 == len(soundStreamer.PCMData) {
					implot.TagXStr(pos, imgui.NewVec4(1, 0, 0, 1), "")
				}
				implot.EndPlot()
			}
			implot.PopStyleVar()
		}
		imgui.EndTable()
	}
}
