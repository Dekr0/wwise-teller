// TODO
package ui

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"
	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
	"golang.design/x/clipboard"
)

func renderRanSeqPlayList(t *be.BankTab, r *wwise.RanSeqCntr) {
	if imgui.TreeNodeExStr("Random / Sequence Container Playlist Setting") {
		imgui.PushItemWidth(160)

		renderRanSeqPlayListSetting(&r.PlayListSetting)

		imgui.PopItemWidth()

		renderRanSeqPlayListTableSet(t, r)

		imgui.TreePop()
	}
}

func renderRanSeqPlayListTableSet(t *be.BankTab, r *wwise.RanSeqCntr) {
	imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
	if imgui.Button("Copy IDs") {
		var builder strings.Builder
		var err error
		for _, p := range r.PlayListItems {
			if _, err = builder.WriteString(strconv.FormatUint(uint64(p.UniquePlayID), 10)+"\n"); err != nil {
				slog.Error("Failed to copy items' ID", "error", err)
				break
			}
		}
		clipboard.Write(clipboard.FmtText, []byte(builder.String()))
	}
	imgui.EndDisabled()

	imgui.SameLine()

	imgui.BeginDisabledV(!GlobalCtx.CopyEnable)
	if imgui.Button("Copy Source IDs") {
		h := t.Bank.HIRC()
		if h == nil {
			panic("A random / sequence container object is being rendered but no HIRC chunk is found.")
		}
		var builder strings.Builder
		var err error
		for _, p := range r.PlayListItems {
			v, in := h.ActorMixerHirc.Load(p.UniquePlayID)
			if !in {
				panic(fmt.Sprintf("Hierarchy object %d doesn't exist", p.UniquePlayID))
			}
			switch sound := v.(wwise.HircObj).(type) {
			case *wwise.Sound:
				if _, err = builder.WriteString(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10) + "\n"); err != nil {
					slog.Error("Failed to copy items' source ID", "error", err)
					break
				}
			}
		}
		clipboard.Write(clipboard.FmtText, []byte(builder.String()))
	}
	imgui.EndDisabled()

	size := imgui.NewVec2(0, 256)
	if imgui.BeginTableV("PLTransfer", 3, 0, size, 0) {
		v := t.ActorMixerViewer
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthStretch, 0, 0)
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableNextRow()

		imgui.TableSetColumnIndex(0)
		renderRanSeqPlayListPendingTable(t, r)

		imgui.TableSetColumnIndex(1)
		if imgui.Button(">>") {
			for i := range r.Container.Children {
				r.AddLeafToPlayList(i)
			}
			v.RanSeqPlaylistStorage.Clear()
		}

		imgui.TableSetColumnIndex(2)
		renderRanSeqPlayListTable(t, r)

		imgui.EndTable()
	}
}

func renderRanSeqPlayListPendingTable(t *be.BankTab, r *wwise.RanSeqCntr) {
	size := imgui.NewVec2(200, 256)
	imgui.BeginChildStrV("PLPendingCell", size, imgui.ChildFlagsNone, imgui.WindowFlagsNone)
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY

	hirc := t.Bank.HIRC()
	size.X = 0

	if imgui.BeginTableV("PLPendTable", 3, flags, size, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Target ID")
		imgui.TableSetupColumn("Source ID")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var toPlayList func() = nil

		for i, child := range r.Container.Children {
			if slices.ContainsFunc(
				r.PlayListItems, func(p wwise.PlayListItem) bool {
					return p.UniquePlayID == child
				},
			) {
				continue
			}
			
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)
			imgui.PushIDStr(fmt.Sprintf("ToPlayList%d", i))
			if imgui.Button(">") {
				toPlayList = bindToRanSeqPlayList(&t.ActorMixerViewer, r, i)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatUint(uint64(child), 10))

			imgui.TableSetColumnIndex(2)
			value, ok := hirc.ActorMixerHirc.Load(child)
			if ok {
				sound := value.(wwise.HircObj).(*wwise.Sound)
				imgui.Text(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10))
			}
		}

		imgui.EndTable()

		if toPlayList != nil {
			toPlayList()
		}
	}
	imgui.EndChild()
}

func bindToRanSeqPlayList(v *be.ActorMixerViewer, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.AddLeafToPlayList(i)
		v.RanSeqPlaylistStorage.Clear()
	}
}

func renderRanSeqPlayListTable(t *be.BankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLCell")
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	size := imgui.NewVec2(0, 256)

	if imgui.BeginTableV("PLTable", 5, flags, size, 0) {
		v := t.ActorMixerViewer
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Sequence")
		imgui.TableSetupColumn("Target ID")
		imgui.TableSetupColumn("Source ID")
		imgui.TableSetupColumn("Weight")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var move func() = nil
		var del func() = nil
		var delSel func() = nil

		flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d
		storageSize := v.RanSeqPlaylistStorage.Size()
		itemCount := int32(len(r.PlayListItems))
		msIO := imgui.BeginMultiSelectV(flags, storageSize, itemCount)
		v.RanSeqPlaylistStorage.ApplyRequests(msIO)

		hirc := t.Bank.HIRC()
		for i, p := range r.PlayListItems {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)

			imgui.PushIDStr(fmt.Sprintf("DelPlayListItem%d", i))
			if imgui.Button("X") {
				del = bindPendRanSeqPlayListItem(&t.ActorMixerViewer, r, i)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			if c := renderRanSeqPlayListItemOrderCombo(i, r); c != nil {
				move = c
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			selected := v.RanSeqPlaylistStorage.Contains(imgui.ID(i))
			label := strconv.FormatUint(uint64(p.UniquePlayID), 10)
			const flags = DefaultTableSelFlags
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(i))
			imgui.SelectableBoolPtrV(label, &selected, flags, imgui.NewVec2(0, 0))

			if c := renderRanSeqPlayListTableCtxMenu(t, r); c != nil {
				delSel = c
			}

			value, ok := hirc.ActorMixerHirc.Load(p.UniquePlayID)
			if ok {
				imgui.TableSetColumnIndex(3)
				sound := value.(wwise.HircObj).(*wwise.Sound)
				imgui.Text(strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10))
			}

			imgui.TableSetColumnIndex(4)
			imgui.SetNextItemWidth(-1)
			imgui.PushIDStr(fmt.Sprintf("PLWeight%d", i))
			imgui.InputScalar("", imgui.DataTypeU32, uintptr(utils.Ptr(&p.Weight)))
			imgui.PopID()
		}

		imgui.EndMultiSelect()
		v.RanSeqPlaylistStorage.ApplyRequests(msIO)

		if move != nil { 
			move()
		}
		if del != nil {
			del()
		}
		if delSel != nil {
			delSel()
		}

		imgui.EndTable()
	}
	imgui.EndChild()
}

func renderRanSeqPlayListItemOrderCombo(i int, r *wwise.RanSeqCntr) func() {
	var move func() = nil

	preview := strconv.FormatUint(uint64(i), 10)
	label := fmt.Sprintf("PLSequence%d", i)
	if imgui.BeginCombo(label, preview) {
		for j := range r.PlayListItems {
			selected := i == j
			label := strconv.FormatUint(uint64(j), 10)
			if imgui.SelectableBoolPtr(label, &selected) {
				move = bindChangeRanSeqPlayListItemOrder(i, j, r)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}

	return move
}

func renderRanSeqPlayListTableCtxMenu(t *be.BankTab, r *wwise.RanSeqCntr) func() {
	var delSel func() = nil

	if imgui.BeginPopupContextItem() {
		if imgui.Button("Delete") {
			delSel = bindPendSelectRanSeqPlayListItem(&t.ActorMixerViewer, r)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	return delSel
}

func bindChangeRanSeqPlayListItemOrder(i, j int, r *wwise.RanSeqCntr) func() {
	return func() {
		r.MovePlayListItem(i, j)
	}
}

func bindPendRanSeqPlayListItem(v *be.ActorMixerViewer, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.RemoveLeafFromPlayList(i)
		v.RanSeqPlaylistStorage.Clear()
	}
}

func bindPendSelectRanSeqPlayListItem(v *be.ActorMixerViewer, r *wwise.RanSeqCntr) func() {
	return func() {
		mut := false
		tids := []uint32{}
		for i, p := range r.PlayListItems {
			if v.RanSeqPlaylistStorage.Contains(imgui.ID(i)) {
				tids = append(tids, p.UniquePlayID)
				mut = true
			}
		}
		if mut {
			v.RanSeqPlaylistStorage.Clear()
		}
		r.RemoveLeafsFromPlayList(tids)
	}
}

func renderRanSeqPlayListSetting(p *wwise.PlayListSetting) {
	size := imgui.NewVec2(0, 144)

	imgui.BeginChildStrV("Play Type", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

	// Play Type (Random)
	size.X = imgui.ContentRegionAvail().X * 0.5
	size.Y = 0
	imgui.BeginChildStrV("RandomMode", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

	if imgui.RadioButtonBool("Random", p.Random()) {
		p.UseRandom()
	}

	imgui.Separator()
	if imgui.RadioButtonBool("Standard", p.RandomModeNormal()) {
		p.UseRandomModeNormal()
	}
	if imgui.RadioButtonBool("Shuffle", p.RandomModeShuffle()) {
		p.UseRandomModeShuffle()
	}

	avoidRepeatCount := int32(p.AvoidRepeatCount)
	imgui.Text("Avoid Reapting Last")
	imgui.SetNextItemWidth(80)
	if imgui.InputInt("##Avoidrepeatinglast", &avoidRepeatCount) {
		if avoidRepeatCount >= 0 && avoidRepeatCount <= 999 {
			p.AvoidRepeatCount = uint16(avoidRepeatCount)
		}
	}
	imgui.SameLine()
	imgui.Text("played")
	imgui.EndChild()
	// End Play Type (Random)

	imgui.SameLine()

	// Play Type (Sequence)
	size.X = 0
	size.Y = 0
	imgui.BeginChildStrV("SequenceMode", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

	if imgui.RadioButtonBool("Sequence", p.Sequence()) {
		p.UseSequence()
	}

	imgui.Separator()

	imgui.Text("At end of playlist")
	if imgui.RadioButtonBool("Restart", !p.RestartBackward()) {
		p.SetRestartBackward(false)
	}
	if imgui.RadioButtonBool("Play in reverse order", p.RestartBackward()) {
		p.SetRestartBackward(true)
	}
	imgui.EndChild()
	// End Play Type (Sequence)
	imgui.EndChild()
	// End Play Type

	// Play Mode
	size.X = 0
	size.Y = 256
	imgui.BeginChildStrV("Play Mode", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)
	if imgui.RadioButtonBool("Step", !p.Continuous()) {
		p.SetContinuous(false)
	}
	imgui.SameLine()
	if imgui.RadioButtonBool("Continuous", p.Continuous()) {
		p.SetContinuous(true)
	}

	resetPlayList := p.ResetPlayListAtEachPlay()
	if imgui.Checkbox("Always reset playlist", &resetPlayList) {
		p.SetResetPlayListAtEachPlay(resetPlayList)
	}

	// Loop
	size.X = imgui.ContentRegionAvail().X * 0.5
	size.Y = 0
	imgui.BeginChildStrV("Loop", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

	imgui.SeparatorText("Loop")
	if imgui.Button("Use Infinite") {
		p.UseInfiniteLoop()
	}
	count := int32(p.LoopCount)
	countMin := int32(p.LoopModMin)
	countMax := int32(p.LoopModMax)
	imgui.SetNextItemWidth(96)
	if imgui.InputInt("# of Loops", &count) {
		if count >= 0 && count <= 32767 {
			p.LoopCount = uint16(count)
		}
	}
	imgui.SeparatorText("Loop Count Randomizer")
	imgui.Text("Min (Negative)")
	imgui.SetNextItemWidth(56)
	if imgui.SliderInt("##LoopModMinSlider", &countMin, 0, 32766) {
		p.LoopModMin = uint16(countMin)
	}
	imgui.SameLine()
	imgui.SetNextItemWidth(96)
	if imgui.InputInt("##LoopModMinInputV", &countMin) {
		if countMin >= 0 && countMin <= 32766 {
			p.LoopModMin = uint16(countMin)
		}
	}
	imgui.Text("Max (Positive)")
	imgui.SetNextItemWidth(56)
	if imgui.SliderInt("##LoopModMaxSlider", &countMax, 0, 32766) {
		p.LoopModMax = uint16(countMax)
	}
	imgui.SameLine()
	imgui.SetNextItemWidth(96)
	if imgui.InputInt("##LoopModMaxInputV", &countMax) {
		if countMax >= 0 && countMax <= 32766 {
			p.LoopModMax = uint16(countMax)
		}
	}
	imgui.EndChild()
	// End Of Loop

	imgui.SameLine()

	// Transition
	size.X = 0
	size.Y = 0
	imgui.BeginChildStrV("Transition", size, imgui.ChildFlagsBorders, imgui.WindowFlagsNone)

	imgui.Text("Transition Type")
	transitionMode := int32(p.TransitionMode)
	if imgui.ComboStrarrV(
		"##Transition Type",
		&transitionMode,
		wwise.TransitionModeString,
		wwise.TransitionModeCount,
		0,
	) {
		p.TransitionMode = uint8(transitionMode)
	}
	
	imgui.Text("Duration")
	imgui.SetNextItemWidth(72)
	imgui.SliderFloat("##DurationSlide", &p.TransitionTime, 0, 10000.0)
	prevTransitionTime := p.TransitionTime
	imgui.SameLine()
	imgui.SetNextItemWidth(72)
	if imgui.InputFloat("##DurationInputV", &p.TransitionTime) {
		if p.TransitionTime < 0.0 || p.TransitionTime > 10000.0 {
			p.TransitionTime = prevTransitionTime
		}
	}

	imgui.SeparatorText("Duration Randomizer")
	imgui.Text("Min")
	imgui.SetNextItemWidth(76)
	imgui.SliderFloat("##TransitionTimeModMinSlider", &p.TransitionTimeModMin, -10000.0, 0)
	imgui.SameLine()
	imgui.SetNextItemWidth(76)
	prevTransitionTimeModMin := p.TransitionTimeModMin
	if imgui.InputFloat("##TransitionTimeModMinInputV", &p.TransitionTimeModMin) {
		if p.TransitionTimeModMin < -10000.0 || p.TransitionTimeModMin > 0 {
			p.TransitionTimeModMin = prevTransitionTimeModMin
		}
	}
	imgui.Text("Max")
	imgui.SetNextItemWidth(76)
	imgui.SliderFloat("##TransitionTimeModMaxSlider", &p.TransitionTimeModMax, 0, 10000.0)
	imgui.SameLine()
	imgui.SetNextItemWidth(76)
	prevTransitionTimeModMax := p.TransitionTimeModMax
	if imgui.InputFloat("##TransitionTimeModMaxInputV", &p.TransitionTimeModMax) {
		if p.TransitionTimeModMax < 0.0 || p.TransitionTimeModMax > 10000.0 {
			p.TransitionTimeModMax = prevTransitionTimeModMax
		}
	}

	imgui.EndChild()
	// End Of Transition
	imgui.EndChild()
	// End Of Play Mode
}
