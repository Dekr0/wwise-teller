package bank_explorer

import (
	"slices"
	"strconv"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type ActorMixerViewer struct {
	HircFilter ActorMixerHircFilter
	RootFilter ActorMixerRootFilter
	ActiveHirc wwise.HircObj

	// Storage
	LinearStorage         *imgui.SelectionBasicStorage
	CntrStorage           *imgui.SelectionBasicStorage
	RanSeqPlaylistStorage *imgui.SelectionBasicStorage
}

type ActorMixerHircFilter struct {
	Id    uint32
	Sid   uint32
	Type  wwise.HircType
	Hircs []wwise.HircObj
}

func (f *ActorMixerHircFilter) Filter(objs []wwise.HircObj) {
	curr := 0 
	prev := len(f.Hircs)
	for _, obj := range objs {
		if !wwise.ActorMixerHircType(obj) {
			continue
		}
		if f.Type > wwise.HircTypeAll && f.Type != obj.HircType() {
			continue
		}
		sound := obj.HircType() == wwise.HircTypeSound
		bySid := f.Type == 0 || f.Type == wwise.HircTypeSound
		if sound && bySid {
			sound := obj.(*wwise.Sound)
			if f.Sid > 0 && !fuzzy.Match(
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
		if f.Id > 0 && !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Hircs) {
			f.Hircs[curr] = obj
		} else {
			f.Hircs = append(f.Hircs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.Hircs = slices.Delete(f.Hircs, curr, prev)
	}
}

type ActorMixerRootFilter struct {
	Id    uint32
	Type  wwise.HircType
	Roots []wwise.HircObj
}

func (f *ActorMixerRootFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.Roots)
	for _, obj := range objs {
		if !wwise.ContainerActorMixerHircType(obj) {
			continue
		}
		if f.Type > wwise.HircTypeAll && f.Type != obj.HircType() {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			continue
		}
		if f.Id > 0 && !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Roots) {
			f.Roots[curr] = obj
		} else {
			f.Roots = append(f.Roots, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.Roots = slices.Delete(f.Roots, curr, prev)
	}
}
