package wwise

import "fmt"

type SoundH struct {
	Id              u32
	SourceData     *SourceData
	PluginParam    *PluginParam
	BaseParameter *BaseParameter
}

func (h *HIRC) SoundH(internalId u32, version u32) (s *SoundH) {
	s = &SoundH{}

	node := h.GetHierarchyNode(internalId)
	s.Id = node.Id

	source := h.GetSourceData(internalId)
	s.SourceData = source
	
	pluginParam := h.GetPluginParam(internalId)
	s.PluginParam = pluginParam

	b := h.BaseParameter(internalId, version)
	s.BaseParameter = b

	return s
}

func (h *HIRC) EncodeSound(e *HircEncoderCtx, internalId u32) {
	s := h.SoundH(internalId, e.Version)
	EncodeSound(e, s)
}

func AssertSound(sound *SoundH, version u32) error {
	hasPlugin := SourceHasPluginParam(sound.SourceData)

	{ // Assertion by correlating different component together
		if !hasPlugin {
			if sound.PluginParam != nil {
				return fmt.Errorf("Source data indicate this source doesn't have plugin parameter.")
			}
		}
	}

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

func SizeOfSound(sound *SoundH, version u32) (size u32) {
	size = SizeOfHierarchyId
	size += SizeOfSourceData(version)
	size += SizeOfBaseParameter(sound.BaseParameter, version)
	if SourceHasPluginParam(sound.SourceData) {
		size += SizeOfPluginParam(sound.PluginParam)
	}
	return size
}

func EncodeSound(e *HircEncoderCtx, sound *SoundH) {
	size := SizeOfSound(sound, e.Version)
	header := HierarchyHeader{ HircTypeSound, size }
	if err := e.Struct(header, SizeOfHierarchyHeader); err != nil {
		panic(fmt.Errorf("(Sound %d) Failed to encode hierarchy header: %w", 
			sound.Id, err,
		))
	}
	curr := e.Count()
	if err := e.Primitive(sound.Id); err != nil {
		panic(fmt.Errorf("(Sound %d) Failed to encode id: %w", sound.Id, err))
	}
	if err := EncodeSourceData(e, sound.SourceData); err != nil {
		panic(fmt.Errorf("(Sound %d) Failed to encode source data: %w", sound.Id, err))
	}
	if SourceHasPluginParam(sound.SourceData) {
		if err := EncodePluginParam(e, sound.PluginParam); err != nil {
			panic(fmt.Errorf("(Sound %d) Failed to encode plugin parameter: %w", sound.Id, err))
		}
	}
	if err := EncodeBaseParameter(e, sound.BaseParameter); err != nil {
		panic(fmt.Errorf("(Sound %d) Failed to encode base parameter: %w", sound.Id, err))
	}
	if err := e.Expect(curr, size); err != nil {
		panic(fmt.Errorf("(Sound %d) %w", sound.Id, err))
	}
}
