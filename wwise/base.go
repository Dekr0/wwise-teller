package wwise

// --- struct definition --- // 

// This is only used for read. How data is stored is completely different
type BaseParameterS struct {
	OverrideParentFx         u8
	OverrideParentFxMetadata u8
	OverrideAttachmentParam    u8
	SettingVector            u8
	AdvanceSettingVector     u8
	VirtualQueueBehavior     u8
	BelowThresholdBehavior   u8
	MaxNumInstance           u16
	OverideBusId             u32
	DirectParentId           u32
	FXs                     *FXs
	FxMetadatas             *FxMetadatas
	Prop                    *Prop
	RProp                   *RProp
	PositionParam           *PositionParam
	AuxParam                *AuxParam
	StateProp               *StateProp
	StateGroup              *StateGroup
	RTPC                     []RTPC
	RTPCGraph                []RTPCGraph
}
