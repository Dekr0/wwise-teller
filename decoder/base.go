package decoder

import (
	"io"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecodeBaseParameter(
	r       io.Reader,
	o       order, 
	version u32, 
	b      *wwise.BaseParameter,
) {
	b.OverrideParentFx = uio.U8P(r, o)
	b.FXs = *AllocDecodeFXs(r, o, version)
	b.OverrideFxMetadata = uio.U8P(r, o)
	b.FxMetadatas = *AllocDecodeFXMetadatas(r, o)
	if version <= 145 {
		b.OverrideAttachmentParam = uio.U8P(r, o)
	}
	b.OverideBusId = uio.U32P(r, o)
	b.DirectParentId = uio.U32P(r, o)
	b.SettingVector = uio.U8P(r, o)
	b.Prop = *AllocDecodeProp(r, o)
	b.RProp = *AllocDecodeRProp(r, o)
	b.PositionParam = *AllocDecodePositionParam(r, o)
	b.AuxParam = *AllocDecodeAuxParam(r, o)
	b.AdvanceSettingVector = uio.U8P(r, o)
	b.VirtualQueueBehavior = uio.U8P(r, o)
	b.MaxNumInstance = uio.U16P(r, o)
	b.BelowThresholdBehavior = uio.U8P(r, o)
	b.HDRSettingVector = uio.U8P(r, o)
	b.StateProp = *AllocDecodeStateProp(r, o)
	b.StateGroup = *AllocDecodeStateGroup(r, o, version)
	b.RTPC = *AllocDecodeRTPC(r, o)
}
