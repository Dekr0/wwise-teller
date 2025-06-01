package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankSourceData(t *BankTab, o *wwise.Sound) {
	if imgui.TreeNodeExStr("Bank Source Data") {
		bsd := &o.BankSourceData
		imgui.Text(fmt.Sprintf("Plugin Type ID: %d", bsd.PluginType()))
		imgui.Text(fmt.Sprintf("Plugin Company ID: %d", bsd.Company()))
		imgui.Text("Stream Type: " + wwise.SourceType[bsd.StreamType])

		renderChangeSourceQuery(t, bsd)
		imgui.SameLine()
		renderChangeSourceTable(t)

		imgui.BeginDisabled()

		languageSpecific := bsd.LanguageSpecific()
		imgui.Checkbox("Language specific", &languageSpecific)

		nonCacheAble := bsd.NonCacheable()
		imgui.Checkbox("Non cacheable", &nonCacheAble)

		imgui.EndDisabled()

		imgui.TreePop()
	}
}

func renderChangeSourceQuery(t *BankTab, bsd *wwise.BankSourceData) {
	size := imgui.NewVec2(imgui.ContentRegionAvail().X * 0.45, 128)
	imgui.BeginChildStrV("ChangeSourceChild", size, 0, 0)

	imgui.Text("Filtered by source ID")
	if imgui.InputScalar("##FilteredBySID", imgui.DataTypeU32, uintptr(utils.Ptr(&t.MediaIndexFilter.Sid))) {
		t.FilterMediaIndices()
	}

	imgui.Text("Source ID")

	var changeSource func() = nil
	preview := strconv.FormatUint(uint64(bsd.SourceID), 10)
	if imgui.BeginComboV("##SourceIDCombo", preview, 0) {
		for _, m := range t.MediaIndexFilter.MediaIndices {
			selected := bsd.SourceID == m.Sid
			preview := strconv.FormatUint(uint64(m.Sid), 10)
			if imgui.SelectableBoolPtr(preview, &selected) {
				changeSource = bindChangeSource(bsd, m.Sid, m.Size)
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}
	if changeSource != nil {
		changeSource()
	}

	imgui.EndChild()
}

func bindChangeSource(bsd *wwise.BankSourceData, sid, inMemorySize uint32) func() {
	return func() {
		bsd.ChangeSource(sid, inMemorySize)
	}
}

func renderChangeSourceTable(t *BankTab) {
	size := imgui.NewVec2(0, 256)
	imgui.BeginChildStrV("ChangeSourceChildTable", size, 0, 0)

	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("SourceTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("Source ID")
		imgui.TableSetupColumn("Media Size")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		clipper:= imgui.NewListClipper()
		clipper.Begin(int32(len(t.MediaIndexFilter.MediaIndices)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				m := t.MediaIndexFilter.MediaIndices[n]
				
				imgui.TableNextRow()

				imgui.TableSetColumnIndex(0)
				imgui.Text(strconv.FormatUint(uint64(m.Sid), 10))

				imgui.TableSetColumnIndex(1)
				imgui.Text(strconv.FormatUint(uint64(m.Size), 10))
			}
		}
		imgui.EndTable()
	}

	imgui.EndChild()
}

