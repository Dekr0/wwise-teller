package wwise

import (
	"fmt"

	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type StateGroup struct {
	NumStateGroups   uio.V128
	StateGroupId    []u32
	StateSyncType   []StateSyncType
	NumStates       []uio.V128
	States          [][]StateGroupState
	StatesProp      [][]StateGroupStateProp
}

type StateGroupState struct {
	Id           u32
	InstanceId   u32 // v <= 145
}

// v > 145
type StateGroupStateProp struct {
	Ids  []u16
	Vals []f32
}

type StateGroupComponent struct {
	StateGroup map[u32]StateGroup
}

// --- allocating and freeing --- //

func AllocStateGroupComponent(size u32) *StateGroupComponent {
	if size <= 0 {
		return &StateGroupComponent{
			StateGroup: make(map[u32]StateGroup),
		}
	}
	return &StateGroupComponent{
		StateGroup: make(map[u32]StateGroup, size),
	}
}

func AllocStateGroup(size uio.V128, version u32) *StateGroup {
	if version <= 145 {
		return &StateGroup{
			NumStateGroups: size,
			StateGroupId: make([]u32, size.V, size.V),
			StateSyncType: make([]StateSyncType, size.V, size.V),
			NumStates: make([]uio.V128, size.V, size.V),
			States: make([][]StateGroupState, size.V, size.V),
		}
	} else {
		return &StateGroup{
			NumStateGroups: size,
			StateGroupId: make([]u32, size.V, size.V),
			StateSyncType: make([]StateSyncType, size.V, size.V),
			NumStates: make([]uio.V128, size.V, size.V),
			States: make([][]StateGroupState, size.V, size.V),
			StatesProp: make([][]StateGroupStateProp, size.V, size.V),
		}
	}
}

// --- assertion --- //

func AssertStateGroup(s StateGroup, version u32) error {
	if s.NumStateGroups.V != u64(len(s.StateGroupId)) {
		return fmt.Errorf("State group counter (%d) does not equal to of # of items in state group (%d)",
			s.NumStateGroups.V, u64(len(s.StateGroupId)),
		)
	}
	if len(s.StateGroupId) != len(s.StateSyncType) {
		return fmt.Errorf("# of state group id (%d) does not equal to # of state sync type (%d)",
			len(s.StateGroupId), len(s.StateSyncType),
		)
	}
	if len(s.StateSyncType) != len(s.NumStates) {
		return fmt.Errorf("# of state sync types (%d) does not equal to state group state count (%d)",
			len(s.StateSyncType), len(s.NumStates),
		)
	}
	if len(s.NumStates) != len(s.States) {
		return fmt.Errorf("# of state group state count (%d) does not equal to state group state array (%d)",
			len(s.NumStates), len(s.States),
		)
	}
	if version <= 145 {
		if s.StatesProp != nil {
			return fmt.Errorf("A state group does not have state group state property in Wwise version 145.")
		}
	} else {
		if s.StatesProp == nil {
			return fmt.Errorf("A state group should have state group state property in Wwise version 145 above.")
		}
		if len(s.States) != len(s.StatesProp) {
			return fmt.Errorf("# of state group state array (%d) does not equal to state group state prop array (%d)",
				len(s.States), len(s.StatesProp),
			)
		}
	}
	for i, numState := range s.NumStates {
		if numState.V != u64(len(s.States[i])) {
			return fmt.Errorf("state group state count (%d) does not equal to # of states (%d) in state group %d",
				numState.V, u64(len(s.States)), i,
			)
		}
		if version > 145 {
			if numState.V != u64(len(s.StatesProp[i])) {
				return fmt.Errorf("state group state count (%d) does not equal to # of state props (%d) in state group %d",
					numState.V, u64(len(s.StatesProp[i])), i,
				)
			}
			for j, stateProp := range s.StatesProp[i] {
				if len(stateProp.Ids) != len(stateProp.Vals) {
					return fmt.Errorf("# of state group state prop ids (%d) does not equal to # of state group state prop values (%d) at state %d in state group %d",
						len(stateProp.Ids), len(stateProp.Vals), j, i,
					)
				}
				if len(stateProp.Ids) > 65536 {
					return fmt.Errorf("# of state group state prop exceed 65536")
				}
			}
		}
	}
	return nil
}

// --- sizing --- //

func SizeOfStateGroup(s StateGroup, version u32) (size u32) {
	size = u32(len(s.NumStateGroups.B))
	size += (Size8 + Size32) * u32(s.NumStateGroups.V)
	for i, numState := range s.NumStates {
		size += u32(len(numState.B))
		if version <= 145 {
			size += u32(numState.V) * SizeOfStateGroupStateLE145
		} else {
			size += u32(numState.V) * Size32
			for _, stateProp := range s.StatesProp[i] {
				size += Size16 + u32(len(stateProp.Ids)) * SizeOfStateGroupStateProp
			}
		}
	}
	return size
}

// --- encoding --- //

func EncodeStateGroup(e *HircEncoderCtx, s StateGroup) error {
	curr := e.Encoder.Count
	size := SizeOfStateGroup(s, e.Version)
	if err := e.Bytes(s.NumStateGroups.B); err != nil {
		return fmt.Errorf("Failed to encode state group count %d: %w",
			s.NumStateGroups.V, err,
		)
	}
	for i, stateGroupId := range s.StateGroupId {
		if err := e.Primitive(stateGroupId); err != nil {
			return fmt.Errorf("Failed to encode %d-th state group id %d: %w",
				i, stateGroupId, err,
			)
		}
		if err := e.Primitive(s.StateSyncType[i]); err != nil {
			return fmt.Errorf("Failed to encode %d-th state sync type %d: %w",
				i, s.StateSyncType[i], err,
			)
		}
		if err := e.Bytes(s.NumStates[i].B); err != nil {
			return fmt.Errorf("Failed to encode %d-th state number count %d: %w",
				i, s.NumStates[i].V, err,
			)
		}
		for j, state := range s.States[i] {
			if err := e.Primitive(state.Id); err != nil {
				return fmt.Errorf("Failed to encode %d-th state id %d in %d-th state group: %w",
					j, state, i, err,
				)
			}
			if e.Version <= 145 {
				if err := e.Primitive(s.States[i][j].InstanceId); err != nil {
					return fmt.Errorf("Failed to encode %d-th state instance id %d in %d-th state group: %w",
						j, s.States[i][j].InstanceId, i, err,
					)
				}
			} else {
				if err := e.Primitive(u16(len(s.StatesProp[i][j].Ids))); err != nil {
					return fmt.Errorf("Failed to encode %d-th state group state property count %d in %d-th state group: %w",
						j, u16(len(s.StatesProp[i][j].Ids)), i, err,
					)
				}
				for k, id := range s.StatesProp[i][j].Ids {
					if err := e.Primitive(id); err != nil {
						return fmt.Errorf("Failed to encode %d-th state id %d in %d-th state group property in %d-th state group: %w",
							k, id, j, i, err,
						)
					}
				}
				for k, val := range s.StatesProp[i][j].Vals {
					if err := e.Primitive(val); err != nil {
						return fmt.Errorf("Failed to encode %d-th state property value %f in %d-th state group property in %d-th state group: %w",
							k, val, j, i, err,
						)
					}
				}
			}
		}
	}
	return e.Expect(curr, size)
}

// --- component wrapper for assertion, encoding, sizing --- //

func (c *StateGroupComponent) AssertStateGroupById(internalId u32, version u32) error {
	s, in := c.StateGroup[internalId]
	if !in {
		return fmt.Errorf("Failed to locate state group")
	}
	return AssertStateGroup(s, version)
}

func (c *StateGroupComponent) SizeOfStateGroupById(internalId u32, version u32) u32 {
	s, in := c.StateGroup[internalId]
	if !in {
		panic("Failed to locate state group")
	}
	return SizeOfStateGroup(s, version)
}

func (c *StateGroupComponent) EncodeStateGroupById(e *HircEncoderCtx, internalId u32) error {
	s, in := c.StateGroup[internalId]
	if !in {
		return fmt.Errorf("Failed to locate state group")
	}
	return EncodeStateGroup(e, s)
}

// --- component getter and setter --- //

func (c *StateGroupComponent) HasStateGroup(internalId u32) (in bool) {
	_, in = c.StateGroup[internalId]
	return in
}

func (c *StateGroupComponent) GetStateGroup(internalId u32) (p StateGroup) {
	p, in := c.StateGroup[internalId]
	if !in {
		panic(fmt.Errorf("Failed to locate state group."))
	}
	return p
}

func (c *StateGroupComponent) AddStateGroup(internalId u32, s StateGroup) {
	if _, in := c.StateGroup[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.StateGroup[internalId] = s
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetStateGroup(internalId u32) StateGroup {
	return h.StateGroupComponent.GetStateGroup(internalId)
}

// --- constant & definition --- //

const SizeOfStateGroupStateLE145 = Size32 * 2
const SizeOfStateGroupStateProp  = Size16 + Size32

type StateSyncType = u8
type StatePropType = u16
