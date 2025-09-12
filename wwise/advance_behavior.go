package wwise

import "sync"

// --- struct definition --- //

type AdvanceBehaviorComponent struct {
	d sync.Mutex

	ActorMixerBaseSettingVector        map[u32]u8
	ActorMixerAdvanceSettingVector     map[u32]u8
	ActorMixerVirtualQueueBehavior     map[u32]u8
	ActorMixerMaxNumInstance           map[u32]u16
	ActorMixerBelowThresholdBehavior   map[u32]u8
}
