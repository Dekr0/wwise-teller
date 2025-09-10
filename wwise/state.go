package wwise

import (
	bin "encoding/binary"
	"fmt"
	"io"
)

const StateBaseSize = HircIdSize + 2

type StateProps struct {
	Ids  []u16
	Vals []f32
}

type StateProp struct {
	Id  u16
	Val f32
}

type StateComponent struct {
	StateProps map[u32]*StateProps
}

// Call this before obtaining size of State and encoding State
func AssertState(s *StateComponent, internalId u32, id u32) {
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

func SizeOfState(s *StateComponent, internalId u32, id u32) (size u32) {
	size = StateBaseSize
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
	AssertState(s, internalId, id)

	var err error

	pos := u32(0)

	size := SizeOfState(s, internalId, id)
	stateProp, in := s.StateProps[internalId]
	if !in {
		panic(fmt.Sprintf("State %d does not have state property", id))
	}
	
	header := HierarchyHeader{ HircTypeState, size }
	if err = bin.Write(w, o, header); err != nil {
		panic(fmt.Errorf("Failed to encode State %d hierarchy header: %w", id, err))
	}
	pos += HircHeaderSize

	if err = bin.Write(w, o, id); err != nil {
		panic(fmt.Errorf("Failed to encode State %d id: %w", id, err))
	}
	pos += Size32

	if err = bin.Write(w, o, u16(len(stateProp.Ids))); err != nil {
		panic(fmt.Errorf("Failed to encode State %d # of state properties: %w", id, err))
	}
	pos += Size16

	ids := stateProp.Ids
	vals := stateProp.Vals
	for i, id := range ids {
		stateProp := StateProp{ id, vals[i] }
		if err = bin.Write(w, o, stateProp); err != nil {
			panic(fmt.Errorf("Failed to encode State %d state property at index %d: %w", id, i, err))
		}
	}
}
