package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

func GetDirAndBank(p string, dirOnly bool) ([]os.DirEntry, error) {
	fd, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	newEntries := make([]os.DirEntry, 0, 1024)
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
			if entry.IsDir() {
				newEntries = append(newEntries, entry)
				continue
			}
			if filepath.Ext(entry.Name()) == ".bnk" && !dirOnly {
				newEntries = append(newEntries, entry)
				continue
			}
		}

		bound -= 1
	}

	if bound <= 0 {
		return nil, fmt.Errorf(
			"Failed to read files from %s: upper bound is reached", p,
		)
	}
	
	return newEntries, nil
}

