package wwise

import (
	"github.com/Dekr0/wwise-teller/wio"
)

type ChannelConfigType uint8

const (
	ChannelConfigTypeAnonymous            ChannelConfigType = 0
	ChannelConfigTypeStandard             ChannelConfigType = 1
	ChannelConfigTypeAmbisonic            ChannelConfigType = 2
	ChannelConfigTypeObjects              ChannelConfigType = 3
	ChannelConfigTypeUseDeviceMain        ChannelConfigType = 14
	ChannelConfigTypeUseDevicePassThrough ChannelConfigType = 15
)

var ChannelConfigTypeName []string = []string{
	"Anonymous",
	"Standard",
	"Ambisonic",
	"Objects",
	"Unknown Channel Configuration Type 4",
	"Unknown Channel Configuration Type 5",
	"Unknown Channel Configuration Type 6",
	"Unknown Channel Configuration Type 7",
	"Unknown Channel Configuration Type 8",
	"Unknown Channel Configuration Type 9",
	"Unknown Channel Configuration Type 10",
	"Unknown Channel Configuration Type 11",
	"Unknown Channel Configuration Type 12",
	"Unknown Channel Configuration Type 13",
	"Use Device Main",
	"Use Device Pass Through",
}

type ChannelMaskType uint32

const (
	ChannelMaskTypeFL  ChannelMaskType = 1 << 0
	ChannelMaskTypeFR  ChannelMaskType = 1 << 1
	ChannelMaskTypeFC  ChannelMaskType = 1 << 2
	ChannelMaskTypeLFE ChannelMaskType = 1 << 3
	ChannelMaskTypeBL  ChannelMaskType = 1 << 4
	ChannelMaskTypeBR  ChannelMaskType = 1 << 5
	ChannelMaskTypeFLC ChannelMaskType = 1 << 6
	ChannelMaskTypeFRC ChannelMaskType = 1 << 7
	ChannelMaskTypeBC  ChannelMaskType = 1 << 8
	ChannelMaskTypeSL  ChannelMaskType = 1 << 9
	ChannelMaskTypeSR  ChannelMaskType = 1 << 10
	ChannelMaskTypeTC  ChannelMaskType = 1 << 11
    ChannelMaskTypeTFL ChannelMaskType = 1 << 12
    ChannelMaskTypeTFC ChannelMaskType = 1 << 13
    ChannelMaskTypeTFR ChannelMaskType = 1 << 14
    ChannelMaskTypeTBL ChannelMaskType = 1 << 15
    ChannelMaskTypeTBC ChannelMaskType = 1 << 16
    ChannelMaskTypeTBR ChannelMaskType = 1 << 17
)

var ChannelMaskTypeNames map[ChannelMaskType]string = map[ChannelMaskType]string{
	ChannelMaskTypeFL: "FL",
	ChannelMaskTypeFR: "FR",
	ChannelMaskTypeFC: "FC",
	ChannelMaskTypeLFE: "LFE",
	ChannelMaskTypeBL: "BL",
	ChannelMaskTypeBR: "BR",
	ChannelMaskTypeFLC: "FLC",
	ChannelMaskTypeFRC: "FRC",
	ChannelMaskTypeBC: "BC",
	ChannelMaskTypeSL: "SL",
	ChannelMaskTypeSR: "SR",
	ChannelMaskTypeTC: "TC",
	ChannelMaskTypeTFL: "TFL",
	ChannelMaskTypeTFC: "TFC",
	ChannelMaskTypeTFR: "TFR",
	ChannelMaskTypeTBL: "TBL",
	ChannelMaskTypeTBC: "TBC",
	ChannelMaskTypeTBR: "TBR",
}

type Bus struct {
	HircObj

	Id                        uint32
	OverrideBusId             uint32
	DeviceShareSetID          uint32 // if OverrideBusId == 0
	PropBundle                PropBundle
	PositioningParam          PositioningParam
	AuxParam                  AuxParam

	VirtualBehaviorBitVector  uint8
	MaxNumInstance            uint16

	ChannelConf               uint32
	HDRBitVector              uint8
	CanSetHDR                 int8
	RecoveryTime              int32
	MaxDuckVolume             float32
	// NumDucks               uint32
	DuckInfoList              []DuckInfo
	BusFxParam                BusFxParam
	OverrideAttachmentParams  uint8 // <= 145
	BusFxMetadataParam        BusFxMetadataParam
	BusRTPC                   RTPC
	StateProp                 StateProp
	StateGroup                StateGroup
}

const SizeOfDuckInfo = 18

type DuckInfo struct {
	BusID         uint32
	DuckVolume    float32
	FadeOutTime   int32
	FadeInTime    int32
	EnumFadeCurve InterpCurveType
	TargetProp    PropType
}

func (h *Bus) KillNewest() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 0)
}

func (h *Bus) SetKillNewest(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 0, set)
}

func (h *Bus) UseVirtualBehavior() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 1)
}

func (h *Bus) SetUseVirtualBehavior(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 1, set)
}

func (h *Bus) IgnoreParentMaxNumInstance() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 2)
}

func (h *Bus) SetIgnoreParentMaxNumInstance(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 2, set)
}

func (h *Bus) BackgroundMusic() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 3)
}

func (h *Bus) SetBackgroundMusic(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 3, set)
}

func (h *Bus) NumChannel() uint8 {
	return uint8((h.ChannelConf >> 0) & 0xFF)
}

func (h *Bus) ChannelConfig() ChannelConfigType {
	return ChannelConfigType((h.ChannelConf >> 8) & 0xF)
}

func (h *Bus) ChannelMask() ChannelMaskType {
	return ChannelMaskType((h.ChannelConf >> 12 ) & 0xFFFFF)
}

func (h *Bus) IsHDRBus() bool {
	return wio.GetBit(h.HDRBitVector, 0)
}

func (h *Bus) SetHDRBus(set bool) {
	h.HDRBitVector = wio.SetBit(h.HDRBitVector, 0, set)
}

func (h *Bus) HDRReleaseModeExponential() bool {
	return wio.GetBit(h.HDRBitVector, 1)
}

func (h *Bus) SetHDRReleaseModeExponential(set bool) {
	h.HDRBitVector = wio.SetBit(h.HDRBitVector, 1, set)
}

func (h *Bus) Encode(v int) []byte {
	dataSize := h.DataSize(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeBus))
	w.Append(dataSize)
	w.Append(h.Id)
	w.Append(h.OverrideBusId)
	if h.OverrideBusId == 0 {
		w.Append(h.DeviceShareSetID)
	}
	w.AppendBytes(h.PropBundle.Encode(v))
	w.AppendBytes(h.PositioningParam.Encode(v))
	w.AppendBytes(h.AuxParam.Encode(v))
	w.Append(h.VirtualBehaviorBitVector)
	w.Append(h.MaxNumInstance)
	w.Append(h.ChannelConf)
	w.Append(h.HDRBitVector)
	w.Append(h.RecoveryTime)
	w.Append(h.MaxDuckVolume)
	w.Append(uint32(len(h.DuckInfoList)))
	for _, i := range h.DuckInfoList {
		w.Append(i)
	}
	w.AppendBytes(h.BusFxParam.Encode(v))
	if v <= 145 {
		w.Append(h.OverrideAttachmentParams)
	}
	w.AppendBytes(h.BusFxMetadataParam.Encode(v))
	w.AppendBytes(h.BusRTPC.Encode(v))
	w.AppendBytes(h.StateProp.Encode(v))
	w.AppendBytes(h.StateGroup.Encode(v))
	return w.BytesAssert(int(size))
}

func (h *Bus) DataSize(v int) uint32 {
	size := uint32(28)
	if v <= 145 {
		size += 1
	}
	size += h.PropBundle.Size(v) + 
		h.PositioningParam.Size(v) +
		h.AuxParam.Size(v) +
		uint32(len(h.DuckInfoList)) * SizeOfDuckInfo +
		h.BusFxParam.Size(v) +
		h.BusFxMetadataParam.Size(v) +
		h.BusRTPC.Size(v) +
		h.StateProp.Size(v) +
		h.StateGroup.Size(v)
	if h.OverrideBusId == 0 {
		size += 4
	}
	return size
}

func (h *Bus) BaseParameter() *BaseParameter { return nil }

func (h *Bus) HircType() HircType { return HircTypeBus }

func (h *Bus) HircID() (uint32, error) { return h.Id, nil }

func (h *Bus) IsCntr() bool { return false }

func (h *Bus) NumLeaf() int { return 0 }

func (h *Bus) ParentID() uint32 { return 0 }

func (h *Bus) AddLeaf(o HircObj) { panic("Bus object cannot add leaf") }

func (h *Bus) RemoveLeaf(o HircObj) { panic("Bus object cannot remove leaf") }

func (h *Bus) Leafs() []uint32 { return []uint32{} }

type BusFxParam struct {
	FxChunk       FxChunk
	FxID_0        uint32  // <= 145
	IsShareSet_0  uint8   // != 0 <= 145
}

func (b *BusFxParam) Encode(v int) []byte {
	size := b.Size(v)
	w := wio.NewWriter(uint64(size))
	w.AppendBytes(b.FxChunk.Encode(v))
	if v <= 145 {
		w.Append(b.FxID_0)
		w.Append(b.IsShareSet_0)
	}
	return w.BytesAssert(int(size))
}

func (b *BusFxParam) Size(v int) uint32 {
	if v <= 145 {
		return b.FxChunk.Size(v) + 5
	} else {
		return b.FxChunk.Size(v)
	}
}

type BusFxMetadataParam struct {
	// NumFx uint8
	FxChunkMetadataItems []FxChunkMetadataItem
}

func (b *BusFxMetadataParam) Encode(v int) []byte {
	size := b.Size(v)
	w := wio.NewWriter(uint64(size))
	w.Append(uint8(len(b.FxChunkMetadataItems)))
	for _, i := range b.FxChunkMetadataItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (b *BusFxMetadataParam) Size(int) uint32 {
	return 1 + uint32(len(b.FxChunkMetadataItems)) * SizeOfFxChunkMetadata
}
