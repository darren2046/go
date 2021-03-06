package golanglibs

import (
	"io/ioutil"

	"github.com/rakyll/statik/fs"
	"github.com/wangtuanjie/ip17mon"
)

var ip17modHadInit bool
var ip17DBFileModifyTime int64

type ipLocationInfo struct {
	Country string
	Region  string
	City    string
	Isp     string
}

// wget https://cdn.jsdelivr.net/npm/qqwry.ipdb/qqwry.ipdb -O files/qqwry.ipdb

func getIPLocation(ip string, dbpath ...string) *ipLocationInfo {
	if len(dbpath) == 0 {
		if !ip17modHadInit {
			statikFS, err := fs.New()
			Panicerr(err)

			ipdbfd, err := statikFS.Open("/qqwry.ipdb")
			Panicerr(err)
			ipdbBytes, err := ioutil.ReadAll(ipdbfd)
			Panicerr(err)
			ipdbfd.Close()

			ip17mon.InitWithIpdb(ipdbBytes)

			ip17modHadInit = true
		}
	} else {
		if ip17DBFileModifyTime == 0 || File(dbpath[0]).Time().Mtime != ip17DBFileModifyTime {
			ip17DBFileModifyTime = File(dbpath[0]).Time().Mtime

			ip17mon.InitWithIpdb([]byte(Open(dbpath[0]).Read().S))
		}
	}

	loc, err := ip17mon.Find(ip)
	Panicerr(err)

	return &ipLocationInfo{
		City:    loc.City,
		Region:  loc.Region,
		Country: loc.Country,
		Isp:     loc.Isp,
	}
}
