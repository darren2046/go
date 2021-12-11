package golanglibs

import (
	"crypto/tls"
	"crypto/x509"
)

type sslStruct struct {
	GetCert func(host string, port ...int) []*x509.Certificate
}

var Ssl sslStruct

func init() {
	Ssl = sslStruct{
		GetCert: getRemoteServerCert,
	}
}

func getRemoteServerCert(host string, port ...int) []*x509.Certificate {
	var p string
	if len(port) == 0 {
		p = "443"
	} else {
		p = Str(port[0])
	}

	conn, err := tls.Dial("tcp", host+":"+Str(p), nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}

	return conn.ConnectionState().PeerCertificates
}
