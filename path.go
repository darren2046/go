package golanglibs

import (
	"os"
	"path"
	"path/filepath"
)

type pathStruct struct {
	Exists    func(path string) bool
	IsFile    func(path string) bool
	IsDir     func(path string) bool
	Basename  func(path string) string
	Basedir   func(path string) string
	Dirname   func(path string) string
	Join      func(args ...string) string
	Abspath   func(path string) string
	IsSymlink func(path string) bool
}

func pathIsSymlink(path string) bool {
	fi, err := os.Lstat(path)
	Panicerr(err)
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	} else {
		return false
	}
}

func abspath(path string) string {
	str, err := filepath.Abs(path)
	Panicerr(err)
	return str
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func pathIsFile(path string) bool {
	fd, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	fm := fd.Mode()
	return fm.IsRegular()
}

func pathIsDir(path string) bool {
	fd, err := os.Stat(path)
	Panicerr(err)
	fm := fd.Mode()
	return fm.IsDir()
}

func pathBasename(path string) string {
	return filepath.Base(path)
}

func pathBasedir(path string) string {
	str, err := filepath.Abs(filepath.Dir(path))
	Panicerr(err)
	return str
}

func pathDirname(path string) string {
	return filepath.Dir(path)
}

func pathJoin(args ...string) string {
	return path.Join(args...)
}
