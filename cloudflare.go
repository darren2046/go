package golanglibs

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

type cloudflareRecord struct {
	ID                       string
	Data                     string
	Name                     string
	TTL                      int
	Type                     string
	Proxiable                bool
	Proxied                  bool
	OriginalCloudflareRecord cloudflare.DNSRecord
}

type cloudflareStruct struct {
	api   *cloudflare.API
	key   string
	email string
}

type cloudflareDomainInfoStruct struct {
	CreatedAt   time.Time
	Domain      string
	DomainID    string
	NameServers []string
	Status      string
}

type cloudflareDomainStruct struct {
	api        *cloudflare.API
	DomainID   string
	DomainName string
}

func getCloudflare(key string, email string) *cloudflareStruct {
	api, err := cloudflare.New(key, email)
	Panicerr(err)
	return &cloudflareStruct{
		api:   api,
		key:   key,
		email: email,
	}
}

func (m *cloudflareStruct) Add(domain string) cloudflare.Zone {
	zone, err := m.api.CreateZone(context.Background(), domain, false, cloudflare.Account{}, "full")
	Panicerr(err)
	return zone
}

func (m *cloudflareStruct) List() (res []cloudflareDomainInfoStruct) {
	zones, err := m.api.ListZones(context.Background())
	Panicerr(err)
	for _, zone := range zones {
		res = append(res, cloudflareDomainInfoStruct{
			CreatedAt:   zone.CreatedOn,
			Domain:      zone.Name,
			DomainID:    zone.ID,
			NameServers: zone.NameServers,
			Status:      zone.Status,
		})
	}
	return
}

func (m *cloudflareStruct) Domain(domainName string) *cloudflareDomainStruct {
	id, err := m.api.ZoneIDByName(domainName)
	Panicerr(err)
	return &cloudflareDomainStruct{
		api:        m.api,
		DomainID:   id,
		DomainName: domainName,
	}
}

func (m *cloudflareDomainStruct) List() (res []cloudflareRecord) {
	records, err := m.api.DNSRecords(context.Background(), m.DomainID, cloudflare.DNSRecord{})
	Panicerr(err)

	var name string
	for _, record := range records {
		if record.Name == record.ZoneName {
			name = "@"
		} else {
			name = record.Name[:len(record.Name)-len(record.ZoneName)-1]
		}
		res = append(res, cloudflareRecord{
			Data:                     record.Content,
			Name:                     name,
			TTL:                      record.TTL,
			Type:                     String(record.Type).Lower().Get(),
			Proxiable:                record.Proxiable,
			Proxied:                  *record.Proxied,
			ID:                       record.ID,
			OriginalCloudflareRecord: record,
		})
	}

	return
}

func (m *cloudflareDomainStruct) Delete(name string) {
	// 虽然之后有log的代码, 但是在这个函数里面关闭log
	// logLevel := lg.getLevel()
	// lg.setLevel("")
	//defer lg.setLevel(logLevel)

	// if name == "@" {
	// 	name = m.DomainName
	// }

	// lg.debug(name)
	for _, v := range m.List() {
		// lg.debug(v)
		if name == v.Name {
			err := m.api.DeleteDNSRecord(context.Background(), m.DomainID, v.ID)
			Panicerr(err)
		}
	}
}

func (m *cloudflareDomainStruct) Add(recordName string, recordType string, recordValue string, proxied ...bool) *cloudflare.DNSRecordResponse {
	if recordName == "@" {
		recordName = m.DomainName
	} else {
		recordName = recordName + "." + m.DomainName
	}
	var prox bool
	if len(proxied) == 0 {
		prox = false
	} else {
		prox = proxied[0]
	}
	var p uint16 = 10
	resp, err := m.api.CreateDNSRecord(context.Background(), m.DomainID, cloudflare.DNSRecord{
		Type:     String(recordType).Upper().Get(),
		Name:     recordName,
		Content:  recordValue,
		TTL:      300,
		Priority: &p,
		Proxied:  &prox,
	})
	Panicerr(err)
	return resp
}

func (m *cloudflareDomainStruct) SetProxied(subdomain string, proxied bool) {
	for _, v := range m.List() {
		//lg.trace(v.Name, domain, proxied)
		if v.Name == subdomain {
			if !v.Proxiable && proxied {
				panic(newerr("类型为" + v.Type + "的dns记录无法设置cdn代理加速"))
			} else if v.Proxied != proxied {
				v.OriginalCloudflareRecord.Proxied = &proxied
				m.api.UpdateDNSRecord(context.Background(), m.DomainID, v.ID, v.OriginalCloudflareRecord)
			}
		}
	}
}

func (m *cloudflareDomainStruct) Update(recordName string, recordType string, recordValue string, proxied ...bool) {
	var prox bool
	if len(proxied) == 0 {
		prox = false
	} else {
		prox = proxied[0]
	}

	for _, v := range m.List() {
		//lg.trace(v.Name, recordName)
		if v.Name == recordName {
			m.Delete(recordName)
			if len(proxied) == 0 {
				prox = v.Proxied
			}
		}
	}

	m.Add(recordName, recordType, recordValue, prox)
}
