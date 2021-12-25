package golanglibs

type socketStruct struct {
	KCP  *kcpStruct
	Smux *smuxStruct
	SSL  *sslStruct
	TCP  *tcpStruct
	UDP  *udpStruct
	UNIX *unixSocketStruct
}

var Socket socketStruct

func init() {
	Socket = socketStruct{
		KCP:  &kcpstruct,
		Smux: &smuxstruct,
		SSL:  &sslstruct,
		TCP:  &tcpstruct,
		UDP:  &udpstruct,
		UNIX: &unixsocketstruct,
	}
}
