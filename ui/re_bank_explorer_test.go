package ui

import (
	"context"
	"fmt"
	// "strconv"
	"testing"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
	// "github.com/AllenDang/cimgui-go/imgui"
)

func renderTreeNodeDry(i *int, hircObjs []wwise.HircObj) bool {
	o := hircObjs[*i]
	*i += 1

	freeFloat := false
	
	var s string
	id, err := o.HircID()
	if err != nil {
		s = fmt.Sprintf("(%d) Unknown Object", *i)
	} else {
		s = fmt.Sprintf("(%d) Object %d (Parent %d)", *i, id, o.ParentID())
	}
	fmt.Println(s)

	if o.ParentID() == 0 {
		freeFloat = true
	}

	for j := 0; j < o.NumLeaf(); {
		if !renderTreeNodeDry(i, hircObjs) {
			j += 1
		}
	}

	return freeFloat
}

func _TestDryDrawTreeNode(t *testing.T) {
	bnk, err := parser.ParseBank(
		"../tests/bnk/content_audio_las_7.bnk", context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	hircObj := bnk.HIRC().HircObjs
	for i < len(hircObj) {
		renderTreeNodeDry(&i, hircObj)
	}
}

// func _renderHircTTable(b *bankTab) {
// 	if !imgui.BeginTableV("LinearTable", 2, 
// 		imgui.TableFlagsResizable | imgui.TableFlagsReorderable | 
// 		imgui.TableFlagsRowBg | 
// 		imgui.TableFlagsBordersH | imgui.TableFlagsBordersV |
// 		imgui.TableFlagsScrollY,
// 		imgui.Vec2{X: 0.0, Y: 0.0}, 0,
// 	) {
// 		return
// 	}
// 	imgui.TableSetupColumn("Hierarchy ID")
// 	imgui.TableSetupColumn("Hierarchy Type")
// 	imgui.TableSetupScrollFreeze(0, 1)
// 	imgui.TableHeadersRow()
// 
// 	hircObjs := b.bank.HIRC().HircObjs
// 	c := imgui.NewListClipper()
// 	c.Begin(int32(len(hircObjs)))
// 	i := 0
// 	di := int32(0)
// 	for c.Step() {
// 		for di < c.DisplayEnd() && i < len(b.roots) {
// 			_renderHircNode(c, &di, hircObjs, b.roots[i])
// 			i += 1
// 		}
// 	}
// 	for i < len(b.roots) {
// 		_renderHircNode(c, &di, hircObjs, b.roots[i])
// 		i += 1
// 	}
// 	imgui.EndTable()
// }

// func _renderHircNode(
// 	c *imgui.ListClipper,
// 	di *int32,
// 	hircObjs []wwise.HircObj,
// 	n *Node,
// ) {
// 	o := hircObjs[n.tid]
// 	visible := *di >= c.DisplayStart() && *di < c.DisplayEnd()
// 	*di += 1
// 
// 	var sid string
// 	id, err := o.HircID()
// 	if err != nil {
// 		sid = fmt.Sprintf("Tree Index %d", n.tid)
// 	} else {
// 		sid = strconv.FormatUint(uint64(id), 10)
// 	}
// 
// 	if visible {
// 		imgui.SetNextItemStorageID(imgui.ID(n.tid))
// 		imgui.TableNextRow()
// 		imgui.TableSetColumnIndex(0)
// 		open := imgui.TreeNodeExStrV(sid, imgui.TreeNodeFlagsSpanAllColumns)
// 		imgui.TableSetColumnIndex(1)
// 		imgui.Text(wwise.HircTypeName[o.HircType()])
// 		if open {
// 			for _, l := range n.leafs {
// 				_renderHircNode(c, di, hircObjs, l)
// 			}
// 			imgui.TreePop()
// 		}
// 	} else if len(n.leafs) > 0 {
// 		// clipped
// 		if imgui.StateStorage().Int(imgui.ID(n.tid)) != 0 { // open?
// 			imgui.TreePushStr(sid)
// 			for _, l := range n.leafs {
// 				_renderHircNode(c, di, hircObjs, l)
// 			}
// 			imgui.TreePop()
// 		}
// 	}
// }

