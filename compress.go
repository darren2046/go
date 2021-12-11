package golanglibs

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"

	"github.com/ulikunitz/xz"
)

type compressStruct struct {
	LzmaCompressString   func(text string) string
	LzmaDecompressString func(text string) string
	ZlibCompressString   func(text string) string
	ZlibDecompressString func(text string) string
}

var compressstruct compressStruct

func init() {
	compressstruct = compressStruct{
		LzmaCompressString:   lzmaCompressString,
		LzmaDecompressString: lzmaDecompressString,
		ZlibCompressString:   zlibCompressString,
		ZlibDecompressString: zlibDecompressString,
	}
}
func lzmaCompressString(text string) string {
	var buf bytes.Buffer

	defer buf.Reset()

	w, err := xz.NewWriter(&buf)
	panicerr(err)

	_, err = io.WriteString(w, text)
	panicerr(err)

	err = w.Close()
	panicerr(err)

	return buf.String()
}

func lzmaDecompressString(text string) string {
	var buf bytes.Buffer
	buf.Write([]byte(text))

	defer buf.Reset()

	r, err := xz.NewReader(&buf)
	panicerr(err)

	dbuf := new(strings.Builder)
	_, err = io.Copy(dbuf, r)
	panicerr(err)
	return dbuf.String()
}

func zlibCompressString(text string) string {
	var buf bytes.Buffer

	w := zlib.NewWriter(&buf)
	_, err := w.Write([]byte(text))
	panicerr(err)

	err = w.Close()
	panicerr(err)

	return buf.String()
}

func zlibDecompressString(text string) string {
	var buf bytes.Buffer
	buf.Write([]byte(text))

	defer buf.Reset()

	r, err := zlib.NewReader(&buf)
	panicerr(err)

	dbuf := new(strings.Builder)
	_, err = io.Copy(dbuf, r)
	panicerr(err)

	return dbuf.String()
}
