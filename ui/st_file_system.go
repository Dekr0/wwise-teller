package ui

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Dekr0/wwise-teller/utils"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type RankDirEntry struct {
	rank int
	entry os.DirEntry
}

type FileSystem struct {
	Pwd      string
	// Filter
	dirOnly  bool
	query    string
	exts     []string
	// Cache
	entries  []os.DirEntry 
	filtered []*RankDirEntry
}

func newFileSystem(
	initialDir string,
	dirOnly    bool,
	exts       []string,
) (*FileSystem, error) {
	entries := make([]os.DirEntry, 0, 1024)
	entries, err := utils.GetDirAndFiles(initialDir, entries)
	if err != nil {
		return nil, err
	}
	filtered := make([]*RankDirEntry, len(entries))
	for i, e := range entries {
		filtered[i] = &RankDirEntry{-1, e}
	}
	fs := &FileSystem{
		Pwd     : initialDir,
		dirOnly : dirOnly,
		query   : "",
		exts    : exts,
		entries : entries,
		filtered: filtered,
	}
	fs.filter()
	return fs, nil
}

func (f *FileSystem) assert() {
	if f.dirOnly && f.exts != nil {
		panic("Enable directory only but extensions filter are provided")
	}
}

func (f *FileSystem) filter() {
	f.assert()
	i := 0
	old := len(f.filtered)
	for _, e := range f.entries {
		if f.dirOnly {
			if !e.IsDir() { continue }
			rank := fuzzy.RankMatch(f.query, e.Name())
			if rank == -1 { continue }
			if i < len(f.filtered) {
				f.filtered[i].rank = rank
				f.filtered[i].entry = e
			} else {
				f.filtered = append(f.filtered, &RankDirEntry{rank, e})
			}
			i += 1
		} else {
			if !e.IsDir() && len(f.exts) > 0 && !slices.ContainsFunc(
				f.exts, 
				func(ext string) bool {
					return strings.Compare(ext, filepath.Ext(e.Name())) == 0
				},
			) {
				continue
			}
			rank := fuzzy.RankMatch(f.query, e.Name())
			if rank == -1 { continue }
			if i < len(f.filtered) {
				f.filtered[i].rank = rank
				f.filtered[i].entry = e
			} else {
				f.filtered = append(f.filtered, &RankDirEntry{rank, e})
			}
			i += 1
		}
	}
	if i < old {
		f.filtered = slices.Delete(f.filtered, i, old)
	}
	slices.SortFunc(f.filtered, func(a *RankDirEntry, b *RankDirEntry) int {
		if a.rank < b.rank {
			return -1
		}
		if a.rank == b.rank {
			return 0
		}
		return 1
	})
}

func (f *FileSystem) cdParent() error {
	pwd := filepath.Dir(f.Pwd)
	if runtime.GOOS == "windows" && pwd == "." {
		return nil
	}
	if pwd == f.Pwd {
		return nil
	}
	var err error
	entries, err := utils.GetDirAndFiles(pwd, f.entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.entries = entries
	}
	f.Pwd = pwd
	f.clearFilter()
	return nil
}

func (f *FileSystem) cd(basename string) error {
	pwd := filepath.Join(f.Pwd, basename)
	var err error
	entries, err := utils.GetDirAndFiles(pwd, f.entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.entries = entries
	}
	f.Pwd = pwd
	f.clearFilter()
	return nil
}

func (f *FileSystem) clearFilter() {
	f.query = ""
	f.filter()
}
