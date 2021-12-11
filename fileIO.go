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
	lock   *lockStruct
}

func (m *fileIOStruct) readlines() chan string {
	if m.reader == nil {
		m.reader = bufio.NewReader(m.fd)
	}

	lines := make(chan string)

	go func() {
		m.lock.acquire()
		defer m.lock.release()
		for {
			line, err := m.reader.ReadString('\n')
			if len(line) == 0 {
				if err != nil {
					close(lines)
					m.close()
					if err == io.EOF {
						return
					}
					_, fn, line, _ := runtime.Caller(0)
					panic(filepath.Base(fn) + ":" + strconv.Itoa(line-7) + " >> " + err.Error())
				}
			}
			line = String(line).Strip("\r\n").Get()
			lines <- line
		}
	}()

	return lines
}

func (m *fileIOStruct) readline() string {
	m.lock.acquire()
	defer m.lock.release()

	b := make([]byte, 1)

	line := ""

	for {
		_, err := io.ReadAtLeast(m.fd, b, 1)
		if err != nil {
			if len(line) != 0 {
				return line
			}
			panicerr(err)
		}
		bs := string(b)
		if bs == "\n" {
			return line
		}
		line = line + bs
	}
}

func (m *fileIOStruct) close() {
	m.fd.Close()
}

func (m *fileIOStruct) write(str interface{}) *fileIOStruct {
	m.lock.acquire()
	defer m.lock.release()
	var buf []byte
	if Typeof(str) == "string" {
		s := str.(string)
		buf = []byte(s)
	} else {
		s := str.([]byte)
		buf = s
	}
	m.fd.Write(buf)
	return m
}

func (m *fileIOStruct) read(num ...int) string {
	m.lock.acquire()
	defer m.lock.release()

	var bytes []byte
	var err error
	if len(num) == 0 {
		bytes, err = ioutil.ReadAll(m.fd)
		panicerr(err)
		m.close()
	} else {
		buffer := make([]byte, num[0])
		i, err := io.ReadAtLeast(m.fd, buffer, num[0])
		if err != nil {
			if !String("EOF").In(err.Error()) {
				panicerr(err)
			}
		}
		bytes = buffer[:i]
	}

	return string(bytes)
}

func (m *fileIOStruct) seek(num int64) {
	_, err := m.fd.Seek(num, 0)
	panicerr(err)
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
	panicerr(err)
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
