package bank_explorer

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type AttenuationFilter struct {
	Id              uint32
	Attenuations []*wwise.Attenuation
}

func (f *AttenuationFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.Attenuations)
	for _, obj := range objs {
		if obj.HircType() != wwise.HircTypeAttenuation {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			slog.Error(
				"Error message before panic",
				"error", "Attenuation struct does not implement HircObj.HircID?",
			)
			panic("Panic Trap")
		}
		if f.Id > 0 && !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Attenuations) {
			f.Attenuations[curr] = obj.(*wwise.Attenuation)
		} else {
			f.Attenuations = append(f.Attenuations, obj.(*wwise.Attenuation))
		}
		curr += 1
	}
	if curr < prev {
		f.Attenuations = slices.Delete(f.Attenuations, curr, prev)
	}
}

type AttenuationViewer struct {
	Filter             AttenuationFilter
	ActiveAttenuation *wwise.Attenuation
}

type FxFilter struct {
	Id     uint32
	Type   wwise.HircType
	Fxs  []wwise.HircObj
}

func (f *FxFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.Fxs)
	for _, obj := range objs {
		if !wwise.FxHircType(obj) {
			continue
		}
		if f.Type > wwise.HircTypeAll && f.Type != obj.HircType() {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			slog.Error(
				"Error message before panic",
				"error", "Attenuation struct does not implement HircObj.HircID?",
			)
			panic("Panic Trap")
		}
		if f.Id > 0 && !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Fxs) {
			f.Fxs[curr] = obj
		} else {
			f.Fxs = append(f.Fxs, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.Fxs = slices.Delete(f.Fxs, curr, prev)
	}
}

type FxViewer struct {
	Filter   FxFilter
	ActiveFx wwise.HircObj
}

type ModulatorFilter struct {
	Id            uint32
	Type          wwise.HircType
	Modulators  []wwise.HircObj
}

func (f *ModulatorFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.Modulators)
	for _, obj := range objs {
		if !wwise.ModulatorType(obj) {
			continue
		}
		if f.Type > wwise.HircTypeAll && f.Type != obj.HircType() {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			slog.Error(
				"Error message before panic",
				"error", "Attenuation struct does not implement HircObj.HircID?",
			)
			panic("Panic Trap")
		}
		if f.Id > 0 && !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Modulators) {
			f.Modulators[curr] = obj
		} else {
			f.Modulators = append(f.Modulators, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.Modulators = slices.Delete(f.Modulators, curr, prev)
	}
}

type ModulatorViewer struct {
	Filter          ModulatorFilter
	ActiveModulator wwise.HircObj
}
