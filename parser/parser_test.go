package parser

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type malformedSoundbank struct {
	name string
	err  error
}

func TestParseBank(t *testing.T) {
	banks, err := os.ReadDir("../tests/bnk")
	if err != nil {
		t.Fatal(err)
	}

	excludes := []*malformedSoundbank{}

	for _, bank := range banks {
		t.Log(bank.Name())
		bnkPath := filepath.Join("../tests/bnk", bank.Name())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
		bnk, err := ParseBank(bnkPath, ctx)
		if err != nil {
			cancel()
			if err == NoBKHD || err == NoDATA || err == NoDIDX || err == NoHIRC {
				excludes = append(excludes, &malformedSoundbank{bank.Name(), err})
				continue
			} else {
				t.Fatal(err)
			}
		}
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*4)
		blob, err := bnk.Encode(ctx)
		if err != nil {
			cancel()
			t.Fatal(bnkPath, err)
		}
		cancel()

		orig, err := os.ReadFile(bnkPath)
		if err != nil {
			t.Fatal(bnkPath, err)
		}

		if bytes.Compare(blob, orig) != 0 {
			if len(blob) > len(orig) {
				for i := range orig {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank.Name(), i, orig[i], blob[i])
					}
				}
			} else {
				for i := range blob {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank.Name(), i, orig[i], blob[i])
					}
				}
			}
		}
	}

	t.Log("Malformed sound bank")
	for _, exclude := range excludes {
		t.Log(exclude.name, exclude.err)
	}
}

func TestParseMusicBank(t *testing.T) {
	banks, err := os.ReadDir("../tests/bnk")
	if err != nil {
		t.Fatal(err)
	}

	excludes := []*malformedSoundbank{}

	for _, bank := range banks {
		if !strings.HasPrefix(bank.Name(), "music_") {
			continue
		}
		t.Log(bank.Name())

		bnkPath := filepath.Join("../tests/bnk", bank.Name())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
		bnk, err := ParseBank(bnkPath, ctx)
		if err != nil {
			cancel()
			if err == NoBKHD || err == NoDATA || err == NoDIDX || err == NoHIRC {
				excludes = append(excludes, &malformedSoundbank{bank.Name(), err})
				continue
			} else {
				t.Fatal(err)
			}
		}
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*360)
		blob, err := bnk.Encode(ctx)
		if err != nil {
			cancel()
			t.Fatal(bnkPath, err)
		}
		cancel()

		orig, err := os.ReadFile(bnkPath)
		if err != nil {
			t.Fatal(bnkPath, err)
		}

		if bytes.Compare(blob, orig) != 0 {
			if len(blob) > len(orig) {
				for i := range orig {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank.Name(), i, orig[i], blob[i])
					}
				}
			} else {
				for i := range blob {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank.Name(), i, orig[i], blob[i])
					}
				}
			}
		}
	}

	t.Log("Malformed sound bank")
	for _, exclude := range excludes {
		t.Log(exclude.name, exclude.err)
	}
}

func TestParseBankCheckSize(t *testing.T) {
	bnk, err := ParseBank("../tests/bnk/wep_mg43.bnk", context.Background())
	if err != nil {
		t.Fatal(err)
	}
	bnk.Encode(context.Background())
}
