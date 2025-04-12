package wwise

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/Dekr0/wwise-teller/assert"
)

var BANK_VERSION = -1

const CHUNK_HEADER_SIZE = 4 + 4

type Bank struct {
	BKHD *BKHD
	DIDX *DIDX
	DATA *DATA
	HIRC *HIRC
}

func NewBank(bkhd *BKHD, didx *DIDX, data *DATA, hirc *HIRC) *Bank {
	return &Bank{BKHD: bkhd, DIDX: didx, DATA: data, HIRC: hirc}
}

func (b *Bank) Write(ctx context.Context, w io.Writer) error {
	return nil
}

func (b *Bank) Encode(ctx context.Context) ([]byte, error) {
	dataChunks := make([][]byte, 4, 4)
	for i := range dataChunks {
		dataChunks[i] = []byte{}
	}

	cBKHDblob := make(chan []byte, 1)
	cDIDXblob := make(chan []byte, 1)
	cDATAblob := make(chan []byte, 1)
	cHIRCblob := make(chan []byte, 1)
	cErr := make(chan error)
	uncollectedBlob := 0

	/* Header section */
	go func() {
		 cBKHDblob <- b.BKHD.Encode()
	}()
	uncollectedBlob += 1

	/* DIDX section */
	if b.DIDX != nil {
		go func() {
			cDIDXblob <- b.DIDX.Encode()
		}()
		uncollectedBlob += 1
	}

	/* DATA section */
	if b.DATA != nil {
		go func() {
			cDATAblob <- b.DATA.Encode()
		}()
		uncollectedBlob += 1
	}

	/* HIRC section */
	go func() {
		HIRCBlob, eErr := b.HIRC.Encode()
		if eErr != nil {
			cErr <- errors.Join(errors.New("Failed to encode HIRC"), eErr)
		} else {
			cHIRCblob <- HIRCBlob
		}
	}()
	uncollectedBlob += 1

	for uncollectedBlob > 0 {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case err := <- cErr:
			return nil, err
		case bkhdBlob := <- cBKHDblob:
			assert.AssertNotNil(bkhdBlob, "bkhdData")
			dataChunks[0] = bkhdBlob
			uncollectedBlob -= 1
			slog.Info("Encoded BKHD section")
		case didxBlob := <- cDIDXblob:
			assert.AssertNotNil(didxBlob, "didxData")
			dataChunks[1] = didxBlob
			uncollectedBlob -= 1
			slog.Info("Encoded DIDX section")
		case dataBlob := <- cDATAblob:
			assert.AssertNotNil(dataBlob, "dataBlob")
			dataChunks[2] = dataBlob
			uncollectedBlob -= 1
			slog.Info("Encoded DATA section")
		case hircBlob := <- cHIRCblob:
			assert.AssertNotNil(hircBlob, "hircBlob")
			dataChunks[3] = hircBlob
			uncollectedBlob -= 1
			slog.Info("Encoded HIRC section")
		}
	}

	for _, dataChunk := range dataChunks {
		slog.Info("Chunk size", "index", len(dataChunk))
	}
	bankData := bytes.Join(dataChunks, []byte{})
	return bankData, nil
}
