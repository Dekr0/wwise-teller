// TODO:
// - Give the option of including META when saving sound bank 
// 		- Save using integration must exclude META
package bank_explorer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

type BankManager struct {
	Banks        sync.Map
	ActiveBank  *BankTab
	InitBank    *BankTab
	ActivePath   string
	WriteLock    atomic.Bool
}

type BankTab struct {
	Bank                  *wwise.Bank

	// Filter
	MediaIndexFilter       MediaIndexFilter

	ActorMixerViewer       ActorMixerViewer
	AttenuationViewer      AttenuationViewer
	BusViewer              BusViewer
	EventViewer            EventViewer
	FxViewer               FxViewer
	GameSyncViewer         GameSyncViewer
	MusicHircViewer        MusicHircViewer

	// Sync
	WriteLock              atomic.Bool
}

func (b *BankTab) ChangeRoot(hid, np, op uint32) {
	b.Bank.HIRC().ChangeRoot(hid, np, op)
	b.FilterActorMixerHircs()
	b.ActorMixerViewer.CntrStorage.Clear()
	b.ActorMixerViewer.RanSeqPlaylistStorage.Clear()
}

func (b *BankTab) RemoveRoot(hid, op uint32) {
	b.Bank.HIRC().RemoveRoot(hid, op)
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

func (b *BankTab) Encode(ctx context.Context) ([]byte, error) {
	b.WriteLock.Store(true)
	defer b.WriteLock.Store(false)
	type result struct {
		data []byte
		err  error
	}
	c := make(chan *result)
	go func() {
		data, err := b.Bank.Encode(ctx)
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
		bank, err := parser.ParseBank(path, ctx)
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
	actorMixerRoots := []wwise.HircObj{}
	attenuations := []*wwise.Attenuation{}
	buses := []wwise.HircObj{}
	events := []*wwise.Event{}
	fxS := []wwise.HircObj{}
	musicHircs := []wwise.HircObj{}
	musicHircRoots := []wwise.HircObj{}
	states := []*wwise.State{}

	if bank.HIRC() != nil {
		hirc := bank.HIRC()

		c := hirc.HierarchyCount()

		actorMixerHircs = make([]wwise.HircObj, 0, c.ActorMixerHircs)
		actorMixerRoots = make([]wwise.HircObj, 0, c.ActorMixerRoots)
		attenuations = make([]*wwise.Attenuation, 0, c.Attenuations)
		buses = make([]wwise.HircObj, 0, c.Buses)
		events = make([]*wwise.Event, 0, c.Events)
		fxS = make([]wwise.HircObj, 0, c.FxS)
		musicHircs = make([]wwise.HircObj, 0, c.MusicHircs)
		musicHircRoots = make([]wwise.HircObj, 0, c.MusicHircRoots)
		states = make([]*wwise.State, 0, c.States)
		for _, o := range hirc.HircObjs {
			if wwise.ActorMixerHircType(o) {
				actorMixerHircs = append(actorMixerHircs, o)
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
		WriteLock: atomic.Bool{},
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
	}

	t.WriteLock.Store(false)

	b.Banks.Store(path, &t)

	return nil
}

func (b *BankManager) CloseBank(del string) {
	b.Banks.Delete(del)
}
