package golanglibs

type sysinfoStruct struct {
	Host    *sysinfoHostStruct
	Process func(pid int) *sysinfoProcessStruct
}

var sysinfostruct sysinfoStruct

func init() {
	sysinfostruct = sysinfoStruct{
		Host:    &sysinfohoststruct,
		Process: getSysinfoProcess,
	}
}
