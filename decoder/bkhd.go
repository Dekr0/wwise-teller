package decoder

import (
	"io"
	"sort"

	def "github.com/Dekr0/unwise/decoder/definition"
	uio "github.com/Dekr0/unwise/io"
	"github.com/Dekr0/unwise/wwise"
)

func DecodeBKHD(
	path     string,
	inReader io.Reader,
	o        order) (
	b *wwise.BKHD, err error,
) {
	size, err := uio.U32(inReader, o)
	if err != nil {
		return nil, err
	}

	r := io.LimitReader(inReader, int64(size))

	b = &wwise.BKHD{}

	version, err := uio.U32(r, o)

	{
		if err != nil {
			return nil, err
		}

		if version == 0 || version == 1 {
			return nil, LegacyBank(path, version)
		}

		_, in := sort.Find(len(def.CustomVersions), func(i int) int {
			if version < def.CustomVersions[i] {
				return -1
			} else if version == def.CustomVersions[i] {
				return 0
			} else {
				return 1
			}
		})
		if in {
			return nil, CustomBank(path, version)
		}

		if version & MaskCustomVersionHigh == MaskCustomVersionResult {
			return nil, UnknownCustomBank(path, version)
		}
		if version & MaskEntryption > 0 {
			return nil, EncryptionBank(path)
		}

		_, in = sort.Find(len(def.Versions), func(i int) int {
			if version < def.Versions[i] {
				return -1
			} else if version == def.Versions[i] {
				return 0
			} else {
				return 1
			}
		})
		if !in {
			return nil, UnsupportKnownBank(path, version)
		}
	}

	b.Version = version

	b.Id, err = uio.U32(r, o)
	if err != nil {
		return nil, err
	}

	b.Language, err = uio.U32(r, o)
	if err != nil {
		return nil, err
	}

	res, err := uio.U32(r, o)
	if err != nil {
		return nil, err
	}

	b.DeviceAllocated = u16(0x0000ffff & res)
	b.Alignment = u16(res >> 16)

	b.Project, err = uio.U32(r, o)
	if err != nil {
		return nil, err
	}

	b.Data, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	
	return b, nil
}
