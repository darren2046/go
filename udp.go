package golanglibs

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type udpStruct struct {
	Listen  func(host string, port int) udpServerSideConn
	Connect func(host string, port int) udpClientSideConn
}

var udpstruct udpStruct

func init() {
	udpstruct = udpStruct{
		Listen:  udpListen,
		Connect: udpConnect,
	}
}

// UDP - Client

type udpClientSideConn struct {
	isclose bool
	conn    net.Conn
}

func udpConnect(host string, port int) udpClientSideConn {
	conn, err := net.Dial("udp", host+":"+Str(port))
	panicerr(err)
	return udpClientSideConn{conn: conn}
}

func (m *udpClientSideConn) send(str string) {
	_, err := fmt.Fprintf(m.conn, str)
	panicerr(err)
}

func (m *udpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}

func (m *udpClientSideConn) recv(buffersize int) string {
	p := make([]byte, buffersize)
	n, err := bufio.NewReader(m.conn).Read(p)
	panicerr(err)
	return string(p[:n])
}

// UDP - Server

type udpServerSideConn struct {
	isclose bool
	conn    *net.UDPConn
}

func udpListen(host string, port int) udpServerSideConn {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	ser, err := net.ListenUDP("udp", &addr)
	panicerr(err)
	return udpServerSideConn{conn: ser}
}

func (m *udpServerSideConn) recvfrom(buffersize int, timeout ...int) (string, *net.UDPAddr) {
	if len(timeout) != 0 {
		m.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	p := make([]byte, buffersize)
	n, remoteaddr, err := m.conn.ReadFromUDP(p)
	panicerr(err)
	return string(p[:n]), remoteaddr
}

func (m *udpServerSideConn) sendto(data string, address *net.UDPAddr, timeout ...int) {
	if len(timeout) != 0 {
		m.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	_, err := m.conn.WriteToUDP([]byte(data), address)
	panicerr(err)
}

func (m *udpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}
