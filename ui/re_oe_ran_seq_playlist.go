package ui

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderPlayListSetting(t *bankTab, r *wwise.RanSeqCntr) {
	if imgui.TreeNodeExStr("Play list setting") {
		imgui.PushItemWidth(160)

		renderPlayListValue(r.PlayListSetting)
		renderPlayListMode(r.PlayListSetting)
		renderPlayListMisc(r.PlayListSetting)

		imgui.PopItemWidth()

		renderPlayListTableSet(t, r)

		imgui.TreePop()
	}
}

func renderPlayListTableSet(t *bankTab, r *wwise.RanSeqCntr) {
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLTransfer", 3, 0, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthStretch, 0, 0)
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
			t.playListStorage.Clear()
		}

		imgui.TableSetColumnIndex(2)
		renderPlayListTable(t, r)

		imgui.EndTable()
	}
}

func renderPlayListValue(p *wwise.PlayListSetting) {
	if imgui.InputScalar("Loop Count", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopCount)),
	) {
	}

	if imgui.InputScalar("Loop Mod Min", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopModMin)),
	) {
	}

	if imgui.InputScalar("Loop Mod Max", imgui.DataTypeU16, uintptr(utils.Ptr(&p.LoopModMax)),
	) {
	}

	if imgui.InputFloat("Transition Time", &p.TransitionTime) {
	}

	if imgui.InputFloat("Transition Time Mod Min", &p.TransitionTimeModMin) {
	}

	if imgui.InputFloat("Transition Time Mod Max", &p.TransitionTimeModMax) {
	}

	if imgui.InputScalar("Avoid Repeat Count", imgui.DataTypeU16, uintptr(utils.Ptr(&p.AvoidRepeatCount))) {
	}
}

func renderPlayListMode(p *wwise.PlayListSetting) {
	tMode := int32(p.TransitionMode)
	if imgui.ComboStrarr("TransitionMode", &tMode, wwise.TransitionModeString, int32(len(wwise.TransitionModeString))) {
		p.TransitionMode = uint8(tMode)
	}

	rMode := int32(p.RandomMode)
	if imgui.ComboStrarr("Random Mode", &rMode, wwise.RandomModeString, int32(len(wwise.RandomModeString))) {
		p.RandomMode = uint8(rMode)
	}

	mode := int32(p.Mode)
	if imgui.ComboStrarr("Mode", &mode, wwise.PlayListModeString, int32(len(wwise.PlayListModeString))) {
		p.Mode = uint8(mode)
	}
}

func renderPlayListMisc(p *wwise.PlayListSetting) {
	usingWeight := p.UsingWeight()
	if imgui.Checkbox("Using Weight", &usingWeight) {
		p.SetUsingWeight(usingWeight)
	}

	resetPlayListAtEachPlay := p.ResetPlayListAtEachPlay()
	if imgui.Checkbox("Reset Playlist At Each Play", &resetPlayListAtEachPlay) {
		p.SetResetPlayListAtEachPlay(resetPlayListAtEachPlay)
	}

	restartBackward := p.RestartBackward()
	if imgui.Checkbox("Restart Backward", &restartBackward) {
		p.SetRestartBackward(restartBackward)
	}

	continuous := p.Continuous()
	if imgui.Checkbox("Continuous", &continuous) {
		p.SetContinuous(continuous)
	}

	global := p.Global()
	if imgui.Checkbox("Global", &global) {
		p.SetGlobal(global)
	}
}

func renderPlayListPendingTable(t *bankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLPendingCell")
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLPendTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Target ID")
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

			imgui.TableSetColumnIndex(0)
			imgui.PushIDStr(fmt.Sprintf("ToPlayList%d", i))
			if imgui.Button(">") {
				toPlayList = bindToPlayList(t, r, i)
			}
			imgui.PopID()

			imgui.TableNextRow()
			imgui.TableSetColumnIndex(1)
			imgui.Text(strconv.FormatUint(uint64(child), 10))
		}

		imgui.EndTable()

		if toPlayList != nil {
			toPlayList()
		}
	}
	imgui.EndChild()
}

func bindToPlayList(t *bankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.AddLeafToPlayList(i)
		t.playListStorage.Clear()
	}
}

func renderPlayListTable(t *bankTab, r *wwise.RanSeqCntr) {
	imgui.BeginChildStr("PLCell")
	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("PLTable", 4, flags, outerSize, 0) {
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0, 0)
		imgui.TableSetupColumn("Sequence")
		imgui.TableSetupColumn("Target ID")
		imgui.TableSetupColumn("Weight")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		var move func() = nil
		var del func() = nil
		var delSel func() = nil

		flags := imgui.MultiSelectFlagsClearOnEscape | imgui.MultiSelectFlagsBoxSelect2d
		storageSize := t.playListStorage.Size()
		itemCount := int32(len(r.PlayListItems))
		msIO := imgui.BeginMultiSelectV(flags, storageSize, itemCount)
		t.playListStorage.ApplyRequests(msIO)

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

			selected := t.playListStorage.Contains(imgui.ID(i))
			label := strconv.FormatUint(uint64(p.UniquePlayID), 10)
			const flags = DefaultTableSelFlags
			imgui.SetNextItemSelectionUserData(imgui.SelectionUserData(i))
			imgui.SelectableBoolPtrV(label, &selected, flags, imgui.NewVec2(0, 0))

			if c := renderPlayListTableCtxMenu(t, r); c != nil {
				delSel = c
			}

			imgui.TableSetColumnIndex(3)
			imgui.SetNextItemWidth(-1)
			imgui.PushIDStr(fmt.Sprintf("PLWeight%d", i))
			if imgui.InputScalar("", imgui.DataTypeU32, uintptr(utils.Ptr(&p.Weight))) {
			}
			imgui.PopID()
		}

		imgui.EndMultiSelect()
		t.playListStorage.ApplyRequests(msIO)

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

func renderPlayListTableCtxMenu(t *bankTab, r *wwise.RanSeqCntr) func() {
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

func bindPendPlayListItem(t *bankTab, r *wwise.RanSeqCntr, i int) func() {
	return func() {
		r.RemoveLeafFromPlayList(i)
		t.playListStorage.Clear()
	}
}

func bindPendSelectPlayListItem(t *bankTab, r *wwise.RanSeqCntr) func() {
	return func() {
		mut := false
		tids := []uint32{}
		for i, p := range r.PlayListItems {
			if t.playListStorage.Contains(imgui.ID(i)) {
				tids = append(tids, p.UniquePlayID)
				mut = true
			}
		}
		if mut {
			t.playListStorage.Clear()
		}
		r.RemoveLeafsFromPlayList(tids)
	}
}
