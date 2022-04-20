package golanglibs

import (
	"testing"
)

func TestElasticsearch(t *testing.T) {
	es := getElasticsearch("http://192.168.168.18:9200")

	idx := "test"
	es.Delete(idx)
	// Time.Sleep(10)

	co := es.Collection(idx)

	co.Index(1, map[string]interface{}{
		"id":   1,
		"text": `test string`,
	})

	co.Refresh()
	// Time.Sleep(10)

	eres := co.Search("text", `test string`, 1, 1, ElasticsearchSearchingConfigStruct{Fuzzy: true})
	Lg.Debug(eres)
	// co.Delete(1)
}
