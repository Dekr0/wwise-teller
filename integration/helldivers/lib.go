package helldivers

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wio"
)

var NotHelldiversGameArchive error = errors.New(
	"Not a game archive used by Helldivers 2",
)
var NotMETA error = errors.New(
	"Not META data section",
)

const (
	AssetTypeSoundBank       = 6006249203084351385
	AssetTypeWwiseDependency = 12624162998411505776
	AssetTypeWwiseStream     = 5785811756662211598
)

const MagicValue uint32 = 0xF0000011

type IntegrationType uint8

const (
	IntegrationTypeHelldivers2 IntegrationType = 0
)

type Asset struct {
	Header      *AssetHeader
	Data        []byte
	StreamData  []byte
	GPURsrcData []byte
	META        []byte
}

type AssetHeader struct {
	FileID        uint64 `json:"fileID"`
	TypeID        uint64 `json:"typeID"`
	DataOffset    uint64 `json:"dataOffset"`
	StreamOffset  uint64 `json:"streamOffset"`
	GPURsrcOffset uint64 `json:"gPURsrcOffset"`
	UnknownU64A   uint64 `json:"unknownU64A"`
	UnknownU64B   uint64 `json:"unknownU64B"`
	DataSize      uint32 `json:"dataSize"`
	StreamSize    uint32 `json:"streamSize"`
	GPURsrcSize   uint32 `json:"gPURsrcSize"`
	UnknownU32A   uint32 `json:"unknownU32A"`
	UnknownU32B   uint32 `json:"unknownU32B"`
	Idx           uint32 `json:"idx"`
}

type META struct {
	T                     [4]byte
	Size                  uint32
	IntegrationType       uint8
	Unknown               [4]byte
	Unk4Data              [56]byte
	SoundBankAssetHeader  *AssetHeader
	UnusedSoundBank16Data [16]byte
	XOR                   [4]byte
	WwiseDependencyHeader *AssetHeader
	WwiseDependencyData   []byte
}

// TODO: concurrency
func ExtractSoundBank(ctx context.Context, path string, dest string, dry bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r := wio.NewInPlaceReader(data, wio.ByteOrder)

	magic := r.U32Unsafe()
	if magic != 0xF0000011 {
		return NotHelldiversGameArchive
	}

	numTypes := r.U32Unsafe()
	numFiles := r.U32Unsafe()

	unknown := r.FourCCNoCopyUnsafe()
	unk4Data := r.ReadNoCopyUnsafe(56)
	r.RelSeekUnsafe(int(32 * numTypes))

	// r.RelSeekUnsafe(4)  // unknown
	// r.RelSeekUnsafe(56) // unk4Data
	// r.RelSeekUnsafe(int(32 * numTypes))

	soundBanks := make(map[uint64]*Asset)
	wwiseDependencies := make(map[uint64]*Asset)
	var w sync.WaitGroup
	e := make(chan error)
	for range numFiles {
		a := &Asset{
			&AssetHeader{
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U64Unsafe(),
				r.U32Unsafe(),
				r.U32Unsafe(),
				r.U32Unsafe(),
				r.U32Unsafe(),
				r.U32Unsafe(),
				r.U32Unsafe(),
			},
			nil, nil, nil, nil,
		}
		h := a.Header
		switch a.Header.TypeID {
		case AssetTypeSoundBank:
			if _, in := soundBanks[h.FileID]; in {
				return fmt.Errorf(
					"Two sound bank headers use the same of ID of %d", h.FileID,
				)
			}
			offset := h.DataOffset + 16
			// Sound Bank Data (without first 16 bytes and XOR encryption)
			a.Data = r.Buff[offset : offset+uint64(h.DataSize-16)]
			// Backup XOR data
			XOR := slices.Clone(a.Data[0x08:0x0C])
			a.Data[0x08] = 0x8D
			a.Data[0x09] = 0x00
			a.Data[0x0A] = 0x00
			a.Data[0x0B] = 0x00
			// Header: 4 bytes header tag + 4 bytes header size + 1 byte
			// integration type
			// Data: 4 bytes unknown + 56 bytes unknown data + 80 bytes asset
			// information + 16 bytes unused + 4 bytes XOR
			a.META = make([]byte, 0, 105)
			a.META = append(a.META, 'M', 'E', 'T', 'A')
			a.META = append(a.META, 0, 0, 0, 0)
			a.META = append(a.META, byte(IntegrationTypeHelldivers2))
			a.META = append(a.META, unknown...)
			a.META = append(a.META, unk4Data...)
			a.META, err = binary.Append(a.META, wio.ByteOrder, h)
			if err != nil {
				return err
			}
			a.META = append(a.META, r.Buff[h.DataOffset:h.DataOffset+16]...)
			a.META = append(a.META, XOR...)
			if dep, in := wwiseDependencies[h.FileID]; in && !dry {
				w.Add(1)
				go func() {
					defer w.Done()
					a.META = append(a.META, dep.META...)
					a.META = append(a.META, dep.Data...)
					metaSize := uint32(len(a.META) - 4 - 4)
					encodedMetaSize, err := binary.Append(
						[]byte{}, wio.ByteOrder, metaSize,
					)
					if err != nil {
						e <- err
					}
					a.META[4] = encodedMetaSize[0]
					a.META[5] = encodedMetaSize[1]
					a.META[6] = encodedMetaSize[2]
					a.META[7] = encodedMetaSize[3]

					st_path := bytes.ReplaceAll(
						bytes.ReplaceAll(dep.Data[5:], []byte{'\u0000'}, []byte{}),
						[]byte{'/'}, []byte{'_'},
					)
					if bytes.Compare(st_path, []byte{}) == 0 {
						e <- fmt.Errorf("Sound bank %d name is empty", h.FileID)
					}
					path := fmt.Sprintf("%s.st_bnk", st_path)
					wf, err := os.OpenFile(
						filepath.Join(dest, path), os.O_WRONLY, 0666,
					)
					if err != nil {
						e <- err
					}
					if _, err := wf.Write(a.Data); err != nil {
						e <- err
					}
					if _, err := wf.Write(a.META); err != nil {
						e <- err
					}
					if err := wf.Close(); err != nil {
						e <- err
					}
				}()
			} else {
				soundBanks[h.FileID] = a
			}
		case AssetTypeWwiseDependency:
			if _, in := wwiseDependencies[h.FileID]; in && !dry {
				return fmt.Errorf(
					"Two Wwise dependency headers use the same of ID of %d",
					h.FileID,
				)
			}
			a.META = make([]byte, 0, 80)
			a.META, err = binary.Append(a.META, wio.ByteOrder, h)
			a.Data = r.Buff[h.DataOffset : h.DataOffset+uint64(h.DataSize)]
			if err != nil {
				return err
			}
			if bnk, in := soundBanks[h.FileID]; in && !dry {
				w.Add(1)
				go func() {
					defer w.Done()
					st_path := bytes.ReplaceAll(
						bytes.ReplaceAll(a.Data[5:], []byte{'\u0000'}, []byte{}),
						[]byte{'/'}, []byte{'_'},
					)
					if bytes.Compare(st_path, []byte{}) == 0 {
						e <- fmt.Errorf("Sound bank %d name is empty", h.FileID)
					}
					path := fmt.Sprintf("%s.st_bnk", st_path)
					bnk.META = append(bnk.META, a.META...)
					bnk.META = append(bnk.META, a.Data...)
					metaSize := uint32(len(bnk.META) - 4 - 4)
					encodedMetaSize, err := binary.Append(
						[]byte{}, wio.ByteOrder, metaSize,
					)
					if err != nil {
						e <- err
					}
					bnk.META[4] = encodedMetaSize[0]
					bnk.META[5] = encodedMetaSize[1]
					bnk.META[6] = encodedMetaSize[2]
					bnk.META[7] = encodedMetaSize[3]

					wf, err := os.OpenFile(filepath.Join(dest, path), IOCW, 0666)
					if err != nil {
						e <- err
					}
					if _, err := wf.Write(bnk.Data); err != nil {
						e <- err
					}
					if _, err := wf.Write(bnk.META); err != nil {
						e <- err
					}
					if err := wf.Close(); err != nil {
						e <- err
					}
				}()
			} else {
				wwiseDependencies[h.FileID] = a
			}
		}
	}
	w.Wait()
	for len(e) > 0 {
		select {
		case err := <-e:
			return err
		}
	}
	return nil
}

func ParseMETA(_ context.Context, r *wio.InPlaceReader) (*META, error) {
	itype, err := r.U8()
	if err != nil {
		return nil, err
	}
	if itype != uint8(IntegrationTypeHelldivers2) {
		return nil, NotHelldiversGameArchive
	}

	unknown, err := r.FourCC()
	if err != nil {
		return nil, err
	}

	unk4Data, err := r.ReadNoCopy(56)
	if err != nil {
		return nil, err
	}

	data, err := r.ReadNoCopy(80)
	if err != nil {
		return nil, err
	}

	var soundBnk AssetHeader
	_, err = binary.Decode(data, wio.ByteOrder, &soundBnk)
	if err != nil {
		return nil, err
	}

	unusedSoundBnk16Data, err := r.ReadNoCopy(16)
	if err != nil {
		return nil, err
	}

	xor, err := r.ReadNoCopy(4)
	if err != nil {
		return nil, err
	}

	data, err = r.ReadNoCopy(80)
	if err != nil {
		return nil, err
	}

	var wwiseDep AssetHeader
	_, err = binary.Decode(data, wio.ByteOrder, &wwiseDep)
	if err != nil {
		return nil, err
	}

	wwiseDepData, err := r.ReadNoCopy(uint(wwiseDep.DataSize))
	if err != nil {
		return nil, err
	}

	if r.Len() != 0 {
		return nil, fmt.Errorf(
			"Data size of Wwise dependency is incorrect. There are more bytes" +
				" to consume.",
		)
	}

	return &META{
		IntegrationType:       itype,
		Unknown:               [4]byte(unknown),
		Unk4Data:              [56]byte(unk4Data),
		SoundBankAssetHeader:  &soundBnk,
		UnusedSoundBank16Data: [16]byte(unusedSoundBnk16Data),
		XOR:                   [4]byte(xor),
		WwiseDependencyHeader: &wwiseDep,
		WwiseDependencyData:   wwiseDepData,
	}, nil
}

func GenHelldiversPatch(
	ctx context.Context,
	bnkData []byte,
	metaCu []byte,
	path string,
) error {
	meta, err := ParseMETA(ctx, wio.NewInPlaceReader(metaCu, wio.ByteOrder))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(
		filepath.Join(path, "9ba626afa44a3aa3.patch_0"),
		os.O_CREATE|os.O_WRONLY,
		0666,
	)
	defer f.Close()
	if err != nil {
		return err
	}

	w := wio.NewBinaryWriteHelper(f)
	if err := writeGameArchiveHeader(ctx, w, meta); err != nil {
		return err
	}
	if err := writeGameArchiveTypeHeader(ctx, w); err != nil {
		return err
	}
	bnkData, err = writeAssetHeaders(ctx, w, bnkData, meta)
	if err != nil {
		return err
	}
	pad := make([]byte, 8, 8)
	if err := w.Bytes(pad); err != nil {
		return err
	}
	if err := writeSoundBank(ctx, w, bnkData, meta); err != nil {
		return err
	}
	if err := w.Bytes(meta.WwiseDependencyData); err != nil {
		return err
	}

	return nil
}

func writeGameArchiveHeader(
	_ context.Context,
	w *wio.BinaryWriteHelper,
	meta *META,
) error {
	if err := w.U32(MagicValue); err != nil {
		return err
	}

	// Number of file types = 2 (Sound bank and dependency)
	if err := w.U32(2); err != nil {
		return err
	}

	// Number of files = 2 x number of sound banks
	if err := w.U32(2); err != nil {
		return err
	}

	if err := w.Bytes(meta.Unknown[:]); err != nil {
		return err
	}

	// if err := w.Bytes(meta.Unk4Data[:]); err != nil {
	// 	return err
	// }
	
	// Require (Reverse engineering black magic)
	Unk4Data := []byte{206,9,245,244,0,0,0,0,12,114,159,158,136,114,184,189,0,160,107,2,0,0,0,0,0,121,81,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	if err := w.Bytes(Unk4Data); err != nil {
		return err
	}

	return nil
}

func writeGameArchiveTypeHeader(_ context.Context, w *wio.BinaryWriteHelper) error {
	// Sound bank type header
	if err := w.U64(0); err != nil {
		return err
	}
	if err := w.U64(AssetTypeSoundBank); err != nil {
		return err
	}
	// Number of sound banks
	if err := w.U64(1); err != nil {
		return err
	}
	// Alignment? 8 x 3 + 2 x 4 = 24 + 8 = 32
	if err := w.U32(16); err != nil {
		return err
	}
	if err := w.U32(64); err != nil {
		return err
	}

	// Wwise dependency type header
	if err := w.U64(0); err != nil {
		return err
	}
	if err := w.U64(AssetTypeWwiseDependency); err != nil {
		return err
	}
	// Number of wwise dependencies
	if err := w.U64(1); err != nil {
		return err
	}
	// Alignment? 8 x 3 + 2 x 4 = 24 + 8 = 32
	if err := w.U32(16); err != nil {
		return err
	}
	if err := w.U32(64); err != nil {
		return err
	}

	return nil
}

func writeAssetHeaders(
	_ context.Context,
	w *wio.BinaryWriteHelper,
	bnkData []byte,
	meta *META,
) ([]byte, error) {
	dataOffset := uint32(160 + w.Tell() + 8)

	// File ID and Type ID should remain same
	meta.SoundBankAssetHeader.Idx = 0
	meta.SoundBankAssetHeader.DataSize = uint32(len(bnkData)) + 16
 	bnkData = utils.Pad16ByteAlign(bnkData)
	meta.SoundBankAssetHeader.StreamSize = 0
	meta.SoundBankAssetHeader.GPURsrcSize = 0
	meta.SoundBankAssetHeader.DataOffset = uint64(dataOffset)
	dataOffset += 16 + uint32(len(bnkData))
	meta.SoundBankAssetHeader.StreamOffset = 0
	meta.SoundBankAssetHeader.GPURsrcOffset = 0
	data, err := binary.Append(nil, wio.ByteOrder, meta.SoundBankAssetHeader)
	if err != nil {
		return nil, err
	}
	if err := w.Bytes(data); err != nil {
		return nil, err
	}

	// File ID, Type ID, File Size should remain same
	meta.WwiseDependencyHeader.Idx = 1
	meta.WwiseDependencyHeader.DataSize = uint32(len(meta.WwiseDependencyData))
	meta.WwiseDependencyData = utils.Pad16ByteAlign(meta.WwiseDependencyData)
	meta.WwiseDependencyHeader.DataOffset = uint64(dataOffset)
	meta.WwiseDependencyHeader.StreamSize = 0
	meta.WwiseDependencyHeader.GPURsrcSize = 0
	dataOffset += uint32(len(meta.WwiseDependencyData))
	meta.WwiseDependencyHeader.StreamOffset = 0
	meta.WwiseDependencyHeader.GPURsrcOffset = 0
	data, err = binary.Append(nil, wio.ByteOrder, meta.WwiseDependencyHeader)
	if err != nil {
		return nil, err
	}
	if err := w.Bytes(data); err != nil {
		return nil, err
	}

	return bnkData, nil
}

func writeSoundBank(
	_ context.Context,
	w *wio.BinaryWriteHelper,
	bnkData []byte,
	meta *META,
) error {
	// Require (Reverse engineering black magic)
	if err := w.Bytes([]byte{216,47,118,120}); err != nil {
		return err
	}
	if err := w.U32(uint32(len(bnkData))); err != nil {
		return err
	}
	if err := w.U64(meta.SoundBankAssetHeader.FileID); err != nil {
		return err
	}
	bnkData[0x08] = meta.XOR[0]
	bnkData[0x09] = meta.XOR[1]
	bnkData[0x0A] = meta.XOR[2]
	bnkData[0x0B] = meta.XOR[3]
	if err := w.Bytes(bnkData); err != nil {
		return err
	}
	return nil
}
