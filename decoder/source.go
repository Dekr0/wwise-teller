package decoder

import (
	"io"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecodeSourceData(
	r           io.Reader,
	o           order,
	version     u32, 
) (sourceData *wwise.SourceData) {
	sourceData = wwise.AllocSourceData()

	sourceData.PluginID = uio.U32P(r, o)
	sourceData.StreamType = uio.U8P(r, o) 
	sourceData.SourceID = uio.U32P(r, o)
	if version <= 150 {
		sourceData.InMemoryMediaSize = uio.U32P(r, o)
	} else {
		sourceData.CacheID = uio.U32P(r, o)
		sourceData.InMemoryMediaSize = uio.U32P(r, o)
	}
	sourceData.SourceBits = uio.U8P(r, o)
	
	return sourceData
}
