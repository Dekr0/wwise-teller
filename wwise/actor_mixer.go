package wwise

import "fmt"

type ActorMixer struct {
	Id u32
	BaseParameter *BaseParameter
	Container *Container
}

func (h *HIRC) ActorMixer(internalId u32, version u32) (a *ActorMixer) {
	a = &ActorMixer{}
	a.Id = h.GetHierarchyNode(internalId).Id
	a.BaseParameter = h.BaseParameter(internalId, version)
	a.Container = h.GetContainer(internalId)
	return a
}

func AssertActorMixer(a *ActorMixer, version u32) error {
	if err := AssertBaseParameter(a.BaseParameter, version); err != nil {
		return fmt.Errorf("Base parameter assertion failed: %w", err)
	}
	return nil
}

func SizeOfActorMixer(a *ActorMixer, version u32) (size u32) {
	size = SizeOfHierarchyId
	size += SizeOfBaseParameter(a.BaseParameter, version)
	size += SizeOfContainer(a.Container)
	return size
}

func EncodeActorMixer(e *HircEncoderCtx, a *ActorMixer) error {
	size := SizeOfActorMixer(a, e.Version)
	header := HierarchyHeader{ HircTypeActorMixer, size }
	if err := e.Struct(header, SizeOfHierarchyHeader); err != nil {
		return fmt.Errorf("Failed to encode hierarchy header: %w", err)
	}
	curr := e.Count()
	if err := e.Primitive(a.Id); err != nil {
		return fmt.Errorf("Failed to encode id: %w", err)
	}
	if err := EncodeBaseParameter(e, a.BaseParameter); err != nil {
		return fmt.Errorf("Failed to encode base parameter: %w", err)
	}
	if err := EncodeContainer(e, a.Container); err != nil {
		return fmt.Errorf("Failed to encode container: %w", err)
	}
	return e.Expect(curr, size)
}
