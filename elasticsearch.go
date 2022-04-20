package golanglibs

import (
	"encoding/json"
)

type ElasticsearchStruct struct {
	baseurl string
}

func getElasticsearch(baseurl string) *ElasticsearchStruct {
	return &ElasticsearchStruct{baseurl: baseurl}
}

func (m *ElasticsearchStruct) Delete(IndexName string) {
	r := httpDelete(String(m.baseurl).Strip("/").S + "/" + IndexName)
	Lg.Debug(r)
}

type ElasticsearchCollectionStruct struct {
	baseurl string
}

func (m *ElasticsearchStruct) Collection(name string) *ElasticsearchCollectionStruct {
	return &ElasticsearchCollectionStruct{baseurl: String(m.baseurl).Strip("/").S + "/" + name}
}

// id是唯一的字符串
func (m *ElasticsearchCollectionStruct) Index(id interface{}, data map[string]interface{}, refresh ...bool) {
	var url string
	if len(refresh) != 0 && refresh[0] {
		url = m.baseurl + "/_doc/" + Str(id) + "?refresh"
	} else {
		url = m.baseurl + "/_doc/" + Str(id)
	}
	r := httpPutJSON(url, data, HttpConfig{TimeoutRetryTimes: -1})
	if r.StatusCode != 201 && r.StatusCode != 200 {
		Lg.Debug(r)
		Panicerr("插入到Elasticsearch出错: 状态码不是201或者200")
	}
}

func (m *ElasticsearchCollectionStruct) Refresh() {
	r := httpPost(m.baseurl+"/_refresh", HttpConfig{TimeoutRetryTimes: -1})
	Lg.Debug(r)
}

type ElasticsearchSearchingConfigStruct struct {
	OrderByKey   string // 根据document里面的这个key排序
	OrderByOrder string // 可选desc或者asc, 默认asc
	Highlight    string // 这个字段里面的value，如果有搜索到关键字则插入html标签以高亮。
	Fuzzy        bool   // 是否模糊搜索，如果为false则所有的词都必须要出现才行
}

type ElasticsearchSearchResultStruct struct {
	Total int
	Data  map[int]map[string]interface{}
}

type ElasticsearchSearchedResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		//MaxScore interface{} `json:"max_score"`
		Hits []struct {
			ID string `json:"_id"`
			// Score  interface{}            `json:"_score"`
			Source map[string]interface{} `json:"_source"`
			// Sort   []int                  `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

func (m *ElasticsearchCollectionStruct) Search(key string, value string, page int, pagesize int, cfg ...ElasticsearchSearchingConfigStruct) *ElasticsearchSearchedResult {
	if page*pagesize > 10000 {
		Panicerr("偏移量不能超过10000: page * pagesize = " + Str(page*pagesize))
	}

	startfrom := (page - 1) * pagesize

	var query map[string]interface{}

	if len(cfg) != 0 {
		if !cfg[0].Fuzzy {
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"match_phrase": map[string]interface{}{
						key: map[string]interface{}{
							"query": value,
							"slop":  100,
						},
					},
				},
				"from": startfrom,
				"size": pagesize,
			}
		} else {
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						key: value,
					},
				},
				"from": startfrom,
				"size": pagesize,
			}
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_phrase": map[string]interface{}{
					key: map[string]interface{}{
						"query": value,
						"slop":  100,
					},
				},
			},
			"from": startfrom,
			"size": pagesize,
		}
	}

	if len(cfg) != 0 {
		if cfg[0].OrderByKey != "" {
			var order string
			if cfg[0].OrderByOrder != "" {
				order = cfg[0].OrderByOrder
			} else {
				order = "asc"
			}
			query["sort"] = []map[string]interface{}{
				{
					cfg[0].OrderByKey: map[string]interface{}{
						"order": order,
					},
				},
			}
		}

		if cfg[0].Highlight != "" {
			query["highlight"] = map[string]interface{}{
				"fields": map[string]interface{}{
					cfg[0].Highlight: map[string]interface{}{},
				},
			}
		}
	}

	Lg.Trace(Json.Dumps(query))
	// Lg.Trace(m.baseurl)

	r := Http.PostJSON(m.baseurl+"/_search", query)
	// Lg.Debug(r)
	if r.StatusCode != 200 {
		Panicerr("在Elasticsearch中搜寻出错：" + r.Content.S)
	}

	res := ElasticsearchSearchedResult{}
	err := json.Unmarshal([]byte(r.Content.S), &res)
	Panicerr(err)

	// Lg.Debug(res)

	return &res
}

func (m *ElasticsearchCollectionStruct) Delete(id interface{}) {
	r := Http.Delete(m.baseurl + "/_doc/" + Str(id))
	if r.StatusCode != 200 {
		Lg.Debug(r)
		Panicerr("在elasticsearch删除id为\"" + Str(id) + "\"的文档出错")
	}
}
