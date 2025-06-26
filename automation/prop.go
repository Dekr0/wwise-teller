package automation

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const BasePropModifierSpecVersion = 0

type BasePropModifierSpec struct {
	Version     uint8           `json:"version"`
	Modifiers []BasePropModifer `json:"modifiers"`
}

type BasePropModifer struct {
	Id                  uint32         `json:"id"`
	RequireProps      []BaseProp       `json:"requireProps"`
	DeleteProps       []wwise.PropType `json:"deleteProps"`
	RequireRangeProps []BaseRangeProp  `json:"requireRangeProps"`
	DeleteRangeProps  []wwise.PropType `json:"deleteRangeProps"`
}

type BaseProp struct {
	Pid wwise.PropType `json:"pid"`
	Val float32        `json:"val"`
}

type BaseRangeProp struct {
	Pid wwise.PropType `json:"pid"`
	Min float32        `json:"min"`
	Max float32        `json:"max"`
}

func ParsePropModifierSpec(fspec string) (*BasePropModifierSpec, error) {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return nil, err
	}
	var spec BasePropModifierSpec
	err = json.Unmarshal(blob, &spec)
	if err != nil {
		return nil, err
	}
	if spec.Version != BasePropModifierSpecVersion {
		return nil, fmt.Errorf("Version spec should be %d!", BasePropModifierSpecVersion)
	}
	return &spec, nil
}

func AutomateBaseProp(bnk *wwise.Bank, fspec string) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	spec, err := ParsePropModifierSpec(fspec)
	if err != nil {
		return nil
	}

	var v any
	var o wwise.HircObj
	var b *wwise.BaseParameter
	buf := make([]byte, 4, 4)
	dbuf := make([]byte, 4, 4)
	for _, m := range spec.Modifiers {
		var in bool
		v, in = h.ActorMixerHirc.Load(m.Id)
		if !in {
			slog.Error(fmt.Sprintf("No actor mixer hierarchy object has ID %d", m.Id))
			continue
		}
		o = v.(wwise.HircObj)
		b = o.BaseParameter()
		if b == nil {
			slog.Error(fmt.Sprintf("%s %d cannot perform base property modification", wwise.PropLabel_140[o.HircType()], m.Id))
			continue
		}
		for _, p := range m.RequireProps {
			if err := wwise.CheckBasePropVal(p.Pid, p.Val); err != nil {
				slog.Error(err.Error())
				continue
			}
			if idx, in := b.PropBundle.HasPid(p.Pid); !in {
				binary.Encode(buf, wio.ByteOrder, &p.Val)
				b.PropBundle.AddWithVal(p.Pid, [4]byte(buf))
			} else {
				b.PropBundle.SetPropByIdxF32(idx, p.Val)
			}
		}
		for _, p := range m.DeleteProps {
			if !slices.Contains(wwise.BasePropType, p) {
				slog.Error(fmt.Sprintf("Invalid base property ID %d", p))
				continue
			}
			b.PropBundle.Remove(p)
		}
		for _, p := range m.RequireRangeProps {
			if err := wwise.CheckBaseRangeProp(p.Pid, p.Min, p.Max); err != nil {
				slog.Error(err.Error())
				continue
			}
			if idx, in := b.RangePropBundle.HasPid(p.Pid); !in {
				binary.Encode(buf, wio.ByteOrder, &p.Min)
				binary.Encode(dbuf, wio.ByteOrder, &p.Max)
				b.RangePropBundle.AddWithVal(p.Pid, [4]byte(buf), [4]byte(dbuf))
			} else {
				b.RangePropBundle.SetPropMinByIdxF32(idx, p.Min)
				b.RangePropBundle.SetPropMinByIdxF32(idx, p.Max)
			}
		}
		for _, p := range m.DeleteRangeProps {
			if !slices.Contains(wwise.BasePropType, p) {
				slog.Error(fmt.Sprintf("Invalid base property ID %d", p))
				continue
			}
			b.RangePropBundle.Remove(p)
		}
	}
	return nil
}
