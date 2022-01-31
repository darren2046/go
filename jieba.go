package golanglibs

import "github.com/yanyiwu/gojieba"

type JiebaStruct struct {
	jieba *gojieba.Jieba
}

func getJieba() *JiebaStruct {
	return &JiebaStruct{
		jieba: gojieba.NewJieba(),
	}
}

func (m *JiebaStruct) Close() {
	m.jieba.Free()
}

func (m *JiebaStruct) Cut(s string) []string {
	return m.jieba.Cut(s, true)
}
