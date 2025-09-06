package decoder

import (
	"io"
)

// @parameter size - for assertion (notice size in here is after consuming id 
// and and action type)
func DecodeAction(size u32, r io.Reader, ver u32) {
}
