package wwise

type BaseParamOverride struct {
	OverrideParentFx         u8
	OverrideAttachmentParam  u8 // <= 145
	OverrideParentFxMetadata u8
	OverrideVector           u8
	OverrideBusId            u32
}

type ParentId struct {
	ParentId u32
}
