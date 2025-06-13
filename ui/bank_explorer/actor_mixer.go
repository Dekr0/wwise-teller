package bank_explorer

import (
	"slices"
	"strconv"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type ActorMixerViewer struct {
	ActorMixerHircFilter ActorMixerHircFilter
	ActorMixerRootFilter ActorMixerRootFilter
	ActiveActorMixerHirc wwise.HircObj

	// Storage
	LinearStorage         *imgui.SelectionBasicStorage
	CntrStorage           *imgui.SelectionBasicStorage
	RanSeqPlaylistStorage *imgui.SelectionBasicStorage
}

type ActorMixerHircFilter struct {
	Id                uint32
	Sid               uint32
	Type              wwise.HircType
	ActorMixerHircs []wwise.HircObj
}

func (f *ActorMixerHircFilter) Filter(objs []wwise.HircObj) {
	curr := 0 
	prev := len(f.ActorMixerHircs)
	for _, obj := range objs {
		if !wwise.ActorMixerHircType(obj) {
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
		if curr < len(f.ActorMixerHircs) {
			f.ActorMixerHircs[curr] = obj
		} else {
			f.ActorMixerHircs = append(f.ActorMixerHircs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.ActorMixerHircs = slices.Delete(f.ActorMixerHircs, curr, prev)
	}
}

type ActorMixerRootFilter struct {
	Id                uint32
	Type              wwise.HircType
	ActorMixerRoots []wwise.HircObj
}

func (f *ActorMixerRootFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.ActorMixerRoots)
	for _, obj := range objs {
		if !wwise.ContainerActorMixerHircType(obj) {
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
		if curr < len(f.ActorMixerRoots) {
			f.ActorMixerRoots[curr] = obj
		} else {
			f.ActorMixerRoots = append(f.ActorMixerRoots, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.ActorMixerRoots = slices.Delete(f.ActorMixerRoots, curr, prev)
	}
}
