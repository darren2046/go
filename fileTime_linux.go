package golanglibs

import "os"

func (f *fileStruct) Time() *fileTimeStruct {
	fi, err := os.Stat(f.filePath)
	Panicerr(err)
	mtime := fi.ModTime().Unix()
	stat := fi.Sys().(*syscall.Stat_t)
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)).Unix()
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec)).Unix()

	return &fileTimeStruct{
		Mtime: mtime,
		ctime: ctime,
		atime: atime,
	}
}
