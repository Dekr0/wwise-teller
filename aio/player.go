package aio

import (
	"os"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

type Player struct {
	Streamer beep.StreamSeekCloser
	cancel   func()
	w        sync.WaitGroup
}

func (p *Player) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if p.cancel != nil {
		p.cancel()
		p.w.Wait()
	}
	var format beep.Format
	p.Streamer, format, err = wav.Decode(f)
	if err != nil {
		return err
	}
	if err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second / 10)); err != nil {
		return err
	}
	return nil
}
