package bank_explorer

import (
	"slices"
	"strconv"

	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type BusFilter struct {
	Id      uint32
	Type    wwise.HircType
	Buses []wwise.HircObj
}

func (f *BusFilter) Filter(objs []wwise.HircObj) {
	curr := 0 
	prev := len(f.Buses)
	for _, obj := range objs {
		if !wwise.BusHircType(obj) {
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
		if curr < len(f.Buses) {
			f.Buses[curr] = obj
		} else {
			f.Buses = append(f.Buses, obj)
		}
		curr += 1
	}
	if curr < prev {
		f.Buses = slices.Delete(f.Buses, curr, prev)
	}
}

type BusViewer struct {
	Filter    BusFilter
	ActiveBus wwise.HircObj
}
