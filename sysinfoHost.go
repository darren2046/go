package golanglibs

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

type sysinfoHostStruct struct {
	Info     func() types.HostInfo
	Memory   func() *types.HostMemoryInfo
	CPUTimes func() types.CPUTimes
}

var sysinfohoststruct sysinfoHostStruct

func init() {
	sysinfohoststruct = sysinfoHostStruct{
		Info:     getSysinfoHostInfo,
		Memory:   getSysinfoHostMemory,
		CPUTimes: getSysinfoHostCPUTimer,
	}
}

func getSysinfoHostInfo() types.HostInfo {
	process, err := sysinfo.Host()
	Panicerr(err)

	return process.Info()
}

func getSysinfoHostMemory() *types.HostMemoryInfo {
	process, err := sysinfo.Host()
	Panicerr(err)

	info, err := process.Memory()
	Panicerr(err)
	return info
}

func getSysinfoHostCPUTimer() types.CPUTimes {
	process, err := sysinfo.Host()
	Panicerr(err)

	info, err := process.CPUTime()
	Panicerr(err)

	return info
}
