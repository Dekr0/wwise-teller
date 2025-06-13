package log

import "github.com/Dekr0/wwise-teller/log"

type GuiLog struct {
	Log   log.InMemoryLog
	Debug bool
	Info  bool
	Warn  bool
	Error bool
}
