package scripts

import (
	"context"
	"testing"

	"github.com/Dekr0/wwise-teller/integration/helldivers"
	"github.com/Dekr0/wwise-teller/parser"
	"github.com/Dekr0/wwise-teller/wwise"
)

func TestRemoveRadioFX(t *testing.T) {
	initBnk := "../tests/default_st_bnks/content_audio_Init.st_bnk"
	bnk, err := parser.ParseBank(initBnk, context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
	hirc := bnk.HIRC()
	if hirc == nil {
		t.Fatal()
	}
	v, in := hirc.Buses.Load(uint32(1149102368))
	if !in {
		t.Fatal()
	}
	bus := v.(*wwise.Bus)
	bus.BusFxParam.FxChunk.BitsFxByPass = 1
	for i := range bus.BusFxParam.FxChunk.FxChunkItems {
		bus.BusFxParam.FxChunk.FxChunkItems[i].BitVector = 1
	}

	bnkData, err := bnk.Encode(context.Background(), true, false)
	if err != nil {
		t.Fatal(err)
	}
	if bnk.META() == nil {
		t.Fatal(err)
	}
	err = helldivers.GenHelldiversPatchStable(bnkData, bnk.META().B, ".")
	if err != nil {
		t.Fatal(err)
	}
}
