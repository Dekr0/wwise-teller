package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeStateGroup(r io.Reader, o order, version u32) (s *wwise.StateGroup) {
	s = wwise.AllocStateGroup(uio.VV128P(r, o), version)
	for i := range s.StateGroupId {
		s.StateGroupId[i] = uio.U32P(r, o)
		s.StateSyncType[i] = uio.U8P(r, o)
		s.NumStates[i] = *uio.VV128P(r, o)
		s.States[i] = make([]wwise.State, s.NumStates[i].V, s.NumStates[i].V)
		if version > 145 {
			s.StatesProp[i] = make([]wwise.StateGroupStateProp, s.NumStates[i].V, s.NumStates[i].V)
		}
		for j := range s.States[i] {
			s.States[i][j].Id = uio.U32P(r, o)
			if version <= 145 {
				s.States[i][j].InstanceId = uio.U32P(r, o)
			} else {
				numStateGroupStateProp := uio.U16P(r, o)
				s.StatesProp[i][j].Ids = make([]u16, numStateGroupStateProp, numStateGroupStateProp)
				s.StatesProp[i][j].Vals = make([]f32, numStateGroupStateProp, numStateGroupStateProp)
				for k := range numStateGroupStateProp {
					s.StatesProp[i][j].Ids[k] = uio.U16P(r, o)
				}
				for k := range numStateGroupStateProp {
					s.StatesProp[i][j].Vals[k] = uio.F32P(r, o)
				}
			}
		}
	}
	return s
}
