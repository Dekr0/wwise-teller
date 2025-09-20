package wwise

type HIRCSpaceSpec struct {
	Sum             u32
	AdvanceBehavior u32
	AuxParam        u32
	Event           u32
	FXs             u32
	HDR             u32
	Override        u32
	PluginParam     u32
	PositionParam   u32
	Prop            u32
	RProp           u32
	RTPC            u32
	RandomSeq       u32
	Layer           u32
	Source          u32
	State           u32
	StateProp       u32
	StateGroup      u32
}

type HierarchyStat struct {
	State             u32
	Sound             u32
	Action            u32
	Event             u32
	RanSeqCntr        u32
	SwitchCntr        u32
	ActorMixer        u32
	Bus               u32
	LayerCntr         u32
	MusicSegment      u32
	MusicTrack        u32
	MusicSwitchCntr   u32
	MusicRanSeqCntr   u32
	Attenuation       u32
	DialogueEvent     u32
	FxShareSet        u32
	FxCustom          u32
	AuxBus            u32
	LFOModulator      u32
	EnvelopeModulator u32
	AudioDevice       u32
	TimeModulator     u32
}

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
	RandomSeqComponent       RandomSeqComponent
	SourceDataComponent      SourceDataComponent
	StateComponent           StateComponent
	StatePropComponent       StatePropComponent
	StateGroupComponent      StateGroupComponent
}

func EstimateHIRCSpace(in *HierarchyStat, out *HIRCSpaceSpec) {
	actorMixer := in.Sound + in.RanSeqCntr + in.SwitchCntr + in.ActorMixer + in.LayerCntr
	buses := in.Bus + in.AuxBus
	FXs := in.FxCustom + in.FxShareSet
	music := in.MusicTrack + in.MusicSegment + in.MusicSwitchCntr + in.MusicRanSeqCntr
	common := actorMixer + buses + music
	out.AdvanceBehavior = common
	out.AuxParam = common
	out.Event = in.Event
	out.FXs = common
	out.HDR = common
	out.Override = common
	out.PluginParam = in.Sound + FXs
	out.PositionParam = common
	out.Prop = common
	out.RProp = common
	out.RTPC = common
	out.RandomSeq = in.RanSeqCntr
	out.Layer = in.LayerCntr
	out.State = in.State
	out.StateProp = common
	out.StateGroup = common
}

func AllocHIRC(s *HIRCSpaceSpec) *HIRC {
	return &HIRC{
		AdvanceBehaviorComponent: *AllocAdvanceBehaviorComponent(s.AdvanceBehavior), 
		AuxParamComponent: *AllocAuxParamComponent(s.AuxParam),
		EventComponet: *AllocEventComponent(s.Event),
		FXsComponent: *AllocFXsComponent(s.FXs),
		FxMetadatasComponent: *AllocFxMetadataComponent(s.FXs),
		Hierarchy: *AllocHierarchy(s.Sum),
		HDRComponent: *AllocateHDRComponent(s.HDR),
		OverrideComponent: *AllocOverrideComponent(s.Override),
		PluginParamComponent: *AllocPluginParamComponent(s.PluginParam),
		PositionParamComponent: *AllocPositionParamComponent(s.PositionParam),
		PropComponent: *AllocPropComponent(s.Prop),
		RPropComponent: *AllocRPropComponent(s.RProp),
		RTPCComponent: *AllocateRTPCComponent(s.RTPC, s.Layer),
		RandomSeqComponent: *AllocRandomSeqComponent(s.RandomSeq),
		SourceDataComponent: *AllocSourceDataComponent(s.Source),
		StateComponent: *AllocStateComponent(s.State),
		StatePropComponent: *AllocStatePropComponent(s.StateProp),
		StateGroupComponent: *AllocStateGroupComponent(s.StateGroup),
	}
}

// Has side effect
func (h *HIRC) AddState(s State) {
	internalId := h.Hierarchy.AddHierarchyNode(s.Id, HircTypeState)
	h.StateComponent.AddStateData(internalId, s.StateProps)
}

// Has side effect
func (h *HIRC) AddSound(data Sound, version u32) {
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeSound)
	h.SourceDataComponent.AddSourceData(internalId, data.SourceData)
	if SourceHasPluginParam(data.SourceData) {
		h.PluginParamComponent.AddPluginParam(internalId, data.PluginParam)
	}
	h.AddBaseParameter(internalId, data.BaseParameter, version)
}

// Has side effect
func (h *HIRC) AddEvent(data Event) {
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeEvent)
	h.EventComponet.AddEventData(internalId, data.EventData)
}

// Has side effect
func (h *HIRC) AddActorMixer(data ActorMixer, version u32) {
	internalId := h.Hierarchy.AddHierarchyNode(data.Id, HircTypeActorMixer)
	h.AddBaseParameter(internalId, data.BaseParameter, version)
	h.AddContainer(internalId, data.Container)
}
