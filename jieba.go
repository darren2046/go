//go:build jieba

package golanglibs

import "github.com/yanyiwu/gojieba"

type JiebaSubStruct struct {
	jieba *gojieba.Jieba
}

func init() {
	Tools.Jieba = getJieba
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
