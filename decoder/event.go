package decoder

import (
	"fmt"
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func DecodeEvent(r io.Reader, o order, ver u32, h *wwise.HIRC, size u32) {
	p := 0

	id := uio.U32PT(r, o, &p)

	data := &wwise.EventData{}
	data.NumActionIds = *uio.VV128PT(r, o, &p)
	v := data.NumActionIds.V
	data.ActionIds = make([]u32, v, v)
	for i := range v {
		data.ActionIds[i] = uio.U32PT(r, o, &p)
	}

	if p != int(size) {
		const sfmt = "After parsing Event %d, expecting debug position to" + 
			         " be %d but receive %d"
		panic(fmt.Sprintf(sfmt, id, size, p))
	}

	wwise.HIRCNewEvent(h, id, data)
}
