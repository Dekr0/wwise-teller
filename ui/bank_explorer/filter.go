package bank_explorer

import (
	"slices"
	"strconv"
	"github.com/Dekr0/wwise-teller/wwise"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

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
