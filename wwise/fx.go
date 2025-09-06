package wwise

// This seesm to be most optimal layout so far despite 21% is wasted.
type FXs struct {
	BypassAll   u8 // > 145 
	FXs       []FX
}

type FX struct {
	// Id should be after Idx but move upforward due to alignment
	Id         u32 
	Idx        u8
	// <= 145
	IsShareSet u8
	IsRender   u8
	// > 145
	BitVector  u8
}

type FxMetadatas struct {
	FxMetadatas    []FxMetadata
}

type FxMetadata struct {
	// The upper 8 bits is IsShareSet; The lower 8 bits is Idx
	Idx        u16
	Id         u32
	// IsShareSet u8 // Layout 1
}
