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
		return nil, fmt.Errorf("Must provide a bank decoder option")
	}

	b = wwise.NewBank()

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %w", p, err)
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, int(opt.DecoderBufferSize))

	var chunkNameBytes []byte = make([]byte, 4, 4)
	_, err = reader.Read(chunkNameBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed read chunk name of first chunk: %w", err)
	}
	chunkName := wwise.ChunkName(chunkNameBytes)

	if chunkName != wwise.ChunkNameBKHD {
		return nil, WrongBKHDPosition(p)
	}

	bkhd, err := DecodeBKHD(p, reader, o)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode BKHD: %w", err)
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
			return nil, fmt.Errorf("Failed to read chunk name of %d-th chunk: %w", pos, err)
		}
		chunkName = wwise.ChunkName(chunkNameBytes)

		chunkSize, err := uio.U32(reader, o)
		if err != nil {
			return nil, fmt.Errorf("Failed to read chunk size of %d-th chunk: %w", pos, err)
		}

		slog.Info(fmt.Sprintf("Parsing %s chunk...", chunkName),
			"position", pos,
			"size", chunkSize,
		)

		if wwise.HasChunk(b, chunkName) {
			return nil, fmt.Errorf("A duplicated %s chunk at position %d", chunkName, pos)
		}
		
		switch chunkName {
		case wwise.ChunkNameDIDX:
			didxdata, err := DecodeDIDX(reader, chunkSize, o)
			if err != nil {
				return nil, fmt.Errorf("Failed to decode DIDX chunk at position %d: %w", pos, err)
			}
			wwise.RegDIDXDATA(b, didxdata, pos)
			slog.Info("Parsed DIDX", "position", pos, "size", chunkSize)
		default:
			encoded := make([]byte, chunkSize, chunkSize)
			if _, err = io.ReadFull(reader, encoded); err != nil {
				return nil, fmt.Errorf(
					"Failed to read %d bytes of data for %s chunk at position %d: %w",
					chunkSize, chunkName, pos, err,
				)
			}
			wwise.NewEncodedChunk(b, chunkName, pos, encoded)
			if chunkName == "DATA" {
				slog.Info("Store encoded DATA chunk and delay its decoding",
					"position", pos,
					"size", chunkSize,
				)
			} else {
				slog.Info(
					fmt.Sprintf("Not decode %s chunk and store encoded data as it is", 
						chunkName,
					),
					"position", pos,
					"size", chunkSize,
				)
			}
		}
		pos += 1
	}

	in := wwise.HasChunk(b, "DATA")
	if opt.IsIncludeDATA() && in {
		_, chunk := wwise.PopEncodedChunk(b, "DATA")
		DecodeDATA(b.DIDXDATA, chunk)
		slog.Info("Parsed DATA chunk", "size", len(chunk))
	}

	return b, nil
}
