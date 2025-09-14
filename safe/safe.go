package safe

import "sync"

type Safe[T any] struct {
	Mutex sync.RWMutex
	Data  T
}
