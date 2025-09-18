package wwise

import "fmt"

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
	FXs map[u32]*FXs
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

type FxMetadataE struct {
	Idx        u8
	Id         u32
	IsShareSet u8
}

type FxMetadatasComponent struct {
	FxMetadatas map[u32]*FxMetadatas
}

// --- struct allocation --- //

func AllocFXsComponent(numActorMixer u32) (c *FXsComponent) {
	c = &FXsComponent{}
	if numActorMixer <= 0 {
		c.FXs = make(map[u32]*FXs)
	} else {
		c.FXs = make(map[u32]*FXs, numActorMixer)
	}
	return c
}

func AllocFxMetadataComponent(numActorMixer u32) (c *FxMetadatasComponent) {
	c = &FxMetadatasComponent{}
	if numActorMixer <= 0 {
		c.FxMetadatas = make(map[u32]*FxMetadatas)
	} else {
		c.FxMetadatas = make(map[u32]*FxMetadatas, numActorMixer)
	}
	return c
}

func AllocFXs(numFX u8) *FXs {
	return &FXs{
		FXs: make([]FX, numFX, numFX),
	}
}

func AllocFxMetadatas(numFxMetadatas u8) *FxMetadatas {
	return &FxMetadatas{
		Idx: make([]u8, numFxMetadatas, numFxMetadatas),
		Id: make([]u32, numFxMetadatas, numFxMetadatas),
		IsShareSet: make([]u8, numFxMetadatas, numFxMetadatas),
	}
}

// --- assertion --- //

func AssertFXs(f *FXs) error {
	if len(f.FXs) > 255 {
		return fmt.Errorf("# of FXs are greater than 255")
	}
	return nil
}

func AssertFxMetadatas(f *FxMetadatas) error {
	if len(f.Idx) > 255 {
		return fmt.Errorf("# of FX Metadata indices are greateer than 255")
	}
	if len(f.Idx) != len(f.Id) {
		return fmt.Errorf(
			"# of index (%d) in FX Metadatas does not equal to # of id (%d) in FX Metadatas",
			len(f.Idx), len(f.Id),
		)
	}
	if len(f.Id) != len(f.IsShareSet) {
		return fmt.Errorf(
			"# of id (%d) in FX Metadata does not equal to # of share set bit vector (%d) in FX Metadatas",
			len(f.Id), len(f.IsShareSet),
		)
	}
	return nil
}

// --- sizing --- //

func SizeOfFXs(f *FXs, version u32) (size u32) {
	size = 1
	if len(f.FXs) <= 0 {
		return size
	}
	if version <= 145 {
		size += 1 + u32(len(f.FXs)) * SizeOfFXLE145
	} else {
		size += 1 + u32(len(f.FXs)) * SizeOfFXG145
	}
	return size
}

func SizeOfFxMetadatas(f *FxMetadatas) u32 {
	return 1 + u32(len(f.Idx)) * Size8 + u32(len(f.Id)) * Size32 + u32(len(f.IsShareSet)) * Size8
}

// --- encoding --- //

func EncodeFXs(e *HircEncoderCtx, f *FXs) (err error) {
	curr := e.Encoder.Count
	size := SizeOfFXs(f, e.Version)
	if err := e.Primitive(u8(len(f.FXs))); err != nil {
		return fmt.Errorf("Failed to encode # of FXs counter value in FXs: %w", err)
	}
	if len(f.FXs) <= 0 {
		return e.Expect(curr, size)
	}
	if err := e.Primitive(f.BypassAll); err != nil {
		return fmt.Errorf("Failed to encode Bypass All value of FXs: %w", err)
	}
	for i, f := range f.FXs {
		if e.Version <= 145 {
			type FX struct {
				Idx        u8
				Id         u32
				IsShareSet u8
				IsRender   u8
			}
			if err := e.Struct(FX{
				f.Idx, f.Id, f.IsShareSet, f.IsRender,
			}, SizeOfFXLE145); err != nil {
				return fmt.Errorf("Failed to encode %d FX in FXs: %w", i, err)
			}
		} else {
			type FX struct {
				Idx       u8
				Id        u32
				BitVector u8
			}
			if err := e.Struct(FX{
				f.Idx, f.Id, f.BitVector,
			}, SizeOfFXG145); err != nil {
				return fmt.Errorf("Failed to encode %d FX in FXs: %w", i, err)
			}
		}
	}
	return e.Expect(curr, size)
}

func EncodeFxMetadatas(e *HircEncoderCtx, f *FxMetadatas) error {
	curr := e.Encoder.Count
	size := SizeOfFxMetadatas(f)
	if err := e.Primitive(u8(len(f.Id))); err != nil {
		return fmt.Errorf("Failed to encode # of FX Metadata value counter: %d", err)
	}
	for i, idx := range f.Idx {
		payload := FxMetadataE{
			idx,
			f.Id[i],
			f.IsShareSet[i],
		}
		if err := e.Struct(payload, SizeOfFXMetadata); err != nil {
			return fmt.Errorf("Failed to encode %d-th FX Metadata: %w", i, err)
		}
	}
	return e.Expect(curr, size)
}

// --- component assertion wrapper --- //

func (c *FXsComponent) AssertFXsById(internalId u32) error {
	f, in := c.FXs[internalId]
	if !in {
		return fmt.Errorf("Failed to locate FXs")
	}
	return AssertFXs(f)
}

func (c *FxMetadatasComponent) AssertFxMetadatasById(internalId u32) error {
	f, in := c.FxMetadatas[internalId];
	if !in {
		return fmt.Errorf("Failed to locate Fx Metadatas")
	}
	return AssertFxMetadatas(f)
}

// --- component sizing wrapper --- //

func (c *FXsComponent) SizeOfFXsById(internalId u32, version u32) u32 {
	f, in := c.FXs[internalId]
	if !in {
		panic("Failed to locate FXs")
	}
	return SizeOfFXs(f, version)
}

func (c *FxMetadatasComponent) SizeOfFxMetadasById(internalId u32) u32 {
	f, in := c.FxMetadatas[internalId]
	if !in {
		panic("Failed to locate FX Metadata")
	}
	return SizeOfFxMetadatas(f)
}

// --- component getter and setter --- //

func (c *FXsComponent) HasFXs(internalId u32) (in bool) {
	_, in = c.FXs[internalId]
	return in
}

func (c *FxMetadatasComponent) HasFxMetadatas(internalId u32) (in bool) {
	_, in = c.FxMetadatas[internalId]
	return in
}

func (c *FXsComponent) GetFXs(internalId u32) (f *FXs) {
	f, in := c.FXs[internalId]
	if !in {
		panic("Failed to locate FXs")
	}
	return f 
}

func (c *FxMetadatasComponent) GetFxMetadatas(internalId u32) (f *FxMetadatas) {
	f, in := c.FxMetadatas[internalId]
	if !in {
		panic("Failed to locate FxMetadatas")
	}
	return f 
}

func (c *FXsComponent) AddFXs(internalId u32, f *FXs) {
	if f == nil {
		panic("FXs is nil")
	}
	if _, in := c.FXs[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.FXs[internalId] = f
}

func (c *FxMetadatasComponent) AddFxMetadatas(internalId u32, f *FxMetadatas) {
	if f == nil {
		panic("FxMetadatas is nil")
	}
	if _, in := c.FxMetadatas[internalId]; in {
		panic(MonotonicIdCollision)
	}
	c.FxMetadatas[internalId] = f
}

// --- HIRC component wrapper --- //

func (h *HIRC) GetFXs(internalId u32) *FXs {
	return h.FXsComponent.GetFXs(internalId)
}

func (h *HIRC) GetFxMetadatas(internalId u32) (f *FxMetadatas) {
	return h.FxMetadatasComponent.GetFxMetadatas(internalId)
}

// --- constant value --- //

const SizeOfFXLE145    = Size32 + 3 * Size8
const SizeOfFXG145     = Size32 + 2 * Size8
const SizeOfFXMetadata = Size32 + 2 * Size8
