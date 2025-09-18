package wwise

import (
	"fmt"
)

// --- struct definition --- //

const SizeOfStateBaseData = SizeOfHierarchyId + Size16

type State struct {
	Id          u32
	StateProps *StateHierarchyProp
}

type StateHierarchyProp struct {
	Ids  []u16
	Vals []f32
}

const SizeOfStateStateProp = Size16 + Size32
type StatePropS struct {
	Id  u16
	Val f32
}

type StateComponent struct {
	StateProps map[u32]*StateHierarchyProp
}

// --- allocation / freeing --- //

func AllocStateComponent(numState u32) *StateComponent {
	if numState <= 0 {
		return &StateComponent{
			StateProps: make(map[u32]*StateHierarchyProp),
		}
	}
	return &StateComponent{
		StateProps: make(map[u32]*StateHierarchyProp, numState),
	}
}

func AllocStateProps(numStateProps u16) *StateHierarchyProp {
	return &StateHierarchyProp{
		Ids: make([]u16, numStateProps, numStateProps),
		Vals: make([]f32, numStateProps, numStateProps),
	} 
}

// --- assertion --- //

// Has no side effect
func AssertState(s *State) error {
	stateProp := s.StateProps
	if len(stateProp.Ids) != len(stateProp.Vals) {
		return fmt.Errorf("# of state property ids does not equal # of state values")
	}
	return nil
}

// --- sizing --- //

// Has no side effect
func SizeOfState(s *State) (size u32) {
	stateProp := s.StateProps
	size = SizeOfStateBaseData
	size += u32(len(stateProp.Ids)) * Size16 + u32(len(stateProp.Vals)) * Size32
	return size
}

// --- encoding --- //

// Has no side effect
func EncodeState(
	e    *HircEncoderCtx,
	s    *State,
) error {
	var err error
	id := s.Id
	stateProps := s.StateProps
	size := SizeOfState(s)
	header := HierarchyHeader{ HircTypeState, size }
	if err = e.Struct(header, SizeOfHierarchyHeader); err != nil {
		return fmt.Errorf("Failed to encode hierarchy header: %w", err)
	}
	curr := e.Encoder.Count
	if err = e.Primitive(id); err != nil {
		return fmt.Errorf("Failed to encode id: %w", err)
	}
	if err = e.Primitive(u16(len(stateProps.Ids))); err != nil {
		return fmt.Errorf("Failed to encode # of state properties: %w", err)
	}
	ids := stateProps.Ids
	vals := stateProps.Vals
	for i, id := range ids {
		stateProp := StatePropS{ id, vals[i] }
		if err = e.Struct(stateProp, SizeOfStateStateProp); err != nil {
			return fmt.Errorf("Failed to encode %d-th state property: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

// --- component getter and setter --- //

// Has no side effect
func (s *StateComponent) HasStateProps(internalId u32) (in bool) {
	_, in = s.StateProps[internalId]
	return in
}

// Has no side effect
func (s *StateComponent) GetStateProps(internalId u32) (p *StateHierarchyProp) {
	p, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate state property"))
	}
	return p
}

// Has side effect
func (s *StateComponent) AddStateData(internalId u32, data *StateHierarchyProp) {
	if data == nil {
		panic("State property is nil")
	}
	if _, in := s.StateProps[internalId]; in {
		panic(MonotonicIdCollision)
	}
	s.StateProps[internalId] = data
}

// --- HIRC component wrapper --- //

// Has no side effect
func (h *HIRC) State(internalId u32) *State {
	node := h.Hierarchy.GetHierarchyNode(internalId)
	p := h.StateComponent.GetStateProps(internalId)
	return &State{ node.Id, p }
}
