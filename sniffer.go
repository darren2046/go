//go:build pcap

package golanglibs

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type networkPacketStruct struct {
	data  string
	sport int
	dport int
	proto string // tcp, udp
	ipv   int    // 4, 6
	sip   string
	dip   string
	smac  string
	dmac  string
}

func init() {
	Funcs.Sniffer = sniffer
	Funcs.ReadPcapFile = readPcapFile
}

func doPacketSource(packetSource *gopacket.PacketSource, pkgchan chan *networkPacketStruct, pcapFileHandler ...*pcap.Handle) {
	for packet := range packetSource.Packets() {
		//print("Packet found: ", packet)
		transportLayer := packet.TransportLayer()
		if transportLayer != nil {
			//print("transportLayer found.")
			pkg := networkPacketStruct{}

			linuxSLLLayer := packet.Layer(layers.LayerTypeLinuxSLL)
			if linuxSLLLayer != nil {
				linuxSLLPacket, _ := linuxSLLLayer.(*layers.LinuxSLL)
				pkg.smac = fmt.Sprintf("%s", linuxSLLPacket.Addr)
			}

			//print(packet)
			ethLayer := packet.Layer(layers.LayerTypeEthernet)
			if ethLayer != nil {
				//print("eth layer found")
				ethernetPacket, _ := ethLayer.(*layers.Ethernet)
				pkg.smac = fmt.Sprintf("%s", ethernetPacket.SrcMAC)
				pkg.dmac = fmt.Sprintf("%s", ethernetPacket.DstMAC)
				//fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer != nil {
				//print("ip layer found")
				ip, ok := ipLayer.(*layers.IPv4)
				if ok {
					pkg.ipv = 4
					pkg.sip = fmt.Sprintf("%s", ip.SrcIP)
					pkg.dip = fmt.Sprintf("%s", ip.DstIP)
				} else {
					pkg.ipv = 6
					ip6, _ := ipLayer.(*layers.IPv6)
					pkg.sip = fmt.Sprintf("%s", ip6.SrcIP)
					pkg.dip = fmt.Sprintf("%s", ip6.DstIP)
				}
			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				//print("tcp layer found")
				pkg.proto = "tcp"
				tcp, _ := tcpLayer.(*layers.TCP)
				pkg.sport = Int(fmt.Sprintf("%d", tcp.SrcPort))
				pkg.dport = Int(fmt.Sprintf("%d", tcp.DstPort))
			}

			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer != nil {
				//print("udp layer found")
				pkg.proto = "udp"
				udp, _ := udpLayer.(*layers.UDP)
				pkg.sport = Int(fmt.Sprintf("%d", udp.SrcPort))
				pkg.dport = Int(fmt.Sprintf("%d", udp.DstPort))
			}

			applicationLayer := packet.TransportLayer()
			if applicationLayer != nil {
				pkg.data = Str(applicationLayer.LayerPayload())
				//print("Data:", pkg.data)
			}

			// if strStartsWith(pkg.data, "GET /action") {
			// 	print(packet)
			// }

			if pkg.data != "" {
				pkgchan <- &pkg
			}
		}
	}
	if len(pcapFileHandler) != 0 {
		pcapFileHandler[0].Close()
	}
	close(pkgchan)
}

func sniffer(interfaceName string, filterString string, promisc ...bool) chan *networkPacketStruct {
	// 4096是读取每一个包的buffer, mtu一般为1500, 所以4096是超出了很多, 除非mtu超出了4096, 才读不全
	// promisc为设置网卡为混杂模式
	// timeout为0.3秒, 是kernel每0.3秒就会吐一次数据给pcap, 如果这个为30秒, 则收到数据包之后会继续等待其他数据包, 30秒再一起吐出来
	var handle *pcap.Handle
	var err error
	if len(promisc) == 0 {
		handle, err = pcap.OpenLive(interfaceName, 4096, false, getTimeDuration(0.3))
	} else {
		handle, err = pcap.OpenLive(interfaceName, 4096, promisc[0], getTimeDuration(0.3))
	}

	Panicerr(err)
	//defer handle.Close()

	err = handle.SetBPFFilter(filterString)
	Panicerr(err)

	pkgchan := make(chan *networkPacketStruct)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go doPacketSource(packetSource, pkgchan)

	return pkgchan
}

func readPcapFile(pcapFile string) chan *networkPacketStruct {
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}

	pkgchan := make(chan *networkPacketStruct)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go doPacketSource(packetSource, pkgchan, handle)

	return pkgchan
}
