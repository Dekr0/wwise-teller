package wwise

import (
	bin "encoding/binary"
	"fmt"
	"io"
)

const SizeOfStateBaseData = SizeOfHierarchyId + Size16

type StateProps struct {
	Ids  []u16
	Vals []f32
}

const SizeOfStateProp = Size16 + Size32
type StateProp struct {
	Id  u16
	Val f32
}

type StateComponent struct {
	StateProps map[u32]*StateProps
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
	w           io.Writer,
	o           order,
	s          *StateComponent,
	version     u32,
	internalId  u32,
	id          u32,
) {
	AssertState(s, version, internalId, id)

	var err error

	pos := u32(0)

	size := SizeOfState(s, version, internalId, id)
	stateProp, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Sprintf("State %d does not have state property", id))
	}
	
	header := HierarchyHeader{ HircTypeState, size }
	if err = bin.Write(w, o, header); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode hierarchy header: %w", id, err))
	}

	if err = bin.Write(w, o, id); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode id: %w", id, err))
	}
	pos += SizeOfHierarchyId

	if err = bin.Write(w, o, u16(len(stateProp.Ids))); err != nil {
		panic(fmt.Errorf("(State %d) Failed to encode # of state properties: %w", id, err))
	}
	pos += Size16

	ids := stateProp.Ids
	vals := stateProp.Vals
	for i, id := range ids {
		stateProp := StateProp{ id, vals[i] }
		if err = bin.Write(w, o, stateProp); err != nil {
			panic(fmt.Errorf("(State %d) Failed to encode %d-th state property: %w", id, i, err))
		}
		pos += SizeOfStateProp
	}

	if pos != size {
		panic(fmt.Errorf("(State %d) Expect encoded data has a size of %d but receive %d", id, size, pos))
	}
}
