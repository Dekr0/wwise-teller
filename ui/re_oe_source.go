package ui

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/cimgui-go/utils"

	be "github.com/Dekr0/wwise-teller/ui/bank_explorer"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderBankSourceData(t *be.BankTab, o *wwise.Sound) {
	if imgui.TreeNodeExStr("Bank Source Data") {
		bsd := &o.BankSourceData

		{
			pluginID := bsd.PluginID
			pluginType := bsd.PluginType()
			pluginCompany := bsd.Company()
			pluginIDFmt, in := wwise.PluginTypeCodecFmt[pluginID]
			if !in {
				pluginIDFmt, in = wwise.PluginTypeFXFmt[pluginID]
				if !in {
					pluginIDFmt = "Unknown"
				}
			}
			pluginTypeFmt := "Unknown"
			if pluginType < uint32(len(wwise.PluginTypeNames)) {
				pluginTypeFmt = wwise.PluginTypeNames[pluginType]
			}
			pluginCompanyFmt, in := wwise.PluginCompanyNames[wwise.PluginCompanyType(pluginCompany)]
			if !in {
				pluginCompanyFmt = "Unknown"
			}
			imgui.SeparatorText("Audio Source Plugin Information")
			imgui.Text(fmt.Sprintf("Plugin Full ID: %d", pluginID))
			imgui.Text(fmt.Sprintf("Plugin Full ID Translation: %s", pluginIDFmt))
			imgui.Text(fmt.Sprintf("Plugin Type ID: %d", pluginType))
			imgui.Text(fmt.Sprintf("Plugin Type ID Translation: %s", pluginTypeFmt))
			if imgui.IsItemHovered() {
				imgui.SetTooltip(
					"If a plugin type ID translation is not \"Codec\", " + 
					"this indicates audio source ID of this sound object might" + 
					" point to a FX hierarchy object",
				)
			}
			imgui.Text(fmt.Sprintf("Plugin Company ID: %d", pluginCompany))
			imgui.Text(fmt.Sprintf("Plugin Company ID Translation: %s", pluginCompanyFmt))
			imgui.Separator()
		}

		curr := int32(bsd.StreamType)
		imgui.Text("Stream Type: ")
		if imgui.IsItemHovered() {
			imgui.SetTooltip("DATA means audio data is stored in the sound bank. Prefetch streams means audio data is completely stored as a separate file in the hard disk. Streaming means a very small amount of audio data is stored in the sound bank, the rest of it is stored as a separate file in the hard disk.")
		}
		imgui.SameLine()
		imgui.SetNextItemWidth(160)
		if imgui.ComboStrarr("##StreamType", &curr, wwise.SourceTypeNames, int32(len(wwise.SourceTypeNames))) {
			bsd.StreamType = wwise.SourceType(curr)
		}

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

func renderChangeSourceQuery(t *be.BankTab, bsd *wwise.BankSourceData) {
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

func renderChangeSourceTable(t *be.BankTab) {
	size := imgui.NewVec2(0, 256)
	imgui.BeginChildStrV("ChangeSourceChildTable", size, 0, 0)

	const flags = DefaultTableFlags | imgui.TableFlagsScrollY
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

