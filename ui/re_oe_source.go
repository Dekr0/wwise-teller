package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"

	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankSourceData(t *bankTab, o *wwise.Sound) {
	if imgui.TreeNodeExStr("Bank Source Data") {
		bsd := o.BankSourceData
		imgui.Text(fmt.Sprintf("Plugin Type ID: %d", bsd.PluginType()))
		imgui.Text(fmt.Sprintf("Plugin Company ID: %d", bsd.Company()))
		imgui.Text("Stream Type: " + wwise.SourceType[bsd.StreamType])

		renderChangeSourceQuery(t, &bsd)
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

func renderChangeSourceQuery(t *bankTab, bsd *wwise.BankSourceData) {
	size := imgui.NewVec2(imgui.ContentRegionAvail().X * 0.30, 128)
	imgui.BeginChildStrV("ChangeSourceChild", size, 0, 0)

	imgui.Text("Filtered by source ID")
	if imgui.InputTextWithHint("##Filtered by source ID", "", &t.rewireSidQuery, 0, nil) {
		t.filterRewireQuery()
	}

	preview := strconv.FormatUint(uint64(bsd.SourceID), 10)
	imgui.Text("Source ID")
	if imgui.BeginComboV("##Source ID", preview, 0) {
		for _, m := range t.filteredMediaIndexs {
			selected := bsd.SourceID == m.Sid
			preview := strconv.FormatUint(uint64(m.Sid), 10)
			if imgui.SelectableBoolPtr(preview, &selected) {
				bsd.SourceID = m.Sid
				bsd.InMemoryMediaSize = m.Size
			}
			if selected {
				imgui.SetItemDefaultFocus()
			}
		}
		imgui.EndCombo()
	}

	imgui.EndChild()
}

func renderChangeSourceTable(t *bankTab) {
	imgui.BeginChildStr("ChangeSourceChildTable")

	const flags = DefaultTableFlags
	outerSize := imgui.NewVec2(0, 0)
	if imgui.BeginTableV("SourceTable", 2, flags, outerSize, 0) {
		imgui.TableSetupColumn("Source ID")
		imgui.TableSetupColumn("Media Size")
		imgui.TableSetupScrollFreeze(0, 1)
		imgui.TableHeadersRow()

		clipper:= imgui.NewListClipper()
		clipper.Begin(int32(len(t.filteredMediaIndexs)))
		for clipper.Step() {
			for n := clipper.DisplayStart(); n < clipper.DisplayEnd(); n++ {
				m := t.filteredMediaIndexs[n]
				
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

