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
