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
	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type BankManager struct {
	Banks sync.Map
	WriteLock *atomic.Bool
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
	Bank *wwise.Bank
	InitBank *wwise.Bank
	idQuery string
	sidQuery string
	typeQuery int32
	parentIdQuery string
	parentTypeQuery wwise.HircType
	rewireSidQuery string
	filtered []wwise.HircObj
	filteredParent []wwise.HircObj
	filteredMediaIndexs []*wwise.MediaIndex

	// Storage
	LinearStorage         *imgui.SelectionBasicStorage
	ActiveHirc             wwise.HircObj
	CntrStorage           *imgui.SelectionBasicStorage
	RanSeqPlaylistStorage *imgui.SelectionBasicStorage

	// Sync
	WriteLock             *atomic.Bool
}

func (b *BankTab) changeRoot(hid, np, op uint32) {
	b.Bank.HIRC().ChangeRoot(hid, np, op)
	b.filter()
	b.CntrStorage.Clear()
	b.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) removeRoot(hid, op uint32) {
	b.Bank.HIRC().RemoveRoot(hid, op)
	b.filter()
	b.CntrStorage.Clear()
	b.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) filter() {
	if b.Bank.HIRC() == nil {
		return
	}
	if !utils.IsDigit(b.idQuery) {
		return
	}
	hirc := b.Bank.HIRC()
	i := 0
	old := len(b.filtered)
	for _, h := range hirc.HircObjs {
		// filter out Event and Action
		if h.HircType() == wwise.HircTypeAction || h.HircType() == wwise.HircTypeEvent {
			continue
		}

		// type filter
		if b.typeQuery > 0 && b.typeQuery != int32(h.HircType()) {
			continue
		}

		// sid filter
		sound := wwise.HircTypeSound == h.HircType()
		filterSid := b.typeQuery == 0 || b.typeQuery == int32(wwise.HircTypeSound)
		if filterSid && sound {
			s := h.(*wwise.Sound)
			if !fuzzy.Match(
				b.sidQuery, 
				strconv.FormatUint(uint64(s.BankSourceData.SourceID), 10),
			) {
				continue
			}
		}

		// uid filter
		id, err := h.HircID()
		// Unused bypass
		if err != nil {
			if i < len(b.filtered) {
				b.filtered[i] = h
			} else {
				b.filtered = append(b.filtered, h)
			}
			i += 1
			continue
		}
		if !fuzzy.Match(b.idQuery, strconv.FormatUint(uint64(id), 10)) {
			continue
		}
		if i < len(b.filtered) {
			b.filtered[i] = h
		} else {
			b.filtered = append(b.filtered, h)
		}
		i += 1
	}
	if i < old {
		b.filtered = slices.Delete(b.filtered, i, old)
	}
}

func (b *BankTab) filterParent() {
	if b.Bank.HIRC() == nil {
		return
	}
	if !utils.IsDigit(b.parentIdQuery) {
		return
	}
	hirc := b.Bank.HIRC()
	i := 0
	old := len(b.filteredParent)
	for _, d := range hirc.HircObjs {
		// filter out Event and Action
		if d.HircType() == wwise.HircTypeAction || d.HircType() == wwise.HircTypeEvent {
			continue
		}

		if !slices.Contains(wwise.ContainerHircType, d.HircType()) {
			continue
		}

		if b.parentTypeQuery > 0 && b.parentTypeQuery != d.HircType() {
			continue
		}

		id, err := d.HircID()
		if err != nil {
			if i < len(b.filteredParent) {
				b.filteredParent[i] = d
			} else {
				b.filteredParent = append(b.filteredParent, d)
			}
			i += 1
			continue
		}
		if !fuzzy.Match(b.parentIdQuery, strconv.FormatUint(uint64(id), 10)) {
			continue
		}
		if i < len(b.filteredParent) {
			b.filteredParent[i] = d
		} else {
			b.filteredParent = append(b.filteredParent, d)
		}
		i += 1
	}
	if i < old {
		b.filteredParent = slices.Delete(b.filteredParent, i, old)
	}
}

func (b *BankTab) filterRewireQuery() {
	if b.Bank.DIDX() == nil {
		return
	}
	if !utils.IsDigit(b.rewireSidQuery) {
		return
	}
	didx := b.Bank.DIDX()
	i := 0
	old := len(b.filteredMediaIndexs)
	for _, m := range didx.MediaIndexs {
		if !fuzzy.Match(b.rewireSidQuery, strconv.FormatUint(uint64(m.Sid), 10)) {
			continue
		}
		if i < len(b.filteredMediaIndexs) {
			b.filteredMediaIndexs[i] = m
		} else {
			b.filteredMediaIndexs = append(b.filteredMediaIndexs, m)
		}
		i += 1
	}
	if i < old {
		b.filteredMediaIndexs = slices.Delete(b.filteredMediaIndexs, i, old)
	}
}

func (b *BankTab) encode(ctx context.Context) ([]byte, error) {
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

func (b *BankManager) openBank(ctx context.Context, path string) error {
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

	filtered := []wwise.HircObj{}
	filteredParent := []wwise.HircObj{}
	if bank.HIRC() != nil {
		hirc := bank.HIRC()
		filtered = make([]wwise.HircObj, len(hirc.HircObjs) - int(hirc.ActionCount.Load()) - int(hirc.EventCount.Load()))
		filteredParent = make([]wwise.HircObj, 0, len(hirc.HircObjs) / 2)
		for i, o := range hirc.HircObjs {
			// filter out Event and Action
			if o.HircType() == wwise.HircTypeAction || o.HircType() == wwise.HircTypeEvent {
				continue
			}
			filtered[i] = o
			if slices.Contains(wwise.ContainerHircType, o.HircType()) {
				filteredParent = append(filteredParent, o)
			}
		}
	}

	filteredSid := []*wwise.MediaIndex{}
	if bank.DIDX() != nil {
		didx := bank.DIDX()
		filteredSid = make([]*wwise.MediaIndex, len(didx.MediaIndexs))
		for i, mediaIndex := range didx.MediaIndexs {
			filteredSid[i] = mediaIndex
		}
	}

	t := &BankTab{
		WriteLock: &atomic.Bool{},
		Bank: bank,
		idQuery: "",
		sidQuery: "",
		typeQuery: 0,
		parentIdQuery: "",
		parentTypeQuery: 0,
		filtered: filtered,
		filteredParent: filteredParent,
		filteredMediaIndexs: filteredSid,
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

func (b *BankManager) closeBank(del string) {
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
