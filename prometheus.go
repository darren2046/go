package golanglibs

import (
	"encoding/json"
)

type prometheusStruct struct {
	url string
}

type prometheusOriginalResultStruct struct {
	Status string
	Data   struct {
		ResultType string
		Result     []struct {
			Metric map[string]interface{}
			Value  []interface{}
		}
	}
}

type prometheusResultStruct struct {
	Label map[string]string
	Value float64
}

func getPrometheus(url string) *prometheusStruct {
	return &prometheusStruct{url: url + "/api/v1/query"}
}

func (m *prometheusStruct) query(query string, time ...float64) (res []prometheusResultStruct) {
	var ttime int64
	if len(time) == 0 {
		ttime = Int64(Time.Now())
	} else {
		ttime = Int64(time[0])
	}
	resp := httpPost(m.url, HttpParam{"query": query, "time": ttime})
	// fmt.Println(resp.content)
	if resp.statusCode != 200 {
		panicerr("查询Prometheus出错, 查询语句: " + query + ", 状态码: " + Str(resp.statusCode) + ", 错误消息:" + resp.content)
	}
	// fmt.Println(resp.content)

	var por prometheusOriginalResultStruct
	err := json.Unmarshal([]byte(resp.content), &por)
	panicerr(err)
	// lg.debug(por)
	if por.Status != "success" {
		panicerr("查询Prometheus出错, 查询语句: " + query + ", Prometheus查询结果状态: " + Str(por.Status))
	}

	var pr prometheusResultStruct
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
