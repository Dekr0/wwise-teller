package parser

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

func ParseState(size uint32, r *wio.Reader, v int) *wwise.State {
	assert.Equal(0, r.Pos(), "State parser position doesn't start at position 0.")
	begin := r.Pos()

	state := wwise.State{
		StateID: r.U32Unsafe(),
		StateProps: make([]struct{PID uint16; Val float32}, r.U16Unsafe()),
	}
	for i := range state.StateProps {
		state.StateProps[i].PID, state.StateProps[i].Val = r.U16Unsafe(), r.F32Unsafe()
	}

	end := r.Pos()
	if begin >= end {
		panic("Reader consumes zero byte.")
	}
	assert.Equal(uint64(size), end-begin,
		"The amount of bytes reader consume does not equal size in the hierarchy header",
	)
	return &state
}

func ParseEvent(size uint32, r *wio.Reader, v int) *wwise.Event {
	assert.Equal(0, r.Pos(), "Layer container parser position doesn't start at position 0.")
	begin := r.Pos()
	e := wwise.Event{}
	e.Id = r.U32Unsafe()
	e.ActionIDs = make([]uint32, r.U8Unsafe())
	for i := range e.ActionIDs {
		e.ActionIDs[i] = r.U32Unsafe()
	}
	end := r.Pos()
	if begin >= end {
		panic("Reader read zero bytes")
	}
	assert.Equal(size, uint32(end-begin),
		"The amount of bytes reader consume doesn't equal to the size in hierarchy header",
	)
	return &e
}

