// TODO
// - Add source ID column
package ui

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderRanSeqPlayList(t *BankTab, r *wwise.RanSeqCntr) {
	if imgui.TreeNodeExStr("Random / Sequence Container Playlist Setting") {
		imgui.PushItemWidth(160)

		renderRanSeqPlayListSetting(r.PlayListSetting)

		imgui.PopItemWidth()

		renderPlayListTableSet(t, r)

		imgui.TreePop()
	}
}

func renderPlayListTableSet(t *BankTab, r *wwise.RanSeqCntr) {
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLTransfer", 3, 0, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthStretch, 0, 0)
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableNextRow()

		imgui.TableSetColumnIndex(0)
		renderPlayListPendingTable(t, r)

		imgui.TableSetColumnIndex(1)
		if imgui.Button(">>") {
			for i := range r.Container.Children {
				r.AddLeafToPlayList(i)
			}
			t.RanSeqPlaylistStorage.Clear()
		}

		imgui.TableSetColumnIndex(2)
		renderPlayListTable(t, r)

		imgui.EndTable()
	}
}

func renderPlayListPendingTable(t *BankTab, r *wwise.RanSeqCntr) {
	size := imgui.NewVec2(200, 0)
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
				r.PlayListItems, func(p *wwise.PlayListItem) bool {
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
				toPlayList = bindToPlayList(t, r, i)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatUint(uint64(child), 10))

			imgui.TableSetColumnIndex(2)
			value, ok := hirc.Sounds.Load(child)
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

func bindToPlayList(t *BankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.AddLeafToPlayList(i)
		t.RanSeqPlaylistStorage.Clear()
	}
}

func renderPlayListTable(t *BankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLCell")
	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLTable", 5, flags, outerSize, 0) {
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
		storageSize := t.RanSeqPlaylistStorage.Size()
		itemCount := int32(len(r.PlayListItems))
		msIO := imgui.BeginMultiSelectV(flags, storageSize, itemCount)
		t.RanSeqPlaylistStorage.ApplyRequests(msIO)

		hirc := t.Bank.HIRC()
		for i, p := range r.PlayListItems {
			imgui.TableNextRow()

			imgui.TableSetColumnIndex(0)
			imgui.SetNextItemWidth(40)

			imgui.PushIDStr(fmt.Sprintf("DelPlayListItem%d", i))
			if imgui.Button("X") {
				del = bindPendPlayListItem(t, r, i)
			}
			imgui.PopID()

			imgui.TableSetColumnIndex(1)
			imgui.SetNextItemWidth(-1)

			if c := renderPlayListItemOrderCombo(i, r); c != nil {
				move = c
			}

			imgui.TableSetColumnIndex(2)
			imgui.SetNextItemWidth(-1)

			selected := t.RanSeqPlaylistStorage.Contains(imgui.ID(i))
			label := strconv.FormatUint(uint64(p.UniquePlayID), 10)
			const flags = DefaultTableSelFlags
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(i))
			imgui.SelectableBoolPtrV(label, &selected, flags, imgui.NewVec2(0, 0))

			if c := renderPlayListTableCtxMenu(t, r); c != nil {
				delSel = c
			}

			value, ok := hirc.Sounds.Load(p.UniquePlayID)
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
		t.RanSeqPlaylistStorage.ApplyRequests(msIO)

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

func renderPlayListItemOrderCombo(i int, r *wwise.RanSeqCntr) func() {
	var move func() = nil

	preview := strconv.FormatUint(uint64(i), 10)
	label := fmt.Sprintf("PLSequence%d", i)
	if imgui.BeginCombo(label, preview) {
		for j := range r.PlayListItems {
			selected := i == j
			label := strconv.FormatUint(uint64(j), 10)
			if imgui.SelectableBoolPtr(label, &selected) {
				move = bindChangePlayListItemOrder(i, j, r)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}

	return move
}

func renderPlayListTableCtxMenu(t *BankTab, r *wwise.RanSeqCntr) func() {
	var delSel func() = nil

	if imgui.BeginPopupContextItem() {
		if imgui.Button("Delete") {
			delSel = bindPendSelectPlayListItem(t, r)
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	return delSel
}

func bindChangePlayListItemOrder(i, j int, r *wwise.RanSeqCntr) func() {
	return func() {
		r.MovePlayListItem(i, j)
	}
}

func bindPendPlayListItem(t *BankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.RemoveLeafFromPlayList(i)
		t.RanSeqPlaylistStorage.Clear()
	}
}

func bindPendSelectPlayListItem(t *BankTab, r *wwise.RanSeqCntr) func() {
	return func() {
		mut := false
		tids := []uint32{}
		for i, p := range r.PlayListItems {
			if t.RanSeqPlaylistStorage.Contains(imgui.ID(i)) {
				tids = append(tids, p.UniquePlayID)
				mut = true
			}
		}
		if mut {
			t.RanSeqPlaylistStorage.Clear()
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
