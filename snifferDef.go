package golanglibs

type networkPacketStruct struct {
	Data      string
	SrcPort   int
	DstPort   int
	Protocol  string // tcp, udp
	IPVersion int    // 4, 6
	SrcIP     string
	DstIP     string
	SrcMac    string
	DstMac    string
}
