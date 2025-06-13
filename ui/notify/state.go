package notify

import (
	"time"
)

type Notfiy struct {
	Msg    string
	Timer *time.Timer
}

type NotifyQ struct {
	Queue []Notfiy
}

func (n *NotifyQ) Q(m string, timeout time.Duration) {
	n.Queue = append(n.Queue, Notfiy{m, time.NewTimer(timeout)})
}
