package golanglibs

type socketStruct struct {
	Kcp kcpStruct
}

var Socket socketStruct

func init() {
	Socket = socketStruct{
		Kcp: kcpstruct,
	}
}
