package helldivers

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"

	"github.com/Dekr0/wwise-teller/utils"
	"github.com/Dekr0/wwise-teller/wio"
	"github.com/Dekr0/wwise-teller/wwise"
)

var NotHelldiversIntegration = errors.New("Not a Helldivers integration")

const SizeOfMetaStable = 12
const SizeOfMetaHeader = wwise.SizeOfChunkHeader + 1
const SizeOfMetaBufferInit = SizeOfMetaStable + SizeOfMetaHeader
type METAStable struct {
	FileID   uint64
	XOR      [4]byte
	WwiseDep []byte
}

func GenHelldiversPatchStable(bnk []byte, meta []byte, path string) error {
	m := METAStable{}
	err := ParseMETAStable(&m, wio.NewInPlaceReader(meta, wio.ByteOrder))
	if err != nil {
		return err
	}

	bnk[0x08] = m.XOR[0]
	bnk[0x09] = m.XOR[1]
	bnk[0x0A] = m.XOR[2]
	bnk[0x0B] = m.XOR[3]

	f, err := os.OpenFile(
		filepath.Join(path, "9ba626afa44a3aa3.patch_0"),
		IOCW,
		0666,
	)
	defer f.Close()
	if err != nil {
		return err
	}

	w := wio.NewBinaryWriteHelper(f)
	if err := WritePatchHeaderStable(w); err != nil {
		return err
	}
	if err := WritePatchTypeHeaderStable(w); err != nil {
		return err
	}

	dataOffset := uint64(w.Tell()) + 160 + 8

	chunks := [][]byte{}

	bnkAssetHeader := AssetHeader{
		FileID: m.FileID,
		TypeID: AssetTypeSoundBank,
		DataOffset: dataOffset,
		StreamOffset: 0,
		GPURsrcOffset: 0,
		UnknownU64A: 0,
		UnknownU64B: 0,
		DataSize: 16 + uint32(len(bnk)),
		StreamSize: 0,
		GPURsrcSize: 0,
		UnknownU32A: 0,
		UnknownU32B: 0,
		Idx: 0,
	}
	buffer := make([]byte, 80, 80)
	binary.Encode(buffer, wio.ByteOrder, bnkAssetHeader)
	if err := w.Bytes(buffer); err != nil {
		return err
	}

	bnkChunk := []byte{216,47,118,120}
	bnkChunk, _ = binary.Append(bnkChunk, wio.ByteOrder, uint32(len(bnk)))
	bnkChunk, _ = binary.Append(bnkChunk, wio.ByteOrder, m.FileID)
	bnkChunk = append(bnkChunk, utils.Pad16ByteAlign(bnk)...)
	chunks = append(chunks, bnkChunk)
	dataOffset += uint64(len(bnkChunk))

	depAssetHeader := AssetHeader{
		FileID: m.FileID,
		TypeID: AssetTypeWwiseDependency,
		DataOffset: dataOffset,
		StreamOffset: 0,
		GPURsrcOffset: 0,
		UnknownU64A: 0,
		UnknownU64B: 0,
		DataSize: uint32(len(m.WwiseDep)),
		StreamSize: 0,
		GPURsrcSize: 0,
		UnknownU32A: 0,
		UnknownU32B: 0,
		Idx: 1,
	}
	binary.Encode(buffer, wio.ByteOrder, depAssetHeader)
	if err := w.Bytes(buffer); err != nil {
		return err
	}

	chunks = append(chunks, utils.Pad16ByteAlign(m.WwiseDep))

	pad := make([]byte, 8, 8)
	if err := w.Bytes(pad); err != nil {
		return err
	}

	for _, c := range chunks {
		if err := w.Bytes(c); err != nil {
			return err
		}
	}

	return f.Close()
}

func ParseMETAStable(m *METAStable, r *wio.InPlaceReader) error {
	itype, err := r.U8()
	if err != nil {
		return err
	}
	if itype != uint8(IntegrationTypeHelldivers2) {
		return NotHelldiversIntegration
	}

	m.FileID, err = r.U64()
	if err != nil {
		return err
	}

	fourCC, err := r.FourCCNoCopy()
	if err != nil {
		return err
	}
	m.XOR = [4]byte(fourCC)

	m.WwiseDep, err = r.ReadAllNoCopy()
	if err != nil {
		return err
	}

	return nil
}

func WritePatchHeaderStable(w *wio.BinaryWriteHelper) error {
	if err := w.U32(MagicValue); err != nil {
		return err
	}
	// Number of file types = 2 (Sound bank and dependency)
	if err := w.U32(2); err != nil {
		return err
	}
	// Number of files = 2 x number of sound banks
	if err := w.U32(2); err != nil {
		return err
	}
	// Unknown 32 bits
	if err := w.U32(0); err != nil {
		return err
	}
	// Require (Reverse engineering black magic)
	Unk4Data := []byte{206,9,245,244,0,0,0,0,12,114,159,158,136,114,184,189,0,160,107,2,0,0,0,0,0,121,81,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	if err := w.Bytes(Unk4Data); err != nil {
		return err
	}
	return nil
}

func WritePatchTypeHeaderStable(w *wio.BinaryWriteHelper) error {
	// Sound bank type header
	if err := w.U64(0); err != nil {
		return err
	}
	if err := w.U64(AssetTypeSoundBank); err != nil {
		return err
	}
	// Number of sound banks
	if err := w.U64(1); err != nil {
		return err
	}
	// Alignment? 8 x 3 + 2 x 4 = 24 + 8 = 32
	if err := w.U32(16); err != nil {
		return err
	}
	if err := w.U32(64); err != nil {
		return err
	}

	// Wwise dependency type header
	if err := w.U64(0); err != nil {
		return err
	}
	if err := w.U64(AssetTypeWwiseDependency); err != nil {
		return err
	}
	// Number of wwise dependencies
	if err := w.U64(1); err != nil {
		return err
	}
	// Alignment? 8 x 3 + 2 x 4 = 24 + 8 = 32
	if err := w.U32(16); err != nil {
		return err
	}
	if err := w.U32(64); err != nil {
		return err
	}
	return nil
}
