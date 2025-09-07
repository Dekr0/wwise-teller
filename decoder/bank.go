package decoder

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func DecodeStreaming(
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
			_, err = reader.Discard(int(size))
			if err != nil {
				return nil, err
			}
		}
	}
	return b, nil
}

func DecodeMem(
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

	mem, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(mem) // len(mem) is equal to bytes.Reader.Size()

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
		
		switch chunkName {
		default:
			size, err := uio.U32(reader, o)
			if err != nil {
				return nil, err
			}

			// Use slicing to avoid copying
			// Total size - bytes.Reader.Len() = current postion
			// Slice the backing buffer from current position to 
			// current position + chunk size
			if _, err = reader.Seek(int64(size), io.SeekCurrent); err != nil {
				return nil, err
			}
		}

		pos += 1
	}

	return b, nil
}
