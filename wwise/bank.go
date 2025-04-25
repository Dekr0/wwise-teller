package wwise

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/Dekr0/wwise-teller/assert"
)

var BankVersion = -1

const chunkHeaderSize = 4 + 4

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

func (bnk *Bank) Encode(ctx context.Context) ([]byte, error) {
	chunks := make([][]byte, 4, 4)
	for i := range chunks {
		chunks[i] = []byte{}
	}

	bkhd := make(chan []byte, 1)
	didx := make(chan []byte, 1)
	data := make(chan []byte, 1)
	hirc := make(chan []byte, 1)
	e := make(chan error)
	pending := 0

	/* Header section */
	go func() { bkhd <- bnk.BKHD.Encode() }()
	pending += 1

	/* DIDX section */
	if bnk.DIDX != nil {
		go func() { didx <- bnk.DIDX.Encode() }()
		pending += 1
	}

	/* DATA section */
	if bnk.DATA != nil {
		go func() { data <- bnk.DATA.Encode() }()
		pending += 1
	}

	/* HIRC section */
	go func() {
		b, err := bnk.HIRC.Encode(ctx)
		if err != nil {
			e <- errors.Join(errors.New("Failed to encode HIRC"), err)
		} else {
			hirc <- b
		}
	}()
	pending += 1

	for pending > 0 {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case err := <- e:
			return nil, err
		case b := <- bkhd:
			assert.NotNil(b, "BKHD")
			chunks[0] = b
			pending -= 1
			slog.Info("Encoded BKHD section")
		case b := <- didx:
			assert.NotNil(b, "DIDX")
			chunks[1] = b
			pending -= 1
			slog.Info("Encoded DIDX section")
		case b := <- data:
			assert.NotNil(b, "DATA")
			chunks[2] = b
			pending -= 1
			slog.Info("Encoded DATA section")
		case b := <- hirc:
			assert.NotNil(b, "HIRC")
			chunks[3] = b
			pending -= 1
			slog.Info("Encoded HIRC section")
		}
	}

	for _, chunk := range chunks {
		slog.Info("Chunk size", "index", len(chunk))
	}

	return bytes.Join(chunks, []byte{}), nil
}
