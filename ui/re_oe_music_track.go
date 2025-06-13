package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"

	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderMusicTrack(t *be.BankTab, init *wwise.Bank, o *wwise.MusicTrack) {
	if imgui.TreeNodeExStr("Music Override Flags") {
		imgui.BeginDisabledV(o.BaseParam.DirectParentId == 0)
		overrideParentMIDITempo := o.OverrideParentMIDITempo()
		if imgui.Checkbox("Override Parent MIDI Tempo", &overrideParentMIDITempo) {
			o.SetOverrideParentMIDITempo(overrideParentMIDITempo)
		}

		overrideParentMIDITarget := o.OverrideParentMIDITarget()
		if imgui.Checkbox("Override Parent MIDI Target", &overrideParentMIDITarget) {
			o.SetOverrideParentMIDITarget(overrideParentMIDITarget)
		}

		midiTargetTypeBus := o.MidiTargetTypeBus()
		if imgui.Checkbox("MIDI Target Type Bus", &midiTargetTypeBus) {
			o.SetMidiTargetTypeBus(midiTargetTypeBus)
		}
		imgui.EndDisabled()
		imgui.TreePop()
	}
	// renderBankSourceData ???
	renderMusicTrackPlayList(t, o)
	// renderClipAutomation(t, o)
	renderBaseParam(t, init, o)
	renderSwitchParam(t, o)
	renderTransitionParam(&o.TransitionParam)
}

func renderMusicTrackPlayList(t *be.BankTab, o *wwise.MusicTrack) {
	if imgui.TreeNodeExStr("Music Track Play List") {
		const flags = DefaultTableFlags
		outerSize := imgui.NewVec2(0, 0)
		if imgui.BeginTableV("MusicTrackPlayListTable", 7, flags, outerSize, 0) {
			imgui.TableSetupColumnV("Track ID", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumnV("Source ID", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumnV("Event ID", imgui.TableColumnFlagsWidthFixed, 0, 0)
			imgui.TableSetupColumn("Play At")
			imgui.TableSetupColumn("Begin Trim Offset")
			imgui.TableSetupColumn("End Trim Offset")
			imgui.TableSetupColumn("Source Duration")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()
			
			for i := range o.PlayListItems {
				p := &o.PlayListItems[i]
				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)
				imgui.SetNextItemWidth(88)
				imgui.Text(strconv.FormatUint(uint64(p.TrackID), 10))

				imgui.TableSetColumnIndex(1)
				imgui.SetNextItemWidth(88)
				imgui.Text(strconv.FormatUint(uint64(p.SourceID), 10))

				imgui.TableSetColumnIndex(2)
				imgui.SetNextItemWidth(88)
				imgui.Text(strconv.FormatUint(uint64(p.EventID), 10))

				imgui.TableSetColumnIndex(3)
				imgui.SetNextItemWidth(-1)
				imgui.BeginDisabledV(!ModifiyEverything)
				playAt := float32(p.PlayAt)
				if imgui.InputFloat(fmt.Sprintf("##%dPlayAt%d", o.Id, i), &playAt) {
					p.PlayAt = float64(playAt)
				}
				imgui.EndDisabled()

				imgui.TableSetColumnIndex(4)
				imgui.SetNextItemWidth(-1)
				beginTrimOffset := float32(p.BeginTrimOffset)
				imgui.BeginDisabledV(!ModifiyEverything)
				if imgui.InputFloat(fmt.Sprintf("##%dBeginTrimOffset%d", o.Id, i), &beginTrimOffset) {
					p.BeginTrimOffset = float64(beginTrimOffset)
				}
				imgui.EndDisabled()

				imgui.TableSetColumnIndex(5)
				imgui.SetNextItemWidth(-1)
				imgui.BeginDisabledV(!ModifiyEverything)
				endTrimOffset := float32(p.EndTrimOffset)
				if imgui.InputFloat(fmt.Sprintf("##%dEndTrimOffset%d", o.Id, i), &endTrimOffset) {
					p.EndTrimOffset = float64(endTrimOffset)
				}
				imgui.EndDisabled()

				imgui.TableSetColumnIndex(6)
				imgui.PushIDStr(fmt.Sprintf("%dSrcDuration%d", o.Id, i))
				imgui.Text(strconv.FormatFloat(p.SrcDuration, 'f', -1, 64))
				imgui.PopID()
			}
			imgui.EndTable()
		}
		imgui.Text(fmt.Sprintf("Number of sub-track: %d", o.NumSubTrack))
		imgui.TreePop()
	}
}

func renderClipAutomation(t *be.BankTab, o *wwise.MusicTrack) {
	if imgui.TreeNodeExStr("Clip Automations") {
		if imgui.Button("Add New Automation") {
			o.AddNewAutomation()
		}
		var rmCA func() = nil
		for i := range o.ClipAutomations {
			c := &o.ClipAutomations[i]

			imgui.PushIDStr(fmt.Sprintf("RemoveAutomation%d", i))
			if imgui.Button("X") {
				rmCA = bindRmCA(o, i)
			}
			imgui.PopID()

			imgui.SameLine()

			if imgui.TreeNodeExStr(
				fmt.Sprintf(
					"%d. %s Automation %d",
					i + 1, wwise.ClipAutomationTypeName[c.AutoType], i,
				),
			) {
				autoType := int32(c.AutoType)
				if imgui.ComboStrarr(
					"Automation Type",
					&autoType,
					wwise.ClipAutomationTypeName,
					int32(len(wwise.ClipAutomationTypeName)),
				) {
					c.AutoType = uint32(autoType)
				}

				const flags = DefaultTableFlags
				outerSize := imgui.NewVec2(0, 0)
				if imgui.Button("Add RTPC Graph Point") {
					c.AddRTPCGraphPoint()
				}

				if imgui.BeginTableV(
					fmt.Sprintf("CARTPC%d", i),
					4,
					flags,
					outerSize, 0,
				) {
					rtpcPts := c.RTPCGraphPoints
					imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
					imgui.TableSetupColumn("From")
					imgui.TableSetupColumn("To")
					imgui.TableSetupColumn("Type")
					imgui.TableSetupScrollFreeze(0, 1)
					imgui.TableHeadersRow()

					var rmRTPCGraphPt func() = nil
					for j := range rtpcPts {
						pt := &rtpcPts[j]

						imgui.TableNextRow()

						imgui.BeginDisabled()
						imgui.TableSetColumnIndex(0)
						imgui.SetNextItemWidth(40)
						imgui.PushIDStr(fmt.Sprintf("CARTPC%dRM%d", i, j))
						if imgui.Button("X") {
							rmRTPCGraphPt = bindRmCARTPCGraphPt(c, j)
						}
						imgui.PopID()

						imgui.TableSetColumnIndex(1)
						imgui.SetNextItemWidth(-1)
						imgui.SliderFloat(
							fmt.Sprintf("##CARTPC%dFrom%d", i, j),
							&pt.From,
							-96.0,
							96.0,
						)

						imgui.TableSetColumnIndex(2)
						imgui.SetNextItemWidth(-1)
						imgui.SliderFloat(
							fmt.Sprintf("##CARTPC%dTo%d", i, j),
							&pt.To,
							-96.0,
							96.0,
						)
						imgui.EndDisabled()

						imgui.TableSetColumnIndex(3)
						imgui.SetNextItemWidth(-1)
						interp := int32(pt.Interp)
						if imgui.ComboStrarr(
							fmt.Sprintf("##CARTPC%dInterp%d", i, j),
							&interp,
							wwise.InterpCurveTypeName,
							int32(wwise.InterpCurveTypeCount),
						) {
							pt.Interp = uint32(interp)
						}
					}

					imgui.EndTable()
					if rmRTPCGraphPt != nil {
						rmRTPCGraphPt()
					}
				}

				imgui.TreePop()
			}
		}
		imgui.TreePop()
		if rmCA != nil {
			rmCA()
		}
	}
}

func bindRmCA(o *wwise.MusicTrack, i int) func() {
	return func() { o.RemoveAutomation(i) }
}

func bindRmCARTPCGraphPt(c *wwise.ClipAutomation, i int) func() {
	return func() { c.RemoveRTPCGraphPoint(i) }
}

func renderTransitionParam(p *wwise.MusicTrackTransitionParam) {
	if imgui.TreeNodeStr("Music Track Transition Parameter") {
		imgui.BeginDisabledV(!ModifiyEverything)
		imgui.SetNextItemWidth(128)
		imgui.InputInt("Source Transition Time", &p.SrcTransitionTime)
		imgui.EndDisabled()

		srcFadeCurve := int32(p.SrcFadeCurve)
		imgui.SetNextItemWidth(256)
		if imgui.ComboStrarr(
			"Source Fade Curve",
			&srcFadeCurve,
			wwise.InterpCurveTypeName,
			int32(wwise.InterpCurveTypeCount),
		) {
			p.SrcFadeCurve = uint32(srcFadeCurve)
		}

		imgui.SetNextItemWidth(128)
		imgui.InputInt("Source Fade Offset", &p.SrcFadeOffset)

		syncType := int32(p.SyncType)
		imgui.SetNextItemWidth(160)
		if imgui.ComboStrarr(
			"Source Sync Type",
			&syncType,
			wwise.SyncTypeName,
			wwise.NumSyncType,
			) {
			p.SyncType = uint32(syncType)
		}

		imgui.Text(fmt.Sprintf("Cue Filter Hash: %d", p.CueFilterHash))

		imgui.BeginDisabledV(!ModifiyEverything)
		imgui.SetNextItemWidth(128)
		imgui.InputInt("Destination Transition Time", &p.DestTransitionTime)
		imgui.EndDisabled()

		destFadeCurve := int32(p.DestFadeCurve)
		imgui.SetNextItemWidth(256)
		if imgui.ComboStrarr(
			"Destination Fade Curve",
			&destFadeCurve,
			wwise.InterpCurveTypeName,
			int32(wwise.InterpCurveTypeCount),
		) {
			p.DestFadeCurve = uint32(destFadeCurve)
		}

		imgui.BeginDisabledV(!ModifiyEverything)
		imgui.SetNextItemWidth(128)
		imgui.InputInt("Destination Fade Offset", &p.DestFadeOffset)
		imgui.EndDisabled()

		imgui.TreePop()
	}
}

func renderSwitchParam(t *be.BankTab, m *wwise.MusicTrack) {
	if imgui.TreeNodeStr("Music Track Switch Parameter") {
		imgui.Text("Group Type: " + wwise.GroupTypeName[m.SwitchParam.GroupType])
		imgui.Text(fmt.Sprintf("Group ID: %d", m.SwitchParam.GroupID))
		imgui.BeginDisabledV(!ModifiyEverything)
		imgui.Text("Default Switch ID")
		imgui.SetNextItemWidth(64)
		imgui.SameLine()
		imgui.InputScalar("##DefaultSwitchID", imgui.DataTypeU32, uintptr(utils.Ptr(&m.SwitchParam.DefaultSwitch)))
		imgui.EndDisabled()

		size := imgui.NewVec2(0, 160)
		const flags = DefaultTableFlags | imgui.TableFlagsScrollY
		if imgui.BeginTableV("MusicTrackSwitchParamTable", 1, flags, size, 0) {
			imgui.TableSetupColumn("Switch ID")
			imgui.TableSetupScrollFreeze(0, 1)
			imgui.TableHeadersRow()
			for _, i := range m.SwitchParam.SwitchAssociates {
				imgui.TableNextRow()
				imgui.TableSetColumnIndex(0)
				imgui.Text(strconv.FormatUint(uint64(i), 10))
			}
			imgui.EndTable()
		}
		imgui.TreePop()
	}
}
