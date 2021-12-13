package golanglibs

import (
	"encoding/json"
	"time"
)

type godaddyRecord struct {
	Data string `json:"data"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
	Type string `json:"type"`
}

type godaddyStruct struct {
	header HttpHeader
}

type godaddyDomainStruct struct {
	header     HttpHeader
	domainName string
}

type godaddyDomainInfoStruct struct {
	CreatedAt              time.Time   `json:"createdAt"`
	Domain                 string      `json:"domain"`
	DomainID               int         `json:"domainId"`
	ExpirationProtected    bool        `json:"expirationProtected"`
	Expires                time.Time   `json:"expires"`
	ExposeWhois            bool        `json:"exposeWhois"`
	HoldRegistrar          bool        `json:"holdRegistrar"`
	Locked                 bool        `json:"locked"`
	NameServers            interface{} `json:"nameServers"`
	Privacy                bool        `json:"privacy"`
	RenewAuto              bool        `json:"renewAuto"`
	RenewDeadline          time.Time   `json:"renewDeadline"`
	Renewable              bool        `json:"renewable"`
	Status                 string      `json:"status"`
	TransferAwayEligibleAt time.Time   `json:"transferAwayEligibleAt"`
	TransferProtected      bool        `json:"transferProtected"`
}

func getGodaddy(key string, secret string) *godaddyStruct {
	return &godaddyStruct{
		header: HttpHeader{
			"Authorization": "sso-key " + key + ":" + secret,
			"Content-Type":  "application/json",
			"Accept":        "application/json",
		},
	}
}

func (m *godaddyStruct) List() (res []godaddyDomainInfoStruct) {
	resp := httpGet("https://api.godaddy.com/v1/domains", m.header)
	if resp.statusCode != 200 {
		Panicerr(resp.content)
	}
	err := json.Unmarshal([]byte(resp.content), &res)
	Panicerr(err)
	return
}

func (m *godaddyStruct) Domain(domainName string) *godaddyDomainStruct {
	return &godaddyDomainStruct{
		header:     m.header,
		domainName: domainName,
	}
}

func (m *godaddyDomainStruct) List() (res []godaddyRecord) {
	resp := httpGet("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", m.header)
	if resp.statusCode != 200 {
		Panicerr(resp.content)
	}
	err := json.Unmarshal([]byte(resp.content), &res)
	Panicerr(err)
	return
}

// 参数留空字符串则忽略这个项目
func (m *godaddyDomainStruct) Delete(name string, dtype string, value string) {
	dtype = String(dtype).Upper().Get()
	var records []godaddyRecord
	for _, v := range m.List() {
		if v.Name != name && !(v.Data == "Parked" && v.Type == "A") {
			if name != "" && dtype != "" && value != "" {
				if !(v.Name == name && v.Type == dtype && v.Data == value) {
					records = append(records, v)
				}
			} else if name == "" && dtype != "" && value != "" {
				if v.Type != dtype && v.Data != value {
					records = append(records, v)
				}
			} else if name != "" && dtype == "" && value != "" {
				if v.Name != name && v.Data != value {
					records = append(records, v)
				}
			} else if name != "" && dtype != "" && value == "" {
				if v.Name == name && v.Type != dtype {
					records = append(records, v)
				}
			} else if name != "" && dtype == "" && value == "" {
				if v.Name != name {
					records = append(records, v)
				}
			} else if name == "" && dtype == "" && value != "" {
				if v.Data != value {
					records = append(records, v)
				}
			} else if name == "" && dtype != "" && value == "" {
				if v.Type != dtype {
					records = append(records, v)
				}
			}
		}
	}

	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		Panicerr(resp.content)
	}
}

func (m *godaddyDomainStruct) Modify(recordName string, srcRecordType string, srcRecordValue string, dstRecordType string, dstRecordValue string) {
	var records []godaddyRecord
	for _, v := range m.List() {
		if !(v.Data == "Parked" && v.Type == "A") {
			records = append(records, v)
		}
	}

	for idx := range records {
		if records[idx].Name == recordName {
			if records[idx].Type == srcRecordType && records[idx].Data == srcRecordValue {
				records[idx].Type = dstRecordType
				records[idx].Data = dstRecordValue
			}
		}
	}
	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		Panicerr(resp.content)
	}
}

func (m *godaddyDomainStruct) Add(recordName string, recordType string, recordValue string) {
	var records []godaddyRecord
	for _, v := range m.List() {
		if !(v.Data == "Parked" && v.Type == "A") {
			records = append(records, v)
		}
	}

	records = append(records, godaddyRecord{
		Data: recordValue,
		Name: recordName,
		TTL:  600,
		Type: String(recordType).Upper().Get(),
	})

	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
	if resp.statusCode != 200 {
		Panicerr(resp.content)
	}
}

// 更新域名的某个记录, 如果不存在则新增
// 如果同一个名字, 多个A记录这样, 应该会出问题
// func (m *godaddyDomainStruct) update(recordName string, recordType string, recordValue string) {
// 	recordType = strUpper(recordType)
// 	var records []godaddyRecord
// 	for _, v := range m.list() {
// 		if !(v.Data == "Parked" && v.Type == "A") {
// 			if !(v.Name == recordName && v.Type == recordType) {
// 				records = append(records, v)
// 			}
// 		}
// 	}

// 	records = append(records, godaddyRecord{
// 		Data: recordValue,
// 		Name: recordName,
// 		TTL:  600,
// 		Type: recordType,
// 	})

// 	resp := httpPutJSON("https://api.godaddy.com/v1/domains/"+m.domainName+"/records", records, m.header)
// 	if resp.statusCode != 200 {
// 		_, fn, line, _ := runtime.Caller(0)
// 		panic(filepath.Base(fn) + ":" + strconv.Itoa(line+1) + " >> " + resp.content + " >> " + fmtDebugStack(string(debug.Stack())))
// 	}
// }
