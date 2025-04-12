package parser

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"

	"github.com/Dekr0/wwise-teller/assert"
	"github.com/Dekr0/wwise-teller/reader"
	"github.com/Dekr0/wwise-teller/wwise"
)

var MissingBKHDError error = errors.New("Sound bank is missing BKHD section")
var MissingDIDXError error = errors.New("Sound bank is missing DIDX section")
var MissingDATAError error = errors.New("Sound bank is missing DATA section")
var MissingHIRCError error = errors.New("Sound bank is missing HIRC section")

var (
	bankCustomVersions []uint32 = []uint32{ 122, 126, 129, 135, 136 }
	bankVersions []uint32 = []uint32 {
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

func ParseBank(filename string, ctx context.Context) (*wwise.Bank, error) {
	f, openErr := os.Open(filename)
	if openErr != nil {
		return nil, openErr
	}

	shallowReader := reader.NewSoundbankReader(f, binary.LittleEndian)

	bnkVersion, checkErr := checkHeader(shallowReader)
	if checkErr != nil {
		return nil, checkErr
	}

	if bnkVersion != 141 {
		return nil, errors.New("Wwise teller currently only targets version 141.")
	}

	/*
	Each channel (except error channel) must be only written once
	*/
	cBKHD := make(chan *wwise.BKHD, 1)
	cDIDX := make(chan *wwise.DIDX, 1)
	cHIRC := make(chan *wwise.HIRC, 1)
	cParserErr  := make(chan error, 4)
	scheduledBKHD := false
	scheduledDIDX := false
	scheduledHIRC := false
	uncollectedResult := 0

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
	bnk := wwise.NewBank(nil, nil, nil, nil)
	var topLevelShallowReadErr error
	for topLevelShallowReadErr == nil {
		chunkTag, shallowReadErr := shallowReader.FourCC()
		if shallowReadErr != nil {
			topLevelShallowReadErr = shallowReadErr
			continue
		}

		chunkSize, shallowReadErr := shallowReader.U32()
		if shallowReadErr != nil {
			topLevelShallowReadErr = shallowReadErr
			continue
		}

		if bytes.Compare(chunkTag, []byte{'B', 'K', 'H', 'D'}) == 0 {
			bkhdReader, shallowReadErr := shallowReader.NewSectionBankReader(uint64(chunkSize))
			if shallowReadErr != nil {
				topLevelShallowReadErr = shallowReadErr
				continue
			}
			slog.Info("Start parsing BKHD section...", "chunkSize", chunkSize)
			go func() {
				/* Should always return */
				bkhd, parserErr := parseBKHD(chunkSize, bkhdReader)
				if parserErr != nil {
					cParserErr <- parserErr
				} else {
					cBKHD <- bkhd
				}
				slog.Info("Finished parsing BKHD", 
					"chunkSize", chunkSize,
					"consumeSize", bkhdReader.Tell(),
				)
			}()
			uncollectedResult += 1
			scheduledBKHD = true
		} else if bytes.Compare(chunkTag, []byte{'D', 'A', 'T', 'A'}) == 0 {
			blob, shallowReadErr := shallowReader.ReadFull(uint64(chunkSize), 0)
			if shallowReadErr != nil {
				topLevelShallowReadErr = shallowReadErr
				continue
			}
			bnk.DATA = wwise.NewDATA(blob)
			slog.Info("Read DATA section")
		} else if bytes.Compare(chunkTag, []byte{'D', 'I', 'D', 'X'}) == 0 {
			didxReader, shallowReadErr := shallowReader.NewSectionBankReader(uint64(chunkSize))
			if shallowReadErr != nil {
				topLevelShallowReadErr = shallowReadErr
				continue
			}
			slog.Info("Start parsing DIDX section...")
			go func() {
				/* Should always return */
				didx, parserErr := parseMediaIndex(chunkSize, didxReader)
				if parserErr != nil {
					cParserErr <- parserErr
				} else {
					cDIDX <- didx
				}
				slog.Info("Finished parsing DIDX", 
					"chunkSize", chunkSize,
					"consumeSize", didxReader.Tell(),
				)
			}()
			uncollectedResult += 1
			scheduledDIDX = true
		} else if bytes.Compare(chunkTag, []byte{'H', 'I', 'R', 'C'}) == 0 {
			hircReader, shallowReadErr := shallowReader.NewSectionBankReader(uint64(chunkSize))
			if shallowReadErr != nil {
				topLevelShallowReadErr = shallowReadErr
				continue
			}
			slog.Info("Start parsing HIRC section...")
			go func() {
				/* Should always return */
				hirc, parserErr := parseHIRC(ctx, chunkSize, hircReader)
				if parserErr != nil {
					cParserErr <- parserErr
				} else {
					cHIRC <- hirc
				}
				slog.Info("Finished parsing HIRC", 
					"chunkSize", chunkSize,
					"consumeSize", hircReader.Tell(),
				)
			}()
			uncollectedResult += 1
			scheduledHIRC = true
		}
	}

	if topLevelShallowReadErr != io.EOF {
		return nil, topLevelShallowReadErr
	}

	if !scheduledBKHD {
		return nil, MissingBKHDError
	}
	if !scheduledDIDX {
		slog.Warn("Sound bank is missing DIDX section. This might not be a " +
			"problem as long as HIRC section exists.", "soundbank", filename)
	}
	if !scheduledHIRC {
		return nil, MissingHIRCError
	}

	for uncollectedResult > 0 {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case err := <- cParserErr:
			return nil, err
		case bkhd := <- cBKHD:
			bnk.BKHD = bkhd
			uncollectedResult -= 1
			slog.Info("Collect BKHD parsing result")
		case didx := <- cDIDX:
			bnk.DIDX = didx
			uncollectedResult -= 1
			slog.Info("Collect DIDX parsing result")
		case hirc := <- cHIRC:
			bnk.HIRC = hirc
			uncollectedResult -= 1
			slog.Info("Collect HIRC parsing result")
		}
	}

	if bnk.BKHD == nil {
		return nil, MissingBKHDError
	}
	if bnk.HIRC == nil {
		return nil, MissingHIRCError
	}

	return bnk, nil
}

func parseBKHD(chunkSize uint32, r *reader.BankReader) (*wwise.BKHD, error) {
	assert.AssertEqual(0, r.Tell(), "Parser for BKHD does not start at byte 0.")

	bkhd := wwise.NewBKHD()

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

	assert.AssertEqual(
		chunkSize,
		uint32(r.Tell()),
		"There are data that is not consumed after parsing all BKHD blob",
	)

	return bkhd, nil
}

func parseMediaIndex(chunkSize uint32, r *reader.BankReader) (*wwise.DIDX, error) {
	assert.AssertEqual(0, r.Tell(), "Parser for DIDX does not start at byte 0.")

	didx := wwise.NewDIDX(chunkSize / 0x0c)

	for range didx.UniqueNumMedias {
		audioSrcId, err := r.U32()
		if err != nil {
			return nil, err
		}
		DATAOffset, err := r.U32()
		if err != nil {
			return nil, err
		}
		DATABlobSize, err := r.U32()
		if err != nil {
			return nil, err
		}
		didx.MediaIndexs = append(didx.MediaIndexs, &wwise.MediaIndex{
			AudioSrcId: audioSrcId,
			DATAOffset: DATAOffset,
			DATABlobSize: DATABlobSize,
		})
	}

	assert.AssertEqual(
		chunkSize,
		uint32(r.Tell()),
		"There are data that is not consumed after parsing all media index blob",
	)

	return didx, nil
}

func checkHeader(r *reader.BankReader) (uint32, error) {
	curr := r.Tell()

	fourcc, err := r.FourCC()
	if err != nil {
		return 0, err
	}
	if bytes.Compare(fourcc, []byte{'A', 'K', 'B', 'K'}) == 0 {
		return 0, fmt.Errorf("AKBK chunk indicate this Wwise sound bank is legacy. Legacy version Wwise bank is not supported.")
	}
	if bytes.Compare(fourcc, []byte{'B', 'K', 'H', 'D'}) != 0 {
		return 0, errors.New("This file is not a Wwise sound bank.")
	}

	_, err = r.U32() // size

	bnkVersion, err := r.U32()
	if err != nil {
		return 0, err
	}
	if bnkVersion == 0 || bnkVersion == 1 {
		_, err := r.U32()
		if err != nil {
			return 0, err
		}
		bnkVersion, err = r.U32()
		if err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("Legacy version %d of Wwise sound bank is not supported.", bnkVersion)
	}

	_, in := sort.Find(len(bankCustomVersions), func(i int) int {
		if bnkVersion < bankCustomVersions[i] {
			return -1
		} else if bnkVersion == bankCustomVersions[i] {
			return 0
		} else {
			return 1
		}
	})
	if in {
		return 0, fmt.Errorf("Custom version %d of Wwise sound bank is not supported yet.", bnkVersion)
	}

	if bnkVersion & 0xFFFF0000 == 0x80000000 {
		bnkVersion = bnkVersion & 0x0000FFFF
		return 0, fmt.Errorf("Unknown custom version %d of Wwise sound bank is not supported yet.", bnkVersion)
	}

	if bnkVersion & 0x0FFFF000 > 0 {
		return 0, fmt.Errorf("Encrypted bank version %d Wwise sound bank. Decryption of Wwise sound bank version is not supported yet.", bnkVersion)
	}

	_, in = sort.Find(len(bankVersions), func(i int) int {
		if bnkVersion < bankVersions[i] {
			return -1
		} else if bnkVersion == bankVersions[i] {
			return 0
		} else {
			return 1
		}
	})
	if !in {
		return 0, fmt.Errorf("Unknown bank version %d Wwise sound bank is not supported yet.", bnkVersion)
	}

	err = r.AbsSeek(curr)

	return bnkVersion, err
}
