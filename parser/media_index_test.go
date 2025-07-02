package parser

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/wwise"
)

func TestMediaIndex(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	bnk, err := ParseBank("../tests/bnks/wep_cr1_adjudicator.bnk", ctx, true)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	fmt.Println(bnk.DIDX().Alignment)
}

func TestAllAlignment(t *testing.T) {
	banks, err := os.ReadDir("../tests/bnks")
	if err != nil {
		t.Fatal(err)
	}

	possibleAlignment := make(map[uint8][]*wwise.DIDX)
	for _, bank := range banks {
		t.Log(bank.Name())
		bnkPath := filepath.Join("../tests/bnks", bank.Name())

		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 360)
		bnk, err := ParseBank(bnkPath, ctx, true)
		if err != nil {
			cancel()
			if err == NoBKHD || err == NoDATA || err == NoDIDX || err == NoHIRC {
				continue
			} else {
				t.Fatal(err)
			}
		}
		cancel()
		if bnk.DIDX() == nil {
			continue
		}
		if d, in := possibleAlignment[bnk.DIDX().Alignment]; !in {
			possibleAlignment[bnk.DIDX().Alignment] = []*wwise.DIDX{bnk.DIDX()}
		} else {
			possibleAlignment[bnk.DIDX().Alignment] = append(d, bnk.DIDX())
		}
	}
	fmt.Println("Possible Alignments:")
	for a, d := range possibleAlignment {
		fmt.Println(a)
		fmt.Printf("Other available alignments for DIDX has max alignment of %d\n", a)
		for _, di := range d {
			fmt.Println(di.AvailableAlignment)
		}
	}
}
