package wwise

import "fmt"

type Sound struct {
	Id              u32
	SourceData     *SourceData
	PluginParam    *PluginParam
	BaseParameter *BaseParameter
}

func (h *HIRC) Sound(internalId u32, version u32) (s *Sound) {
	s = &Sound{}
	s.Id = h.GetHierarchyNode(internalId).Id
	s.SourceData = h.GetSourceData(internalId)
	if SourceHasPluginParam(s.SourceData) {
		s.PluginParam = h.GetPluginParam(internalId)
	}
	s.BaseParameter = h.BaseParameter(internalId, version)
	return s
}

func (h *HIRC) EncodeSound(e *HircEncoderCtx, internalId u32) {
	s := h.Sound(internalId, e.Version)
	EncodeSound(e, s)
}

func AssertSound(sound *Sound, version u32) error {
	hasPlugin := SourceHasPluginParam(sound.SourceData)
	{ // Assertion by correlating different component together
		if !hasPlugin {
			if sound.PluginParam != nil {
				return fmt.Errorf("Source data indicate this source doesn't have plugin parameter but receive non nil plugin parameter")
			}
		}
	}
	if hasPlugin {
		if sound.PluginParam == nil {
			return fmt.Errorf("Source data indicate this source has plugin parameter but plugin parameter is nil")
		}
		if err := AssertPluginParm(sound.PluginParam); err != nil {
			return fmt.Errorf("Source plugin parameter assertion failed: %w", err)
		}
	}
	if err := AssertBaseParameter(sound.BaseParameter, version); err != nil {
		return fmt.Errorf("Base parameter assertion failed: %w", err)
	}
	return nil
}

func SizeOfSound(sound *Sound, version u32) (size u32) {
	size = SizeOfHierarchyId
	size += SizeOfSourceData(version)
	size += SizeOfBaseParameter(sound.BaseParameter, version)
	if SourceHasPluginParam(sound.SourceData) {
		size += SizeOfPluginParam(sound.PluginParam)
	}
	return size
}

func EncodeSound(e *HircEncoderCtx, sound *Sound) error {
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
