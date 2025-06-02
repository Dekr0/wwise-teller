package ui

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
	EventFilter   EventFilter
	ActiveEvent  *wwise.Event
	ActiveAction *wwise.Action   
}
