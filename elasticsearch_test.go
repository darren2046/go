package golanglibs

import (
	"testing"
)

func TestElasticsearch(t *testing.T) {
	es := getElasticsearch("http://192.168.168.18:9200")
	co := es.Collection("test")
	co.Index(1, map[string]interface{}{
		"id":    1,
		"title": "Chinese copywriting for Chinese audience",
	})
	co.Search("title", "Chinese", 1, 1, ElasticsearchSearchingConfigStruct{
		OrderByKey:   "id",
		OrderByOrder: "desc",
		Highlight:    "title",
	})
}
