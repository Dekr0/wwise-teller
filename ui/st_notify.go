package ui

import (
	"time"

	"github.com/Dekr0/wwise-teller/log"
)

type notfiy struct {
	message string
	timer *time.Timer
}

type NotifyQ struct {
	queue []*notfiy
}

type GuiLog struct {
	log *log.InMemoryLog
	debug bool
	info bool
	warn bool
	error bool
}
