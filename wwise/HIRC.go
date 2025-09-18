package wwise

type HIRC struct {
	AdvanceBehaviorComponent AdvanceBehaviorComponent
	AuxParamComponent        AuxParamComponent
	EventComponet            EventComponent
	FXsComponent             FXsComponent
	FxMetadatasComponent     FxMetadatasComponent
	HDRComponent             HDRComponent
	Hierarchy                Hierarchy
	OverrideComponent        OverrideComponent
	PluginParamComponent     PluginParamComponent
	PositionParamComponent   PositionParamComponent
	PropComponent            PropComponent
	RPropComponent           RPropComponent
	RTPCComponent            RTPCComponent
	SourceDataComponent      SourceDataComponent
	StateComponent           StateComponent
	StatePropComponent       StatePropComponent
	StateGroupComponent      StateGroupComponent
}

func AllocHIRC(numHirc u32) *HIRC {
	return &HIRC{
		AdvanceBehaviorComponent: *AllocAdvanceBehaviorComponent(numHirc / 2), 
		AuxParamComponent: *AllocAuxParamComponent(numHirc),
		EventComponet: *AllocEventComponent(numHirc / 4),
		FXsComponent: *AllocFXsComponent(0),
		FxMetadatasComponent: *AllocFxMetadataComponent(0),
		Hierarchy: *AllocHierarchy(numHirc),
		HDRComponent: *AllocateHDRComponent(numHirc),
		OverrideComponent: *AllocOverrideComponent(0),
		PluginParamComponent: *AllocPluginParamComponent(0),
		PositionParamComponent: *AllocPositionParamComponent(numHirc),
		PropComponent: *AllocPropComponent(numHirc),
		RPropComponent: *AllocRPropComponent(numHirc),
		RTPCComponent: *AllocateRTPCComponent(numHirc, 0),
		SourceDataComponent: *AllocSourceDataComponent(0),
		StateComponent: *AllocStateComponent(0),
		StatePropComponent: *AllocStatePropComponent(numHirc),
		StateGroupComponent: *AllocStateGroupComponent(numHirc),
	}
}

// Has side effect
func (h *HIRC) AddState(id u32, data *StateHierarchyProp) {
	if data == nil {
		panic("State property is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(id, HircTypeState)
	h.StateComponent.AddStateData(internalId, data)
}

// Has side effect
func (h *HIRC) AddSound(data *SoundH, version u32) {
	if data == nil {
		panic("Sound data is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeSound)
	h.SourceDataComponent.AddSourceData(internalId, data.SourceData)
	h.PluginParamComponent.AddPluginParam(internalId, data.PluginParam)
	h.AddBaseParameter(internalId, data.BaseParameter, version)
}

func (h *HIRC) AddBaseParameter(internalId u32, b *BaseParameter, version u32) {
	h.OverrideComponent.AddOverrideParentFx(internalId, b.OverrideParentFx)
	h.FXsComponent.AddFXs(internalId, b.FXs)
	h.OverrideComponent.AddOverrideFxMetadata(internalId, b.OverrideFxMetadata)
	h.FxMetadatasComponent.AddFxMetadatas(internalId, b.FxMetadatas)
	if version <= 145 {
		h.OverrideComponent.AddOverrideAttachmentParam(internalId, b.OverrideAttachmentParam)
	}
	h.OverrideComponent.AddOverrideBusId(internalId, b.OverideBusId)
	h.Hierarchy.AddDirectParentId(internalId, b.DirectParentId)
	h.AdvanceBehaviorComponent.AddBaseSettingVector(internalId, b.SettingVector)
	h.PropComponent.AddProp(internalId, b.Prop)
	h.RPropComponent.AddRProp(internalId, b.RProp)
	h.PositionParamComponent.AddPositionParam(internalId, b.PositionParam)
	h.AuxParamComponent.AddAuxParam(internalId, b.AuxParam)
	h.AdvanceBehaviorComponent.AddAdvanceSettingVector(internalId, b.AdvanceSettingVector)
	h.AdvanceBehaviorComponent.AddVirtualQueueBehavior(internalId, b.VirtualQueueBehavior)
	h.AdvanceBehaviorComponent.AddMaxNumInstance(internalId, b.MaxNumInstance)
	h.AdvanceBehaviorComponent.AddBelowThresholdBehavior(internalId, b.BelowThresholdBehavior)
	h.HDRComponent.AddHDRSettingVector(internalId, b.HDRSettingVector)
	h.StatePropComponent.AddStateProp(internalId, b.StateProp)
	h.StateGroupComponent.AddStateGroup(internalId, b.StateGroup)
	h.RTPCComponent.AddBaseRTPC(internalId, b.RTPC)
}

// Has side effect
func (h *HIRC) AddEvent(id u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(id, HircTypeEvent)
	h.EventComponet.AddEventData(internalId, data)
}
