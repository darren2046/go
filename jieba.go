package golanglibs

import (
	"github.com/go-ego/gse"
)

type JiebaStruct struct {
	seg gse.Segmenter
}

func getJieba() *JiebaStruct {
	var seg gse.Segmenter
	seg.LoadDictEmbed()
	seg.LoadDictEmbed("zh_s")
	seg.LoadDictEmbed("zh_t")
	seg.LoadDictEmbed("jp")

	return &JiebaStruct{
		seg: seg,
	}
}

func (m *JiebaStruct) LoadDict(path ...string) {
	m.seg.LoadDict(path...)
}

func (m *JiebaStruct) Cut(s string) []string {
	return m.seg.Cut(s, true)
}

func (m *JiebaStruct) AddWord(text string) {
	m.seg.AddToken(text, 100)
}
