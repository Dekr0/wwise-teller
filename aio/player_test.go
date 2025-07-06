package aio

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

var testWavesDir string = os.Getenv("TESTS")

func TestPlayerOpenResume(t *testing.T) {
	if err := speaker.Init(DefaultSampleRate, DefaultSampleRate.N(time.Second / 10)); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(testWavesDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join(testWavesDir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		done := make(chan bool)
		streamer, format, err := wav.Decode(f)
		if err != nil {
			t.Fatal(err)
		}
		resample := beep.Resample(4, format.SampleRate, DefaultSampleRate, streamer)
		speaker.Play(beep.Seq(resample, beep.Callback(func() { done <- true })))
		<- done
		streamer.Seek(0)
		speaker.Play(beep.Seq(resample, beep.Callback(func() { done <- true })))
		<- done
		return
	}
}
