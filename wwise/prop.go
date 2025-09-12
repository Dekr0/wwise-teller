package wwise

import "sync"

// --- struct definition --- //

type Prop struct {
	Ids  []u32
	Vals []byte
}

type PropComponent struct {
	dMu sync.Mutex

	ActorMixerProp map[u32]*PropComponent
}

type RProp struct {
	Ids  []u32
	Mins []byte
	Maxs []byte
}

type RPropComponent struct {
	dMu sync.Mutex

	ActorMixerRProp map[u32]*RProp
}
