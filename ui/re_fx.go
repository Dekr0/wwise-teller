package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderFxTable(t *be.BankTab) {
	focusTable := false

	useViUp()
	useViShiftUp()
	useViDown()
	useViShiftDown()

	imgui.SeparatorText("Filter")

	filterState := &t.FxViewer.Filter

	imgui.SetNextItemShortcut(DefaultSearchSC)
	imgui.SetNextItemWidth(96)
	if imgui.InputScalar(
		"By ID",
		imgui.DataTypeU32,
		uintptr(utils.Ptr(&filterState.Id)),
	) {
		t.FilterFxS()
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(128)
	preview := wwise.HircTypeName[filterState.Type]
	if imgui.BeginCombo("By Type", preview) {
		var filter func() = nil
		for _, _type := range wwise.FxHircTypes {
			selected := filterState.Type == _type
			preview = wwise.HircTypeName[_type]
			if imgui.SelectableBoolPtr(preview, &selected) {
				filterState.Type = _type
				filter = t.FilterFxS
			}
		}
		imgui.EndCombo()
		if filter != nil {
			filter()
		}
	}
	imgui.SeparatorText("")

	if imgui.Shortcut(UnFocusQuerySC) {
		focusTable = true
		imgui.SetKeyboardFocusHere()
	}

	if imgui.BeginTableV("FxSTable", 3, DefaultTableFlagsY, DefaultSize, 0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Type")
		imgui.TableSetupColumn("FX Type")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()
		if focusTable {
			imgui.SetKeyboardFocusHere()
		}

		clipper := imgui.NewListClipper()
		clipper.Begin(int32(len(filterState.Fxs)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				fx := filterState.Fxs[n]
				
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				if t.FxViewer.ActiveFx == nil {
					t.FxViewer.ActiveFx = fx
				}
				idA, err := fx.HircID()
				if err != nil { panic(err) }
				idB, err := t.FxViewer.ActiveFx.HircID()
				if err != nil { panic(err) }
				selected := idA == idB
				if imgui.SelectableBoolPtrV(
					strconv.FormatUint(uint64(idA), 10),
					&selected,
					DefaultSelectableFlags,
					DefaultSize,
				) {
					t.FxViewer.ActiveFx = fx
				}

				imgui.TableSetColumnIndex(1)
				imgui.Text(wwise.HircTypeName[fx.HircType()])

				imgui.TableSetColumnIndex(2)
				switch sfx := fx.(type) {
				case *wwise.FxCustom:
					name, in := wwise.PluginNameLUT[int32(sfx.PluginTypeId)]
					if !in {
						imgui.Text(fmt.Sprintf("Plugin ID %d", sfx.PluginTypeId))
					} else {
						imgui.Text(name)
					}
				case *wwise.FxShareSet:
					name, in := wwise.PluginNameLUT[int32(sfx.PluginTypeId)]
					if !in {
						imgui.Text(fmt.Sprintf("Plugin ID %d", sfx.PluginTypeId))
					} else {
						imgui.Text(name)
					}
				}
			}
		}
		imgui.EndTable()
	}
}

func renderFXViewer(t *be.BankTab, open *bool) {
	if !*open {
		return
	}
	imgui.BeginV("FX", open, imgui.WindowFlagsNone)
	defer imgui.End()
	if !*open {
		return
	}
	if t == nil || t.Bank == nil || t.Bank.HIRC() == nil || t.SounBankLock.Load() {
		return
	}
	if t.FxViewer.ActiveFx != nil {
		switch f := t.FxViewer.ActiveFx.(type) {
		case *wwise.FxShareSet:
			renderFxShareSet(f)
		case *wwise.FxCustom:
			renderFxCustom(f)
		default:
			panic("Panic trap")
		}
	}
}

func renderFxShareSet(f *wwise.FxShareSet) {
	name := fmt.Sprintf("Plugin %d", f.PluginTypeId)
	if n, in := wwise.PluginNameLUT[int32(f.PluginTypeId)]; in {
		name = n
	}
	imgui.SeparatorText(fmt.Sprintf("Fx Share Set %d - %s", f.Id, name))

	company := "Unknown"
	if n, in := wwise.PluginCompanyNames[f.PluginCompany()]; in {
		company = n
	}
	imgui.Text(fmt.Sprintf("Plugin Company: %s", company))

	renderFxParam(f.PluginParam.PluginParamData)
}

func renderFxCustom(f *wwise.FxCustom) {
	name := fmt.Sprintf("Plugin %d", f.PluginTypeId)
	if n, in := wwise.PluginNameLUT[int32(f.PluginTypeId)]; in {
		name = n
	}
	imgui.SeparatorText(fmt.Sprintf("Fx Share Set %d - %s", f.Id, name))

	company := "Unknown"
	if n, in := wwise.PluginCompanyNames[f.PluginCompany()]; in {
		company = n
	}
	imgui.Text(fmt.Sprintf("Plugin Company: %s", company))

	renderFxParam(f.PluginParam.PluginParamData)
}

func renderFxParam(f wwise.FxParam) {
	switch f := f.(type) {
	case *wwise.ParametricEQ:
		renderParametricEQ(f)
	case *wwise.MeterFX:
		renderFXMeter(f)
	case *wwise.PeakLimiter:
		renderPeakLimiter(f)
	case *wwise.GainFX:
		renderGainFX(f)
	case *wwise.Compressor:
		renderCompressor(f)
	}
}

func renderParametricEQ(f *wwise.ParametricEQ) {
	size := imgui.NewVec2(160, 160)
	for i := range f.EQBand {
		stack := fmt.Sprintf("Band %d", i + 1)

		imgui.BeginChildStrV(stack, size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

		imgui.SeparatorText(fmt.Sprintf("Band %d", i + 1))

		b := &f.EQBand[i]

		enabled := b.OnOff != 0
		imgui.PushIDStr(fmt.Sprintf("%sEnable", stack))
		if imgui.Checkbox("Enable", &enabled) {
			if enabled {
				b.OnOff = 1
			} else {
				b.OnOff = 0
			}
		}
		imgui.PopID()

		imgui.BeginDisabledV(!enabled)

		filterType := int32(b.FilterType)
		imgui.PushIDStr(fmt.Sprintf("%sCurve", stack))
		if imgui.ComboStrarrV("Curve", &filterType, wwise.EQFilterNames, int32(wwise.EQFilterTypeCount), 0) {
			b.FilterType = wwise.EQFilterType(filterType)
		}
		imgui.PopID()

		imgui.BeginDisabledV(
			b.FilterType == wwise.EQFilterTypeLowPass  || 
			b.FilterType == wwise.EQFilterTypeHiPass   || 
			b.FilterType == wwise.EQFilterTypeBandPass ||
			b.FilterType == wwise.EQFilterTypeNotch,
		)
		imgui.PushIDStr(fmt.Sprintf("%sGain", stack))
		imgui.SliderFloat("Gain", &b.Gain, -24, 24)
		imgui.PopID()
		imgui.EndDisabled()

		imgui.PushIDStr(fmt.Sprintf("%sFreq.", stack))
		imgui.SliderFloat("Freq.", &b.Frequency, 20, 20000)
		imgui.PopID()

		imgui.BeginDisabledV(
			b.FilterType == wwise.EQFilterTypeLowPass  || 
			b.FilterType == wwise.EQFilterTypeHiPass   || 
			b.FilterType == wwise.EQFilterTypeLowShelf ||
			b.FilterType == wwise.EQFilterTypeHiShelf,
		)
		imgui.PushIDStr(fmt.Sprintf("%sQ", stack))
		imgui.SliderFloat("Q", &b.QFactor, 0.5, 100)
		imgui.PopID()
		imgui.EndDisabled()

		imgui.EndDisabled()
		imgui.EndChild()
		if i < 2 {
			imgui.SameLine()
		}
	}

	imgui.SetNextItemWidth(96)
	imgui.SliderFloat("Output Gain", &f.OutputLevel, -24, 24)

	lfe := f.ProcessLFE != 0
	if imgui.Checkbox("Process LFE", &lfe) {
		if lfe {
			f.ProcessLFE = 1
		} else {
			f.ProcessLFE = 0
		}
	}
}

func renderFXMeter(f *wwise.MeterFX) {
	size := imgui.NewVec2(170, 48)
	{
		imgui.BeginChildStrV("ModeScope", size, 0, imgui.WindowFlagsNone)
		if f.Mode != nil {
			mode := int32(*f.Mode)
			imgui.Text("Mode ")
			imgui.SameLine()
			imgui.SetNextItemWidth(96)
			if imgui.ComboStrarrV("##Mode", &mode, wwise.MeterModeNames, int32(wwise.MeterModeCount), 0) {
				*f.Mode = wwise.MeterMode(mode)
			}
		}
		if f.Scope != nil {
			scope := int32(*f.Scope)
			imgui.Text("Scope")
			imgui.SameLine()
			imgui.SetNextItemWidth(110)
			if imgui.ComboStrarrV("##Scope", &scope, wwise.MeterScopeNames, int32(wwise.MeterScopeCount), 0) {
				*f.Scope = wwise.MeterScope(scope)
			}
		}
		imgui.EndChild()
	}

	imgui.SameLine()
	size.X, size.Y = 196, 110
	{
		imgui.BeginChildStrV("Dynamics", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Dynamics")
		imgui.SliderFloat("Attack", &f.Attack, 0, 10)
		if f.InfiniteHold != nil {
			inifiniteHold := *f.InfiniteHold == 1
			if imgui.Checkbox("Infinite Hold", &inifiniteHold) {
				if inifiniteHold {
					*f.InfiniteHold = 1
				} else {
					*f.InfiniteHold = 0
				}
			}
		}
		imgui.SliderFloat("Hold", &f.Hold, 0, 10)
		imgui.SliderFloat("Release", &f.Release, 0, 10)
		imgui.EndChild()
	}

	imgui.SameLine()
	size.X, size.Y = 250, 130
	{
		imgui.BeginChildStrV("Output Game Parameter", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
		imgui.SeparatorText("Output Game Parameter")
		imgui.Text(fmt.Sprintf("Game Parameter ID %d", f.GameParamID))
		imgui.SliderFloat("Min", &f.Min, -96.3, 0)
		imgui.SliderFloat("Max", &f.Max, -96.3, 12)
		applyDownStreamVol := f.ApplyDownstreamVolume == 1
		if imgui.Checkbox("Apply downstream volume", &applyDownStreamVol) {
			if applyDownStreamVol {
				f.ApplyDownstreamVolume = 1
			} else {
				f.ApplyDownstreamVolume = 0
			}
		}
		imgui.EndChild()
	}
}

func renderGainFX(f *wwise.GainFX) {
	imgui.PushItemWidth(96.0)
	imgui.SliderFloat("Full-Band Channels", &f.FullbandGain, -96.3, 24.0)
	imgui.SliderFloat("LFE", &f.LFEGain, -96.3, 24.0)
	imgui.PopItemWidth()
}

func renderCompressor(f *wwise.Compressor) {
	imgui.PushItemWidth(96.0)
	imgui.SliderFloat("Threshold", &f.Threshold, -96.3, 0.0)
	imgui.SliderFloat("Ratio", &f.Ratio, 1, 50)
	imgui.SliderFloat("Attack", &f.Attack, 0, 2)
	imgui.SliderFloat("Release", &f.Release, 0, 2)
	imgui.SliderFloat("Output Gain", &f.OutputGain, -24, 24)
	imgui.PopItemWidth()

	lfe := f.ProcessLFE == 1
	if imgui.Checkbox("Process LFE", &lfe) {
		if lfe {
			f.ProcessLFE = 1
		} else {
			f.ProcessLFE = 0
		}
	}

	channelLink := f.ChannelLink == 1
	if imgui.Checkbox("Channel link", &channelLink) {
		if lfe {
			f.ChannelLink = 1
		} else {
			f.ChannelLink = 0
		}
	}
}

func renderPeakLimiter(f *wwise.PeakLimiter) {
	imgui.PushItemWidth(96.0)
	imgui.SliderFloat("Threshold", &f.Threshold, -96.3, 0.0)
	imgui.SliderFloat("Ratio", &f.Ratio, 1, 50)
	imgui.SliderFloat("Look ahead time", &f.LookAhead, 0.001, 0.02)
	imgui.SliderFloat("Release", &f.Release, 0.001, 0.5)
	imgui.SliderFloat("Output Gain", &f.OutputLevel, -24, 24)
	imgui.PopItemWidth()

	lfe := f.ProcessLFE == 1
	if imgui.Checkbox("Process LFE", &lfe) {
		if lfe {
			f.ProcessLFE = 1
		} else {
			f.ProcessLFE = 0
		}
	}

	channelLink := f.ChannelLink == 1
	if imgui.Checkbox("Channel link", &channelLink) {
		if lfe {
			f.ChannelLink = 1
		} else {
			f.ChannelLink = 0
		}
	}
}
