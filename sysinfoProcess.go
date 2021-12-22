package golanglibs

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

type sysinfoProcessStruct struct {
	pid int
}

func getSysinfoProcess(pid int) *sysinfoProcessStruct {
	return &sysinfoProcessStruct{pid: pid}
}

func (p *sysinfoProcessStruct) Info() types.ProcessInfo {
	process, err := sysinfo.Process(p.pid)
	Panicerr(err)

	info, err := process.Info()
	Panicerr(err)
	return info
}

func (p *sysinfoProcessStruct) Memory() types.MemoryInfo {
	process, err := sysinfo.Process(p.pid)
	Panicerr(err)

	info, err := process.Memory()
	Panicerr(err)
	return info
}

func (p *sysinfoProcessStruct) User() types.UserInfo {
	process, err := sysinfo.Process(p.pid)
	Panicerr(err)

	info, err := process.User()
	Panicerr(err)
	return info
}

func (p *sysinfoProcessStruct) Parent() *sysinfoProcessStruct {
	process, err := sysinfo.Process(p.pid)
	Panicerr(err)

	pp, err := process.Parent()
	Panicerr(err)

	return getSysinfoProcess(pp.PID())
}

func (p *sysinfoProcessStruct) CPUTimes() types.CPUTimes {
	process, err := sysinfo.Process(p.pid)
	Panicerr(err)

	info, err := process.CPUTime()
	Panicerr(err)
	return info
}
