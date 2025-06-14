package bank_explorer

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type EventFilter struct {
	Id        uint32
	Events []*wwise.Event
}

func (f *EventFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.Events)
	for _, obj := range objs {
		if obj.HircType() != wwise.HircTypeEvent {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			slog.Error(
				"Error message before panic",
				"error", "Event struct does not implement HircObj.HircID?",
			)
			panic("Panic Trap")
		}
		if !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.Events) {
			f.Events[curr] = obj.(*wwise.Event)
		} else {
			f.Events = append(f.Events, obj.(*wwise.Event))
		}
		curr += 1
	}
	if curr < prev {
		f.Events = slices.Delete(f.Events, curr, prev)
	}
}

type EventViewer struct {
	Filter   EventFilter
	ActiveEvent  *wwise.Event
	ActiveAction *wwise.Action   
}

type StateFilter struct {
	Id        uint32
	States []*wwise.State
}

func (f *StateFilter) Filter(objs []wwise.HircObj) {
	curr := 0
	prev := len(f.States)
	for _, obj := range objs {
		if obj.HircType() != wwise.HircTypeState {
			continue
		}
		id, err := obj.HircID()
		if err != nil {
			slog.Error(
				"Error message before panic",
				"error", "State struct does not implement HircObj.HircID?",
			)
			panic("Panic Trap")
		}
		if !fuzzy.Match(
			strconv.FormatUint(uint64(f.Id), 10),
			strconv.FormatUint(uint64(id), 10),
		) {
			continue
		}
		if curr < len(f.States) {
			f.States[curr] = obj.(*wwise.State)
		} else {
			f.States = append(f.States, obj.(*wwise.State))
		}
		curr += 1
	}
	if curr < prev {
		f.States = slices.Delete(f.States, curr, prev)
	}
}

type GameSyncViewer struct {
	Filter StateFilter
	ActiveState *wwise.State
}
