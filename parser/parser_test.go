package parser

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

/*
func TestCheckHeader(t *testing.T) {
	banks, err := os.ReadDir("../tests/bnk")
	if err != nil {
		t.Fatal(err)
	}
	for _, bank := range banks {
		f, err := os.Open(path.Join("../tests/bnk", bank.Name()))
		if err != nil {
			t.Fatal(err)
		}
		r := reader.NewSoundbankReader(f, binary.LittleEndian)
		version, err := checkHeader(r)
		if err != nil {
			t.Fatal(err)
		}
		if version != 141 {
			t.Fail()
		}
		t.Log(version)
	}
}
*/

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

		bnk, err := ParseBank(bnkPath, context.Background())
		if err != nil {
			if err == MissingBKHDError || err == MissingDATAError || err == MissingDIDXError || err == MissingHIRCError {
				excludes = append(excludes, &malformedSoundbank{bank.Name(), err})
				continue
			} else {
				t.Fatal(err)
			}
		}

		blob, err := bnk.Encode(context.Background())
		if err != nil {
			t.Fatal(bnkPath, err)
		}

		orig, err := os.ReadFile(bnkPath)
		if err != nil {
			t.Fatal(bnkPath, err)
		}

		if bytes.Compare(blob, orig) != 0 {
			if len(blob) > len(orig) {
				for i := range orig {
					if blob[i] != orig[i] {
						t.Fatalf("%s: Byte difference at %d. Original: %d, Received: %d\n", bank.Name(), i, orig[i], blob[i])
					}
				}
			} else {
				for i := range blob {
					if blob[i] != orig[i] {
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

/*
func TestParseBankEdgeCase(t *testing.T) {
	edgeCases := []string{
		"../tests/bnk/content_audio_haz_explosivemushroom.bnk",
	}
	for _, edgeCase := range edgeCases {
		t.Log(fmt.Sprintf("Parsing %s", edgeCase))
		bnk, err := ParseBank(edgeCase, context.Background())
		if err != nil {
			t.Fatal(edgeCase, err)
		}

		blob, err := bnk.Encode(context.Background())
		if err != nil {
			t.Fatal(edgeCase, err)
		}

		orig, err := os.ReadFile(edgeCase)
		if err != nil {
			t.Fatal(edgeCase, err)
		}

		if bytes.Compare(blob, orig) != 0 {
			if len(blob) > len(orig) {
				for i := range orig {
					if blob[i] != orig[i] {
						t.Fatalf(edgeCase, "Byte difference at %d. Original: %d, Received: %d\n", i, orig[i], blob[i])
					}
				}
			} else {
				for i := range blob {
					if blob[i] != orig[i] {
						t.Fatalf(edgeCase, "Byte difference at %d. Original: %d, Received: %d\n", i, orig[i], blob[i])
					}
				}
			}
		}
	}
}
*/
