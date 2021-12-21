package golanglibs

import (
	"github.com/elastic/go-sysinfo/types"
)

type sysinfoStruct struct {
	HostInfo func() types.HostInfo
}

var sysinfostruct sysinfoStruct

func init() {
	sysinfostruct = sysinfoStruct{
		HostInfo: getSysinfoHostInfo,
	}
}
