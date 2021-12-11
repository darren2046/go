package golanglibs

import (
	"crypto/tls"
	"crypto/x509"
	"net"
)

type sslStruct struct {
	Listen        func(host string, port int, key string, crt string) *tcpServerSideListener
	ServerWrapper func(conn net.Conn, key string, crt string) *tcpServerSideConn
	Connect       func(host string, port int, cfg ...sslCfg) *sslClientSideConn
	ClientWrapper func(conn net.Conn, cfg ...sslCfg) *sslClientSideConn
}

var sslstruct sslStruct

func init() {
	sslstruct = sslStruct{
		Listen:        sslListen,
		ServerWrapper: sslServerWrapper,
		Connect:       sslConnect,
		ClientWrapper: sslClientWrapper,
	}
}

// SSL - Server
// 只实现了一个Listener， 其他的方法是tcp的方法

func sslListen(host string, port int, key string, crt string) *tcpServerSideListener {
	cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
	panicerr(err)

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp4", host+":"+Str(port), tlsCfg)
	panicerr(err)

	return &tcpServerSideListener{listener: listener}
}

func sslServerWrapper(conn net.Conn, key string, crt string) *tcpServerSideConn {
	cert, err := tls.X509KeyPair([]byte(crt), []byte(key))
	panicerr(err)

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	tconn := tls.Server(conn, tlsCfg)

	return &tcpServerSideConn{conn: tconn}
}

// SSL - Client

type sslClientSideConn struct {
	conn *tls.Conn
}

type sslCfg struct {
	InsecureSkipVerify  bool     // true为跳过证书验证
	additionRootCA      []string // 额外的用来验证证书的CA证书
	domain              string   // 需要认证的域名, 也会在请求证书的时候提供
	withoutSystemRootCA bool     // true为不使用系统内置的CA
}

func sslConnect(host string, port int, cfg ...sslCfg) *sslClientSideConn {
	servAddr := host + ":" + Str(port)

	tcfg := tls.Config{}
	if len(cfg) != 0 {
		if cfg[0].InsecureSkipVerify {
			tcfg.InsecureSkipVerify = cfg[0].InsecureSkipVerify
		}
		if len(cfg[0].additionRootCA) != 0 {
			var rootCAs *x509.CertPool
			if cfg[0].withoutSystemRootCA {
				rootCAs = x509.NewCertPool()
			} else {
				rootCAs, _ = x509.SystemCertPool()
				if rootCAs == nil {
					rootCAs = x509.NewCertPool()
				}
			}
			for _, ca := range cfg[0].additionRootCA {
				rootCAs.AppendCertsFromPEM([]byte(ca))
			}
			tcfg.RootCAs = rootCAs
		}
		if cfg[0].domain != "" {
			tcfg.ServerName = cfg[0].domain
		}
	}

	conn, err := tls.Dial("tcp", servAddr, &tcfg)
	panicerr(err)
	return &sslClientSideConn{conn: conn}
}

func sslClientWrapper(conn net.Conn, cfg ...sslCfg) *sslClientSideConn {
	tcfg := tls.Config{}
	if len(cfg) != 0 {
		if cfg[0].InsecureSkipVerify {
			tcfg.InsecureSkipVerify = cfg[0].InsecureSkipVerify
		}
		if len(cfg[0].additionRootCA) != 0 {
			var rootCAs *x509.CertPool
			if cfg[0].withoutSystemRootCA {
				rootCAs = x509.NewCertPool()
			} else {
				rootCAs, _ = x509.SystemCertPool()
				if rootCAs == nil {
					rootCAs = x509.NewCertPool()
				}
			}
			for _, ca := range cfg[0].additionRootCA {
				rootCAs.AppendCertsFromPEM([]byte(ca))
			}
			tcfg.RootCAs = rootCAs
		}
		if cfg[0].domain != "" {
			tcfg.ServerName = cfg[0].domain
		}
	}

	tconn := tls.Client(conn, &tcfg)
	return &sslClientSideConn{conn: tconn}
}

func (m *sslClientSideConn) send(str string) {
	_, err := m.conn.Write([]byte(str))
	panicerr(err)
}

func (m *sslClientSideConn) recv(buffersize int) string {
	reply := make([]byte, buffersize)
	n, err := m.conn.Read(reply)
	panicerr(err)
	return string(reply[:n])
}

func (m *sslClientSideConn) close() {
	err := m.conn.Close()
	panicerr(err)
}
