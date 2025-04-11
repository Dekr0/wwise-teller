package wwise

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
)

type BaseParameter struct {
	BitIsOverrideParentFx uint8
	FxChunk *FxChunk
	FxChunkMetadata *FxChunkMetadata
	BitOverrideAttachmentParams uint8
	OverrideBusId uint32
	DirectParentId uint32
	ByBitVectorA uint8
	PropBundle *PropBundle
	RangePropBundle *RangePropBundle
	PositioningParam *PositioningParam
	AuxParam *AuxParam
	AdvanceSetting *AdvanceSetting
	StateProp *StateProp
	StateGroup *StateGroup
	RTPC *RTPC
}

func (b *BaseParameter) Encode() []byte {
	fxChunkBlob := b.FxChunk.Encode()
	fxChunkMetadataBlob := b.FxChunkMetadata.Encode()
	propBundleBlob := b.PropBundle.Encode()
	rangeBundleBlob := b.RangePropBundle.Encode()
	positioningParamBlob := b.PositioningParam.Encode()
	auxParamBlob := b.AuxParam.Encode()
	statePropBlob := b.StateProp.Enocde()
	stateGroupBlob := b.StateGroup.Encode()
	rtpcBlob := b.RTPC.Encode()

	blobSize := 1 + len(fxChunkBlob) + len(fxChunkMetadataBlob) + 
				1 + 4 + 4 + 1 + 
			    len(propBundleBlob) + len(rangeBundleBlob) + 
				len(positioningParamBlob) + len(auxParamBlob) + ADVANCE_SETTING_FIELD_SIZE +
				len(statePropBlob) + len(stateGroupBlob) + len(rtpcBlob)
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))

	bw.AppendByte(b.BitIsOverrideParentFx)
	bw.AppendBytes(fxChunkBlob)
	bw.AppendBytes(fxChunkMetadataBlob)
	bw.AppendByte(b.BitOverrideAttachmentParams)
	bw.Append(b.OverrideBusId)
	bw.Append(b.DirectParentId)
	bw.AppendByte(b.ByBitVectorA)
	bw.AppendBytes(propBundleBlob)
	bw.AppendBytes(rangeBundleBlob)
	bw.AppendBytes(positioningParamBlob)
	bw.AppendBytes(auxParamBlob)
	bw.Append(b.AdvanceSetting)
	bw.AppendBytes(statePropBlob)
	bw.AppendBytes(stateGroupBlob)
	bw.AppendBytes(rtpcBlob)

	return bw.Flush(blobSize)
}

type FxChunk struct {
	UniqueNumFx uint8
	BitsFxByPass uint8
	FxChunkItems []*FxChunkItem
}

func NewFxChunk() *FxChunk {
	return &FxChunk{0, 0, []*FxChunkItem{}}
}

func (f *FxChunk) Encode() []byte {
	assert.AssertEqual(
		int(f.UniqueNumFx),
		len(f.FxChunkItems),
		"Unique FX counter doesn't equal to # of FX Item",
	)
	if f.UniqueNumFx <= 0 {
		return []byte{ f.UniqueNumFx }
	}

	blobSize := 1 + 1 + f.UniqueNumFx * FX_CHUNK_ITEM_FIELD_SIZE
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(f.UniqueNumFx)
	bw.AppendByte(f.BitsFxByPass)
	for _, fi := range f.FxChunkItems {
		bw.Append(fi)
	}
	return bw.Flush(int(blobSize))
}

const FX_CHUNK_ITEM_FIELD_SIZE = 7
type FxChunkItem struct {
	UniqueFxIndex uint8 // u8i
	FxId uint32 // tid
	BitIsShareSet uint8 // U8x
	BitIsRendered uint8 // U8x
}

type FxChunkMetadata struct {
	BitIsOverrideParentMetadata uint8
	UniqueNumFxMetadata uint8
	FxMetaDataChunkItems []*FxChunkMetadataItem
}

func NewFxChunkMetadata() *FxChunkMetadata {
	return &FxChunkMetadata{0, 0, []*FxChunkMetadataItem{}}
}

func (f *FxChunkMetadata) Encode() []byte {
	assert.AssertEqual(
		int(f.UniqueNumFxMetadata),
		len(f.FxMetaDataChunkItems),
		"Unique FX counter doesn't equal to # of FX Item",
	)
	blobSize := 1 + 1 + f.UniqueNumFxMetadata * FX_CHUNK_METADATA_FIELD_SIZE
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(f.BitIsOverrideParentMetadata)
	bw.AppendByte(f.UniqueNumFxMetadata)
	for _, fi := range f.FxMetaDataChunkItems {
		bw.Append(fi)
	}
	return bw.Flush(int(blobSize))
}

const FX_CHUNK_METADATA_FIELD_SIZE = 6
type FxChunkMetadataItem struct {
	UniqueFxIndex uint8 // u8i
	FxId uint32 // tid
	BitIsShareSet uint8 // U8x
}

const PROP_VALUE_FIELD_SIZE = 4
type PropBundle struct {
	CProps uint8 // u8i
	PIds []uint8 // CProps * u8i
	PValues [][]byte // CProps * (Union[tid, uni / float32])
}

func NewPropBundle() *PropBundle {
	return &PropBundle{0, []uint8{}, [][]byte{}}
}

func CreatePropValue[T comparable](val T) []byte {
	bw := reader.NewFixedSizeBlobWriter(4)
	bw.Append(val)
	return bw.Flush(4)
}

func (p *PropBundle) Encode() []byte {
	assert.AssertEqual(
		int(p.CProps), len(p.PIds), "Property counter does not equal to # of property IDs",
	)
	assert.AssertEqual(
		len(p.PIds), len(p.PValues), "# of property IDs does not equal to # of property values",
	)

	blobSize := 1 + len(p.PIds) + PROP_VALUE_FIELD_SIZE * len(p.PValues)
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))

	bw.AppendByte(p.CProps)
	bw.Append(p.PIds)
	for _, pValue := range p.PValues {
		bw.Append(pValue)
	}

	return bw.Flush(blobSize)
}

func (p *PropBundle) HasProp(pID uint8) bool {
	_, found := sort.Find(len(p.PIds), func(i int) int {
		if pID < p.PIds[i] {
			return -1
		} else if pID == p.PIds[i] {
			return 0
		} else {
			return 1
		}
	})
	return found
}

func (p *PropBundle) CheckCount(count uint8) bool {
	if p.CProps != count {
		return false
	}
	if len(p.PIds) != int(count) {
		return false
	}
	if len(p.PValues) != int(count) {
		return false
	}
	return true
}

func (p *PropBundle) UpdatePropUint32(pId uint8, val uint32) {
	insert, found := sort.Find(len(p.PIds), func(i int) int {
		if pId < p.PIds[i] {
			return -1
		} else if pId == p.PIds[i] {
			return 0
		} else {
			return 1
		}
	})
	bw := reader.NewFixedSizeBlobWriter(4)
	bw.Append(val)
	blob := bw.Flush(4)
	if !found {
		p.PIds = slices.Insert(p.PIds, insert, pId)
		p.PValues = slices.Insert(p.PValues, insert, blob)
	} else {
		p.PValues[insert] = blob
	}
}

func (p *PropBundle) UpdatePropFloat32(pId uint8, val float32) {
	insert, found := sort.Find(len(p.PIds), func(i int) int {
		if pId < p.PIds[i] {
			return -1
		} else if pId == p.PIds[i] {
			return 0
		} else {
			return 1
		}
	})
	bw := reader.NewFixedSizeBlobWriter(4)
	bw.Append(val)
	blob := bw.Flush(4)
	if !found {
		p.PIds = slices.Insert(p.PIds, insert, pId)
		p.PValues = slices.Insert(p.PValues, insert, blob)
	} else {
		p.PValues[insert] = blob
	}
}

func (p *PropBundle) RemoveProp(pId uint8) error {
	remove, found := sort.Find(len(p.PIds), func(i int) int {
		if pId < p.PIds[i] {
			return -1
		} else if pId == p.PIds[i] {
			return 0
		} else {
			return 1
		}
	})
	if !found {
		return fmt.Errorf("Failed to find property ID %d", pId)
	}
	p.PIds = slices.Delete(p.PIds, remove, remove+1)
	return nil
}

type RangePropBundle struct {
	CProps uint8 // u8i
	PIds []uint8 // CProps * u8i
	RangeValues []*RangeValue // CProps * sizeof(RangeValue)
}

func NewRangePropBundle() *RangePropBundle {
	return &RangePropBundle{0, []uint8{}, []*RangeValue{}}
}

func CreateRangeValue[T comparable](min T, max T) *RangeValue {
	bw := reader.NewFixedSizeBlobWriter(4)
	bw.Append(min)
	minBlob := bw.Flush(4)
	bw = reader.NewFixedSizeBlobWriter(4)
	bw.Append(max)
	maxBlob := bw.Flush(4)
	return &RangeValue{minBlob, maxBlob}
}

func (r *RangePropBundle) Encode() []byte {
	assert.AssertEqual(
		int(r.CProps), len(r.PIds), "Property counter does not equal to # of property IDs",
	)
	assert.AssertEqual(
		len(r.PIds), len(r.RangeValues), "# of property IDs does not equal to # of property values",
	)

	blobSize := 1 + len(r.PIds) + RANGE_VALUE_FIELD_SIZE * len(r.RangeValues)
	
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(r.CProps)
	bw.Append(r.PIds)
	for _, rangeValue := range r.RangeValues { bw.AppendBytes(rangeValue.Encode()) }

	return bw.Flush(blobSize)
}

const RANGE_VALUE_FIELD_SIZE = 4 + 4
type RangeValue struct {
	Min []byte // Union[tid, uni / float32]
	Max []byte // Union[tid, uni / float32]
}

func (r *RangeValue) Encode() []byte {
	assert.AssertEqual(4, len(r.Min), "Min of range value has incorrect size")
	assert.AssertEqual(4, len(r.Max), "Max of range value has incorrect size")

	blob := make([]byte, 0, RANGE_VALUE_FIELD_SIZE)
	blob = append(blob, r.Min...) 
	blob = append(blob, r.Max...)

	assert.AssertEqual(
		RANGE_VALUE_FIELD_SIZE, len(blob),
		"Encoded data of RangeValue has incorrect size",
	)
	return blob
}

type PositioningParam struct {
	BitsPositioning uint8 // U8x
	Bits3D uint8 // U8x
	PathMode uint8 // U8x
	TransitionTime int32 // s32
	NumPositionVertices uint32 // u32
	PositionVertices []*PositionVertex // NumPositionVertices * sizeof(PositionVertex)
	NumPositionPlayListItem uint32 // u32
	PositionPlayListItems []*PositionPlayListItem // NumPositionPlayListItem * sizeof(PositionPlayListItem)
	Ak3DAutomationParams []*Ak3DAutomationParam // NumPositionPlayListItem * sizeof(Ak3DAutomationParams)
}

func NewPositioningParam() *PositioningParam {
	return &PositioningParam{
		0, 0, 0, 0, 
		0, []*PositionVertex{}, 
		0, []*PositionPlayListItem{}, []*Ak3DAutomationParam{},
	}
}

func (p *PositioningParam) HasPositioning() bool {
	return (p.BitsPositioning >> 0) & 1 != 0
}

func (p *PositioningParam) Has3D() bool {
	if !p.HasPositioning() {
		return false
	}
	return (p.BitsPositioning >> 1) & 1 != 0
}

func (p *PositioningParam) HasPositioningAnd3D() bool {
	return p.HasPositioning() && p.Has3D()
}

func (p *PositioningParam) Get3DPositionType() (uint8, error) {
	if !p.HasPositioning() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter is not enable.",
		)
	}
	if !p.Has3D() {
		return 0, errors.New(
			"Failed to get 3D position type: positioning parameter does not enable 3D setting",
		)
	}
	return (p.BitsPositioning >> 5) & 3, nil
}

func (p *PositioningParam) HasAutomation() bool {
	if !p.HasPositioningAnd3D() {
		return false
	}
	_3DPositioningType, err := p.Get3DPositionType()
	assert.AssertNil(err, "Error of Get3DPositionType")
	return p.HasPositioningAnd3D() && _3DPositioningType != 0
}

func (p *PositioningParam) Encode() []byte {
	p.validateIntegrity()

	if !p.HasPositioning() || !p.HasPositioningAnd3D() {
		bw := reader.NewFixedSizeBlobWriter(1)
		bw.AppendByte(p.BitsPositioning)
		return bw.Flush(1)
	}

	if !p.HasAutomation() {
		bw := reader.NewFixedSizeBlobWriter(2)
		bw.AppendByte(p.BitsPositioning)
		bw.AppendByte(p.Bits3D)
		return bw.Flush(2)
	}

	blobSize := 1 + 1 + 1 + 4 + 
				4 + len(p.PositionVertices) * POSITION_VERTEX_FIELD_SIZE +
				4 + len(p.PositionPlayListItems) * POSITION_PLAY_LIST_ITEM_FIELD_SIZE +
				len(p.PositionPlayListItems) * AK_3D_AUTOMATION_PARAM_FIELD_SIZE

	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(p.BitsPositioning)
	bw.Append(p.Bits3D)
	bw.Append(p.PathMode)
	bw.Append(p.TransitionTime)
	bw.Append(p.NumPositionVertices)
	for _, v := range p.PositionVertices { bw.Append(v) }
	bw.Append(p.NumPositionPlayListItem)
	for _, i := range p.PositionPlayListItems { bw.Append(i) }
	for _, p := range p.Ak3DAutomationParams { bw.Append(p) }

	return bw.Flush(blobSize)
}

/* Will Panic */
func (p *PositioningParam) validateIntegrity() {
	/* TODO, document assertion */
	if !p.HasPositioning() {
		assert.AssertTrue(((p.BitsPositioning >> 1) & 1) == 0, "") // No 3D
		assert.AssertTrue(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.AssertTrue(p.Bits3D == 0, "")
		assert.AssertTrue(p.PathMode == 0, "")
		assert.AssertTrue(p.TransitionTime == 0, "")
		assert.AssertTrue(p.NumPositionVertices == 0, "")
		assert.AssertTrue(len(p.PositionVertices) == 0, "")
		assert.AssertTrue(p.NumPositionPlayListItem == 0, "")
		assert.AssertTrue(len(p.PositionPlayListItems) == 0, "")
		assert.AssertTrue(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	if !p.HasPositioningAnd3D() {
		assert.AssertTrue(((p.BitsPositioning >> 1) & 1) == 0, "") // No 3D
		assert.AssertTrue(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.AssertTrue(p.Bits3D == 0, "")
		assert.AssertTrue(p.PathMode == 0, "")
		assert.AssertTrue(p.TransitionTime == 0, "")
		assert.AssertTrue(p.NumPositionVertices == 0, "")
		assert.AssertTrue(len(p.PositionVertices) == 0, "")
		assert.AssertTrue(p.NumPositionPlayListItem == 0, "")
		assert.AssertTrue(len(p.PositionPlayListItems) == 0, "")
		assert.AssertTrue(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	if !p.HasAutomation() {
		assert.AssertTrue(((p.BitsPositioning >> 5) & 3) == 0, "") // 3D Position Type is 0
		assert.AssertTrue(p.PathMode == 0, "")
		assert.AssertTrue(p.TransitionTime == 0, "")
		assert.AssertTrue(p.NumPositionVertices == 0, "")
		assert.AssertTrue(len(p.PositionVertices) == 0, "")
		assert.AssertTrue(p.NumPositionPlayListItem == 0, "")
		assert.AssertTrue(len(p.PositionPlayListItems) == 0, "")
		assert.AssertTrue(len(p.Ak3DAutomationParams) == 0, "")
		return
	}
	assert.AssertTrue(
		len(p.PositionVertices) == int(p.NumPositionVertices),
		"Position vertices counter doesn't equal of # of position vertices",
	)
	assert.AssertTrue(
		len(p.PositionPlayListItems) == int(p.NumPositionPlayListItem),
		"Position play list item counter doesn't equal of # of position play list items",
	)
	assert.AssertTrue(
		len(p.Ak3DAutomationParams) == int(p.NumPositionPlayListItem),
		"Position play list item counter doesn't equal of # of 3D automation parameters",
	)
}

const POSITION_VERTEX_FIELD_SIZE = 4 * 4
type PositionVertex struct {
	X float32 // f32
	Y float32 // f32
	Z float32 // f32
	Duration int32 // s32
}

const POSITION_PLAY_LIST_ITEM_FIELD_SIZE = 4 * 2
type PositionPlayListItem struct {
	UniqueVerticesOffset uint32 // U32
	INumVertices uint32 // u32
}

const AK_3D_AUTOMATION_PARAM_FIELD_SIZE = 4 * 3
type Ak3DAutomationParam struct {
	XRange float32 // f32
	YRange float32 // f32
	ZRange float32 // f32
}

type AuxParam struct {
	AuxBitVector uint8 // U8x
	AuxIds []uint32 // 4 * tid
	ReflectionAuxBus uint32 // tid
}

func NewAuxParam() *AuxParam {
	return &AuxParam{0, []uint32{}, 0}
}

func (a *AuxParam) HasAux() bool {
	return a.AuxBitVector & 0b00001000 != 0
}

func (a *AuxParam) Encode() []byte {
	a.validateIntegrity()

	if !a.HasAux() {
		bw := reader.NewFixedSizeBlobWriter(1 + 4)
		bw.AppendByte(a.AuxBitVector)
		bw.Append(a.ReflectionAuxBus)
		return bw.Flush(5)
	}

	blobSize := 1 + len(a.AuxIds) * 4 + 4
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(a.AuxBitVector)
	for _, id := range a.AuxIds { bw.Append(id) }
	bw.Append(a.ReflectionAuxBus)

	return bw.Flush(blobSize)
}

func (a *AuxParam) validateIntegrity() {
	if !a.HasAux() {
		assert.AssertTrue(
			len(a.AuxIds) <= 0, 
			"Aux bit vector indicate there is no auxiliary send but # of Aux IDs" +
			" is not equal to 0.",
		)
		return
	}
	assert.AssertTrue(
		len(a.AuxIds) == 4,
		"Each auxiliary parameter should only have 4 auxiliary IDs",
	)
}

const ADVANCE_SETTING_FIELD_SIZE = 1 + 1 + 2 + 1 + 1
type AdvanceSetting struct {
	AdvanceSettingBitVector uint8 // U8x
	VirtualQueueBehavior uint8 // U8x
	MaxNumInstance uint16 // u16
	BelowThresholdBehavior uint8 // U8x
	HDRBitVector uint8 // U8x
}

type StateProp struct {
	NumStateProps uint8
	StatePropItems []*StatePropItem
}

func NewStateProp() *StateProp {
	return &StateProp{0, []*StatePropItem{}}
}

func (s *StateProp) Enocde() []byte {
	assert.AssertEqual(
		int(s.NumStateProps),
		len(s.StatePropItems),
		"State props counter doesn't equal to # of state prop items",
	)
	blobSize := 1 + len(s.StatePropItems) * STATE_PROP_ITEM_FIELD_SIZE
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(s.NumStateProps)
	for _, si := range s.StatePropItems {
		bw.Append(si)
	}
	return bw.Flush(blobSize)
}

const STATE_PROP_ITEM_FIELD_SIZE = 1 * 3
type StatePropItem struct {
	PropertyId uint8 // var (at least 1 byte / 8 bits)
	AccumType uint8 // U8x
	InDb uint8 // U8x
}

type StateGroup struct {
	NumStateGroups uint8
	StateGroupItems []*StateGroupItem
}

func NewStateGroup() *StateGroup {
	return &StateGroup{0, []*StateGroupItem{}}
}

func (s *StateGroup) Encode() []byte {
	assert.AssertEqual(
		int(s.NumStateGroups),
		len(s.StateGroupItems),
		"State groups counter doesn't equal to # of state group items",
	)
	stateGroupItemBlobs := make([][]byte, len(s.StateGroupItems))
	blobSize := 1
	for i, si := range s.StateGroupItems {
		stateGroupItemBlobs[i] = si.Encode()
		blobSize += len(stateGroupItemBlobs[i])
	}
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.AppendByte(s.NumStateGroups)
	bw.AppendBytes(bytes.Join(stateGroupItemBlobs, []byte{}))
	return bw.Flush(blobSize)
}

type StateGroupItem struct {
	StateGroupID uint32 // tid
	StateSyncType uint8 // U8x
	NumStates uint8 // var (assume at least 1 byte / 8 bits, can be more)
	States []*StateGroupItemState // NumStates * sizeof(StateGroupItemState)
}

func NewStateGroupItem() *StateGroupItem {
	return &StateGroupItem{0, 0, 0, []*StateGroupItemState{}}
}

func (s * StateGroupItem) Encode() []byte {
	assert.AssertEqual(
		int(s.NumStates),
		len(s.States),
		"State counter doesn't equal to # of states",
	)
	blobSize := 4 + 1 + 1 + STATE_GROUP_ITEM_STATE_SIZE * len(s.States)
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(s.StateGroupID)
	bw.AppendByte(s.StateSyncType)
	bw.AppendByte(s.NumStates)
	for _, state := range s.States {
		bw.Append(state)
	}
	return bw.Flush(blobSize)
}

const STATE_GROUP_ITEM_STATE_SIZE = 4 * 2 
type StateGroupItemState struct {
	StateID uint32 // tid
	StateInstanceID uint32 // tid
}

type RTPC struct {
	NumRTPC uint16
	RTPCItems []*RTPCItem
}

func NewRTPC() *RTPC {
	return &RTPC{0, []*RTPCItem{}}
}

func (r *RTPC) Encode() []byte {
	assert.AssertEqual(
		int(r.NumRTPC),
		len(r.RTPCItems),
		"RTPC counter doesn't equal to # of RTPC items",
	)
	rtpcItemBlobs := make([][]byte, r.NumRTPC)
	blobSize := 2
	for i, ri := range r.RTPCItems {
		rtpcItemBlobs[i] = ri.Encode()
		blobSize += len(rtpcItemBlobs[i])
	}
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(r.NumRTPC)
	bw.AppendBytes(bytes.Join(rtpcItemBlobs, []byte{}))
	return bw.Flush(blobSize)
}

type RTPCItem struct {
	RTPCID uint32 // tid
	RTPCType uint8 // U8x
	RTPCAccum uint8 // U8x
	ParamID uint8 // var (assume at least 1 byte / 8 bits, can be more)
	RTPCCurveID uint32 // sid
	Scaling uint8 // U8x
	NumRTPCGraphPoints uint16 // u16
	RTPCGraphPoints []*RTPCGraphPoint 
}

func NewRTPCItem() *RTPCItem {
	return &RTPCItem{0, 0, 0, 0, 0, 0, 0, []*RTPCGraphPoint{}}
}

func (r *RTPCItem) Encode() []byte {
	assert.AssertEqual(
		int(r.NumRTPCGraphPoints),
		len(r.RTPCGraphPoints),
		"RTPC graph point counter doesn't equal to the # of RTPC graph points",
	)
	blobSize := 4 + 1 + 1 + 1 + 4 + 1 + 2 + len(r.RTPCGraphPoints) * RTPC_GRAPH_POINT_FIELD_SIZE
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(r.RTPCID)
	bw.AppendByte(r.RTPCType)
	bw.AppendByte(r.RTPCAccum)
	bw.AppendByte(r.ParamID)
	bw.Append(r.RTPCCurveID)
	bw.AppendByte(r.Scaling)
	bw.Append(r.NumRTPCGraphPoints)
	for _, pt := range r.RTPCGraphPoints {
		bw.Append(pt)
	}
	return bw.Flush(blobSize)
}

const RTPC_GRAPH_POINT_FIELD_SIZE = 4 * 3
type RTPCGraphPoint struct {
	From float32 // f32 
	To float32 // f32
	Interp uint32 // U32
}

type BankSourceData struct {
	PluginID uint32 // U32
	StreamType uint8 // U8x
	SourceID uint32 // tid
	InMemoryMediaSize uint32 // U32
	SourceBits uint8 // U8x
	PluginParam *PluginParam
}

func (b *BankSourceData) GetPluginType() uint32 {
	return (b.PluginID >> 0) & 0x000F
}

func (b *BankSourceData) GetCompany() uint32 {
	return (b.PluginID >> 4) & 0x03FF
}

func (b *BankSourceData) IsLanguageSpecific() bool {
	return b.SourceBits & 0b00000001 != 0
}

func (b *BankSourceData) IsPrefetch() bool {
	return b.SourceBits & 0b00000010 != 0
}

func (b *BankSourceData) IsNonCachable() bool {
	return b.SourceBits & 0b00001000 != 0
}

func (b *BankSourceData) HasSource() bool {
	return b.SourceBits & 0b10000000 != 0
}

func (b *BankSourceData) HasParam() bool {
	return (b.PluginID & 0x0F) == 2
}

func (b *BankSourceData) Encode() []byte {
	b.validateIntegrity()
	blobSize := 4 + 1 + 4 + 4 + 1 + b.PluginParam.PluginParamSize
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(b.PluginID)
	bw.AppendByte(b.StreamType)
	bw.Append(b.SourceID)
	bw.Append(b.InMemoryMediaSize)
	bw.Append(b.SourceBits)
	bw.AppendBytes(b.PluginParam.Encode())
	return bw.Flush(int(blobSize))
}

func (b *BankSourceData) validateIntegrity() {
	if b.GetPluginType() != 2 {
		assert.AssertTrue(
			b.PluginParam.PluginParamSize <= 0,
			"Plugin type indicate that there's no plugin parameter data.",
		)
		return
	}
	/* This make no sense??? */
	if b.PluginID == 0 {
		assert.AssertTrue(
			b.PluginParam.PluginParamSize <= 0,
			"Plugin ID is zero. There's no plugin parameter data.",
		)
	}
}

type PluginParam struct {
	PluginParamSize uint32 // U32
	PluginParamData []byte
}

func (p *PluginParam) Encode() []byte {
	assert.AssertEqual(
		int(p.PluginParamSize),
		len(p.PluginParamData),
		"Plugin parameter size doesn't equal # of bytes in plugin parameter data",
	)
	blobSize := 4 + len(p.PluginParamData)
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(p.PluginParamSize)
	bw.AppendBytes(p.PluginParamData)
	return bw.Flush(blobSize)
}

type CntrChildren struct {
	NumChild uint32 // u32
	Children []uint32 // NUmChild * sizeof(tid)
}

func NewCntrChildren() *CntrChildren {
	return &CntrChildren{0, []uint32{}}
}

func (c *CntrChildren) Encode() []byte {
	assert.AssertEqual(
		int(c.NumChild),
		len(c.Children),
		"Container children counter does not equal to # of children.",
	)
	blobSize := 4 + 4 * c.NumChild
	bw := reader.NewFixedSizeBlobWriter(uint64(blobSize))
	bw.Append(c.NumChild)
	bw.Append(c.Children)
	return bw.Flush(int(blobSize))
}
