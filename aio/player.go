package aio

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

var BeepEnable bool = false

const DefaultMaxLRUWEMPlayerManagerSize = 8
const MaxCreatePlayerRetries = 8
const DefaultResampleQuality = 4
const DefaultSampleRate beep.SampleRate = 48000

func InitBeep() error {
	err := speaker.Init(DefaultSampleRate, DefaultSampleRate.N(time.Second / 10))
	if err == nil {
		BeepEnable = true
	}
	return err
}

type PlayersManager struct {
	CreateLock atomic.Bool
	Max        int32
	Lock       sync.Mutex
	BlackList  map[string]uint32
	Active    *Player
	Players []*Player
}

func (l *PlayersManager) PauseExceptActive() {
	l.Lock.Lock()
	active := l.Active
	for _, p := range l.Players {
		if p != active {
			p.Ctrl.Paused = true
		}
	}
	l.Lock.Unlock()
}

func (l *PlayersManager) SetActivePlayer(path string) {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	i := slices.IndexFunc(l.Players, func(p *Player) bool {
		return strings.EqualFold(path, p.Path)
	})
	if i == -1 {
		return
	}
	if l.Active != nil {
		speaker.Lock()
		l.Active.Ctrl.Paused = true
		speaker.Unlock()
	}
	l.Active = l.Players[i]
}

func (l *PlayersManager) SeekPlayer(path string, seek int) {}

func (l *PlayersManager) LoopActivePlayer() error {
	if l.Active == nil {
		return nil
	}

	if l.Active.Loop {
		l.ResumeActivePlayer()
	}

	var err error

	speaker.Lock()

	l.PauseExceptActive()
	p := l.Active
	p.Ctrl.Paused = true
	p.Ctrl.Streamer = nil
	if err = p.Streamer.Seek(0); err == nil {
		looper, err := beep.Loop2(p.Streamer)
		if err != nil {
			speaker.Unlock()
			return err
		}
		p.Ctrl = &beep.Ctrl{Streamer: looper, Paused: false}
		defer speaker.Play(p.Ctrl)
		p.Loop = true
	}
	speaker.Unlock()

	return err
}

func (l *PlayersManager) PlayActivePlayer() error {
	if l.Active == nil {
		return nil
	}

	var err error

	speaker.Lock()

	l.PauseExceptActive()
	p := l.Active
	p.Ctrl.Paused = true
	if p.Streamer.Position() == p.Streamer.Len() {
		if err = p.Streamer.Seek(0); err == nil {
			p.Resampler = *beep.Resample(
				DefaultResampleQuality,
				p.Format.SampleRate,
				DefaultSampleRate,
				p.Streamer,
			)
			p.Ctrl = &beep.Ctrl{Streamer: &p.Resampler, Paused: false}
			defer speaker.Play(p.Ctrl)
		}
	} else {
		if err = p.Streamer.Seek(0); err == nil {
			p.Ctrl.Paused = false
		}
	}
	speaker.Unlock()

	return err
}

func (l *PlayersManager) PauseActivePlayer() {
	p := l.Active
	if p == nil {
		return
	}
	speaker.Lock()
	p.Ctrl.Paused = true
	speaker.Unlock()
}

func (l *PlayersManager) ResumeActivePlayer() {
	p := l.Active
	if p == nil {
		return
	}
	speaker.Lock()
	p.Ctrl.Paused = false
	speaker.Unlock()
}

func (l *PlayersManager) HasPlayer(path string) int {
	l.Lock.Lock()
	defer l.Lock.Unlock()

	count, in := l.BlackList[path]
	if in && count >= MaxCreatePlayerRetries {
		return -2
	}

	return slices.IndexFunc(l.Players, func(p *Player) bool {
		return strings.EqualFold(path, p.Path)
	})
}

func (l *PlayersManager) NewPlayer(path string) error {
	l.Lock.Lock()
	defer l.Lock.Unlock()

	i := slices.IndexFunc(l.Players, func(p *Player) bool {
		return strings.EqualFold(path, p.Path)
	})
	if i != -1 {
		return fmt.Errorf("There's already an audio player associated with %s.", path)
	}

	c, in := l.BlackList[path]
	if in && c >= MaxCreatePlayerRetries {
		return fmt.Errorf("%s has been black listed because it fails to create a new audio player too many times.", path)
	}

	f, err := os.Open(path)
	if err != nil {
		if !in {
			l.BlackList[path] = 1
		} else {
			l.BlackList[path] = c + 1
		}
		return err 
	}

	s, format, err := wav.Decode(f)
	if err != nil {
		if !in {
			l.BlackList[path] = 1
		} else {
			l.BlackList[path] = c + 1
		}
		return err 
	}

	r := beep.Resample(DefaultResampleQuality, format.SampleRate, DefaultSampleRate, s)
	p := &Player{
		Path: path,
		Streamer: s,
		Resampler: *r,
		Ctrl: &beep.Ctrl{Streamer: r, Paused: true},
		Format: format,
	}

	speaker.Play(p.Ctrl)
	m := int(l.Max)
	if len(l.Players) >= m {
		del := m - 1
		delP := l.Players[del]
		if delP == l.Active {
			del -= 1
			delP = l.Players[del]
		}
		delP.Streamer.Close()
		l.Players = slices.Delete(l.Players, int(del - 1), int(del))
	}
	l.Players = slices.Insert(l.Players, 0, p)
	return nil
}

func (l *PlayersManager) NumPlayers() int {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	return len(l.Players)
}

func (l *PlayersManager) PlayerInfo(path string) (int, int, *beep.Format, bool) {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	i := slices.IndexFunc(l.Players, func(p *Player) bool {
		return strings.EqualFold(p.Path, path)
	})
	if i == -1 {
		return -1, -1, nil, false
	}
	p := l.Players[i]
	speaker.Lock()
	pos := p.Streamer.Position()
	speaker.Unlock()
	return pos, p.Streamer.Len(), &p.Format, true
}

func (l *PlayersManager) Debug() []string {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	infos := []string{}
	for _, p := range l.Players {
		infos = append(infos, p.Path)
	}
	return infos
}

type Player struct {
	Path      string
	Streamer  beep.StreamSeekCloser
	Resampler beep.Resampler
	Ctrl      *beep.Ctrl
	Format    beep.Format
	Loop      bool
}
