package decoder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"sync"

	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func AllocDecode(
	ctx  context.Context, 
	p    string, 
	o    order,
	bankOpt *BankDecodeOption,
	hircOpt *HircDecodeOption,
) (b *wwise.Bank, err error) {
	if bankOpt == nil {
		return nil, fmt.Errorf("Must provide a bank decoder option")
	}
	if hircOpt == nil {
		return nil, fmt.Errorf("Mus provide a HIRC decoder option")
	}

	if bankOpt.NumDecoder <= 0 {
		bankOpt.NumDecoder = 1
	}

	headers, err := PrefetchChunkHeaders(ctx, p, o)
	if err != nil {
		return nil, err
	}
	
	var tagBytes []byte = make([]byte, 4, 4)
	var tag string
	{
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf("Failed to open %s: %w", p, err)
		}
		defer f.Close()

		header, in := headers[wwise.TagBKHD]
		if !in {
			panic("Chunk headers prefetch process did not locate BKHD!")
		}
		delete(headers, wwise.TagBKHD)
		reader := bufio.NewReaderSize(f, int(header.Size + uio.Size32 * 2))

		_, err = io.ReadFull(reader, tagBytes)
		if err != nil {
			return nil, fmt.Errorf("Failed read chunk tag of first chunk: %w", err)
		}
		tag = wwise.Tag(tagBytes)

		if tag != wwise.TagBKHD {
			return nil, WrongBKHDPosition(p)
		}

		bkhd, err := AllocDecodeBKHD(p, reader, o)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode BKHD: %w", err)
		}

		b = wwise.AllocBank()
		b.RegBKHD(bkhd)

		slog.Info("Parsed BKHD")
		f.Close()
	}

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var mu sync.Mutex
	var wa sync.WaitGroup
	decoderC := make(chan func() error)
	errorC := make(chan error, len(headers))
	subroutine := func() {
		defer wa.Done()
		for {
			select {
			case <- subCtx.Done():
				return
			case f := <- decoderC:
				errorC <- f()
			default:
			}
		}
	}

	for range bankOpt.NumDecoder {
		wa.Add(1)
		go subroutine()
	}

	for tag, header := range headers {
		f, err := os.Open(p)
		if err != nil {
			return nil, fmt.Errorf(
				"Failed to open a new file descriptor to decode %s chunk: %w", 
				tag, err,
			)
		}

		tagCopy := tag
		headerCopy := header

		switch tag {
		case wwise.TagDIDX:
			decoder := func() error {
				defer f.Close()
				sr := io.NewSectionReader(f, int64(headerCopy.Pos), int64(headerCopy.Size))
				var bReader *bufio.Reader
				if bankOpt.DecoderBufferSize > headerCopy.Size {
					bReader = bufio.NewReaderSize(sr, int(headerCopy.Size) * 2)
				} else {
					bReader = bufio.NewReaderSize(sr, int(bankOpt.DecoderBufferSize))
				}

				audioStore, err := AllocDecodeDIDX(bReader, headerCopy.Size, o)
				if err != nil {
					return fmt.Errorf("Failed to decode DIDX chunk: %w", err)
				}
				mu.Lock()
				b.RegDIDXDATA(audioStore, headerCopy.Idx)
				mu.Unlock()
				slog.Info("Decoded DIDX", "position", headerCopy.Idx, "size", headerCopy.Size)
				return nil
			}
			decoderC <- decoder
		case wwise.TagHIRC:
			decoder := func() error {
				defer f.Close()

				stat := wwise.HierarchyStat{}
				if _, err := f.Seek(headerCopy.Pos, io.SeekStart); err != nil {
					return fmt.Errorf("Failed to seek to HIRC data chunk before prefetching HIRC metadata for space allocation: %w", err)
				}
				if err := PrefetchHIRCMetadata(ctx, f, &stat, o, headerCopy.Size); err != nil {
					return fmt.Errorf("Failed to prefetch HIRC metadata for space allocation: %w", err)
				}
				if _, err := f.Seek(0, io.SeekStart); err != nil {
					return fmt.Errorf("Failed to seek file descriptor back origin after prefetching HIRC metadata for space allocation: %w", err)
				}

				spec := wwise.HIRCSpaceSpec{}
				wwise.EstimateHIRCSpace(&stat, &spec)

				sr := io.NewSectionReader(f, int64(headerCopy.Pos), int64(headerCopy.Size))
				var bReader *bufio.Reader
				if bankOpt.DecoderBufferSize > headerCopy.Size {
					bReader = bufio.NewReaderSize(sr, int(headerCopy.Size) * 2)
				} else {
					bReader = bufio.NewReaderSize(sr, int(bankOpt.DecoderBufferSize))
				}

				h, err := AllocDecodeHIRC(ctx, hircOpt, &spec, bReader, o, header.Size, b.BKHD.Version)
				if err != nil {
					return fmt.Errorf("Failed to decode HIRC chunk: %w", err)
				}

				mu.Lock()
				b.RegHIRC(h, header.Idx)
				slog.Info("Decoded HIRC chunk")
				mu.Unlock()
				return nil
			}
			decoderC <- decoder
		default:
			decoder := func() error {
				defer f.Close()
				sr := io.NewSectionReader(f, int64(headerCopy.Pos), int64(headerCopy.Size))
				var bReader *bufio.Reader
				if bankOpt.DecoderBufferSize > headerCopy.Size {
					bReader = bufio.NewReaderSize(sr, int(headerCopy.Size) * 2)
				} else {
					bReader = bufio.NewReaderSize(sr, int(bankOpt.DecoderBufferSize))
				}

				encoded := make([]byte, headerCopy.Size, headerCopy.Size)
				if _, err := io.ReadFull(bReader, encoded); err != nil {
					return fmt.Errorf("Failed to read %d bytes of data for %s chunk at position %d: %w",
						headerCopy.Size, tagCopy, headerCopy.Idx, err,
					)
				}
				mu.Lock()
				b.Chunk.AddEncodedChunk(tagCopy, headerCopy.Idx, encoded)
				mu.Unlock()
				if tagCopy == "DATA" {
					slog.Info("Stored encoded DATA chunk. Decoding of DATA is delayed at the end.",
						"position", headerCopy.Idx,
						"size", headerCopy.Size,
					)
				} else {
					slog.Info(fmt.Sprintf("%s chunk wasn't decoded, and it was stored as raw encoded data", tag),
						"position", headerCopy.Idx,
						"size", headerCopy.Size,
					)
				}
				return nil
			}
			decoderC <- decoder
		}
	}

	numDecoded := 1
	for numDecoded < len(headers) {
		select {
		case <- subCtx.Done():
			return nil, fmt.Errorf("Sound bank %s decoding process cancel: %w", p, subCtx.Err())
		case err := <- errorC:
			if err != nil {
				return nil, fmt.Errorf("A decoder error occur: %w", err)
			}
			numDecoded++
		default:
		}
	}
	cancel()
	wa.Wait()

	in := b.Chunk.HasChunk(wwise.TagDATA)
	if bankOpt.IsIncludeDATA() && in {
		_, chunk := b.Chunk.PopEncodedChunk(wwise.TagDATA)
		AllocDecodeDATA(b.AudioStore, chunk)
		slog.Info("Parsed DATA chunk", "size", len(chunk))
	}

	return b, nil
}

type ChunkHeader struct {
	Idx  u8
	Size u32
	// Position offset will include 4 bytes of size value. Use `Size` to get size 
	// instead of reading it
	Pos  int64
	Tag  string
}

// TODO: provide an additional metadata chunk or a sepearate metadata file to 
// tell the location of if chunks. An incoming problem -> what if the metadata 
// information is not correct.
func PrefetchChunkHeaders(
	ctx  context.Context, 
	p    string, 
	o    order,
) (headers map[string]ChunkHeader, err error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %w", p, err)
	}
	defer f.Close()

	pos := int64(0)
	
	var tagBuf []byte = make([]byte, 4, 4)
	nread, err := io.ReadFull(f, tagBuf)
	if err != nil {
		return nil, fmt.Errorf("Failed to read chunk tag of first chunk: %w", err)
	}
	pos += int64(nread)

	tag := wwise.Tag(tagBuf)
	if tag != wwise.TagBKHD {
		return nil, WrongBKHDPosition(p)
	}
	size, err := uio.U32(f, o)
	if err != nil {
		return nil, fmt.Errorf("Failed to read BKHD chunk size: %w", err)
	}
	pos += uio.Size32

	headers = make(map[string]ChunkHeader, len(wwise.KnownTag))
	headers[tag] = ChunkHeader{
		Tag: tag, 
		Idx: 0,
		Size: size,
		Pos: pos,
	}

	newOffset, err := f.Seek(int64(size), io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("Failed to skip BKHD chunk (%d bytes ahead) to the next chunk: %w", size, err)
	}
	pos = newOffset

	idx := u8(1)
	for {
		nread, err = io.ReadFull(f, tagBuf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("Failed to read chunk tag of next chunk: %w", err)
		}
		pos += int64(nread)

		tag = wwise.Tag(tagBuf)
		if !slices.Contains(wwise.KnownTag, tag) {
			return nil, fmt.Errorf("Unknown chunk tag %s with size of %d", tag, size)
		}
		if _, in := headers[tag]; in {
			return nil, fmt.Errorf("Duplicate chunk %s", tag)
		}
		header := ChunkHeader{
			Tag: tag, 
			Idx: idx,
			Size: 0,
			Pos: pos,
		}

		size, err = uio.U32(f, o)
		if err != nil {
			return nil, fmt.Errorf("Failed to read chunk size of %s chunk: %w", tag, err)
		}
		newOffset, err = f.Seek(int64(size), io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf(
				"Failed to skip %s chunk (%d bytes ahead) to the next chunk: %w", 
				tag, size, err,
			)
		}
		pos = newOffset

		header.Pos, header.Size = header.Pos + uio.Size32, size
		headers[tag] = header
		idx++
	}
	return headers, nil
}
