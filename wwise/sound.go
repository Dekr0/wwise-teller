package wwise

import "fmt"

type Sound struct {
	Id              u32
	SourceData      SourceData
	PluginParam     PluginParam
	BaseParameter   BaseParameter
}

func (h *HIRC) Sound(internalId u32, version u32, inOut *Sound) {
	inOut.Id = h.GetHierarchyNode(internalId).Id
	inOut.SourceData = h.GetSourceData(internalId)
	if SourceHasPluginParam(inOut.SourceData) {
		inOut.PluginParam = h.GetPluginParam(internalId)
	}
	inOut.BaseParameter = h.BaseParameter(internalId, version)
}

func AssertSound(sound Sound, version u32) error {
	hasPlugin := SourceHasPluginParam(sound.SourceData)
	if hasPlugin {
		if err := AssertPluginParm(sound.PluginParam); err != nil {
			return fmt.Errorf("Source plugin parameter assertion failed: %w", err)
		}
	}
	if err := AssertBaseParameter(sound.BaseParameter, version); err != nil {
		return fmt.Errorf("Base parameter assertion failed: %w", err)
	}
	return nil
}

func SizeOfSound(sound Sound, version u32) (size u32) {
	size = SizeOfHierarchyId
	size += SizeOfSourceData(version)
	size += SizeOfBaseParameter(sound.BaseParameter, version)
	if SourceHasPluginParam(sound.SourceData) {
		size += SizeOfPluginParam(sound.PluginParam)
	}
	return size
}

func EncodeSound(e *HircEncoderCtx, sound Sound) error {
	size := SizeOfSound(sound, e.Version)
	header := HierarchyHeader{ HircTypeSound, size }
	if err := e.Struct(header, SizeOfHierarchyHeader); err != nil {
		return fmt.Errorf("Failed to encode hierarchy header: %w", err)
	}
	curr := e.Count()
	if err := e.Primitive(sound.Id); err != nil {
		return fmt.Errorf("Failed to encode id: %w", err)
	}
	if err := EncodeSourceData(e, sound.SourceData); err != nil {
		return fmt.Errorf("Failed to encode source data: %w", err)
	}
	if SourceHasPluginParam(sound.SourceData) {
		if err := EncodePluginParam(e, sound.PluginParam); err != nil {
			return fmt.Errorf("Failed to encode plugin parameter: %w", err)
		}
	}
	if err := EncodeBaseParameter(e, sound.BaseParameter); err != nil {
		return fmt.Errorf("Failed to encode base parameter: %w", err)
	}
	return e.Expect(curr, size)
}
