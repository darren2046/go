package golanglibs

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
)

type hashStruct struct {
	Md5sum   func(str string) string
	Md5File  func(path string) string
	Sha1sum  func(str string) string
	Sha1File func(path string) string
}

var Hash hashStruct

func init() {
	Hash = hashStruct{
		Md5sum:   md5sum,
		Md5File:  md5File,
		Sha1sum:  sha1sum,
		Sha1File: sha1File,
	}
}

func sha1File(path string) string {
	data, err := ioutil.ReadFile(path)
	panicerr(err)
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func sha1sum(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func md5sum(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func md5File(path string) string {
	data, err := ioutil.ReadFile(path)
	panicerr(err)
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
