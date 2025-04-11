package reader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var BYTE_ORDER = binary.LittleEndian

/*
TYPE_SID: ShortID (uint32_t)
TYPE_TID: target ShortID (uint32_t) same thing but for easier understanding of output
TYPE_UNI: union (float / int32_t)
TYPE_D64: double
TYPE_F32: float
TYPE_4CC: FourCC
TYPE_S64: int64_t
TYPE_U64: uint64_t
TYPE_S32: int32_t
TYPE_U32: uint32_t
TYPE_S16: int16_t
TYPE_U16: uint16_t
TYPE_S8 : int8_t
TYPE_U8 : uint8_t
TYPE_VAR: variable size #u8/u16/u32
TYPE_GAP: byte gap
TYPE_STR: string
TYPE_STZ: string (null-terminated)
*/

/*
A wrapper around io.ReadSeeker for reading binary data in different section of
a sound bank, or the entire section of a sound bank.
It maintains byte order consistency and the absolute location (from origin, 0
byte).
*/
type BankReader struct {
	cursorPos uint64
	reader io.ReadSeeker
	byteOrder binary.ByteOrder
}

func NewSoundbankReader(f io.ReadSeeker, o binary.ByteOrder) *BankReader {
	return &BankReader{0, f, o}
}

func (s *BankReader) GetByteOrder() binary.ByteOrder {
	return s.byteOrder
}

func (s *BankReader) Tell() uint64 {
	return s.cursorPos
}

func (s *BankReader) AbsSeek(pos uint64) error {
	absPos, err := s.reader.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return err
	}
	s.cursorPos = uint64(absPos)
	return nil
}

func (s *BankReader) RelSeek(pos int64) error {
	absPos, err := s.reader.Seek(pos, io.SeekCurrent)
	if err != nil {
		return err
	}
	s.cursorPos = uint64(absPos)
	return nil
}

/*
Read `size` amount of bytes and create a new bank reader. This new bank reader 
will operate on byte slice reader that support read and seek.
*/
func (s *BankReader) NewSectionBankReader(size uint64) (*BankReader, error) {
	section := make([]byte, size)
	nread, err := s.reader.Read(section)
	if err != nil {
		return nil, err
	}
	s.cursorPos += uint64(nread)
	return &BankReader{ 0, bytes.NewReader(section), s.byteOrder }, nil
}

func (s *BankReader) ReadFullUnsafe(size uint64, reserved uint64) []byte {
	b, err := s.ReadFull(size, reserved)
	if err != nil {
		panic(b)
	}
	return b
}

/*
Reserved provide an minor optimization on encoding field that require header 
field attached at the beginning. Use this when reader is operating on data that 
is not yet in the memory. 

TODO: Create a different version of bank reader to avoid situation when reader 
is operating in-memory data, especially reading large amount of data since it 
will cause copy.
*/
func (s *BankReader) ReadFull(size uint64, reserved uint64) ([]byte, error) {
	b := make([]byte, size, size + reserved)
	nread, err := s.reader.Read(b)
	if err != nil {
		return nil, err
	}
	if uint64(nread) > size {
		err := fmt.Errorf(
			"BankReader.ReadFull: io.Reader read more than the specified size." +
			" Expected read size %d, actual read size %d.", size, nread,
		)
		return nil, err
	}
	s.cursorPos += size
	return b, nil
}

/*
Read data from the current point in time until it encounter an error / EOF or 
successfully return all data.
*/
func (s *BankReader) ReadAll() ([]byte, error) {
	b, err := io.ReadAll(s.reader)
	if err != nil {
		return nil, err
	}
	s.cursorPos += uint64(len(b))
	return b, nil
}

func (s *BankReader) U8Unsafe() uint8 {
	v, err := s.U8()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) U8() (uint8, error) {
	var v uint8
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 1
	return v, err
}


func (s *BankReader) I8Unsafe() int8 {
	v, err := s.I8()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) I8() (int8, error) {
	var v int8
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 1
	return v, err
}

func (s *BankReader) U16Unsafe() uint16 {
	v, err := s.U16()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) U16() (uint16, error) {
	var v uint16
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 2
	return v, err
}

func (s *BankReader) I16Unsafe() int16 {
	v, err := s.I16()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) I16() (int16, error) {
	var v int16
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 2
	return v, err
}

func (s *BankReader) U32Unsafe() uint32 {
	v, err := s.U32()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) U32() (uint32, error) {
	var v uint32
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 4
	return v, err
}

func (s *BankReader) I32Unsafe() int32 {
	v, err := s.I32()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) I32() (int32, error) {
	var v int32
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 4
	return v, err
}

func (s *BankReader) F32Unsafe() float32 {
	v, err := s.F32()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) F32() (float32, error) {
	var v float32
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 4
	return v, err
}

func (s *BankReader) FourCC() ([]byte, error) {
	buf := make([]byte, 4, 4)
	err := binary.Read(s.reader, s.byteOrder, &buf)
	s.cursorPos += 4
	return buf, err
}

func (s *BankReader) U64Unsafe() uint64 {
	v, err := s.U64()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) U64() (uint64, error) {
	var v uint64
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 8
	return v, err
}

func (s *BankReader) I64Unsafe() int64 {
	v, err := s.I64()
	if err != nil {
		panic(err)
	}
	return v
}

func (s *BankReader) I64() (int64, error) {
	var v int64
	err := binary.Read(s.reader, s.byteOrder, &v)
	s.cursorPos += 8
	return v, err
}
