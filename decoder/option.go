package decoder

type DecoderOption struct {
	DecoderBufferSize u32
	option            u8
}

const DecodeBufferSize  = 4096
const MaskMETA u8 = 0b1000_0000
const MaskDATA u8 = 0b0000_0100

func (o *DecoderOption) IncludeDATA() {
	o.option |= MaskDATA
}

func (o *DecoderOption) ExcludeDATA() {
	o.option = o.option | (^MaskDATA)
}

func (o *DecoderOption) IncludeMETA() {
	o.option |= MaskMETA
}

func (o *DecoderOption) ExcludeMETA() {
	o.option = o.option | (^MaskMETA)
}
