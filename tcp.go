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
	conn    net.Conn
	isclose bool
}

type tcpServerSideListener struct {
	listener net.Listener
	isclose  bool
}

func tcpListen(host string, port int) *tcpServerSideListener {
	l, err := net.Listen("tcp", host+":"+Str(port))
	panicerr(err)

	return &tcpServerSideListener{listener: l}
}

func (m *tcpServerSideListener) accept() chan *tcpServerSideConn {
	ch := make(chan *tcpServerSideConn)

	go func() {
		for {
			c, err := m.listener.Accept()

			if err != nil {
				if String("use of closed network connection").In(err.Error()) {
					break
				}
				panicerr(err)
			}
			ct := &tcpServerSideConn{conn: c, isclose: false}
			ch <- ct
		}
	}()

	return ch
}

func (m *tcpServerSideListener) close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}
}

func (m *tcpServerSideConn) close() {
	if !m.isclose {
		m.isclose = true
		m.conn.Close()
	}
}

func (m *tcpServerSideConn) send(str string) {
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *tcpServerSideConn) recv(buffersize int) string {
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

// TCP - Client

type tcpClientSideConn struct {
	conn    *net.TCPConn
	isclose bool
}

func tcpConnect(host string, port int, timeout ...int) *tcpClientSideConn {
	servAddr := host + ":" + Str(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	panicerr(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	panicerr(err)
	return &tcpClientSideConn{conn: conn, isclose: false}
}

func (m *tcpClientSideConn) send(str string, timeout ...int) {
	if len(timeout) != 0 {
		m.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *tcpClientSideConn) recv(buffersize int, timeout ...int) string {
	if len(timeout) != 0 {
		m.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout[0]) * time.Second))
	}
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

func (m *tcpClientSideConn) close() {
	if !m.isclose {
		m.isclose = true
		err := m.conn.Close()
		panicerr(err)
	}
}
