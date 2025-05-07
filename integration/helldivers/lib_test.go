package helldivers

import (
	"context"
	"encoding/binary"
	"os"
	"testing"
	"time"

	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wio"
)

func _TestExtractSoundBank(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian
	target := "/mnt/Program Files/Steam/steamapps/common/Helldivers 2/data/be6c260fadcb8719"
	if err := ExtractSoundBank(nil, target, ".", true); err != nil {
		t.Fatal(err)
	}
}

func TestExtractSoundBankDiff(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	teller := "./inplace/teller"
	stat, err := os.Lstat(teller)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Teller size: ", stat.Size())

	eve := "./inplace/eve"
	stat, err = os.Lstat(eve)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Eve size: ", stat.Size())

	//fmt.Println("eve:")
	//if err := ExtractSoundBank(nil, eve, ".", false); err != nil {
	//	t.Fatal(err)
	//}

	//fmt.Println("teller:")
	//if err := ExtractSoundBank(nil, teller, ".", true); err != nil {
	//	t.Fatal(err)
	//}
}

func _TestParseMeta(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "./content_audio_obj_hellbomb.st_bnk"

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	bnk, err := parser.ParseBank(target, timeout)
	if err != nil {
		t.Fatal(err)
	}

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

func _TestGenHelldiversPatch(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "./guard_dog.st_bnk"

	timeout, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	bnk, err := parser.ParseBank(target, timeout)
	if err != nil {
		t.Fatal(err)
	}

	bnkData, err := bnk.Encode(timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err := GenHelldiversPatch(timeout, bnkData, bnk.META().Data, "."); err != nil {
		t.Fatal(err)
	}
}
