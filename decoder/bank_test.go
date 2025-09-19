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

var TestLogger *slog.Logger = slog.New(slog.NewTextHandler(
	os.Stdout,
	&slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	},
))

func TestDecodeComplex(t *testing.T) {
	slog.SetDefault(TestLogger)

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.NumDecoder = 4
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}

	SoundBanksDir = "/mnt/d/wwise-teller/storage/soundbanks/vanilla/latest"

	const bankName = "content_audio_weapons_superearth.st_bnk"
	inputBank := filepath.Join(SoundBanksDir, bankName)
	bnk, err := decoder.AllocDecode(
		t.Context(),
		inputBank,
		binary.LittleEndian,
		&bankDecodeOpt,
		&hircDecodeOpt,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed decode %s bank: %w", inputBank, err))
	}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	outputBank := filepath.Join("output", bankName)
	f, err := os.Create(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to create bank %s: %w", outputBank, err))
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
		t.Fatal(fmt.Errorf("Failed to close %s: %w", outputBank, err))
	}

	data, err := os.ReadFile(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", inputBank, err))
	}

	expectData, err := os.ReadFile(inputBank)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("Diff test fail")
	}
}

func TestDecodeAll(t *testing.T) {
	slog.SetDefault(TestLogger)

	entries, err := os.ReadDir(SoundBanksDir)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed read information about sound bank directory: %w", err))
	}

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.ExcludeDATA()
	bankDecodeOpt.NumDecoder = 4
	hircDecodeOpt := decoder.HircDecodeOption{}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	for _, entry := range entries {
		bankName := entry.Name()
		outputBank := filepath.Join("output", bankName)
		inputBank := filepath.Join(SoundBanksDir, bankName)
		t.Run(fmt.Sprintf("Running decoding test on %s", bankName), func(t *testing.T) {
			bnk, err := decoder.AllocDecode(
				t.Context(),
				inputBank,
				binary.LittleEndian,
				&bankDecodeOpt,
				&hircDecodeOpt,
			)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed decode %s bank: %w", inputBank, err))
			}

			f, err := os.Create(outputBank)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed to create bank %s: %w", outputBank, err))
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
				t.Fatal(fmt.Errorf("Failed to close %s: %w", outputBank, err))
			}

			data, err := os.ReadFile(outputBank)
			if err != nil {
				t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", outputBank, err))
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

// Test sound banks that can be easily failed
func TestDecodeFault(t *testing.T) {
	slog.SetDefault(TestLogger)

	SoundBanksDir = "/mnt/d/wwise-teller/storage/soundbanks/vanilla/latest"

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.NumDecoder = 4
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}

	const bankName = "content_audio_music_mission_tutorial.st_bnk"
	inputBank := filepath.Join(SoundBanksDir, bankName)
	bnk, err := decoder.AllocDecode(
		t.Context(),
		inputBank,
		binary.LittleEndian,
		&bankDecodeOpt,
		&hircDecodeOpt,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed decode %s bank: %w", inputBank, err))
	}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	outputBank := filepath.Join("output", bankName)
	f, err := os.Create(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to create bank %s: %w", outputBank, err))
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
		t.Fatal(fmt.Errorf("Failed to close %s: %w", outputBank, err))
	}

	data, err := os.ReadFile(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", outputBank, err))
	}

	expectData, err := os.ReadFile(inputBank)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("Diff test fail")
	}

}

func TestDecodeInit(t *testing.T) {
	slog.SetDefault(TestLogger)

	bankDecodeOpt := decoder.BankDecodeOption{}
	bankDecodeOpt.NumDecoder = 4
	bankDecodeOpt.ExcludeDATA()
	hircDecodeOpt := decoder.HircDecodeOption{}

	const bankName = "content_audio_Init.st_bnk"
	inputBank := filepath.Join(SoundBanksDir, bankName)
	bnk, err := decoder.AllocDecode(
		t.Context(),
		inputBank,
		binary.LittleEndian,
		&bankDecodeOpt,
		&hircDecodeOpt,
	)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed decode %s bank: %w", inputBank, err))
	}

	bankEncodeOption := wwise.EncodeBankOpt{}
	wwise.IncludeEncodedMETA(&bankEncodeOption)
	hircEncodeOption := wwise.EncodeHircOpt{}

	outputBank := filepath.Join("output", bankName)
	f, err := os.Create(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to create bank %s: %w", outputBank, err))
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
		t.Fatal(fmt.Errorf("Failed to close %s: %w", outputBank, err))
	}

	data, err := os.ReadFile(outputBank)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read data from bank %s: %w", outputBank, err))
	}

	expectData, err := os.ReadFile(inputBank)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(data, expectData) != 0 {
		t.Fatal("Diff test fail")
	}
}
