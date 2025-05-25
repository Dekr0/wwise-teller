package parser

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

var NoBKHD error = errors.New("Sound bank is missing BKHD section")
var NoDIDX error = errors.New("Sound bank is missing DIDX section")
var NoDATA error = errors.New("Sound bank is missing DATA section")
var NoHIRC error = errors.New("Sound bank is missing HIRC section")

var (
	CustomVersions []uint32 = []uint32{122, 126, 129, 135, 136}
	Versions       []uint32 = []uint32{
		//  --,  // 0x-- Wwise 2016.1~3
		//  14,  // 0x0E Wwise 2007.1/2?
		26,  // 0x1A Wwise 2007.3?
		29,  // 0x1D Wwise 2007.4?
		34,  // 0x22 Wwise 2008.1?
		35,  // 0x23 Wwise 2008.2?
		36,  // 0x24 Wwise 2008.3?
		38,  // 0x26 Wwise 2008.4
		44,  // 0x2C Wwise 2009.1?
		45,  // 0x2D Wwise 2009.2?
		46,  // 0x2E Wwise 2009.3
		48,  // 0x30 Wwise 2010.1
		52,  // 0x34 Wwise 2010.2
		53,  // 0x35 Wwise 2010.3
		56,  // 0x38 Wwise 2011.1
		62,  // 0x3E Wwise 2011.2
		65,  // 0x41 Wwise 2011.3?
		70,  // 0x46 Wwise 2012.1?
		72,  // 0x48 Wwise 2012.2
		88,  // 0x58 Wwise 2013.1/2
		89,  // 0x59 Wwise 2013.2-B?
		112, // 0x70 Wwise 2014.1
		113, // 0x71 Wwise 2015.1
		118, // 0x76 Wwise 2016.1
		120, // 0x78 Wwise 2016.2
		122, // 0x7A Wwise 2017.1-B?
		125, // 0x7D Wwise 2017.1
		126, // 0x7E Wwise 2017.1-B?
		128, // 0x80 Wwise 2017.2
		129, // 0x81 Wwise 2017.2-B?
		132, // 0x84 Wwise 2018.1
		134, // 0x86 Wwise 2019.1
		135, // 0x87 Wwise 2019.1-B?
		135, // 0x87 Wwise 2019.2
		136, // 0x88 Wwise 2019.2-B?
		140, // 0x8c Wwise 2021.1
		141, // 0x8d Wwise 2021.1-B?
		144, // 0x90 Wwise 2022.1-B
		145, // 0x91 Wwise 2022.1
		150, // 0x96 Wwise 2023.1
		152, // 0x98 Wwise 2024.1-B
	}
)

type DecodeResult struct {
	c wwise.Chunk
	e error
}

func ParseBank(path string, ctx context.Context) (*wwise.Bank, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bankReader := wio.NewReader(f, binary.LittleEndian)

	version, err := CheckHeader(bankReader)
	if err != nil {
		return nil, err
	}

	if version != 141 {
		return nil, errors.New("Wwise teller currently only targets version 141.")
	}

	c := make(chan *DecodeResult)
	hasBKHD := false
	hasDIDX := false
	hasHIRC := false
	pending := 0
	I := uint8(0)

	/*
		Parallel parsing

		Assumption:
		Ideally, a sound bank is made out of different sections / chunks. Each
		section / chunk usually follow this pattern:
		1 - section / chunk tag
		2 - section / chunk size
		3 - data

		Implementation:
		1 - Read off a chunk tag and its associated chunk size.
		2 - Read all bytes in this chunk from the file.
		3 - Create a new reader that operates on this byte slice, and use it for
		parsing. The parsing happen in a go routine

	*/
	bnk := wwise.NewBank()
	err = nil
	var tag []byte
	var size uint32
	for err == nil {
		tag, err = bankReader.FourCC()
		if err != nil {
			continue
		}

		if bytes.Compare(tag, []byte{'B', 'K', 'H', 'D'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var reader *wio.Reader
			reader, err = bankReader.NewBufferReader(uint64(size))
			if err != nil {
				continue
			}
			slog.Debug("Start parsing BKHD section...", "size", size)
			go BKHDRoutine(ctx, reader, c, I, tag, size)
			I += 1
			pending += 1
			hasBKHD = true
		} else if bytes.Compare(tag, []byte{'D', 'A', 'T', 'A'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewDATA(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read DATA section", "size", size)
		} else if bytes.Compare(tag, []byte{'D', 'I', 'D', 'X'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var reader *wio.Reader
			reader, err = bankReader.NewBufferReader(uint64(size))
			if err != nil {
				continue
			}
			slog.Debug("Start parsing DIDX section...", "chunkSize", size)
			go DIDXRoutine(ctx, reader, c, I, tag, size)
			I += 1
			pending += 1
			hasDIDX = true
		} else if bytes.Compare(tag, []byte{'E', 'N', 'V', 'S'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewENVS(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read ENVS section", "size", size)
		} else if bytes.Compare(tag, []byte{'F', 'X', 'P', 'R'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewFXPR(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read FXPR section", "size", size)
		} else if bytes.Compare(tag, []byte{'H', 'I', 'R', 'C'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var reader *wio.Reader
			reader, err = bankReader.NewBufferReader(uint64(size))
			if err != nil {
				continue
			}
			slog.Debug("Start parsing HIRC section...", "chunkSize", size)
			go HIRCRoutine(ctx, reader, c, I, tag, size)
			I += 1
			pending += 1
			hasHIRC = true
		} else if bytes.Compare(tag, []byte{'I', 'N', 'I', 'T'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewINIT(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read INIT section", "size", size)
		} else if bytes.Compare(tag, []byte{'P', 'L', 'A', 'T'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewPLAT(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read PLAT section", "size", size)
		} else if bytes.Compare(tag, []byte{'S', 'T', 'I', 'D'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewSTID(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read STID section", "size", size)
		} else if bytes.Compare(tag, []byte{'S', 'T', 'M', 'G'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewSTMG(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read STMG section", "size", size)
		} else if bytes.Compare(tag, []byte{'M', 'E', 'T', 'A'}) == 0 {
			size, err = bankReader.U32()
			if err != nil {
				continue
			}
			var blob []byte
			blob, err = bankReader.ReadN(uint64(size), 0)
			if err != nil {
				continue
			}
			if err := bnk.AddChunk(wwise.NewMETA(I, tag, blob)); err != nil {
				return nil, err
			}
			I += 1
			slog.Debug("Read META section", "size", size)
		} else {
			tagHex := uint32(0)
			binary.Decode(tag, wio.ByteOrder, tagHex)
			slog.Info("Unknown Chunk Tag", "Tag in hex", tagHex)
		}
	}

	if !hasBKHD {
		return nil, NoBKHD
	}
	if !hasDIDX {
		slog.Warn("Sound bank is missing DIDX section. This might not be a "+
			"problem as long as HIRC section exists.", "soundbank", path)
	}
	if !hasHIRC {
		return nil, NoHIRC
	}

	for pending > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res := <-c:
			if res.e != nil {
				return nil, err
			}

			if err := bnk.AddChunk(res.c); err != nil {
				return nil, err
			}
			pending -= 1
			slog.Debug(fmt.Sprintf("Parsed %s section", res.c.Tag()))
		}
	}

	if bnk.BKHD() == nil {
		return nil, NoBKHD
	}
	if bnk.HIRC() == nil {
		return nil, NoHIRC
	}

	return bnk, nil
}

func BKHDRoutine(
	ctx context.Context,
	r *wio.Reader,
	c chan *DecodeResult,
	I uint8,
	T []byte,
	size uint32,
) {
	cu, e := ParseBKHD(r, I, T, size)
	c <- &DecodeResult{cu, e}
}

func DIDXRoutine(
	ctx context.Context,
	r *wio.Reader,
	c chan *DecodeResult,
	I uint8,
	T []byte,
	size uint32,
) {
	cu, e := ParseDIDX(r, I, T, size)
	c <- &DecodeResult{cu, e}
}

func HIRCRoutine(
	ctx context.Context,
	r *wio.Reader,
	c chan *DecodeResult,
	I uint8,
	T []byte,
	size uint32,
) {
	cu, e := ParseHIRC(ctx, r, I, T, size)
	c <- &DecodeResult{cu, e}
}

func ParseBKHD(r *wio.Reader, I uint8, T []byte, size uint32) (
	*wwise.BKHD, error,
) {
	assert.Equal(0, r.Pos(), "Parser for BKHD does not start at byte 0.")

	bkhd := wwise.NewBKHD(I, T)

	bankGeneratorVersion, err := r.U32()
	if err != nil {
		return nil, err
	}
	bkhd.BankGenerationVersion = bankGeneratorVersion

	bkhd.SoundbankID, err = r.U32()
	if err != nil {
		return nil, err
	}

	bkhd.LanguageID, err = r.U32()
	if err != nil {
		return nil, err
	}

	altValues, err := r.U32()
	if err != nil {
		return nil, err
	}

	bkhd.Alignment = uint16(altValues & 0xFFFF)
	bkhd.DeviceAllocated = uint16((altValues >> 16) & 0xFFFF)

	bkhd.ProjectID, err = r.U32()
	if err != nil {
		return nil, err
	}

	bkhd.Undefined, err = r.ReadAll()
	if err != nil {
		return nil, err
	}

	assert.Equal(
		size,
		uint32(r.Pos()),
		"There are data that is not consumed after parsing all BKHD blob",
	)

	return bkhd, nil
}

func ParseDIDX(r *wio.Reader, I uint8, T []byte, size uint32) (
	*wwise.DIDX, error,
) {
	assert.Equal(0, r.Pos(), "Parser for DIDX does not start at byte 0.")

	num := size / 0x0c
	didx := wwise.NewDIDX(I, T, num)
	for range num {
		sid, err := r.U32()
		if err != nil {
			return nil, err
		}
		offset, err := r.U32()
		if err != nil {
			return nil, err
		}
		size, err := r.U32()
		if err != nil {
			return nil, err
		}
		didx.MediaIndexs = append(didx.MediaIndexs, &wwise.MediaIndex{
			Sid:    sid,
			Offset: offset,
			Size:   size,
		})
	}

	assert.Equal(
		size,
		uint32(r.Pos()),
		"There are data that is not consumed after parsing all media index blob",
	)

	return didx, nil
}

func CheckHeader(r *wio.Reader) (uint32, error) {
	curr := r.Pos()

	tag, err := r.FourCC()
	if err != nil {
		return 0, err
	}
	if bytes.Compare(tag, []byte{'A', 'K', 'B', 'K'}) == 0 {
		return 0, fmt.Errorf("AKBK chunk indicate this Wwise sound bank is legacy. Legacy version Wwise bank is not supported.")
	}
	if bytes.Compare(tag, []byte{'B', 'K', 'H', 'D'}) != 0 {
		return 0, errors.New("This file is not a Wwise sound bank.")
	}

	_, err = r.U32() // size

	version, err := r.U32()
	if err != nil {
		return 0, err
	}
	if version == 0 || version == 1 {
		_, err := r.U32()
		if err != nil {
			return 0, err
		}
		version, err = r.U32()
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("Legacy version %d of Wwise sound bank is not supported.", version)
	}

	_, in := sort.Find(len(CustomVersions), func(i int) int {
		if version < CustomVersions[i] {
			return -1
		} else if version == CustomVersions[i] {
			return 0
		} else {
			return 1
		}
	})
	if in {
		return 0, fmt.Errorf("Custom version %d of Wwise sound bank is not supported yet.", version)
	}

	if version & 0xFFFF0000 == 0x80000000 {
		version = version & 0x0000FFFF
		return 0, fmt.Errorf("Unknown custom version %d of Wwise sound bank is not supported yet.", version)
	}

	if version & 0x0FFFF000 > 0 {
		return 0, fmt.Errorf("Encrypted bank version %d Wwise sound bank. Decryption of Wwise sound bank version is not supported yet.", version)
	}

	_, in = sort.Find(len(Versions), func(i int) int {
		if version < Versions[i] {
			return -1
		} else if version == Versions[i] {
			return 0
		} else {
			return 1
		}
	})
	if !in {
		return 0, fmt.Errorf("Unknown bank version %d Wwise sound bank is not supported yet.", version)
	}

	err = r.SeekStart(curr)

	return version, err
}
