package wwise

// --- struct definition --- //

type AdvanceBehaviorComponent struct {
	// bit 0 - Priority Override Parent
	// bit 1 - Priority Apply Distance Factor
	// bit 2 - Override MIDI Events Behavior
	// bit 3 - Override MIDI Note Tracking
	// bit 4 - Enable MIDI Note Tracking
	// bit 5 - Is MIDI Break Loop On Note Off
	BaseSettingVector      map[u32]u8
	// bit 0 kill newest
	// bit 1 use virtual queue behavior
	// bit 3 ignore parent max number instance
	// bit 4 override parent virtual voice option
	AdvanceSettingVector   map[u32]u8
	// 0x00
	// 0x01 From Elapsed Time 
	VirtualQueueBehavior   map[u32]VirtualQueueBehavior
	MaxNumInstance         map[u32]u16
	// Ox00 Continue To Play
	BelowThresholdBehavior map[u32]BelowThresholdBehavior
}

// --- allocating / freeing --- //

func AllocAdvanceBehaviorComponent(size u32) *AdvanceBehaviorComponent {
	if size <= 0 {
		return &AdvanceBehaviorComponent{
			make(map[u32]u8),
			make(map[u32]u8),
			make(map[u32]VirtualQueueBehavior),
			make(map[u32]u16),
			make(map[u32]BelowThresholdBehavior),
		}
	}
	return &AdvanceBehaviorComponent{
		make(map[u32]u8, size),
		make(map[u32]u8, size),
		make(map[u32]VirtualQueueBehavior, size),
		make(map[u32]u16, size),
		make(map[u32]BelowThresholdBehavior, size),
	}
}

// --- Component Getter and Setter --- //

func (c *AdvanceBehaviorComponent) GetBaseSettingVector(internalId u32) u8 {
	if baseSettingVector, in := c.BaseSettingVector[internalId]; !in {
		panic("Failed to locate base setting vector")
	} else {
		return baseSettingVector
	}
}

func (c *AdvanceBehaviorComponent) GetAdvanceSettingVector(internalId u32) u8 {
	if advanceSettingVector, in := c.AdvanceSettingVector[internalId]; !in {
		panic("Failed to locate advance setting vector")
	} else {
		return advanceSettingVector
	}
}

func (c *AdvanceBehaviorComponent) GetVirtualQueueBehavior(internalId u32) VirtualQueueBehavior {
	if virtualQueueBehavior, in := c.VirtualQueueBehavior[internalId]; !in {
		panic("Failed to locate virtual queue behavior")
	} else {
		return virtualQueueBehavior
	}
}

func (c *AdvanceBehaviorComponent) GetMaxNumInstance(internalId u32) u16 {
	if maxNumInstance, in := c.MaxNumInstance[internalId]; !in {
		panic("Failed to locate max number instance")
	} else {
		return maxNumInstance
	}
}

func (c *AdvanceBehaviorComponent) GetBelowThresholdBehavior(internalId u32) BelowThresholdBehavior {
	if belowThresholdBehavior, in := c.BelowThresholdBehavior[internalId]; !in {
		panic("Failed to locate below threshold behavior")
	} else {
		return belowThresholdBehavior
	}
}

func (c *AdvanceBehaviorComponent) AddBaseSettingVector(internalId u32, b u8) {
	if _, in := c.BaseSettingVector[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.BaseSettingVector[internalId] = b
}

func (c *AdvanceBehaviorComponent) AddAdvanceSettingVector(internalId u32, v u8) {
	if _, in := c.AdvanceSettingVector[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.AdvanceSettingVector[internalId] = v
}

func (c *AdvanceBehaviorComponent) AddVirtualQueueBehavior(internalId u32, v VirtualQueueBehavior) {
	if _, in := c.VirtualQueueBehavior[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.VirtualQueueBehavior[internalId] = v
}

func (c *AdvanceBehaviorComponent) AddMaxNumInstance(internalId u32, m u16) {
	if _, in := c.MaxNumInstance[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.MaxNumInstance[internalId] = m
}

func (c *AdvanceBehaviorComponent) AddBelowThresholdBehavior(internalId u32, b BelowThresholdBehavior) {
	if _, in := c.BelowThresholdBehavior[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.BelowThresholdBehavior[internalId] = b
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetBaseSettingVector(internalId u32) u8 {
	return h.AdvanceBehaviorComponent.GetBaseSettingVector(internalId)
}

func (h *HIRC) GetAdvanceSettingVector(internalId u32) u8 {
	return h.AdvanceBehaviorComponent.GetAdvanceSettingVector(internalId)
}

func (h *HIRC) GetVirtualQueueBehavior(internalId u32) VirtualQueueBehavior {
	return h.AdvanceBehaviorComponent.GetVirtualQueueBehavior(internalId)
}

func (h *HIRC) GetMaxNumInstance(internalId u32) u16 {
	return h.AdvanceBehaviorComponent.GetMaxNumInstance(internalId)
}

func (h *HIRC) GetBelowThresholdBehavior(internalId u32) BelowThresholdBehavior {
	return h.AdvanceBehaviorComponent.GetBelowThresholdBehavior(internalId)
}

// --- constant and definition --- //

type VirtualQueueBehavior   = u8
type BelowThresholdBehavior = u8
