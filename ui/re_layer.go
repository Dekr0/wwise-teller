package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
)

func renderLayer(layerCntr *wwise.LayerCntr) {
	if imgui.TreeNodeExStr("Layer / Blend Tracks") {
		if imgui.BeginTabBarV("LayerTabBar", DefaultTabFlags) {
			for _, layer := range layerCntr.Layers  {
				if imgui.BeginTabItem(fmt.Sprintf("Layer / Blend Track %d", layer.Id)) {
					imgui.EndTabItem()
				}
			}
			imgui.EndTabBar()
		}
		imgui.TreePop()
	}
}
