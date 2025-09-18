package wwise

// --- struct definition --- //

type OverrideComponent struct {
	OverrideParentFx        map[u32]u8
	OverrideFxMetadata      map[u32]u8
	OverrideAttachmentParam map[u32]u8
	OverrideBusId           map[u32]u32
}

type Override struct {
	OverrideParentFx        u8
	OverrideFxMetadata      u8
	OverrideAttachmentParam u8
	OverrideBusId           u32
}

// --- allocation --- //

func AllocOverrideComponent(size u32) (a *OverrideComponent) {
	a = &OverrideComponent{}
	if size <= 0 {
		a.OverrideParentFx = make(map[u32]u8)
		a.OverrideFxMetadata = make(map[u32]u8)
		a.OverrideAttachmentParam = make(map[u32]u8)
		a.OverrideBusId = make(map[u32]u32)
	} else {
		a.OverrideParentFx = make(map[u32]u8, size)
		a.OverrideFxMetadata = make(map[u32]u8, size)
		a.OverrideAttachmentParam = make(map[u32]u8, size)
		a.OverrideBusId = make(map[u32]u32, size)
	}
	return a
}

// --- component getter and setter --- //

func (c *OverrideComponent) GetOverrideParentFx(internalId u32) u8 {
	if o, in := c.OverrideParentFx[internalId]; !in {
		panic("Failed to locate override parent FX.")
	} else {
		return o
	}
}

func (c *OverrideComponent) GetOverrideFxMetadata(internalId u32) u8 {
	if o, in := c.OverrideFxMetadata[internalId]; !in {
		panic("Failed to locate override FX metadata.")
	} else {
		return o
	}
}

func (c *OverrideComponent) GetOverrideAttachmentParam(internalId u32) u8 {
	if o, in := c.OverrideAttachmentParam[internalId]; !in {
		panic("Failed to locate override attachment parameter.")
	} else {
		return o
	}
}

func (c *OverrideComponent) GetOverrideBusId(internalId u32) u32 {
	if o, in := c.OverrideBusId[internalId]; !in {
		panic("Failed to locate override bus id.")
	} else {
		return o
	}
}

// Has side effect
func (c *OverrideComponent) AddOverrideFxMetadata(internalId u32, o u8) {
	if _, in := c.OverrideFxMetadata[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.OverrideFxMetadata[internalId] = o
}

// Has side effect
func (c *OverrideComponent) AddOverrideAttachmentParam(internalId u32, o u8) {
	if _, in := c.OverrideAttachmentParam[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.OverrideAttachmentParam[internalId] = o
}

// Has side effect
func (c *OverrideComponent) AddOverrideBusId(internalId u32, o u32) {
	if _, in := c.OverrideBusId[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.OverrideBusId[internalId] = o
}

// Has side effect
func (c *OverrideComponent) AddOverrideParentFx(internalId u32, o u8) {
	if _, in := c.OverrideParentFx[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.OverrideParentFx[internalId] = o
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetOverrideParentFx(internalId u32) u8 {
	return h.OverrideComponent.GetOverrideParentFx(internalId)
}

func (h *HIRC) GetOverrideFxMetadata(internalId u32) u8 {
	return h.OverrideComponent.GetOverrideFxMetadata(internalId)
}

func (h *HIRC) GetOverrideAttachmentParam(internalId u32) u8 {
	return h.OverrideComponent.GetOverrideAttachmentParam(internalId)
}

func (h *HIRC) GetOverrideBusId(internalId u32) u32 {
	return h.OverrideComponent.GetOverrideBusId(internalId)
}

// --- constant --- //

const SizeOfActorMixerOverrideLE145 = 7
const SizeOfActorMixerOverrideG145  = 6
