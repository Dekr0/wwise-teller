package helldivers

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wio"
)

func TestExtractSoundBank(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "./9ba626afa44a3aa3.patch_0"

	if err := ExtractSoundBank(nil, target, "."); err != nil {
		t.Fatal(err)
	}
}

func TestParseMeta(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "./content_audio_obj_hellbomb.st_bnk"

	timeout, cancel := context.WithTimeout(context.Background(), time.Second * 4)
	defer cancel()

	bnk, err := parser.ParseBank(target, timeout)
	if err != nil { t.Fatal(err) }

	metaCu := bnk.META()
	if metaCu == nil {
		t.Fail()
		return
	}
	_, err = ParseMETA(nil, wio.NewInPlaceReader(metaCu.Data, wio.ByteOrder))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenHelldiversPatch(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "./content_audio_weapons_superearth.st_bnk"

	timeout, cancel := context.WithTimeout(context.Background(), time.Second * 4)
	defer cancel()

	bnk, err := parser.ParseBank(target, timeout)
	if err != nil { t.Fatal(err) }
	
	bnkData, err := bnk.Encode(timeout)
	if err != nil { t.Fatal(err) }

	if err := GenHelldiversPatch(timeout, bnkData, bnk.META().Data, "."); err != nil {
		t.Fatal(err)
	}
}
