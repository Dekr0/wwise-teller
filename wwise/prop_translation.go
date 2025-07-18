package wwise

import "fmt"

type PropType uint8

func ForwardTranslateProp(p PropType, v int) uint8 {
	if v < 150 {
		tp, in := ForwardTranslationV128[p]
		if !in {
			panic(fmt.Sprintf("Failed to translate version 128 property ID %d", p))
		}
		return tp
	}
	if v >= 150 && v < 154 {
		panic("Translation is not implemented between version 150 (inclusive) and version 154 (exclusive)")
	}
	if v >= 154 {
		tp, in := ForwardTranslationV154[p]
		if !in {
			panic(fmt.Sprintf("Failed to translate version 154 property ID %d", p))
		}
		return tp
	}
	panic(fmt.Sprintf("Forward Translation is not implemented for version %d", v))
}

func InverseTranslateProp(p uint8, v int) PropType {
	if v < 150 {
		tp, in := InverseTranslationV128[p]
		if !in {
			panic(fmt.Sprintf("Failed to inverse translate version 128 property ID %d", p))
		}
		return tp
	}
	if v >= 150 && v < 154 {
		panic("Inverse translation is not implemented between version 150 (inclusive) and version 154 (exclusive)")
	}
	if v >= 154 {
		tp, in := InverseTranslationV154[p]
		if !in {
			panic(fmt.Sprintf("Failed to inverse translate version 154 property ID %d", p))
		}
		return tp
	}
	panic(fmt.Sprintf("Inverse Translation is not implemented for version %d", v))
}

func PropLabel(p PropType) string {
	name, in := TranslateName[p]
	if !in {
		return fmt.Sprintf("Unknown %d", p)
	}
	return name
}

var BasePropType []PropType = []PropType{
	TVolume,
	TPitch,
	TLPF,
	THPF,
	TMakeUpGain,
	TGameAuxSendVolume,
	TInitialDelay,
}

var BaseRangePropType []PropType = []PropType {
	TVolume,
	TPitch,
	TLPF,
	THPF,
	TMakeUpGain,
	TInitialDelay,
}

var UserAuxSendVolumePropType []PropType = []PropType {
	TUserAuxSendVolume0,
	TUserAuxSendVolume1,
	TUserAuxSendVolume2,
	TUserAuxSendVolume3,
}
