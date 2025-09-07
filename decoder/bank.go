package decoder

import (
	"bufio"
	"context"
	"io"
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
		opt = &DecoderOption{DecoderBufferSize: DecodeBufferSize}
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
	wwise.BankAddBKHD(b, bkhd)

	pos := uint8(1)
	for {
		_, err = reader.Read(chunkNameBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chunkName = wwise.ChunkName(chunkNameBytes)
		
		switch chunkName {
		default:
			size, err := uio.U32(reader, o)
			if err != nil {
				return nil, err
			}
			encoded := make([]byte, size, size)
			if _, err = reader.Read(encoded); err != nil {
				return nil, err
			}
			if err = wwise.BankAddEncodedChunk(b, chunkName, pos, encoded); err != nil {
				return nil, err
			}
		}
		pos += 1
	}
	return b, nil
}
