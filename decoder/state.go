package decoder

import (
	"fmt"
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func DecodeState(r io.Reader, o order, ver u32, h *wwise.HIRC, size u32) {
	p := 0

	id := uio.U32PT(r, o, &p)

	numStateProps := uio.U16PT(r, o, &p)
	data := &wwise.StateProp{
		Ids: make([]u16, numStateProps, numStateProps),
		Vals: make([]f32, numStateProps, numStateProps),
	}
	for i := range numStateProps {
		data.Ids[i] = uio.U16PT(r, o, &p)
		data.Vals[i] = uio.F32PT(r, o, &p)
	}

	if p != int(size) {
		const sfmt = "After parsing State %d, expecting debug position to" + 
			         " be %d but receive %d"
		panic(fmt.Sprintf(sfmt, id, size, p))
	}

	wwise.HIRCNewState(h, id, data)
}
