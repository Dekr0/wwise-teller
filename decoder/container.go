package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeContainer(r io.Reader, o order) (c *wwise.Container) {
	c = wwise.AllocContainer(uio.U32P(r, o))
	for i := range c.Ids {
		c.Ids[i] = uio.U32P(r, o)
	}
	return c
}
