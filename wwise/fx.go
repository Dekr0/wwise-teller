package wwise

// --- struct definition --- //

type FXs struct {
	BypassAll   u8 // > 145 
	FXs       []FX
}

type FX struct {
	Id         u32 
	Idx        u8
	IsShareSet u8 // <= 145
	IsRender   u8 // <= 145
	BitVector  u8 // > 145
}

type FXsComponent struct {
	ActorMixerFXs map[u32]*FXs
}

type FxMetadatas struct {
	Idx        []u8
	Id         []u32
	IsShareSet []u8
}

type FxMetadataS struct {
	Idx        u8
	IsShareSet u8
	Id         u32
}

type FxMetadataComponent struct {
	ActorMixerFxMetadata map[u32]*FxMetadatas
}

// --- struct allocation --- //
func AllocFXs(numFX u8) *FXs {
	return &FXs{
		FXs: make([]FX, numFX, numFX),
	}
}

func AllocFxMetadatas(numFxMetadata u8) *FxMetadatas {
	return &FxMetadatas{
		Idx: make([]u8, numFxMetadata, numFxMetadata),
		Id: make([]u32, numFxMetadata, numFxMetadata),
		IsShareSet: make([]u8, numFxMetadata, numFxMetadata),
	}
}
