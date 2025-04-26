package helldivers

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
)

var NotHelldiversGameArchive error = errors.New(
	"Not a game archive used by Helldivers 2",
)

const (
	AssetTypeSoundBank       = 6006249203084351385
	AssetTypeWwiseDependency = 12624162998411505776
	AssetTypeWwiseStream     = 5785811756662211598
)

type Asset struct {
	Header      *AssetHeader
	Data        []byte
	StreamData  []byte
	GPURsrcData []byte
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
func ExtractSoundBank(ctx context.Context, path string) error {
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

	// var w sync.WaitGroup

	// e := make(chan error)

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
			nil, nil, nil,
		}
		h := a.Header
		switch a.Header.TypeID {
		case AssetTypeSoundBank:
			if _, in := soundBanks[h.FileID]; in {
				return fmt.Errorf(
					"Two sound bank headers use the same of ID of %d", h.FileID,
				)
			}
			// Sound Bank Data (without first 16 bytes and XOR encryption)
			a.Data = r.Buff[h.DataOffset + 16:h.DataOffset + uint64(h.DataSize) - 16]

			XOR := slices.Clone(a.Data[0x08:0x0B])

			a.Data[0x08] = 0x8D
			a.Data[0x09] = 0x00
			a.Data[0x0A] = 0x00
			a.Data[0x0B] = 0x00

			// Backing up data for resurrection
			// Asset header
			// Padding
			a.Data = append(a.Data, make([]byte, 32)...)
			
			a.Data, err = binary.Append(a.Data, wio.ByteOrder, h)
			if err != nil {
				return err
			}
			// first 16 bytes unused data and XOR encryption
			a.Data = append(a.Data, r.Buff[h.DataOffset:h.DataOffset+16]...)
			a.Data = append(a.Data, XOR...)

			// Sound bank come first, then Wwise dependency
			if dep, in := wwiseDependencies[h.FileID]; in {
				a.Data = append(a.Data, dep.Data...)

				st_path := bytes.ReplaceAll(
					bytes.ReplaceAll(dep.Data[5:], []byte{'\u0000'}, []byte{}),
					[]byte{'/'}, []byte{'_'},
				)

				if bytes.Compare(st_path, []byte{}) == 0 {
					return fmt.Errorf("Sound bank %d name is empty", h.FileID)
				}

				path := fmt.Sprintf("%s.st_bnk", st_path)

				if err := os.WriteFile(path, a.Data, 0666); err != nil {
					return err
				}

				// w.Add(1)
				// go func() {
				// 	defer w.Wait()
				// 	// Append encoded Wwise dependency data
				// 	a.Data = append(a.Data, dep.Data...)

				// 	st_path := bytes.ReplaceAll(
				// 		bytes.ReplaceAll(dep.Data[4:], []byte{'\u0000'}, []byte{}),
				// 		[]byte{'/'}, []byte{'_'},
				// 	)

				// 	if bytes.Compare(st_path, []byte{}) == 0 {
				// 		e <- fmt.Errorf("Sound bank %d name is empty", h.FileID)
				// 		return
				// 	}

				// 	path := fmt.Sprintf("%s.st_bnk", st_path)

				// 	if err := os.WriteFile(path, a.Data, 0666); err != nil {
				// 		e <- err
				// 	}
				// }()
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

			// Wwise dependency comes first, then sound bank
			if bnk, in := soundBanks[h.FileID]; in {
				// Compute name first before position get lost after attach
				// sound bank data
				st_path := bytes.ReplaceAll(
					bytes.ReplaceAll(a.Data[5:], []byte{'\u0000'}, []byte{}),
					[]byte{'/'}, []byte{'_'},
				)

				if bytes.Compare(st_path, []byte{}) == 0 {
					return fmt.Errorf("Sound bank %d name is empty", h.FileID)
				}

				path := fmt.Sprintf("%s.st_bnk", st_path)

				// Attach sound bank in the front
				a.Data = append(bnk.Data, a.Data...)

				if err := os.WriteFile(path, a.Data, 0666); err != nil {
					return err
				}

				// w.Add(1)
				// go func() {
				// 	defer w.Wait()

				// 	// Compute name first before position get lost after attach
				// 	// sound bank data
				// 	st_path := bytes.ReplaceAll(
				// 		bytes.ReplaceAll(a.Data[4:], []byte{'\u0000'}, []byte{}),
				// 		[]byte{'/'}, []byte{'_'},
				// 	)

				// 	if bytes.Compare(st_path, []byte{}) == 0 {
				// 		e <- fmt.Errorf("Sound bank %d name is empty", h.FileID)
				// 		return
				// 	}

				// 	path := fmt.Sprintf("%s.st_bnk", st_path)

				// 	// Attach sound bank in the front
				// 	a.Data = append(bnk.Data, a.Data...)

				// 	if err := os.WriteFile(path, a.Data, 0666); err != nil {
				// 		e <- err
				// 	}
				// }()
			} else {
				wwiseDependencies[h.FileID] = a
			}
		}
	}

	return nil
}
