package golanglibs

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/rakyll/statik/fs"
)

type statikFileStruct struct {
	path   string
	fd     http.File
	mode   string
	reader *bufio.Reader
}

func (m *statikFileStruct) Readlines() chan string {
	if m.reader == nil {
		m.reader = bufio.NewReader(m.fd)
	}

	lines := make(chan string)

	go func() {
		for {
			line, err := m.reader.ReadString('\n')
			if len(line) == 0 {
				if err != nil {
					close(lines)
					m.Close()
					if err == io.EOF {
						return
					}
					_, fn, line, _ := runtime.Caller(0)
					panic(filepath.Base(fn) + ":" + strconv.Itoa(line-7) + " >> " + err.Error())
				}
			}
			line = line[:len(line)-1]
			lines <- line
		}
	}()

	return lines
}

func (m *statikFileStruct) Readline() string {
	b := make([]byte, 1)

	line := ""

	for {
		_, err := io.ReadAtLeast(m.fd, b, 1)
		if err != nil {
			if len(line) != 0 {
				return line
			}
			Panicerr(err)
		}
		bs := string(b)
		if bs == "\n" {
			return line
		}
		line = line + bs
	}
}

func (m *statikFileStruct) Close() {
	m.fd.Close()
}

func (m *statikFileStruct) Read(num ...int) string {
	var bytes []byte
	var err error
	if len(num) == 0 {
		bytes, err = ioutil.ReadAll(m.fd)
		Panicerr(err)
		m.Close()
	} else {
		buffer := make([]byte, num[0])
		_, err := io.ReadAtLeast(m.fd, buffer, num[0])
		Panicerr(err)
		bytes = buffer
	}

	return string(bytes)
}

func (m *statikFileStruct) Seek(num int64) {
	_, err := m.fd.Seek(num, 0)
	Panicerr(err)
}

func statikOpen(path string) *statikFileStruct {
	statikFS, err := fs.New()
	Panicerr(err)

	if !String(path).StartsWith("/") {
		path = "/" + path
	}

	fd, err := statikFS.Open(path)
	Panicerr(err)

	return &statikFileStruct{
		path: path,
		fd:   fd,
		mode: "r",
	}
}
