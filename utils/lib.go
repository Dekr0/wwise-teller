package utils

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/shirou/gopsutil/v4/disk"
)

const MaxInt = int32(^uint32(0) >> 1)

var Vols []string = []string{}

func GetHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	home, err = filepath.Abs(home)
	if err != nil {
		return "", err
	}
	return home, err
}

// Write bytes to a file. Create if it doesn't exist. Otherwise, keep increase 
// post fix number counter until it exceeds maximum limit.
func SaveFileWithRetry(b []byte, path string) error {
	stat, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(path, b, 0666); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	if stat.IsDir() {
		return fmt.Errorf(
			"Failed to save file as %s: target destination is a directory",
			path,
		)
	}

	for i := range 128 {
		postFixPath := fmt.Sprintf("%s_%d", path, i)
		stat, err = os.Lstat(postFixPath)
		if err != nil {
			if os.IsNotExist(err) {
				if err := os.WriteFile(path, b, 0666); err != nil {
					return err
				}
				return nil
			}
			return err
		}
	}

	return errors.New("Maximum save file retry exceeds")
}

// Expect a pre-allocated buffer slice for memory reused. Panic if buffer is nil
func GetDirAndFiles(p string, buffer []os.DirEntry) (
	[]os.DirEntry, error,
) {
	if buffer == nil {
		panic("GetDirAndFiles expect a pre-allocated buffer slice")
	}

	fd, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	buffer = slices.Delete(buffer, 0, len(buffer))
	bound := 128
	for bound > 0 {
		entries, err := fd.ReadDir(1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, entry := range entries {
			buffer = append(buffer, entry)
			continue
		}

		bound -= 1
	}

	if bound <= 0 {
		return nil, fmt.Errorf(
			"Failed to read files from %s: upper bound is reached", p,
		)
	}
	
	return buffer, nil
}

func IsDigit(s string) bool {
	for _, c := range s {
		if c < '0' && c > '9' {
			return false
		}
	}
	return true
}

func Pad16ByteAlign(data []byte) []byte {
	pad := (int(math.Ceil(float64(len(data)) / float64(16))) * 16) - len(data)
	return append(data, make([]byte, pad, pad)...)
}

func ScanMountPoint() error {
	if runtime.GOOS != "windows" {
		return nil
	}
	stats, err := disk.Partitions(true)
	if err != nil {
		return err
	}
	Vols = []string{}
	for _, stat := range stats {
		Vols = append(Vols, stat.Mountpoint)
	}
	return nil
}
