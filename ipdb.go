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

func getIPLocation(ip string, dbpath ...string) *ipLocationInfo {
	if len(dbpath) == 0 {
		if ip17modHadInit == false {
			statikFS, err := fs.New()
			panicerr(err)

			ipdbfd, err := statikFS.Open("/qqwry.ipdb")
			panicerr(err)
			ipdbBytes, err := ioutil.ReadAll(ipdbfd)
			panicerr(err)
			ipdbfd.Close()

			ip17mon.InitWithIpdb(ipdbBytes)

			ip17modHadInit = true
		}
	} else {
		if ip17DBFileModifyTime == 0 || File(dbpath[0]).Time().mtime != ip17DBFileModifyTime {
			ip17DBFileModifyTime = File(dbpath[0]).Time().mtime

			ip17mon.InitWithIpdb([]byte(open(dbpath[0]).read()))
		}
	}

	loc, err := ip17mon.Find(ip)
	panicerr(err)

	return &ipLocationInfo{
		City:    loc.City,
		Region:  loc.Region,
		Country: loc.Country,
		Isp:     loc.Isp,
	}
}
