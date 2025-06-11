// All parser used for decoding HIRC chunks expect an io.Reader that operates
// on in memory buffer. All parser will have a side effect of which it will
// advance the cursor position of the accepted io.Reader.
// All hierarchy project parser only consume all data excluding hierarchy object
// header data (hierarchy object type [u8] and hierarchy data size [u32])

package parser

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync/atomic"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const MaxNumParseRoutine = 6

type ParserResult struct {
	i   uint32
	obj wwise.HircObj
}

type Parser func(uint32, *wio.Reader) wwise.HircObj

func ParseHIRC(ctx context.Context, r *wio.Reader, I uint8, T []byte, size uint32) (
	*wwise.HIRC, error,
) {
	assert.Equal(0, r.Pos(), "Parser for HIRC does not start at byte 0.")

	numHircItem := r.U32Unsafe()

	hirc := wwise.NewHIRC(I, T, numHircItem)

	/* sync signal */
	sem := make(chan struct{}, MaxNumParseRoutine)
	i := uint32(0)
	parsed := atomic.Uint32{}

	slog.Debug("Start scanning through all hierarchy object, and scheduling parser",
		"numHircItem", numHircItem,
	)

	for parsed.Load() < numHircItem {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if i >= numHircItem {
			continue
		}
		eHircType := r.U8Unsafe()
		dwSectionSize := r.U32Unsafe()
		if SkipHircObjType(wwise.HircType(eHircType)) {
			unknown := wwise.NewUnknown(
				wwise.HircType(eHircType),
				dwSectionSize,
				r.ReadNUnsafe(uint64(dwSectionSize), 4),
			)
			hirc.HircObjs[i] = unknown

			i += 1
			parsed.Add(1)
			slog.Debug("Skipped hierarchy object",
				"index", i,
				"hircType", eHircType,
				"dwSectionSize", dwSectionSize,
				"readerPosition", r.Pos(),
			)
			continue
		}
		select {
		case sem <- struct{}{}:
			switch wwise.HircType(eHircType) {
			case wwise.HircTypeState:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseState,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeSound:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseSound,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeAction:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseAction,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeEvent:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseEvent,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeRanSeqCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseRanSeqCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeSwitchCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseSwitchCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeActorMixer:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseActorMixer,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeBus:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseBus,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeLayerCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseLayerCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicSegment:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicSegment,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicTrack:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicTrack,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicSwitchCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicSwitchCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeMusicRanSeqCntr:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseMusicRanSeqCntr,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeAttenuation:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseAttenuation,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeFxShareSet:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseFxShareSet,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeFxCustom:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseFxCustom,
					hirc,
					sem,
					&parsed,
				)
			case wwise.HircTypeAuxBus:
				go ParserRoutine(
					dwSectionSize,
					uint32(i),
					r.NewBufferReaderUnsafe(uint64(dwSectionSize)),
					ParseAuxBus,
					hirc,
					sem,
					&parsed,
				)
			default:
				panic("Assertion Trap")
			}
			i += 1
			slog.Debug(
				fmt.Sprintf("Scheduled %s parser", wwise.HircTypeName[eHircType]),
				"index", i,
				"hircType", eHircType,
				"dwSectionSize", dwSectionSize,
				"readerPosition", r.Pos(),
			)
		default:
			var obj wwise.HircObj
			switch wwise.HircType(eHircType) {
			case wwise.HircTypeState:
				obj = ParseState(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeSound:
				obj = ParseSound(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeAction:
				obj = ParseAction(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeEvent:
				obj = ParseEvent(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeRanSeqCntr:
				obj = ParseRanSeqCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeSwitchCntr:
				obj = ParseSwitchCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeActorMixer:
				obj = ParseActorMixer(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeBus:
				obj = ParseBus(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeLayerCntr:
				obj = ParseLayerCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicSegment:
				obj = ParseMusicSegment(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicTrack:
				obj = ParseMusicTrack(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicSwitchCntr:
				obj = ParseMusicSwitchCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeMusicRanSeqCntr:
				obj = ParseMusicRanSeqCntr(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeAttenuation:
				obj = ParseAttenuation(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeFxShareSet:
				obj = ParseFxShareSet(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeFxCustom:
				obj = ParseFxCustom(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			case wwise.HircTypeAuxBus:
				obj = ParseAuxBus(dwSectionSize, r.NewBufferReaderUnsafe(uint64(dwSectionSize)))
			default:
				panic("Assertion Trap")
			}
			AddHircObj(hirc, uint32(i), obj)
			i += 1
			parsed.Add(1)
		}
	}

	assert.Equal(
		size,
		uint32(r.Pos()),
		"There are data that is not consumed after parsing all HIRC blob",
	)

	return hirc, nil
}

// Side effect: It will modify HIRC. Specifically, HIRC.HircObjs and maps for
// different types of hierarchy objects.
func AddHircObj(h *wwise.HIRC, i uint32, obj wwise.HircObj) {
	t := obj.HircType()
	id, err := obj.HircID()
	if err != nil {
		panic(err)
	}
	switch t {
	case wwise.HircTypeAudioDevice:
		if _, in := h.AudioDevices.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Audio Device %d.", id))
		}
	case wwise.HircTypeBus:
		if _, in := h.Buses.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Bus %d.", id))
		}
	case wwise.HircTypeAttenuation:
		if _, in := h.Attenuations.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Attenuations %d.", id))
		}
	case wwise.HircTypeFxShareSet:
		if _, in := h.FxShareSets.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Fx Share Set %d.", id))
		}
	case wwise.HircTypeFxCustom:
		if _, in := h.FxCustoms.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Fx Custom %d.", id))
		}
	case wwise.HircTypeAuxBus:
		if _, in := h.AuxBuses.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Aux Bus %d.", id))
		}
	case wwise.HircTypeLFOModulator:
		if _, in := h.LFOModulators.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate LFO Modulator %d.", id))
		}
	case wwise.HircTypeEnvelopeModulator:
		if _, in := h.EnvelopeModulator.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Envelope Modulator %d.", id))
		}
	case wwise.HircTypeTimeModulator:
		if _, in := h.TimeModulator.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Time Modulator %d.", id))
		}
	case wwise.HircTypeSound,
		 wwise.HircTypeRanSeqCntr,
		 wwise.HircTypeSwitchCntr,
		 wwise.HircTypeActorMixer,
		 wwise.HircTypeLayerCntr,
		 wwise.HircTypeDialogueEvent:
		if _, in := h.ActorMixerHirc.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Actor Mixer Hierarchy %d.", id))
		}
	case wwise.HircTypeAction:
		if _, in := h.Actions.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Action %d.", id))
		}
	case wwise.HircTypeEvent:
		if _, in := h.Events.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate Event %d.", id))
		}
	case wwise.HircTypeState:
		if _, in := h.States.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate State %d.", id))
		}
	case wwise.HircTypeMusicSegment,
		 wwise.HircTypeMusicTrack,
		 wwise.HircTypeMusicSwitchCntr,
		 wwise.HircTypeMusicRanSeqCntr:
		if _, in := h.MusicHirc.LoadOrStore(id, obj); in {
			panic(fmt.Sprintf("Duplicate State %d.", id))
		}
	default:
		panic("Panic Trap")
	}
	h.HircObjs[i] = obj
	slog.Debug(fmt.Sprintf("Collected %s parser", wwise.HircTypeName[obj.HircType()]))
}

func SkipHircObjType(t wwise.HircType) bool {
	_, find := sort.Find(len(wwise.KnownHircTypes), func(i int) int {
		if t < wwise.KnownHircTypes[i] {
			return -1
		}
		if t == wwise.KnownHircTypes[i] {
			return 0
		}
		return 1
	})
	return !find
}

func ParserRoutine[T wwise.HircObj](
	size uint32,
	i uint32,
	r *wio.Reader,
	f func(uint32, *wio.Reader) T,
	h *wwise.HIRC,
	sem chan struct{},
	parsed *atomic.Uint32,
) {
	AddHircObj(h, i, f(size, r))
	parsed.Add(1)
	<-sem
}
