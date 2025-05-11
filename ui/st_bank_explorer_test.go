package ui

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"testing"
// 
// 	"github.com/Dekr0/wwise-teller/parser"
// 	"github.com/Dekr0/wwise-teller/wwise"
// )

// func TestBuildTree(t *testing.T) {
// 	bnk, err := parser.ParseBank(
// 		"../tests/bnk/content_audio_stratagems_supply_pack.bnk",
// 		context.Background(),
// 	)
// 	if err != nil { t.Fatal(err) }
// 
// 	b := bankTab{bank: bnk}
// 	b.buildTree()
// 	for _, root := range b.roots {
// 		testRenderHircNode(bnk.HIRC().HircObjs, root, 0)
// 	}
// }

// func testRenderHircNode(hircObjs []wwise.HircObj, n *Node, tab int) {
// 	o := hircObjs[n.tid]
// 
// 	var sid string
// 	id, err := o.HircID()
// 	if err != nil {
// 		sid = fmt.Sprintf("Tree Index %d", n.tid)
// 	} else {
// 		sid = strconv.FormatUint(uint64(id), 10)
// 	}
// 	for range tab { fmt.Printf(" ") }
// 	fmt.Printf("%s (%s) (Parent %d)\n", sid, wwise.HircTypeName[o.HircType()], o.ParentID())
// 	for _, l := range n.leafs {
// 		testRenderHircNode(hircObjs, l, tab + 2)
// 	}
// }
