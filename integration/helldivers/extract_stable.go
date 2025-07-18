package helldivers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

const IOCW = os.O_CREATE | os.O_WRONLY

func ExtractSoundBankStable(path string, dest string, dry bool) error {
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

	r.FourCCNoCopyUnsafe()
	r.ReadNoCopyUnsafe(56)
	r.RelSeekUnsafe(int(32 * numTypes))

	// r.RelSeekUnsafe(4)  // unknown
	// r.RelSeekUnsafe(56) // unk4Data
	// r.RelSeekUnsafe(int(32 * numTypes))

	soundBanks := make(map[uint64]*Asset)
	wwiseDependencies := make(map[uint64]*Asset)

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
			soundBanks[h.FileID] = a

			offset := h.DataOffset + 16
			length := offset + uint64(h.DataSize) - 16

			a.Data = r.Buff[offset:length]

			XOR := slices.Clone(a.Data[0x08:0x0C])

			a.Data[0x08] = 0x9A
			a.Data[0x09] = 0x00
			a.Data[0x0A] = 0x00
			a.Data[0x0B] = 0x00

			a.META = make([]byte, 0, SizeOfMetaBufferInit)
			a.META = append(a.META, 'M', 'E', 'T', 'A') // 4
			a.META = append(a.META, 0, 0, 0, 0) // 8
			a.META = append(a.META, byte(IntegrationTypeHelldivers2)) // 13
			a.META, err = binary.Append(a.META, wio.ByteOrder, h.FileID) // 17 
			if err != nil {
				return err
			}
			a.META = append(a.META, XOR...) // 21
			if dep, in := wwiseDependencies[h.FileID]; in && !dry {
				if len(a.META) != 21 {
					panic(fmt.Sprintf("Sound bank %d has not filled up the necessary information on META chunk (chunk tag + chunk size + integration type + file id). The current META size chunk is %d.", h.FileID, len(a.META)))
				}

				a.META = append(a.META, dep.Data...) // 21 + sizeof(wwise dependency)
				size := len(a.META) - wwise.SizeOfChunkHeader

				buffer := make([]byte, 4, 4)
				_, err := binary.Encode(buffer, wio.ByteOrder, uint32(size))
				if err != nil {
					return err
				}

				a.META[4] = buffer[0]
				a.META[5] = buffer[1]
				a.META[6] = buffer[2]
				a.META[7] = buffer[3]

				data := dep.Data[5:]
				path := bytes.ReplaceAll(data, []byte{'\u0000'}, []byte{})
				path = bytes.ReplaceAll(path, []byte{'/'}, []byte{'_'})

				if bytes.Compare(path, []byte{}) == 0 {
					path = fmt.Appendf(path, "%d", a.Header.FileID)
				}

				name := fmt.Sprintf("%s.st_bnk", path)
				w, err := os.OpenFile(filepath.Join(dest, name), IOCW, 0666)
				if err != nil {
					return err
				}
				if _, err := w.Write(a.Data); err != nil {
					return err
				}
				if _, err := w.Write(a.META); err != nil {
					return err
				}
				if err := w.Close(); err != nil {
					return err
				}
			}
		case AssetTypeWwiseDependency:
			if _, in := wwiseDependencies[h.FileID]; in && !dry {
				return fmt.Errorf(
					"Two Wwise dependency headers use the same of ID of %d",
					h.FileID,
				)
			}
			wwiseDependencies[h.FileID] = a

			offset := h.DataOffset
			length := offset + uint64(h.DataSize)

			a.Data = r.Buff[offset:length]
			if bnk, in := soundBanks[h.FileID]; in && !dry {
				data := a.Data[5:]

				path := bytes.ReplaceAll(data, []byte{'\u0000'}, []byte{})
				path = bytes.ReplaceAll(path, []byte{'/'}, []byte{'_'})

				if bytes.Compare(path, []byte{}) == 0 {
					path = fmt.Appendf(path, "%d", a.Header.FileID)
				}
				name := fmt.Sprintf("%s.st_bnk", path)

				if len(bnk.META) != 21 {
					panic(fmt.Sprintf("Sound bank %d has not filled up the necessary information on META chunk (chunk tag + chunk size + integration type + file id). The current META size chunk is %d.", h.FileID, len(bnk.META)))
				}

				bnk.META = append(bnk.META, a.Data...)

				size := len(bnk.META) - wwise.SizeOfChunkHeader

				buffer := make([]byte, 4, 4)
				_, err := binary.Encode(buffer, wio.ByteOrder, uint32(size))
				if err != nil {
					return err
				}
				bnk.META[4] = buffer[0]
				bnk.META[5] = buffer[1]
				bnk.META[6] = buffer[2]
				bnk.META[7] = buffer[3]

				w, err := os.OpenFile(filepath.Join(dest, name), IOCW, 0666)
				if err != nil {
					return err
				}
				if _, err := w.Write(bnk.Data); err != nil {
					return err
				}
				if _, err := w.Write(bnk.META); err != nil {
					return err
				}
				if err := w.Close(); err != nil {
					return err
				}
			}
		}
	}

	for fid, bnk := range soundBanks {
		if _, in := wwiseDependencies[fid]; in {
			continue
		}
		if len(bnk.META) != 21 {
			panic(fmt.Sprintf("Sound bank %d has not filled up the necessary information on META chunk (chunk tag + chunk size + integration type + file id). The current META size chunk is %d.", fid, len(bnk.META)))
		}
		size := len(bnk.META) - wwise.SizeOfChunkHeader
		buffer := make([]byte, 4, 4)
		_, err := binary.Encode(buffer, wio.ByteOrder, uint32(size))
		if err != nil {
			return err
		}
		bnk.META[4] = buffer[0]
		bnk.META[5] = buffer[1]
		bnk.META[6] = buffer[2]
		bnk.META[7] = buffer[3]
		w, err := os.OpenFile(filepath.Join(dest, fmt.Sprintf("%d.st_bnk", fid)), IOCW, 0666)
		if err != nil {
			return err
		}
		if _, err := w.Write(bnk.Data); err != nil {
			return err
		}
		if _, err := w.Write(bnk.META); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			slog.Info(fmt.Sprintf("Failed to closed %d sound bank", fid), "error", err)
		}
	}

	return nil
}

func GetHelldivers2Data() (string, error) {
	data := os.Getenv("HD2DATA")
	if !filepath.IsAbs(data) {
		return "", fmt.Errorf("Helldivers 2 data directory path %s is not an absolute path.", data)
	}
	stat, err := os.Lstat(data)
	if err != nil {
		return "", fmt.Errorf("Invalid Helldivers 2 data directory path %s: %w", data, err)
	}
	if !stat.IsDir() {
		return "", fmt.Errorf("Helldivers 2 data directory path %s points to a file", data)
	}
	return data, nil
}
