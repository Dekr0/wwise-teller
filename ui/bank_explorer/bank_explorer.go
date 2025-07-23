// TODO:
// - Give the option of including META when saving sound bank
//   - Save using integration must exclude META
package bank_explorer

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/ui/audio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const MaxWEMExportRetrys  = 8
const MaxInitStreamerRetrys = 8

type BankManager struct {
	Banks        sync.Map
	ActiveBank  *BankTab
	InitBank    *BankTab
	SetNextBank *BankTab
	ActivePath   string
	WriteLock    atomic.Bool
}

type BankTabEnum int8

const (
	BankTabNone        BankTabEnum = -1
	BankTabActorMixer  BankTabEnum = 0
	BankTabMusic       BankTabEnum = 1
	BankTabAttenuation BankTabEnum = 2
	BankTabBuses       BankTabEnum = 3
	BankTabFX          BankTabEnum = 4
	BankTabModulator   BankTabEnum = 5
	BankTabEvents      BankTabEnum = 6
	BankTabGameSync    BankTabEnum = 7
)

// TODO: Examine structure size
type BankTab struct {
	Bank              *wwise.Bank

	// Filter
	MediaIndexFilter  MediaIndexFilter

	ActorMixerViewer  ActorMixerViewer
	AttenuationViewer AttenuationViewer
	BusViewer         BusViewer
	EventViewer       EventViewer
	FxViewer          FxViewer
	GameSyncViewer    GameSyncViewer
	ModulatorViewer   ModulatorViewer
	MusicHircViewer   MusicHircViewer

	// Sync
	Focus             BankTabEnum 
	SounBankLock      atomic.Bool
	WEMStreamRef      map[string]uint32
	WEMExportLock     atomic.Bool
	WEMExportCache    sync.Map
	ErrorAudioSources sync.Map
	ErrorStreamers    sync.Map
	Session           audio.Session
}

func (b *BankTab) BusyWEMExport() bool {
	return b.WEMExportLock.Load()
}

func (b *BankTab) LockWEMExport() {
	b.WEMExportLock.Store(true)
}

func (b *BankTab) UnlockWEMExport() {
	b.WEMExportLock.Store(false)
}

func (b *BankTab) UpdateErrorAudioSource(sid uint32) {
	if actual, loaded := b.ErrorAudioSources.LoadOrStore(sid, 1); loaded {
		b.ErrorAudioSources.Store(sid, actual.(int) + 1)
	}
}

func (b *BankTab) UpdateErrorStreamers(id uint32) {
	if actual, loaded := b.ErrorStreamers.LoadOrStore(id, 1); loaded {
		b.ErrorStreamers.Store(id, actual.(int) + 1)
	}
}

func (b *BankTab) Version() int {
	return int(b.Bank.BKHD().BankGenerationVersion)
}

func (b *BankTab) ChangeRoot(hid, np, op uint32) {
	b.Bank.HIRC().ChangeRoot(hid, np, op, true)
	b.FilterActorMixerHircs()
	b.ActorMixerViewer.CntrStorage.Clear()
	b.ActorMixerViewer.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) RemoveRoot(hid, op uint32) {
	b.Bank.HIRC().RemoveRoot(hid, op, true)
	b.FilterActorMixerHircs()
	b.ActorMixerViewer.CntrStorage.Clear()
	b.ActorMixerViewer.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) FilterActorMixerHircs() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.ActorMixerViewer.HircFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterActorMixerRoots() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.ActorMixerViewer.RootFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterAttenuations() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.AttenuationViewer.Filter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterModulator() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.ModulatorViewer.Filter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterBuses() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.BusViewer.Filter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterEvents() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.EventViewer.Filter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterFxS() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.FxViewer.Filter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterMusicHircs() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.MusicHircViewer.HircFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterMusicHircRoots() {
	if b.Bank.HIRC() == nil {
		return
	}
	b.MusicHircViewer.HircRootFilter.Filter(b.Bank.HIRC().HircObjs)
}

func (b *BankTab) FilterMediaIndices() {
	if b.Bank.DIDX() == nil {
		return
	}
	b.MediaIndexFilter.Filter(b.Bank.DIDX().MediaIndexs)
}

func (b *BankTab) SetActiveBus(id uint32) bool {
	if b.Bank.HIRC() == nil {
		return false
	}
	hirc := b.Bank.HIRC()
	if v, ok := hirc.Buses.Load(id); !ok {
		if v, ok := hirc.AuxBuses.Load(id); !ok {
			return false
		} else {
			b.BusViewer.ActiveBus = v.(wwise.HircObj)
			b.OpenBusHircNode(id)
			return true
		}
	} else {
		b.BusViewer.ActiveBus = v.(wwise.HircObj)
		b.OpenBusHircNode(id)
		return true
	}
}

func (b *BankTab) SetActiveFX(id uint32) {
	if b.Bank.HIRC() == nil {
		return
	}
	hirc := b.Bank.HIRC()
	if v, ok := hirc.FxCustoms.Load(id); !ok {
		if v, ok := hirc.FxShareSets.Load(id); !ok {
			return
		} else {
			b.FxViewer.ActiveFx = v.(wwise.HircObj)
			b.Focus = BankTabFX
		}
	} else {
		b.FxViewer.ActiveFx = v.(wwise.HircObj)
		b.Focus = BankTabFX
	}
}

func (b *BankTab) SetActiveActorMixerHirc(id uint32) {
	hirc := b.Bank.HIRC()
	if hirc == nil {
		return
	}
	v, ok := hirc.ActorMixerHirc.Load(id)
	if !ok {
		return
	}
	h := v.(wwise.HircObj)
	b.ActorMixerViewer.ActiveHirc = h
	b.ActorMixerViewer.CntrStorage.Clear()
	b.ActorMixerViewer.RanSeqPlaylistStorage.Clear()
	b.ActorMixerViewer.LinearStorage.Clear()
	b.ActorMixerViewer.SetSelected(id, true)
}

func (b *BankTab) OpenActorMixerHircNode(id uint32) {
	h := b.Bank.HIRC()
	if h == nil {
		panic("An actor mixer hierarchy node attempts to expand but no HIRC chunk is found")
	}
	h.OpenActorMixerHircNode(id)
}

func (b *BankTab) SearchNearestEventAction(id uint32) (find bool) {
	find = false
	h := b.Bank.HIRC()
	actionId := uint32(0)
	h.Actions.Range(func(key, value any) bool {
		action := value.(*wwise.Action)
		if action.IdExt == id {
			actionId = key.(uint32)
			return false
		}
		return true
	})
	if actionId == 0 {
		return find
	}

	eventId := uint32(0)
	h.Events.Range(func(key, value any) bool {
		event := value.(*wwise.Event)
		idx := slices.Index(event.ActionIDs, actionId)
		if idx == -1 {
			return true
		}
		eventId = key.(uint32)
		return false 
	})

	v, ok := h.Actions.Load(actionId)
	if ok {
		b.EventViewer.ActiveAction = v.(*wwise.Action)
		find = true
	}
	v, ok = h.Events.Load(eventId)
	if ok {
		b.EventViewer.ActiveEvent = v.(*wwise.Event)
		find = true
	}

	return find
}

func (b *BankTab) OpenBusHircNode(id uint32) {
	h := b.Bank.HIRC()
	if h == nil {
		panic("An bus hierarchy node attempts to expand but no HIRC chunk is found")
	}
	h.OpenBusHircNode(id)
}

func (b *BankTab) Encode(ctx context.Context, excludeMeta bool) ([]byte, error) {
	b.SounBankLock.Store(true)
	defer b.SounBankLock.Store(false)
	type result struct {
		data []byte
		err  error
	}
	c := make(chan *result)
	go func() {
		data, err := b.Bank.Encode(ctx, excludeMeta, false)
		c <- &result{data, err}
	}()

	select {
	case <- ctx.Done():
		<- c
		return nil, ctx.Err()
	case r := <- c:
		return r.data, r.err
	}
}

func (b *BankManager) OpenBank(ctx context.Context, path string) error {
	type result struct {
		bank *wwise.Bank
		err error
	}

	c := make(chan *result, 1)

	if _, in := b.Banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}
	
	go func() {
		bank, err := parser.ParseBank(path, ctx, false)
		c <- &result{bank, err}
	}()

	var bank *wwise.Bank
	select {
	case res := <- c:
		if res.bank == nil {
			return res.err
		}
		bank = res.bank
	case <- ctx.Done():
		return ctx.Err()
	}

	if _, in := b.Banks.Load(path); in {
		return fmt.Errorf("Sound bank %s is already open", path)
	}

	actorMixerHircs := []wwise.HircObj{}
	actorMixerSelectionHash := map[uint32]uint32{}
	actorMixerRoots := []wwise.HircObj{}
	attenuations := []*wwise.Attenuation{}
	buses := []wwise.HircObj{}
	events := []*wwise.Event{}
	fxS := []wwise.HircObj{}
	modulator := []wwise.HircObj{}
	musicHircs := []wwise.HircObj{}
	musicHircRoots := []wwise.HircObj{}
	states := []*wwise.State{}

	if bank.HIRC() != nil {
		hirc := bank.HIRC()

		c := hirc.HierarchyCount()

		actorMixerHircs = make([]wwise.HircObj, 0, c.ActorMixerHircs)
		actorMixerSelectionHash = make(map[uint32]uint32, c.ActorMixerHircs)
		actorMixerRoots = make([]wwise.HircObj, 0, c.ActorMixerRoots)
		attenuations = make([]*wwise.Attenuation, 0, c.Attenuations)
		buses = make([]wwise.HircObj, 0, c.Buses)
		events = make([]*wwise.Event, 0, c.Events)
		fxS = make([]wwise.HircObj, 0, c.FxS)
		modulator = make([]wwise.HircObj, 0, c.Modulators)
		musicHircs = make([]wwise.HircObj, 0, c.MusicHircs)
		musicHircRoots = make([]wwise.HircObj, 0, c.MusicHircRoots)
		states = make([]*wwise.State, 0, c.States)
		for i, o := range hirc.HircObjs {
			if wwise.ActorMixerHircType(o) {
				actorMixerHircs = append(actorMixerHircs, o)
				id, err := o.HircID()
				if err != nil {
					panic("All actor mixer hierarchy should be able to return a hierarchy ID without error")
				}
				if id == 85813280 {
					fmt.Println(i)
				}
				if _, in := actorMixerSelectionHash[id]; in {
					slog.Warn(fmt.Sprintf("ID %d is duplicated in actor mixer hierarchy", id))
				} else {
					if id == 85813280 {
						fmt.Println(i)
					}
					actorMixerSelectionHash[id] = uint32(i)
				}
				if wwise.ContainerActorMixerHircType(o) {
					actorMixerRoots = append(actorMixerRoots, o)
				}
			} else if wwise.MusicHircType(o) {
				musicHircs = append(musicHircs, o)
				if wwise.ContainerMusicHircType(o) {
					musicHircRoots = append(musicHircRoots, o)
				}
			} else if wwise.BusHircType(o) {
				buses = append(buses, o)
			} else if wwise.FxHircType(o) {
				fxS = append(fxS, o)
			} else if wwise.ModulatorType(o) {
				modulator = append(modulator, o)
			} else {
				switch t := o.(type) {
				case *wwise.Attenuation:
					attenuations = append(attenuations, t)
				case *wwise.Event:
					events = append(events, t)
				case *wwise.State:
					states = append(states, t)
				}
			}
		}
	}

	indices := []*wwise.MediaIndex{}
	if bank.DIDX() != nil {
		didx := bank.DIDX()
		indices = make([]*wwise.MediaIndex, len(didx.MediaIndexs))
		for i, index := range didx.MediaIndexs {
			indices[i] = &index
		}
	}

	t := BankTab{
		Bank: bank,

		MediaIndexFilter: MediaIndexFilter{
			Sid: 0,
			MediaIndices: indices,
		},

		ActorMixerViewer: ActorMixerViewer{
			HircFilter: ActorMixerHircFilter{
				Id: 0,
				Sid: 0,
				Type : wwise.HircTypeAll,
				Hircs: actorMixerHircs,
			},
			RootFilter: ActorMixerRootFilter{
				Id: 0,
				Type: wwise.HircTypeAll,
				Roots: actorMixerRoots,
			},
			LinearStorage: imgui.NewSelectionBasicStorage(),
			SelectionHash: actorMixerSelectionHash,
			CntrStorage: imgui.NewSelectionBasicStorage(),
			RanSeqPlaylistStorage: imgui.NewSelectionBasicStorage(),
		},
		AttenuationViewer: AttenuationViewer{
			Filter: AttenuationFilter{
				Id: 0,
				Attenuations: attenuations,
			},
			ActiveAttenuation: nil,
		},
		BusViewer: BusViewer{
			Filter: BusFilter{
				Id: 0,
				Type: wwise.HircTypeAll,
				Buses: buses,
			},
			ActiveBus: nil,
		},
		EventViewer: EventViewer{
			Filter: EventFilter{
				Id: 0,
				Events: events,
			},
			ActiveEvent: nil,
			ActiveAction: nil,
		},
		FxViewer: FxViewer{
			Filter: FxFilter{
				Id: 0,
				Type: wwise.HircTypeAll,
				Fxs: fxS,
			},
			ActiveFx: nil,
		},
		GameSyncViewer: GameSyncViewer{
			Filter: StateFilter{
				Id: 0,
				States: states,
			},
			ActiveState: nil,
		},
		ModulatorViewer: ModulatorViewer{
			Filter: ModulatorFilter{
				Id: 0,
				Modulators: modulator,
			},
			ActiveModulator: nil,
		},
		MusicHircViewer: MusicHircViewer{
			HircFilter: MusicHircFilter{
				Id: 0,
				Type: wwise.HircTypeAll,
				MusicHircs: musicHircs,
			},
			HircRootFilter: MusicHircRootFilter{
				Id: 0,
				Type: wwise.HircTypeAll,
				MusicHircRoots: musicHircRoots,
			},
			LinearStorage: imgui.NewSelectionBasicStorage(),
			CntrStorage: imgui.NewSelectionBasicStorage(),
		},
		Focus: BankTabNone,
		SounBankLock: atomic.Bool{},
		WEMExportLock: atomic.Bool{},
		Session: audio.NewSession(),
	}

	t.SounBankLock.Store(false)
	t.WEMExportLock.Store(false)

	b.Banks.Store(path, &t)

	return nil
}

func (b *BankManager) CloseBank(del string) {
	if v, in := b.Banks.Load(del); in {
		bnkTab := v.(*BankTab)
		bnkTab.Session.Shutdown()
	}
	b.Banks.Delete(del)
}
