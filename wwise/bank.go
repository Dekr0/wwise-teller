package wwise

type Chunk struct {
	handle  uint32
	idx     uint8
}

type Bank struct {
	Chunks  []Chunk
}
