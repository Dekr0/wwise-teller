package decoder

import (
	"io"

	"github.com/Dekr0/unwise/wwise"
	uio "github.com/Dekr0/unwise/io"
)

func AllocDecodeRTPC(r io.Reader, o order) (rt *wwise.RTPC) {
	rt = wwise.AllocateRTPC(u32(uio.U16P(r, o)))
	for i := range rt.Id {
		rt.Id[i] = uio.U32P(r, o)
		rt.Type[i] = uio.U8P(r, o)
		rt.Accum[i] = uio.U8P(r, o)
		rt.ParamId[i] = *uio.VV128P(r, o)
		rt.CurveId[i] = uio.U32P(r, o)
		rt.Scaling[i] = uio.U8P(r, o)
		rt.Graph[i] = *wwise.AllocateRTPCGraph(u32(uio.U16P(r, o)))
		for j := range rt.Graph[i].PointX {
			frame := wwise.RTPCGraphFrame{}
			uio.StructP(r, o, &frame)
			rt.Graph[i].PointX[j] = frame.PointX
			rt.Graph[i].PointY[j] = frame.PointY
			rt.Graph[i].PointLerp[j] = frame.PointLerp
		}
	}
	return rt
}
