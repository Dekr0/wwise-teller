package aio

import (
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

var BeepEnable bool = false

const DefaultResampleQuality = 4
const DefaultSampleRate beep.SampleRate = 48000

func InitBeep() error {
	err := speaker.Init(DefaultSampleRate, DefaultSampleRate.N(time.Second / 10))
	if err == nil {
		BeepEnable = true
	}
	return err
}
