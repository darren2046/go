package golanglibs

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
	"github.com/h2non/filetype"
	"github.com/hpcloud/tail"
	"github.com/miekg/dns"
	ua "github.com/mileusna/useragent"
	"github.com/mmcdole/gofeed"
)

type funcsStruct struct {
	Nslookup               func(name string, querytype string, dnsService ...string) [][]string
	FakeNameEnglish        func() string
	FakeNameChinese        func() string
	FileType               func(fpath string) string
	Inotify                func(path string) chan *fsnotifyFileEventStruct
	IPLocation             func(ip string, dbpath ...string) *ipLocationInfo
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
	IpInNet                func(ip string, Net string, mask ...string) bool
	Int2ip                 func(ipnr int64) string
	Ip2int                 func(ipnr string) int64
	Zh2PinYin              func(zh string) (ress []string)
	Fmtsize                func(num uint64) string
	Sniffer                func(interfaceName string, filterString string, promisc ...bool) chan *networkPacketStruct
	ReadPcapFile           func(pcapFile string) chan *networkPacketStruct
}

var Funcs funcsStruct

func init() {
	Funcs = funcsStruct{
		Nslookup:               nslookup,
		FakeNameEnglish:        fakeNameEnglish,
		FakeNameChinese:        fakeNameChinese,
		FileType:               fileType,
		Inotify:                inotify,
		IPLocation:             getIPLocation,
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
		// BaiduTranslateAnyToZH:  baiduTranslateAnyToZH,
		ParseUserAgent: parseUserAgent,
		Wget:           wget,
		Whois:          whois,
		IpInNet:        ipInNet,
		Int2ip:         int2ip,
		Ip2int:         ip2int,
		Zh2PinYin:      zh2PinYin,
		Fmtsize:        fmtsize,
	}
}

func fmtsize(num uint64) string {
	return humanize.Bytes(num)
}

func getRemoteServerCert(host string, port ...int) []*x509.Certificate {
	var p string
	if len(port) == 0 {
		p = "443"
	} else {
		p = Str(port[0])
	}

	conn, err := tls.Dial("tcp", host+":"+Str(p), nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}

	return conn.ConnectionState().PeerCertificates
}

func int2ip(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

func ip2int(ipnr string) int64 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func ipInNet(ip string, Net string, mask ...string) bool {
	if len(mask) != 0 {
		ip := net.ParseIP(mask[0])
		addr := ip.To4()
		cidrsuffix, _ := net.IPv4Mask(addr[0], addr[1], addr[2], addr[3]).Size()
		Net = Net + "/" + Str(cidrsuffix)
	}

	_, ipnetA, _ := net.ParseCIDR(Net)
	ipB := net.ParseIP(ip)

	if ipnetA.Contains(ipB) {
		return true
	} else {
		return false
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
				Panicerr(err)
			}
		}
	}
}

func parseUserAgent(UserAgent string) ua.UserAgent {
	return ua.Parse(UserAgent)
}

func fileType(fpath string) string {
	kind, err := filetype.Match([]byte(Open(fpath).Read().S))

	Panicerr(err)

	if kind == filetype.Unknown {
		return ""
	}
	return kind.Extension
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

	Panicerr(err)

	return dst
}
