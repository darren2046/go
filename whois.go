package golanglibs

import (
	wwwhoisgo "github.com/likexian/whois"
)

func whois(s string, servers ...string) string {
	result, err := wwwhoisgo.Whois(s, servers...)
	Panicerr(err)
	return result
}

// var domainSuffixWhoisServerMap map[string]string

// func getDomainWhoisInfo(domain string) string {
// 	darr := String(domain).Split(".")
// 	suffix := darr[len(darr)-1]
// 	var dws string

// 	if domainSuffixWhoisServerMap == nil {
// 		domainSuffixWhoisServerMap = make(map[string]string)
// 	}

// 	if Map(domainSuffixWhoisServerMap).Has(suffix) {
// 		dws = domainSuffixWhoisServerMap[suffix]
// 	} else {
// 		//lg.trace("尝试查找顶级域" + suffix + "的whois服务器")
// 		s := Re.FindAll("whois:(.+)", whois(suffix, "whois.iana.org"))
// 		if len(s) == 0 {
// 			Panicerr("找不到顶级域" + suffix + "的whois服务器")
// 		}
// 		dws = String(s[0][1]).Strip().Get()

// 		domainSuffixWhoisServerMap[suffix] = dws
// 	}

// 	return whois(domain, dws)
// }
