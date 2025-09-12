package wwise

import "sync"

// --- struct definition --- //

type OverrideComponent struct {
	dMu sync.Mutex

	ActorMixerOverrideParentFx         map[u32]u8
	ActorMixerOverrideFxMetadata       map[u32]u8
	ActorMixerOverrideAttachmentParam  map[u32]u8
	ActorMixerOverrideBusId            map[u32]u32
	ActorMixerDirectParentId           map[u32]u32
}
