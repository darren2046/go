package golanglibs

import (
	"github.com/hpcloud/tail"
)

func tailf(path string, startFromEndOfFile ...bool) chan *tail.Line {
	if !pathExists(path) {
		Panicerr("Cannot open file \"" + path + "\": no such file or directory")
	}
	var t *tail.Tail
	var err error

	if len(startFromEndOfFile) == 0 {
		t, err = tail.TailFile(path, tail.Config{Follow: true, Poll: true, Logger: tail.DiscardingLogger})
	} else {
		if startFromEndOfFile[0] == false {
			t, err = tail.TailFile(path, tail.Config{Follow: true, Poll: true, Logger: tail.DiscardingLogger})
		} else {
			t, err = tail.TailFile(path, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, Poll: true, Logger: tail.DiscardingLogger})
		}
	}

	Panicerr(err)

	return t.Lines
}
