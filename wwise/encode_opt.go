package wwise

const MaskMETA u8 = 0b1000_0000
const MaskDATA u8 = 0b0000_0100
const MaskHIRC u8 = 0b0000_1000

type EncodeBankOpt struct {
	option  u8
}

func IncludeEncodedMETA(o *EncodeBankOpt) {
	o.option |= MaskMETA
}

func ExcludeEncodedMETA(o *EncodeBankOpt) {
	o.option |= (^MaskMETA)
}

func IsIncludeEncodedMETA(o *EncodeBankOpt) bool {
	return o.option & MaskMETA > 0
}

type EncodeHircOpt struct {}
