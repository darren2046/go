package golanglibs

import (
	"io"
	"os"
)

type fileStruct struct {
	filePath string
}

func File(filePath string) *fileStruct {
	return &fileStruct{filePath: filePath}
}

type fileTimeStruct struct {
	Ctime int64
	Mtime int64
	Atime int64
}

// Get file details
func (f *fileStruct) Stat() os.FileInfo {
	ff, err := os.Stat(f.filePath)
	Panicerr(err)
	return ff
}

// Get file size
func (f *fileStruct) Size() int64 {
	info, err := os.Stat(f.filePath)
	Panicerr(err)
	return info.Size()
}

// Touch a file like touch command
func (f *fileStruct) Touch() {
	fd, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	Panicerr(err)
	fd.Close()
}

func (f *fileStruct) Chmod(mode uint32) bool {
	if len(Str(mode)) == 3 {
		mode = Uint32("0" + Str(mode))
	}
	return os.Chmod(f.filePath, os.FileMode(mode)) == nil
}

func (f *fileStruct) Chown(uid, gid int) bool {
	return os.Chown(f.filePath, uid, gid) == nil
}

func (f *fileStruct) Mtime() int64 {
	fd, err := os.Open(f.filePath)
	Panicerr(err)
	defer fd.Close()
	fileinfo, err := fd.Stat()
	Panicerr(err)
	return fileinfo.ModTime().Unix()
}

func (f *fileStruct) Unlink() {
	err := os.RemoveAll(f.filePath)
	Panicerr(err)
}

func (f *fileStruct) Copy(dest string) {
	fd1, err := os.Open(f.filePath)
	Panicerr(err)
	defer fd1.Close()
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	Panicerr(err)
	defer fd2.Close()
	_, err = io.Copy(fd2, fd1)
	Panicerr(err)
}

func (f *fileStruct) Move(newPosition string) {
	err := os.Rename(f.filePath, newPosition)
	Panicerr(err)
	f.filePath = newPosition
}
