package golanglibs

import (
	"strings"

	"github.com/antchfx/jsonquery"
)

type xpathJsonStruct struct {
	doc *jsonquery.Node
}

func getXPathJson(jsonstr string) *xpathJsonStruct {
	doc, err := jsonquery.Parse(strings.NewReader(jsonstr))
	Panicerr(err)
	return &xpathJsonStruct{
		doc: doc,
	}
}

func (m *xpathJsonStruct) Exists(expr string) bool {
	return Try(func() {
		m.First(expr).Text()
	}).Error == nil
}

func (m *xpathJsonStruct) First(expr string) (res *xpathJsonStruct) {
	return &xpathJsonStruct{
		doc: jsonquery.FindOne(m.doc, expr),
	}
}

func (m *xpathJsonStruct) Find(expr string) (res []*xpathJsonStruct) {
	for _, doc := range jsonquery.Find(m.doc, expr) {
		res = append(res, &xpathJsonStruct{doc: doc})
	}
	return
}

func (m *xpathJsonStruct) Text() string {
	return m.doc.InnerText()
}
