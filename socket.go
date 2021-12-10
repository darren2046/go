package golanglibs

import (
	"context"
	"net"
	"time"
)

type socketStruct struct {
	Gethostbyname func(hostname string, dnsserver ...string) (res []string)
	Getcname      func(hostname string, dnsserver ...string) (res string)
}

var Socket socketStruct

func init() {
	Socket = socketStruct{
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
