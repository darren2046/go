package golanglibs

import "os"

func (f *fileStruct) Time() *fileTimeStruct {
	fi, err := os.Stat(f.filePath)
	Panicerr(err)
	mtime := fi.ModTime().Unix()

	return &fileTimeStruct{
		Mtime: mtime,
	}
}
