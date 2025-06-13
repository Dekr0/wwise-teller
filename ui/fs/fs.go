package fs

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
	Rank int
	Entry os.DirEntry
}

type FileSystem struct {
	Pwd      string
	// Filter
	DirOnly  bool
	Query    string
	Exts     []string
	// Cache
	Entries  []os.DirEntry 
	Filtered []*RankDirEntry
}

func NewFileSystem(
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
		DirOnly : dirOnly,
		Query   : "",
		Exts    : exts,
		Entries : entries,
		Filtered: filtered,
	}
	fs.Filter()
	return fs, nil
}

func (f *FileSystem) Assert() {
	if f.DirOnly && f.Exts != nil {
		panic("Enable directory only but extensions filter are provided")
	}
}

func (f *FileSystem) CD(basename string) error {
	pwd := filepath.Join(f.Pwd, basename)
	var err error
	entries, err := utils.GetDirAndFiles(pwd, f.Entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.Entries = entries
	}
	f.Pwd = pwd
	f.ClearFilter()
	return nil
}

func (f *FileSystem) Parent() error {
	pwd := filepath.Dir(f.Pwd)
	if runtime.GOOS == "windows" && pwd == "." {
		return nil
	}
	if pwd == f.Pwd {
		return nil
	}
	var err error
	entries, err := utils.GetDirAndFiles(pwd, f.Entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.Entries = entries
	}
	f.Pwd = pwd
	f.ClearFilter()
	return nil
}

func (f *FileSystem) ClearFilter() {
	f.Query = ""
	f.Filter()
}

func (f *FileSystem) Filter() {
	f.Assert()
	i := 0
	old := len(f.Filtered)
	for _, e := range f.Entries {
		if f.DirOnly {
			if !e.IsDir() { continue }
			rank := fuzzy.RankMatch(f.Query, e.Name())
			if rank == -1 { continue }
			if i < len(f.Filtered) {
				f.Filtered[i].Rank = rank
				f.Filtered[i].Entry = e
			} else {
				f.Filtered = append(f.Filtered, &RankDirEntry{rank, e})
			}
			i += 1
		} else {
			if !e.IsDir() && len(f.Exts) > 0 && !slices.ContainsFunc(
				f.Exts, 
				func(ext string) bool {
					return strings.Compare(ext, filepath.Ext(e.Name())) == 0
				},
			) {
				continue
			}
			rank := fuzzy.RankMatch(f.Query, e.Name())
			if rank == -1 { continue }
			if i < len(f.Filtered) {
				f.Filtered[i].Rank = rank
				f.Filtered[i].Entry = e
			} else {
				f.Filtered = append(f.Filtered, &RankDirEntry{rank, e})
			}
			i += 1
		}
	}
	if i < old {
		f.Filtered = slices.Delete(f.Filtered, i, old)
	}
	slices.SortFunc(f.Filtered, func(a *RankDirEntry, b *RankDirEntry) int {
		if a.Rank < b.Rank {
			return -1
		}
		if a.Rank == b.Rank {
			return 0
		}
		return 1
	})
}

func (f *FileSystem) Refresh() error {
	entries, err := utils.GetDirAndFiles(f.Pwd, f.Entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.Entries = entries
	}
	f.ClearFilter()
	return nil
}

func (f *FileSystem) SwitchVolume(vol string) error {
	vol = vol + "/"
	if runtime.GOOS != "windows" {
		panic("FileSystem.switchVol must be called only on Windows.")
	}
	entries, err := utils.GetDirAndFiles(vol, f.Entries)
	if err != nil {
		return err
	}
	if entries != nil {
		f.Entries = entries
	}
	f.Pwd = vol
	f.ClearFilter()
	return nil
}

func (f *FileSystem) Vol() string {
	if runtime.GOOS != "windows" {
		panic("Current OS is not Windows")
	}
	return filepath.VolumeName(f.Pwd)
}
