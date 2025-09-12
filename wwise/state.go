package wwise

import (
	"fmt"
	"sync"
)

// --- struct definition --- //

const SizeOfStateBaseData = SizeOfHierarchyId + Size16

type StateProps struct {
	Ids  []u16
	Vals []f32
}

const SizeOfStateProp = Size16 + Size32
type StatePropS struct {
	Id  u16
	Val f32
}

type StateComponent struct {
	dMu sync.Mutex

	StateProps map[u32]*StateProps
}

// --- allocation / freeing --- //

func AllocStateProps(numStateProps u16) *StateProps {
	return &StateProps{
		Ids: make([]u16, numStateProps, numStateProps),
		Vals: make([]f32, numStateProps, numStateProps),
	} 
}

// --- function --- //

func AddStateData(s *StateComponent, internalId u32, data *StateProps) {
	if data == nil {
		panic("State property is nil")
	}

	s.dMu.Lock()
	defer s.dMu.Unlock()
	if _, in := s.StateProps[internalId]; in {
		panic(MonotonicIdCollision)
	}
	s.StateProps[internalId] = data
}

// Call this before obtaining size of State and encoding State
func AssertState(s *StateComponent, version u32, internalId u32, id u32) {
	stateProp, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Sprintf("State %d does not have state property", id))
	}
	if len(stateProp.Ids) != len(stateProp.Vals) {
		panic(fmt.Sprintf(
			"State %d: # of state property ids does not equal # of state values",
			id,
		))
	}
}

func SizeOfState(s *StateComponent, version u32, internalId u32, id u32) (size u32) {
	size = SizeOfStateBaseData
	stateProp, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Sprintf("State %d does not have state property", id))
	}

	size += u32(len(stateProp.Ids)) * Size16 + u32(len(stateProp.Vals)) * Size32

	return size
}

func EncodeState(
	e          *HircEncoderCtx,
	s          *StateComponent,
	internalId  u32,
	id          u32,
) {
	AssertState(s, e.Version, internalId, id)

	var err error

	size := SizeOfState(s, e.Version, internalId, id)
	stateProp, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Sprintf("State %d does not have state property", id))
	}
	
	header := HierarchyHeader{ HircTypeState, size }
	if err = HIRCEncodeStruct(e, header, SizeOfHierarchyHeader); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode hierarchy header: %w", id, err))
	}

	curr := e.Encoder.Count

	if err = HIRCEncode(e, id); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode id: %w", id, err))
	}

	if err = HIRCEncode(e, u16(len(stateProp.Ids))); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode # of state properties: %w", id, err))
	}

	ids := stateProp.Ids
	vals := stateProp.Vals
	for i, id := range ids {
		stateProp := StatePropS{ id, vals[i] }
		if err = HIRCEncodeStruct(e, stateProp, SizeOfStateProp); err != nil {
			panic(fmt.Errorf("(State %d) Failed to encode %d-th state property: %w", id, i, err))
		}
	}

	if err := HIRCAssertEncodeLimit(e, curr, size); err != nil {
		panic(fmt.Errorf("(State %d) %w", id, err))
	}
}
