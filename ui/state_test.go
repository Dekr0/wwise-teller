package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExplorer(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	f, err := newFileSystem(home, false, []string{".bnk", ".st_bnk"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n")
	t.Log(f.Pwd)
	for _, e := range f.entries {
		t.Log(e)
	}

	err = f.cdParent()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n")
	t.Log(f.Pwd)
	for _, e := range f.entries {
		t.Log(e)
	}

	err = f.cdParent()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n")
	t.Log(f.Pwd)
	for _, e := range f.entries {
		t.Log(e)
	}

	err = f.cdParent()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n")
	t.Log(f.Pwd)
	for _, e := range f.entries {
		t.Log(e)
	}
}

func TestCdParentWindows(t *testing.T) {
	pwd := "D:/a/b"
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
}

func TestCdParentUnix(t *testing.T) {
	pwd := "/home/user/.config"
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
	pwd = filepath.Dir(pwd)
	t.Log(pwd)
}
