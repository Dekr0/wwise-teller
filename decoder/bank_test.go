package decoder_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dekr0/unwise/decoder"
	"github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

var SoundBanksDir string = os.Getenv("SOUNDBANKS")

func TestDecodeComplex(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})))

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}
	hircDecodeOpt.NumRoutine = 4

	const bank = "content_audio_weapons_superearth.st_bnk"
	bnk, err := decoder.AllocDecode(
		t.Context(),
		filepath.Join(SoundBanksDir, bank),
		binary.LittleEndian,
		&bankDecodeOpt,
		&hircDecodeOpt,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed decode %s bank: %w", bank, err))
	}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	f, err := os.Create(bank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to create bank %s: %w", bank, err))
	}
	w := bufio.NewWriterSize(f, decoder.PageSize32k)

	err = wwise.EncodeBank(
		t.Context(),
		&io.EncoderCtx{
			Writer: w,
			Order:  binary.LittleEndian,
			Count:  0,
		},
		bnk,
		&bankEncodeOption,
		&hircEncodeOption,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to encode bank: %w", err))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(fmt.Errorf("Failed to flush: %w", err))
	}
	if err := f.Close(); err != nil {
		t.Fatal(fmt.Errorf("Failed to close %s: %w", bank, err))
	}

	data, err := os.ReadFile(bank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", bank, err))
	}

	expectData, err := os.ReadFile(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("Diff test fail")
	}

	os.Remove(bank)
}

func TestDecodeAll(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})))

	entries, err := os.ReadDir(SoundBanksDir)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed read information about sound bank directory: %w", err))
	}

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}
	hircDecodeOpt.NumRoutine = 0

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	for _, entry := range entries {
		bank := entry.Name()
		outputBank := filepath.Join("output", bank)
		inputBank := filepath.Join(SoundBanksDir, bank)
		t.Run(fmt.Sprintf("Running decoding test on %s", bank), func(t *testing.T) {
			bnk, err := decoder.AllocDecode(
				t.Context(),
				inputBank,
				binary.LittleEndian,
				&bankDecodeOpt,
				&hircDecodeOpt,
			)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed decode %s bank: %w", bank, err))
			}

			f, err := os.Create(outputBank)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed to create bank %s: %w", bank, err))
			}

			w := bufio.NewWriterSize(f, decoder.PageSize32k)
			if err = wwise.EncodeBank(
				t.Context(),
				&io.EncoderCtx{
					Writer: w,
					Order:  binary.LittleEndian,
					Count:  0,
				},
				bnk,
				&bankEncodeOption,
				&hircEncodeOption,
			); err != nil {
				t.Fatal(fmt.Errorf("Failed to encode bank: %w", err))
			}
			if err := w.Flush(); err != nil {
				t.Fatal(fmt.Errorf("Failed to flush: %w", err))
			}
			if err := f.Close(); err != nil {
				t.Fatal(fmt.Errorf("Failed to close %s: %w", bank, err))
			}

			data, err := os.ReadFile(outputBank)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", bank, err))
			}

			expectData, err := os.ReadFile(inputBank)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(data, expectData) != 0 {
				t.Fatal("Diff test fail")
			}

			os.Remove(outputBank)
		})
	}
}

func TestDecodeInit(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})))

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}
	hircDecodeOpt.NumRoutine = 4

	const bank = "content_audio_Init.st_bnk"
	bnk, err := decoder.AllocDecode(
		t.Context(),
		filepath.Join(SoundBanksDir, bank),
		binary.LittleEndian,
		&bankDecodeOpt,
		&hircDecodeOpt,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed decode %s bank: %w", bank, err))
	}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	f, err := os.Create(bank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to create bank %s: %w", bank, err))
	}
	w := bufio.NewWriterSize(f, decoder.PageSize32k)

	err = wwise.EncodeBank(
		t.Context(),
		&io.EncoderCtx{
			Writer: w,
			Order:  binary.LittleEndian,
			Count:  0,
		},
		bnk,
		&bankEncodeOption,
		&hircEncodeOption,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to encode bank: %w", err))
	}
	if err := w.Flush(); err != nil {
		t.Fatal(fmt.Errorf("Failed to flush: %w", err))
	}
	if err := f.Close(); err != nil {
		t.Fatal(fmt.Errorf("Failed to close %s: %w", bank, err))
	}

	data, err := os.ReadFile(bank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", bank, err))
	}

	expectData, err := os.ReadFile(filepath.Join(SoundBanksDir, bank))
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("Diff test fail")
	}

	os.Remove(bank)
}
