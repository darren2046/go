package golanglibs

import (
	"testing"
)

func TestElasticsearch(t *testing.T) {
	es := getElasticsearch("http://192.168.168.18:9200")
	co := es.Collection("telegram_history_content")
	// co.Index(1, map[string]interface{}{
	// 	"id":    1,
	// 	"title": "Chinese copywriting for Chinese audience",
	// })
	co.Search("text", `fhs 下浮\{出U\}  @peii_usdt`, 1000, 10, ElasticsearchSearchingConfigStruct{Fuzzy: true})
	// co.Delete(1)
}
