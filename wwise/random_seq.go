package wwise

type RandomSeqSetting struct {
	TransitionMode       u8
	RandomMode           u8
	Mode                 u8
	// bit 0 Is using weight
	// bit 1 Reset playList at each play
	// bit 2 Is restart backward
	// bit 3 Is continuous
	// bit 4 Is global
	SettingVector        u8
	LoopCount            u16
	LoopModMin           u16
	LoopModMax           u16
	AvoidRepeatCount     u16
	TransitionTime       f32
	TransitionTimeModMin f32
	TransitionTimeModMax f32
}

type RandomSeqComponent struct {
	RandomSeqSetting map[u32]*RandomSeqSetting
}

// --- allocating --- //

func AllocRandomSeqComponent(size u32) *RandomSeqComponent {
	if size <= 0 {
		return &RandomSeqComponent{
			make(map[u32]*RandomSeqSetting),
		}
	}
	return &RandomSeqComponent{
		make(map[u32]*RandomSeqSetting, size),
	}
}
