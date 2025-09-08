package decoder

type DecoderOption struct {
	option              u8
	DecodedChunkRoutine u8
	DecoderBufferSize   u32
}

const PageSize1k   = 1024
const PageSize2k   = PageSize1k * 2
const PageSize4k   = PageSize2k * 2
const PageSize8k   = PageSize4k * 2
const PageSize16k  = PageSize8k * 2
const PageSize32k  = PageSize16k * 2
const PageSize64k  = PageSize32k * 2
const PageSize128k = PageSize64k * 2

const DecodeBufferSize = PageSize32k
const MaskMETA u8 = 0b1000_0000
const MaskDATA u8 = 0b0000_0100

func (o *DecoderOption) IncludeDATA() {
	o.option |= MaskDATA
}

func (o *DecoderOption) ExcludeDATA() {
	o.option = o.option | (^MaskDATA)
}

func (o *DecoderOption) IsIncludeDATA() bool {
	return o.option & MaskDATA > 0
}

func (o *DecoderOption) IncludeMETA() {
	o.option |= MaskMETA
}

func (o *DecoderOption) ExcludeMETA() {
	o.option = o.option | (^MaskMETA)
}

func (o *DecoderOption) IsIncludeMETA() bool {
	return o.option & MaskMETA > 0
}
