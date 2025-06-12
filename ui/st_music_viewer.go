package ui

import (
	"slices"
	"strconv"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type MusicHircViewer struct {
	MusicHircFilter        MusicHircFilter
	MusicHircRootFilter    MusicHircRootFilter
	ActiveMusicHirc        wwise.HircObj             

	// Storage
	LinearStorage         *imgui.SelectionBasicStorage
	CntrStorage           *imgui.SelectionBasicStorage
}

type MusicHircFilter struct {
	Id                uint32
	Type              wwise.HircType
	MusicHircs      []wwise.HircObj
}

func (f *MusicHircFilter) Filter(objs []wwise.HircObj) {
	curr := 0 
	prev := len(f.MusicHircs)
	for _, obj := range objs {
		if !wwise.MusicHircType(obj) {
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
		if curr < len(f.MusicHircs) {
			f.MusicHircs[curr] = obj
		} else {
			f.MusicHircs = append(f.MusicHircs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.MusicHircs = slices.Delete(f.MusicHircs, curr, prev)
	}
}

type MusicHircRootFilter struct {
	Id                uint32
	Type              wwise.HircType
	MusicHircRoots  []wwise.HircObj
}

func (f *MusicHircRootFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.MusicHircRoots)
	for _, obj := range objs {
		if !wwise.ContainerMusicHircType(obj) {
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
		if curr < len(f.MusicHircRoots) {
			f.MusicHircRoots[curr] = obj
		} else {
			f.MusicHircRoots = append(f.MusicHircRoots, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.MusicHircRoots = slices.Delete(f.MusicHircRoots, curr, prev)
	}
}
