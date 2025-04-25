package log

import (
	"container/ring"
	"io"
)

const DefaultSize = 1000

type InMemoryLog struct {
	io.Writer
	Logs *ring.Ring
}

func (i *InMemoryLog) Write(p []byte) (n int, err error) {
	i.Logs.Value = string(p)
	i.Logs = i.Logs.Next()
	return len(p), nil
}
