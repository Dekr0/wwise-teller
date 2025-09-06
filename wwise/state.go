package wwise

import "sync"

type StateProp struct {
	Ids  []u16
	Vals []f32
}

type StateComponent struct {
	mu sync.Mutex

	StateProps map[u32]*StateProp
}
