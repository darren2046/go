package golanglibs

import (
	"encoding/csv"
	"os"
)

type csvStruct struct {
	Reader func(fpath string) *csvReaderStruct
	Writer func(fpath string, mode string) *csvWriterStruct
}

var csvstruct *csvStruct

func init() {
	csvstruct = &csvStruct{
		Reader: getCSVReader,
		Writer: getCSVWriter,
	}
}

type csvReaderStruct struct {
	reader  *csv.Reader
	headers []string
	isclose bool
	fd      *os.File
}

func getCSVReader(fpath string) *csvReaderStruct {
	fd := Open(fpath).fd

	reader := csv.NewReader(fd)
	headers, err := reader.Read()
	if err != nil {
		Panicerr("Error while reading headers:" + Str(err))
	}
	return &csvReaderStruct{
		reader:  reader,
		headers: headers,
		isclose: false,
		fd:      fd,
	}
}

func (m *csvReaderStruct) Read() (res map[string]string) {
	if m.isclose {
		Panicerr("文件已關閉")
	}

	row, err := m.reader.Read()
	if err != nil {
		m.isclose = true
		m.fd.Close()
		Panicerr(err)
	}

	if len(row) != len(m.headers) {
		m.isclose = true
		m.fd.Close()
		Panicerr("數據的個數跟表頭的個數不同")
	}

	res = make(map[string]string)
	for i := range row {
		res[m.headers[i]] = row[i]
	}

	return
}

func (m *csvReaderStruct) Readrows() chan map[string]string {
	if m.isclose {
		Panicerr("文件已關閉")
	}

	reschan := make(chan map[string]string)

	go func() {
		if err := Try(func() {
			for {
				reschan <- m.Read()
			}
		}).Error; err != nil {
			close(reschan)
		}
	}()

	return reschan
}

func (m *csvReaderStruct) Close() {
	m.fd.Close()
	m.isclose = true
}

type csvWriterStruct struct {
	writer  *csv.Writer
	headers []string
	fd      *os.File
	isclose bool
	mode    string
}

// mode可以是a或者w，跟打開文件一樣
// 如果是w的話會需要優先設置header才能wirte數據
func getCSVWriter(fpath string, mode string) *csvWriterStruct {
	headers := []string{}
	if mode == "a" {
		fd := Open(fpath).fd

		reader := csv.NewReader(fd)
		var err error
		headers, err = reader.Read()
		if err != nil {
			Panicerr("Error while reading headers:" + Str(err))
		}
	}

	fd := Open(fpath, mode).fd
	writer := csv.NewWriter(fd)
	return &csvWriterStruct{
		writer:  writer,
		fd:      fd,
		isclose: false,
		headers: headers,
		mode:    mode,
	}
}

func (m *csvWriterStruct) Flush() {
	if m.isclose {
		Panicerr("文件已關閉")
	}

	m.writer.Flush()
}

func (m *csvWriterStruct) SetHeaders(headers []string) {
	if m.isclose {
		Panicerr("文件已關閉")
	}

	if m.mode == "w" && (len(m.headers) == 0 || len(m.headers) != len(headers)) {
		m.writer.Write(headers)
		m.writer.Flush()
	}
	m.headers = headers
}

func (m *csvWriterStruct) Write(record map[string]string) {
	if m.isclose {
		Panicerr("文件已關閉")
	}

	if len(m.headers) == 0 {
		Panicerr("需要先設置表頭，請使用SetHeader()方法設置")
	}
	row := []string{}
	for _, field := range m.headers {
		if Map(record).Has(field) {
			row = append(row, record[field])
		} else {
			row = append(row, "")
		}
	}
	err := m.writer.Write(row)
	m.writer.Flush()
	Panicerr(err)
}

func (m *csvWriterStruct) Close() {
	m.isclose = true
	m.writer.Flush()
	m.fd.Close()
}
