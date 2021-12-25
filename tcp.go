package golanglibs

import (
	"net"
	"time"
)

type tcpStruct struct {
	Listen  func(host string, port int) *tcpServerSideListener
	Connect func(host string, port int, timeout ...int) *tcpClientSideConn
}

var tcpstruct tcpStruct

func init() {
	tcpstruct = tcpStruct{
		Listen:  tcpListen,
		Connect: tcpConnect,
	}
}

// TCP - Server

type tcpServerSideConn struct {
	Conn    net.Conn
	isclose bool
}

type tcpServerSideListener struct {
	listener net.Listener
	isclose  bool
}

func tcpListen(host string, port int) *tcpServerSideListener {
	l, err := net.Listen("tcp", host+":"+Str(port))
	Panicerr(err)

	return &tcpServerSideListener{listener: l}
}

func (m *tcpServerSideListener) Accept() chan *tcpServerSideConn {
	ch := make(chan *tcpServerSideConn)

	go func() {
		for {
			c, err := m.listener.Accept()

			if err != nil {
				if String("use of closed network connection").In(err.Error()) {
					break
				}
				Panicerr(err)
			}
			ct := &tcpServerSideConn{Conn: c, isclose: false}
			ch <- ct
		}
	}()

	return ch
}

func (m *tcpServerSideListener) Close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}
}

func (m *tcpServerSideConn) Close() {
	if !m.isclose {
		m.isclose = true
		m.Conn.Close()
	}
}

func (m *tcpServerSideConn) Send(str string) {
	_, err := m.Conn.Write([]byte(str))
	Panicerr(err)
}

func (m *tcpServerSideConn) Recv(buffersize int) string {
	reply := make([]byte, buffersize)
	n, err := m.Conn.Read(reply)
	Panicerr(err)
	return string(reply[:n])
}

// TCP - Client

type tcpClientSideConn struct {
	Conn    *net.TCPConn
	isclose bool
}

func tcpConnect(host string, port int, timeout ...int) *tcpClientSideConn {
	servAddr := host + ":" + Str(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	Panicerr(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	Panicerr(err)
	return &tcpClientSideConn{Conn: conn, isclose: false}
}

func (m *tcpClientSideConn) Send(str string, timeout ...int) {
	if len(timeout) != 0 {
		m.Conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	_, err := m.Conn.Write([]byte(str))
	Panicerr(err)
}

func (m *tcpClientSideConn) Recv(buffersize int, timeout ...int) string {
	if len(timeout) != 0 {
		m.Conn.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	reply := make([]byte, buffersize)
	n, err := m.Conn.Read(reply)
	Panicerr(err)
	return string(reply[:n])
}

func (m *tcpClientSideConn) Close() {
	if !m.isclose {
		m.isclose = true
		err := m.Conn.Close()
		Panicerr(err)
	}
}
