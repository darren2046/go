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
	panicerr(err)
	return &xpathJsonStruct{
		doc: doc,
	}
}

func (m *xpathJsonStruct) exists(expr string) bool {
	return try(func() {
		m.first(expr).text()
	}).Error == nil
}

func (m *xpathJsonStruct) first(expr string) (res *xpathJsonStruct) {
	return &xpathJsonStruct{
		doc: jsonquery.FindOne(m.doc, expr),
	}
}

func (m *xpathJsonStruct) find(expr string) (res []*xpathJsonStruct) {
	for _, doc := range jsonquery.Find(m.doc, expr) {
		res = append(res, &xpathJsonStruct{doc: doc})
	}
	return
}

func (m *xpathJsonStruct) text() string {
	return m.doc.InnerText()
}
