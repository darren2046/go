package golanglibs

import (
	"net"
)

type unixSocketStruct struct {
	Listen  func(path string) *unixSocketServerSideListener
	Connect func(path string) *tcpServerSideConn
}

var unixsocketstruct unixSocketStruct

func init() {
	unixsocketstruct = unixSocketStruct{
		Listen:  unixSocketListen,
		Connect: unixSocketConnect,
	}
}

type unixSocketServerSideListener struct {
	listener net.Listener
	isclose  bool
}

func unixSocketListen(path string) *unixSocketServerSideListener {
	l, err := net.Listen("unix", path)
	Panicerr(err)

	return &unixSocketServerSideListener{listener: l}
}

func (m *unixSocketServerSideListener) Close() {
	if !m.isclose {
		m.isclose = true
		m.listener.Close()
	}
}

func (m *unixSocketServerSideListener) Accept() chan *tcpServerSideConn {
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

// --- ----

func unixSocketConnect(path string) *tcpServerSideConn {
	conn, err := net.Dial("unix", path)
	Panicerr(err)
	return &tcpServerSideConn{Conn: conn, isclose: false}
}
