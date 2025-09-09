package decoder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func Decode(
	ctx  context.Context, 
	p    string, 
	o    order,
	opt *DecoderOption,
) (b *wwise.Bank, err error) {
	if opt == nil {
		opt = &DecoderOption{
			DecoderBufferSize: DecodeBufferSize,
			DecodedChunkRoutine: 4,
		}
		opt.IncludeDATA()
		opt.IncludeMETA()
	}

	b = wwise.NewBank()

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, int(opt.DecoderBufferSize))

	var chunkNameBytes []byte = make([]byte, 4, 4)
	_, err = reader.Read(chunkNameBytes)
	if err != nil {
		return nil, err
	}
	chunkName := wwise.ChunkName(chunkNameBytes)

	if chunkName != wwise.ChunkNameBKHD {
		return nil, WrongBKHDPosition(p)
	}

	bkhd, err := DecodeBKHD(p, reader, o)
	if err != nil {
		return nil, err
	}
	wwise.RegBKHD(b, bkhd)

	slog.Info("Parsed BKHD")

	pos := u8(1)
	for {
		_, err = reader.Read(chunkNameBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chunkName = wwise.ChunkName(chunkNameBytes)

		chunkSize, err := uio.U32(reader, o)
		if err != nil {
			return nil, err
		}
		slog.Debug(fmt.Sprintf("Locate %s chunk (size %d)", chunkName, chunkSize))

		if wwise.HasChunk(b, chunkName) {
			return nil, fmt.Errorf("Duplicated chunk %s", chunkName)
		}
		
		switch chunkName {
		case wwise.ChunkNameDIDX:
			didxdata, err := DecodeDIDX(reader, chunkSize, o)
			if err != nil {
				return nil, err
			}
			if err = wwise.RegDIDXDATA(b, didxdata, pos); err != nil {
				return nil, err
			}
			slog.Info("Parsed DIDX")
		default:
			encoded := make([]byte, chunkSize, chunkSize)
			if _, err = io.ReadFull(reader, encoded); err != nil {
				return nil, err
			}
			if err = wwise.NewEncodedChunk(b, chunkName, pos, encoded); err != nil {
				return nil, err
			}
			slog.Warn(fmt.Sprintf("Skipping (storing as encoded chunk) chunk %s (size = %d)", chunkName, chunkSize))
		}
		pos += 1
	}

	in := wwise.HasChunk(b, "DATA")
	if opt.IsIncludeDATA() && in {
		_, chunk := wwise.PopEncodedChunk(b, "DATA")
		DecodeDATA(b.DIDXDATA, chunk)
		slog.Info("Parsed DATA chunk")
	}

	return b, nil
}
