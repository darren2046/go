package golanglibs

import (
	"encoding/json"
)

type prometheusClientStruct struct {
	url string
}

type prometheusClientOriginalResultStruct struct {
	Status string
	Data   struct {
		ResultType string
		Result     []struct {
			Metric map[string]interface{}
			Value  []interface{}
		}
	}
}

type prometheusClientResultStruct struct {
	Label map[string]string
	Value float64
}

func getPrometheusClient(url string) *prometheusClientStruct {
	return &prometheusClientStruct{url: url + "/api/v1/query"}
}

func (m *prometheusClientStruct) Query(query string, time ...float64) (res []prometheusClientResultStruct) {
	var ttime int64
	if len(time) == 0 {
		ttime = Int64(Time.Now())
	} else {
		ttime = Int64(time[0])
	}
	resp := httpPost(m.url, HttpParam{"query": query, "time": ttime})
	// fmt.Println(resp.content)
	if resp.StatusCode != 200 {
		Panicerr("查询prometheus出错, 查询语句: " + query + ", 状态码: " + Str(resp.StatusCode) + ", 错误消息:" + resp.Content.S)
	}
	// fmt.Println(resp.content)

	var por prometheusClientOriginalResultStruct
	err := json.Unmarshal([]byte(resp.Content.S), &por)
	Panicerr(err)
	// lg.debug(por)
	if por.Status != "success" {
		Panicerr("查询prometheus出错, 查询语句: " + query + ", prometheus查询结果状态: " + Str(por.Status))
	}

	var pr prometheusClientResultStruct
	for _, i := range por.Data.Result {
		pr.Label = make(map[string]string)
		for k, v := range i.Metric {
			pr.Label[k] = Str(v)
		}
		pr.Value = Float64(i.Value[1])
		res = append(res, pr)
	}
	return
}
