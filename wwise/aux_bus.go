package wwise

import "github.com/Dekr0/wwise-teller/wio"

type AuxBus struct {
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

func (h *AuxBus) KillNewest() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 0)
}

func (h *AuxBus) SetKillNewest(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 0, set)
}

func (h *AuxBus) UseVirtualBehavior() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 1)
}

func (h *AuxBus) SetUseVirtualBehavior(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 1, set)
}

func (h *AuxBus) IgnoreParentMaxNumInstance() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 2)
}

func (h *AuxBus) SetIgnoreParentMaxNumInstance(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 2, set)
}

func (h *AuxBus) BackgroundMusic() bool {
	return wio.GetBit(h.VirtualBehaviorBitVector, 3)
}

func (h *AuxBus) SetBackgroundMusic(set bool) {
	h.VirtualBehaviorBitVector = wio.SetBit(h.VirtualBehaviorBitVector, 3, set)
}

func (h *AuxBus) NumChannel() uint8 {
	return uint8((h.ChannelConf >> 0) & 0xFF)
}

func (h *AuxBus) ChannelConfig() ChannelConfigType {
	return ChannelConfigType((h.ChannelConf >> 8) & 0xF)
}

func (h *AuxBus) ChannelMask() ChannelMaskType {
	return ChannelMaskType((h.ChannelConf >> 12 ) & 0xFFFFF)
}

func (h *AuxBus) IsHDRBus() bool {
	return wio.GetBit(h.HDRBitVector, 0)
}

func (h *AuxBus) SetHDRBus(set bool) {
	h.HDRBitVector = wio.SetBit(h.HDRBitVector, 0, set)
}

func (h *AuxBus) HDRReleaseModeExponential() bool {
	return wio.GetBit(h.HDRBitVector, 1)
}

func (h *AuxBus) SetHDRReleaseModeExponential(set bool) {
	h.HDRBitVector = wio.SetBit(h.HDRBitVector, 1, set)
}

func (h *AuxBus) Encode(v int) []byte {
	dataSize := h.DataSize(v)
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeAuxBus))
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

func (h *AuxBus) DataSize(v int) uint32 {
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

func (h *AuxBus) BaseParameter() *BaseParameter { return nil }

func (h *AuxBus) HircType() HircType { return HircTypeAuxBus }

func (h *AuxBus) HircID() (uint32, error) { return h.Id, nil }

func (h *AuxBus) IsCntr() bool { return false }

func (h *AuxBus) NumLeaf() int { return 0 }

func (h *AuxBus) ParentID() uint32 { return 0 }

func (h *AuxBus) AddLeaf(o HircObj) { panic("Auxiliary Bus object cannot add leaf") }

func (h *AuxBus) RemoveLeaf(o HircObj) { panic("Auxiliary Bus object cannot remove leaf") }

func (h *AuxBus) Leafs() []uint32 { return []uint32{} }
