package wwise

import "github.com/Dekr0/wwise-teller/wio"

const (
	ChannelConfigTypeAnonymous            = 0
	ChannelConfigTypeStandard             = 1
	ChannelConfigTypeAmbisonic            = 2
	ChannelConfigTypeObjects              = 3
	ChannelConfigTypeUseDeviceMain        = 14
	ChannelConfigTypeUseDevicePassThrough = 15
)
var ChannelConfigTypeName []string = []string{
	"Anonymous",
	"Standard",
	"Ambisonic",
	"Objects",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"",
	"Use Device Main",
	"Use Device Pass Through",
}

type Bus struct {
	HircObj

	Id                        uint32
	OverrideBusId             uint32
	DeviceShareSetID          uint32 // if OverrideBusId == 0
	PropBundle                PropBundle
	PositioningParam          PositioningParam
	AuxParam                  AuxParam

	// 0 Kill Newest
	// 1 Use Virtual Behavior
	// 2 Ignore Parent Max Number Instance
	// 3 Background Music
	VirtualBehaviorBitVector  uint8
	MaxNumInstance            uint16

	// (>> 0) & 0xFF Number Channels
	// (>> 8) & 0xF  Enum Config Type
	// (>> 12) & 0xFFFFF Channel Mask
	// 		(1 << 0):  "FL", # front left
    // 		(1 << 1):  "FR", # front right
    // 		(1 << 2):  "FC", # front center
    // 		(1 << 3):  "LFE", # low frequency effects
    // 		(1 << 4):  "BL", # back left
    // 		(1 << 5):  "BR", # back right
    // 		(1 << 6):  "FLC", # front left center
    // 		(1 << 7):  "FRC", # front right center
    // 		(1 << 8):  "BC", # back center
    // 		(1 << 9):  "SL", # side left
    // 		(1 << 10): "SR", # side right

    // 		(1 << 11): "TC", # top center
    // 		(1 << 12): "TFL", # top front left
    // 		(1 << 13): "TFC", # top front center
    // 		(1 << 14): "TFR", # top front right
    // 		(1 << 15): "TBL", # top back left
    // 		(1 << 16): "TBC", # top back center
    // 		(1 << 17): "TBR", # top back left
	ChannelConfig             uint32
	// 0 Is HDR Bus
	// 1 HDR Release Mode Exponential
	HDRBitVector              uint8
	RecoveryTime              int32
	MaxDuckVolume             float32
	// NumDucks               uint32
	DuckInfoList              []DuckInfo
	BusFxParam                BusFxParam
	OverrideAttachmentParams  uint8
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

func (h *Bus) Encode() []byte {
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeBus))
	w.Append(dataSize)
	w.Append(h.Id)
	return w.BytesAssert(int(size))
}

func (h *Bus) DataSize() uint32 {
	return 0
}

func (h *Bus) BaseParameter() *BaseParameter { return nil }

func (h *Bus) HircType() HircType { return HircTypeEvent }

func (h *Bus) HircID() (uint32, error) { return h.Id, nil }

func (h *Bus) IsCntr() bool { return false }

func (h *Bus) NumLeaf() int { return 0 }

func (h *Bus) ParentID() int { return 0 }

func (h *Bus) AddLeaf(o HircObj) { panic("") }

func (h *Bus) RemoveLeaf(o HircObj) { panic("") }

type BusFxParam struct {
	FxChunk       FxChunk
	FxID_0        uint32
	IsShareSet_0  uint8  // !=0
}

func (b *BusFxParam) Encode() []byte {
	size := b.Size()
	w := wio.NewWriter(uint64(size))
	w.AppendBytes(b.FxChunk.Encode())
	w.Append(b.FxID_0)
	w.Append(b.IsShareSet_0)
	return w.BytesAssert(int(size))
}

func (b *BusFxParam) Size() uint32 {
	return b.FxChunk.Size() + 5
}

type BusMetaFxParam struct {
	// NumFx uint8
}
