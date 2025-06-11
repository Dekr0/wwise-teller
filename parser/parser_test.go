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
	banks, err := os.ReadDir("../tests/bnks")
	if err != nil {
		t.Fatal(err)
	}

	excludes := []*malformedSoundbank{}

	for _, bank := range banks {
		t.Log(bank.Name())
		bnkPath := filepath.Join("../tests/bnks", bank.Name())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 360)
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

		ctx, cancel = context.WithTimeout(context.Background(), time.Second * 4)
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
	banks, err := os.ReadDir("../tests/bnks")
	if err != nil {
		t.Fatal(err)
	}

	excludes := []*malformedSoundbank{}

	for _, bank := range banks {
		if !strings.HasPrefix(bank.Name(), "music_") {
			continue
		}
		t.Log(bank.Name())

		bnkPath := filepath.Join("../tests/bnks", bank.Name())

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

func TestFaulty(t *testing.T) {
	banks := []string{
		"../tests/bnks/Helldiver_Standard_VO.bnk",
	}
	
	excludes := []*malformedSoundbank{}

	for _, bank := range banks {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*360)
		bnk, err := ParseBank(bank, ctx)
		if err != nil {
			cancel()
			if err == NoBKHD || err == NoDATA || err == NoDIDX || err == NoHIRC {
				excludes = append(excludes, &malformedSoundbank{bank, err})
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
			t.Fatal(bank, err)
		}
		cancel()

		orig, err := os.ReadFile(bank)
		if err != nil {
			t.Fatal(bank, err)
		}

		if bytes.Compare(blob, orig) != 0 {
			if len(blob) > len(orig) {
				for i := range orig {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank, i, orig[i], blob[i])
					}
				}
			} else {
				for i := range blob {
					if blob[i] != orig[i] {
						l, _ := bnk.HIRC().Encode(context.Background())
						fmt.Println(len(l))
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank, i, orig[i], blob[i])
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

func TestTreeRendering(t *testing.T) {
	bnk, err := ParseBank("../tests/bnks/wep_ar19_liberator.bnk", context.Background())
	if err != nil {
		t.Fatal(err)
	}
	bnk.HIRC().BuildTree()
}
