package aio

import (
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

type LRUWEMPlayersManger struct {
	Max             int
	Lock            sync.Mutex
	LRUWEMPlayers []Player
}

func (l *LRUWEMPlayersManger) HasPlayer(path string) (bool) {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	i := slices.IndexFunc(l.LRUWEMPlayers, func(p Player) bool {
		return strings.EqualFold(path, p.Path)
	})
	return i != -1
}

func (l *LRUWEMPlayersManger) Player(path string) (*Player, error) {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	i := slices.IndexFunc(l.LRUWEMPlayers, func(p Player) bool {
		return strings.EqualFold(path, p.Path)
	})
	if i != -1 {
		player := l.LRUWEMPlayers[i]
		l.LRUWEMPlayers = slices.Insert(l.LRUWEMPlayers, 0, player)
		return &l.LRUWEMPlayers[i], nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}
	player := Player{
		Path: path,
		Streamer: streamer,
		Format: format,
	}
	if len(l.LRUWEMPlayers) >= int(l.Max) {
		player = l.LRUWEMPlayers[l.Max - 1]
		player.Streamer.Close()
		l.LRUWEMPlayers = slices.Delete(l.LRUWEMPlayers, l.Max - 1, l.Max - 1)
		l.LRUWEMPlayers = slices.Insert(l.LRUWEMPlayers, 0, player)
	}
	return &player, nil
}

type Player struct {
	Path      string
	Streamer  beep.StreamSeekCloser
	Format    beep.Format
}

/*
func (p *Player) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if p.Streamer != nil {
		p.Streamer.Close()
		p.Format = nil
	}

	var format beep.Format
	p.Streamer, format, err = wav.Decode(f)
	if err != nil {
		return err
	}
	p.Format = &format
	if err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second / 10)); err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Open %s", path), "number of samples", p.Streamer.Len())
	
	return nil
}
*/

func (p *Player) Resume() error {
	return speaker.Resume()
}

func (p *Player) Suspend() error {
	return speaker.Suspend()
}

func (p *Player) Len() int {
	if p.Streamer != nil {
		return p.Streamer.Len()
	}
	return -1
}

func (p *Player) Position() int {
	if p.Streamer != nil {
		speaker.Lock()
		defer speaker.Unlock()
		return p.Streamer.Position()
	}
	return -1
}

func (p *Player) HasStream() bool {
	return p.Streamer != nil
}
