package helldivers

import (
	"encoding/binary"
	"testing"

	"github.com/Dekr0/wwise-teller/wio"
)

func TestExtractSoundBank(t *testing.T) {
	wio.ByteOrder = binary.LittleEndian

	target := "/mnt/Program Files/Steam/steamapps/common/Helldivers 2/data/e75f556a740e00c9"

	ExtractSoundBank(nil, target)
}
