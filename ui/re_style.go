package ui

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

const DefaultTableFlags = imgui.TableFlagsResizable   |
				          imgui.TableFlagsReorderable |
			              imgui.TableFlagsRowBg       |
			              imgui.TableFlagsBordersH    |
				          imgui.TableFlagsBordersV

const DefaultTableFlagsY = DefaultTableFlags | imgui.TableFlagsScrollY

const DefaultTableSelFlags = imgui.SelectableFlagsSpanAllColumns | 
					 		 imgui.SelectableFlagsAllowOverlap

const DefaultMultiSelectFlags = imgui.MultiSelectFlagsClearOnEscape | 
		                        imgui.MultiSelectFlagsBoxSelect2d

const DefaultTabFlags = imgui.TabBarFlagsReorderable | 
		         imgui.TabBarFlagsTabListPopupButton | 
	             imgui.TabBarFlagsFittingPolicyScroll

var DefaultSize = imgui.NewVec2(0, 0)
