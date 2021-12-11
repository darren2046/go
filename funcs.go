package golanglibs

import (
	"context"
	"crypto/x509"
	"net"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/h2non/filetype"
	"github.com/hpcloud/tail"
	"github.com/icrowley/fake"
	"github.com/miekg/dns"
	ua "github.com/mileusna/useragent"
	"github.com/mmcdole/gofeed"
)

type funcsStruct struct {
	Nslookup               func(name string, querytype string, dnsService ...string) [][]string
	FakeName               func() string
	FileType               func(fpath string) string
	Inotify                func(path string) chan *fsnotifyFileEventStruct
	IPLocation             func(ip string, dbpath ...string) *ipLocationInfo
	Gethostbyname          func(hostname string, dnsserver ...string) (res []string)
	Getcname               func(hostname string, dnsserver ...string) (res string)
	HightLightHTMLForCode  func(code string, codeType ...string) (html string)
	Markdown2html          func(md string) string
	CPUUsagePerProgress    func() (res map[int64]progressCPUUsageStruct)
	ResizeImg              func(srcPath string, dstPath string, width int, height ...int)
	GetRSS                 func(url string, config ...rssConfig) *gofeed.Feed
	GbkToUtf8              func(s string) string
	Utf8ToGbk              func(s string) string
	GetSnowflakeID         func(nodeNumber ...int) int64
	GetRemoteServerSSLCert func(host string, port ...int) []*x509.Certificate
	Tailf                  func(path string, startFromEndOfFile ...bool) chan *tail.Line
	BaiduTranslateAnyToZH  func(text string) string
	ParseUserAgent         func(UserAgent string) ua.UserAgent
	Wget                   func(url string, cfg ...WgetCfg) (filename string)
	Whois                  func(s string, servers ...string) string
}

var Funcs funcsStruct

func init() {
	Funcs = funcsStruct{
		Nslookup:               nslookup,
		FakeName:               fakeName,
		FileType:               fileType,
		Inotify:                inotify,
		IPLocation:             getIPLocation,
		Gethostbyname:          gethostbyname,
		Getcname:               getcname,
		HightLightHTMLForCode:  getHightLightHTML,
		Markdown2html:          md2html,
		CPUUsagePerProgress:    getSystemProgressCPUUsage,
		ResizeImg:              resizeImg,
		GetRSS:                 getRSS,
		GbkToUtf8:              gbkToUtf8,
		Utf8ToGbk:              utf8ToGbk,
		GetSnowflakeID:         getSnowflakeID,
		GetRemoteServerSSLCert: getRemoteServerCert,
		Tailf:                  tailf,
		BaiduTranslateAnyToZH:  baiduTranslateAnyToZH,
		ParseUserAgent:         parseUserAgent,
		Wget:                   wget,
		Whois:                  whois,
	}
}

type WgetCfg struct {
	savepath string // 保存的本地路径, 可以是目录或者完整文件路径
	retry    int    // 出错尝试次数, -1为一直尝试直到成功
}

func wget(url string, cfg ...WgetCfg) (filename string) {
	savepath := "."
	retry := 0
	if len(cfg) != 0 {
		if cfg[0].savepath != "" {
			savepath = cfg[0].savepath
		}
		retry = cfg[0].retry
	}
	if retry < 0 {
		for {
			resp, err := grab.Get(savepath, url)
			if err == nil {
				return resp.Filename
			}
		}
	} else {
		i := 0
		for {
			resp, err := grab.Get(savepath, url)
			if err == nil {
				return resp.Filename
			}
			i++

			if i > retry && err != nil {
				panicerr(err)
			}
		}
	}
}

func parseUserAgent(UserAgent string) ua.UserAgent {
	return ua.Parse(UserAgent)
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
	kind, err := filetype.Match([]byte(Open(fpath).read()))

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
