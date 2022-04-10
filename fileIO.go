package golanglibs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type fileIOStruct struct {
	path   string
	fd     *os.File
	mode   string
	reader *bufio.Reader
	lock   *LockStruct
}

func (m *fileIOStruct) Readlines() chan *StringStruct {
	if m.reader == nil {
		m.reader = bufio.NewReader(m.fd)
	}

	lines := make(chan *StringStruct)

	go func() {
		m.lock.Acquire()
		defer m.lock.Release()
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
			lines <- String(line).Strip("\r\n")
		}
	}()

	return lines
}

func (m *fileIOStruct) Readline() *StringStruct {
	m.lock.Acquire()
	defer m.lock.Release()

	b := make([]byte, 1)

	line := ""

	for {
		_, err := io.ReadAtLeast(m.fd, b, 1)
		if err != nil {
			if len(line) != 0 {
				return String(line)
			}
			Panicerr(err)
		}
		bs := string(b)
		if bs == "\n" {
			return String(line)
		}
		line = line + bs
	}
}

func (m *fileIOStruct) Close() {
	m.fd.Close()
}

// str can be string, *StringStruct, []byte
func (m *fileIOStruct) Write(str interface{}) *fileIOStruct {
	m.lock.Acquire()
	defer m.lock.Release()
	var buf []byte
	switch vv := str.(type) {
	case string:
		buf = []byte(vv)
	case *StringStruct:
		buf = []byte(vv.S)
	case []byte:
		buf = vv
	}
	m.fd.Write(buf)
	return m
}

func (m *fileIOStruct) Read(num ...int) *StringStruct {
	m.lock.Acquire()
	defer m.lock.Release()

	var bytes []byte
	var err error
	if len(num) == 0 {
		bytes, err = ioutil.ReadAll(m.fd)
		Panicerr(err)
		m.Close()
	} else {
		buffer := make([]byte, num[0])
		i, err := io.ReadAtLeast(m.fd, buffer, num[0])
		if err != nil {
			if !String("EOF").In(err.Error()) {
				Panicerr(err)
			}
		}
		bytes = buffer[:i]
	}

	return String(string(bytes))
}

func (m *fileIOStruct) Seek(num int64) {
	_, err := m.fd.Seek(num, 0)
	Panicerr(err)
}

func Open(args ...string) *fileIOStruct {
	path := args[0]
	var mode string
	if len(args) != 1 {
		mode = args[1]
	} else {
		mode = "r"
	}
	var fd *os.File
	var err error
	if string(mode[0]) == "r" {
		fd, err = os.Open(path)
	}
	if string(mode[0]) == "a" {
		fd, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if string(mode[0]) == "w" {
		fd, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	}
	Panicerr(err)
	return &fileIOStruct{
		path: path,
		fd:   fd,
		mode: mode,
		lock: getLock(),
	}
}

// func getStdin() *fileIOStruct {
// 	return &fileIOStruct{fd: os.Stdin, mode: "r"}
// }
