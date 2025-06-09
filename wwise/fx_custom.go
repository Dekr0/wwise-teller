package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

// import "github.com/Dekr0/wwise-teller/wio"

type PluginType uint16

const (
	PluginTypeNone            PluginType = 0
	PluginTypeCodec           PluginType = 1
	PluginTypeSource          PluginType = 2
	PluginTypeEffect          PluginType = 3
	PluginTypeMotionDevice    PluginType = 4
	PluginTypeMotionSource    PluginType = 5
	PluginTypeMixer           PluginType = 6
	PluginTypeSink            PluginType = 7
	PluginTypeGlobalExtension PluginType = 8
	PluginTypeMetaData        PluginType = 9
	PluginTypeInteralType     PluginType = 10 // Above this point is internal class plugin
	PluginTypeInvalid         PluginType = 0xFFFF
)

var PluginTypeNames []string = []string{
	"None",
	"Codec",
	"Source",
	"Effect",
	"Motion Device",
	"Motion Source",
	"Mixer",
	"Sink",
	"Global Extension",
	"Metadata",
	"Internal Hierarchy Class Types",
}

type PluginCompanyType uint16

const (
	PluginCompanyTypeAudiokinetic PluginCompanyType = 0
	PluginCompanyTypeAudiokineticExternal PluginCompanyType = 1
	PluginCompanyTypePluginDevMin PluginCompanyType = 64
		
	PluginCompanyTypePluginDevMax PluginCompanyType = 255
	PluginCompanyTypeMcDSP PluginCompanyType = 256
	PluginCompanyTypeWaveArts PluginCompanyType = 257
	PluginCompanyTypePhoneticArts PluginCompanyType = 258
	PluginCompanyTypeiZotope PluginCompanyType = 259
	PluginCompanyTypeCrankcaseAudio PluginCompanyType = 261
	PluginCompanyTypeIOSONO PluginCompanyType = 262
	PluginCompanyTypeAuroTechnologies PluginCompanyType = 263
	PluginCompanyTypeDolby PluginCompanyType = 264
	PluginCompanyTypeTwoBigEars PluginCompanyType = 265
	PluginCompanyTypeOculus PluginCompanyType = 266
	PluginCompanyTypeBlueRippleSound PluginCompanyType = 267
	PluginCompanyTypeEnzienAudio PluginCompanyType = 268
	PluginCompanyTypeKrotosDehumanizer PluginCompanyType = 269
	PluginCompanyTypeNurulize PluginCompanyType = 270
	PluginCompanyTypeSuperPowered PluginCompanyType = 271
	PluginCompanyTypeGoogle PluginCompanyType = 272
	PluginCompanyTypeNVIDIA PluginCompanyType = 273
	PluginCompanyTypeReserved PluginCompanyType = 274
	PluginCompanyTypeMicrosoft PluginCompanyType = 275
	PluginCompanyTypeYAMAHA PluginCompanyType = 276
	PluginCompanyTypeVisiSonics PluginCompanyType = 277
 
	// Unoffical
  	PluginCompanyTypeUbisoft PluginCompanyType = 128
  	PluginCompanyTypeCDProjektRED PluginCompanyType = 666
	PluginCompanyTypeInvalid PluginCompanyType = 0xFFFF
)

var PluginCompanyNames map[PluginCompanyType]string = map[PluginCompanyType]string {
	0: "Audiokinetic",
	1: "AudiokineticExternal",
	64: "PluginDevMin",
		
	255: "PluginDevMax",
	256: "McDSP",
	257: "WaveArts",
	258: "PhoneticArts",
	259: "iZotope",
	261: "CrankcaseAudio",
	262: "IOSONO",
	263: "AuroTechnologies",
	264: "Dolby",
	265: "TwoBigEars",
	266: "Oculus",
	267: "BlueRippleSound",
	268: "EnzienAudio",
	269: "KrotosDehumanizer",
	270: "Nurulize",
	271: "SuperPowered",
	272: "Google",
	273: "NVIDIA",
	274: "Reserved",
	275: "Microsoft",
	276: "YAMAHA",
	277: "VisiSonics",
	128: "Ubisoft",
	666: "CDProjektRED",
}

type FxCustom struct {
	Id              uint32
	PluginTypeId    uint32
	// Present if PluginID >= 0
	PluginParam     *PluginParam
	// uNumBankData uint8
	MediaMap        []MediaMapItem
	RTPC            RTPC
	StateProp       StateProp
	StateGroup      StateGroup
	// NumValues    uint16
	PluginProps     []PluginProp
}

const SizeOfMediaMapItem = 5
type MediaMapItem struct {
	Index    uint8
	SourceId uint32
}

const SizeOfPluginProp = 6
type PluginProp struct {
	PropertyID RTPCParameterType
	RTPCAccum  RTPCAccumType
	Value      float32
}

func (h *FxCustom) HasParam() bool {
	return h.PluginTypeId >= 0
}

func (h *FxCustom) assert() {
	if !h.HasParam() {
		assert.Nil(h.PluginParam,
			"Plugin Type ID indicate that there's no plugin parameter data.",
		)
	}
}

func (h *FxCustom) Encode() []byte {
	h.assert()
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.Append(HircTypeFxCustom)
	w.Append(dataSize)
	w.Append(h.Id)
	w.Append(h.PluginTypeId)
	if h.PluginParam != nil {
		w.AppendBytes(h.PluginParam.Encode())
	}
	w.Append(uint8(len(h.MediaMap)))
	for _, i := range h.MediaMap {
		w.Append(i)
	}
	w.AppendBytes(h.RTPC.Encode())
	w.AppendBytes(h.StateProp.Encode())
	w.AppendBytes(h.StateGroup.Encode())
	w.Append(uint16(len(h.PluginProps)))
	for _, p := range h.PluginProps {
		w.Append(p)
	}
	return w.BytesAssert(int(size))
}

func (h *FxCustom) DataSize() uint32 {
	size := 8 + 1 + uint32(len(h.MediaMap)) * SizeOfMediaMapItem + h.RTPC.Size() + h.StateProp.Size() + h.StateGroup.Size() + 2 + uint32(len(h.PluginProps)) * SizeOfPluginProp
	if h.PluginParam != nil {
		size += h.PluginParam.Size()
	}
	return size
}

func (h *FxCustom) PluginType() PluginType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginTypeInvalid
	}
	return PluginType((h.PluginTypeId >> 0) & 0x000F)
}

func (h *FxCustom) PluginCompany () PluginCompanyType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginCompanyTypeInvalid
	}
	return PluginCompanyType((h.PluginTypeId >> 4) & 0x03FF)
}

func (h *FxCustom) BaseParameter() *BaseParameter { return nil }

func (h *FxCustom) HircType() HircType { return HircTypeFxCustom }

func (h *FxCustom) HircID() (uint32, error) { return h.Id, nil }

func (h *FxCustom) IsCntr() bool { return false }

func (h *FxCustom) NumLeaf() int { return 0 }

func (h *FxCustom) ParentID() uint32 { return 0 }

func (h *FxCustom) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *FxCustom) RemoveLeaf(o HircObj) { panic("Panic Trap") }
