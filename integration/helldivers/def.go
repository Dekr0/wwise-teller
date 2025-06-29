package helldivers

import (
	"errors"
)

var NotHelldiversGameArchive error = errors.New(
	"Not a game archive used by Helldivers 2",
)
var NotMETA error = errors.New(
	"Not META data section",
)

const (
	AssetTypeSoundBank       = 6006249203084351385
	AssetTypeWwiseDependency = 12624162998411505776
	AssetTypeWwiseStream     = 5785811756662211598
)

const MagicValue uint32 = 0xF0000011

type IntegrationType uint8

const (
	IntegrationTypeHelldivers2 IntegrationType = 2
)

type Asset struct {
	Header      *AssetHeader
	Data        []byte
	StreamData  []byte
	GPURsrcData []byte
	META        []byte
}

type AssetHeader struct {
	FileID        uint64 `json:"fileID"`
	TypeID        uint64 `json:"typeID"`
	DataOffset    uint64 `json:"dataOffset"`
	StreamOffset  uint64 `json:"streamOffset"`
	GPURsrcOffset uint64 `json:"gPURsrcOffset"`
	UnknownU64A   uint64 `json:"unknownU64A"`
	UnknownU64B   uint64 `json:"unknownU64B"`
	DataSize      uint32 `json:"dataSize"`
	StreamSize    uint32 `json:"streamSize"`
	GPURsrcSize   uint32 `json:"gPURsrcSize"`
	UnknownU32A   uint32 `json:"unknownU32A"`
	UnknownU32B   uint32 `json:"unknownU32B"`
	Idx           uint32 `json:"idx"`
}
