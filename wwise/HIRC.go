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
func (h *HIRC) AddSound(data *Sound, version u32) {
	if data == nil {
		panic("Sound data is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeSound)
	h.SourceDataComponent.AddSourceData(internalId, data.SourceData)
	if SourceHasPluginParam(data.SourceData) {
		h.PluginParamComponent.AddPluginParam(internalId, data.PluginParam)
	}
	h.AddBaseParameter(internalId, data.BaseParameter, version)
}

// Has side effect
func (h *HIRC) AddEvent(id u32, data *EventData) {
	if data == nil {
		panic("Event data is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(id, HircTypeEvent)
	h.EventComponet.AddEventData(internalId, data)
}

// Has side effect
func (h *HIRC) AddActorMixer(data *ActorMixer, version u32) {
	if data == nil {
		panic("Actor mixer is nil")
	}
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeActorMixer)
	h.AddBaseParameter(internalId, data.BaseParameter, version)
	h.AddContainer(internalId, data.Container)
}
