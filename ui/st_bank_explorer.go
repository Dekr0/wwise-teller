// TODO:
// - Give the option of including META when saving sound bank 
// 		- Save using integration must exclude META
package ui

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type BankManager struct {
	Banks        sync.Map
	ActiveBank  *BankTab
	ActivePath   string
	WriteLock    atomic.Bool
}

type HircFilter struct {
	Id         uint32
	Sid        uint32
	Type       wwise.HircType
	HircObjs []wwise.HircObj
}

func (f *HircFilter) Filter(objs []wwise.HircObj) {
	curr := 0 
	prev := len(f.HircObjs)
	for _, obj := range objs {
		if wwise.NonHircType(obj) {
			continue
		}
		if f.Type > 0 && f.Type != obj.HircType() {
			continue
		}
		sound := obj.HircType() == wwise.HircTypeSound
		bySid := f.Type == 0 || f.Type == wwise.HircTypeSound
		if sound && bySid {
			sound := obj.(*wwise.Sound)
			if !fuzzy.Match(
				strconv.FormatUint(uint64(f.Sid), 10),
				strconv.FormatUint(uint64(sound.BankSourceData.SourceID), 10),
			) {
				continue
			}
		}
		id, err := obj.HircID()
		if err != nil {
			continue
		}
		if !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.HircObjs) {
			f.HircObjs[curr] = obj
		} else {
			f.HircObjs = append(f.HircObjs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.HircObjs = slices.Delete(f.HircObjs, curr, prev)
	}
}

type HircRootFilter struct {
	Id         uint32
	Type       wwise.HircType
	HircObjs []wwise.HircObj
}

func (f *HircRootFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.HircObjs)
	for _, obj := range objs {
		if wwise.NonHircType(obj) {
			continue
		}
		if !wwise.ContainerHircType(obj) {
			continue
		}
		if f.Type > 0 && f.Type != obj.HircType() {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			continue
		}
		if !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.HircObjs) {
			f.HircObjs[curr] = obj
		} else {
			f.HircObjs = append(f.HircObjs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.HircObjs = slices.Delete(f.HircObjs, curr, prev)
	}
}

type MediaIndexFilter struct {
	Sid             uint32
	MediaIndices []*wwise.MediaIndex
}

func (f *MediaIndexFilter) Filter(indices []wwise.MediaIndex) {
	curr := 0
	prev := len(f.MediaIndices)
	for _, index := range indices {
		if !fuzzy.Match(
			strconv.FormatUint(uint64(f.Sid), 10),
			strconv.FormatUint(uint64(index.Sid), 10),
		) {
			continue
		}
		if curr < len(f.MediaIndices) {
			f.MediaIndices[curr] = &index
		} else {
			f.MediaIndices = append(f.MediaIndices, &index)
		}
		curr += 1
	}
	if curr < prev {
		f.MediaIndices = slices.Delete(f.MediaIndices, curr, prev)
	}
}

type BankTab struct {
	Bank                  *wwise.Bank
	InitBank              *wwise.Bank

	// Filter
	HircFilter             HircFilter
	HircRootFilter         HircRootFilter
	MediaIndexFilter       MediaIndexFilter

	AttenuationViewer      AttenuationViewer
	EventViewer            EventViewer
	GameSyncViewer         GameSyncViewer

	// Storage
	ActiveHirc             wwise.HircObj
	LinearStorage         *imgui.SelectionBasicStorage


	CntrStorage           *imgui.SelectionBasicStorage
	RanSeqPlaylistStorage *imgui.SelectionBasicStorage

	// Sync
	WriteLock              atomic.Bool
}

func (b *BankTab) ChangeRoot(hid, np, op uint32) {
	b.Bank.HIRC().ChangeRoot(hid, np, op)
	b.FilterHircs()
	b.CntrStorage.Clear()
	b.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) RemoveRoot(hid, op uint32) {
	b.Bank.HIRC().RemoveRoot(hid, op)
	b.FilterHircs()
	b.CntrStorage.Clear()
	b.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) FilterHircs() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.HircFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterRoots() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.HircRootFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterEvents() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.EventViewer.EventFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterMediaIndices() {
	if b.Bank.DIDX() == nil {
		return
	}
	b.MediaIndexFilter.Filter(b.Bank.DIDX().MediaIndexs)
}

func (b *BankTab) Encode(ctx context.Context) ([]byte, error) {
	b.WriteLock.Store(true)
	defer b.WriteLock.Store(false)
	type result struct {
		data []byte
		err  error
	}
	c := make(chan *result)
	go func() {
		data, err := b.Bank.Encode(ctx)
		c <- &result{data, err}
	}()

	select {
	case <- ctx.Done():
		<- c
		return nil, ctx.Err()
	case r := <- c:
		return r.data, r.err
	}
}

func (b *BankManager) OpenBank(ctx context.Context, path string) error {
	type result struct {
		bank *wwise.Bank
		err error
	}

	c := make(chan *result, 1)

	if _, in := b.Banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}
	
	go func() {
		bank, err := parser.ParseBank(path, ctx)
		c <- &result{bank, err}
	}()

	var bank *wwise.Bank
	select {
	case res := <- c:
		if res.bank == nil {
			return res.err
		}
		bank = res.bank
	case <- ctx.Done():
		return ctx.Err()
	}

	if _, in := b.Banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}

	objs := []wwise.HircObj{}
	roots := []wwise.HircObj{}
	events := []*wwise.Event{}
	states := []*wwise.State{}
	attenuations := []*wwise.Attenuation{}

	if bank.HIRC() != nil {
		hirc := bank.HIRC()
		objs = make([]wwise.HircObj, 0, hirc.NumHirc())
		roots = make([]wwise.HircObj, 0, len(hirc.HircObjs) / 2)
		for _, o := range hirc.HircObjs {
			if !wwise.NonHircType(o) {
				objs = append(objs, o)
				if wwise.ContainerHircType(o) {
					roots = append(roots, o)
				}
			} else {
				switch t := o.(type) {
				case *wwise.Attenuation:
					attenuations = append(attenuations, t)
				case *wwise.Event:
					events = append(events, t)
				case *wwise.State:
					states = append(states, t)
				}
			}
		}
	}

	indices := []*wwise.MediaIndex{}
	if bank.DIDX() != nil {
		didx := bank.DIDX()
		indices = make([]*wwise.MediaIndex, len(didx.MediaIndexs))
		for i, index := range didx.MediaIndexs {
			indices[i] = &index
		}
	}

	t := &BankTab{
		WriteLock: atomic.Bool{},
		Bank: bank,

		HircFilter: HircFilter{
			Id: 0,
			Sid: 0,
			Type : wwise.HircTypeAll,
			HircObjs: objs,
		},
		HircRootFilter: HircRootFilter{
			Id: 0,
			Type: wwise.HircTypeAll,
			HircObjs: roots,
		},
		MediaIndexFilter: MediaIndexFilter{
			Sid: 0,
			MediaIndices: indices,
		},

		AttenuationViewer: AttenuationViewer{
			AttenuationFilter: AttenuationFilter{
				Id: 0,
				Attenuations: attenuations,
			},
			ActiveAttenuation: nil,
		},
		EventViewer: EventViewer{
			EventFilter: EventFilter{
				Id: 0,
				Events: events,
			},
			ActiveEvent: nil,
			ActiveAction: nil,
		},
		GameSyncViewer: GameSyncViewer{
			StateFilter: StateFilter{
				Id: 0,
				States: states,
			},
			ActiveState: nil,
		},

		ActiveHirc: nil,
		LinearStorage: imgui.NewSelectionBasicStorage(),
		CntrStorage: imgui.NewSelectionBasicStorage(),
		RanSeqPlaylistStorage: imgui.NewSelectionBasicStorage(),
	}
	// t.buildTree()
	t.WriteLock.Store(false)

	b.Banks.Store(path, t)

	return nil
}

func (b *BankManager) CloseBank(del string) {
	b.Banks.Delete(del)
}

// type Node struct {
// 	tid   uint
// 	leafs []*Node
// }
// 
// func (b *bankTab) buildTree() {
// 	hircObjs := b.bank.HIRC().HircObjs
// 	b.roots = []*Node{}
// 	tid := 0
// 	for tid < len(hircObjs) {
// 		root, _ := b.buildRoot(&tid, hircObjs)
// 		b.roots = append(b.roots, root)
// 	}
// }

// func (b *bankTab) buildRoot(tid *int, hircObjs []wwise.HircObj) (*Node, bool) {
// 	o := hircObjs[*tid]
// 	n := &Node{
// 		tid: uint(*tid),
// 		leafs: make([]*Node, 0, o.NumChild()),
// 	}
// 	*tid += 1
// 
// 	rootless := false
// 
// 	if o.ParentID() == 0 { rootless = true }
// 
// 	for j := 0; j < o.NumChild(); {
// 		leaf, rootless := b.buildRoot(tid, hircObjs)
// 		n.leafs = append(n.leafs, leaf)
// 		if !rootless { j += 1 }
// 	}
// 
// 	return n, rootless
// }
