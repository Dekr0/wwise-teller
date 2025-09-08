package decoder

import (
	"fmt"

	"github.com/Dekr0/unwise/wwise"
)

// Has side effect
// Thread safe
func DecodeDATA(d *wwise.DIDXDATA, chunk []byte) {
	wwise.LockDIDXDATA(d)
	defer wwise.UnlockDIDXDATA(d)

	sourceIds := d.SourceIds
	offsets := d.Offsets
	sizes := d.Sizes

	audioData := make(map[u32][]byte, len(sourceIds))

	offsetNotAligned := u32(0)
	for _, sourceId := range sourceIds {
		offset, in := offsets[sourceId]
		if !in {
			panic(fmt.Sprintf("%d has meida index but it has no offset value", sourceId))
		}

		size, in := sizes[sourceId]
		if !in {
			panic(fmt.Sprintf("%d has meida index but it has no size value", sourceId))
		}

		if _, in := audioData[sourceId]; in {
			panic(fmt.Sprintf("%d has audio data before storing new audio data", sourceId))
		}

		audioData[sourceId] = chunk[offset:offset + size]
		
		offsets[sourceId] = offsetNotAligned

		offsetNotAligned += size
	}
}
