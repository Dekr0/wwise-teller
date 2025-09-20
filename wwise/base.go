package wwise

import "fmt"

// --- struct definition --- //

// This is only used for read. How data is stored is completely different
type BaseParameter struct {
	OverrideParentFx        u8
	OverrideFxMetadata      u8
	OverrideAttachmentParam u8 // <= 145
	SettingVector           u8
	AdvanceSettingVector    u8
	VirtualQueueBehavior    VirtualQueueBehavior
	BelowThresholdBehavior  BelowThresholdBehavior
	HDRSettingVector        u8
	MaxNumInstance          u16
	OverideBusId            u32
	DirectParentId          u32
	FXs                     FXs
	FxMetadatas             FxMetadatas
	Prop                    Prop
	RProp                   RProp
	PositionParam           PositionParam
	AuxParam                AuxParam
	StateProp               StateProp
	StateGroup              StateGroup
	RTPC                    RTPC
}

func (h *HIRC) AddBaseParameter(internalId u32, b BaseParameter, version u32) {
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

func (h *HIRC) BaseParameter(internalId u32, version u32) (b BaseParameter) {
	b = BaseParameter{}
	b.OverrideParentFx = h.GetOverrideParentFx(internalId)
	b.OverrideFxMetadata = h.GetOverrideFxMetadata(internalId)
	if version <= 145 {
		b.OverrideAttachmentParam = h.GetOverrideAttachmentParam(internalId)
	}
	b.SettingVector = h.GetBaseSettingVector(internalId)
	b.AdvanceSettingVector = h.GetAdvanceSettingVector(internalId)
	b.VirtualQueueBehavior = h.GetVirtualQueueBehavior(internalId)
	b.BelowThresholdBehavior = h.GetBelowThresholdBehavior(internalId)
	b.HDRSettingVector = h.GetHDRSettingVector(internalId)
	b.MaxNumInstance = h.GetMaxNumInstance(internalId)
	b.OverideBusId = h.GetOverrideBusId(internalId)
	b.DirectParentId = h.GetDirectParentId(internalId)
	b.FXs = h.GetFXs(internalId)
	b.FxMetadatas = h.GetFxMetadatas(internalId)
	b.Prop = h.GetProp(internalId)
	b.RProp = h.GetRProp(internalId)
	b.AuxParam = h.GetAuxParam(internalId)
	b.PositionParam = h.GetPositionParam(internalId)
	b.StateProp = h.GetStateProp(internalId)
	b.StateGroup = h.GetStateGroup(internalId)
	b.RTPC = h.RTPCComponent.GetBaseRTPC(internalId)

	return b
}

func AssertBaseParameter(b BaseParameter, version u32) error {
	{ // Assertion for each data struct
		if err := AssertFXs(b.FXs); err != nil {
			return fmt.Errorf("FXs assertion failed: %w", err)
		}
		if err := AssertFxMetadatas(b.FxMetadatas); err != nil {
			return fmt.Errorf("Fx Metadata assertion failed: %w", err)
		}
		if err := AssertProp(b.Prop); err != nil {
			return fmt.Errorf("Property assertion failed: %w", err)
		}
		if err := AssertRProp(b.RProp); err != nil {
			return fmt.Errorf("Property assertion failed: %w", err)
		}
		if err := AssertPositionParam(b.PositionParam, b.DirectParentId == 0); err != nil {
			return fmt.Errorf("Position parameter assertion failed: %w", err)
		}
		if err := AssertAuxParam(b.AuxParam); err != nil {
			return fmt.Errorf("Auxiliary parameter assertion failed: %w", err)
		}
		if err := AssertStateProp(b.StateProp); err != nil {
			return fmt.Errorf("State property assertion failed: %w", err)
		}
		if err := AssertStateGroup(b.StateGroup, version); err != nil {
			return fmt.Errorf("State group assertion failed: %w", err)
		}
		if err := AssertRTPC(b.RTPC); err != nil {
			return fmt.Errorf("RTPC assertion failed: %w", err)
		}
	}
	{ // Assertion for relation between data struct. Example, position parameter 
      // can relate to attenuation id property in base parameter.

	}
	return nil
}

func SizeOfBaseParameter(b BaseParameter, version u32) (size u32) {
	if version <= 145 {
		size = 18
	} else {
		size = 17
	}
	size += SizeOfFXs(b.FXs, version) + 
		    SizeOfFxMetadatas(b.FxMetadatas) +
		    SizeOfProp(b.Prop) +
		    SizeOfRProp(b.RProp) +
		    SizeOfPositionParam(b.PositionParam) +
		    SizeOfAuxParam(b.AuxParam) +
		    SizeOfStateProp(b.StateProp) +
		    SizeOfStateGroup(b.StateGroup, version) +
		    SizeOfRTPC(b.RTPC)
	return size
}

func EncodeBaseParameter(e *HircEncoderCtx, b BaseParameter) error {
	curr := e.Encoder.Count
	size := SizeOfBaseParameter(b, e.Version)
	if err := e.Primitive(b.OverrideParentFx); err != nil {
		return fmt.Errorf("Failed to encode override parent FX: %w", err)
	}
	if err := EncodeFXs(e, b.FXs); err != nil {
		return fmt.Errorf("Failed to encode FXs: %w", err)
	}
	if err := e.Primitive(b.OverrideFxMetadata); err != nil {
		return fmt.Errorf("Failed to encode override Fx metadata: %w", err)
	}
	if err := EncodeFxMetadatas(e, b.FxMetadatas); err != nil {
		return fmt.Errorf("Failed to encode Fx metadatas: %w", err)
	}
	if err := e.Primitive(b.OverideBusId); err != nil {
		return fmt.Errorf("Failed to encode overbus id %d: %w", b.OverideBusId, err)
	}
	if err := e.Primitive(b.DirectParentId); err != nil {
		return fmt.Errorf("Failed to encode direct parent id %d: %w", b.DirectParentId, err)
	}
	if err := e.Primitive(b.SettingVector); err != nil {
		return fmt.Errorf("Failed to encode base setting vector %d: %w", b.SettingVector, err)
	}
	if err := EncodeProp(e, b.Prop); err != nil {
		return fmt.Errorf("Failed to encode property: %w", err)
	}
	if err := EncodeRProp(e, b.RProp); err != nil {
		return fmt.Errorf("Failed to encode range based property: %w", err)
	}
	if err := EncodePositionParam(e, b.PositionParam); err != nil {
		return fmt.Errorf("Failed to encode position parameter: %w", err)
	}
	if err := EncodeAuxParam(e, b.AuxParam); err != nil {
		return fmt.Errorf("Failed to encode auxiliary parameter: %w", err)
	}
	if err := e.Primitive(b.AdvanceSettingVector); err != nil {
		return fmt.Errorf("Failed to encode advance setting vector %d: %w", b.AdvanceSettingVector, err)
	}
	if err := e.Primitive(b.VirtualQueueBehavior); err != nil {
		return fmt.Errorf("Failed to encode virtual queue behavior %d: %w", b.VirtualQueueBehavior, err)
	}
	if err := e.Primitive(b.MaxNumInstance); err != nil {
		return fmt.Errorf("Failed to encode max number instance %d: %w", b.MaxNumInstance, err)
	}
	if err := e.Primitive(b.BelowThresholdBehavior); err != nil {
		return fmt.Errorf("Failed to encode below threshold behavior %d: %w", b.BelowThresholdBehavior, err)
	}
	if err := e.Primitive(b.HDRSettingVector); err != nil {
		return fmt.Errorf("Failed to encode HDR setting vector %d: %w", b.HDRSettingVector, err)
	}
	if err := EncodeStateProp(e, b.StateProp); err != nil {
		return fmt.Errorf("Failed to encode state property: %w", err)
	}
	if err := EncodeStateGroup(e, b.StateGroup); err != nil {
		return fmt.Errorf("Failed to encode state group: %w", err)
	}
	if err := EncodeRTPC(e, b.RTPC); err != nil {
		return fmt.Errorf("Failed to encode RTPC: %w", err)
	}
	return e.Expect(curr, size)
}
