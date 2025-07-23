package automation

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const BasePropModifierSpecVersion    = 0
const RanSeqModifierSpecVersion      = 0
const BulkProcessBasePropSpecVersion = 0

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
	MaxNumInstances        int16           `json:"maxNumInstances"`
	Id                     uint32          `json:"id"`
}

// Mainly for normalizing sound
type BulkBasePropModifierSpec struct {
	Version                uint8          `json:"version"`
	IDs                  []uint32         `json:"ids"`
	RequirePropIds       []wwise.PropType `json:"requirePropIds"`
	RequirePropVals      []float32        `json:"requirePropVals"`
	RequireRangePropIds  []wwise.PropType `json:"requireRangePropIds"`
	RequireRangePropVals []BaseRangeProp  `json:"requireRangePropVals"`
	MaxNumInstances      int16
}

type RanSeqModifierSpec struct {
	Version          uint8          `json:"version"`
	RanSeqModifier []RanSeqModifier `json:"modifiers"`
}

type RanSeqModifier struct {
	Seq                     bool   `json:"Seq"`
	Shuffle                 bool   `json:"Shuffle"`
	Weighted                bool   `json:"Weighted"`
	Continuous              bool   `json:"Continuous"`
	ResetPlayListAtEachPlay bool   `json:"ResetPlayListAtEachPlay"`
	RestartBackward         bool   `json:"RestartBackward"`
	LoopCount               uint16 `json:"LoopCount"`
	Id                      uint32
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

func ProcessBaseProps(bnk *wwise.Bank, fspec string) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	spec, err := ParsePropModifierSpec(fspec)
	if err != nil {
		return err
	}

	if len(spec.Modifiers) <= 0 {
		return nil
	}

	var w sync.WaitGroup
	sem := make(chan struct{}, 4)

	ver := int(bnk.BKHD().BankGenerationVersion)

	proc := func(m *BasePropModifer, main bool) {
		if !main {
			defer func() {
				<- sem
				w.Done()
			}()
		}

		buf := make([]byte, 4, 4)
		dbuf := make([]byte, 4, 4)

		v, in := h.ActorMixerHirc.Load(m.Id)
		if !in {
			slog.Error(fmt.Sprintf("No actor mixer hierarchy object has ID %d", m.Id))
			return
		}

		o := v.(wwise.HircObj)
		b := o.BaseParameter()
		if b == nil {
			slog.Error(fmt.Sprintf("%s %d cannot perform base property modification", wwise.HircTypeName[o.HircType()], m.Id))
			return
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

		if m.MaxNumInstances >= 0 {
			b.AdvanceSetting.SetIgnoreParentMaxNumInst(true)
			b.AdvanceSetting.MaxNumInstance = uint16(m.MaxNumInstances)
		}
	}

	unique := make(map[uint32]struct{}, len(spec.Modifiers))
	for i := range spec.Modifiers {
		m := &spec.Modifiers[i]
		if _, in := unique[m.Id]; in {
			continue
		}
		unique[m.Id] = struct{}{}
		select {
		case sem <- struct{}{}:
			w.Add(1)
			go proc(m, false)
		default:
			proc(m, true)
		}
	}

	w.Wait()

	close(sem)

	return nil
}

func BulkProcessBaseProp(bnk *wwise.Bank, fspec string) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	var spec BulkBasePropModifierSpec
	err := ParseBulkBasePropModifierSpec(&spec, fspec)
	if err != nil {
		return err
	}

	if len(spec.IDs) <= 0 {
		slog.Warn("No hierarchy IDs are provided. Do nothing")
		return nil
	}

	var w sync.WaitGroup
	sem := make(chan struct{}, 4)

	ver := int(bnk.BKHD().BankGenerationVersion)

	proc := func(id uint32, main bool) {
		if !main {
			defer func() {
				<- sem
				w.Done()
			}()
		}

		v, in := h.ActorMixerHirc.Load(id)
		if !in {
			slog.Error(fmt.Sprintf("No hierarchy object has ID %d", id))
			return
		}

		o := v.(wwise.HircObj)
		b := o.BaseParameter()
		if b == nil {
			slog.Error(fmt.Sprintf("Cannot modifiy common property of hierarchy object %d (type %s)", id, wwise.HircTypeName[o.HircType()]))
			return
		}

		buf := make([]byte, 4, 4)
		dbuf := make([]byte, 4, 4)

		for j := range spec.RequirePropIds {
			pid, p := spec.RequireRangePropIds[j], spec.RequireRangePropVals[j]
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

		for j := range spec.RequireRangePropIds {
			pid, p := spec.RequireRangePropIds[j], spec.RequireRangePropVals[j]
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
	}

	unique := make([]uint32, 0, len(spec.IDs))
	for _, id := range spec.IDs {
		in := slices.Contains(unique, id)
		if in {
			continue
		}
		unique = append(unique, id)
		select {
		case sem <- struct{}{}:
			w.Add(1)
			go proc(id, false)
		default:
			proc(id, true)
		}
	}

	w.Wait()

	close(sem)

	return nil
}

func ParseBulkBasePropModifierSpec(spec *BulkBasePropModifierSpec, fspec string) error {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return fmt.Errorf("Failed to open bulk process base property script %s: %w", fspec, err)
	}
	err = json.Unmarshal(blob, spec)
	if err != nil {
		return fmt.Errorf("Failed to decode bulk process base property script %s: %w", fspec, err)
	}
	if spec.Version != BulkProcessBasePropSpecVersion {
		return fmt.Errorf("Version spec should be %d!", BulkProcessBasePropSpecVersion)
	}
	return nil
}

func ParseRanSeqModifierSpec(spec *RanSeqModifierSpec, fspec string) error {
	blob, err := os.ReadFile(fspec)
	if err != nil {
		return fmt.Errorf("Failed to open random / sequence process script %s: %w", fspec, err)
	}
	err = json.Unmarshal(blob, spec)
	if err != nil {
		return fmt.Errorf("Failed to decode random / sequence process script %s: %w", fspec, err)
	}
	if spec.Version != RanSeqModifierSpecVersion {
		return fmt.Errorf("Version spec should be %d!", RanSeqModifierSpecVersion)
	}
	return nil
}

func ProcessRanSeq(bnk *wwise.Bank, fspec string) error {
	h := bnk.HIRC()
	if h == nil {
		return wwise.NoHIRC
	}

	var spec RanSeqModifierSpec
	err := ParseRanSeqModifierSpec(&spec, fspec)
	if err != nil {
		return err
	}

	if len(spec.RanSeqModifier) <= 0 {
		return nil
	}

	var w sync.WaitGroup
	sem := make(chan struct{}, 4)

	proc := func(m *RanSeqModifier, main bool) {
		if !main {
			defer func() {
				<- sem
				w.Done()
			}()
		}

		v, in := h.ActorMixerHirc.Load(uint32(m.Id))
		if !in {
			slog.Error(fmt.Sprintf("No random / sequence container has ID %d", m.Id))
			return
		}

		r := v.(*wwise.RanSeqCntr)
		if m.Seq {
			r.PlayListSetting.Mode = wwise.ModeSequence
			r.ResetPlayListToLeafOrder()
		} else {
			r.PlayListSetting.Mode = wwise.ModeRandom
		}

		if r.PlayListSetting.Random() {
			if m.Shuffle {
				r.PlayListSetting.RandomMode = wwise.RandomModeShuffle
			} else {
				r.PlayListSetting.RandomMode = wwise.RandomModeNormal
			}
			r.PlayListSetting.SetUsingWeight(m.Weighted)
		}

		r.PlayListSetting.SetContinuous(m.Continuous)
		if r.PlayListSetting.Continuous() {
			r.PlayListSetting.LoopCount = m.LoopCount
		}

		r.PlayListSetting.SetResetPlayListAtEachPlay(m.ResetPlayListAtEachPlay)
		r.PlayListSetting.SetRestartBackward(m.RestartBackward)
	}

	unique := make(map[uint32]struct{}, len(spec.RanSeqModifier))
	for i := range spec.RanSeqModifier {
		m := &spec.RanSeqModifier[i]

		if _, in := unique[m.Id]; in {
			continue
		}
		unique[m.Id] = struct{}{} 

		select {
		case sem <- struct{}{}:
			w.Add(1)
			go proc(m ,false)
		default:
			proc(m, true)
		}
	}

	w.Wait()

	close(sem)

	return nil
}
