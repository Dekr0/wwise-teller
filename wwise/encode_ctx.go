package wwise

import (
	uio "github.com/Dekr0/unwise/io"
)

type HircEncoderCtx struct {
	Encoder *uio.EncoderCtx
	Version  u32
}

func (e *HircEncoderCtx) Count() u32 {
	return e.Encoder.Count
}

func (e *HircEncoderCtx) Primitive(data any) (err error) {
	return uio.Encode(e.Encoder, data)
}

func (e *HircEncoderCtx) Bytes(data []byte) (err error) {
	return uio.EncodeBytes(e.Encoder, data)
}

func (e *HircEncoderCtx) Struct(data any, size u32) (err error) {
	return uio.EncodeStruct(e.Encoder, data, size)
}

func (e *HircEncoderCtx) Expect(prev u32, expect u32) (err error) {
	return uio.AssertEncodeLimit(e.Encoder, prev, expect)
}
