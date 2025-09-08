package wwise

import (
	"fmt"
	"sync"
)

type DATA struct {
	mu sync.Mutex

	AudioData [][]byte
	Indexing  map[u32][]byte
}

func NewDATA(size u32) *DATA {
	return &DATA{
		AudioData: make([][]byte, 0, size),
		Indexing: make(map[u32][]byte, size),
	}
}

// Has side effect
// Thread safe
// Use this if assuming there will be no duplicate in DIDX entry (e.g., decoding 
// phase)
func NewAudioData(d *DATA, sourceId u32, data []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.AudioData = append(d.AudioData, data)
	d.Indexing[sourceId] = data
}

func NewAudioDataCheck(d *DATA, sourceId u32, data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, in := d.Indexing[sourceId]; in {
		return fmt.Errorf("Duplicate audio data for source id %d", sourceId)
	}

	d.AudioData = append(d.AudioData, data)
	d.Indexing[sourceId] = data

	return nil
}
