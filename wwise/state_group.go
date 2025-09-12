package wwise

import (
	"sync"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type StateGroup struct {
	NumStateGroups   uio.V128
	StateGroupId   []u32
	StateSyncType  []u8
	NumStates      []uio.V128
	States         [][]State
}

type StatePropType = u16

type State struct {
	Id           u32
	InstanceId   u32
	PropType     StatePropType 
	PropVal    []byte
}

type StateGroupComponent struct {
	dMu sync.Mutex

	ActorMixerStateGroup map[u32]*StateGroup
}
