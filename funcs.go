package golanglibs

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/icrowley/fake"
	"github.com/miekg/dns"
)

type funcsStruct struct {
	Nslookup      func(name string, querytype string, dnsService ...string) [][]string
	FakeName      func() string
	FileType      func(fpath string) string
	Inotify       func(path string) chan *fsnotifyFileEventStruct
	IPLocation    func(ip string, dbpath ...string) *ipLocationInfo
	Gethostbyname func(hostname string, dnsserver ...string) (res []string)
	Getcname      func(hostname string, dnsserver ...string) (res string)
}

var Funcs funcsStruct

func init() {
	Funcs = funcsStruct{
		Nslookup:      nslookup,
		FakeName:      fakeName,
		FileType:      fileType,
		Inotify:       inotify,
		IPLocation:    getIPLocation,
		Gethostbyname: gethostbyname,
		Getcname:      getcname,
	}
}

func getcname(hostname string, dnsserver ...string) (res string) {
	hostname = String(hostname).RStrip(".").Get()
	var err error
	if len(dnsserver) == 0 {
		res, err = net.LookupCNAME(hostname)
		panicerr(err)
		if String(res).RStrip(".").Get() == hostname {
			res = ""
		}
		return
	} else {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(8000),
				}
				if !String(":").In(dnsserver[0]) {
					dnsserver[0] = dnsserver[0] + ":53"
				}
				return d.DialContext(ctx, "udp", dnsserver[0])
			},
		}
		res, err = r.LookupCNAME(context.Background(), hostname)
		panicerr(err)
		if String(hostname).RStrip(".").Get() == hostname {
			res = ""
		}
		return
	}
}

func gethostbyname(hostname string, dnsserver ...string) (res []string) {
	if len(dnsserver) == 0 {
		ips, err := net.LookupIP(hostname)
		panicerr(err)
		if ips != nil {
			for _, v := range ips {
				if v.To4() != nil || v.To16() != nil {
					res = append(res, v.String())
				}
			}
		}
		return
	} else {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(8000),
				}
				if !String(":").In(dnsserver[0]) {
					dnsserver[0] = dnsserver[0] + ":53"
				}
				return d.DialContext(ctx, "udp", dnsserver[0])
			},
		}
		ips, err := r.LookupHost(context.Background(), hostname)
		panicerr(err)
		for _, ip := range ips {
			if net.ParseIP(ip) != nil {
				res = append(res, ip)
			}
		}
	}
	return
}

func fileType(fpath string) string {
	kind, err := filetype.Match([]byte(open(fpath).read()))

	panicerr(err)

	if kind == filetype.Unknown {
		return ""
	}
	return kind.Extension
}

func fakeName() string {
	return fake.FullName()
}

func nslookup(name string, querytype string, dnsService ...string) [][]string {
	c := dns.Client{
		Timeout: 5 * time.Second,
	}

	var server string
	if len(dnsService) == 0 {
		server = "8.8.8.8"
	} else {
		server = dnsService[0]
	}

	querytype = String(querytype).Lower().Get()
	var qtype uint16
	if querytype == "ns" {
		qtype = dns.TypeNS
	} else if querytype == "a" {
		qtype = dns.TypeA
	} else if querytype == "txt" {
		qtype = dns.TypeTXT
	} else if querytype == "cname" {
		qtype = dns.TypeCNAME
	} else if querytype == "aaaa" {
		qtype = dns.TypeAAAA
	} else if querytype == "soa" {
		qtype = dns.TypeSOA
	} else if querytype == "mx" {
		qtype = dns.TypeMX
	}

	var err error
	var dst [][]string
	for i := 0; i < 3; i++ {
		m := dns.Msg{}
		if !String(name).EndsWith(".") {
			name = name + "."
		}
		m.SetQuestion(name, qtype)
		r, _, err := c.Exchange(&m, server+":53")
		if err != nil {
			time.Sleep(1 * time.Second * time.Duration(i+1))
			continue
		}

		for _, ans := range r.Answer {
			s := strings.Split(ans.String(), "\t")
			dst = append(dst, []string{s[0], s[3], s[4]})
		}
		err = nil
		break
	}

	panicerr(err)

	return dst
}
