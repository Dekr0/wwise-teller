package audio

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/Dekr0/wwise-teller/aio"
	dwav "github.com/go-audio/wav"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

const DefaultResampleQuality = 4
const DefaultSampleRate beep.SampleRate = 48000

const MaxNumStreamers = 8

type Session struct {
	InitLock       atomic.Bool
	ActiveID       uint32
	Mutex          sync.Mutex
	Ctrl          *beep.Ctrl
	Streamers     []Streamer
}

func NewSession() Session {
	return Session{Streamers: make([]Streamer, 0, MaxNumStreamers)}
}

func (s *Session) Busy() bool {
	return s.InitLock.Load()
}

// Lock Session when it's initializing a streamer. When a session is locked, any 
// request of initializing new streamer is denied.
func (s *Session) Lock() {
	s.InitLock.Store(true)
}

// Unlock Session after it finishes initializing a streamer so that new request 
// of initializing new stream is allowed.
func (s *Session) Unlock() {
	s.InitLock.Store(false)
}

func (s *Session) Streamer(id uint32) (Streamer, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	i := slices.IndexFunc(s.Streamers, func(s Streamer) bool {
		return id == s.Id()
	})
	if i == -1 {
		return nil, false
	}
	streamer := s.Streamers[i]
	s.Streamers[0], s.Streamers[i] = streamer, s.Streamers[0]
	return streamer, true
}

func (s *Session) NewSoundStreamerFile(ctx context.Context, id uint32, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	output, df, err := aio.FFMPEGDownsample(ctx, file)
	if err != nil {
		return err
	}
	defer os.Remove(output)
	err = s.NewSoundStreamer(id, f, df)
	return err
}

func (s *Session) NewSoundStreamer(id uint32, r io.Reader, d io.ReadSeekCloser) error {
	defer d.Close()
	s.Mutex.Lock()
	i := slices.IndexFunc(s.Streamers, func(s Streamer) bool {
		return s.Id() == id
	})
	s.Mutex.Unlock()
	if i != -1 {
		return fmt.Errorf("There's already a sound streamer for sound %d", id)
	}

	streamer, format, err := wav.Decode(r)
	if err != nil {
		return err
	}

	decoder := dwav.NewDecoder(d)
	decoder.ReadInfo()
	pcmData := make([][]int64, 0, decoder.NumChans)
	// TODO Deterministic allocation for each channel
	for range decoder.NumChans {
		pcmData = append(pcmData, make([]int64, 0, 256))
	}
	if err := decoder.FwdToPCM(); err != nil {
		return err
	}
	// Buffering?
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}
	c := uint16(0)
	for _, sample := range buf.Data {
		pcmData[c] = append(pcmData[c], int64(sample))
		c += 1
		if c >= decoder.NumChans {
			c = 0
		}
	}

	soundStreamer := SoundStreamer{
		SoundId: id,
		Format: &format,
		Streamer: streamer,
		PCMData: pcmData,
	}

	s.Mutex.Lock()
	if len(s.Streamers) >= MaxNumStreamers {
		last := MaxNumStreamers - 1
		lastStreamer := s.Streamers[last]
		if s.Ctrl != nil && lastStreamer.Id() == s.ActiveID {
			slog.Info("Fallback")
			fallbackStreamer := s.Streamers[last - 1]
			fallbackStreamer.Close()
			s.Streamers = slices.Delete(s.Streamers, last - 1, last)
		} else {
			lastStreamer.Close()
			s.Streamers = slices.Delete(s.Streamers, last, MaxNumStreamers)
		}
	}
	s.Streamers = slices.Insert(s.Streamers, 0, Streamer(&soundStreamer))
	s.Mutex.Unlock()

	return nil
}

func (s *Session) Play(id uint32) error {
	s.Mutex.Lock()
	i := slices.IndexFunc(s.Streamers, func(s Streamer) bool {
		return s.Id() == id
	})
	if i == -1 {
		return nil
	}
	defer s.Mutex.Unlock()

	streamer := s.Streamers[i]
	if s.Ctrl != nil {
		speaker.Lock()
		s.Ctrl.Paused = true
		s.Ctrl = nil
		speaker.Unlock()
	}
	if err := streamer.RewindStart(); err != nil {
		return err
	}
	s.Ctrl = &beep.Ctrl{
		Streamer: beep.Resample(
			DefaultResampleQuality,
			streamer.UWFormat().SampleRate,
			DefaultSampleRate,
			streamer.UWStreamer(),
		),
		Paused: false,
	}
	speaker.Play(s.Ctrl)
	s.ActiveID = streamer.Id()
	return nil
}

func (s *Session) Pause(id uint32) {
	if id == s.ActiveID {
		speaker.Lock()
		s.Ctrl.Paused = true
		speaker.Unlock()
		return
	}
}

func (s *Session) Resume(id uint32) error {
	s.Mutex.Lock()
	i := slices.IndexFunc(s.Streamers, func(s Streamer) bool {
		return s.Id() == id
	})
	if i == -1 {
		s.Mutex.Unlock()
		return nil
	}
	s.Mutex.Unlock()

	streamer := s.Streamers[i]

	if s.Ctrl != nil {
		speaker.Lock()
		s.Ctrl.Paused = true
		s.Ctrl = nil
		speaker.Unlock()
	}
	if streamer.UWStreamer().Position() == streamer.UWStreamer().Len() {
		if err := streamer.RewindStart(); err != nil {
			return err
		}
	}
	s.Ctrl = &beep.Ctrl{
		Streamer: beep.Resample(
			DefaultResampleQuality,
			streamer.UWFormat().SampleRate,
			DefaultSampleRate,
			streamer.UWStreamer(),
		),
		Paused: false,
	}
	speaker.Play(s.Ctrl)
	s.ActiveID = streamer.Id()
	return nil
}

type Streamer interface {
	Close() error
	Id() uint32
	RewindStart() error
	UWStreamer() beep.StreamSeekCloser
	UWFormat() *beep.Format
}

type SoundStreamer struct {
	SoundId        uint32
	Format   *beep.Format
	Streamer  beep.StreamSeekCloser
	PCMData   [][]int64
}

func (s *SoundStreamer) Close() error {
	if err := s.Streamer.Close(); err != nil {
		slog.Error(fmt.Sprintf("Failed to close sound stream for sound %d", s.SoundId), "error", err)
	}
	s.Streamer = nil
	return nil
}

func (s *SoundStreamer) Id() uint32 {
	return s.SoundId
}

func (s *SoundStreamer) RewindStart() error {
	return s.Streamer.Seek(0)
}

func (s *SoundStreamer) UWFormat() *beep.Format {
	return s.Format
}

func (s *SoundStreamer) UWStreamer() beep.StreamSeekCloser {
	return s.Streamer
}

