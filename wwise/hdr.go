package wwise

// --- struct definition --- //

type HDRComponent struct {
	// bit 0 Override HDR Envelope
	// bit 1 Override Analysis
	// bit 2 Normalize Loudness
	// bit 3 Enable Envelope
	HDRSettingVector map[u32]u8
}

// --- allocating / freeing --- //

func AllocateHDRComponent(size u32) *HDRComponent {
	if size <= 0 {
		return &HDRComponent{make(map[u32]u8)}
	}
	return &HDRComponent{make(map[u32]u8, size)}
}

// --- component getter and setter --- //

func (c *HDRComponent) GetHDRSettingVector(internalId u32) u8 {
	if h, in := c.HDRSettingVector[internalId]; !in {
		panic("Failed to locate HDR setting vector")
	} else {
		return h
	}
}

func (c *HDRComponent) AddHDRSettingVector(internalId u32, vector u8) {
	if _, in := c.HDRSettingVector[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.HDRSettingVector[internalId] = vector
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetHDRSettingVector(internalId u32) u8 {
	return h.HDRComponent.GetHDRSettingVector(internalId)
}
