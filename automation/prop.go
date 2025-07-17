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
	RequirePropIds       []wwise.PropType  `json:"requirePropIds"`
	RequirePropVals      []float32         `json:"requirePropVals"`
	DeleteProps          []wwise.PropType  `json:"deleteProps"`
	RequireRangePropIds  []wwise.PropType  `json:"requireRangePropIds"`
	RequireRangePropVals []BaseRangeProp   `json:"requireRangePropVals"`
	DeleteRangeProps     []wwise.PropType  `json:"deleteRangeProps"`
	// Set this to any negative value if it's unused
	HDRActiveRange         float32         `json:"HDRActiveRange"`
	Id                     uint32          `json:"id"`
}

type BaseRangeProp struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
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

// TODO: All modifiers should be run in parallel.
func ProcessBaseProps(bnk *wwise.Bank, fspec string) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	spec, err := ParsePropModifierSpec(fspec)
	if err != nil {
		return err
	}

	var v any
	var o wwise.HircObj
	var b *wwise.BaseParameter
	buf := make([]byte, 4, 4)
	dbuf := make([]byte, 4, 4)
	ver := int(bnk.BKHD().BankGenerationVersion)
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
			slog.Error(fmt.Sprintf("%s %d cannot perform base property modification", wwise.HircTypeName[o.HircType()], m.Id))
			continue
		}
		for i := range m.RequirePropIds {
			pid, val := m.RequirePropIds[i], m.RequirePropVals[i]
			if err := wwise.CheckBasePropVal(pid, val); err != nil {
				slog.Error(err.Error())
				continue
			}
			if idx, in := b.PropBundle.HasPid(pid, ver); !in {
				binary.Encode(buf, wio.ByteOrder, &val)
				b.PropBundle.AddWithVal(pid, [4]byte(buf), ver)
			} else {
				b.PropBundle.SetPropByIdxF32(idx, val)
			}
		}
		if m.HDRActiveRange >= 0.0 {
			// Enable HDR Envelope and set HDR Active range
			b.SetEnableEnvelope(true, ver)
			i, _ := b.PropBundle.HDRActiveRange(ver)
			if i != -1 {
				b.PropBundle.SetPropByIdxF32(i, m.HDRActiveRange)
			}
		}
		for _, p := range m.DeleteProps {
			if !slices.Contains(wwise.BasePropType, p) {
				slog.Error(fmt.Sprintf("Invalid base property ID %d", p))
				continue
			}
			b.PropBundle.Remove(p, ver)
		}
		for i := range m.RequireRangePropIds {
			pid, p := m.RequireRangePropIds[i], m.RequireRangePropVals[i]
			if err := wwise.CheckBaseRangeProp(pid, p.Min, p.Max); err != nil {
				slog.Error(err.Error())
				continue
			}
			if idx, in := b.RangePropBundle.HasPid(pid, ver); !in {
				binary.Encode(buf, wio.ByteOrder, &p.Min)
				binary.Encode(dbuf, wio.ByteOrder, &p.Max)
				b.RangePropBundle.AddWithVal(pid, [4]byte(buf), [4]byte(dbuf), ver)
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
			b.RangePropBundle.Remove(p, ver)
		}
	}
	return nil
}
