package wwise

import (
	uio "github.com/Dekr0/unwise/io"
)

// --- struct definition --- //

type StateProp struct {
	NumStateProp   uio.V128
	PropId       []uio.V128
	AccumType    []u8
	InDb         []u8
}

type StatePropComponent struct {
	ActorMixerStateProp map[u32]*StateProp
}
