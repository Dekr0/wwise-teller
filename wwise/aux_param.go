package wwise

import "sync"

// --- struct definition --- //

type AuxParam struct {
	SettingVector    u8
	AuxIds        [4]u32
	ReflectionAux    u32
}

type AuxParamComponent struct {
	dMu sync.Mutex

	ActorMixerAuxParam map[u32]*AuxParam
}
