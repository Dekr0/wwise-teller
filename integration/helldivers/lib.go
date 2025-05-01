package helldivers

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"context"
	"encoding/binary"
	"path/filepath"
	"sync"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

var NotHelldiversGameArchive error = errors.New(
	"Not a game archive used by Helldivers 2",
)

const (
	AssetTypeSoundBank       = 6006249203084351385
	AssetTypeWwiseDependency = 12624162998411505776
	AssetTypeWwiseStream     = 5785811756662211598
)

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
	UnknownU64B   uint64 `json:"unknownU64B"`
	UnknownU64A   uint64 `json:"unknownU64A"`
	DataSize      uint32 `json:"dataSize"`
	StreamSize    uint32 `json:"streamSize"`
	GPURsrcSize   uint32 `json:"gPURsrcSize"`
	UnknownU32A   uint32 `json:"unknownU32A"`
	UnknownU32B   uint32 `json:"unknownU32B"`
	Idx           uint32 `json:"idx"`
}

// TODO: concurrency
func ExtractSoundBank(ctx context.Context, path string, dest string) error {
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
	r.RelSeekUnsafe(4)  // unknown
	r.RelSeekUnsafe(56) // unk4Data
	r.RelSeekUnsafe(int(32 * numTypes))

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
			a.Data = r.Buff[offset:offset + uint64(h.DataSize - 16)]
			// Backup XOR data
			XOR := slices.Clone(a.Data[0x08:0x0B])
			a.Data[0x08] = 0x8D
			a.Data[0x09] = 0x00
			a.Data[0x0A] = 0x00
			a.Data[0x0B] = 0x00
			// Header: 4 bytes header tag + 4 bytes header size + 1 byte 
			// integration type
			// Data: 80 bytes asset information + 16 bytes unused + 4 bytes XOR
			a.META = make([]byte, 0, 105)
			a.META = append(a.META, 'M', 'E', 'T', 'A')
			a.META = append(a.META, 0, 0, 0, 0)
			a.META = append(a.META, byte(IntegrationTypeHelldivers2))
			a.META, err = binary.Append(a.META, wio.ByteOrder, h)
			if err != nil {
				return err
			}
			a.META = append(a.META, r.Buff[h.DataOffset:h.DataOffset+16]...)
			a.META = append(a.META, XOR...)
			if dep, in := wwiseDependencies[h.FileID]; in {
				w.Add(1)
				go func() {
					defer w.Done()
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
			if _, in := wwiseDependencies[h.FileID]; in {
				return fmt.Errorf(
					"Two Wwise dependency headers use the same of ID of %d",
					h.FileID,
				)
			}
			a.Data = r.Buff[h.DataOffset:h.DataOffset + uint64(h.DataSize)]
			if bnk, in := soundBanks[h.FileID]; in {
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

					wf, err := os.OpenFile(
						filepath.Join(dest, path), os.O_CREATE | os.O_WRONLY, 0666,
					)
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
		case err := <- e:
			return err
		}
	}
	return nil
}

func GenHelldiversPatch(ctx context.Context, bank *wwise.Bank) error {
	return nil
}
