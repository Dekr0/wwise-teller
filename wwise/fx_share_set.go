package wwise

import (
	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

var PluginNameLUT map[int32]string = map[int32]string{

	  -1:         "None", // found in early banks with no id
	  0x00000000: "None", // found in early banks with no id
	
	  // AKCODECID
	  0x00000001: "BANK",
	  0x00010001: "PCM",
	  0x00020001: "ADPCM",
	  0x00030001: "XMA",
	  0x00040001: "VORBIS",
	  0x00050001: "WIIADPCM",
	  // 0x00060001: "?",
	  0x00070001: "PCMEX", // "Standard PCM WAV file parser for Wwise Authoring" (same as PCM with another codec)
	  0x00080001: "EXTERNAL_SOURCE", // "unknown encoding" (.wem can be anything defined at runtime)
	  0x00090001: "XWMA",
	  0x000A0001: "AAC",
	  0x000B0001: "FILE_PACKAGE", // "File package files generated by the File Packager utility."
	  0x000C0001: "ATRAC9",
	  0x000D0001: "VAG/HE-VAG",
	  0x000E0001: "PROFILERCAPTURE", // "Profiler capture file (.prof)"
	  0x000F0001: "ANALYSISFILE",
	  0x00100001: "MIDI", // .wmid (modified .mid)
	  0x00110001: "OPUSNX", // originally just OPUS
	  0x00120001: "CAF", // unused?
	  0x00130001: "OPUS",
	  0x00140001: "OPUS_WEM",
	  0x00150001: "OPUS_WEM", // "Memory stats file as written through AK::MemoryMgr::DumpToFile()"
	  0x00160001: "SONY360", // unused/internal? '360 Reality Audio', MPEG-H derived?
	
	  // other types
	  0x00640002: "Wwise Sine", // AkSineTone
	  0x00650002: "Wwise Silence", // AkSilenceGenerator
	  0x00660002: "Wwise Tone Generator", // AkToneGen
	  0x00670003: "Wwise ?", // [The Lord of the Rings: Conquest (Wii)]
	  0x00680003: "Wwise ?", // [KetnetKick 2 (PC), The Lord of the Rings: Conquest (Wii)]
	  0x00690003: "Wwise Parametric EQ", // AkParametricEQ
	  0x006A0003: "Wwise Delay", // AkDelay
	  0x006C0003: "Wwise Compressor", // AkCompressor
	  0x006D0003: "Wwise Expander", //
	  0x006E0003: "Wwise Peak Limiter", // AkPeakLimiter
	  0x006F0003: "Wwise ?", // [Tony Hawk's Shred (Wii)]
	  0x00700003: "Wwise ?", // [Tony Hawk's Shred (Wii)]
	  0x00730003: "Wwise Matrix Reverb", // AkMatrixReverb
	  0x00740003: "SoundSeed Impact", //
	  0x00760003: "Wwise RoomVerb", // AkRoomVerb
	  0x00770002: "SoundSeed Air Wind", // AkSoundSeedAir
	  0x00780002: "SoundSeed Air Woosh", // AkSoundSeedAir
	  0x007D0003: "Wwise Flanger", // AkFlanger
	  0x007E0003: "Wwise Guitar Distortion", // AkGuitarDistortion
	  0x007F0003: "Wwise Convolution Reverb", // AkConvolutionReverb
	  0x00810003: "Wwise Meter", // AkSoundEngineDLL
	  0x00820003: "Wwise Time Stretch", // AkTimeStretch
	  0x00830003: "Wwise Tremolo", // AkTremolo
	  0x00840003: "Wwise Recorder", //
	  0x00870003: "Wwise Stereo Delay", // AkStereoDelay
	  0x00880003: "Wwise Pitch Shifter", // AkPitchShifter
	  0x008A0003: "Wwise Harmonizer", // AkHarmonizer
	  0x008B0003: "Wwise Gain", // AkGain
	  0x00940002: "Wwise Synth One", // AkSynthOne
	  0x00AB0003: "Wwise Reflect", // AkReflect
	
	  0x00AE0007: "System", // DefaultSink
	  0x00B00007: "Communication", // DefaultSink
	  0x00B10007: "Controller Headphones", // DefaultSink
	  0x00B30007: "Controller Speaker", // DefaultSink
	  0x00B50007: "No Output", // DefaultSink
	  0x03840009: "Wwise System Output Settings", // DefaultSink
	  0x00B70002: "SoundSeed Grain", //
	  0x00BA0003: "Mastering Suite", // MasteringSuite
	  0x00C80002: "Wwise Audio Input", // AkAudioInput
	  0x01950002: "Wwise Motion Generator", // AkMotion (used in CAkSound, v128>= / v130<=?)
	  0x01950005: "Wwise Motion Generator", // AkMotion (used in CAkFeedbackNode, v125<=)
	  0x01990002: "Wwise Motion Source", // AkMotion (used in CAkSound, v132>=)
	  0x01990005: "Wwise Motion Source?", // AkMotion
	  0x01FB0007: "Wwise Motion ?", // AkMotion
	
	  0x044C1073: "Auro Headphone", // Auro
	
	  // other companies (IDs can be repeated)
	  0x00671003: "McDSP ML1", // McDSP
	  0x006E1003: "McDSP FutzBox", //
	
	  0x00021033: "iZotope Hybrid Reverb", //
	  0x00031033: "iZotope Trash Distortion", //
	  0x00041033: "iZotope Trash Delay", //
	  0x00051033: "iZotope Trash Dynamics Mono", //
	  0x00061033: "iZotope Trash Filters", //
	  0x00071033: "iZotope Trash Box Modeler", //
	  0x00091033: "iZotope Trash Multiband Distortion", //
	
	  0x006E0403: "Platinum MatrixSurroundMk2", // PgMatrixSurroundMk2
	  0x006F0403: "Platinum LoudnessMeter", // PgLoudnessMeter
	  0x00710403: "Platinum SpectrumViewer", // PgSpectrumViewer
	  0x00720403: "Platinum EffectCollection", // PgEffectCollection
	  0x00730403: "Platinum MeterWithFilter", // PgMeterWithFilter
	  0x00740403: "Platinum Simple3D", // PgSimple3D
	  0x00750403: "Platinum Upmixer", // PgUpmixer
	  0x00760403: "Platinum Reflection", // PgReflection
	  0x00770403: "Platinum Downmixer", // PgDownmixer
	  0x00780403: "Platinum Flex?", // PgFlex? [Nier Automata] 
	
	  0x00020403: "Codemasters ? Effect", //  [Dirt Rally (PS4)]
	
	  0x00640332: "Ubisoft ?", // [Mario + Rabbids DLC 3]
	  0x04F70803: "Ubisoft ? Effect", // [AC Valhalla]
	  0x04F80806: "Ubisoft ? Mixer", // [AC Valhalla]
	  0x04F90803: "Ubisoft ? Effect", // [AC Valhalla]
	
	  0x00AA1137: "Microsoft Spatial Sound", // MSSpatial
	
	  0x000129A3: "CPR Simple Delay", // CDPSimpleDelay
	  0x000229A2: "CPR Voice Broadcast Receive ?", // CDPVoiceBroadcastReceive
	  0x000329A3: "CPR Voice Broadcast Send ?", // CDPVoiceBroadcastSend
	  0x000429A2: "CPR Voice Broadcast Receive ?", // CDPVoiceBroadcastReceive
	  0x000529A3: "CPR Voice Broadcast Send ?", // CDPVoiceBroadcastSend
	
	  0x01A01052: "Crankcase REV Model Player", // CrankcaseAudioREVModelPlayer
}

type FxShareSet struct {
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

func (h *FxShareSet) HasParam() bool {
	return h.PluginTypeId >= 0
}

func (h *FxShareSet) assert() {
	if !h.HasParam() {
		assert.Nil(h.PluginParam,
			"Plugin Type ID indicate that there's no plugin parameter data.",
		)
	}
}

func (h *FxShareSet) Encode() []byte {
	h.assert()
	dataSize := h.DataSize()
	size := SizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.Append(HircTypeFxShareSet)
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

func (h *FxShareSet) DataSize() uint32 {
	size := 8 + 1 + uint32(len(h.MediaMap)) * SizeOfMediaMapItem + h.RTPC.Size() + h.StateProp.Size() + h.StateGroup.Size() + 2 + uint32(len(h.PluginProps)) * SizeOfPluginProp
	if h.PluginParam != nil {
		size += h.PluginParam.Size()
	}
	return size
}

func (h *FxShareSet) PluginType() PluginType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginTypeInvalid
	}
	return PluginType((h.PluginTypeId >> 0) & 0x000F)
}

func (h *FxShareSet) PluginCompany () PluginCompanyType {
	if h.PluginTypeId == 0xFFFFFFFF {
		return PluginCompanyTypeInvalid
	}
	return PluginCompanyType((h.PluginTypeId >> 4) & 0x03FF)
}

func (h *FxShareSet) BaseParameter() *BaseParameter { return nil }

func (h *FxShareSet) HircType() HircType { return HircTypeFxShareSet }

func (h *FxShareSet) HircID() (uint32, error) { return h.Id, nil }

func (h *FxShareSet) IsCntr() bool { return false }

func (h *FxShareSet) NumLeaf() int { return 0 }

func (h *FxShareSet) ParentID() uint32 { return 0 }

func (h *FxShareSet) AddLeaf(o HircObj) { panic("Panic Trap") }

func (h *FxShareSet) RemoveLeaf(o HircObj) { panic("Panic Trap") }

func (h *FxShareSet) Leafs() []uint32 { return []uint32{} }
