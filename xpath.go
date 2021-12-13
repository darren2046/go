package golanglibs

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type xpathStruct struct {
	doc *html.Node
}

func getXPath(htmlString string) *xpathStruct {
	doc, err := htmlquery.Parse(strings.NewReader(htmlString))
	Panicerr(err)
	return &xpathStruct{
		doc: doc,
	}
}

func (m *xpathStruct) First(expr string) (res *xpathStruct) {
	return &xpathStruct{
		doc: htmlquery.FindOne(m.doc, expr),
	}
}

func (m *xpathStruct) Find(expr string) (res []*xpathStruct) {
	for _, doc := range htmlquery.Find(m.doc, expr) {
		res = append(res, &xpathStruct{doc: doc})
	}
	return
}

func (m *xpathStruct) Text() string {
	return htmlquery.InnerText(m.doc)
}

func (m *xpathStruct) GetAttr(attr string) string {
	return htmlquery.SelectAttr(m.doc, attr)
}

func (m *xpathStruct) Html() string {
	return htmlquery.OutputHTML(m.doc, true)
}
