package utils

import (
	"hash/fnv"

	"github.com/google/uuid"
)

const FNV30Mask = (1 << 30) - 1

func FNV32(data []byte) (uint32, error) {
	f := fnv.New32()
	_, err := f.Write(data)
	if err != nil {
		return 0, err
	}
	return f.Sum32(), nil
}

func FNV30(data []byte) (uint32, error) {
	h, err := FNV32(data)
	if err != nil {
		return 0, err
	}
	return (h >> 30) ^ (h & FNV30Mask), nil
}

func ShortID() (uint32, error) {
	guid := uuid.New()
	data, err := guid.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return FNV30(data)
}
