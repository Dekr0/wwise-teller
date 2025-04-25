package wwise

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
)

const maxEncodeRoutine = 8

// # of hierarchy object (uint32)
const sizeOfHIRCHeader = 4

type HircType uint8

const (
	HircTypeSound HircType = 0x02
	HircRanSeqCntr HircType = 0x05
	HircSwitchCntr HircType = 0x06
	HircTypeActorMixer HircType = 0x07
	HircTypeLayerCntr HircType = 0x09
)

var KnownHircType []HircType = []HircType{
	HircTypeSound,
	HircRanSeqCntr,
	HircSwitchCntr,
	HircTypeActorMixer,
	HircTypeLayerCntr,
}

var HircTypeName []string = []string{
	"",
	"State",
	"Sound",
	"Action",
	"Event",
	"Random / Sequence Container",
	"Switch Container",
	"Actor Mixer",
	"Bus",
	"Layer Container",
	"Music Segment",
	"Music Track",
	"Music Switch Container",
	"Music Random / Sequence Container",
	"Attenuation",
	"Dialogue Event",
    "FX Share Set",
    "FX Custom",
    "Aux Bus",
    "LFO Modulator",
    "Envelope Modulator",
    "Audio Device",
    "Time Modulator",
}

type HIRC struct {
 	// Used for memory allocation upfront during encoding 
	// Future notes: Should I track individual byte change whenever a hierarchy 
	// obj changes?
	oldChunkSize uint32

	// Currently, I don't know the algorithm of how Wwise encode its hierarchy 
	// tree. 
	// It's probably some sort of modified DFS since there are lot of places 
	// where child nodes come first, and then the parent node of those child 
	// nodes come right after it. 
	// So far now, I will book keep the linear order of hierarchy tree as I 
	// parse them linearly through.
	HircObjs    []HircObj
	
	// Map for different types of hierarchy objects. Each object is a pointer 
	// to a specific hierarchy object, which is also in `HircObjs`.
	ActorMixers map[uint32]*ActorMixer
	LayerCntrs  map[uint32]*LayerCntr
	SwitchCntrs map[uint32]*SwitchCntr
	RanSeqCntrs map[uint32]*RanSeqCntr
	Sounds      map[uint32]*Sound
}

func NewHIRC(size uint32, numHircItem uint32) *HIRC {
	return &HIRC{
		oldChunkSize: size,
		HircObjs: make([]HircObj, numHircItem),
		ActorMixers: make(map[uint32]*ActorMixer),
		LayerCntrs: make(map[uint32]*LayerCntr),
		SwitchCntrs: make(map[uint32]*SwitchCntr),
		RanSeqCntrs: make(map[uint32]*RanSeqCntr),
		Sounds: make(map[uint32]*Sound),
	}
}

func (h *HIRC) encode(ctx context.Context) ([]byte, error) {
	type result struct {
		i int
		b []byte
	}

	// No initialization since I want it to crash and catch encoding bugs
	results := make([][]byte, len(h.HircObjs))

	// sync signal
	c := make(chan *result, maxEncodeRoutine)

	// limit # of go routines running at the same time
	sem := make(chan struct{} , maxEncodeRoutine)

	done := 0
	i := 0
	for done < len(h.HircObjs) {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case r := <- c:
			results[r.i] = r.b
			done += 1
		case sem <- struct{}{}:
			if i < len(h.HircObjs) {
				j := i
				go func() {
					c <- &result{j, h.HircObjs[j].Encode()}
					<- sem
				}()	
				i += 1
			}
		default:
			if i < len(h.HircObjs) {
				results[i] = h.HircObjs[i].Encode()
				done += 1
				i += 1
			}
		}
	}
	
	return bytes.Join(results, []byte{}), nil
}

func (h *HIRC) Encode(ctx context.Context) ([]byte, error) {
	b, err := h.encode(ctx)
	if err != nil {
		return nil, err
	}

	dataSize := uint32(sizeOfHIRCHeader + len(b))
	size := chunkHeaderSize + dataSize
	w := wio.NewWriter(uint64(chunkHeaderSize + dataSize))
	w.AppendBytes([]byte{'H', 'I', 'R', 'C'})
	w.Append(dataSize)
	w.Append(uint32(len(h.HircObjs)))
	w.AppendBytes(b)
	return w.BytesAssert(int(size)), nil 
}

type HircObj interface {
	Encode() []byte
	HircID() (uint32, error)
	HircType() HircType 
}

const sizeOfHircObjHeader = 1 + 4

type HircObjHeader struct {
	Type HircType // U8x
	Size uint32 // U32
}

type ActorMixer struct {
	Id uint32
	BaseParam *BaseParameter
	Container *Container
}

func (a *ActorMixer) Encode() []byte {
	dataSize := a.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeActorMixer))
	w.Append(dataSize)
	w.Append(a.Id)
	w.AppendBytes(a.BaseParam.Encode())
	w.AppendBytes(a.Container.Encode())
	return w.BytesAssert(int(size))
}

func (a *ActorMixer) DataSize() uint32 {
	return uint32(4 + a.BaseParam.Size() + a.Container.Size())
}

func (a *ActorMixer) HircID() (uint32, error) {
	return a.Id, nil
}

func (a *ActorMixer) HircType() HircType {
	return HircTypeActorMixer 
}

type LayerCntr struct {
	Id uint32
	BaseParam *BaseParameter
	Container *Container

	// NumLayers uint32 // u32
	
	Layers []*Layer
	IsContinuousValidation uint8 // U8x
}

func (l *LayerCntr) Encode() []byte {
	dataSize := l.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeLayerCntr))
	w.Append(dataSize)
	w.Append(l.Id)
	w.AppendBytes(l.BaseParam.Encode())
	w.AppendBytes(l.Container.Encode())
	w.Append(uint32(len(l.Layers)))
	for _, i := range l.Layers {
		w.AppendBytes(i.Encode())
	}
	w.AppendByte(l.IsContinuousValidation)
	return w.BytesAssert(int(size))
}

func (l *LayerCntr) DataSize() uint32 {
	size := 4 + l.BaseParam.Size() + l.Container.Size() + 4
	for _, i := range l.Layers {
		size += i.Size()
	}
	return size + 1
}


func (l *LayerCntr) HircID() (uint32, error) {
	return l.Id, nil
}

func (l *LayerCntr) HircType() HircType {
	return HircTypeLayerCntr
}

type RanSeqCntr struct {
	Id uint32
	BaseParam *BaseParameter
	Container *Container
	PlayListSetting *PlayListSetting

	// NumPlayListItem u16

	PlayListItems []*PlayListItem 
}

func (r *RanSeqCntr) Encode() []byte {
	dataSize := r.DataSize()
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircRanSeqCntr))
	w.Append(dataSize)
	w.Append(r.Id)
	w.AppendBytes(r.BaseParam.Encode())
	w.Append(r.PlayListSetting)
	w.AppendBytes(r.Container.Encode())
	w.Append(uint16(len(r.PlayListItems)))
	for _, i := range r.PlayListItems {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (r *RanSeqCntr) DataSize() uint32 {
	return uint32(4 + r.BaseParam.Size() + r.Container.Size() + sizeOfPlayListSetting + 2 + uint32(len(r.PlayListItems)) * sizeOfPlayListItem)
}

func (r *RanSeqCntr) HircID() (uint32, error) {
	return r.Id, nil
}

func (r *RanSeqCntr) HircType() HircType {
	return HircRanSeqCntr
}

type SwitchCntr struct {
	Id uint32
	BaseParam *BaseParameter
	GroupType uint8 // U8x
	GroupID uint32 // tid
	DefaultSwitch uint32 // tid
	IsContinuousValidation uint8 // U8x
	Container *Container

	// NumSwitchGroups uint32 // u32

	SwitchGroups []*SwitchGroupItem

	// NumSwitchParams uint32 // u32

	SwitchParams []*SwitchParam
}

func (s *SwitchCntr) Encode() []byte {
	baseParamData := s.BaseParam.Encode()
	cntrData := s.Container.Encode()
	switchGroupDataSize := uint32(4)
	for _, i := range s.SwitchGroups {
		switchGroupDataSize += i.Size()
	}
	dataSize := 4 + uint32(len(baseParamData)) + 1 + 4 + 4 + 1 + 
				uint32(len(cntrData)) + switchGroupDataSize + 4 + 
				uint32(len(s.SwitchParams)) * sizeOfSwitchParam
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircSwitchCntr))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(baseParamData)
	w.AppendByte(s.GroupType)
	w.Append(s.GroupID)
	w.Append(s.DefaultSwitch)
	w.AppendByte(s.IsContinuousValidation)
	w.AppendBytes(cntrData)
	w.Append(uint32(len(s.SwitchGroups)))
	for _, i := range s.SwitchGroups {
		w.AppendBytes(i.Encode())
	}
	w.Append(uint32(len(s.SwitchParams)))
	for _, i := range s.SwitchParams {
		w.Append(i)
	}
	return w.BytesAssert(int(size))
}

func (s *SwitchCntr) HircID() (uint32, error) {
	return s.Id, nil
}

func (s *SwitchCntr) HircType() HircType {
	return HircSwitchCntr
}

type Sound struct {
	Id uint32
	BankSourceData *BankSourceData
	BaseParam *BaseParameter
}

func (s *Sound) Encode() []byte {
	b := s.BankSourceData.Encode()
	b = append(b, s.BaseParam.Encode()...)
	dataSize := uint32(4 + len(b))
	size := sizeOfHircObjHeader + dataSize
	w := wio.NewWriter(uint64(size))
	w.AppendByte(uint8(HircTypeSound))
	w.Append(dataSize)
	w.Append(s.Id)
	w.AppendBytes(b)
	return w.BytesAssert(int(size))
}

func (s *Sound) HircID() (uint32, error) {
	return s.Id, nil
}

func (s *Sound) HircType() HircType {
	return HircTypeSound
}

type Unknown struct {
	Header *HircObjHeader
	b []byte
}

func NewUnknown(t HircType, s uint32, b []byte) *Unknown {
	return &Unknown{
		Header: &HircObjHeader{Type: t, Size: s},
		b: b,
	}
}

func (u *Unknown) Encode() []byte {
	assert.Equal(
		u.Header.Size,
		uint32(len(u.b)),
		"Header size does not equal to actual data size",
	)

	bw := wio.NewWriter(uint64(sizeOfHircObjHeader + len(u.b)))
	
	/* Header */
	bw.Append(u.Header)
	bw.AppendBytes(u.b)

	return bw.Bytes() 
}

func (u *Unknown) HircID() (uint32, error) {
	return 0, fmt.Errorf("Hierarchy object type %d has yet implement GetHircID.", u.Header.Type)
}

func (u *Unknown) HircType() HircType {
	return u.Header.Type
}
