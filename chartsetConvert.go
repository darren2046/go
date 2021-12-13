package golanglibs

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func gbkToUtf8(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	Panicerr(err)
	return string(d)
}

func utf8ToGbk(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, err := ioutil.ReadAll(reader)
	Panicerr(err)
	return string(d)
}
